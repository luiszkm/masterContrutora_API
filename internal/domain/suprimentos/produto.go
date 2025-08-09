// file: internal/domain/suprimentos/material.go
package suprimentos

import (
	"time"
)

// Material representa um item de suprimento que pode ser usado em orçamentos.
type Produto struct {
	ID              string     `json:"id" db:"id"`
	Nome            string     `json:"nome" db:"nome"`
	Descricao       *string    `json:"descricao,omitempty" db:"descricao"` // MUDANÇA: Agora é um ponteiro
	UnidadeDeMedida string     `json:"unidadeDeMedida" db:"unidade_de_medida"`
	Categoria       string     `json:"categoria" db:"categoria"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

// SoftDelete marca o produto como deletado
func (p *Produto) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.UpdatedAt = now
}

// IsDeleted verifica se o produto foi soft deleted
func (p *Produto) IsDeleted() bool {
	return p.DeletedAt != nil
}
