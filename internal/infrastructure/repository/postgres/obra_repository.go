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
	"github.com/luiszkm/masterCostrutora/internal/platform/bus/db"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
)

var ErrNaoEncontrado = errors.New("recurso não encontrado")

// ObraRepositoryPostgres é a implementação do repositório de Obras para o PostgreSQL.
type ObraRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}
type ListarObrasFiltros struct {
	Status        string
	Pagina        int
	TamanhoPagina int
}

func NovaObraRepository(db *pgxpool.Pool, logger *slog.Logger) *ObraRepositoryPostgres {
	return &ObraRepositoryPostgres{db: db, logger: logger}
}

func (r *ObraRepositoryPostgres) ListarObras(ctx context.Context, filtros common.ListarFiltros) ([]*dto.ObraListItemDTO, *common.PaginacaoInfo, error) {
	const op = "repository.postgres.ListarObras"

	args := pgx.NamedArgs{}
	whereClauses := []string{"o.deleted_at IS NULL"}

	if filtros.Status != "" {
		whereClauses = append(whereClauses, "o.status = @status")
		args["status"] = filtros.Status
	}

	whereString := " WHERE " + strings.Join(whereClauses, " AND ")

	// Query para contar o total de itens
	countQuery := "SELECT COUNT(o.id) FROM obras o" + whereString
	var totalItens int
	err := r.db.QueryRow(ctx, countQuery, args).Scan(&totalItens)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: erro ao contar obras: %w", op, err)
	}

	paginacao := common.NewPaginacaoInfo(totalItens, filtros.Pagina, filtros.TamanhoPagina)
	if totalItens == 0 {
		return []*dto.ObraListItemDTO{}, paginacao, nil
	}

	// Query para buscar os dados da página, incluindo a etapa atual e a evolução.
	query := `
		SELECT
			o.id,
			o.nome,
			o.cliente,
			o.status,
			COALESCE(etapa_atual.nome, 'N/A') AS etapa,
			CONCAT(ROUND(COALESCE(etapa_stats.percentual_concluido, 0)), ' %') AS evolucao
		FROM
			obras o
		LEFT JOIN LATERAL (
			SELECT e.nome
			FROM etapas e
			WHERE e.obra_id = o.id AND e.status = 'Em Andamento'
			ORDER BY e.data_inicio_prevista DESC
			LIMIT 1
		) etapa_atual ON true
		LEFT JOIN LATERAL (
			SELECT
				(CAST(COUNT(CASE WHEN e.status = 'Concluída' THEN 1 END) AS FLOAT) / GREATEST(COUNT(e.id), 1)) * 100 AS percentual_concluido
			FROM etapas e
			WHERE e.obra_id = o.id
		) etapa_stats ON true
	` + whereString + ` ORDER BY o.nome ASC LIMIT @limit OFFSET @offset`

	args["limit"] = filtros.TamanhoPagina
	offset := (filtros.Pagina - 1) * filtros.TamanhoPagina
	args["offset"] = offset

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: erro ao listar obras: %w", op, err)
	}
	defer rows.Close()

	obrasList, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[dto.ObraListItemDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*dto.ObraListItemDTO{}, paginacao, nil
		}
		return nil, nil, fmt.Errorf("%s: falha ao escanear obras: %w", op, err)
	}

	return obrasList, paginacao, nil
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
        COALESCE(ps.custo_real, 0)
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
func (r *ObraRepositoryPostgres) Salvar(ctx context.Context, dbtx db.DBTX, obra *obras.Obra) error {
	const op = "repository.postgres.Salvar"
	query := `INSERT INTO obras (id, nome, cliente, endereco, data_inicio, data_fim, descricao, status,
	                           valor_contrato_total, valor_recebido, tipo_cobranca, data_assinatura_contrato)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := dbtx.Exec(ctx, query,
		obra.ID, obra.Nome, obra.Cliente, obra.Endereco,
		obra.DataInicio, obra.DataFim, obra.Descricao, obra.Status,
		obra.ValorContratoTotal, obra.ValorRecebido, obra.TipoCobranca, obra.DataAssinaturaContrato,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// BuscarPorID implementa a interface obras.Repository.
func (r *ObraRepositoryPostgres) BuscarPorID(ctx context.Context, id string) (*obras.Obra, error) {
	const op = "repository.postgres.BuscarPorID"

	query := `SELECT id, nome, cliente, endereco, data_inicio, data_fim, status, descricao,
	                 valor_contrato_total, valor_recebido, tipo_cobranca, data_assinatura_contrato 
	          FROM obras WHERE id = $1 AND deleted_at IS NULL`

	var obra obras.Obra
	err := r.db.QueryRow(ctx, query, id).Scan(
		&obra.ID, &obra.Nome, &obra.Cliente, &obra.Endereco, 
		&obra.DataInicio, &obra.DataFim, &obra.Status, &obra.Descricao,
		&obra.ValorContratoTotal, &obra.ValorRecebido, &obra.TipoCobranca, &obra.DataAssinaturaContrato,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNaoEncontrado
		}
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

func (r *ObraRepositoryPostgres) Atualizar(ctx context.Context, obra *obras.Obra) error {
	const op = "repository.postgres.obra.Atualizar"

	query := `
		UPDATE obras
		SET nome = $1, cliente = $2, endereco = $3, data_inicio = $4, data_fim = $5, status = $6, descricao = $7,
		    valor_contrato_total = $8, valor_recebido = $9, tipo_cobranca = $10, data_assinatura_contrato = $11
		WHERE id = $12 AND deleted_at IS NULL
	`

	cmd, err := r.db.Exec(ctx, query,
		obra.Nome,
		obra.Cliente,
		obra.Endereco,
		obra.DataInicio,
		obra.DataFim,
		obra.Status,
		obra.Descricao,
		obra.ValorContratoTotal,
		obra.ValorRecebido,
		obra.TipoCobranca,
		obra.DataAssinaturaContrato,
		obra.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}

	r.logger.InfoContext(ctx, "obra atualizada com sucesso", "obra_id", obra.ID)
	return nil
}
