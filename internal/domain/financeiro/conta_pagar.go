package financeiro

import (
	"errors"
	"time"
)

// StatusContaPagar representa os possíveis status de uma conta a pagar
const (
	StatusContaPagarPendente  = "PENDENTE"
	StatusContaPagarPago      = "PAGO"
	StatusContaPagarVencido   = "VENCIDO"
	StatusContaPagarParcial   = "PARCIAL"
	StatusContaPagarCancelado = "CANCELADO"
)

// TipoContaPagar representa os tipos de conta a pagar
const (
	TipoContaPagarFornecedor = "FORNECEDOR"
	TipoContaPagarServico    = "SERVICO"
	TipoContaPagarMaterial   = "MATERIAL"
	TipoContaPagarOutros     = "OUTROS"
)

// ContaPagar representa uma conta a pagar no sistema
type ContaPagar struct {
	ID                string     `json:"id"`
	FornecedorID      *string    `json:"fornecedorId,omitempty"`        // Referência ao fornecedor
	ObraID            *string    `json:"obraId,omitempty"`              // Referência à obra (opcional)
	OrcamentoID       *string    `json:"orcamentoId,omitempty"`         // Referência ao orçamento que originou
	FornecedorNome    string     `json:"fornecedorNome"`                // Nome do fornecedor
	TipoContaPagar    string     `json:"tipoContaPagar"`                // FORNECEDOR, SERVICO, MATERIAL, OUTROS
	Descricao         string     `json:"descricao"`                     // Descrição da conta
	ValorOriginal     float64    `json:"valorOriginal"`                 // Valor original
	ValorPago         float64    `json:"valorPago"`                     // Valor já pago
	DataVencimento    time.Time  `json:"dataVencimento"`                // Data de vencimento
	DataPagamento     *time.Time `json:"dataPagamento,omitempty"`       // Data do pagamento
	Status            string     `json:"status"`                        // Status da conta
	FormaPagamento    *string    `json:"formaPagamento,omitempty"`      // Como foi pago
	Observacoes       *string    `json:"observacoes,omitempty"`         // Observações gerais
	NumeroDocumento   *string    `json:"numeroDocumento,omitempty"`     // Número da nota fiscal/documento
	NumeroCompraNF    *string    `json:"numeroCompraNf,omitempty"`      // Número da compra/NF
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

// ValorSaldo retorna o saldo a pagar
func (cp *ContaPagar) ValorSaldo() float64 {
	return cp.ValorOriginal - cp.ValorPago
}

// PercentualPago calcula o percentual já pago
func (cp *ContaPagar) PercentualPago() float64 {
	if cp.ValorOriginal == 0 {
		return 0
	}
	return (cp.ValorPago / cp.ValorOriginal) * 100
}

// EstaVencido verifica se a conta está vencida
func (cp *ContaPagar) EstaVencido() bool {
	return time.Now().After(cp.DataVencimento) && cp.Status != StatusContaPagarPago
}

// DiasVencimento retorna quantos dias está vencida (negativo se ainda não venceu)
func (cp *ContaPagar) DiasVencimento() int {
	diff := time.Now().Sub(cp.DataVencimento)
	return int(diff.Hours() / 24)
}

// RegistrarPagamento registra um pagamento (total ou parcial)
func (cp *ContaPagar) RegistrarPagamento(valor float64, formaPagamento, observacoes *string) error {
	if valor <= 0 {
		return errors.New("valor deve ser positivo")
	}

	novoValorPago := cp.ValorPago + valor
	if novoValorPago > cp.ValorOriginal {
		return errors.New("valor pago não pode exceder o valor original")
	}

	cp.ValorPago = novoValorPago
	now := time.Now()
	cp.DataPagamento = &now
	cp.UpdatedAt = now

	if formaPagamento != nil {
		cp.FormaPagamento = formaPagamento
	}
	
	if observacoes != nil {
		cp.Observacoes = observacoes
	}

	// Atualiza status baseado no valor pago
	if cp.ValorPago >= cp.ValorOriginal {
		cp.Status = StatusContaPagarPago
	} else if cp.ValorPago > 0 {
		cp.Status = StatusContaPagarParcial
	}

	return nil
}

// MarcarComoVencido marca a conta como vencida
func (cp *ContaPagar) MarcarComoVencido() {
	if cp.Status == StatusContaPagarPendente && cp.EstaVencido() {
		cp.Status = StatusContaPagarVencido
		cp.UpdatedAt = time.Now()
	}
}

// Cancelar cancela a conta a pagar
func (cp *ContaPagar) Cancelar(motivo *string) error {
	if cp.Status == StatusContaPagarPago {
		return errors.New("não é possível cancelar uma conta já paga")
	}
	
	cp.Status = StatusContaPagarCancelado
	cp.UpdatedAt = time.Now()
	
	if motivo != nil {
		cp.Observacoes = motivo
	}
	
	return nil
}

// Validar valida os dados da conta a pagar
func (cp *ContaPagar) Validar() error {
	if cp.FornecedorNome == "" {
		return errors.New("fornecedorNome é obrigatório")
	}
	if cp.TipoContaPagar == "" {
		return errors.New("tipoContaPagar é obrigatório")
	}
	if cp.TipoContaPagar != TipoContaPagarFornecedor && 
		cp.TipoContaPagar != TipoContaPagarServico && 
		cp.TipoContaPagar != TipoContaPagarMaterial &&
		cp.TipoContaPagar != TipoContaPagarOutros {
		return errors.New("tipoContaPagar deve ser FORNECEDOR, SERVICO, MATERIAL ou OUTROS")
	}
	if cp.Descricao == "" {
		return errors.New("descrição é obrigatória")
	}
	if cp.ValorOriginal <= 0 {
		return errors.New("valorOriginal deve ser positivo")
	}
	if cp.DataVencimento.IsZero() {
		return errors.New("dataVencimento é obrigatória")
	}
	return nil
}

// ParcelaContaPagar representa uma parcela de uma conta a pagar
type ParcelaContaPagar struct {
	ID              string     `json:"id"`
	ContaPagarID    string     `json:"contaPagarId"`
	NumeroParcela   int        `json:"numeroParcela"`
	ValorParcela    float64    `json:"valorParcela"`
	DataVencimento  time.Time  `json:"dataVencimento"`
	DataPagamento   *time.Time `json:"dataPagamento,omitempty"`
	ValorPago       float64    `json:"valorPago"`
	Status          string     `json:"status"`
	FormaPagamento  *string    `json:"formaPagamento,omitempty"`
	Observacoes     *string    `json:"observacoes,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// ValorSaldoParcela retorna o saldo desta parcela
func (p *ParcelaContaPagar) ValorSaldoParcela() float64 {
	return p.ValorParcela - p.ValorPago
}

// EstaVencida verifica se a parcela está vencida
func (p *ParcelaContaPagar) EstaVencida() bool {
	return time.Now().After(p.DataVencimento) && p.Status != StatusContaPagarPago
}

// RegistrarPagamentoParcela registra pagamento da parcela
func (p *ParcelaContaPagar) RegistrarPagamentoParcela(valor float64, formaPagamento, observacoes *string) error {
	if valor <= 0 {
		return errors.New("valor deve ser positivo")
	}

	novoValorPago := p.ValorPago + valor
	if novoValorPago > p.ValorParcela {
		return errors.New("valor pago não pode exceder o valor da parcela")
	}

	p.ValorPago = novoValorPago
	now := time.Now()
	p.DataPagamento = &now
	p.UpdatedAt = now

	if formaPagamento != nil {
		p.FormaPagamento = formaPagamento
	}
	
	if observacoes != nil {
		p.Observacoes = observacoes
	}

	// Atualizar status
	if p.ValorPago >= p.ValorParcela {
		p.Status = StatusContaPagarPago
	} else if p.ValorPago > 0 {
		p.Status = StatusContaPagarParcial
	}

	return nil
}