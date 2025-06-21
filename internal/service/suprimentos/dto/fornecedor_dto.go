package dto

type AtualizarFornecedorInput struct {
	Nome          string
	CNPJ          string
	Categoria     string
	Contato       string
	Email         string
	Status        string
	Website       *string
	Endereco      *string
	NomeAtendente *string
	Avaliacao     *float64
	Observacoes   *string
}
