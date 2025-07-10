// file: internal/domain/suprimentos/fornecedor.go
package suprimentos

type Fornecedor struct {
	ID              string
	Nome            string
	CNPJ            string
	Categorias      []Categoria
	Contato         string
	Email           string
	Status          string
	Website         *string
	Endereco        *string
	NomeAtendente   *string
	Avaliacao       *float64
	Observacoes     *string
	OrcamentosCount int `json:"orcamentosCount"`
}
