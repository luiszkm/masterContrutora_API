package suprimentos

import (
	"context"
	"time"

	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

// Repository define o contrato para a persistência de Fornecedores.
type FornecedorRepository interface {
	Salvar(ctx context.Context, fornecedor *Fornecedor, categoriaIDs []string) error
	Atualizar(ctx context.Context, fornecedor *Fornecedor, categoriaIDs *[]string) error
	BuscarPorID(ctx context.Context, id string) (*Fornecedor, error)
	ListarTodos(ctx context.Context) ([]*Fornecedor, error)
	Deletar(ctx context.Context, id string) error
}

// MaterialRepository define o contrato para a persistência de produto.
type ProdutoRepository interface {
	Salvar(ctx context.Context, material *Produto) error
	ListarTodos(ctx context.Context) ([]*Produto, error)
	BuscarPorID(ctx context.Context, id string) (*Produto, error)
	BuscarPorNome(ctx context.Context, nome string) (*Produto, error) // NOVO
	SoftDelete(ctx context.Context, id string) error                  // NOVO
	Atualizar(ctx context.Context, produto *Produto) error           // NOVO

}

// OrcamentoRepository define o contrato para a persistência de Orçamentos.
// Note que ele lida com o agregado completo (Orçamento + Itens).
type OrcamentoRepository interface {
	Salvar(ctx context.Context, orcamento *Orcamento) error
	ListarPorEtapaID(ctx context.Context, etapaID string) ([]*Orcamento, error)
	BuscarPorID(ctx context.Context, id string) (*Orcamento, error)
	ListarOrcamentos(ctx context.Context, filtros common.ListarFiltros) ([]*dto.OrcamentoListItemDTO, *common.PaginacaoInfo, error)
	Atualizar(ctx context.Context, orcamento *Orcamento) error                              // NOVO MÉTODO
	ContarPorMesAno(ctx context.Context, ano int, mes time.Month) (int, error)              // NOVO
	BuscarPorDetalhesID(ctx context.Context, id string) (*dto.OrcamentoDetalhadoDTO, error) // Assinatura atualizada
	SoftDelete(ctx context.Context, id string) error                                        // NOVO
	CompararPorCategoria(ctx context.Context, categoria string, limite int) ([]*dto.OrcamentoComparacao, error) // NOVO: Comparar orçamentos por categoria

}

type CategoriaRepository interface {
	Salvar(ctx context.Context, categoria *Categoria) error
	BuscarPorID(ctx context.Context, id string) (*Categoria, error)
	ListarTodas(ctx context.Context) ([]*Categoria, error)
	Atualizar(ctx context.Context, categoria *Categoria) error
	Deletar(ctx context.Context, id string) error
}
