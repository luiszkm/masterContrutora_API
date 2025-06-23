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
	ID                  string    `json:"id"`
	FuncionarioID       string    `json:"funcionarioId"`
	ObraID              string    `json:"obraId"`
	PeriodoInicio       time.Time `json:"periodoInicio"`       // Data de início do período
	PeriodoFim          time.Time `json:"periodoFim"`          // Data de fim do período
	Diaria              float64   `json:"diaria"`              // Valor da diária
	DiasTrabalhados     int       `json:"diasTrabalhados"`     // Número de dias trabalhados no período
	Adicionais          float64   `json:"adicionais"`          // Valor de adicionais (bônus, horas extras, etc.)
	Descontos           float64   `json:"descontos"`           // Valor de descontos (faltas, atrasos, etc.)
	Adiantamentos       float64   `json:"adiantamentos"`       // Valor de adiantamentos já pagos
	ValorTotalCalculado float64   `json:"valorTotalCalculado"` // Valor total calculado do apontamento
	Status              string    `json:"status"`              // Status do apontamento (EM_ABERTO, APROVADO_PARA_PAGAMENTO, PAGO)
	CreatedAt           time.Time `json:"createdAt"`           // Data de criação do apontamento
	UpdatedAt           time.Time `json:"updatedAt"`           // Data da última atualização do apontamento
	FuncionarioNome     string    `json:"funcionarioNome"`     // Nome do funcionário (opcional, para exibição)
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

func (a *ApontamentoQuinzenal) recalcularTotal() {
	valorDias := float64(a.DiasTrabalhados) * a.Descontos
	a.ValorTotalCalculado = valorDias + a.Adicionais - a.Descontos - a.Adiantamentos
}

func (a *ApontamentoQuinzenal) AtualizarValores(diasTrabalhados int,
	adicionais, descontos, adiantamentos, valorDiaria float64,
	periodoInicio, periodoFim time.Time, obraId string,
) error {
	if a.Status != StatusApontamentoEmAberto {
		return errors.New("só é possível editar um apontamento que está 'Em Aberto'")
	}

	a.Diaria = valorDiaria
	a.DiasTrabalhados = diasTrabalhados
	a.Adicionais = adicionais
	a.Descontos = descontos
	a.Adiantamentos = adiantamentos
	a.UpdatedAt = time.Now()
	a.PeriodoInicio = periodoInicio
	a.PeriodoFim = periodoFim
	a.ObraID = obraId

	// Chama o método privado para recalcular o total
	a.recalcularTotal()

	return nil
}

func (a *ApontamentoQuinzenal) AprovarEPagar() error {
	if a.Status != StatusApontamentoEmAberto {
		return errors.New("só é possível usar o pagamento direto em um apontamento que está 'Em Aberto'")
	}
	a.Status = StatusApontamentoPago
	a.UpdatedAt = time.Now()
	return nil
}
