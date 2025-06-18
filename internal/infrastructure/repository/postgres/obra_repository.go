// file: internal/repository/postgres/obra_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
)

var ErrNaoEncontrado = errors.New("recurso não encontrado")

// ObraRepositoryPostgres é a implementação do repositório de Obras para o PostgreSQL.
type ObraRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NovaObraRepository(db *pgxpool.Pool, logger *slog.Logger) *ObraRepositoryPostgres {
	return &ObraRepositoryPostgres{db: db, logger: logger}
}

func (r *ObraRepositoryPostgres) ListarObras(ctx context.Context, filtros common.ListarFiltros) ([]*dto.ObraListItemDTO, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.ListarObras"

	var args []interface{}
	var countArgs []interface{} // Slice separado para os argumentos da contagem
	paramCount := 1

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("SELECT id, nome, cliente, status FROM obras WHERE deleted_at IS NULL")

	countQueryBuilder := strings.Builder{}
	countQueryBuilder.WriteString("SELECT COUNT(*) FROM obras WHERE deleted_at IS NULL")

	if filtros.Status != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", paramCount))
		countQueryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", paramCount))
		// Adiciona o argumento a ambos os slices
		args = append(args, filtros.Status)
		countArgs = append(countArgs, filtros.Status)
		paramCount++
	}

	var totalItens int
	err := r.db.QueryRow(ctx, countQueryBuilder.String(), countArgs...).Scan(&totalItens)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: erro ao contar obras: %w", op, err)
	}

	paginacao := common.NewPaginacaoInfo(totalItens, filtros.Pagina, filtros.TamanhoPagina)
	// Se não houver itens, retornamos uma lista vazia e a paginação correta.
	if totalItens == 0 {
		return []*dto.ObraListItemDTO{}, paginacao, nil
	}

	offset := (filtros.Pagina - 1) * filtros.TamanhoPagina
	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY nome ASC LIMIT $%d OFFSET $%d", paramCount, paramCount+1))
	args = append(args, filtros.TamanhoPagina, offset)

	rows, err := r.db.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: erro ao listar obras: %w", op, err)
	}
	defer rows.Close()

	// SIMPLIFICAÇÃO: Usamos pgx.CollectRows para escanear diretamente para a struct do DTO.
	obrasList, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[dto.ObraListItemDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*dto.ObraListItemDTO{}, paginacao, nil
		}
		return nil, nil, fmt.Errorf("%s: falha ao escanear obras: %w", op, err)
	}

	// Não precisamos mais do loop de conversão no final.
	return obrasList, paginacao, nil
}

type ListarObrasFiltros struct {
	Status        string
	Pagina        int
	TamanhoPagina int
}

// BuscarDashboardPorID implements obras.Querier.
func (r *ObraRepositoryPostgres) BuscarDashboardPorID(ctx context.Context, id string) (*dto.ObraDashboard, error) {
	const op = "repository.postgres.BuscarDashboardPorID"

	// Usamos CTEs (WITH clauses) para tornar a query mais legível e modular.
	// Cada CTE calcula um pedaço da informação que precisamos.
	query := `
    WITH etapa_stats AS (
        -- CTE para calcular o percentual de conclusão
        SELECT
            obra_id,
            CAST(COUNT(CASE WHEN status = 'Concluída' THEN 1 END) AS FLOAT) / GREATEST(COUNT(*), 1) * 100 AS percentual_concluido
        FROM etapas
        GROUP BY obra_id
    ),
    alocacao_stats AS (
        -- CTE para contar os funcionários atualmente alocados
        SELECT
            obra_id,
            COUNT(*) AS funcionarios_alocados
        FROM alocacoes
        WHERE data_fim_alocacao IS NULL OR data_fim_alocacao >= CURRENT_DATE
        GROUP BY obra_id
    ),
    etapa_atual AS (
        -- CTE para encontrar a etapa que está "Em Andamento"
        SELECT
            obra_id,
            nome,
            data_fim_prevista
        FROM etapas
        WHERE status = 'Em Andamento'
        ORDER BY data_inicio_prevista DESC
        LIMIT 1
    ),
    orcamento_stats AS (
        -- CTE para calcular os dados financeiros a partir dos orçamentos.
   		SELECT
				e.obra_id,
				COALESCE(SUM(o.valor_total) FILTER (WHERE o.status IN ('Aprovado', 'Pago')), 0) AS orcamento_total_aprovado
			FROM orcamentos o
			JOIN etapas e ON o.etapa_id = e.id
			GROUP BY e.obra_id
    ),
	pagamento_stats AS (
			-- NOVA CTE: Calcula o custo real somando os pagamentos efetivados.
			SELECT
				obra_id,
				COALESCE(SUM(valor_calculado), 0) AS custo_real
			FROM registros_pagamento
			GROUP BY obra_id
		)
    -- Query Principal que junta tudo. Começa aqui, FORA da CTE anterior.
    SELECT
        o.id,
        o.nome,
        o.status,
        ea.nome,
        ea.data_fim_prevista,
        COALESCE(es.percentual_concluido, 0),
        COALESCE(als.funcionarios_alocados, 0),
        COALESCE(os.orcamento_total_aprovado, 0),
        COALESCE(os.custo_total_realizado, 0)
    FROM
        obras o
    LEFT JOIN etapa_stats es ON o.id = es.obra_id
    LEFT JOIN alocacao_stats als ON o.id = als.obra_id
    LEFT JOIN etapa_atual ea ON o.id = ea.obra_id
    LEFT JOIN orcamento_stats os ON o.id = os.obra_id
	LEFT JOIN pagamento_stats ps ON o.id = ps.obra_id
    WHERE
        o.id = $1
	AND o.deleted_at IS NULL
`
	row := r.db.QueryRow(ctx, query, id)

	var dashboard dto.ObraDashboard
	var etapaAtualNome sql.NullString
	var dataFimPrevistaEtapa sql.NullTime

	err := row.Scan(
		&dashboard.ObraID,
		&dashboard.NomeObra,
		&dashboard.StatusObra,
		&etapaAtualNome,
		&dataFimPrevistaEtapa,
		&dashboard.PercentualConcluido,
		&dashboard.FuncionariosAlocados,
		&dashboard.OrcamentoTotalAprovado,
		&dashboard.CustoTotalRealizado,
	)
	if err != nil {
		// Tratamento específico para o erro "não encontrado".
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNaoEncontrado
		}
		// Para todos os outros erros, adicionamos nosso contexto.
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Após o Scan, convertemos os tipos Null* para os ponteiros do nosso DTO.
	if etapaAtualNome.Valid {
		dashboard.EtapaAtualNome = &etapaAtualNome.String
	}
	if dataFimPrevistaEtapa.Valid {
		dashboard.DataFimPrevistaEtapa = &dataFimPrevistaEtapa.Time
		// Calculamos os dias restantes apenas se a data for válida.
		dias := int(time.Until(*dashboard.DataFimPrevistaEtapa).Hours() / 24)
		dashboard.DiasParaPrazoEtapa = &dias
	}

	dashboard.BalancoFinanceiro = dashboard.OrcamentoTotalAprovado - dashboard.CustoTotalRealizado
	dashboard.UltimaAtualizacao = time.Now()
	dashboard.BalancoFinanceiro = dashboard.OrcamentoTotalAprovado - dashboard.CustoTotalRealizado
	dashboard.UltimaAtualizacao = time.Now()

	return &dashboard, nil
}

// Salvar agora insere a obra no banco de dados.
func (r *ObraRepositoryPostgres) Salvar(ctx context.Context, obra *obras.Obra) error {
	const op = "repository.postgres.Salvar"

	// A query SQL para inserir uma nova obra.
	// Usamos $1, $2, etc., como placeholders para prevenir SQL Injection.
	query := `INSERT INTO obras (id, nome, cliente, endereco, data_inicio, data_fim, status)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	// Passamos o contexto para a chamada do banco, permitindo o cancelamento.
	_, err := r.db.Exec(ctx, query,
		obra.ID,
		obra.Nome,
		obra.Cliente,
		obra.Endereco,
		obra.DataInicio,
		obra.DataFim, // Será NULL se o time.Time estiver zerado
		obra.Status,
	)

	if err != nil {
		// Adicionamos contexto ao erro para facilitar a depuração.
		return fmt.Errorf("%s: %w", op, err)
	}

	r.logger.InfoContext(ctx, "obra salva no banco de dados com sucesso", "obra_id", obra.ID)
	return nil
}

// BuscarPorID implementa a interface obras.Repository.
func (r *ObraRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*obras.Obra, error) {
	const op = "repository.postgres.BuscarPorID"

	query := `SELECT id, nome, cliente, endereco, data_inicio, data_fim, status FROM obras WHERE id = $1`

	var obra obras.Obra
	err := r.db.QueryRow(ctx, query, id).Scan(&obra.ID, &obra.Nome, &obra.Cliente, &obra.Endereco, &obra.DataInicio, &obra.DataFim, &obra.Status)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &obra, nil
}

func (r *ObraRepositoryPostgres) Deletar(ctx context.Context, id string) error {
	const op = "repository.postgres.obra.Deletar"
	query := `UPDATE obras SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado // Usa nosso erro padrão se a obra não existir ou já estiver deletada
	}
	return nil
}
