// file: internal/infrastructure/repository/postgres/apontamento_repository.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
)

type ApontamentoRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoApontamentoRepository(db *pgxpool.Pool, logger *slog.Logger) *ApontamentoRepositoryPostgres {
	return &ApontamentoRepositoryPostgres{db: db, logger: logger}
}

func (r *ApontamentoRepositoryPostgres) Salvar(ctx context.Context, a *pessoal.ApontamentoQuinzenal) error {
	const op = "repository.postgres.apontamento.Salvar"
	query := `
		INSERT INTO apontamentos_quinzenais (
			id, funcionario_id, obra_id, periodo_inicio, periodo_fim,
			dias_trabalhados, adicionais, descontos, adiantamentos,
			valor_total_calculado, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		a.ID, a.FuncionarioID, a.ObraID, a.PeriodoInicio, a.PeriodoFim,
		a.DiasTrabalhados, a.Adicionais, a.Descontos, a.Adiantamentos,
		a.ValorTotalCalculado, a.Status, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		// TODO: Tratar erro de violação da constraint UNIQUE(funcionario_id, periodo_inicio, periodo_fim)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *ApontamentoRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*pessoal.ApontamentoQuinzenal, error) {
	const op = "repository.postgres.apontamento.BuscarPorID"
	query := `
		SELECT id, funcionario_id, obra_id, periodo_inicio, periodo_fim,
			   dias_trabalhados, adicionais, descontos, adiantamentos,
			   valor_total_calculado, status, created_at, updated_at
		FROM apontamentos_quinzenais WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)
	var a pessoal.ApontamentoQuinzenal

	err := row.Scan(
		&a.ID, &a.FuncionarioID, &a.ObraID, &a.PeriodoInicio, &a.PeriodoFim,
		&a.DiasTrabalhados, &a.Adicionais, &a.Descontos, &a.Adiantamentos,
		&a.ValorTotalCalculado, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &a, nil
}

func (r *ApontamentoRepositoryPostgres) Atualizar(ctx context.Context, a *pessoal.ApontamentoQuinzenal) error {
	const op = "repository.postgres.apontamento.Atualizar"
	query := `
		UPDATE apontamentos_quinzenais SET
			dias_trabalhados = $1, adicionais = $2, descontos = $3, adiantamentos = $4,
			valor_total_calculado = $5, status = $6, updated_at = $7
		WHERE id = $8`

	cmd, err := r.db.Exec(ctx, query,
		a.DiasTrabalhados, a.Adicionais, a.Descontos, a.Adiantamentos,
		a.ValorTotalCalculado, a.Status, a.UpdatedAt, a.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
