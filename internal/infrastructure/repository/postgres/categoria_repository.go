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

var (
	ErrCategoriaJaExiste = fmt.Errorf("categoria já existe")
)

// CategoriaRepositoryPostgres implementa a interface suprimentos.CategoriaRepository.
type CategoriaRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NovoCategoriaRepository é o construtor para o repositório de categorias.
func NovoCategoriaRepository(db *pgxpool.Pool, logger *slog.Logger) *CategoriaRepositoryPostgres {
	return &CategoriaRepositoryPostgres{db: db, logger: logger}
}

// Salvar insere uma nova categoria no banco de dados.
func (r *CategoriaRepositoryPostgres) Salvar(ctx context.Context, c *suprimentos.Categoria) error {
	const op = "repository.postgres.categoria.Salvar"
	query := `INSERT INTO categorias (id, nome) VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, query, c.ID, c.Nome)
	if err != nil {
		// TODO: Adicionar tratamento para erro de violação de constraint UNIQUE do nome.
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// BuscarPorID encontra uma categoria pelo seu ID.
func (r *CategoriaRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*suprimentos.Categoria, error) {
	const op = "repository.postgres.categoria.BuscarPorID"
	query := `SELECT id, nome FROM categorias WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var c suprimentos.Categoria
	err := row.Scan(&c.ID, &c.Nome)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &c, nil
}

// ListarTodos busca todas as categorias no banco de dados.
func (r *CategoriaRepositoryPostgres) ListarTodos(ctx context.Context) ([]*suprimentos.Categoria, error) {
	const op = "repository.postgres.categoria.ListarTodos"
	query := `SELECT id, nome FROM categorias ORDER BY nome ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	categorias, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[suprimentos.Categoria])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*suprimentos.Categoria{}, nil // Retorna lista vazia se não houver resultados
		}
		return nil, fmt.Errorf("%s: falha ao escanear categorias: %w", op, err)
	}
	return categorias, nil
}

// Atualizar modifica os dados de uma categoria existente.
func (r *CategoriaRepositoryPostgres) Atualizar(ctx context.Context, c *suprimentos.Categoria) error {
	const op = "repository.postgres.categoria.Atualizar"
	query := `UPDATE categorias SET nome = $1 WHERE id = $2`

	cmd, err := r.db.Exec(ctx, query, c.Nome, c.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

// Deletar remove permanentemente uma categoria do banco de dados.
func (r *CategoriaRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.categoria.Deletar"
	query := `DELETE FROM categorias WHERE id = $1`

	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		// Se houver uma violação de chave estrangeira, o banco retornará um erro.
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
