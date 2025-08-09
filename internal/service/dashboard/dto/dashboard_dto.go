package dto

import "time"

// ResumoGeralDTO representa um resumo executivo dos principais KPIs
type ResumoGeralDTO struct {
	TotalObras              int     `json:"totalObras"`
	ObrasEmAndamento        int     `json:"obrasEmAndamento"`
	TotalFuncionarios       int     `json:"totalFuncionarios"`
	FuncionariosAtivos      int     `json:"funcionariosAtivos"`
	TotalFornecedores       int     `json:"totalFornecedores"`
	FornecedoresAtivos      int     `json:"fornecedoresAtivos"`
	SaldoFinanceiroAtual    float64 `json:"saldoFinanceiroAtual"`
	TotalInvestido          float64 `json:"totalInvestido"`
	ProgressoMedioObras     float64 `json:"progressoMedioObras"`
	ObrasEmAtraso           int     `json:"obrasEmAtraso"`
	PercentualAtraso        float64 `json:"percentualAtraso" db:"-"`
}

// AlertasDTO representa alertas e notificações importantes
type AlertasDTO struct {
	ObrasComAtraso       []string `json:"obrasComAtraso" db:"-"`
	FornecedoresInativos []string `json:"fornecedoresInativos" db:"-"`
	FuncionariosSemApontamento []string `json:"funcionariosSemApontamento" db:"-"`
	OrcamentosPendentes  int      `json:"orcamentosPendentes" db:"-"`
	PagamentosPendentes  int      `json:"pagamentosPendentes" db:"-"`
}

// DashboardGeralDTO representa o dashboard completo da aplicação
type DashboardGeralDTO struct {
	ResumoGeral       *ResumoGeralDTO            `json:"resumoGeral" db:"-"`
	Alertas           *AlertasDTO                `json:"alertas" db:"-"`
	Financeiro        *DashboardFinanceiroDTO    `json:"financeiro" db:"-"`
	Obras             *DashboardObrasDTO         `json:"obras" db:"-"`
	Funcionarios      *DashboardFuncionariosDTO  `json:"funcionarios" db:"-"`
	Fornecedores      *DashboardFornecedoresDTO  `json:"fornecedores" db:"-"`
	UltimaAtualizacao time.Time                  `json:"ultimaAtualizacao" db:"-"`
	VersaoCache       string                     `json:"versaoCache,omitempty" db:"-"` // Para controle de cache
}

// ParametrosDashboardDTO representa os parâmetros de filtro para o dashboard
type ParametrosDashboardDTO struct {
	DataInicio    *time.Time `json:"dataInicio,omitempty" db:"-"`
	DataFim       *time.Time `json:"dataFim,omitempty" db:"-"`
	ObraIDs       []string   `json:"obraIds,omitempty" db:"-"`
	FornecedorIDs []string   `json:"fornecedorIds,omitempty" db:"-"`
	IncluirInativos bool     `json:"incluirInativos" db:"-"`
	Secoes        []string   `json:"secoes,omitempty" db:"-"` // ["financeiro", "obras", "funcionarios", "fornecedores"]
}