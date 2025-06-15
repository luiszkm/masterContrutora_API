// file: internal/domain/obras/obra.go
package obras

import (
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
	ID         string     `json:"id"`
	Nome       string     `json:"nome"`
	Cliente    string     `json:"cliente"`
	Endereco   string     `json:"endereco"`
	DataInicio time.Time  `json:"dataInicio"`
	DataFim    time.Time  `json:"dataFim,omitempty"`
	Status     Status     `json:"status"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"` // Marca de exclusão lógica
}
