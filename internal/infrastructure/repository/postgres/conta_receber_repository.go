package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
)

type ContaReceberRepositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NovoContaReceberRepositoryPostgres(dbpool *pgxpool.Pool) *ContaReceberRepositoryPostgres {
	return &ContaReceberRepositoryPostgres{dbpool: dbpool}
}

func (r *ContaReceberRepositoryPostgres) Salvar(ctx context.Context, dbtx db.DBTX, conta *financeiro.ContaReceber) error {
	const op = "repository.postgres.conta_receber.Salvar"

	query := `
		INSERT INTO contas_receber (
			id, obra_id, cronograma_recebimento_id, cliente, tipo_conta_receber,
			descricao, valor_original, valor_recebido, data_vencimento, 
			data_recebimento, status, forma_pagamento, observacoes, 
			numero_documento, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	// Se não há transação, usar o pool
	if dbtx == nil {
		dbtx = r.dbpool
	}

	_, err := dbtx.Exec(ctx, query,
		conta.ID,
		conta.ObraID,
		conta.CronogramaRecebimentoID,
		conta.Cliente,
		conta.TipoContaReceber,
		conta.Descricao,
		conta.ValorOriginal,
		conta.ValorRecebido,
		conta.DataVencimento,
		conta.DataRecebimento,
		conta.Status,
		conta.FormaPagamento,
		conta.Observacoes,
		conta.NumeroDocumento,
		conta.CreatedAt,
		conta.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ContaReceberRepositoryPostgres) Atualizar(ctx context.Context, conta *financeiro.ContaReceber) error {
	const op = "repository.postgres.conta_receber.Atualizar"

	query := `
		UPDATE contas_receber 
		SET cliente = $2,
			tipo_conta_receber = $3,
			descricao = $4,
			valor_original = $5,
			valor_recebido = $6,
			data_vencimento = $7,
			data_recebimento = $8,
			status = $9,
			forma_pagamento = $10,
			observacoes = $11,
			numero_documento = $12,
			updated_at = $13
		WHERE id = $1
	`

	result, err := r.dbpool.Exec(ctx, query,
		conta.ID,
		conta.Cliente,
		conta.TipoContaReceber,
		conta.Descricao,
		conta.ValorOriginal,
		conta.ValorRecebido,
		conta.DataVencimento,
		conta.DataRecebimento,
		conta.Status,
		conta.FormaPagamento,
		conta.Observacoes,
		conta.NumeroDocumento,
		conta.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: conta a receber não encontrada", op)
	}

	return nil
}

func (r *ContaReceberRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*financeiro.ContaReceber, error) {
	const op = "repository.postgres.conta_receber.BuscarPorID"

	query := `
		SELECT id, obra_id, cronograma_recebimento_id, cliente, tipo_conta_receber,
			   descricao, valor_original, valor_recebido, data_vencimento, 
			   data_recebimento, status, forma_pagamento, observacoes, 
			   numero_documento, created_at, updated_at
		FROM contas_receber 
		WHERE id = $1
	`

	row := r.dbpool.QueryRow(ctx, query, id)

	conta := &financeiro.ContaReceber{}
	err := row.Scan(
		&conta.ID,
		&conta.ObraID,
		&conta.CronogramaRecebimentoID,
		&conta.Cliente,
		&conta.TipoContaReceber,
		&conta.Descricao,
		&conta.ValorOriginal,
		&conta.ValorRecebido,
		&conta.DataVencimento,
		&conta.DataRecebimento,
		&conta.Status,
		&conta.FormaPagamento,
		&conta.Observacoes,
		&conta.NumeroDocumento,
		&conta.CreatedAt,
		&conta.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return conta, nil
}

func (r *ContaReceberRepositoryPostgres) ListarPorObraID(ctx context.Context, obraID string) ([]*financeiro.ContaReceber, error) {
	const op = "repository.postgres.conta_receber.ListarPorObraID"

	query := `
		SELECT id, obra_id, cronograma_recebimento_id, cliente, tipo_conta_receber,
			   descricao, valor_original, valor_recebido, data_vencimento, 
			   data_recebimento, status, forma_pagamento, observacoes, 
			   numero_documento, created_at, updated_at
		FROM contas_receber 
		WHERE obra_id = $1
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasReceber(ctx, rows, op)
}

func (r *ContaReceberRepositoryPostgres) ListarVencidas(ctx context.Context) ([]*financeiro.ContaReceber, error) {
	const op = "repository.postgres.conta_receber.ListarVencidas"

	query := `
		SELECT id, obra_id, cronograma_recebimento_id, cliente, tipo_conta_receber,
			   descricao, valor_original, valor_recebido, data_vencimento, 
			   data_recebimento, status, forma_pagamento, observacoes, 
			   numero_documento, created_at, updated_at
		FROM contas_receber 
		WHERE data_vencimento < CURRENT_DATE
		  AND status NOT IN ('RECEBIDO', 'CANCELADO')
		  AND valor_recebido < valor_original
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasReceber(ctx, rows, op)
}

func (r *ContaReceberRepositoryPostgres) ListarVencidasPorPeriodo(ctx context.Context, dataInicio, dataFim time.Time) ([]*financeiro.ContaReceber, error) {
	const op = "repository.postgres.conta_receber.ListarVencidasPorPeriodo"

	query := `
		SELECT id, obra_id, cronograma_recebimento_id, cliente, tipo_conta_receber,
			   descricao, valor_original, valor_recebido, data_vencimento, 
			   data_recebimento, status, forma_pagamento, observacoes, 
			   numero_documento, created_at, updated_at
		FROM contas_receber 
		WHERE data_vencimento BETWEEN $1 AND $2
		  AND status NOT IN ('RECEBIDO', 'CANCELADO')
		  AND valor_recebido < valor_original
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasReceber(ctx, rows, op)
}

func (r *ContaReceberRepositoryPostgres) ListarPorStatus(ctx context.Context, status string) ([]*financeiro.ContaReceber, error) {
	const op = "repository.postgres.conta_receber.ListarPorStatus"

	query := `
		SELECT id, obra_id, cronograma_recebimento_id, cliente, tipo_conta_receber,
			   descricao, valor_original, valor_recebido, data_vencimento, 
			   data_recebimento, status, forma_pagamento, observacoes, 
			   numero_documento, created_at, updated_at
		FROM contas_receber 
		WHERE status = $1
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasReceber(ctx, rows, op)
}

func (r *ContaReceberRepositoryPostgres) ListarPorCliente(ctx context.Context, cliente string) ([]*financeiro.ContaReceber, error) {
	const op = "repository.postgres.conta_receber.ListarPorCliente"

	query := `
		SELECT id, obra_id, cronograma_recebimento_id, cliente, tipo_conta_receber,
			   descricao, valor_original, valor_recebido, data_vencimento, 
			   data_recebimento, status, forma_pagamento, observacoes, 
			   numero_documento, created_at, updated_at
		FROM contas_receber 
		WHERE UPPER(cliente) LIKE UPPER($1)
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, "%"+cliente+"%")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasReceber(ctx, rows, op)
}

func (r *ContaReceberRepositoryPostgres) Listar(ctx context.Context, filtros common.ListarFiltros) ([]*financeiro.ContaReceber, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.conta_receber.Listar"

	// Query base
	baseQuery := `
		FROM contas_receber cr
		WHERE 1=1
	`

	// Adicionar filtros baseados nos parâmetros
	args := []interface{}{}
	whereClause := ""
	argCount := 0

	// Filtros por status
	if filtros.Status != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND cr.status = $%d", argCount)
		args = append(args, filtros.Status)
	}

	// Filtros por obra
	if filtros.ObraID != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND cr.obra_id = $%d", argCount)
		args = append(args, filtros.ObraID)
	}

	// Query para contar total
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var total int64
	err := r.dbpool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: erro ao contar registros: %w", op, err)
	}

	// Query para buscar dados
	dataQuery := `
		SELECT id, obra_id, cronograma_recebimento_id, cliente, tipo_conta_receber,
			   descricao, valor_original, valor_recebido, data_vencimento, 
			   data_recebimento, status, forma_pagamento, observacoes, 
			   numero_documento, created_at, updated_at
	` + baseQuery + whereClause + `
		ORDER BY cr.data_vencimento ASC
		LIMIT $` + fmt.Sprintf("%d", argCount+1) + ` OFFSET $` + fmt.Sprintf("%d", argCount+2)

	// Paginação
	limite := 50
	offset := 0
	if filtros.TamanhoPagina > 0 {
		limite = filtros.TamanhoPagina
	}
	if filtros.Pagina > 1 {
		offset = (filtros.Pagina - 1) * limite
	}
	args = append(args, limite, offset)

	rows, err := r.dbpool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	contas, err := r.scanContasReceber(ctx, rows, op)
	if err != nil {
		return nil, nil, err
	}

	// Informações de paginação
	paginaAtual := 1
	if filtros.Pagina > 0 {
		paginaAtual = filtros.Pagina
	}
	
	totalPaginas := 1
	if limite > 0 {
		totalPaginas = int((total + int64(limite) - 1) / int64(limite))
	}
	
	paginacao := &common.PaginacaoInfo{
		TotalItens:    int(total),
		PaginaAtual:   paginaAtual,
		TotalPaginas:  totalPaginas,
		TamanhoPagina: limite,
	}

	return contas, paginacao, nil
}

func (r *ContaReceberRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.conta_receber.Deletar"

	query := `DELETE FROM contas_receber WHERE id = $1`

	result, err := r.dbpool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: conta a receber não encontrada", op)
	}

	return nil
}

// scanContasReceber é um método auxiliar para escanear múltiplas contas a receber
func (r *ContaReceberRepositoryPostgres) scanContasReceber(ctx context.Context, rows pgx.Rows, op string) ([]*financeiro.ContaReceber, error) {
	var contas []*financeiro.ContaReceber

	for rows.Next() {
		conta := &financeiro.ContaReceber{}
		err := rows.Scan(
			&conta.ID,
			&conta.ObraID,
			&conta.CronogramaRecebimentoID,
			&conta.Cliente,
			&conta.TipoContaReceber,
			&conta.Descricao,
			&conta.ValorOriginal,
			&conta.ValorRecebido,
			&conta.DataVencimento,
			&conta.DataRecebimento,
			&conta.Status,
			&conta.FormaPagamento,
			&conta.Observacoes,
			&conta.NumeroDocumento,
			&conta.CreatedAt,
			&conta.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: erro ao escanear conta a receber: %w", op, err)
		}
		contas = append(contas, conta)
	}

	return contas, nil
}