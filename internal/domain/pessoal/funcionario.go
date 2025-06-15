// file: internal/domain/pessoal/funcionario.go
package pessoal

import (
	"time"
)

type Funcionario struct {
	ID              string
	Nome            string
	CPF             string
	Cargo           string
	DataContratacao time.Time
	// Salario e Diaria podem ser tipos de dinheiro mais complexos no futuro.
	Salario   float64
	Diaria    float64
	Status    string     // "Ativo", "Inativo"
	DeletedAt *time.Time // Marca de exclusão lógica
}
