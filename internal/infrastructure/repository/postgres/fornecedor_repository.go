// file: internal/repository/postgres/fornecedor_repository.go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
)

// FornecedorRepositoryPostgres implementa a interface suprimentos.Repository.
type FornecedorRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoFornecedorRepository(db *pgxpool.Pool, logger *slog.Logger) suprimentos.FornecedorRepository {
	return &FornecedorRepositoryPostgres{db: db, logger: logger}
}

// Atualizar implements suprimentos.FornecedorRepository.
func (r *FornecedorRepositoryPostgres) Atualizar(ctx context.Context, f *suprimentos.Fornecedor, categoriaIDs *[]string) error {
	const op = "repository.postgres.fornecedor.Atualizar"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx)

	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	if f.Nome != "" {
		setClauses = append(setClauses, fmt.Sprintf("nome = $%d", argIndex))
		args = append(args, f.Nome)
		argIndex++
	}
	if f.CNPJ != "" {
		setClauses = append(setClauses, fmt.Sprintf("cnpj = $%d", argIndex))
		args = append(args, f.CNPJ)
		argIndex++
	}
	if f.Contato != "" {
		setClauses = append(setClauses, fmt.Sprintf("contato = $%d", argIndex))
		args = append(args, f.Contato)
		argIndex++
	}
	if f.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, f.Email)
		argIndex++
	}
	if f.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, f.Status)
		argIndex++
	}
	if f.Endereco != nil {
		setClauses = append(setClauses, fmt.Sprintf("endereco = $%d", argIndex))
		args = append(args, f.Endereco)
		argIndex++
	}
	if f.Avaliacao != nil {
		setClauses = append(setClauses, fmt.Sprintf("avaliacao = $%d", argIndex))
		args = append(args, f.Avaliacao)
		argIndex++
	}
	if f.Observacoes != nil {
		setClauses = append(setClauses, fmt.Sprintf("observacoes = $%d", argIndex))
		args = append(args, f.Observacoes)
		argIndex++
	}

	if len(setClauses) > 0 {
		query := fmt.Sprintf(`UPDATE fornecedores SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), argIndex)
		args = append(args, f.ID)

		if _, err := tx.Exec(ctx, query, args...); err != nil {
			return fmt.Errorf("%s: falha ao atualizar fornecedor: %w", op, err)
		}
	}

	if categoriaIDs != nil {
		queryDelete := `DELETE FROM fornecedor_categorias WHERE fornecedor_id = $1`
		if _, err := tx.Exec(ctx, queryDelete, f.ID); err != nil {
			return fmt.Errorf("%s: falha ao limpar categorias: %w", op, err)
		}

		if len(*categoriaIDs) > 0 {
			batch := &pgx.Batch{}
			queryInsert := `INSERT INTO fornecedor_categorias (fornecedor_id, categoria_id) VALUES ($1, $2)`
			for _, catID := range *categoriaIDs {
				batch.Queue(queryInsert, f.ID, catID)
			}
			if err := tx.SendBatch(ctx, batch).Close(); err != nil {
				return fmt.Errorf("%s: falha ao inserir categorias: %w", op, err)
			}
		}
	}

	return tx.Commit(ctx)
}

// Deletar implements suprimentos.FornecedorRepository.
func (r *FornecedorRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.fornecedor.Deletar"
	query := `
		UPDATE fornecedores
		SET deleted_at = NOW(), 
			status = 'Inativo'
		 WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: falha ao atualizar fornecedor: %w", op, err)
	}
	return nil

}

// BuscarPorID implements suprimentos.FornecedorRepository.
func (r *FornecedorRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*suprimentos.Fornecedor, error) {
	const op = "repository.postgres.fornecedor.BuscarPorID"

	// A query agora usa LEFT JOIN e json_agg para buscar as categorias em uma única chamada.
	query := `
		SELECT
			f.id, f.nome, f.cnpj, f.contato, f.email, f.status,
			COALESCE(
				json_agg(json_build_object('ID', c.id, 'Nome', c.nome)) FILTER (WHERE c.id IS NOT NULL),
				'[]'
			) as categorias
		FROM fornecedores f
		LEFT JOIN fornecedor_categorias fc ON f.id = fc.fornecedor_id
		LEFT JOIN categorias c ON fc.categoria_id = c.id
		WHERE f.id = $1 AND f.deleted_at IS NULL
		GROUP BY f.id`

	var f suprimentos.Fornecedor
	var categoriasJSON []byte // Recebe o resultado do json_agg como bytes

	err := r.db.QueryRow(ctx, query, id).Scan(&f.ID, &f.Nome, &f.CNPJ, &f.Contato, &f.Email, &f.Status, &categoriasJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: falha ao escanear fornecedor: %w", op, err)
	}

	// Decodifica o JSON das categorias para a struct Go.
	if err := json.Unmarshal(categoriasJSON, &f.Categorias); err != nil {
		return nil, fmt.Errorf("%s: falha ao decodificar JSON das categorias: %w", op, err)
	}

	return &f, nil
}

func (r *FornecedorRepositoryPostgres) Salvar(ctx context.Context, f *suprimentos.Fornecedor, categoriaIDs []string) error {
	const op = "repository.postgres.fornecedor.Salvar"

	// Inicia a transação
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: falha ao iniciar transação: %w", op, err)
	}
	defer tx.Rollback(ctx) // Garante o rollback em caso de erro

	// 1. Insere o registro principal na tabela 'fornecedores'
	queryFornecedor := `
		INSERT INTO fornecedores (id, nome, cnpj, contato, email, status, endereco, avaliacao, observacoes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(ctx, queryFornecedor, f.ID, f.Nome, f.CNPJ, f.Contato, f.Email, f.Status, f.Endereco, f.Avaliacao, f.Observacoes)
	if err != nil {
		return fmt.Errorf("%s: falha ao inserir fornecedor: %w", op, err)
	}
	// 2. Insere as associações na tabela de junção 'fornecedor_categorias'
	if len(categoriaIDs) > 0 {
		queryCategorias := `INSERT INTO fornecedor_categorias (fornecedor_id, categoria_id) VALUES ($1, $2)`
		batch := &pgx.Batch{}
		for _, catID := range categoriaIDs {
			batch.Queue(queryCategorias, f.ID, catID)
		}

		br := tx.SendBatch(ctx, batch)
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("%s: falha ao executar lote de associações de categoria: %w", op, err)
		}

		if err := br.Close(); err != nil {
			return fmt.Errorf("%s: falha ao inserir associações de categoria: %w", op, err)
		}
	}

	// Se tudo deu certo, confirma a transação (COMMIT).
	return tx.Commit(ctx)
}

func (r *FornecedorRepositoryPostgres) ListarTodos(ctx context.Context) ([]*suprimentos.Fornecedor, error) {
	const op = "repository.postgres.fornecedor.ListarTodos"
	query := `
		SELECT
			f.id, f.nome, f.cnpj, f.contato, f.email, f.status, f.endereco, f.avaliacao, f.observacoes,
			COUNT(DISTINCT o.id) as orcamentos_count, -- Contagem de orçamentos distintos
			COALESCE(
				json_agg(json_build_object('ID', c.id, 'Nome', c.nome)) FILTER (WHERE c.id IS NOT NULL),
				'[]'
			) as categorias
		FROM fornecedores f
		LEFT JOIN fornecedor_categorias fc ON f.id = fc.fornecedor_id
		LEFT JOIN categorias c ON fc.categoria_id = c.id
		LEFT JOIN orcamentos o ON f.id = o.fornecedor_id 
		WHERE f.deleted_at IS NULL
		GROUP BY f.id 
		ORDER BY f.nome ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var fornecedores []*suprimentos.Fornecedor
	for rows.Next() {
		var f suprimentos.Fornecedor
		var categoriasJSON []byte

		if err := rows.Scan(
			&f.ID, &f.Nome, &f.CNPJ, &f.Contato, &f.Email, &f.Status, &f.Endereco, &f.Avaliacao, &f.Observacoes,
			&f.OrcamentosCount, // NOVO CAMPO NO SCAN
			&categoriasJSON,
		); err != nil {
			return nil, fmt.Errorf("%s: falha ao escanear linha de fornecedor: %w", op, err)
		}

		if err := json.Unmarshal(categoriasJSON, &f.Categorias); err != nil {
			return nil, fmt.Errorf("%s: falha ao decodificar JSON das categorias para %s: %w", op, f.Nome, err)
		}
		fornecedores = append(fornecedores, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: erro ao iterar sobre as linhas: %w", op, err)
	}

	return fornecedores, nil
}
