// file: internal/service/obras/dto/alocacao_dto.go
package dto

// AlocarFuncionarioInput é o DTO para alocar um funcionário a uma obra.
type AlocarFuncionarioInput struct {
	FuncionarioID      string `json:"funcionarioId"`
	DataInicioAlocacao string `json:"dataInicioAlocacao"` // Formato "YYYY-MM-DD"
}
