package financeiro

import (
	"errors"
	"time"
)

// StatusContaReceber representa os possíveis status de uma conta a receber
const (
	StatusContaReceberPendente = "PENDENTE"
	StatusContaReceberRecebido = "RECEBIDO"
	StatusContaReceberVencido  = "VENCIDO"
	StatusContaReceberParcial  = "PARCIAL"
	StatusContaReceberCancelado = "CANCELADO"
)

// TipoContaReceber representa os tipos de conta a receber
const (
	TipoContaReceberObra      = "OBRA"
	TipoContaReceberServico   = "SERVICO"
	TipoContaReceberOutros    = "OUTROS"
)

// ContaReceber representa uma conta a receber no sistema
type ContaReceber struct {
	ID                      string     `json:"id"`
	ObraID                  *string    `json:"obraId,omitempty"`              // Referência à obra (opcional)
	CronogramaRecebimentoID *string    `json:"cronogramaRecebimentoId,omitempty"` // Referência ao cronograma
	Cliente                 string     `json:"cliente"`                       // Nome do cliente
	TipoContaReceber        string     `json:"tipoContaReceber"`              // OBRA, SERVICO, OUTROS
	Descricao               string     `json:"descricao"`                     // Descrição da conta
	ValorOriginal           float64    `json:"valorOriginal"`                 // Valor original
	ValorRecebido           float64    `json:"valorRecebido"`                 // Valor já recebido
	DataVencimento          time.Time  `json:"dataVencimento"`                // Data de vencimento
	DataRecebimento         *time.Time `json:"dataRecebimento,omitempty"`     // Data do recebimento
	Status                  string     `json:"status"`                        // Status da conta
	FormaPagamento          *string    `json:"formaPagamento,omitempty"`      // Como foi pago
	Observacoes             *string    `json:"observacoes,omitempty"`         // Observações gerais
	NumeroDocumento         *string    `json:"numeroDocumento,omitempty"`     // Número do documento/nota fiscal
	CreatedAt               time.Time  `json:"createdAt"`
	UpdatedAt               time.Time  `json:"updatedAt"`
}

// ValorSaldo retorna o saldo a receber
func (cr *ContaReceber) ValorSaldo() float64 {
	return cr.ValorOriginal - cr.ValorRecebido
}

// PercentualRecebido calcula o percentual já recebido
func (cr *ContaReceber) PercentualRecebido() float64 {
	if cr.ValorOriginal == 0 {
		return 0
	}
	return (cr.ValorRecebido / cr.ValorOriginal) * 100
}

// EstaVencido verifica se a conta está vencida
func (cr *ContaReceber) EstaVencido() bool {
	return time.Now().After(cr.DataVencimento) && cr.Status != StatusContaReceberRecebido
}

// DiasVencimento retorna quantos dias está vencida (negativo se ainda não venceu)
func (cr *ContaReceber) DiasVencimento() int {
	diff := time.Now().Sub(cr.DataVencimento)
	return int(diff.Hours() / 24)
}

// RegistrarRecebimento registra um recebimento (total ou parcial)
func (cr *ContaReceber) RegistrarRecebimento(valor float64, formaPagamento, observacoes *string) error {
	if valor <= 0 {
		return errors.New("valor deve ser positivo")
	}

	novoValorRecebido := cr.ValorRecebido + valor
	if novoValorRecebido > cr.ValorOriginal {
		return errors.New("valor recebido não pode exceder o valor original")
	}

	cr.ValorRecebido = novoValorRecebido
	now := time.Now()
	cr.DataRecebimento = &now
	cr.UpdatedAt = now

	if formaPagamento != nil {
		cr.FormaPagamento = formaPagamento
	}
	
	if observacoes != nil {
		cr.Observacoes = observacoes
	}

	// Atualiza status baseado no valor recebido
	if cr.ValorRecebido >= cr.ValorOriginal {
		cr.Status = StatusContaReceberRecebido
	} else if cr.ValorRecebido > 0 {
		cr.Status = StatusContaReceberParcial
	}

	return nil
}

// MarcarComoVencido marca a conta como vencida
func (cr *ContaReceber) MarcarComoVencido() {
	if cr.Status == StatusContaReceberPendente && cr.EstaVencido() {
		cr.Status = StatusContaReceberVencido
		cr.UpdatedAt = time.Now()
	}
}

// Cancelar cancela a conta a receber
func (cr *ContaReceber) Cancelar(motivo *string) error {
	if cr.Status == StatusContaReceberRecebido {
		return errors.New("não é possível cancelar uma conta já recebida")
	}
	
	cr.Status = StatusContaReceberCancelado
	cr.UpdatedAt = time.Now()
	
	if motivo != nil {
		cr.Observacoes = motivo
	}
	
	return nil
}

// Validar valida os dados da conta a receber
func (cr *ContaReceber) Validar() error {
	if cr.Cliente == "" {
		return errors.New("cliente é obrigatório")
	}
	if cr.TipoContaReceber == "" {
		return errors.New("tipoContaReceber é obrigatório")
	}
	if cr.TipoContaReceber != TipoContaReceberObra && 
		cr.TipoContaReceber != TipoContaReceberServico && 
		cr.TipoContaReceber != TipoContaReceberOutros {
		return errors.New("tipoContaReceber deve ser OBRA, SERVICO ou OUTROS")
	}
	if cr.Descricao == "" {
		return errors.New("descrição é obrigatória")
	}
	if cr.ValorOriginal <= 0 {
		return errors.New("valorOriginal deve ser positivo")
	}
	if cr.DataVencimento.IsZero() {
		return errors.New("dataVencimento é obrigatória")
	}
	return nil
}