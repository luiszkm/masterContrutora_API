package obras

import (
	"context"

	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
)

type ObrasRepository interface {
	// Salvar persiste uma nova Obra ou atualiza uma existente.
	Salvar(ctx context.Context, obra *Obra) error
	// BuscarPorID encontra uma obra pelo seu identificador.
	BuscarPorID(ctx context.Context, id string) (*Obra, error)
	BuscarDashboardPorID(ctx context.Context, id string) (*dto.ObraDashboard, error)
	ListarObras(ctx context.Context) ([]*dto.ObraListItemDTO, error)
}

type AlocacaoRepository interface {
	Salvar(ctx context.Context, alocacao *Alocacao) error
}

type EtapaRepository interface {
	Salvar(ctx context.Context, etapa *Etapa) error
	BuscarPorID(ctx context.Context, etapaID string) (*Etapa, error) // NOVO
	Atualizar(ctx context.Context, etapa *Etapa) error               // NOVO
}
