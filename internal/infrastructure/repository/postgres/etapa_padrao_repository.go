// file: internal/infrastructure/repository/postgres/etapa_padrao_repository.go
package postgres

import (
	"context"
	"log/slog"

	// ... imports
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
)

type EtapaPadraoRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoEtapaPadraoRepository(db *pgxpool.Pool, logger *slog.Logger) *EtapaPadraoRepositoryPostgres {
	return &EtapaPadraoRepositoryPostgres{db: db, logger: logger}
}

func (r *EtapaPadraoRepositoryPostgres) Salvar(ctx context.Context, e *obras.EtapaPadrao) error {
	query := `INSERT INTO etapas_padrao (id, nome, descricao, ordem, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, e.ID, e.Nome, e.Descricao, e.Ordem, e.CreatedAt, e.UpdatedAt)
	return err
}

func (r *EtapaPadraoRepositoryPostgres) ListarTodas(ctx context.Context) ([]*obras.EtapaPadrao, error) {
	query := `SELECT id, nome, descricao, ordem, created_at, updated_at FROM etapas_padrao ORDER BY ordem, nome ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[obras.EtapaPadrao])
}

func (r *EtapaPadraoRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*obras.EtapaPadrao, error) {
	query := `SELECT id, nome, descricao, ordem, created_at, updated_at FROM etapas_padrao WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var etapa obras.EtapaPadrao
	err := row.Scan(&etapa.ID, &etapa.Nome, &etapa.Descricao, &etapa.Ordem, &etapa.CreatedAt, &etapa.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &etapa, nil
}

func (r *EtapaPadraoRepositoryPostgres) Atualizar(ctx context.Context, e *obras.EtapaPadrao) error {
	query := `UPDATE etapas_padrao SET nome = $1, descricao = $2, ordem = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx, query, e.Nome, e.Descricao, e.Ordem, e.UpdatedAt, e.ID)
	return err
}

func (r *EtapaPadraoRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	query := `DELETE FROM etapas_padrao WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
