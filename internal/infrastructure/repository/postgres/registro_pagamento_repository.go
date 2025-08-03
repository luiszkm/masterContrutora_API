// file: internal/repository/postgres/registro_pagamento_repository.go
package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
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

func (r *RegistroPagamentoRepositoryPostgres) ListarPagamentos(ctx context.Context, filtros common.ListarFiltros) ([]*financeiro.RegistroDePagamento, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.pagamento.ListarPagamentos"

	args := pgx.NamedArgs{}
	whereClauses := []string{}

	if filtros.FuncionarioID != "" {
		whereClauses = append(whereClauses, "rp.funcionario_id = @funcionarioID")
		args["funcionarioID"] = filtros.FuncionarioID
	}
	if filtros.ObraID != "" {
		whereClauses = append(whereClauses, "rp.obra_id = @obraID")
		args["obraID"] = filtros.ObraID
	}

	whereString := ""
	if len(whereClauses) > 0 {
		whereString = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Query de contagem
	countQuery := "SELECT COUNT(*) FROM registros_pagamento rp" + whereString
	var totalItens int
	if err := r.db.QueryRow(ctx, countQuery, args).Scan(&totalItens); err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao contar pagamentos: %w", op, err)
	}

	paginacao := common.NewPaginacaoInfo(totalItens, filtros.Pagina, filtros.TamanhoPagina)
	if totalItens == 0 {
		return []*financeiro.RegistroDePagamento{}, paginacao, nil
	}

	// Query principal
	query := `
		SELECT 
			rp.id, rp.funcionario_id, rp.obra_id, rp.periodo_referencia, 
			rp.valor_calculado, rp.data_de_efetivacao, rp.conta_bancaria_id
		FROM registros_pagamento rp` + whereString + `
		ORDER BY rp.data_de_efetivacao DESC
		LIMIT @limit OFFSET @offset`

	args["limit"] = filtros.TamanhoPagina
	offset := (filtros.Pagina - 1) * filtros.TamanhoPagina
	args["offset"] = offset

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	pagamentos, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[financeiro.RegistroDePagamento])
	if err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao escanear pagamentos: %w", op, err)
	}

	return pagamentos, paginacao, nil
}
