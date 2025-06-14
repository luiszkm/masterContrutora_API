// file: internal/repository/postgres/obra_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
)

// ObraRepositoryPostgres é a implementação do repositório de Obras para o PostgreSQL.
type ObraRepositoryPostgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

var ErrNaoEncontrado = errors.New("recurso não encontrado")

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
		)
		-- Query Principal que junta tudo
		SELECT
			o.id,
			o.nome,
			o.status,
			ea.nome,
			ea.data_fim_prevista,
			COALESCE(es.percentual_concluido, 0),
			COALESCE(als.funcionarios_alocados, 0)
		FROM
			obras o
		LEFT JOIN etapa_stats es ON o.id = es.obra_id
		LEFT JOIN alocacao_stats als ON o.id = als.obra_id
		LEFT JOIN etapa_atual ea ON o.id = ea.obra_id
		WHERE
			o.id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var dashboard dto.ObraDashboard
	// Usamos tipos sql.Null* para escanear colunas que podem ser nulas (vindas dos LEFT JOINs).
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

	// TODO: Substituir por cálculos reais quando os agregados Financeiro e Suprimentos existirem.
	dashboard.CustoTotalRealizado = 150000.75
	dashboard.OrcamentoTotalAprovado = 200000.00
	dashboard.BalancoFinanceiro = dashboard.OrcamentoTotalAprovado - dashboard.CustoTotalRealizado
	dashboard.UltimaAtualizacao = time.Now()

	return &dashboard, nil
}

// NovaObraRepository cria uma nova instância do repositório de obras.
func NovaObraRepository(db *pgxpool.Pool, logger *slog.Logger) *ObraRepositoryPostgres {
	return &ObraRepositoryPostgres{
		db:     db,
		logger: logger,
	}
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
