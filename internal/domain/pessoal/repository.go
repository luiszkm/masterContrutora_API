package pessoal

import "context"

// Repository define o contrato para a persistência de Funcionarios.
type Repository interface {
	Salvar(ctx context.Context, funcionario *Funcionario) error
}
