// file: internal/repository/postgres/etapa_repository.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

// EtapaRepositoryPostgres implementa a persistência para o agregado Etapa.
type EtapaRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoEtapaRepository(db *pgxpool.Pool, logger *slog.Logger) obras.EtapaRepository {
	return &EtapaRepositoryPostgres{
		db:     db,
		logger: logger,
	}
}

func (r *EtapaRepositoryPostgres) Salvar(ctx context.Context, dbtx db.DBTX, etapa *obras.Etapa) error {
	const op = "repository.postgres.etapa.Salvar"
	query := `
		INSERT INTO etapas (id, obra_id, nome, data_inicio_prevista, data_fim_prevista, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := dbtx.Exec(ctx, query,
		etapa.ID, etapa.ObraID, etapa.Nome,
		etapa.DataInicioPrevista, etapa.DataFimPrevista, etapa.Status,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *EtapaRepositoryPostgres) BuscarPorID(ctx context.Context, etapaID string) (*obras.Etapa, error) {
	const op = "repository.postgres.etapa.BuscarPorID"
	query := `SELECT id, obra_id, nome, data_inicio_prevista, data_fim_prevista, status FROM etapas WHERE id = $1`
	row := r.db.QueryRow(ctx, query, etapaID)

	var etapa obras.Etapa
	err := row.Scan(
		&etapa.ID,
		&etapa.ObraID,
		&etapa.Nome,
		&etapa.DataInicioPrevista,
		&etapa.DataFimPrevista,
		&etapa.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &etapa, nil
}

func (r *EtapaRepositoryPostgres) Atualizar(ctx context.Context, etapa *obras.Etapa) error {
	const op = "repository.postgres.etapa.Atualizar"
	query := `UPDATE etapas SET nome = $1, data_inicio_prevista = $2, data_fim_prevista = $3, status = $4 WHERE id = $5`

	cmd, err := r.db.Exec(ctx, query,
		etapa.Nome,
		etapa.DataInicioPrevista,
		etapa.DataFimPrevista,
		etapa.Status,
		etapa.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	// Verifica se alguma linha foi realmente atualizada
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}

	return nil
}

func (r *EtapaRepositoryPostgres) ListarPorObraID(ctx context.Context, obraID string) ([]*obras.Etapa, error) {
	const op = "repository.postgres.etapa.ListarPorObraID"
	query := `
		SELECT id, obra_id, nome, data_inicio_prevista, data_fim_prevista, status
		FROM etapas
		WHERE obra_id = $1
		ORDER BY data_inicio_prevista, nome ASC
	`
	rows, err := r.db.Query(ctx, query, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	etapas, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[obras.Etapa])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*obras.Etapa{}, nil // Retorna lista vazia se não houver etapas
		}
		return nil, fmt.Errorf("%s: falha ao escanear etapas: %w", op, err)
	}
	return etapas, nil
}
