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
	PercentualAtraso        float64 `json:"percentualAtraso"`
}

// AlertasDTO representa alertas e notificações importantes
type AlertasDTO struct {
	ObrasComAtraso       []string `json:"obrasComAtraso"`
	FornecedoresInativos []string `json:"fornecedoresInativos"`
	FuncionariosSemApontamento []string `json:"funcionariosSemApontamento"`
	OrcamentosPendentes  int      `json:"orcamentosPendentes"`
	PagamentosPendentes  int      `json:"pagamentosPendentes"`
}

// DashboardGeralDTO representa o dashboard completo da aplicação
type DashboardGeralDTO struct {
	ResumoGeral       *ResumoGeralDTO            `json:"resumoGeral"`
	Alertas           *AlertasDTO                `json:"alertas"`
	Financeiro        *DashboardFinanceiroDTO    `json:"financeiro"`
	Obras             *DashboardObrasDTO         `json:"obras"`
	Funcionarios      *DashboardFuncionariosDTO  `json:"funcionarios"`
	Fornecedores      *DashboardFornecedoresDTO  `json:"fornecedores"`
	UltimaAtualizacao time.Time                  `json:"ultimaAtualizacao"`
	VersaoCache       string                     `json:"versaoCache,omitempty"` // Para controle de cache
}

// ParametrosDashboardDTO representa os parâmetros de filtro para o dashboard
type ParametrosDashboardDTO struct {
	DataInicio    *time.Time `json:"dataInicio,omitempty"`
	DataFim       *time.Time `json:"dataFim,omitempty"`
	ObraIDs       []string   `json:"obraIds,omitempty"`
	FornecedorIDs []string   `json:"fornecedorIds,omitempty"`
	IncluirInativos bool     `json:"incluirInativos"`
	Secoes        []string   `json:"secoes,omitempty"` // ["financeiro", "obras", "funcionarios", "fornecedores"]
}