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

func NovoFuncionarioRepository(db *pgxpool.Pool, logger *slog.Logger) pessoal.FuncionarioRepository {
	return &FuncionarioRepositoryPostgres{db: db, logger: logger}
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
func (r *FuncionarioRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.funcionario.Deletar"
	query := `UPDATE funcionarios SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

func (r *FuncionarioRepositoryPostgres) Atualizar(ctx context.Context, f *pessoal.Funcionario) error {
	const op = "repository.postgres.funcionario.Atualizar"
	query := `
		UPDATE funcionarios
		SET nome = $1, cpf = $2, cargo = $3, data_contratacao = $4, salario = $5, diaria = $6, status = $7
		WHERE id = $8 AND deleted_at IS NULL
	`
	cmd, err := r.db.Exec(ctx, query, f.Nome, f.CPF, f.Cargo, f.DataContratacao, f.Salario, f.Diaria, f.Status, f.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
func (r *FuncionarioRepositoryPostgres) Listar(ctx context.Context) ([]*pessoal.Funcionario, error) {
	const op = "repository.postgres.funcionario.Listar"
	query := `
		SELECT id, nome, cpf, cargo, data_contratacao, salario, diaria, status
		FROM funcionarios
		WHERE deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var funcionarios []*pessoal.Funcionario
	for rows.Next() {
		var f pessoal.Funcionario
		if err := rows.Scan(&f.ID, &f.Nome, &f.CPF, &f.Cargo, &f.DataContratacao, &f.Salario, &f.Diaria, &f.Status); err != nil {
			return nil, fmt.Errorf("%s: erro ao ler linha: %w", op, err)
		}
		funcionarios = append(funcionarios, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre linhas: %w", op, err)
	}

	return funcionarios, nil
}
