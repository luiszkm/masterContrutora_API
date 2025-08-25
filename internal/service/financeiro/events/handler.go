package events

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/service/financeiro/dto"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
)

// ContaReceberService interface para o service de contas a receber
type ContaReceberService interface {
	CriarConta(ctx context.Context, input dto.CriarContaReceberInput) (*dto.ContaReceberOutput, error)
	BuscarPorCronogramaRecebimentoID(ctx context.Context, cronogramaID string) (*dto.ContaReceberOutput, error)
	SincronizarRecebimento(ctx context.Context, contaID string, valor float64, observacoes string) (*dto.ContaReceberOutput, error)
}

// ContaPagarService interface para o service de contas a pagar
type ContaPagarService interface {
	CriarConta(ctx context.Context, input dto.CriarContaPagarInput) (*dto.ContaPagarOutput, error)
	CriarContaDeOrcamento(ctx context.Context, input dto.CriarContaPagarDeOrcamentoInput, orcamento interface{}) (*dto.ContaPagarOutput, error)
	CancelarContaDeOrcamento(ctx context.Context, orcamentoID string) error
	CancelarContaDeApontamento(ctx context.Context, apontamentoID string) error
}

// FinanceiroEventHandler processa eventos relacionados ao módulo financeiro
type FinanceiroEventHandler struct {
	contaReceberService ContaReceberService
	contaPagarService   ContaPagarService
	logger              *slog.Logger
}

func NovoFinanceiroEventHandler(
	contaReceberService ContaReceberService,
	contaPagarService ContaPagarService,
	logger *slog.Logger,
) *FinanceiroEventHandler {
	return &FinanceiroEventHandler{
		contaReceberService: contaReceberService,
		contaPagarService:   contaPagarService,
		logger:              logger.With("handler", "FinanceiroEventHandler"),
	}
}

// HandleCronogramaRecebimentoCriado processa evento de cronograma criado
func (h *FinanceiroEventHandler) HandleCronogramaRecebimentoCriado(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.CronogramaRecebimentoCriadoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de cronograma inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando criação de cronograma", 
		"obra_id", payload.ObraID, 
		"quantidade_etapas", payload.QuantidadeEtapas,
		"valor_total", payload.ValorTotalPrevisto)

	// Para cada cronograma criado, criar uma conta a receber correspondente
	for i, cronogramaID := range payload.CronogramasIds {
		// Calcular valor proporcional (assumindo divisão igual por etapa)
		valorEtapa := payload.ValorTotalPrevisto / float64(payload.QuantidadeEtapas)
		
		input := dto.CriarContaReceberInput{
			ObraID:                  &payload.ObraID,
			CronogramaRecebimentoID: &cronogramaID,
			Cliente:                 payload.Cliente,
			TipoContaReceber:        "OBRA",
			Descricao:               payload.ObraNome + " - Etapa " + string(rune('1'+i)),
			ValorOriginal:           valorEtapa,
			DataVencimento:          payload.PrimeiroVencimento,
		}

		conta, err := h.contaReceberService.CriarConta(ctx, input)
		if err != nil {
			h.logger.ErrorContext(ctx, "falha ao criar conta a receber a partir do cronograma", 
				"cronograma_id", cronogramaID, 
				"obra_id", payload.ObraID,
				"erro", err)
			continue
		}

		h.logger.InfoContext(ctx, "conta a receber criada a partir do cronograma", 
			"conta_id", conta.ID,
			"cronograma_id", cronogramaID,
			"valor", valorEtapa)
	}
}

// HandleRecebimentoRealizado processa evento de recebimento realizado
func (h *FinanceiroEventHandler) HandleRecebimentoRealizado(ctx context.Context, evento bus.Evento) {
	h.logger.InfoContext(ctx, "🔧 NOVA VERSÃO: HandleRecebimentoRealizado executando")
	payload, ok := evento.Payload.(events.RecebimentoRealizadoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de recebimento inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando recebimento realizado", 
		"cliente", payload.Cliente,
		"valor", payload.ValorRecebido,
		"cronograma_id", payload.CronogramaRecebimentoID)

	// Se não há cronograma ID, não podemos sincronizar
	if payload.CronogramaRecebimentoID == nil {
		h.logger.InfoContext(ctx, "recebimento sem cronograma vinculado, nada para sincronizar")
		return
	}

	h.logger.InfoContext(ctx, "DEBUG: cronograma ID encontrado, iniciando sincronização", 
		"cronograma_id", *payload.CronogramaRecebimentoID)

	cronogramaID := *payload.CronogramaRecebimentoID

	// Buscar conta a receber por cronograma ID
	conta, err := h.contaReceberService.BuscarPorCronogramaRecebimentoID(ctx, cronogramaID)
	if err != nil {
		h.logger.ErrorContext(ctx, "falha ao buscar conta a receber para sincronização", 
			"cronograma_id", cronogramaID,
			"erro", err)
		return
	}

	// Sincronizar recebimento na conta a receber (sem disparar eventos para evitar loop)
	observacoes := fmt.Sprintf("SYNC:Cronograma:%s", cronogramaID)
	contaAtualizada, err := h.contaReceberService.SincronizarRecebimento(ctx, conta.ID, payload.ValorRecebido, observacoes)
	if err != nil {
		h.logger.ErrorContext(ctx, "falha ao sincronizar recebimento na conta a receber", 
			"conta_id", conta.ID,
			"cronograma_id", cronogramaID,
			"valor", payload.ValorRecebido,
			"erro", err)
		return
	}

	h.logger.InfoContext(ctx, "conta a receber sincronizada com sucesso", 
		"conta_id", conta.ID,
		"cronograma_id", cronogramaID,
		"valor", payload.ValorRecebido,
		"novo_status", contaAtualizada.Status,
		"valor_total_recebido", contaAtualizada.ValorRecebido)

	// TODO: Implementar criação de MovimentacaoFinanceira quando a entidade existir
}

// HandleOrcamentoStatusAtualizado processa quando orçamento é aprovado (cria conta a pagar)
func (h *FinanceiroEventHandler) HandleOrcamentoStatusAtualizado(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.OrcamentoStatusAtualizadoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de orçamento inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando atualização de status de orçamento", 
		"orcamento_id", payload.OrcamentoID,
		"status_anterior", payload.StatusAnterior,
		"novo_status", payload.NovoStatus,
		"valor", payload.Valor)

	// Quando orçamento é aprovado, criar conta a pagar automaticamente
	if payload.NovoStatus == "Aprovado" {
		h.logger.InfoContext(ctx, "orçamento aprovado - criando conta a pagar", 
			"orcamento_id", payload.OrcamentoID,
			"valor", payload.Valor)
		
		// Criar input para conta a pagar baseado no orçamento
		input := dto.CriarContaPagarDeOrcamentoInput{
			OrcamentoID:        payload.OrcamentoID,
			ValorOrcamento:     payload.Valor, // ✅ Usando valor real do orçamento
			DataVencimento:     time.Now().AddDate(0, 0, 30), // 30 dias para vencimento por padrão
			NumeroDocumento:    nil, // Será preenchido quando tiver a nota fiscal
			Observacoes:        func() *string { s := "Conta gerada automaticamente do orçamento aprovado"; return &s }(),
			DividirParcelas:    false, // Por padrão, não dividir em parcelas
		}

		conta, err := h.contaPagarService.CriarContaDeOrcamento(ctx, input, nil)
		if err != nil {
			h.logger.ErrorContext(ctx, "falha ao criar conta a pagar a partir do orçamento", 
				"orcamento_id", payload.OrcamentoID,
				"erro", err)
			return
		}

		h.logger.InfoContext(ctx, "conta a pagar criada automaticamente", 
			"conta_id", conta.ID,
			"orcamento_id", payload.OrcamentoID,
			"valor", conta.ValorOriginal)
	}

	// Quando orçamento muda de "Aprovado" para "Cancelado", cancelar conta a pagar
	if payload.StatusAnterior == "Aprovado" && payload.NovoStatus == "Cancelado" {
		h.logger.InfoContext(ctx, "orçamento cancelado após aprovação - cancelando conta a pagar", 
			"orcamento_id", payload.OrcamentoID)
		
		if err := h.contaPagarService.CancelarContaDeOrcamento(ctx, payload.OrcamentoID); err != nil {
			h.logger.ErrorContext(ctx, "falha ao cancelar conta a pagar do orçamento cancelado", 
				"orcamento_id", payload.OrcamentoID,
				"erro", err)
			return
		}

		h.logger.InfoContext(ctx, "conta a pagar cancelada automaticamente devido ao cancelamento do orçamento", 
			"orcamento_id", payload.OrcamentoID)
	}
}


// HandleApontamentoAprovado cria conta a pagar quando apontamento é aprovado
func (h *FinanceiroEventHandler) HandleApontamentoAprovado(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.ApontamentoAprovadoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de apontamento aprovado inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando apontamento aprovado", 
		"apontamento_id", payload.ApontamentoID,
		"funcionario", payload.FuncionarioNome,
		"valor", payload.ValorCalculado)

	// Criar conta a pagar para o funcionário
	input := dto.CriarContaPagarInput{
		FornecedorNome:  payload.FuncionarioNome,
		TipoContaPagar:  "FUNCIONARIO",
		Descricao:       fmt.Sprintf("Pagamento de funcionário - %s (%s)", payload.FuncionarioNome, payload.PeriodoReferencia),
		ValorOriginal:   payload.ValorCalculado,
		DataVencimento:  payload.DataVencimentoPrevisto,
		ObraID:          &payload.ObraID,
		NumeroDocumento: &payload.ApontamentoID, // Referência ao apontamento
		Observacoes:     func() *string { s := fmt.Sprintf("Conta gerada automaticamente do apontamento aprovado %s", payload.ApontamentoID); return &s }(),
	}

	conta, err := h.contaPagarService.CriarConta(ctx, input)
	if err != nil {
		h.logger.ErrorContext(ctx, "falha ao criar conta a pagar para apontamento", 
			"apontamento_id", payload.ApontamentoID,
			"funcionario_id", payload.FuncionarioID,
			"erro", err)
		return
	}

	h.logger.InfoContext(ctx, "conta a pagar criada para apontamento aprovado", 
		"conta_id", conta.ID,
		"apontamento_id", payload.ApontamentoID,
		"funcionario", payload.FuncionarioNome,
		"valor", payload.ValorCalculado)
}

// HandleApontamentoCancelado cancela conta a pagar quando apontamento é cancelado
func (h *FinanceiroEventHandler) HandleApontamentoCancelado(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.ApontamentoCanceladoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de apontamento cancelado inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando apontamento cancelado", 
		"apontamento_id", payload.ApontamentoID,
		"funcionario", payload.FuncionarioNome,
		"status_anterior", payload.StatusAnterior,
		"motivo", payload.MotivoCancelamento)

	// Cancelar conta a pagar associada ao apontamento
	if err := h.contaPagarService.CancelarContaDeApontamento(ctx, payload.ApontamentoID); err != nil {
		h.logger.ErrorContext(ctx, "falha ao cancelar conta a pagar do apontamento cancelado", 
			"apontamento_id", payload.ApontamentoID,
			"funcionario_id", payload.FuncionarioID,
			"erro", err)
		return
	}

	h.logger.InfoContext(ctx, "conta a pagar cancelada automaticamente devido ao cancelamento do apontamento", 
		"apontamento_id", payload.ApontamentoID,
		"funcionario", payload.FuncionarioNome)
}

// HandleOrcamentoExcluido cancela conta a pagar quando orçamento é excluído
func (h *FinanceiroEventHandler) HandleOrcamentoExcluido(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.OrcamentoExcluidoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de orçamento excluído inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando orçamento excluído", 
		"orcamento_id", payload.OrcamentoID,
		"status_anterior", payload.Status,
		"valor", payload.Valor,
		"motivo", payload.MotivoCancelamento)

	// Cancelar conta a pagar associada ao orçamento
	if err := h.contaPagarService.CancelarContaDeOrcamento(ctx, payload.OrcamentoID); err != nil {
		h.logger.ErrorContext(ctx, "falha ao cancelar conta a pagar do orçamento excluído", 
			"orcamento_id", payload.OrcamentoID,
			"erro", err)
		return
	}

	h.logger.InfoContext(ctx, "conta a pagar cancelada automaticamente devido à exclusão do orçamento", 
		"orcamento_id", payload.OrcamentoID,
		"valor", payload.Valor)
}

// ConfigurarEventHandlers configura os handlers de eventos
func ConfigurarEventHandlers(eventBus bus.EventBus, handler *FinanceiroEventHandler) {
	// Eventos de cronograma de recebimento
	eventBus.Subscrever(events.CronogramaRecebimentoCriado, handler.HandleCronogramaRecebimentoCriado)
	eventBus.Subscrever(events.RecebimentoRealizado, handler.HandleRecebimentoRealizado)
	
	// Eventos de orçamento (integração com Suprimentos)
	eventBus.Subscrever(events.OrcamentoStatusAtualizado, handler.HandleOrcamentoStatusAtualizado)
	eventBus.Subscrever(events.OrcamentoExcluido, handler.HandleOrcamentoExcluido)
	
	// Eventos de apontamento (integração com Pessoal)
	eventBus.Subscrever(events.ApontamentoAprovado, handler.HandleApontamentoAprovado)
	eventBus.Subscrever(events.ApontamentoCancelado, handler.HandleApontamentoCancelado)
}