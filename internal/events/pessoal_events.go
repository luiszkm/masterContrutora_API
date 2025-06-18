// file: internal/events/pessoal_events.go
package events

import "time"

const (
	PagamentoApontamentoRealizado = "pessoal:pagamento_apontamento_realizado"
)

// PagamentoApontamentoRealizadoPayload são os dados que o evento carrega.
type PagamentoApontamentoRealizadoPayload struct {
	FuncionarioID     string
	ObraID            string
	PeriodoReferencia string
	ValorCalculado    float64
	DataDeEfetivacao  time.Time
	ContaBancariaID   string
}
