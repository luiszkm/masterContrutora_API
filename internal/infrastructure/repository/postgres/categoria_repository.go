// file: internal/infrastructure/repository/postgres/categoria_repository.go
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

type CategoriaRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoCategoriaRepository(db *pgxpool.Pool, logger *slog.Logger) *CategoriaRepositoryPostgres {
	return &CategoriaRepositoryPostgres{db: db, logger: logger}
}

func (r *CategoriaRepositoryPostgres) Salvar(ctx context.Context, c *suprimentos.Categoria) error {
	const op = "repository.postgres.categoria.Salvar"
	query := `INSERT INTO categorias (id, nome, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, c.ID, c.Nome, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *CategoriaRepositoryPostgres) Atualizar(ctx context.Context, c *suprimentos.Categoria) error {
	const op = "repository.postgres.categoria.Atualizar"
	query := `UPDATE categorias SET nome = $1, updated_at = $2 WHERE id = $3`
	cmd, err := r.db.Exec(ctx, query, c.Nome, c.UpdatedAt, c.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

func (r *CategoriaRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*suprimentos.Categoria, error) {
	const op = "repository.postgres.categoria.BuscarPorID"
	query := `SELECT id, nome, created_at, updated_at FROM categorias WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var c suprimentos.Categoria
	err := row.Scan(&c.ID, &c.Nome, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &c, nil
}

func (r *CategoriaRepositoryPostgres) ListarTodas(ctx context.Context) ([]*suprimentos.Categoria, error) {
	const op = "repository.postgres.categoria.ListarTodas"
	query := `SELECT id, nome, created_at, updated_at FROM categorias ORDER BY nome ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	categorias, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[suprimentos.Categoria])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*suprimentos.Categoria{}, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return categorias, nil
}

func (r *CategoriaRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.categoria.Deletar"
	// ATENÇÃO: Este é um "hard delete". Se uma categoria estiver em uso por um fornecedor,
	// esta operação falhará devido a restrições de chave estrangeira.
	query := `DELETE FROM categorias WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
