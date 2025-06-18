// file: internal/domain/financeiro/pagamento.go
package financeiro

import (
	"time"
)

// RegistroDePagamento é o agregado que representa uma transação financeira.
type RegistroDePagamento struct {
	ID                string
	FuncionarioID     string
	ObraID            string
	PeriodoReferencia string // Ex: "Junho/2025"
	ValorCalculado    float64
	DataDeEfetivacao  time.Time
	ContaBancariaID   string // ID da conta da empresa de onde o dinheiro saiu
}
