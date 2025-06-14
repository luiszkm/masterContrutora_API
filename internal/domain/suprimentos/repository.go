package suprimentos

import "context"

// Repository define o contrato para a persistência de Fornecedores.
type FornecedorRepository interface {
	Salvar(ctx context.Context, fornecedor *Fornecedor) error
	ListarTodos(ctx context.Context) ([]*Fornecedor, error)
	BuscarPorID(ctx context.Context, id string) (*Fornecedor, error)
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
	ListarTodos(ctx context.Context) ([]*Orcamento, error)
	Atualizar(ctx context.Context, orcamento *Orcamento) error // NOVO
}
