// file: internal/service/pessoal/dto/apontamento_input.go
package dto

type CriarApontamentoInput struct {
	FuncionarioID string
	ObraID        string
	PeriodoInicio string // Formato "YYYY-MM-DD"
	PeriodoFim    string // Formato "YYYY-MM-DD"
}
