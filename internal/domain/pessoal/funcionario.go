// file: internal/domain/pessoal/funcionario.go
package pessoal

import (
	"time"
)

type Funcionario struct {
	ID                  string     `json:"id"`
	Nome                string     `json:"nome"`
	CPF                 string     `json:"cpf"`
	Telefone            string     `json:"telefone"`
	Cargo               string     `json:"cargo"`
	Email               string     `json:"email,omitempty"` // Email do funcionário, opcional
	Departamento        string     `json:"departamento"`
	DataContratacao     time.Time  `json:"dataContratacao"`
	ValorDiaria         float64    `json:"valorDiaria"`                  // Valor da diária do funcionário, usado para calcular o custo diário
	ChavePix            string     `json:"chavePix"`                     // Chave PIX do funcionário para pagamentos
	Status              string     `json:"status"`                       // Ativo, Inativo, Desligado
	DesligamentoData    *time.Time `json:"desligamentoData,omitempty"`   // Data de desligamento, se aplicável
	MotivoDesligamento  string     `json:"motivoDesligamento,omitempty"` // Motivo do desligamento, se aplicável
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	Diaria              float64    `json:"diaria"`                        // Valor da diária do funcionário, usado para calcular o custo diário
	AvaliacaoDesempenho string     `json:"avaliacaoDesempenho,omitempty"` // Avaliação de desempenho do funcionário
	Observacoes         string     `json:"observacoes,omitempty"`         // Observações adicionais sobre o funcionário
}
