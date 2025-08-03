package financeiro

import (
	"context"

	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

// Repository define o contrato para a persistência de pagamentos.
type Repository interface {
	Salvar(ctx context.Context, db db.DBTX, pagamento *RegistroDePagamento) error
	ListarPagamentos(ctx context.Context, filtros common.ListarFiltros) ([]*RegistroDePagamento, *common.PaginacaoInfo, error)
}
