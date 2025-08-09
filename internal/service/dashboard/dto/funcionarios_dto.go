package dto

import "time"

// ProdutividadeFuncionarioItemDTO representa a produtividade de um funcionário
type ProdutividadeFuncionarioItemDTO struct {
	FuncionarioID       string  `json:"funcionarioId" db:"funcionario_id"`
	NomeFuncionario     string  `json:"nomeFuncionario" db:"nome_funcionario"`
	Cargo               string  `json:"cargo" db:"cargo"`
	DiasTrabalhados     int     `json:"diasTrabalhados" db:"dias_trabalhados"`
	MediaDiasPorPeriodo float64 `json:"mediaDiasPorPeriodo" db:"media_dias_por_periodo"`
	IndiceProdutividade float64 `json:"indiceProdutividade" db:"-"` // Calculado baseado na média
	ObrasAlocadas       int     `json:"obrasAlocadas" db:"obras_alocadas"`
}

// ProdutividadeFuncionariosDTO representa dados gerais de produtividade
type ProdutividadeFuncionariosDTO struct {
	MediaGeralProdutividade float64                            `json:"mediaGeralProdutividade" db:"-"`
	TotalFuncionarios       int                                `json:"totalFuncionarios" db:"-"`
	FuncionariosAtivos      int                                `json:"funcionariosAtivos" db:"-"`
	ProdutividadePorFuncionario []*ProdutividadeFuncionarioItemDTO `json:"produtividadePorFuncionario" db:"-"`
	Top5Produtivos          []*ProdutividadeFuncionarioItemDTO `json:"top5Produtivos" db:"-"`
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
	CustoTotalMaoObra     float64                       `json:"custoTotalMaoObra" db:"-"`
	CustoMedioFuncionario float64                       `json:"custoMedioFuncionario" db:"-"`
	CustoMedioObra        float64                       `json:"custoMedioObra" db:"-"`
	CustosPorFuncionario  []*CustoMaoObraFuncionarioDTO `json:"custosPorFuncionario" db:"-"`
	CustosPorObra         []*CustoMaoObraPorObraDTO     `json:"custosPorObra" db:"-"`
}

// TopFuncionarioDTO representa um funcionário no ranking de avaliações
type TopFuncionarioDTO struct {
	FuncionarioID       string  `json:"funcionarioId" db:"funcionario_id"`
	NomeFuncionario     string  `json:"nomeFuncionario" db:"nome_funcionario"`
	Cargo               string  `json:"cargo" db:"cargo"`
	AvaliacaoDesempenho string  `json:"avaliacaoDesempenho" db:"avaliacao_desempenho"`
	NotaAvaliacao       float64 `json:"notaAvaliacao" db:"-"` // Calculada a partir da avaliação textual
	DiasTrabalhadosTotal int    `json:"diasTrabalhadosTotal" db:"dias_trabalhados_total"`
	ObrasParticipadas   int     `json:"obrasParticipadas" db:"obras_participadas"`
	DataContratacao     time.Time `json:"dataContratacao" db:"data_contratacao"`
}

// TopFuncionariosDTO representa o ranking dos melhores funcionários
type TopFuncionariosDTO struct {
	Top5Funcionarios []*TopFuncionarioDTO `json:"top5Funcionarios" db:"-"`
	CriterioAvaliacao string              `json:"criterioAvaliacao" db:"-"`
}

// DashboardFuncionariosDTO agrega todas as informações de funcionários do dashboard
type DashboardFuncionariosDTO struct {
	Produtividade     *ProdutividadeFuncionariosDTO `json:"produtividade" db:"-"`
	CustosMaoObra     *CustosMaoObraDTO            `json:"custosMaoObra" db:"-"`
	TopFuncionarios   *TopFuncionariosDTO          `json:"topFuncionarios" db:"-"`
	UltimaAtualizacao time.Time                    `json:"ultimaAtualizacao" db:"-"`
}