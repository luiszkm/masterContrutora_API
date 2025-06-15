// file: internal/repository/postgres/fornecedor_repository.go
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

// FornecedorRepositoryPostgres implementa a interface suprimentos.Repository.
type FornecedorRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoFornecedorRepository(db *pgxpool.Pool, logger *slog.Logger) suprimentos.FornecedorRepository {
	return &FornecedorRepositoryPostgres{db: db, logger: logger}
}

// Atualizar implements suprimentos.FornecedorRepository.
func (r *FornecedorRepositoryPostgres) Atualizar(ctx context.Context, fornecedor *suprimentos.Fornecedor) error {
	const op = "repository.postgres.fornecedor.Atualizar"
	query := `
		UPDATE fornecedores
		SET nome = $1, cnpj = $2, categoria = $3, contato = $4, email = $5, status = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query, fornecedor.Nome, fornecedor.CNPJ, fornecedor.Categoria, fornecedor.Contato, fornecedor.Email, fornecedor.Status, fornecedor.ID)
	if err != nil {
		return fmt.Errorf("%s: falha ao atualizar fornecedor: %w", op, err)
	}
	return nil

}

// Deletar implements suprimentos.FornecedorRepository.
func (r *FornecedorRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.fornecedor.Deletar"
	query := `
		UPDATE fornecedores
		SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: falha ao atualizar fornecedor: %w", op, err)
	}
	return nil

}

// BuscarPorID implements suprimentos.FornecedorRepository.
func (r *FornecedorRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*suprimentos.Fornecedor, error) {
	const op = "repository.postgres.fornecedor.BuscarPorID"
	query := `SELECT id, nome, cnpj, categoria, contato, email, status FROM fornecedores WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var f suprimentos.Fornecedor
	err := row.Scan(&f.ID, &f.Nome, &f.CNPJ, &f.Categoria, &f.Contato, &f.Email, &f.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: fornecedor não encontrado com ID %s", op, id)
		}
		return nil, fmt.Errorf("%s: falha ao escanear fornecedor: %w", op, err)
	}

	return &f, nil
}

func (r *FornecedorRepositoryPostgres) Salvar(ctx context.Context, f *suprimentos.Fornecedor) error {
	const op = "repository.postgres.fornecedor.Salvar"
	query := `
		INSERT INTO fornecedores (id, nome, cnpj, categoria, contato, email, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, f.ID, f.Nome, f.CNPJ, f.Categoria, f.Contato, f.Email, f.Status)
	if err != nil {
		// TODO: Tratar erro de violação de constraint UNIQUE do CNPJ com um erro customizado.
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *FornecedorRepositoryPostgres) ListarTodos(ctx context.Context) ([]*suprimentos.Fornecedor, error) {
	const op = "repository.postgres.fornecedor.ListarTodos"
	query := `SELECT id, nome, cnpj, categoria, contato, email, status FROM fornecedores ORDER BY nome ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	fornecedores, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*suprimentos.Fornecedor, error) {
		var f suprimentos.Fornecedor
		err := row.Scan(&f.ID, &f.Nome, &f.CNPJ, &f.Categoria, &f.Contato, &f.Email, &f.Status)
		return &f, err
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*suprimentos.Fornecedor{}, nil // Retorna lista vazia se não houver resultados
		}
		return nil, fmt.Errorf("%s: falha ao escanear fornecedores: %w", op, err)
	}

	return fornecedores, nil
}
