package dto

import "time"

// FluxoCaixaDTO representa os dados de fluxo de caixa por período
type FluxoCaixaDTO struct {
	Periodo      time.Time `json:"periodo" db:"periodo"`
	Entradas     float64   `json:"entradas" db:"entradas"`
	Saidas       float64   `json:"saidas" db:"saidas"`
	SaldoLiquido float64   `json:"saldoLiquido" db:"saldo_liquido"`
}

// FluxoCaixaResumoDTO representa o resumo geral do fluxo de caixa
type FluxoCaixaResumoDTO struct {
	TotalEntradas     float64            `json:"totalEntradas"`
	TotalSaidas       float64            `json:"totalSaidas"`
	SaldoAtual        float64            `json:"saldoAtual"`
	FluxoPorPeriodo   []*FluxoCaixaDTO   `json:"fluxoPorPeriodo"`
	TendenciaMensal   string             `json:"tendenciaMensal"` // "crescente", "decrescente", "estavel"
}

// DistribuicaoDespesasItemDTO representa um item da distribuição de despesas
type DistribuicaoDespesasItemDTO struct {
	Categoria      string  `json:"categoria" db:"categoria"`
	Valor          float64 `json:"valor" db:"valor"`
	Percentual     float64 `json:"percentual"`
	QuantidadeItens int    `json:"quantidadeItens" db:"quantidade_itens"`
}

// DistribuicaoDespesasDTO representa a distribuição completa de despesas
type DistribuicaoDespesasDTO struct {
	TotalGasto       float64                        `json:"totalGasto"`
	Distribuicao     []*DistribuicaoDespesasItemDTO `json:"distribuicao"`
	MaiorCategoria   string                         `json:"maiorCategoria"`
	ValorMaiorCategoria float64                    `json:"valorMaiorCategoria"`
}

// DashboardFinanceiroDTO agrega todas as informações financeiras do dashboard
type DashboardFinanceiroDTO struct {
	FluxoCaixa           *FluxoCaixaResumoDTO     `json:"fluxoCaixa"`
	DistribuicaoDespesas *DistribuicaoDespesasDTO `json:"distribuicaoDespesas"`
	UltimaAtualizacao    time.Time                `json:"ultimaAtualizacao"`
}