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
	TotalEntradas     float64            `json:"totalEntradas" db:"-"`
	TotalSaidas       float64            `json:"totalSaidas" db:"-"`
	SaldoAtual        float64            `json:"saldoAtual" db:"-"`
	FluxoPorPeriodo   []*FluxoCaixaDTO   `json:"fluxoPorPeriodo" db:"-"`
	TendenciaMensal   string             `json:"tendenciaMensal" db:"-"` // "crescente", "decrescente", "estavel"
}

// DistribuicaoDespesasItemDTO representa um item da distribuição de despesas
type DistribuicaoDespesasItemDTO struct {
	Categoria      string  `json:"categoria" db:"categoria"`
	Valor          float64 `json:"valor" db:"valor"`
	Percentual     float64 `json:"percentual" db:"-"`
	QuantidadeItens int    `json:"quantidadeItens" db:"quantidade_itens"`
}

// DistribuicaoDespesasDTO representa a distribuição completa de despesas
type DistribuicaoDespesasDTO struct {
	TotalGasto       float64                        `json:"totalGasto" db:"-"`
	Distribuicao     []*DistribuicaoDespesasItemDTO `json:"distribuicao" db:"-"`
	MaiorCategoria   string                         `json:"maiorCategoria" db:"-"`
	ValorMaiorCategoria float64                    `json:"valorMaiorCategoria" db:"-"`
}

// DashboardFinanceiroDTO agrega todas as informações financeiras do dashboard
type DashboardFinanceiroDTO struct {
	FluxoCaixa           *FluxoCaixaResumoDTO     `json:"fluxoCaixa" db:"-"`
	DistribuicaoDespesas *DistribuicaoDespesasDTO `json:"distribuicaoDespesas" db:"-"`
	UltimaAtualizacao    time.Time                `json:"ultimaAtualizacao" db:"-"`
}