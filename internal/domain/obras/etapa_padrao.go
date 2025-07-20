// file: internal/domain/obras/etapa_padrao.go
package obras

import "time"

// EtapaPadrao representa um item no catálogo de etapas.
type EtapaPadrao struct {
	ID        string    `json:"id"`
	Nome      string    `json:"nome"`
	Descricao *string   `json:"descricao,omitempty"`
	Ordem     int       `json:"ordem"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
