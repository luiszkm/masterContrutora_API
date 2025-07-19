// file: internal/service/suprimentos/dto/categoria_dto.go
package dto

// DTO para a criação de uma nova categoria.
type CriarCategoriaInput struct {
	Nome string `json:"nome"`
}

// DTO para a atualização de uma categoria existente.
type AtualizarCategoriaInput struct {
	Nome string `json:"nome"`
}
