// file: internal/repository/postgres/registro_pagamento_repository.go
package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

// RegistroPagamentoRepositoryPostgres implementa a interface financeiro.Repository.
type RegistroPagamentoRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoRegistroPagamentoRepository(db *pgxpool.Pool, logger *slog.Logger) *RegistroPagamentoRepositoryPostgres {
	return &RegistroPagamentoRepositoryPostgres{db: db, logger: logger}
}

func (r *RegistroPagamentoRepositoryPostgres) Salvar(ctx context.Context, dbtx db.DBTX, p *financeiro.RegistroDePagamento) error {
	const op = "repository.postgres.pagamento.Salvar"
	query := `
		INSERT INTO registros_pagamento (id, funcionario_id, obra_id, periodo_referencia, valor_calculado, data_de_efetivacao, conta_bancaria_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := dbtx.Exec(ctx, query,
		p.ID,
		p.FuncionarioID,
		p.ObraID,
		p.PeriodoReferencia,
		p.ValorCalculado,
		p.DataDeEfetivacao,
		p.ContaBancariaID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
