package dashboard

import (
	"context"
	"time"

	"github.com/luiszkm/masterCostrutora/internal/service/dashboard/dto"
)

// Querier define as operações de consulta para o dashboard
type Querier interface {
	// Seção Financeira
	ObterFluxoCaixa(ctx context.Context, dataInicio, dataFim time.Time) ([]*dto.FluxoCaixaDTO, error)
	ObterFluxoCaixaResumo(ctx context.Context, dataInicio, dataFim time.Time) (*dto.FluxoCaixaResumoDTO, error)
	ObterDistribuicaoDespesas(ctx context.Context, dataInicio, dataFim time.Time) (*dto.DistribuicaoDespesasDTO, error)

	// Seção Obras
	ObterProgressoObras(ctx context.Context) (*dto.ProgressoObrasDTO, error)
	ObterDistribuicaoObras(ctx context.Context) (*dto.DistribuicaoObrasDTO, error)
	ObterTendenciasObras(ctx context.Context, mesesAtras int) (*dto.TendenciasObrasDTO, error)

	// Seção Funcionários
	ObterProdutividadeFuncionarios(ctx context.Context) (*dto.ProdutividadeFuncionariosDTO, error)
	ObterCustosMaoObra(ctx context.Context, dataInicio, dataFim time.Time) (*dto.CustosMaoObraDTO, error)
	ObterTopFuncionarios(ctx context.Context, limite int) (*dto.TopFuncionariosDTO, error)

	// Seção Fornecedores
	ObterFornecedoresPorCategoria(ctx context.Context) (*dto.FornecedoresPorCategoriaDTO, error)
	ObterTopFornecedores(ctx context.Context, limite int) (*dto.TopFornecedoresDTO, error)
	ObterGastosFornecedores(ctx context.Context, dataInicio, dataFim time.Time, limite int) (*dto.GastosFornecedoresDTO, error)
	ObterEstatisticasGeraisFornecedores(ctx context.Context) (*dto.EstatisticasGeraisFornecedoresDTO, error)

	// Resumo Geral e Alertas
	ObterResumoGeral(ctx context.Context) (*dto.ResumoGeralDTO, error)
	ObterAlertas(ctx context.Context) (*dto.AlertasDTO, error)
}