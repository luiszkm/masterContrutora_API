// file: internal/service/suprimentos/dto/input.go
package dto

// CadastrarFornecedorInput é o DTO para o caso de uso de cadastro.
type CadastrarFornecedorInput struct {
	Nome         string
	CNPJ         string
	CategoriaIDs []string
	Contato      string
	Email        string
	Endereco     *string  // NOVO
	Observacoes  *string  // NOVO
	Avaliacao    *float64 // NOVO
}

type CadastrarProdutoInput struct {
	Nome            string
	Descricao       string
	UnidadeDeMedida string
	Categoria       string
}
