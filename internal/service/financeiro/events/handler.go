package events

import (
	"context"
	"log/slog"
	"time"

	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/service/financeiro/dto"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
)

// ContaReceberService interface para o service de contas a receber
type ContaReceberService interface {
	CriarConta(ctx context.Context, input dto.CriarContaReceberInput) (*dto.ContaReceberOutput, error)
}

// ContaPagarService interface para o service de contas a pagar
type ContaPagarService interface {
	CriarContaDeOrcamento(ctx context.Context, input dto.CriarContaPagarDeOrcamentoInput, orcamento interface{}) (*dto.ContaPagarOutput, error)
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
	payload, ok := evento.Payload.(events.RecebimentoRealizadoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de recebimento inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando recebimento realizado", 
		"cliente", payload.Cliente,
		"valor", payload.ValorRecebido,
		"cronograma_id", payload.CronogramaRecebimentoID)

	// Este evento pode disparar outros processos, como:
	// - Atualização de movimentações financeiras
	// - Notificações para o cliente
	// - Relatórios de recebimento
	// Por enquanto, apenas logamos

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
}

// HandlePagamentoApontamentoRealizado processa pagamento de apontamento
func (h *FinanceiroEventHandler) HandlePagamentoApontamentoRealizado(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.PagamentoApontamentoRealizadoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de pagamento de apontamento inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando pagamento de apontamento", 
		"funcionario_id", payload.FuncionarioID,
		"obra_id", payload.ObraID,
		"valor", payload.ValorCalculado)

	// Este evento já é processado pelo handler existente em financeiro/handler.go
	// Aqui poderíamos adicionar lógica adicional se necessário

	// TODO: Implementar criação de MovimentacaoFinanceira quando a entidade existir
}

// ConfigurarEventHandlers configura os handlers de eventos
func ConfigurarEventHandlers(eventBus bus.EventBus, handler *FinanceiroEventHandler) {
	// Eventos de cronograma de recebimento
	eventBus.Subscrever(events.CronogramaRecebimentoCriado, handler.HandleCronogramaRecebimentoCriado)
	eventBus.Subscrever(events.RecebimentoRealizado, handler.HandleRecebimentoRealizado)
	
	// Eventos de orçamento (integração com Suprimentos)
	eventBus.Subscrever(events.OrcamentoStatusAtualizado, handler.HandleOrcamentoStatusAtualizado)
	
	// Eventos de pagamento (integração com Pessoal)
	eventBus.Subscrever(events.PagamentoApontamentoRealizado, handler.HandlePagamentoApontamentoRealizado)
}