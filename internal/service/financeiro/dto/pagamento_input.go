// file: internal/service/financeiro/dto/input.go
package dto

// RegistrarPagamentoInput é o DTO para o caso de uso de registro de pagamento.
type RegistrarPagamentoInput struct {
	FuncionarioID     string
	ObraID            string
	PeriodoReferencia string
	ValorCalculado    float64
	ContaBancariaID   string
}
