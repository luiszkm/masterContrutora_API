// file: internal/events/orcamento_events.go
package events

const (
	// OrcamentoStatusAtualizado é o nome do nosso tópico de evento.
	OrcamentoStatusAtualizado = "orcamento:status_atualizado"
	// OrcamentoExcluido é disparado quando um orçamento é soft deleted
	OrcamentoExcluido = "orcamento:excluido"
)

// OrcamentoStatusAtualizadoPayload contém os dados que queremos enviar no evento.
type OrcamentoStatusAtualizadoPayload struct {
	OrcamentoID   string
	EtapaID       string
	StatusAnterior string  // Status antes da mudança
	NovoStatus    string  // Status após a mudança
	Valor         float64
}

// OrcamentoExcluidoPayload contém dados do orçamento excluído
type OrcamentoExcluidoPayload struct {
	OrcamentoID        string  `json:"orcamentoId"`
	EtapaID            string  `json:"etapaId"`
	Status             string  `json:"status"`        // Status no momento da exclusão
	Valor              float64 `json:"valor"`
	MotivoCancelamento string  `json:"motivoCancelamento"`
}
