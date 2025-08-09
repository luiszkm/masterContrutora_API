// file: internal/service/obras/dto/alocacao_dto.go
package dto

// AlocarFuncionariosInput é o DTO para alocar um ou mais funcionários a uma obra.
type AlocarFuncionariosInput struct {
	FuncionarioIDs     []string `json:"funcionarioIds"`
	DataInicioAlocacao string   `json:"dataInicioAlocacao"` // Formato "YYYY-MM-DD"
}
