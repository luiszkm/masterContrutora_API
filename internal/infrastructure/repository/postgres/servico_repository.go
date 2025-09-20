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
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
)

type ServicoRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovoServicoRepository(db *pgxpool.Pool, logger *slog.Logger) *ServicoRepositoryPostgres {
	return &ServicoRepositoryPostgres{db: db, logger: logger}
}

func (r *ServicoRepositoryPostgres) Salvar(ctx context.Context, s *suprimentos.Servico) error {
	const op = "repository.postgres.servico.Salvar"
	query := `INSERT INTO servicos (id, nome, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, s.ID, s.Nome, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *ServicoRepositoryPostgres) Atualizar(ctx context.Context, s *suprimentos.Servico) error {
	const op = "repository.postgres.servico.Atualizar"
	query := `UPDATE servicos SET nome = $1, updated_at = $2 WHERE id = $3`
	cmd, err := r.db.Exec(ctx, query, s.Nome, s.UpdatedAt, s.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

func (r *ServicoRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*suprimentos.Servico, error) {
	const op = "repository.postgres.servico.BuscarPorID"
	query := `SELECT id, nome, created_at, updated_at FROM servicos WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var s suprimentos.Servico
	err := row.Scan(&s.ID, &s.Nome, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &s, nil
}

func (r *ServicoRepositoryPostgres) ListarTodos(ctx context.Context) ([]*suprimentos.Servico, error) {
	const op = "repository.postgres.servico.ListarTodos"
	query := `SELECT id, nome, created_at, updated_at FROM servicos ORDER BY nome ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	servicos, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[suprimentos.Servico])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*suprimentos.Servico{}, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return servicos, nil
}

func (r *ServicoRepositoryPostgres) Listar(ctx context.Context, filtros common.ListarFiltros) ([]*suprimentos.Servico, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.servico.Listar"

	pagina := filtros.Pagina
	tamanhoPagina := filtros.TamanhoPagina
	if pagina <= 0 {
		pagina = 1
	}
	if tamanhoPagina <= 0 {
		tamanhoPagina = 10
	}

	whereConditions := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1

	whereClause := strings.Join(whereConditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM servicos
		WHERE %s`, whereClause)

	var totalItens int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItens)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: falha ao contar servicos: %w", op, err)
	}

	offset := (pagina - 1) * tamanhoPagina
	query := fmt.Sprintf(`
		SELECT id, nome, created_at, updated_at
		FROM servicos
		WHERE %s
		ORDER BY nome ASC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, tamanhoPagina, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	servicos, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[suprimentos.Servico])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			paginacaoInfo := common.NewPaginacaoInfo(0, pagina, tamanhoPagina)
			return []*suprimentos.Servico{}, paginacaoInfo, nil
		}
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	paginacaoInfo := common.NewPaginacaoInfo(totalItens, pagina, tamanhoPagina)

	return servicos, paginacaoInfo, nil
}

func (r *ServicoRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.servico.Deletar"
	query := `DELETE FROM servicos WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}