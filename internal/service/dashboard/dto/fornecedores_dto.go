package dto

import "time"

// FornecedorPorCategoriaItemDTO representa fornecedores por categoria
type FornecedorPorCategoriaItemDTO struct {
	CategoriaID         string  `json:"categoriaId" db:"categoria_id"`
	CategoriaNome       string  `json:"categoriaNome" db:"categoria_nome"`
	QuantidadeFornecedores int  `json:"quantidadeFornecedores" db:"quantidade_fornecedores"`
	Percentual          float64 `json:"percentual" db:"-"`
	AvaliacaoMedia      float64 `json:"avaliacaoMedia" db:"avaliacao_media"`
}

// FornecedoresPorCategoriaDTO representa a distribuição de fornecedores por categoria
type FornecedoresPorCategoriaDTO struct {
	TotalFornecedores     int                              `json:"totalFornecedores" db:"-"`
	TotalCategorias       int                              `json:"totalCategorias" db:"-"`
	DistribuicaoPorCategoria []*FornecedorPorCategoriaItemDTO `json:"distribuicaoPorCategoria" db:"-"`
	CategoriaMaisPopular  string                           `json:"categoriaMaisPopular" db:"-"`
	CategoriaComMelhorAvaliacao string                     `json:"categoriaComMelhorAvaliacao" db:"-"`
}

// TopFornecedorDTO representa um fornecedor no ranking de avaliações
type TopFornecedorDTO struct {
	FornecedorID    string    `json:"fornecedorId" db:"fornecedor_id"`
	NomeFornecedor  string    `json:"nomeFornecedor" db:"nome_fornecedor"`
	CNPJ            string    `json:"cnpj" db:"cnpj"`
	Avaliacao       float64   `json:"avaliacao" db:"avaliacao"`
	Status          string    `json:"status" db:"status"`
	TotalOrcamentos int       `json:"totalOrcamentos" db:"total_orcamentos"`
	ValorTotalGasto float64   `json:"valorTotalGasto" db:"valor_total_gasto"`
	UltimoOrcamento *time.Time `json:"ultimoOrcamento" db:"ultimo_orcamento"`
	Categorias      []string  `json:"categorias" db:"-"` // Lista de categorias do fornecedor
}

// TopFornecedoresDTO representa o ranking dos melhores fornecedores
type TopFornecedoresDTO struct {
	Top5Fornecedores    []*TopFornecedorDTO `json:"top5Fornecedores" db:"-"`
	CriterioAvaliacao   string              `json:"criterioAvaliacao" db:"-"`
	AvaliacaoMedia      float64             `json:"avaliacaoMedia" db:"-"`
	FornecedoresAtivos  int                 `json:"fornecedoresAtivos" db:"-"`
}

// GastoFornecedorItemDTO representa gastos com um fornecedor específico
type GastoFornecedorItemDTO struct {
	FornecedorID       string    `json:"fornecedorId" db:"fornecedor_id"`
	NomeFornecedor     string    `json:"nomeFornecedor" db:"nome_fornecedor"`
	Avaliacao          float64   `json:"avaliacao" db:"avaliacao"`
	ValorTotalGasto    float64   `json:"valorTotalGasto" db:"valor_total_gasto"`
	QuantidadeOrcamentos int     `json:"quantidadeOrcamentos" db:"quantidade_orcamentos"`
	TicketMedio        float64   `json:"ticketMedio" db:"-"` // Calculado
	UltimoOrcamento    *time.Time `json:"ultimoOrcamento" db:"ultimo_orcamento"`
	Percentual         float64   `json:"percentual" db:"-"` // Percentual do total gasto
}

// GastosFornecedoresDTO representa os gastos com os top fornecedores
type GastosFornecedoresDTO struct {
	TotalGastoFornecedores float64                   `json:"totalGastoFornecedores" db:"-"`
	GastoMedioFornecedor   float64                   `json:"gastoMedioFornecedor" db:"-"`
	Top10Gastos            []*GastoFornecedorItemDTO `json:"top10Gastos" db:"-"`
	FornecedorMaiorGasto   string                    `json:"fornecedorMaiorGasto" db:"-"`
	ValorMaiorGasto        float64                   `json:"valorMaiorGasto" db:"-"`
}

// EstatisticasGeraisFornecedoresDTO representa estatísticas gerais
type EstatisticasGeraisFornecedoresDTO struct {
	TotalFornecedores      int     `json:"totalFornecedores" db:"-"`
	FornecedoresAtivos     int     `json:"fornecedoresAtivos" db:"-"`
	FornecedoresInativos   int     `json:"fornecedoresInativos" db:"-"`
	AvaliacaoMediaGeral    float64 `json:"avaliacaoMediaGeral" db:"-"`
	TempoMedioContrato     int     `json:"tempoMedioContrato" db:"-"` // Em dias
}

// DashboardFornecedoresDTO agrega todas as informações de fornecedores do dashboard
type DashboardFornecedoresDTO struct {
	FornecedoresPorCategoria *FornecedoresPorCategoriaDTO      `json:"fornecedoresPorCategoria" db:"-"`
	TopFornecedores          *TopFornecedoresDTO               `json:"topFornecedores" db:"-"`
	GastosFornecedores       *GastosFornecedoresDTO            `json:"gastosFornecedores" db:"-"`
	EstatisticasGerais       *EstatisticasGeraisFornecedoresDTO `json:"estatisticasGerais" db:"-"`
	UltimaAtualizacao        time.Time                         `json:"ultimaAtualizacao" db:"-"`
}