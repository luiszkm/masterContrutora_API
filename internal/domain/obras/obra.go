// file: internal/domain/obras/obra.go
package obras

import (
	"errors"
	"time"
)

// Status representa os possíveis estados de uma obra.
type Status string

const (
	StatusEmPlanejamento Status = "Em Planejamento"
	StatusEmAndamento    Status = "Em Andamento"
	StatusConcluida      Status = "Concluída"
	StatusCancelada      Status = "Cancelada"
)

// Obra representa o agregado principal do contexto de Obras, conforme ADR-003.
type Obra struct {
	ID         string     `json:"id" db:"id"`
	Nome       string     `json:"nome" db:"nome"`
	Cliente    string     `json:"cliente" db:"cliente"`
	Endereco   string     `json:"endereco" db:"endereco"`
	DataInicio time.Time  `json:"dataInicio" db:"data_inicio"`
	DataFim    *time.Time `json:"dataFim,omitempty" db:"data_fim"`
	Status     Status     `json:"status" db:"status"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"` // Marca de exclusão lógica
	Descricao  *string    `json:"descricao,omitempty" db:"descricao"` // Descrição opcional da obra
	
	// Campos Financeiros
	ValorContratoTotal    float64 `json:"valorContratoTotal" db:"valor_contrato_total"`    // Valor total do contrato
	ValorRecebido         float64 `json:"valorRecebido" db:"valor_recebido"`         // Valor já recebido
	TipoCobranca          string  `json:"tipoCobranca" db:"tipo_cobranca"`          // "VISTA", "PARCELADO", "ETAPAS"
	DataAssinaturaContrato *time.Time `json:"dataAssinaturaContrato,omitempty" db:"data_assinatura_contrato"` // Data da assinatura do contrato
}

// TipoCobranca representa os tipos de cobrança possíveis
const (
	TipoCobrancaVista     = "VISTA"
	TipoCobrancaParcelado = "PARCELADO"
	TipoCobrancaEtapas    = "ETAPAS"
)

// ValorSaldo calcula o saldo a receber da obra
func (o *Obra) ValorSaldo() float64 {
	return o.ValorContratoTotal - o.ValorRecebido
}

// PercentualRecebido calcula o percentual já recebido
func (o *Obra) PercentualRecebido() float64 {
	if o.ValorContratoTotal == 0 {
		return 0
	}
	return (o.ValorRecebido / o.ValorContratoTotal) * 100
}

// RegistrarRecebimento atualiza o valor recebido
func (o *Obra) RegistrarRecebimento(valor float64) error {
	if valor <= 0 {
		return errors.New("valor deve ser positivo")
	}
	
	novoValorRecebido := o.ValorRecebido + valor
	if novoValorRecebido > o.ValorContratoTotal {
		return errors.New("valor recebido não pode exceder o valor do contrato")
	}
	
	o.ValorRecebido = novoValorRecebido
	return nil
}
