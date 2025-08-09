// file: internal/platform/bus/eventbus.go
package bus

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Evento define a estrutura básica de um evento no nosso sistema.
type Evento struct {
	Nome    string
	Payload any
}

// HandlerFunc é o tipo da função que irá tratar um evento.
type HandlerFunc func(ctx context.Context, evento Evento)

// EventBus gerencia a subscrição e publicação de eventos de forma assíncrona.
type EventBus struct {
	handlers map[string][]HandlerFunc
	mu       sync.RWMutex
	logger   *slog.Logger
}

func NovoEventBus(logger *slog.Logger) *EventBus {
	return &EventBus{
		handlers: make(map[string][]HandlerFunc),
		logger:   logger,
	}
}

// Subscrever adiciona um novo handler para um tópico de evento.
func (b *EventBus) Subscrever(nomeEvento string, handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[nomeEvento] = append(b.handlers[nomeEvento], handler)
}

// Publicar envia um evento para todos os handlers subscritos, cada um em sua própria goroutine.
func (b *EventBus) Publicar(ctx context.Context, evento Evento) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, ok := b.handlers[evento.Nome]; ok {
		for _, handler := range handlers {
			// Executa cada handler de forma assíncrona.
			go func(h HandlerFunc) {
				// Criar contexto com timeout para handlers de eventos
				// Isso evita que o contexto seja cancelado quando a requisição HTTP termina
				handlerCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				
				b.logger.InfoContext(handlerCtx, "processando evento", "evento", evento.Nome)
				// Em um sistema real, adicionaríamos retentativas e DLQ aqui (ADR-007)
				h(handlerCtx, evento)
			}(handler)
		}
	}
}
