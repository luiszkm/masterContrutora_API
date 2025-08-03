package dto

import "time"

// FornecedorPorCategoriaItemDTO representa fornecedores por categoria
type FornecedorPorCategoriaItemDTO struct {
	CategoriaID         string  `json:"categoriaId" db:"categoria_id"`
	CategoriaNome       string  `json:"categoriaNome" db:"categoria_nome"`
	QuantidadeFornecedores int  `json:"quantidadeFornecedores" db:"quantidade_fornecedores"`
	Percentual          float64 `json:"percentual"`
	AvaliacaoMedia      float64 `json:"avaliacaoMedia" db:"avaliacao_media"`
}

// FornecedoresPorCategoriaDTO representa a distribuição de fornecedores por categoria
type FornecedoresPorCategoriaDTO struct {
	TotalFornecedores     int                              `json:"totalFornecedores"`
	TotalCategorias       int                              `json:"totalCategorias"`
	DistribuicaoPorCategoria []*FornecedorPorCategoriaItemDTO `json:"distribuicaoPorCategoria"`
	CategoriaMaisPopular  string                           `json:"categoriaMaisPopular"`
	CategoriaComMelhorAvaliacao string                     `json:"categoriaComMelhorAvaliacao"`
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
	Categorias      []string  `json:"categorias"` // Lista de categorias do fornecedor
}

// TopFornecedoresDTO representa o ranking dos melhores fornecedores
type TopFornecedoresDTO struct {
	Top5Fornecedores    []*TopFornecedorDTO `json:"top5Fornecedores"`
	CriterioAvaliacao   string              `json:"criterioAvaliacao"`
	AvaliacaoMedia      float64             `json:"avaliacaoMedia"`
	FornecedoresAtivos  int                 `json:"fornecedoresAtivos"`
}

// GastoFornecedorItemDTO representa gastos com um fornecedor específico
type GastoFornecedorItemDTO struct {
	FornecedorID       string    `json:"fornecedorId" db:"fornecedor_id"`
	NomeFornecedor     string    `json:"nomeFornecedor" db:"nome_fornecedor"`
	Avaliacao          float64   `json:"avaliacao" db:"avaliacao"`
	ValorTotalGasto    float64   `json:"valorTotalGasto" db:"valor_total_gasto"`
	QuantidadeOrcamentos int     `json:"quantidadeOrcamentos" db:"quantidade_orcamentos"`
	TicketMedio        float64   `json:"ticketMedio"` // Calculado
	UltimoOrcamento    *time.Time `json:"ultimoOrcamento" db:"ultimo_orcamento"`
	Percentual         float64   `json:"percentual"` // Percentual do total gasto
}

// GastosFornecedoresDTO representa os gastos com os top fornecedores
type GastosFornecedoresDTO struct {
	TotalGastoFornecedores float64                   `json:"totalGastoFornecedores"`
	GastoMedioFornecedor   float64                   `json:"gastoMedioFornecedor"`
	Top10Gastos            []*GastoFornecedorItemDTO `json:"top10Gastos"`
	FornecedorMaiorGasto   string                    `json:"fornecedorMaiorGasto"`
	ValorMaiorGasto        float64                   `json:"valorMaiorGasto"`
}

// EstatisticasGeraisFornecedoresDTO representa estatísticas gerais
type EstatisticasGeraisFornecedoresDTO struct {
	TotalFornecedores      int     `json:"totalFornecedores"`
	FornecedoresAtivos     int     `json:"fornecedoresAtivos"`
	FornecedoresInativos   int     `json:"fornecedoresInativos"`
	AvaliacaoMediaGeral    float64 `json:"avaliacaoMediaGeral"`
	TempoMedioContrato     int     `json:"tempoMedioContrato"` // Em dias
}

// DashboardFornecedoresDTO agrega todas as informações de fornecedores do dashboard
type DashboardFornecedoresDTO struct {
	FornecedoresPorCategoria *FornecedoresPorCategoriaDTO      `json:"fornecedoresPorCategoria"`
	TopFornecedores          *TopFornecedoresDTO               `json:"topFornecedores"`
	GastosFornecedores       *GastosFornecedoresDTO            `json:"gastosFornecedores"`
	EstatisticasGerais       *EstatisticasGeraisFornecedoresDTO `json:"estatisticasGerais"`
	UltimaAtualizacao        time.Time                         `json:"ultimaAtualizacao"`
}