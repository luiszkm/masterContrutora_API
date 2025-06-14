// file: internal/events/orcamento_events.go
package events

const (
	// OrcamentoStatusAtualizado é o nome do nosso tópico de evento.
	OrcamentoStatusAtualizado = "orcamento:status_atualizado"
)

// OrcamentoStatusAtualizadoPayload contém os dados que queremos enviar no evento.
type OrcamentoStatusAtualizadoPayload struct {
	OrcamentoID string
	EtapaID     string
	NovoStatus  string
	Valor       float64
}
