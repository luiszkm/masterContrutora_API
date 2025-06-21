// file: internal/domain/pessoal/apontamento.go
package pessoal

import (
	"errors"
	"time"
)

// Constantes para os status do ciclo de vida, para evitar "magic strings".
const (
	StatusApontamentoEmAberto              = "EM_ABERTO"
	StatusApontamentoAprovadoParaPagamento = "APROVADO_PARA_PAGAMENTO"
	StatusApontamentoPago                  = "PAGO"
)

// ApontamentoQuinzenal representa os dados transacionais de uma quinzena de trabalho.
type ApontamentoQuinzenal struct {
	ID                  string
	FuncionarioID       string
	ObraID              string
	PeriodoInicio       time.Time
	PeriodoFim          time.Time
	Diaria              float64 // Valor da diária
	DiasTrabalhados     int
	Adicionais          float64
	Descontos           float64
	Adiantamentos       float64
	ValorTotalCalculado float64
	Status              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// --- MÉTODOS DE NEGÓCIO (Rich Domain Model) ---

// Aprovar valida e executa a transição de estado de EM_ABERTO para APROVADO_PARA_PAGAMENTO.
func (a *ApontamentoQuinzenal) Aprovar() error {
	if a.Status != StatusApontamentoEmAberto {
		return errors.New("só é possível aprovar um apontamento que está 'Em Aberto'")
	}
	a.Status = StatusApontamentoAprovadoParaPagamento
	a.UpdatedAt = time.Now()
	return nil
}

// RegistrarPagamento valida e finaliza o ciclo de vida, movendo para PAGO.
func (a *ApontamentoQuinzenal) RegistrarPagamento() error {
	if a.Status != StatusApontamentoAprovadoParaPagamento {
		return errors.New("só é possível pagar um apontamento que está 'Aprovado para Pagamento'")
	}
	a.Status = StatusApontamentoPago
	a.UpdatedAt = time.Now()
	return nil
}

// EditarDiasTrabalhados protege o invariante de que um apontamento pago não pode ser alterado.
func (a *ApontamentoQuinzenal) EditarDiasTrabalhados(novosDias int) error {
	if a.Status == StatusApontamentoPago {
		return errors.New("não é possível editar um apontamento que já foi pago")
	}
	a.DiasTrabalhados = novosDias
	a.UpdatedAt = time.Now()
	// TODO: Adicionar uma chamada para recalcular o ValorTotalCalculado aqui.
	return nil
}
