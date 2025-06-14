package pessoal

import "context"

// Repository define o contrato para a persistência de Funcionarios.
type FuncionarioRepository interface {
	Salvar(ctx context.Context, funcionario *Funcionario) error
	BuscarPorID(ctx context.Context, funcionarioID string) (*Funcionario, error)
}
