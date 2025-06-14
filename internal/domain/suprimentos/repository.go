package suprimentos

import "context"

// Repository define o contrato para a persistência de Fornecedores.
type FornecedorRepository interface {
	Salvar(ctx context.Context, fornecedor *Fornecedor) error
	ListarTodos(ctx context.Context) ([]*Fornecedor, error)
}

// MaterialRepository define o contrato para a persistência de Materiais.
type MaterialRepository interface {
	Salvar(ctx context.Context, material *Material) error
	ListarTodos(ctx context.Context) ([]*Material, error)
}
