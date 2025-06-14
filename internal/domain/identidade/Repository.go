package identidade

import "context"

type Repository interface {
	Salvar(ctx context.Context, usuario *Usuario) error
	BuscarPorEmail(ctx context.Context, email string) (*Usuario, error)
}
