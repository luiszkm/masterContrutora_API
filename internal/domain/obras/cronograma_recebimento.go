package obras

import (
	"errors"
	"time"
)

// StatusRecebimento representa os possíveis status de um cronograma de recebimento
const (
	StatusRecebimentoPendente = "PENDENTE"
	StatusRecebimentoRecebido = "RECEBIDO"
	StatusRecebimentoVencido  = "VENCIDO"
	StatusRecebimentoParcial  = "PARCIAL"
)

// CronogramaRecebimento representa uma etapa de recebimento de uma obra
type CronogramaRecebimento struct {
	ID                string     `json:"id"`
	ObraID            string     `json:"obraId"`
	NumeroEtapa       int        `json:"numeroEtapa"`
	DescricaoEtapa    string     `json:"descricaoEtapa"`
	ValorPrevisto     float64    `json:"valorPrevisto"`
	DataVencimento    time.Time  `json:"dataVencimento"`
	Status            string     `json:"status"`
	DataRecebimento   *time.Time `json:"dataRecebimento,omitempty"`
	ValorRecebido     float64    `json:"valorRecebido"`
	ObservacoesRecebimento *string `json:"observacoesRecebimento,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

// ValorSaldo retorna o valor que ainda falta receber nesta etapa
func (cr *CronogramaRecebimento) ValorSaldo() float64 {
	return cr.ValorPrevisto - cr.ValorRecebido
}

// PercentualRecebido calcula o percentual já recebido desta etapa
func (cr *CronogramaRecebimento) PercentualRecebido() float64 {
	if cr.ValorPrevisto == 0 {
		return 0
	}
	return (cr.ValorRecebido / cr.ValorPrevisto) * 100
}

// EstaVencido verifica se o cronograma está vencido
func (cr *CronogramaRecebimento) EstaVencido() bool {
	return time.Now().After(cr.DataVencimento) && cr.Status != StatusRecebimentoRecebido
}

// RegistrarRecebimento registra um recebimento (total ou parcial)
func (cr *CronogramaRecebimento) RegistrarRecebimento(valor float64, observacoes *string) error {
	if valor <= 0 {
		return errors.New("valor deve ser positivo")
	}

	novoValorRecebido := cr.ValorRecebido + valor
	if novoValorRecebido > cr.ValorPrevisto {
		return errors.New("valor recebido não pode exceder o valor previsto")
	}

	cr.ValorRecebido = novoValorRecebido
	now := time.Now()
	cr.DataRecebimento = &now
	cr.UpdatedAt = now
	
	if observacoes != nil {
		cr.ObservacoesRecebimento = observacoes
	}

	// Atualiza status baseado no valor recebido
	if cr.ValorRecebido >= cr.ValorPrevisto {
		cr.Status = StatusRecebimentoRecebido
	} else if cr.ValorRecebido > 0 {
		cr.Status = StatusRecebimentoParcial
	}

	return nil
}

// MarcarComoVencido marca o cronograma como vencido
func (cr *CronogramaRecebimento) MarcarComoVencido() {
	if cr.Status == StatusRecebimentoPendente && cr.EstaVencido() {
		cr.Status = StatusRecebimentoVencido
		cr.UpdatedAt = time.Now()
	}
}

// Validar valida os dados do cronograma
func (cr *CronogramaRecebimento) Validar() error {
	if cr.ObraID == "" {
		return errors.New("obraId é obrigatório")
	}
	if cr.NumeroEtapa <= 0 {
		return errors.New("numeroEtapa deve ser positivo")
	}
	if cr.DescricaoEtapa == "" {
		return errors.New("descricaoEtapa é obrigatória")
	}
	if cr.ValorPrevisto <= 0 {
		return errors.New("valorPrevisto deve ser positivo")
	}
	if cr.DataVencimento.IsZero() {
		return errors.New("dataVencimento é obrigatória")
	}
	return nil
}