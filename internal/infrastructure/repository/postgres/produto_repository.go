// file: internal/infrastructure/repository/postgres/produto_repository.go
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

// ProdutoRepositoryPostgres implementa a interface suprimentos.ProdutoRepository.
type ProdutoRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoProdutoRepository(db *pgxpool.Pool, logger *slog.Logger) suprimentos.ProdutoRepository {
	return &ProdutoRepositoryPostgres{db: db, logger: logger}
}

func (r *ProdutoRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*suprimentos.Produto, error) {
	const op = "repository.postgres.produto.BuscarPorID"
	query := `SELECT id, nome, descricao, unidade_de_medida, categoria, created_at, updated_at, deleted_at FROM produtos WHERE id = $1 AND deleted_at IS NULL`
	row := r.db.QueryRow(ctx, query, id)

	var p suprimentos.Produto
	err := row.Scan(&p.ID, &p.Nome, &p.Descricao, &p.UnidadeDeMedida, &p.Categoria, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: falha ao escanear produto: %w", op, err)
	}
	return &p, nil
}

func (r *ProdutoRepositoryPostgres) Salvar(ctx context.Context, p *suprimentos.Produto) error {
	const op = "repository.postgres.produto.Salvar"
	query := `
		INSERT INTO produtos (id, nome, descricao, unidade_de_medida, categoria, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, p.ID, p.Nome, p.Descricao, p.UnidadeDeMedida, p.Categoria, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *ProdutoRepositoryPostgres) ListarTodos(ctx context.Context) ([]*suprimentos.Produto, error) {
	const op = "repository.postgres.produto.ListarTodos"
	query := `SELECT id, nome, descricao, unidade_de_medida, categoria, created_at, updated_at, deleted_at FROM produtos WHERE deleted_at IS NULL ORDER BY nome ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	produtos, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[suprimentos.Produto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*suprimentos.Produto{}, nil
		}
		return nil, fmt.Errorf("%s: falha ao escanear produtos: %w", op, err)
	}
	return produtos, nil
}

func (r *ProdutoRepositoryPostgres) BuscarPorNome(ctx context.Context, nome string) (*suprimentos.Produto, error) {
	const op = "repository.postgres.produto.BuscarPorNome"
	query := `SELECT id, nome, descricao, unidade_de_medida, categoria, created_at, updated_at, deleted_at FROM produtos WHERE nome = $1 AND deleted_at IS NULL`
	row := r.db.QueryRow(ctx, query, nome)

	var p suprimentos.Produto
	err := row.Scan(&p.ID, &p.Nome, &p.Descricao, &p.UnidadeDeMedida, &p.Categoria, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado // Usa o erro padrão para indicar que não foi encontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &p, nil
}

func (r *ProdutoRepositoryPostgres) Atualizar(ctx context.Context, p *suprimentos.Produto) error {
	const op = "repository.postgres.produto.Atualizar"
	query := `UPDATE produtos SET nome = $1, descricao = $2, unidade_de_medida = $3, categoria = $4, updated_at = $5 WHERE id = $6 AND deleted_at IS NULL`
	cmd, err := r.db.Exec(ctx, query, p.Nome, p.Descricao, p.UnidadeDeMedida, p.Categoria, p.UpdatedAt, p.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

func (r *ProdutoRepositoryPostgres) SoftDelete(ctx context.Context, id string) error {
	const op = "repository.postgres.produto.SoftDelete"
	query := `UPDATE produtos SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
