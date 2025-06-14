// file: internal/service/suprimentos/dto/input.go
package dto

// CadastrarFornecedorInput é o DTO para o caso de uso de cadastro.
type CadastrarFornecedorInput struct {
	Nome      string
	CNPJ      string
	Categoria string
	Contato   string
	Email     string
}

// CadastrarMaterialInput é o DTO para o caso de uso de cadastro de material.
type CadastrarMaterialInput struct {
	Nome            string
	Descricao       string
	UnidadeDeMedida string
	Categoria       string
}
