package dto

import "time"

// ProdutividadeFuncionarioItemDTO representa a produtividade de um funcionário
type ProdutividadeFuncionarioItemDTO struct {
	FuncionarioID       string  `json:"funcionarioId" db:"funcionario_id"`
	NomeFuncionario     string  `json:"nomeFuncionario" db:"nome_funcionario"`
	Cargo               string  `json:"cargo" db:"cargo"`
	DiasTrabalhados     int     `json:"diasTrabalhados" db:"dias_trabalhados"`
	MediaDiasPorPeriodo float64 `json:"mediaDiasPorPeriodo" db:"media_dias_por_periodo"`
	IndiceProdutividade float64 `json:"indiceProdutividade"` // Calculado baseado na média
	ObrasAlocadas       int     `json:"obrasAlocadas" db:"obras_alocadas"`
}

// ProdutividadeFuncionariosDTO representa dados gerais de produtividade
type ProdutividadeFuncionariosDTO struct {
	MediaGeralProdutividade float64                            `json:"mediaGeralProdutividade"`
	TotalFuncionarios       int                                `json:"totalFuncionarios"`
	FuncionariosAtivos      int                                `json:"funcionariosAtivos"`
	ProdutividadePorFuncionario []*ProdutividadeFuncionarioItemDTO `json:"produtividadePorFuncionario"`
	Top5Produtivos          []*ProdutividadeFuncionarioItemDTO `json:"top5Produtivos"`
}

// CustoMaoObraFuncionarioDTO representa custos por funcionário
type CustoMaoObraFuncionarioDTO struct {
	FuncionarioID   string  `json:"funcionarioId" db:"funcionario_id"`
	NomeFuncionario string  `json:"nomeFuncionario" db:"nome_funcionario"`
	Cargo           string  `json:"cargo" db:"cargo"`
	CustoTotal      float64 `json:"custoTotal" db:"custo_total"`
	CustoMedio      float64 `json:"custoMedio" db:"custo_medio"`
	ValorDiaria     float64 `json:"valorDiaria" db:"valor_diaria"`
	PeriodosTrabalho int    `json:"periodosTrabalho" db:"periodos_trabalho"`
}

// CustoMaoObraPorObraDTO representa custos de mão de obra por obra
type CustoMaoObraPorObraDTO struct {
	ObraID      string  `json:"obraId" db:"obra_id"`
	NomeObra    string  `json:"nomeObra" db:"nome_obra"`
	CustoTotal  float64 `json:"custoTotal" db:"custo_total"`
	CustoMedio  float64 `json:"custoMedio" db:"custo_medio"`
	NumFuncionarios int `json:"numFuncionarios" db:"num_funcionarios"`
}

// CustosMaoObraDTO representa dados gerais de custos de mão de obra
type CustosMaoObraDTO struct {
	CustoTotalMaoObra     float64                       `json:"custoTotalMaoObra"`
	CustoMedioFuncionario float64                       `json:"custoMedioFuncionario"`
	CustoMedioObra        float64                       `json:"custoMedioObra"`
	CustosPorFuncionario  []*CustoMaoObraFuncionarioDTO `json:"custosPorFuncionario"`
	CustosPorObra         []*CustoMaoObraPorObraDTO     `json:"custosPorObra"`
}

// TopFuncionarioDTO representa um funcionário no ranking de avaliações
type TopFuncionarioDTO struct {
	FuncionarioID       string  `json:"funcionarioId" db:"funcionario_id"`
	NomeFuncionario     string  `json:"nomeFuncionario" db:"nome_funcionario"`
	Cargo               string  `json:"cargo" db:"cargo"`
	AvaliacaoDesempenho string  `json:"avaliacaoDesempenho" db:"avaliacao_desempenho"`
	NotaAvaliacao       float64 `json:"notaAvaliacao"` // Calculada a partir da avaliação textual
	DiasTrabalhadosTotal int    `json:"diasTrabalhadosTotal" db:"dias_trabalhados_total"`
	ObrasParticipadas   int     `json:"obrasParticipadas" db:"obras_participadas"`
	DataContratacao     time.Time `json:"dataContratacao" db:"data_contratacao"`
}

// TopFuncionariosDTO representa o ranking dos melhores funcionários
type TopFuncionariosDTO struct {
	Top5Funcionarios []*TopFuncionarioDTO `json:"top5Funcionarios"`
	CriterioAvaliacao string              `json:"criterioAvaliacao"`
}

// DashboardFuncionariosDTO agrega todas as informações de funcionários do dashboard
type DashboardFuncionariosDTO struct {
	Produtividade     *ProdutividadeFuncionariosDTO `json:"produtividade"`
	CustosMaoObra     *CustosMaoObraDTO            `json:"custosMaoObra"`
	TopFuncionarios   *TopFuncionariosDTO          `json:"topFuncionarios"`
	UltimaAtualizacao time.Time                    `json:"ultimaAtualizacao"`
}