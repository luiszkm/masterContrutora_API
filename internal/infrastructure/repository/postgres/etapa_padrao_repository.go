// file: internal/infrastructure/repository/postgres/etapa_padrao_repository.go
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
)

type EtapaPadraoRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoEtapaPadraoRepository(db *pgxpool.Pool, logger *slog.Logger) *EtapaPadraoRepositoryPostgres {
	return &EtapaPadraoRepositoryPostgres{db: db, logger: logger}
}

func (r *EtapaPadraoRepositoryPostgres) Salvar(ctx context.Context, e *obras.EtapaPadrao) error {
	query := `INSERT INTO etapas_padrao (id, nome, descricao, ordem, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, e.ID, e.Nome, e.Descricao, e.Ordem, e.CreatedAt, e.UpdatedAt)
	return err
}

func (r *EtapaPadraoRepositoryPostgres) ListarTodas(ctx context.Context) ([]*obras.EtapaPadrao, error) {
	query := `SELECT id, nome, descricao, ordem, created_at, updated_at FROM etapas_padrao ORDER BY ordem, nome ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[obras.EtapaPadrao])
}

func (r *EtapaPadraoRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*obras.EtapaPadrao, error) {
	query := `SELECT id, nome, descricao, ordem, created_at, updated_at FROM etapas_padrao WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var etapa obras.EtapaPadrao
	err := row.Scan(&etapa.ID, &etapa.Nome, &etapa.Descricao, &etapa.Ordem, &etapa.CreatedAt, &etapa.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &etapa, nil
}

func (r *EtapaPadraoRepositoryPostgres) Atualizar(ctx context.Context, e *obras.EtapaPadrao) error {
	query := `UPDATE etapas_padrao SET nome = $1, descricao = $2, ordem = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx, query, e.Nome, e.Descricao, e.Ordem, e.UpdatedAt, e.ID)
	return err
}

func (r *EtapaPadraoRepositoryPostgres) Listar(ctx context.Context, filtros common.ListarFiltros) ([]*obras.EtapaPadrao, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.etapa_padrao.Listar"

	// Definir valores padrão de paginação
	pagina := filtros.Pagina
	tamanhoPagina := filtros.TamanhoPagina
	if pagina <= 0 {
		pagina = 1
	}
	if tamanhoPagina <= 0 {
		tamanhoPagina = 10
	}

	// Para etapas padrão, não há muitos filtros específicos
	whereConditions := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1
	
	whereClause := strings.Join(whereConditions, " AND ")

	// Query para contar o total de registros
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM etapas_padrao 
		WHERE %s`, whereClause)

	var totalItens int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItens)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao contar etapas padrão: %w", op, err)
	}

	// Query principal com paginação
	offset := (pagina - 1) * tamanhoPagina
	query := fmt.Sprintf(`
		SELECT id, nome, descricao, ordem, created_at, updated_at 
		FROM etapas_padrao 
		WHERE %s
		ORDER BY ordem, nome ASC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, tamanhoPagina, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	etapas, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[obras.EtapaPadrao])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			paginacaoInfo := common.NewPaginacaoInfo(0, pagina, tamanhoPagina)
			return []*obras.EtapaPadrao{}, paginacaoInfo, nil
		}
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	// Criar info de paginação
	paginacaoInfo := common.NewPaginacaoInfo(totalItens, pagina, tamanhoPagina)

	return etapas, paginacaoInfo, nil
}

func (r *EtapaPadraoRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	query := `DELETE FROM etapas_padrao WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
