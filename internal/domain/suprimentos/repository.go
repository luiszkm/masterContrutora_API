package suprimentos

import (
	"context"

	"github.com/luiszkm/masterCostrutora/internal/domain/common"
)

// Repository define o contrato para a persistência de Fornecedores.
type FornecedorRepository interface {
	Salvar(ctx context.Context, fornecedor *Fornecedor, categoriaIDs []string) error
	Atualizar(ctx context.Context, fornecedor *Fornecedor, categoriaIDs *[]string) error
	BuscarPorID(ctx context.Context, id string) (*Fornecedor, error)
	ListarTodos(ctx context.Context) ([]*Fornecedor, error)
	Deletar(ctx context.Context, id string) error
}

// MaterialRepository define o contrato para a persistência de Materiais.
type MaterialRepository interface {
	Salvar(ctx context.Context, material *Material) error
	ListarTodos(ctx context.Context) ([]*Material, error)
	BuscarPorID(ctx context.Context, id string) (*Material, error)
}

// OrcamentoRepository define o contrato para a persistência de Orçamentos.
// Note que ele lida com o agregado completo (Orçamento + Itens).
type OrcamentoRepository interface {
	Salvar(ctx context.Context, orcamento *Orcamento) error
	ListarPorEtapaID(ctx context.Context, etapaID string) ([]*Orcamento, error)
	BuscarPorID(ctx context.Context, id string) (*Orcamento, error)
	ListarTodos(ctx context.Context, filtros common.ListarFiltros) ([]*Orcamento, *common.PaginacaoInfo, error)
	Atualizar(ctx context.Context, orcamento *Orcamento) error // NOVO
}

type CategoriaRepository interface {
	Salvar(ctx context.Context, categoria *Categoria) error
	BuscarPorID(ctx context.Context, id string) (*Categoria, error)
	ListarTodas(ctx context.Context) ([]*Categoria, error)
	Atualizar(ctx context.Context, categoria *Categoria) error
	Deletar(ctx context.Context, id string) error
}
