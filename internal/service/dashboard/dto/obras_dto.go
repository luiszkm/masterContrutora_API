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
	ProgressoMedio      float64                  `json:"progressoMedio" db:"-"`
	ObrasEmAndamento    int                      `json:"obrasEmAndamento" db:"-"`
	ObrasConcluidas     int                      `json:"obrasConcluidas" db:"-"`
	TotalObras          int                      `json:"totalObras" db:"-"`
	ProgressoPorObra    []*ProgressoObraItemDTO  `json:"progressoPorObra" db:"-"`
}

// DistribuicaoObraItemDTO representa um item da distribuição por status
type DistribuicaoObraItemDTO struct {
	Status      string  `json:"status" db:"status"`
	Quantidade  int     `json:"quantidade" db:"quantidade"`
	Percentual  float64 `json:"percentual" db:"-"`
	ValorTotal  float64 `json:"valorTotal" db:"valor_total"`
}

// DistribuicaoObrasDTO representa a distribuição de obras por status/tipo
type DistribuicaoObrasDTO struct {
	TotalObras         int                        `json:"totalObras" db:"-"`
	DistribuicaoPorStatus []*DistribuicaoObraItemDTO `json:"distribuicaoPorStatus" db:"-"`
	StatusMaisComum    string                     `json:"statusMaisComum" db:"-"`
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
	ObrasEmAtraso         int                     `json:"obrasEmAtraso" db:"-"`
	ObrasNoPrazo          int                     `json:"obrasNoPrazo" db:"-"`
	PercentualAtraso      float64                 `json:"percentualAtraso" db:"-"`
	TendenciaMensal       []*TendenciaObraItemDTO `json:"tendenciaMensal" db:"-"`
	PrevisaoConclusaoMes  int                     `json:"previsaoConclusaoMes" db:"-"`
	TendenciaGeral        string                  `json:"tendenciaGeral" db:"-"` // "melhorando", "piorando", "estavel"
}

// DashboardObrasDTO agrega todas as informações de obras do dashboard
type DashboardObrasDTO struct {
	Progresso         *ProgressoObrasDTO     `json:"progresso" db:"-"`
	Distribuicao      *DistribuicaoObrasDTO  `json:"distribuicao" db:"-"`
	Tendencias        *TendenciasObrasDTO    `json:"tendencias" db:"-"`
	UltimaAtualizacao time.Time              `json:"ultimaAtualizacao" db:"-"`
}