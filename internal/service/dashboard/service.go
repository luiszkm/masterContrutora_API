package dashboard

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/luiszkm/masterCostrutora/internal/domain/dashboard"
	"github.com/luiszkm/masterCostrutora/internal/service/dashboard/dto"
	"github.com/luiszkm/masterCostrutora/pkg/logging"
)

// Service representa o serviço de dashboard
type Service struct {
	querier       dashboard.Querier
	logger        *slog.Logger
	dashLogger    *logging.DashboardLogger
}

// NovoServicoDashboard cria uma nova instância do serviço de dashboard
func NovoServicoDashboard(querier dashboard.Querier, logger *slog.Logger, dashLogger *logging.DashboardLogger) *Service {
	return &Service{
		querier:    querier,
		logger:     logger,
		dashLogger: dashLogger,
	}
}

// ObterDashboardCompleto retorna o dashboard completo com todas as seções
func (s *Service) ObterDashboardCompleto(ctx context.Context, parametros dto.ParametrosDashboardDTO) (*dto.DashboardGeralDTO, error) {
	const op = "service.dashboard.ObterDashboardCompleto"
	startTime := time.Now()

	// Log início da operação
	s.dashLogger.LogDashboardServiceCall(ctx, "geral", "ObterDashboardCompleto", startTime, nil, map[string]interface{}{
		"parametros": parametros,
	})

	// Define período padrão se não fornecido
	dataInicio := parametros.DataInicio
	dataFim := parametros.DataFim
	if dataInicio == nil {
		inicio := time.Now().AddDate(0, -6, 0) // 6 meses atrás
		dataInicio = &inicio
	}
	if dataFim == nil {
		fim := time.Now()
		dataFim = &fim
	}

	// Validar parâmetros
	if dataFim.Before(*dataInicio) {
		err := fmt.Errorf("data de fim não pode ser anterior à data de início")
		s.dashLogger.LogDashboardValidation(ctx, "geral", "periodo", map[string]interface{}{
			"dataInicio": dataInicio,
			"dataFim": dataFim,
		}, err.Error(), nil)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	dashboard := &dto.DashboardGeralDTO{
		UltimaAtualizacao: time.Now(),
	}

	var err error
	var queryCount int

	// Obter resumo geral sempre
	queryStart := time.Now()
	dashboard.ResumoGeral, err = s.querier.ObterResumoGeral(ctx)
	if err != nil {
		s.dashLogger.LogDashboardError(ctx, "geral", "ObterResumoGeral", err, map[string]interface{}{
			"operation": op,
		})
		return nil, fmt.Errorf("%s: falha ao obter resumo geral: %w", op, err)
	}
	s.dashLogger.LogDashboardQuery(ctx, "geral", "ObterResumoGeral", time.Since(queryStart), 1, nil)
	queryCount++

	// Obter alertas sempre
	queryStart = time.Now()
	dashboard.Alertas, err = s.querier.ObterAlertas(ctx)
	if err != nil {
		s.dashLogger.LogDashboardError(ctx, "geral", "ObterAlertas", err, map[string]interface{}{
			"operation": op,
		})
		return nil, fmt.Errorf("%s: falha ao obter alertas: %w", op, err)
	}
	s.dashLogger.LogDashboardQuery(ctx, "geral", "ObterAlertas", time.Since(queryStart), len(dashboard.Alertas.ObrasComAtraso), nil)
	queryCount++

	// Obter seções específicas ou todas se não especificado
	secoes := parametros.Secoes
	if len(secoes) == 0 {
		secoes = []string{"financeiro", "obras", "funcionarios", "fornecedores"}
	}

	for _, secao := range secoes {
		switch secao {
		case "financeiro":
			dashboard.Financeiro, err = s.ObterDashboardFinanceiro(ctx, *dataInicio, *dataFim)
			if err != nil {
				s.logger.ErrorContext(ctx, "falha ao obter dashboard financeiro", "erro", err)
				return nil, fmt.Errorf("%s: falha ao obter dashboard financeiro: %w", op, err)
			}
		case "obras":
			dashboard.Obras, err = s.ObterDashboardObras(ctx)
			if err != nil {
				s.logger.ErrorContext(ctx, "falha ao obter dashboard obras", "erro", err)
				return nil, fmt.Errorf("%s: falha ao obter dashboard obras: %w", op, err)
			}
		case "funcionarios":
			dashboard.Funcionarios, err = s.ObterDashboardFuncionarios(ctx, *dataInicio, *dataFim)
			if err != nil {
				s.logger.ErrorContext(ctx, "falha ao obter dashboard funcionários", "erro", err)
				return nil, fmt.Errorf("%s: falha ao obter dashboard funcionários: %w", op, err)
			}
		case "fornecedores":
			dashboard.Fornecedores, err = s.ObterDashboardFornecedores(ctx, *dataInicio, *dataFim)
			if err != nil {
				s.logger.ErrorContext(ctx, "falha ao obter dashboard fornecedores", "erro", err)
				return nil, fmt.Errorf("%s: falha ao obter dashboard fornecedores: %w", op, err)
			}
		}
	}

	// Log performance final
	totalDuration := time.Since(startTime)
	s.dashLogger.LogDashboardPerformance(ctx, "geral", totalDuration, queryCount, map[string]interface{}{
		"secoes": secoes,
	})

	// Log dados retornados
	s.dashLogger.LogDashboardData(ctx, "geral", "dashboard_completo", 1, false, map[string]interface{}{
		"secoesIncluidas": len(secoes),
		"totalQueries": queryCount,
	})

	s.dashLogger.LogDashboardServiceCall(ctx, "geral", "ObterDashboardCompleto", startTime, nil, map[string]interface{}{
		"success": true,
		"duration": totalDuration.String(),
	})
	
	s.logger.InfoContext(ctx, "dashboard completo obtido com sucesso", "secoes", secoes)
	return dashboard, nil
}

// ObterDashboardFinanceiro retorna apenas a seção financeira do dashboard
func (s *Service) ObterDashboardFinanceiro(ctx context.Context, dataInicio, dataFim time.Time) (*dto.DashboardFinanceiroDTO, error) {
	const op = "service.dashboard.ObterDashboardFinanceiro"

	financeiro := &dto.DashboardFinanceiroDTO{
		UltimaAtualizacao: time.Now(),
	}

	var err error

	// Obter fluxo de caixa
	financeiro.FluxoCaixa, err = s.querier.ObterFluxoCaixaResumo(ctx, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter fluxo de caixa: %w", op, err)
	}

	// Obter distribuição de despesas
	financeiro.DistribuicaoDespesas, err = s.querier.ObterDistribuicaoDespesas(ctx, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter distribuição de despesas: %w", op, err)
	}

	return financeiro, nil
}

// ObterDashboardObras retorna apenas a seção de obras do dashboard
func (s *Service) ObterDashboardObras(ctx context.Context) (*dto.DashboardObrasDTO, error) {
	const op = "service.dashboard.ObterDashboardObras"

	obras := &dto.DashboardObrasDTO{
		UltimaAtualizacao: time.Now(),
	}

	var err error

	// Obter progresso das obras
	obras.Progresso, err = s.querier.ObterProgressoObras(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter progresso das obras: %w", op, err)
	}

	// Obter distribuição das obras
	obras.Distribuicao, err = s.querier.ObterDistribuicaoObras(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter distribuição das obras: %w", op, err)
	}

	// Obter tendências (últimos 6 meses)
	obras.Tendencias, err = s.querier.ObterTendenciasObras(ctx, 6)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter tendências das obras: %w", op, err)
	}

	return obras, nil
}

// ObterDashboardFuncionarios retorna apenas a seção de funcionários do dashboard
func (s *Service) ObterDashboardFuncionarios(ctx context.Context, dataInicio, dataFim time.Time) (*dto.DashboardFuncionariosDTO, error) {
	const op = "service.dashboard.ObterDashboardFuncionarios"

	funcionarios := &dto.DashboardFuncionariosDTO{
		UltimaAtualizacao: time.Now(),
	}

	var err error

	// Obter produtividade
	funcionarios.Produtividade, err = s.querier.ObterProdutividadeFuncionarios(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter produtividade dos funcionários: %w", op, err)
	}

	// Obter custos de mão de obra
	funcionarios.CustosMaoObra, err = s.querier.ObterCustosMaoObra(ctx, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter custos de mão de obra: %w", op, err)
	}

	// Obter top 5 funcionários
	funcionarios.TopFuncionarios, err = s.querier.ObterTopFuncionarios(ctx, 5)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter top funcionários: %w", op, err)
	}

	return funcionarios, nil
}

// ObterDashboardFornecedores retorna apenas a seção de fornecedores do dashboard
func (s *Service) ObterDashboardFornecedores(ctx context.Context, dataInicio, dataFim time.Time) (*dto.DashboardFornecedoresDTO, error) {
	const op = "service.dashboard.ObterDashboardFornecedores"

	fornecedores := &dto.DashboardFornecedoresDTO{
		UltimaAtualizacao: time.Now(),
	}

	var err error

	// Obter fornecedores por categoria
	fornecedores.FornecedoresPorCategoria, err = s.querier.ObterFornecedoresPorCategoria(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter fornecedores por categoria: %w", op, err)
	}

	// Obter top 5 fornecedores
	fornecedores.TopFornecedores, err = s.querier.ObterTopFornecedores(ctx, 5)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter top fornecedores: %w", op, err)
	}

	// Obter gastos com fornecedores (top 10)
	fornecedores.GastosFornecedores, err = s.querier.ObterGastosFornecedores(ctx, dataInicio, dataFim, 10)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter gastos com fornecedores: %w", op, err)
	}

	// Obter estatísticas gerais
	fornecedores.EstatisticasGerais, err = s.querier.ObterEstatisticasGeraisFornecedores(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: falha ao obter estatísticas gerais dos fornecedores: %w", op, err)
	}

	return fornecedores, nil
}

// ObterFluxoCaixa retorna apenas os dados de fluxo de caixa
func (s *Service) ObterFluxoCaixa(ctx context.Context, dataInicio, dataFim time.Time) (*dto.FluxoCaixaResumoDTO, error) {
	const op = "service.dashboard.ObterFluxoCaixa"

	fluxo, err := s.querier.ObterFluxoCaixaResumo(ctx, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return fluxo, nil
}