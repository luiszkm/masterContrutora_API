// file: internal/service/identidade/dto/input.go
package dto

// RegistrarUsuarioInput é o DTO para o caso de uso de registro no serviço.
type RegistrarUsuarioInput struct {
	Nome  string
	Email string
	Senha string
}

// LoginInput é o DTO para o caso de uso de login no serviço.
type LoginInput struct {
	Email string
	Senha string
}
