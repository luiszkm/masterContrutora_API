// file: internal/service/obras/events/handler.go
package events

import (
	"context"
	"log/slog"

	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
)

// ObrasEventHandler lida com eventos destinados ao contexto de Obras.
type ObrasEventHandler struct {
	logger *slog.Logger
}

func NovoObrasEventHandler(logger *slog.Logger) *ObrasEventHandler {
	return &ObrasEventHandler{logger: logger}
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
