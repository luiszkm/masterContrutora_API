package events

import (
	"context"
	"log/slog"

	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

// OrcamentoService interface para atualizar orçamentos
type OrcamentoService interface {
	AtualizarStatusOrcamento(ctx context.Context, orcamentoID string, input dto.AtualizarStatusOrcamentoInput) (*suprimentos.Orcamento, error)
}

// SuprimentosEventHandler lida com eventos destinados ao contexto de Suprimentos.
type SuprimentosEventHandler struct {
	orcamentoService OrcamentoService
	logger           *slog.Logger
}

func NovoSuprimentosEventHandler(orcamentoService OrcamentoService, logger *slog.Logger) *SuprimentosEventHandler {
	return &SuprimentosEventHandler{
		orcamentoService: orcamentoService,
		logger:           logger.With("handler", "SuprimentosEventHandler"),
	}
}

// HandleContaPagarPaga processa evento de conta a pagar paga para atualizar orçamento
func (h *SuprimentosEventHandler) HandleContaPagarPaga(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.ContaPagarPagaPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de conta pagar paga inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando pagamento de conta a pagar", 
		"conta_id", payload.ContaPagarID,
		"orcamento_id", payload.OrcamentoID,
		"valor_pago", payload.ValorPago,
		"status", payload.Status)

	// Se não há orçamento vinculado, não há o que atualizar
	if payload.OrcamentoID == nil {
		h.logger.InfoContext(ctx, "conta a pagar sem orçamento vinculado, nada para atualizar", 
			"conta_id", payload.ContaPagarID)
		return
	}

	// Apenas atualizar orçamento quando a conta for totalmente paga
	if payload.Status != "PAGO" {
		h.logger.InfoContext(ctx, "conta a pagar não está totalmente paga, orçamento não será atualizado", 
			"conta_id", payload.ContaPagarID,
			"orcamento_id", *payload.OrcamentoID,
			"status", payload.Status)
		return
	}

	orcamentoID := *payload.OrcamentoID

	h.logger.InfoContext(ctx, "conta a pagar totalmente paga, atualizando orçamento para 'Pago'", 
		"conta_id", payload.ContaPagarID,
		"orcamento_id", orcamentoID,
		"valor_total_pago", payload.ValorTotalPago)

	// Atualizar status do orçamento para "Pago"
	input := dto.AtualizarStatusOrcamentoInput{
		Status: "Pago",
	}

	orcamentoAtualizado, err := h.orcamentoService.AtualizarStatusOrcamento(ctx, orcamentoID, input)
	if err != nil {
		h.logger.ErrorContext(ctx, "falha ao atualizar status do orçamento", 
			"orcamento_id", orcamentoID,
			"conta_id", payload.ContaPagarID,
			"novo_status", "Pago",
			"erro", err)
		return
	}

	h.logger.InfoContext(ctx, "orçamento atualizado para 'Pago' com sucesso", 
		"orcamento_id", orcamentoID,
		"conta_id", payload.ContaPagarID,
		"novo_status", orcamentoAtualizado.Status,
		"valor_total", orcamentoAtualizado.ValorTotal)
}

// ConfigurarEventHandlers configura os handlers de eventos
func ConfigurarEventHandlers(eventBus bus.EventBus, handler *SuprimentosEventHandler) {
	// Eventos de contas a pagar (integração com Financeiro)
	eventBus.Subscrever(events.ContaPagarPaga, handler.HandleContaPagarPaga)
}