package financeiro

import "context"

// Repository define o contrato para a persistência de pagamentos.
type Repository interface {
	Salvar(ctx context.Context, pagamento *RegistroDePagamento) error
}
