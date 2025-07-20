// file: internal/infrastructure/repository/postgres/querier_obra.go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
	"golang.org/x/sync/errgroup"
)

// BuscarDetalhesPorID orquestra a busca de todos os dados relacionados em paralelo.
func (q *ObraRepositoryPostgres) BuscarDetalhesPorID(ctx context.Context, obraID string) (*dto.ObraDetalhadaDTO, error) {
	const op = "querier.postgres.obra.BuscarDetalhesPorID"

	// Usamos um errgroup para executar as buscas em paralelo.
	// Se qualquer uma das goroutines retornar um erro, o contexto é cancelado
	// e o erro é propagado.
	g, gCtx := errgroup.WithContext(ctx)

	var obraBase *obras.Obra
	var etapas []dto.EtapaDTO
	var funcionarios []dto.FuncionarioAlocadoDTO
	var fornecedores []dto.FornecedorDTO
	var orcamentos []dto.OrcamentoDTO
	var produtos []dto.ProdutoDto
	// Goroutine para buscar a obra principal
	g.Go(func() error {
		var err error
		obraBase, err = q.fetchObraBase(gCtx, obraID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNaoEncontrado // Retorna um erro específico para 404
			}
			return fmt.Errorf("%s: falha ao buscar obra base: %w", op, err)
		}
		return nil
	})

	// Goroutine para buscar as etapas
	g.Go(func() error {
		var err error
		etapas, err = q.fetchEtapas(gCtx, obraID)
		if err != nil {
			return fmt.Errorf("%s: falha ao buscar etapas: %w", op, err)
		}
		return nil
	})

	// Goroutine para buscar os funcionários alocados
	g.Go(func() error {
		var err error
		funcionarios, err = q.fetchFuncionariosAlocados(gCtx, obraID)
		if err != nil {
			return fmt.Errorf("%s: falha ao buscar funcionários: %w", op, err)
		}
		return nil
	})
	g.Go(func() error {
		var err error
		fornecedores, err = q.fetchFornecedores(gCtx, obraID)
		if err != nil {
			q.logger.ErrorContext(gCtx, "falha na goroutine fetchFornecedores", "erro", err)
			return fmt.Errorf("%s: falha ao buscar fornecedores: %w", op, err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		orcamentos, err = q.fetchOrcamentos(gCtx, obraID)
		if err != nil {
			q.logger.ErrorContext(gCtx, "falha na goroutine fetchOrcamentos", "erro", err)
			return fmt.Errorf("%s: falha ao buscar orçamentos: %w", op, err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		produtos, err = q.fetchProdutos(gCtx, obraID)
		if err != nil {
			q.logger.ErrorContext(gCtx, "falha na goroutine fetchProdutos", "erro", err)
			return fmt.Errorf("%s: falha ao buscar produtos: %w", op, err)
		}
		return nil
	})
	// Aguarda todas as goroutines terminarem.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Monta o DTO final com os resultados coletados.
	return &dto.ObraDetalhadaDTO{
		ID:           obraBase.ID,
		Nome:         obraBase.Nome,
		Cliente:      obraBase.Cliente,
		Endereco:     obraBase.Endereco,
		DataInicio:   obraBase.DataInicio,
		DataFim:      obraBase.DataFim,
		Status:       string(obraBase.Status),
		Descricao:    obraBase.Descricao, // Adiciona descrição opcional
		Etapas:       etapas,
		Funcionarios: funcionarios,
		Fornecedores: fornecedores,
		Orcamentos:   orcamentos,
		Produtos:     produtos,
	}, nil
}

func (q *ObraRepositoryPostgres) fetchObraBase(ctx context.Context, obraID string) (*obras.Obra, error) {
	query := `SELECT id, nome, cliente, endereco,descricao, data_inicio, data_fim, status, deleted_at FROM obras WHERE id = $1`
	row, err := q.db.Query(ctx, query, obraID)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	return pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[obras.Obra])
}

func (q *ObraRepositoryPostgres) fetchEtapas(ctx context.Context, obraID string) ([]dto.EtapaDTO, error) {
	query := `SELECT id, nome, data_inicio_prevista, data_fim_prevista, status FROM etapas WHERE obra_id = $1 ORDER BY data_inicio_prevista` // Corrected line
	rows, err := q.db.Query(ctx, query, obraID)                                                                                              // Corrected line
	if err != nil {                                                                                                                          // Corrected line
		return nil, err // Corrected line
	} // Corrected line
	defer rows.Close()                                                // Corrected line
	return pgx.CollectRows(rows, pgx.RowToStructByName[dto.EtapaDTO]) // Corrected line
}

func (q *ObraRepositoryPostgres) fetchFuncionariosAlocados(ctx context.Context, obraID string) ([]dto.FuncionarioAlocadoDTO, error) {
	query := `
		SELECT
			a.funcionario_id,
			f.nome as nome_funcionario,
			a.data_inicio_alocacao
		FROM alocacoes a
		JOIN funcionarios f ON a.funcionario_id = f.id
		WHERE a.obra_id = $1
		ORDER BY f.nome
	`
	rows, err := q.db.Query(ctx, query, obraID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[dto.FuncionarioAlocadoDTO])
}
func (q *ObraRepositoryPostgres) fetchFornecedores(ctx context.Context, obraID string) ([]dto.FornecedorDTO, error) {
	// A lógica é: encontre todas as etapas da obra, depois todos os orçamentos dessas etapas,
	// e finalmente os fornecedores distintos desses orçamentos.
	query := `
		SELECT DISTINCT f.id, f.nome
		FROM fornecedores f
		JOIN orcamentos o ON f.id = o.fornecedor_id
		JOIN etapas e ON o.etapa_id = e.id
		WHERE e.obra_id = $1
		ORDER BY f.nome` // Corrected line
	rows, err := q.db.Query(ctx, query, obraID) // Corrected line
	if err != nil {                             // Corrected line
		return nil, err // Corrected line
	} // Corrected line
	defer rows.Close()                                                     // Corrected line
	return pgx.CollectRows(rows, pgx.RowToStructByName[dto.FornecedorDTO]) // Corrected line
}

func (q *ObraRepositoryPostgres) fetchOrcamentos(ctx context.Context, obraID string) ([]dto.OrcamentoDTO, error) {
	query := `
		SELECT o.id, o.numero, o.valor_total, o.status
		FROM orcamentos o
		JOIN etapas e ON o.etapa_id = e.id
		WHERE e.obra_id = $1
		ORDER BY o.data_emissao DESC`
	rows, err := q.db.Query(ctx, query, obraID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[dto.OrcamentoDTO])
}

func (q *ObraRepositoryPostgres) fetchProdutos(ctx context.Context, obraID string) ([]dto.ProdutoDto, error) {
	// A lógica é: encontre as etapas da obra, os orçamentos, os itens de orçamento,
	// e finalmente os produtos distintos desses itens.
	query := `
		SELECT DISTINCT m.id, m.nome
		FROM produtos m
		JOIN orcamento_itens oi ON m.id = oi.produto_id
		JOIN orcamentos o ON oi.orcamento_id = o.id
		JOIN etapas e ON o.etapa_id = e.id
		WHERE e.obra_id = $1
		ORDER BY m.nome` // Corrected line
	rows, err := q.db.Query(ctx, query, obraID) // Corrected line
	if err != nil {                             // Corrected line
		return nil, err // Corrected line
	} // Corrected line
	defer rows.Close()                                                  // Corrected line
	return pgx.CollectRows(rows, pgx.RowToStructByName[dto.ProdutoDto]) // Corrected line
}
