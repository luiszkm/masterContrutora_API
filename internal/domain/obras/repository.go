package obras

import "context"

type Repository interface {
	// Salvar persiste uma nova Obra ou atualiza uma existente.
	Salvar(ctx context.Context, obra *Obra) error
	// BuscarPorID encontra uma obra pelo seu identificador.
	BuscarPorID(ctx context.Context, id string) (*Obra, error)
}
