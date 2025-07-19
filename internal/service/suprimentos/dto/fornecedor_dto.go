package dto

type AtualizarFornecedorInput struct {
	Nome         *string   `json:"nome"`
	CNPJ         *string   `json:"cnpj"`
	Contato      *string   `json:"contato"`
	Email        *string   `json:"email"`
	Status       *string   `json:"status"`
	Endereco     *string   `json:"endereco"`
	Avaliacao    *float64  `json:"avaliacao"`
	Observacoes  *string   `json:"observacoes"`
	CategoriaIDs *[]string `json:"categoriaIds"` // aqui também precisa ser ponteiro
}
