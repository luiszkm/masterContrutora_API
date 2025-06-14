// file: internal/domain/suprimentos/material.go
package suprimentos

// Material representa um item de suprimento que pode ser usado em orçamentos.
type Material struct {
	ID              string
	Nome            string
	Descricao       string
	UnidadeDeMedida string // Ex: "m³", "kg", "unidade"
	Categoria       string // Ex: "Estrutural", "Acabamento", "Elétrico"
}
