package financeiro

import (
	"context"

	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

// Repository define o contrato para a persistência de pagamentos.
type Repository interface {
	Salvar(ctx context.Context, db db.DBTX, pagamento *RegistroDePagamento) error // Modificado
}
