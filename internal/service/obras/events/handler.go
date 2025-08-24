// file: internal/service/obras/events/handler.go
package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
)

// CronogramaService interface para atualizar cronogramas
type CronogramaService interface {
	RegistrarRecebimento(ctx context.Context, cronogramaID string, input dto.RegistrarRecebimentoInput) (*dto.CronogramaRecebimentoOutput, error)
	BuscarPorID(ctx context.Context, id string) (*dto.CronogramaRecebimentoOutput, error)
	SincronizarRecebimento(ctx context.Context, cronogramaID string, input dto.RegistrarRecebimentoInput) (*dto.CronogramaRecebimentoOutput, error)
}

// ObrasEventHandler lida com eventos destinados ao contexto de Obras.
type ObrasEventHandler struct {
	cronogramaService CronogramaService
	logger            *slog.Logger
}

func NovoObrasEventHandler(cronogramaService CronogramaService, logger *slog.Logger) *ObrasEventHandler {
	return &ObrasEventHandler{
		cronogramaService: cronogramaService,
		logger:            logger,
	}
}

// HandleOrcamentoStatusAtualizado é o método que será subscrito ao evento.
func (h *ObrasEventHandler) HandleOrcamentoStatusAtualizado(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.OrcamentoStatusAtualizadoPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento inválido", "evento", evento.Nome)
		return
	}

	h.logger.Info("EVENTO RECEBIDO PELO CONTEXTO DE OBRAS!", "novo_status", payload.NovoStatus, "orcamento_id", payload.OrcamentoID)

	// TODO: Lógica futura aqui.
	// Por exemplo, poderíamos usar o payload.EtapaID para encontrar a Obra
	// e forçar a atualização de um modelo de leitura (dashboard) em cache.
}

// HandleContaReceberPaga processa evento de conta a receber paga para sincronizar cronograma
func (h *ObrasEventHandler) HandleContaReceberPaga(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.ContaReceberPagaPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento inválido", "evento", evento.Nome)
		return
	}

	// Se não há cronograma vinculado, não há o que sincronizar
	if payload.CronogramaRecebimentoID == nil {
		h.logger.InfoContext(ctx, "conta recebida sem cronograma vinculado, nada para sincronizar", 
			"conta_id", payload.ContaReceberID)
		return
	}

	cronogramaID := *payload.CronogramaRecebimentoID

	h.logger.InfoContext(ctx, "sincronizando recebimento no cronograma", 
		"conta_id", payload.ContaReceberID,
		"cronograma_id", cronogramaID, 
		"valor", payload.ValorRecebido)

	// Verificar se cronograma existe antes de tentar registrar recebimento
	_, err := h.cronogramaService.BuscarPorID(ctx, cronogramaID)
	if err != nil {
		h.logger.ErrorContext(ctx, "cronograma não encontrado para sincronização", 
			"cronograma_id", cronogramaID, 
			"conta_id", payload.ContaReceberID,
			"erro", err)
		return
	}

	// Preparar input para registrar recebimento no cronograma
	// IMPORTANTE: Incluir marcação para evitar loop infinito
	input := dto.RegistrarRecebimentoInput{
		Valor:       payload.ValorRecebido,
		Observacoes: func() *string { s := fmt.Sprintf("SYNC:ContaReceber:%s", payload.ContaReceberID); return &s }(),
	}

	// Sincronizar recebimento no cronograma (sem disparar eventos)
	cronogramaAtualizado, err := h.cronogramaService.SincronizarRecebimento(ctx, cronogramaID, input)
	if err != nil {
		h.logger.ErrorContext(ctx, "falha ao sincronizar recebimento no cronograma", 
			"cronograma_id", cronogramaID,
			"conta_id", payload.ContaReceberID, 
			"valor", payload.ValorRecebido,
			"erro", err)
		return
	}

	h.logger.InfoContext(ctx, "cronograma sincronizado com sucesso", 
		"cronograma_id", cronogramaID,
		"conta_id", payload.ContaReceberID,
		"valor", payload.ValorRecebido,
		"novo_status", cronogramaAtualizado.Status,
		"valor_total_recebido", cronogramaAtualizado.ValorRecebido)
}
