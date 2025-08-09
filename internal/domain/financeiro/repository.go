package financeiro

import (
	"context"
	"time"

	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

// Repository define o contrato para a persistência de pagamentos.
type Repository interface {
	Salvar(ctx context.Context, db db.DBTX, pagamento *RegistroDePagamento) error
	ListarPagamentos(ctx context.Context, filtros common.ListarFiltros) ([]*RegistroDePagamento, *common.PaginacaoInfo, error)
}

// ContaReceberRepository define o contrato para persistência de contas a receber
type ContaReceberRepository interface {
	Salvar(ctx context.Context, db db.DBTX, conta *ContaReceber) error
	Atualizar(ctx context.Context, conta *ContaReceber) error
	BuscarPorID(ctx context.Context, id string) (*ContaReceber, error)
	ListarPorObraID(ctx context.Context, obraID string) ([]*ContaReceber, error)
	ListarVencidas(ctx context.Context) ([]*ContaReceber, error)
	ListarVencidasPorPeriodo(ctx context.Context, dataInicio, dataFim time.Time) ([]*ContaReceber, error)
	ListarPorStatus(ctx context.Context, status string) ([]*ContaReceber, error)
	ListarPorCliente(ctx context.Context, cliente string) ([]*ContaReceber, error)
	Listar(ctx context.Context, filtros common.ListarFiltros) ([]*ContaReceber, *common.PaginacaoInfo, error)
	Deletar(ctx context.Context, id string) error
}

// ContaPagarRepository define o contrato para persistência de contas a pagar
type ContaPagarRepository interface {
	Salvar(ctx context.Context, db db.DBTX, conta *ContaPagar) error
	Atualizar(ctx context.Context, conta *ContaPagar) error
	BuscarPorID(ctx context.Context, id string) (*ContaPagar, error)
	ListarPorObraID(ctx context.Context, obraID string) ([]*ContaPagar, error)
	ListarPorFornecedorID(ctx context.Context, fornecedorID string) ([]*ContaPagar, error)
	ListarPorOrcamentoID(ctx context.Context, orcamentoID string) ([]*ContaPagar, error)
	ListarVencidas(ctx context.Context) ([]*ContaPagar, error)
	ListarVencidasPorPeriodo(ctx context.Context, dataInicio, dataFim time.Time) ([]*ContaPagar, error)
	ListarPorStatus(ctx context.Context, status string) ([]*ContaPagar, error)
	ListarPorFornecedor(ctx context.Context, fornecedorNome string) ([]*ContaPagar, error)
	Listar(ctx context.Context, filtros common.ListarFiltros) ([]*ContaPagar, *common.PaginacaoInfo, error)
	Deletar(ctx context.Context, id string) error
}

// ParcelaContaPagarRepository define o contrato para persistência de parcelas
type ParcelaContaPagarRepository interface {
	Salvar(ctx context.Context, db db.DBTX, parcela *ParcelaContaPagar) error
	SalvarMuitas(ctx context.Context, db db.DBTX, parcelas []*ParcelaContaPagar) error
	Atualizar(ctx context.Context, parcela *ParcelaContaPagar) error
	BuscarPorID(ctx context.Context, id string) (*ParcelaContaPagar, error)
	ListarPorContaPagarID(ctx context.Context, contaPagarID string) ([]*ParcelaContaPagar, error)
	ListarVencidas(ctx context.Context) ([]*ParcelaContaPagar, error)
	Deletar(ctx context.Context, id string) error
}
