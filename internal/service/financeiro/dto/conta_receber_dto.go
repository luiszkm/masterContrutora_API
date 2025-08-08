package dto

import "time"

// CriarContaReceberInput representa o input para criar uma conta a receber
type CriarContaReceberInput struct {
	ObraID                  *string    `json:"obraId,omitempty"`
	CronogramaRecebimentoID *string    `json:"cronogramaRecebimentoId,omitempty"`
	Cliente                 string     `json:"cliente" validate:"required"`
	TipoContaReceber        string     `json:"tipoContaReceber" validate:"required,oneof=OBRA SERVICO OUTROS"`
	Descricao               string     `json:"descricao" validate:"required"`
	ValorOriginal           float64    `json:"valorOriginal" validate:"required,gt=0"`
	DataVencimento          time.Time  `json:"dataVencimento" validate:"required"`
	NumeroDocumento         *string    `json:"numeroDocumento,omitempty"`
}

// AtualizarContaReceberInput representa o input para atualizar uma conta a receber
type AtualizarContaReceberInput struct {
	Cliente         *string    `json:"cliente,omitempty"`
	TipoContaReceber *string   `json:"tipoContaReceber,omitempty" validate:"omitempty,oneof=OBRA SERVICO OUTROS"`
	Descricao       *string    `json:"descricao,omitempty"`
	ValorOriginal   *float64   `json:"valorOriginal,omitempty" validate:"omitempty,gt=0"`
	DataVencimento  *time.Time `json:"dataVencimento,omitempty"`
	NumeroDocumento *string    `json:"numeroDocumento,omitempty"`
	Observacoes     *string    `json:"observacoes,omitempty"`
}

// RegistrarRecebimentoContaInput representa o input para registrar um recebimento
type RegistrarRecebimentoContaInput struct {
	Valor           float64 `json:"valor" validate:"required,gt=0"`
	FormaPagamento  *string `json:"formaPagamento,omitempty"`
	ContaBancariaID *string `json:"contaBancariaId,omitempty"`
	Observacoes     *string `json:"observacoes,omitempty"`
}

// ContaReceberOutput representa o output de uma conta a receber
type ContaReceberOutput struct {
	ID                      string     `json:"id"`
	ObraID                  *string    `json:"obraId,omitempty"`
	CronogramaRecebimentoID *string    `json:"cronogramaRecebimentoId,omitempty"`
	Cliente                 string     `json:"cliente"`
	TipoContaReceber        string     `json:"tipoContaReceber"`
	Descricao               string     `json:"descricao"`
	ValorOriginal           float64    `json:"valorOriginal"`
	ValorRecebido           float64    `json:"valorRecebido"`
	ValorSaldo              float64    `json:"valorSaldo"`
	PercentualRecebido      float64    `json:"percentualRecebido"`
	DataVencimento          time.Time  `json:"dataVencimento"`
	DataRecebimento         *time.Time `json:"dataRecebimento,omitempty"`
	Status                  string     `json:"status"`
	FormaPagamento          *string    `json:"formaPagamento,omitempty"`
	Observacoes             *string    `json:"observacoes,omitempty"`
	NumeroDocumento         *string    `json:"numeroDocumento,omitempty"`
	EstaVencido             bool       `json:"estaVencido"`
	DiasVencimento          int        `json:"diasVencimento"`
	CreatedAt               time.Time  `json:"createdAt"`
	UpdatedAt               time.Time  `json:"updatedAt"`
}

// FiltrosContaReceberInput representa filtros para listagem
type FiltrosContaReceberInput struct {
	Cliente          *string    `json:"cliente,omitempty"`
	Status           *string    `json:"status,omitempty"`
	TipoContaReceber *string    `json:"tipoContaReceber,omitempty"`
	ObraID           *string    `json:"obraId,omitempty"`
	DataVencimentoInicio *time.Time `json:"dataVencimentoInicio,omitempty"`
	DataVencimentoFim    *time.Time `json:"dataVencimentoFim,omitempty"`
	ApenasVencidas   bool       `json:"apenasVencidas"`
}

// ResumoContasReceberOutput resumo das contas a receber
type ResumoContasReceberOutput struct {
	TotalContas           int     `json:"totalContas"`
	TotalValorOriginal    float64 `json:"totalValorOriginal"`
	TotalValorRecebido    float64 `json:"totalValorRecebido"`
	TotalValorSaldo       float64 `json:"totalValorSaldo"`
	ContasPendentes       int     `json:"contasPendentes"`
	ContasVencidas        int     `json:"contasVencidas"`
	ContasRecebidas       int     `json:"contasRecebidas"`
	PercentualRecebimento float64 `json:"percentualRecebimento"`
}