package pessoal

import "context"

// Repository define o contrato para a persistência de Funcionarios.
type FuncionarioRepository interface {
	Salvar(ctx context.Context, funcionario *Funcionario) error
	BuscarPorID(ctx context.Context, funcionarioID string) (*Funcionario, error)
	Deletar(ctx context.Context, id string) error // NOVO
	Listar(ctx context.Context) ([]*Funcionario, error)
	Atualizar(ctx context.Context, funcionario *Funcionario) error // NOVO
}
