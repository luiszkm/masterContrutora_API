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

type ContaPagarRepositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NovoContaPagarRepositoryPostgres(dbpool *pgxpool.Pool) *ContaPagarRepositoryPostgres {
	return &ContaPagarRepositoryPostgres{dbpool: dbpool}
}

func (r *ContaPagarRepositoryPostgres) Salvar(ctx context.Context, dbtx db.DBTX, conta *financeiro.ContaPagar) error {
	const op = "repository.postgres.conta_pagar.Salvar"

	query := `
		INSERT INTO contas_pagar (
			id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			categoria, descricao, valor_original, valor_pago, data_vencimento, 
			data_pagamento, status, forma_pagamento, observacoes, 
			numero_documento, numero_compra_nf, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	// Se não há transação, usar o pool
	if dbtx == nil {
		dbtx = r.dbpool
	}

	_, err := dbtx.Exec(ctx, query,
		conta.ID,
		conta.FornecedorID,
		conta.ObraID,
		conta.OrcamentoID,
		conta.FornecedorNome,
		conta.TipoContaPagar,
		conta.Categoria,
		conta.Descricao,
		conta.ValorOriginal,
		conta.ValorPago,
		conta.DataVencimento,
		conta.DataPagamento,
		conta.Status,
		conta.FormaPagamento,
		conta.Observacoes,
		conta.NumeroDocumento,
		conta.NumeroCompraNF,
		conta.CreatedAt,
		conta.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ContaPagarRepositoryPostgres) Atualizar(ctx context.Context, conta *financeiro.ContaPagar) error {
	const op = "repository.postgres.conta_pagar.Atualizar"

	query := `
		UPDATE contas_pagar 
		SET fornecedor_nome = $2,
			tipo_conta_pagar = $3,
			categoria = $4,
			descricao = $5,
			valor_original = $6,
			valor_pago = $7,
			data_vencimento = $8,
			data_pagamento = $9,
			status = $10,
			forma_pagamento = $11,
			observacoes = $12,
			numero_documento = $13,
			numero_compra_nf = $14,
			updated_at = $15
		WHERE id = $1
	`

	result, err := r.dbpool.Exec(ctx, query,
		conta.ID,
		conta.FornecedorNome,
		conta.TipoContaPagar,
		conta.Categoria,
		conta.Descricao,
		conta.ValorOriginal,
		conta.ValorPago,
		conta.DataVencimento,
		conta.DataPagamento,
		conta.Status,
		conta.FormaPagamento,
		conta.Observacoes,
		conta.NumeroDocumento,
		conta.NumeroCompraNF,
		conta.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: conta a pagar não encontrada", op)
	}

	return nil
}

func (r *ContaPagarRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*financeiro.ContaPagar, error) {
	const op = "repository.postgres.conta_pagar.BuscarPorID"

	query := `
		SELECT id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			   categoria, descricao, valor_original, valor_pago, data_vencimento, 
			   data_pagamento, status, forma_pagamento, observacoes, 
			   numero_documento, numero_compra_nf, created_at, updated_at
		FROM contas_pagar 
		WHERE id = $1
	`

	row := r.dbpool.QueryRow(ctx, query, id)

	conta := &financeiro.ContaPagar{}
	err := row.Scan(
		&conta.ID,
		&conta.FornecedorID,
		&conta.ObraID,
		&conta.OrcamentoID,
		&conta.FornecedorNome,
		&conta.TipoContaPagar,
		&conta.Categoria,
		&conta.Descricao,
		&conta.ValorOriginal,
		&conta.ValorPago,
		&conta.DataVencimento,
		&conta.DataPagamento,
		&conta.Status,
		&conta.FormaPagamento,
		&conta.Observacoes,
		&conta.NumeroDocumento,
		&conta.NumeroCompraNF,
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

func (r *ContaPagarRepositoryPostgres) ListarPorObraID(ctx context.Context, obraID string) ([]*financeiro.ContaPagar, error) {
	const op = "repository.postgres.conta_pagar.ListarPorObraID"

	query := `
		SELECT id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			   descricao, valor_original, valor_pago, data_vencimento, 
			   data_pagamento, status, forma_pagamento, observacoes, 
			   numero_documento, numero_compra_nf, created_at, updated_at
		FROM contas_pagar 
		WHERE obra_id = $1
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, obraID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasPagar(ctx, rows, op)
}

func (r *ContaPagarRepositoryPostgres) ListarPorFornecedorID(ctx context.Context, fornecedorID string) ([]*financeiro.ContaPagar, error) {
	const op = "repository.postgres.conta_pagar.ListarPorFornecedorID"

	query := `
		SELECT id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			   descricao, valor_original, valor_pago, data_vencimento, 
			   data_pagamento, status, forma_pagamento, observacoes, 
			   numero_documento, numero_compra_nf, created_at, updated_at
		FROM contas_pagar 
		WHERE fornecedor_id = $1
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, fornecedorID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasPagar(ctx, rows, op)
}

func (r *ContaPagarRepositoryPostgres) ListarPorOrcamentoID(ctx context.Context, orcamentoID string) ([]*financeiro.ContaPagar, error) {
	const op = "repository.postgres.conta_pagar.ListarPorOrcamentoID"

	query := `
		SELECT id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			   descricao, valor_original, valor_pago, data_vencimento, 
			   data_pagamento, status, forma_pagamento, observacoes, 
			   numero_documento, numero_compra_nf, created_at, updated_at
		FROM contas_pagar 
		WHERE orcamento_id = $1
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, orcamentoID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasPagar(ctx, rows, op)
}

func (r *ContaPagarRepositoryPostgres) ListarVencidas(ctx context.Context) ([]*financeiro.ContaPagar, error) {
	const op = "repository.postgres.conta_pagar.ListarVencidas"

	query := `
		SELECT id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			   descricao, valor_original, valor_pago, data_vencimento, 
			   data_pagamento, status, forma_pagamento, observacoes, 
			   numero_documento, numero_compra_nf, created_at, updated_at
		FROM contas_pagar 
		WHERE data_vencimento < CURRENT_DATE
		  AND status NOT IN ('PAGO', 'CANCELADO')
		  AND valor_pago < valor_original
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasPagar(ctx, rows, op)
}

func (r *ContaPagarRepositoryPostgres) ListarVencidasPorPeriodo(ctx context.Context, dataInicio, dataFim time.Time) ([]*financeiro.ContaPagar, error) {
	const op = "repository.postgres.conta_pagar.ListarVencidasPorPeriodo"

	query := `
		SELECT id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			   descricao, valor_original, valor_pago, data_vencimento, 
			   data_pagamento, status, forma_pagamento, observacoes, 
			   numero_documento, numero_compra_nf, created_at, updated_at
		FROM contas_pagar 
		WHERE data_vencimento BETWEEN $1 AND $2
		  AND status NOT IN ('PAGO', 'CANCELADO')
		  AND valor_pago < valor_original
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasPagar(ctx, rows, op)
}

func (r *ContaPagarRepositoryPostgres) ListarPorStatus(ctx context.Context, status string) ([]*financeiro.ContaPagar, error) {
	const op = "repository.postgres.conta_pagar.ListarPorStatus"

	query := `
		SELECT id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			   descricao, valor_original, valor_pago, data_vencimento, 
			   data_pagamento, status, forma_pagamento, observacoes, 
			   numero_documento, numero_compra_nf, created_at, updated_at
		FROM contas_pagar 
		WHERE status = $1
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasPagar(ctx, rows, op)
}

func (r *ContaPagarRepositoryPostgres) ListarPorFornecedor(ctx context.Context, fornecedorNome string) ([]*financeiro.ContaPagar, error) {
	const op = "repository.postgres.conta_pagar.ListarPorFornecedor"

	query := `
		SELECT id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			   descricao, valor_original, valor_pago, data_vencimento, 
			   data_pagamento, status, forma_pagamento, observacoes, 
			   numero_documento, numero_compra_nf, created_at, updated_at
		FROM contas_pagar 
		WHERE UPPER(fornecedor_nome) LIKE UPPER($1)
		ORDER BY data_vencimento ASC
	`

	rows, err := r.dbpool.Query(ctx, query, "%"+fornecedorNome+"%")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return r.scanContasPagar(ctx, rows, op)
}

func (r *ContaPagarRepositoryPostgres) Listar(ctx context.Context, filtros common.ListarFiltros) ([]*financeiro.ContaPagar, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.conta_pagar.Listar"

	// Query base
	baseQuery := `
		FROM contas_pagar cp
		WHERE 1=1
	`

	args := []interface{}{}
	whereClause := ""
	argCount := 0

	// Filtros por status
	if filtros.Status != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND cp.status = $%d", argCount)
		args = append(args, filtros.Status)
	}

	// Filtros por obra
	if filtros.ObraID != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND cp.obra_id = $%d", argCount)
		args = append(args, filtros.ObraID)
	}

	// Filtros por fornecedor
	if filtros.FornecedorID != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND cp.fornecedor_id = $%d", argCount)
		args = append(args, filtros.FornecedorID)
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
		SELECT id, fornecedor_id, obra_id, orcamento_id, fornecedor_nome, tipo_conta_pagar,
			   descricao, valor_original, valor_pago, data_vencimento, 
			   data_pagamento, status, forma_pagamento, observacoes, 
			   numero_documento, numero_compra_nf, created_at, updated_at
	` + baseQuery + whereClause + `
		ORDER BY cp.data_vencimento ASC
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

	contas, err := r.scanContasPagar(ctx, rows, op)
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

func (r *ContaPagarRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.conta_pagar.Deletar"

	query := `DELETE FROM contas_pagar WHERE id = $1`

	result, err := r.dbpool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: conta a pagar não encontrada", op)
	}

	return nil
}

// scanContasPagar é um método auxiliar para escanear múltiplas contas a pagar
func (r *ContaPagarRepositoryPostgres) scanContasPagar(ctx context.Context, rows pgx.Rows, op string) ([]*financeiro.ContaPagar, error) {
	var contas []*financeiro.ContaPagar

	for rows.Next() {
		conta := &financeiro.ContaPagar{}
		err := rows.Scan(
			&conta.ID,
			&conta.FornecedorID,
			&conta.ObraID,
			&conta.OrcamentoID,
			&conta.FornecedorNome,
			&conta.TipoContaPagar,
			&conta.Descricao,
			&conta.ValorOriginal,
			&conta.ValorPago,
			&conta.DataVencimento,
			&conta.DataPagamento,
			&conta.Status,
			&conta.FormaPagamento,
			&conta.Observacoes,
			&conta.NumeroDocumento,
			&conta.NumeroCompraNF,
			&conta.CreatedAt,
			&conta.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: erro ao escanear conta a pagar: %w", op, err)
		}
		contas = append(contas, conta)
	}

	return contas, nil
}