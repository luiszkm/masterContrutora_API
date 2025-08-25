// file: internal/events/pessoal_events.go
package events

import "time"

const (
	ApontamentoAprovado           = "pessoal:apontamento_aprovado"
	ApontamentoCancelado          = "pessoal:apontamento_cancelado"
)

// ApontamentoAprovadoPayload contém dados do apontamento aprovado
type ApontamentoAprovadoPayload struct {
	ApontamentoID       string    `json:"apontamentoId"`
	FuncionarioID       string    `json:"funcionarioId"`
	FuncionarioNome     string    `json:"funcionarioNome"`
	ObraID              string    `json:"obraId"`
	ObraNome            string    `json:"obraNome"`
	PeriodoReferencia   string    `json:"periodoReferencia"`
	ValorCalculado      float64   `json:"valorCalculado"`
	DataAprovacao       time.Time `json:"dataAprovacao"`
	DataVencimentoPrevisto time.Time `json:"dataVencimentoPrevisto"` // Quando deve ser pago
	UsuarioID           string    `json:"usuarioId"`
}

// ApontamentoCanceladoPayload contém dados do apontamento cancelado
type ApontamentoCanceladoPayload struct {
	ApontamentoID       string    `json:"apontamentoId"`
	FuncionarioID       string    `json:"funcionarioId"`
	FuncionarioNome     string    `json:"funcionarioNome"`
	ObraID              string    `json:"obraId"`
	ObraNome            string    `json:"obraNome"`
	PeriodoReferencia   string    `json:"periodoReferencia"`
	ValorCalculado      float64   `json:"valorCalculado"`
	StatusAnterior      string    `json:"statusAnterior"`     // EM_ABERTO ou APROVADO_PARA_PAGAMENTO
	DataCancelamento    time.Time `json:"dataCancelamento"`
	MotivoCancelamento  string    `json:"motivoCancelamento,omitempty"`
	UsuarioID           string    `json:"usuarioId"`
}

