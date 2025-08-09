// file: internal/domain/suprimentos/fornecedor.go
package suprimentos

type Fornecedor struct {
	ID              string      `json:"id"`
	Nome            string      `json:"nome"`
	CNPJ            string      `json:"cnpj"`
	Contato         *string     `json:"contato,omitempty"` // MUDANÇA: Agora é um ponteiro
	Email           *string     `json:"email,omitempty"`   // MUDANÇA: Agora é um ponteiro
	Status          string      `json:"status"`
	Endereco        *string     `json:"endereco,omitempty"`
	Avaliacao       *float64    `json:"avaliacao,omitempty"`
	Observacoes     *string     `json:"observacoes,omitempty"`
	Categorias      []Categoria `json:"categorias"`
	OrcamentosCount int         `json:"orcamentosCount"`
}
