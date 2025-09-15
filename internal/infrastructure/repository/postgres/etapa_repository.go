// file: internal/repository/postgres/etapa_repository.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

// EtapaRepositoryPostgres implementa a persistência para o agregado Etapa.
type EtapaRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoEtapaRepository(db *pgxpool.Pool, logger *slog.Logger) obras.EtapaRepository {
	return &EtapaRepositoryPostgres{
		db:     db,
		logger: logger,
	}
}

func (r *EtapaRepositoryPostgres) Salvar(ctx context.Context, dbtx db.DBTX, etapa *obras.Etapa) error {
	const op = "repository.postgres.etapa.Salvar"
	query := `
		INSERT INTO etapas (id, obra_id, nome, ordem, data_inicio_prevista, data_fim_prevista, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := dbtx.Exec(ctx, query,
		etapa.ID, etapa.ObraID, etapa.Nome, etapa.Ordem,
		etapa.DataInicioPrevista, etapa.DataFimPrevista, etapa.Status,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *EtapaRepositoryPostgres) BuscarPorID(ctx context.Context, etapaID string) (*obras.Etapa, error) {
	const op = "repository.postgres.etapa.BuscarPorID"
	query := `SELECT id, obra_id, nome, ordem, data_inicio_prevista, data_fim_prevista, status FROM etapas WHERE id = $1`
	row := r.db.QueryRow(ctx, query, etapaID)

	var etapa obras.Etapa
	err := row.Scan(
		&etapa.ID,
		&etapa.ObraID,
		&etapa.Nome,
		&etapa.Ordem,
		&etapa.DataInicioPrevista,
		&etapa.DataFimPrevista,
		&etapa.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &etapa, nil
}

func (r *EtapaRepositoryPostgres) Atualizar(ctx context.Context, etapa *obras.Etapa) error {
	const op = "repository.postgres.etapa.Atualizar"
	query := `UPDATE etapas SET nome = $1, data_inicio_prevista = $2, data_fim_prevista = $3, status = $4 WHERE id = $5`

	cmd, err := r.db.Exec(ctx, query,
		etapa.Nome,
		etapa.DataInicioPrevista,
		etapa.DataFimPrevista,
		etapa.Status,
		etapa.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	// Verifica se alguma linha foi realmente atualizada
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}

	return nil
}

func (r *EtapaRepositoryPostgres) ListarPorObraID(ctx context.Context, obraID string) ([]*obras.Etapa, error) {
	const op = "repository.postgres.etapa.ListarPorObraID"
	query := `
		SELECT id, obra_id, nome, ordem, data_inicio_prevista, data_fim_prevista, status
		FROM etapas
		WHERE obra_id = $1
		ORDER BY ordem ASC, nome ASC
	`
	rows, err := r.db.Query(ctx, query, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	etapas, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[obras.Etapa])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*obras.Etapa{}, nil // Retorna lista vazia se não houver etapas
		}
		return nil, fmt.Errorf("%s: falha ao escanear etapas: %w", op, err)
	}
	return etapas, nil
}

func (r *EtapaRepositoryPostgres) ListarPorObraIDPaginado(ctx context.Context, obraID string, filtros common.ListarFiltros) ([]*obras.Etapa, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.etapa.ListarPorObraIDPaginado"

	// Definir valores padrão de paginação
	pagina := filtros.Pagina
	tamanhoPagina := filtros.TamanhoPagina
	if pagina <= 0 {
		pagina = 1
	}
	if tamanhoPagina <= 0 {
		tamanhoPagina = 10
	}

	// Construir a cláusula WHERE baseada nos filtros
	whereConditions := []string{"obra_id = $1"}
	args := []interface{}{obraID}
	argIndex := 2

	if filtros.Status != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filtros.Status)
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Query para contar o total de registros
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM etapas 
		WHERE %s`, whereClause)

	var totalItens int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItens)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao contar etapas: %w", op, err)
	}

	// Query principal com paginação
	offset := (pagina - 1) * tamanhoPagina
	query := fmt.Sprintf(`
		SELECT id, obra_id, nome, ordem, data_inicio_prevista, data_fim_prevista, status
		FROM etapas
		WHERE %s
		ORDER BY ordem ASC, nome ASC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, tamanhoPagina, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	etapas, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[obras.Etapa])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			paginacaoInfo := common.NewPaginacaoInfo(0, pagina, tamanhoPagina)
			return []*obras.Etapa{}, paginacaoInfo, nil
		}
		return nil, nil, fmt.Errorf("%s: falha ao escanear etapas: %w", op, err)
	}

	// Criar info de paginação
	paginacaoInfo := common.NewPaginacaoInfo(totalItens, pagina, tamanhoPagina)

	return etapas, paginacaoInfo, nil
}
