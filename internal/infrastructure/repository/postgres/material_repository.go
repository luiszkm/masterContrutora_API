// file: internal/repository/postgres/material_repository.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
)

// MaterialRepositoryPostgres implementa a interface suprimentos.MaterialRepository.
type MaterialRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoMaterialRepository(db *pgxpool.Pool, logger *slog.Logger) *MaterialRepositoryPostgres {
	return &MaterialRepositoryPostgres{db: db, logger: logger}
}

func (r *MaterialRepositoryPostgres) Salvar(ctx context.Context, m *suprimentos.Material) error {
	const op = "repository.postgres.material.Salvar"
	query := `
		INSERT INTO materiais (id, nome, descricao, unidade_de_medida, categoria)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, m.ID, m.Nome, m.Descricao, m.UnidadeDeMedida, m.Categoria)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *MaterialRepositoryPostgres) ListarTodos(ctx context.Context) ([]*suprimentos.Material, error) {
	const op = "repository.postgres.material.ListarTodos"
	query := `SELECT id, nome, descricao, unidade_de_medida, categoria FROM materiais ORDER BY nome ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	materiais, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*suprimentos.Material, error) {
		var m suprimentos.Material
		err := row.Scan(&m.ID, &m.Nome, &m.Descricao, &m.UnidadeDeMedida, &m.Categoria)
		return &m, err
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*suprimentos.Material{}, nil
		}
		return nil, fmt.Errorf("%s: falha ao escanear materiais: %w", op, err)
	}
	return materiais, nil
}
