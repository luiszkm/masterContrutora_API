package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

type CronogramaRecebimentoRepositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NovoCronogramaRecebimentoRepositoryPostgres(dbpool *pgxpool.Pool) *CronogramaRecebimentoRepositoryPostgres {
	return &CronogramaRecebimentoRepositoryPostgres{dbpool: dbpool}
}

func (r *CronogramaRecebimentoRepositoryPostgres) Salvar(ctx context.Context, dbtx db.DBTX, cronograma *obras.CronogramaRecebimento) error {
	const op = "repository.postgres.cronograma_recebimento.Salvar"

	query := `
		INSERT INTO cronograma_recebimentos (
			id, obra_id, numero_etapa, descricao_etapa, valor_previsto, 
			data_vencimento, status, data_recebimento, valor_recebido, 
			observacoes_recebimento, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := dbtx.Exec(ctx, query,
		cronograma.ID,
		cronograma.ObraID,
		cronograma.NumeroEtapa,
		cronograma.DescricaoEtapa,
		cronograma.ValorPrevisto,
		cronograma.DataVencimento,
		cronograma.Status,
		cronograma.DataRecebimento,
		cronograma.ValorRecebido,
		cronograma.ObservacoesRecebimento,
		cronograma.CreatedAt,
		cronograma.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *CronogramaRecebimentoRepositoryPostgres) SalvarMuitos(ctx context.Context, dbtx db.DBTX, cronogramas []*obras.CronogramaRecebimento) error {
	const op = "repository.postgres.cronograma_recebimento.SalvarMuitos"

	if len(cronogramas) == 0 {
		return nil
	}

	query := `
		INSERT INTO cronograma_recebimentos (
			id, obra_id, numero_etapa, descricao_etapa, valor_previsto, 
			data_vencimento, status, data_recebimento, valor_recebido, 
			observacoes_recebimento, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	// Execute cada insert individualmente
	for i, cronograma := range cronogramas {
		_, err := dbtx.Exec(ctx, query,
			cronograma.ID,
			cronograma.ObraID,
			cronograma.NumeroEtapa,
			cronograma.DescricaoEtapa,
			cronograma.ValorPrevisto,
			cronograma.DataVencimento,
			cronograma.Status,
			cronograma.DataRecebimento,
			cronograma.ValorRecebido,
			cronograma.ObservacoesRecebimento,
			cronograma.CreatedAt,
			cronograma.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("%s: falha na inserção do cronograma %d: %w", op, i, err)
		}
	}

	return nil
}

func (r *CronogramaRecebimentoRepositoryPostgres) Atualizar(ctx context.Context, cronograma *obras.CronogramaRecebimento) error {
	const op = "repository.postgres.cronograma_recebimento.Atualizar"

	query := `
		UPDATE cronograma_recebimentos 
		SET descricao_etapa = $2,
			valor_previsto = $3,
			data_vencimento = $4,
			status = $5,
			data_recebimento = $6,
			valor_recebido = $7,
			observacoes_recebimento = $8,
			updated_at = $9
		WHERE id = $1
	`

	result, err := r.dbpool.Exec(ctx, query,
		cronograma.ID,
		cronograma.DescricaoEtapa,
		cronograma.ValorPrevisto,
		cronograma.DataVencimento,
		cronograma.Status,
		cronograma.DataRecebimento,
		cronograma.ValorRecebido,
		cronograma.ObservacoesRecebimento,
		cronograma.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: cronograma não encontrado", op)
	}

	return nil
}

func (r *CronogramaRecebimentoRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*obras.CronogramaRecebimento, error) {
	const op = "repository.postgres.cronograma_recebimento.BuscarPorID"

	query := `
		SELECT id, obra_id, numero_etapa, descricao_etapa, valor_previsto, 
			   data_vencimento, status, data_recebimento, valor_recebido, 
			   observacoes_recebimento, created_at, updated_at
		FROM cronograma_recebimentos 
		WHERE id = $1
	`

	row := r.dbpool.QueryRow(ctx, query, id)

	cronograma := &obras.CronogramaRecebimento{}
	err := row.Scan(
		&cronograma.ID,
		&cronograma.ObraID,
		&cronograma.NumeroEtapa,
		&cronograma.DescricaoEtapa,
		&cronograma.ValorPrevisto,
		&cronograma.DataVencimento,
		&cronograma.Status,
		&cronograma.DataRecebimento,
		&cronograma.ValorRecebido,
		&cronograma.ObservacoesRecebimento,
		&cronograma.CreatedAt,
		&cronograma.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cronograma, nil
}

func (r *CronogramaRecebimentoRepositoryPostgres) ListarPorObraID(ctx context.Context, obraID string) ([]*obras.CronogramaRecebimento, error) {
	const op = "repository.postgres.cronograma_recebimento.ListarPorObraID"

	query := `
		SELECT id, obra_id, numero_etapa, descricao_etapa, valor_previsto, 
			   data_vencimento, status, data_recebimento, valor_recebido, 
			   observacoes_recebimento, created_at, updated_at
		FROM cronograma_recebimentos 
		WHERE obra_id = $1
		ORDER BY numero_etapa ASC
	`

	rows, err := r.dbpool.Query(ctx, query, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var cronogramas []*obras.CronogramaRecebimento
	for rows.Next() {
		cronograma := &obras.CronogramaRecebimento{}
		err := rows.Scan(
			&cronograma.ID,
			&cronograma.ObraID,
			&cronograma.NumeroEtapa,
			&cronograma.DescricaoEtapa,
			&cronograma.ValorPrevisto,
			&cronograma.DataVencimento,
			&cronograma.Status,
			&cronograma.DataRecebimento,
			&cronograma.ValorRecebido,
			&cronograma.ObservacoesRecebimento,
			&cronograma.CreatedAt,
			&cronograma.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: erro ao escanear cronograma: %w", op, err)
		}
		cronogramas = append(cronogramas, cronograma)
	}

	return cronogramas, nil
}

func (r *CronogramaRecebimentoRepositoryPostgres) ListarVencidosPorPeriodo(ctx context.Context, dataInicio, dataFim time.Time) ([]*obras.CronogramaRecebimento, error) {
	const op = "repository.postgres.cronograma_recebimento.ListarVencidosPorPeriodo"

	query := `
		SELECT cr.id, cr.obra_id, cr.numero_etapa, cr.descricao_etapa, cr.valor_previsto, 
			   cr.data_vencimento, cr.status, cr.data_recebimento, cr.valor_recebido, 
			   cr.observacoes_recebimento, cr.created_at, cr.updated_at
		FROM cronograma_recebimentos cr
		WHERE cr.data_vencimento BETWEEN $1 AND $2
		  AND cr.status IN ('PENDENTE', 'VENCIDO', 'PARCIAL')
		  AND cr.valor_recebido < cr.valor_previsto
		ORDER BY cr.data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var cronogramas []*obras.CronogramaRecebimento
	for rows.Next() {
		cronograma := &obras.CronogramaRecebimento{}
		err := rows.Scan(
			&cronograma.ID,
			&cronograma.ObraID,
			&cronograma.NumeroEtapa,
			&cronograma.DescricaoEtapa,
			&cronograma.ValorPrevisto,
			&cronograma.DataVencimento,
			&cronograma.Status,
			&cronograma.DataRecebimento,
			&cronograma.ValorRecebido,
			&cronograma.ObservacoesRecebimento,
			&cronograma.CreatedAt,
			&cronograma.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: erro ao escanear cronograma vencido: %w", op, err)
		}
		cronogramas = append(cronogramas, cronograma)
	}

	return cronogramas, nil
}

func (r *CronogramaRecebimentoRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.cronograma_recebimento.Deletar"

	query := `DELETE FROM cronograma_recebimentos WHERE id = $1`

	result, err := r.dbpool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: cronograma não encontrado", op)
	}

	return nil
}