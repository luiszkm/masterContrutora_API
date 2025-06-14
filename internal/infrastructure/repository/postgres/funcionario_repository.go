// file: internal/repository/postgres/funcionario_repository.go
package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
)

type FuncionarioRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// BuscarPorID implements obras.PessoalFinder.
func (r *FuncionarioRepositoryPostgres) BuscarPorID(ctx context.Context, funcionarioID string) (*pessoal.Funcionario, error) {

	const op = "repository.postgres.funcionario.BuscarPorID"
	query := `
		SELECT id, nome, cpf, cargo, data_contratacao, salario, diaria, status
		FROM funcionarios
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, funcionarioID)

	var f pessoal.Funcionario
	err := row.Scan(&f.ID, &f.Nome, &f.CPF, &f.Cargo, &f.DataContratacao, &f.Salario, &f.Diaria, &f.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%s: funcionário não encontrado com ID %s", op, funcionarioID)
		}
		return nil, fmt.Errorf("%s: erro ao buscar funcionário: %w", op, err)
	}
	return &f, nil
}

func NovoFuncionarioRepository(db *pgxpool.Pool, logger *slog.Logger) *FuncionarioRepositoryPostgres {
	return &FuncionarioRepositoryPostgres{db: db, logger: logger}
}

func (r *FuncionarioRepositoryPostgres) Salvar(ctx context.Context, f *pessoal.Funcionario) error {
	const op = "repository.postgres.funcionario.Salvar"
	query := `
		INSERT INTO funcionarios (id, nome, cpf, cargo, data_contratacao, salario, diaria, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query, f.ID, f.Nome, f.CPF, f.Cargo, f.DataContratacao, f.Salario, f.Diaria, f.Status)
	if err != nil {
		// TODO: Tratar erro de violação de constraint UNIQUE do CPF
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
