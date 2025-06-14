// file: internal/repository/postgres/orcamento_repository.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
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

func (r *OrcamentoRepositoryPostgres) Atualizar(ctx context.Context, o *suprimentos.Orcamento) error {
	const op = "repository.postgres.orcamento.Atualizar"

	// Por enquanto, nosso caso de uso só atualiza o status, mas a query está pronta para mais.
	query := `UPDATE orcamentos SET status = $1 WHERE id = $2`

	cmd, err := r.db.Exec(ctx, query, o.Status, o.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
func (r *OrcamentoRepositoryPostgres) BuscarPorID(ctx context.Context, orcamentoID string) (*suprimentos.Orcamento, error) {
	const op = "repository.postgres.orcamento.BuscarPorID"

	// Usaremos uma transação para garantir a consistência da leitura
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// 1. Busca o registro principal do orçamento
	queryOrcamento := `SELECT id, numero, etapa_id, fornecedor_id, valor_total, status, data_emissao FROM orcamentos WHERE id = $1`
	row := tx.QueryRow(ctx, queryOrcamento, orcamentoID)

	var o suprimentos.Orcamento
	if err := row.Scan(&o.ID, &o.Numero, &o.EtapaID, &o.FornecedorID, &o.ValorTotal, &o.Status, &o.DataEmissao); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: falha ao buscar orçamento: %w", op, err)
	}

	// 2. Busca todos os itens associados a este orçamento
	queryItens := `SELECT id, orcamento_id, material_id, quantidade, valor_unitario FROM orcamento_itens WHERE orcamento_id = $1`
	rowsItens, err := tx.Query(ctx, queryItens, orcamentoID)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao buscar itens do orçamento: %w", op, err)
	}

	itens, err := pgx.CollectRows(rowsItens, func(row pgx.CollectableRow) (suprimentos.ItemOrcamento, error) {
		var item suprimentos.ItemOrcamento
		err := row.Scan(&item.ID, &item.OrcamentoID, &item.MaterialID, &item.Quantidade, &item.ValorUnitario)
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao escanear itens do orçamento: %w", op, err)
	}

	o.Itens = itens

	// Finaliza a transação (apenas leitura, então commit ou rollback não alteram dados)
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: falha ao finalizar transação de busca: %w", op, err)
	}

	return &o, nil
}

// ListarPorEtapaID implements suprimentos.OrcamentoRepository.
func (r *OrcamentoRepositoryPostgres) ListarPorEtapaID(ctx context.Context, etapaID string) ([]*suprimentos.Orcamento, error) {
	panic("unimplemented")
}

// ListarTodos implements suprimentos.OrcamentoRepository.
func (r *OrcamentoRepositoryPostgres) ListarTodos(ctx context.Context) ([]*suprimentos.Orcamento, error) {
	panic("unimplemented")
}

// Salvar usa uma transação para garantir atomicidade.
func (r *OrcamentoRepositoryPostgres) Salvar(ctx context.Context, o *suprimentos.Orcamento) error {
	const op = "repository.postgres.orcamento.Salvar"

	// Inicia a transação
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	// Garante que a transação seja desfeita (ROLLBACK) em caso de erro em qualquer ponto.
	defer tx.Rollback(ctx)

	// 1. Insere o registro principal na tabela 'orcamentos'
	queryOrcamento := `
		INSERT INTO orcamentos (id, numero, etapa_id, fornecedor_id, valor_total, status, data_emissao)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, queryOrcamento, o.ID, o.Numero, o.EtapaID, o.FornecedorID, o.ValorTotal, o.Status, o.DataEmissao)
	if err != nil {
		return fmt.Errorf("%s: falha ao inserir orçamento: %w", op, err)
	}

	// 2. Insere cada item na tabela 'orcamento_itens'
	queryItem := `
		INSERT INTO orcamento_itens (id, orcamento_id, material_id, quantidade, valor_unitario)
		VALUES ($1, $2, $3, $4, $5)
	`
	// Usamos um "batch" para inserir múltiplos itens de forma eficiente dentro da mesma transação.
	batch := &pgx.Batch{}
	for _, item := range o.Itens {
		batch.Queue(queryItem, item.ID, item.OrcamentoID, item.MaterialID, item.Quantidade, item.ValorUnitario)
	}

	batchResult := tx.SendBatch(ctx, batch)
	if err := batchResult.Close(); err != nil {
		return fmt.Errorf("%s: falha ao inserir itens do orçamento: %w", op, err)
	}

	// Se tudo deu certo até aqui, confirma a transação (COMMIT).
	return tx.Commit(ctx)
}
