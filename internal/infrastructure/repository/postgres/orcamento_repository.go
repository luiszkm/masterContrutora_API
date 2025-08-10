// file: internal/repository/postgres/orcamento_repository.go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

type OrcamentoRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoOrcamentoRepository(db *pgxpool.Pool, logger *slog.Logger) suprimentos.OrcamentoRepository {
	return &OrcamentoRepositoryPostgres{
		db:     db,
		logger: logger,
	}
}
func (r *OrcamentoRepositoryPostgres) BuscarPorDetalhesID(ctx context.Context, id string) (*dto.OrcamentoDetalhadoDTO, error) {
	const op = "repository.postgres.orcamento.BuscarPorDetalhesID"

	query := `
		SELECT
			o.id, o.numero, o.valor_total, o.status, o.data_emissao, o.observacoes, o.condicoes_pagamento,
			json_build_object('id', ob.id, 'nome', ob.nome) as obra,
			json_build_object('id', e.id, 'nome', e.nome) as etapa,
			json_build_object('id', f.id, 'nome', f.nome) as fornecedor,
			COALESCE(json_agg(
				json_build_object(
					'id', oi.id,
					'produtoId', oi.produto_id,
					'produtoNome', p.nome,
					'unidadeDeMedida', p.unidade_de_medida,
					'categoria', p.categoria, -- CAMPO CRÍTICO ADICIONADO AQUI
					'quantidade', oi.quantidade,
					'valorUnitario', oi.valor_unitario
				)
			) FILTER (WHERE oi.id IS NOT NULL), '[]') as itens
		FROM orcamentos o
		JOIN etapas e ON o.etapa_id = e.id
		JOIN obras ob ON e.obra_id = ob.id
		JOIN fornecedores f ON o.fornecedor_id = f.id
		LEFT JOIN orcamento_itens oi ON o.id = oi.orcamento_id
		LEFT JOIN produtos p ON oi.produto_id = p.id
		WHERE o.id = $1
		GROUP BY o.id, ob.id, e.id, f.id
	`
	row := r.db.QueryRow(ctx, query, id)
	var d dto.OrcamentoDetalhadoDTO
	var obraJSON, etapaJSON, fornecedorJSON, itensJSON []byte

	err := row.Scan(
		&d.ID, &d.Numero, &d.ValorTotal, &d.Status, &d.DataEmissao, &d.Observacoes, &d.CondicoesPagamento,
		&obraJSON, &etapaJSON, &fornecedorJSON, &itensJSON,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: falha ao escanear orçamento: %w", op, err)
	}

	// A lógica de decodificação do JSON permanece a mesma.
	if err := json.Unmarshal(obraJSON, &d.Obra); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(etapaJSON, &d.Etapa); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fornecedorJSON, &d.Fornecedor); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(itensJSON, &d.Itens); err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *OrcamentoRepositoryPostgres) Atualizar(ctx context.Context, o *suprimentos.Orcamento) error {
	const op = "repository.postgres.orcamento.Atualizar"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// 1. Atualiza o registo principal na tabela 'orcamentos'.
	queryUpdateOrcamento := `
		UPDATE orcamentos SET
			etapa_id = $1, fornecedor_id = $2, valor_total = $3, observacoes = $4,
			condicoes_pagamento = $5, status = $6, updated_at = NOW(), data_aprovacao = $7
		WHERE id = $8
	`
	cmd, err := tx.Exec(ctx, queryUpdateOrcamento, o.EtapaID, o.FornecedorID, o.ValorTotal, o.Observacoes, o.CondicoesPagamento, o.Status, o.DataAprovacao, o.ID)
	if err != nil {
		return fmt.Errorf("%s: falha ao atualizar orçamento: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}

	// 2. Deleta todos os itens antigos.
	queryDeleteItens := `DELETE FROM orcamento_itens WHERE orcamento_id = $1`
	if _, err := tx.Exec(ctx, queryDeleteItens, o.ID); err != nil {
		return fmt.Errorf("%s: falha ao limpar itens antigos: %w", op, err)
	}

	// 3. Insere a nova lista de itens.
	if len(o.Itens) > 0 {
		queryInsertItem := `
			INSERT INTO orcamento_itens (id, orcamento_id, produto_id, quantidade, valor_unitario)
			VALUES ($1, $2, $3, $4, $5)
		`
		batch := &pgx.Batch{}
		for _, item := range o.Itens {
			batch.Queue(queryInsertItem, item.ID, item.OrcamentoID, item.ProdutoID, item.Quantidade, item.ValorUnitario)
		}
		br := tx.SendBatch(ctx, batch)
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("%s: falha ao executar lote de inserção de novos itens: %w", op, err)
		}
		if err := br.Close(); err != nil {
			return fmt.Errorf("%s: falha ao fechar lote de inserção de itens: %w", op, err)
		}
	}

	return tx.Commit(ctx)
}

// BuscarPorID implements suprimentos.OrcamentoRepository.
func (r *OrcamentoRepositoryPostgres) BuscarPorID(ctx context.Context, orcamentoID string) (*suprimentos.Orcamento, error) {
	const op = "repository.postgres.orcamento.BuscarPorID"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// 1. Busca o registo principal do orçamento (query já estava correta).
	queryOrcamento := `
		SELECT id, numero, etapa_id, fornecedor_id, valor_total, status, 
		       data_emissao, data_aprovacao, observacoes, condicoes_pagamento,
		       created_at, updated_at, deleted_at
		FROM orcamentos WHERE id = $1 AND deleted_at IS NULL`

	row := tx.QueryRow(ctx, queryOrcamento, orcamentoID)
	var o suprimentos.Orcamento

	if err := row.Scan(
		&o.ID, &o.Numero, &o.EtapaID, &o.FornecedorID, &o.ValorTotal, &o.Status,
		&o.DataEmissao, &o.DataAprovacao, &o.Observacoes, &o.CondicoesPagamento,
		&o.CreatedAt, &o.UpdatedAt, &o.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: falha ao buscar orçamento: %w", op, err)
	}

	// --- CORREÇÃO APLICADA AQUI ---
	// 2. Busca todos os itens, fazendo JOIN com a tabela produtos para obter os detalhes.
	queryItens := `
		SELECT
			oi.id, oi.orcamento_id, oi.produto_id, oi.quantidade, oi.valor_unitario,
			p.nome as produto_nome, p.unidade_de_medida, p.categoria
		FROM orcamento_itens oi
		JOIN produtos p ON oi.produto_id = p.id
		WHERE oi.orcamento_id = $1
	`
	rowsItens, err := tx.Query(ctx, queryItens, orcamentoID)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao buscar itens do orçamento: %w", op, err)
	}

	// Usamos RowToAddrOfStructByName para um mapeamento mais robusto.
	itens, err := pgx.CollectRows(rowsItens, pgx.RowToAddrOfStructByName[suprimentos.ItemOrcamento])
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao escanear itens do orçamento: %w", op, err)
	}
	// --- FIM DA CORREÇÃO ---

	// Convert []*suprimentos.ItemOrcamento to []suprimentos.ItemOrcamento
	o.Itens = make([]suprimentos.ItemOrcamento, len(itens))
	for i, itemPtr := range itens {
		if itemPtr != nil {
			o.Itens[i] = *itemPtr
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: falha ao finalizar transação de busca: %w", op, err)
	}

	return &o, nil
}

// ListarPorEtapaID implements suprimentos.OrcamentoRepository.
func (r *OrcamentoRepositoryPostgres) ListarPorEtapaID(ctx context.Context, etapaID string) ([]*suprimentos.Orcamento, error) {
	panic("unimplemented")
}

func (r *OrcamentoRepositoryPostgres) ListarOrcamentos(ctx context.Context, filtros common.ListarFiltros) ([]*dto.OrcamentoListItemDTO, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.orcamento.ListarTodos"

	args := pgx.NamedArgs{}
	whereClauses := []string{}

	if filtros.Status != "" {
		whereClauses = append(whereClauses, "o.status = @status")
		args["status"] = filtros.Status
	}
	if filtros.FornecedorID != "" {
		whereClauses = append(whereClauses, "o.fornecedor_id = @fornecedorID")
		args["fornecedorID"] = filtros.FornecedorID
	}
	if filtros.ObraID != "" {
		whereClauses = append(whereClauses, "e.obra_id = @obraID")
		args["obraID"] = filtros.ObraID
	}

	// A cláusula FROM agora inclui todos os JOINs necessários, incluindo produtos para categorias
	fromClause := `
		FROM orcamentos o
		JOIN etapas e ON o.etapa_id = e.id
		JOIN obras ob ON e.obra_id = ob.id
		JOIN fornecedores f ON o.fornecedor_id = f.id
		LEFT JOIN orcamento_itens oi ON o.id = oi.orcamento_id
		LEFT JOIN produtos p ON oi.produto_id = p.id
	`
	// Adiciona filtro para soft delete
	whereClauses = append(whereClauses, "o.deleted_at IS NULL")
	
	whereString := ""
	if len(whereClauses) > 0 {
		whereString = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Query de contagem
	countQuery := "SELECT COUNT(DISTINCT o.id)" + fromClause + whereString
	var totalItens int
	if err := r.db.QueryRow(ctx, countQuery, args).Scan(&totalItens); err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao contar orçamentos: %w", op, err)
	}

	paginacao := common.NewPaginacaoInfo(totalItens, filtros.Pagina, filtros.TamanhoPagina)
	if totalItens == 0 {
		return []*dto.OrcamentoListItemDTO{}, paginacao, nil
	}

	// --- QUERY PRINCIPAL CORRIGIDA E MELHORADA COM CATEGORIAS ---
	query := `
		SELECT
			o.id,
			o.numero,
			o.valor_total,
			o.status,
			o.data_emissao,
			o.fornecedor_id,
			f.nome as fornecedor_nome,
			ob.id as obra_id,
			ob.nome as obra_nome,
			COUNT(oi.id) as itens_count,
			COALESCE(array_agg(DISTINCT p.categoria) FILTER (WHERE p.categoria IS NOT NULL), ARRAY[]::text[]) as categorias
		` + fromClause + whereString + `
		GROUP BY o.id, f.nome, ob.id, ob.nome
		ORDER BY o.data_emissao DESC
		LIMIT @limit OFFSET @offset`

	args["limit"] = filtros.TamanhoPagina
	offset := (filtros.Pagina - 1) * filtros.TamanhoPagina
	args["offset"] = offset

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	orcamentos, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.OrcamentoListItemDTO])
	if err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao escanear orçamentos: %w", op, err)
	}

	return orcamentos, paginacao, nil
}

// Salvar usa uma transação para garantir atomicidade.
func (r *OrcamentoRepositoryPostgres) Salvar(ctx context.Context, o *suprimentos.Orcamento) error {
	const op = "repository.postgres.orcamento.Salvar"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// 1. Insere o registro principal na tabela 'orcamentos'
	queryOrcamento := `
		INSERT INTO orcamentos (id, numero, etapa_id, fornecedor_id, valor_total, status, data_emissao, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(ctx, queryOrcamento, o.ID, o.Numero, o.EtapaID, o.FornecedorID, o.ValorTotal, o.Status, o.DataEmissao, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		return fmt.Errorf("%s: falha ao inserir orçamento: %w", op, err)
	}

	// 2. Insere cada item na tabela 'orcamento_itens'
	if len(o.Itens) > 0 {
		queryItem := `
			INSERT INTO orcamento_itens (id, orcamento_id, produto_id, quantidade, valor_unitario)
			VALUES ($1, $2, $3, $4, $5)
		`
		batch := &pgx.Batch{}
		for _, item := range o.Itens {
			batch.Queue(queryItem, item.ID, item.OrcamentoID, item.ProdutoID, item.Quantidade, item.ValorUnitario)
		}

		batchResult := tx.SendBatch(ctx, batch)

		// --- CORREÇÃO CRÍTICA APLICADA AQUI ---
		// Executa o lote e verifica se houve erros na inserção dos itens.
		if _, err := batchResult.Exec(); err != nil {
			return fmt.Errorf("%s: falha ao executar lote de inserção de itens do orçamento: %w", op, err)
		}
		// --- FIM DA CORREÇÃO ---

		if err := batchResult.Close(); err != nil {
			return fmt.Errorf("%s: falha ao fechar lote de inserção de itens: %w", op, err)
		}
	}

	// Se tudo deu certo, confirma a transação (COMMIT).
	return tx.Commit(ctx)
}

func (r *OrcamentoRepositoryPostgres) ContarPorMesAno(ctx context.Context, ano int, mes time.Month) (int, error) {
	const op = "repository.postgres.orcamento.ContarPorMesAno"

	query := `
		SELECT COUNT(*)
		FROM orcamentos
		WHERE EXTRACT(YEAR FROM data_emissao) = $1
		  AND EXTRACT(MONTH FROM data_emissao) = $2
	`
	var count int
	err := r.db.QueryRow(ctx, query, ano, mes).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return count, nil
}

func (r *OrcamentoRepositoryPostgres) SoftDelete(ctx context.Context, id string) error {
	const op = "repository.postgres.orcamento.SoftDelete"
	query := `UPDATE orcamentos SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

// CompararPorCategoria retorna os orçamentos mais baratos para uma categoria específica
func (r *OrcamentoRepositoryPostgres) CompararPorCategoria(ctx context.Context, categoria string, limite int) ([]*dto.OrcamentoComparacao, error) {
	const op = "repository.postgres.orcamento.CompararPorCategoria"

	query := `
		SELECT DISTINCT 
			o.id, 
			o.numero, 
			f.nome as fornecedor_nome, 
			o.valor_total, 
			o.status, 
			o.data_emissao,
			COUNT(oi.id) as itens_categoria
		FROM orcamentos o
		JOIN fornecedores f ON o.fornecedor_id = f.id
		JOIN orcamento_itens oi ON o.id = oi.orcamento_id
		JOIN produtos p ON oi.produto_id = p.id
		WHERE p.categoria = $1
		  AND o.deleted_at IS NULL
		  AND o.status IN ('Aprovado', 'Em Aberto')
		GROUP BY o.id, f.nome, o.numero, o.valor_total, o.status, o.data_emissao
		ORDER BY o.valor_total ASC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, categoria, limite)
	if err != nil {
		return nil, fmt.Errorf("%s: falha na consulta: %w", op, err)
	}
	defer rows.Close()

	var orcamentos []*dto.OrcamentoComparacao
	for rows.Next() {
		var orc dto.OrcamentoComparacao
		err := rows.Scan(
			&orc.ID,
			&orc.Numero,
			&orc.FornecedorNome,
			&orc.ValorTotal,
			&orc.Status,
			&orc.DataEmissao,
			&orc.ItensCategoria,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: falha ao escanear resultado: %w", op, err)
		}
		orcamentos = append(orcamentos, &orc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro durante iteração: %w", op, err)
	}

	return orcamentos, nil
}
