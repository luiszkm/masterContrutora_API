package dto

import "time"

// CriarContaPagarInput representa o input para criar uma conta a pagar
type CriarContaPagarInput struct {
	FornecedorID    *string    `json:"fornecedorId,omitempty"`
	ObraID          *string    `json:"obraId,omitempty"`
	OrcamentoID     *string    `json:"orcamentoId,omitempty"`
	FornecedorNome  string     `json:"fornecedorNome" validate:"required"`
	TipoContaPagar  string     `json:"tipoContaPagar" validate:"required,oneof=FORNECEDOR SERVICO MATERIAL OUTROS"`
	Categoria       string     `json:"categoria" validate:"required,oneof=ORCAMENTO APONTAMENTO MANUAL OUTROS"`
	Descricao       string     `json:"descricao" validate:"required"`
	ValorOriginal   float64    `json:"valorOriginal" validate:"required,gt=0"`
	DataVencimento  time.Time  `json:"dataVencimento" validate:"required"`
	NumeroDocumento *string    `json:"numeroDocumento,omitempty"`
	NumeroCompraNF  *string    `json:"numeroCompraNf,omitempty"`
	Observacoes     *string    `json:"observacoes,omitempty"`
}

// AtualizarContaPagarInput representa o input para atualizar uma conta a pagar
type AtualizarContaPagarInput struct {
	FornecedorNome *string    `json:"fornecedorNome,omitempty"`
	TipoContaPagar *string    `json:"tipoContaPagar,omitempty" validate:"omitempty,oneof=FORNECEDOR SERVICO MATERIAL OUTROS"`
	Categoria      *string    `json:"categoria,omitempty" validate:"omitempty,oneof=ORCAMENTO APONTAMENTO MANUAL OUTROS"`
	Descricao      *string    `json:"descricao,omitempty"`
	ValorOriginal  *float64   `json:"valorOriginal,omitempty" validate:"omitempty,gt=0"`
	DataVencimento *time.Time `json:"dataVencimento,omitempty"`
	NumeroDocumento *string   `json:"numeroDocumento,omitempty"`
	NumeroCompraNF *string    `json:"numeroCompraNf,omitempty"`
	Observacoes    *string    `json:"observacoes,omitempty"`
}

// RegistrarPagamentoContaPagarInput representa o input para registrar um pagamento
type RegistrarPagamentoContaPagarInput struct {
	Valor           float64 `json:"valor" validate:"required,gt=0"`
	FormaPagamento  *string `json:"formaPagamento,omitempty"`
	ContaBancariaID *string `json:"contaBancariaId,omitempty"`
	Observacoes     *string `json:"observacoes,omitempty"`
}

// ContaPagarOutput representa o output de uma conta a pagar
type ContaPagarOutput struct {
	ID              string     `json:"id"`
	FornecedorID    *string    `json:"fornecedorId,omitempty"`
	ObraID          *string    `json:"obraId,omitempty"`
	OrcamentoID     *string    `json:"orcamentoId,omitempty"`
	FornecedorNome  string     `json:"fornecedorNome"`
	TipoContaPagar  string     `json:"tipoContaPagar"`
	Categoria       string     `json:"categoria"`
	Descricao       string     `json:"descricao"`
	ValorOriginal   float64    `json:"valorOriginal"`
	ValorPago       float64    `json:"valorPago"`
	ValorSaldo      float64    `json:"valorSaldo"`
	PercentualPago  float64    `json:"percentualPago"`
	DataVencimento  time.Time  `json:"dataVencimento"`
	DataPagamento   *time.Time `json:"dataPagamento,omitempty"`
	Status          string     `json:"status"`
	FormaPagamento  *string    `json:"formaPagamento,omitempty"`
	Observacoes     *string    `json:"observacoes,omitempty"`
	NumeroDocumento *string    `json:"numeroDocumento,omitempty"`
	NumeroCompraNF  *string    `json:"numeroCompraNf,omitempty"`
	EstaVencido     bool       `json:"estaVencido"`
	DiasVencimento  int        `json:"diasVencimento"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// FiltrosContaPagarInput representa filtros para listagem
type FiltrosContaPagarInput struct {
	FornecedorNome       *string    `json:"fornecedorNome,omitempty"`
	Status               *string    `json:"status,omitempty"`
	TipoContaPagar       *string    `json:"tipoContaPagar,omitempty"`
	ObraID               *string    `json:"obraId,omitempty"`
	FornecedorID         *string    `json:"fornecedorId,omitempty"`
	DataVencimentoInicio *time.Time `json:"dataVencimentoInicio,omitempty"`
	DataVencimentoFim    *time.Time `json:"dataVencimentoFim,omitempty"`
	ApenasVencidas       bool       `json:"apenasVencidas"`
}

// ResumoContasPagarOutput resumo das contas a pagar
type ResumoContasPagarOutput struct {
	TotalContas        int     `json:"totalContas"`
	TotalValorOriginal float64 `json:"totalValorOriginal"`
	TotalValorPago     float64 `json:"totalValorPago"`
	TotalValorSaldo    float64 `json:"totalValorSaldo"`
	ContasPendentes    int     `json:"contasPendentes"`
	ContasVencidas     int     `json:"contasVencidas"`
	ContasPagas        int     `json:"contasPagas"`
	PercentualPago     float64 `json:"percentualPago"`
}

// CriarContaPagarDeOrcamentoInput para criação automática a partir de orçamento
type CriarContaPagarDeOrcamentoInput struct {
	OrcamentoID       string     `json:"orcamentoId" validate:"required"`
	DataVencimento    time.Time  `json:"dataVencimento" validate:"required"`
	NumeroDocumento   *string    `json:"numeroDocumento,omitempty"`
	NumeroCompraNF    *string    `json:"numeroCompraNf,omitempty"`
	Observacoes       *string    `json:"observacoes,omitempty"`
	DividirParcelas   bool       `json:"dividirParcelas"`
	QuantidadeParcelas *int      `json:"quantidadeParcelas,omitempty" validate:"omitempty,min=1,max=60"`
}

// ParcelaContaPagarOutput representa o output de uma parcela
type ParcelaContaPagarOutput struct {
	ID             string     `json:"id"`
	ContaPagarID   string     `json:"contaPagarId"`
	NumeroParcela  int        `json:"numeroParcela"`
	ValorParcela   float64    `json:"valorParcela"`
	DataVencimento time.Time  `json:"dataVencimento"`
	DataPagamento  *time.Time `json:"dataPagamento,omitempty"`
	ValorPago      float64    `json:"valorPago"`
	ValorSaldo     float64    `json:"valorSaldo"`
	Status         string     `json:"status"`
	FormaPagamento *string    `json:"formaPagamento,omitempty"`
	Observacoes    *string    `json:"observacoes,omitempty"`
	EstaVencida    bool       `json:"estaVencida"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}