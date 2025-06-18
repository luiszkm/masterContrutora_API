// file: internal/domain/pessoal/funcionario.go
package pessoal

import (
	"time"
)

type Funcionario struct {
	ID                  string
	Nome                string
	CPF                 string
	Telefone            string
	Cargo               string
	Departamento        string
	DataContratacao     time.Time
	ValorDiaria         float64 // O valor contratual da diária
	ChavePix            string
	Status              string     // "Ativo", "Inativo", "Desligado"
	DesligamentoData    *time.Time // Ponteiro para aceitar data nula
	MotivoDesligamento  string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Diaria              float64 // Valor da diária do funcionário, usado para calcular o custo diário
	AvaliacaoDesempenho string  // Avaliação de desempenho do funcionário
	Observacoes         string  // Observações adicionais sobre o funcionário
}
