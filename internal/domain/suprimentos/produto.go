// file: internal/domain/suprimentos/material.go
package suprimentos

// Material representa um item de suprimento que pode ser usado em orçamentos.
type Produto struct {
	ID              string
	Nome            string
	Descricao       string
	UnidadeDeMedida string
	Categoria       string
}
