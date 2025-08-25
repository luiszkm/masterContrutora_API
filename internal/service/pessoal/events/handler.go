package events

import (
	"context"
	"log/slog"

	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
)

// ApontamentoService interface para o service de apontamentos
type ApontamentoService interface {
	MarcarComoPago(ctx context.Context, apontamentoID string) error
}

// PessoalEventHandler processa eventos relacionados ao módulo pessoal
type PessoalEventHandler struct {
	apontamentoService ApontamentoService
	logger             *slog.Logger
}

func NovoPessoalEventHandler(
	apontamentoService ApontamentoService,
	logger *slog.Logger,
) *PessoalEventHandler {
	return &PessoalEventHandler{
		apontamentoService: apontamentoService,
		logger:             logger.With("handler", "PessoalEventHandler"),
	}
}

// HandleContaPagarPaga processa evento de conta a pagar paga para atualizar apontamento
func (h *PessoalEventHandler) HandleContaPagarPaga(ctx context.Context, evento bus.Evento) {
	payload, ok := evento.Payload.(events.ContaPagarPagaPayload)
	if !ok {
		h.logger.ErrorContext(ctx, "payload de evento de conta a pagar paga inválido", "evento", evento.Nome)
		return
	}

	h.logger.InfoContext(ctx, "processando conta a pagar paga", 
		"conta_id", payload.ContaPagarID,
		"tipo", payload.TipoContaPagar,
		"fornecedor", payload.FornecedorNome,
		"numero_documento", payload.NumeroDocumento,
		"valor_pago", payload.ValorPago,
		"status", payload.Status)

	// Só processa contas de funcionário que estão totalmente pagas
	if payload.TipoContaPagar != "FUNCIONARIO" {
		h.logger.InfoContext(ctx, "conta a pagar não é de funcionário, ignorando evento", 
			"tipo", payload.TipoContaPagar)
		return
	}

	if payload.Status != "PAGO" {
		h.logger.InfoContext(ctx, "conta a pagar não está totalmente paga, ignorando evento", 
			"status", payload.Status)
		return
	}

	// Verifica se há número do documento (apontamento ID)
	if payload.NumeroDocumento == nil || *payload.NumeroDocumento == "" {
		h.logger.WarnContext(ctx, "conta a pagar de funcionário sem número do documento (apontamento ID)", 
			"conta_id", payload.ContaPagarID)
		return
	}

	apontamentoID := *payload.NumeroDocumento

	h.logger.InfoContext(ctx, "atualizando apontamento para PAGO", 
		"apontamento_id", apontamentoID,
		"conta_pagar_id", payload.ContaPagarID)

	// Marca o apontamento como pago
	if err := h.apontamentoService.MarcarComoPago(ctx, apontamentoID); err != nil {
		h.logger.ErrorContext(ctx, "falha ao marcar apontamento como pago", 
			"apontamento_id", apontamentoID,
			"conta_pagar_id", payload.ContaPagarID,
			"erro", err)
		return
	}

	h.logger.InfoContext(ctx, "apontamento marcado como PAGO com sucesso", 
		"apontamento_id", apontamentoID,
		"conta_pagar_id", payload.ContaPagarID,
		"valor_pago", payload.ValorPago)
}

// ConfigurarEventHandlers configura os handlers de eventos
func ConfigurarEventHandlers(eventBus bus.EventBus, handler *PessoalEventHandler) {
	// Eventos de conta a pagar (integração com Financeiro)
	eventBus.Subscrever(events.ContaPagarPaga, handler.HandleContaPagarPaga)
}