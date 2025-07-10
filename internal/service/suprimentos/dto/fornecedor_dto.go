package dto

type AtualizarFornecedorInput struct {
	Nome         *string
	CNPJ         *string
	CategoriaIDs []string
	Contato      *string
	Email        *string
	Status       *string
	Endereco     *string
	Avaliacao    *float64
	Observacoes  *string
}
