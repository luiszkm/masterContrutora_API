// file: internal/repository/postgres/alocacao_repository.go
package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
)

type AlocacaoRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoAlocacaoRepository(db *pgxpool.Pool, logger *slog.Logger) obras.AlocacaoRepository {
	return &AlocacaoRepositoryPostgres{db: db, logger: logger}
}

func (r *AlocacaoRepositoryPostgres) Salvar(ctx context.Context, a *obras.Alocacao) error {
	const op = "repository.postgres.alocacao.Salvar"
	query := `
		INSERT INTO alocacoes (id, obra_id, funcionario_id, data_inicio_alocacao, data_fim_alocacao)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, a.ID, a.ObraID, a.FuncionarioID, a.DataInicioAlocacao, a.DataFimAlocacao)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// SalvarMuitos insere múltiplas alocações de forma eficiente usando um batch.
func (r *AlocacaoRepositoryPostgres) SalvarMuitos(ctx context.Context, alocacoes []*obras.Alocacao) error {
	const op = "repository.postgres.alocacao.SalvarMuitos"

	batch := &pgx.Batch{}
	query := `
		INSERT INTO alocacoes (id, obra_id, funcionario_id, data_inicio_alocacao, data_fim_alocacao)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, a := range alocacoes {
		batch.Queue(query, a.ID, a.ObraID, a.FuncionarioID, a.DataInicioAlocacao, a.DataFimAlocacao)
	}

	br := r.db.SendBatch(ctx, batch)
	// É crucial fechar o batch result para verificar por erros.
	if err := br.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *AlocacaoRepositoryPostgres) ExistemAlocacoesAtivasParaFuncionario(ctx context.Context, funcionarioID string) (bool, error) {
	const op = "repository.postgres.alocacao.ExistemAlocacoesAtivasParaFuncionario"
	query := `SELECT EXISTS(SELECT 1 FROM alocacoes WHERE funcionario_id = $1 AND (data_fim_alocacao IS NULL OR data_fim_alocacao >= CURRENT_DATE))`

	var existe bool
	err := r.db.QueryRow(ctx, query, funcionarioID).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return existe, nil
}
