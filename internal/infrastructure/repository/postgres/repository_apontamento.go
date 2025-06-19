// file: internal/infrastructure/repository/postgres/apontamento_repository.go
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

func (r *ApontamentoRepositoryPostgres) Listar(ctx context.Context, filtros common.ListarFiltros) ([]*pessoal.ApontamentoQuinzenal, *common.PaginacaoInfo, error) {
	// A query base para buscar todos os apontamentos
	baseQuery := "FROM apontamentos_quinzenais"
	// Para os filtros, passamos um mapa que será preenchido
	filterArgs := make(map[string]interface{})
	if filtros.Status != "" {
		filterArgs["status"] = filtros.Status
	}

	return r.listarComFiltros(ctx, baseQuery, filterArgs, filtros)
}

func (r *ApontamentoRepositoryPostgres) ListarPorFuncionarioID(ctx context.Context, funcionarioID string, filtros common.ListarFiltros) ([]*pessoal.ApontamentoQuinzenal, *common.PaginacaoInfo, error) {
	// A query base agora filtra por funcionário
	baseQuery := "FROM apontamentos_quinzenais WHERE funcionario_id = @funcionarioID"
	filterArgs := map[string]interface{}{"funcionarioID": funcionarioID}
	if filtros.Status != "" {
		filterArgs["status"] = filtros.Status
	}

	return r.listarComFiltros(ctx, baseQuery, filterArgs, filtros)
}

// listarComFiltros é uma função helper interna para não duplicar a lógica de paginação.
func (r *ApontamentoRepositoryPostgres) listarComFiltros(ctx context.Context, baseQuery string, filterArgs map[string]interface{}, filtros common.ListarFiltros) ([]*pessoal.ApontamentoQuinzenal, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.apontamento.listarComFiltros"

	args := pgx.NamedArgs(filterArgs)

	countQueryBuilder := strings.Builder{}
	countQueryBuilder.WriteString("SELECT COUNT(*) ")
	countQueryBuilder.WriteString(baseQuery)

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("SELECT id, funcionario_id, obra_id, periodo_inicio, periodo_fim, dias_trabalhados, adicionais, descontos, adiantamentos, valor_total_calculado, status, created_at, updated_at ")
	queryBuilder.WriteString(baseQuery)

	if status, ok := filterArgs["status"]; ok {
		countQueryBuilder.WriteString(" AND status = @status")
		queryBuilder.WriteString(" AND status = @status")
		args["status"] = status
	}

	var totalItens int
	err := r.db.QueryRow(ctx, countQueryBuilder.String(), args).Scan(&totalItens)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao contar apontamentos: %w", op, err)
	}

	paginacao := common.NewPaginacaoInfo(totalItens, filtros.Pagina, filtros.TamanhoPagina)
	if totalItens == 0 {
		return []*pessoal.ApontamentoQuinzenal{}, paginacao, nil
	}

	offset := (filtros.Pagina - 1) * filtros.TamanhoPagina
	queryBuilder.WriteString(" ORDER BY periodo_inicio DESC, created_at DESC LIMIT @limit OFFSET @offset")
	args["limit"] = filtros.TamanhoPagina
	args["offset"] = offset

	rows, err := r.db.Query(ctx, queryBuilder.String(), args)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	apontamentos, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[pessoal.ApontamentoQuinzenal])
	if err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao escanear apontamentos: %w", op, err)
	}

	return apontamentos, paginacao, nil
}
