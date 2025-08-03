package dto

import "time"

// ProgressoObraItemDTO representa o progresso de uma obra individual
type ProgressoObraItemDTO struct {
	ObraID              string  `json:"obraId" db:"obra_id"`
	NomeObra            string  `json:"nomeObra" db:"nome_obra"`
	PercentualConcluido float64 `json:"percentualConcluido" db:"percentual_concluido"`
	EtapasConcluidas    int     `json:"etapasConcluidas" db:"etapas_concluidas"`
	EtapasTotal         int     `json:"etapasTotal" db:"etapas_total"`
	Status              string  `json:"status" db:"status"`
	DataInicio          time.Time `json:"dataInicio" db:"data_inicio"`
	DataFimPrevista     *time.Time `json:"dataFimPrevista" db:"data_fim_prevista"`
}

// ProgressoObrasDTO representa o progresso geral das obras
type ProgressoObrasDTO struct {
	ProgressoMedio      float64                  `json:"progressoMedio"`
	ObrasEmAndamento    int                      `json:"obrasEmAndamento"`
	ObrasConcluidas     int                      `json:"obrasConcluidas"`
	TotalObras          int                      `json:"totalObras"`
	ProgressoPorObra    []*ProgressoObraItemDTO  `json:"progressoPorObra"`
}

// DistribuicaoObraItemDTO representa um item da distribuição por status
type DistribuicaoObraItemDTO struct {
	Status      string  `json:"status" db:"status"`
	Quantidade  int     `json:"quantidade" db:"quantidade"`
	Percentual  float64 `json:"percentual"`
	ValorTotal  float64 `json:"valorTotal" db:"valor_total"`
}

// DistribuicaoObrasDTO representa a distribuição de obras por status/tipo
type DistribuicaoObrasDTO struct {
	TotalObras         int                        `json:"totalObras"`
	DistribuicaoPorStatus []*DistribuicaoObraItemDTO `json:"distribuicaoPorStatus"`
	StatusMaisComum    string                     `json:"statusMaisComum"`
}

// TendenciaObraItemDTO representa dados de tendência por período
type TendenciaObraItemDTO struct {
	Periodo          time.Time `json:"periodo" db:"periodo"`
	ObrasIniciadas   int       `json:"obrasIniciadas" db:"obras_iniciadas"`
	ObrasConcluidas  int       `json:"obrasConcluidas" db:"obras_concluidas"`
	ObrasEmAtraso    int       `json:"obrasEmAtraso" db:"obras_em_atraso"`
}

// TendenciasObrasDTO representa as tendências e análises de prazos
type TendenciasObrasDTO struct {
	ObrasEmAtraso         int                     `json:"obrasEmAtraso"`
	ObrasNoPrazo          int                     `json:"obrasNoPrazo"`
	PercentualAtraso      float64                 `json:"percentualAtraso"`
	TendenciaMensal       []*TendenciaObraItemDTO `json:"tendenciaMensal"`
	PrevisaoConclusaoMes  int                     `json:"previsaoConclusaoMes"`
	TendenciaGeral        string                  `json:"tendenciaGeral"` // "melhorando", "piorando", "estavel"
}

// DashboardObrasDTO agrega todas as informações de obras do dashboard
type DashboardObrasDTO struct {
	Progresso         *ProgressoObrasDTO     `json:"progresso"`
	Distribuicao      *DistribuicaoObrasDTO  `json:"distribuicao"`
	Tendencias        *TendenciasObrasDTO    `json:"tendencias"`
	UltimaAtualizacao time.Time              `json:"ultimaAtualizacao"`
}