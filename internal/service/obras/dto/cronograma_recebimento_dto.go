package dto

import "time"

// CriarCronogramaRecebimentoInput representa o input para criar um cronograma de recebimento
type CriarCronogramaRecebimentoInput struct {
	ObraID         string    `json:"obraId" validate:"required"`
	NumeroEtapa    int       `json:"numeroEtapa" validate:"required,min=1"`
	DescricaoEtapa string    `json:"descricaoEtapa" validate:"required"`
	ValorPrevisto  float64   `json:"valorPrevisto" validate:"required,gt=0"`
	DataVencimento time.Time `json:"dataVencimento" validate:"required"`
}

// AtualizarCronogramaRecebimentoInput representa o input para atualizar um cronograma
type AtualizarCronogramaRecebimentoInput struct {
	DescricaoEtapa string    `json:"descricaoEtapa,omitempty"`
	ValorPrevisto  *float64  `json:"valorPrevisto,omitempty" validate:"omitempty,gt=0"`
	DataVencimento *time.Time `json:"dataVencimento,omitempty"`
}

// RegistrarRecebimentoInput representa o input para registrar um recebimento
type RegistrarRecebimentoInput struct {
	Valor       float64 `json:"valor" validate:"required,gt=0"`
	Observacoes *string `json:"observacoes,omitempty"`
}

// CronogramaRecebimentoOutput representa o output de um cronograma
type CronogramaRecebimentoOutput struct {
	ID                     string     `json:"id"`
	ObraID                 string     `json:"obraId"`
	NumeroEtapa            int        `json:"numeroEtapa"`
	DescricaoEtapa         string     `json:"descricaoEtapa"`
	ValorPrevisto          float64    `json:"valorPrevisto"`
	DataVencimento         time.Time  `json:"dataVencimento"`
	Status                 string     `json:"status"`
	DataRecebimento        *time.Time `json:"dataRecebimento,omitempty"`
	ValorRecebido          float64    `json:"valorRecebido"`
	ValorSaldo             float64    `json:"valorSaldo"`
	PercentualRecebido     float64    `json:"percentualRecebido"`
	ObservacoesRecebimento *string    `json:"observacoesRecebimento,omitempty"`
	EstaVencido            bool       `json:"estaVencido"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}

// CriarCronogramaEmLoteInput permite criar múltiplos cronogramas de uma vez
type CriarCronogramaEmLoteInput struct {
	ObraID             string                            `json:"obraId" validate:"required"`
	Cronogramas        []CriarCronogramaRecebimentoInput `json:"cronogramas" validate:"required,dive"`
	SubstituirExistente bool                             `json:"substituirExistente"` // Se true, remove cronogramas existentes
}

// ResumoFinanceiroObraOutput resumo financeiro de uma obra
type ResumoFinanceiroObraOutput struct {
	ObraID                 string  `json:"obraId"`
	ObraNome               string  `json:"obraNome"`
	ValorContratoTotal     float64 `json:"valorContratoTotal"`
	ValorRecebido          float64 `json:"valorRecebido"`
	ValorSaldo             float64 `json:"valorSaldo"`
	PercentualRecebido     float64 `json:"percentualRecebido"`
	TipoCobranca           string  `json:"tipoCobranca"`
	DataAssinaturaContrato *time.Time `json:"dataAssinaturaContrato,omitempty"`
	
	// Estatísticas do cronograma
	TotalEtapas            int     `json:"totalEtapas"`
	EtapasRecebidas        int     `json:"etapasRecebidas"`
	EtapasVencidas         int     `json:"etapasVencidas"`
	ProximoVencimento      *time.Time `json:"proximoVencimento,omitempty"`
	ValorProximoVencimento float64 `json:"valorProximoVencimento"`
}