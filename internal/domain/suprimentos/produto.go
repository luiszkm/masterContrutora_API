// file: internal/domain/suprimentos/material.go
package suprimentos

// Material representa um item de suprimento que pode ser usado em orçamentos.
type Produto struct {
	ID              string  `json:"id" db:"id"`
	Nome            string  `json:"nome" db:"nome"`
	Descricao       *string `json:"descricao,omitempty" db:"descricao"` // MUDANÇA: Agora é um ponteiro
	UnidadeDeMedida string  `json:"unidadeDeMedida" db:"unidade_de_medida"`
	Categoria       string  `json:"categoria" db:"categoria"`
}
