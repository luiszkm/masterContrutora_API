// file: internal/events/pessoal_events.go
package events

import "time"

const (
	ApontamentoAprovado           = "pessoal:apontamento_aprovado"
	PagamentoApontamentoRealizado = "pessoal:pagamento_apontamento_realizado"
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

// PagamentoApontamentoRealizadoPayload são os dados que o evento carrega.
type PagamentoApontamentoRealizadoPayload struct {
	FuncionarioID     string
	ObraID            string
	PeriodoReferencia string
	ValorCalculado    float64
	DataDeEfetivacao  time.Time
	ContaBancariaID   string
}
