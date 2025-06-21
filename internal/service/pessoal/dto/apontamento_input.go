// file: internal/service/pessoal/dto/apontamento_input.go
package dto

type CriarApontamentoInput struct {
	FuncionarioID   string
	ObraID          string
	PeriodoInicio   string  // Formato "YYYY-MM-DD"
	PeriodoFim      string  // Formato "YYYY-MM-DD"
	Diaria          float64 // Valor da diária
	DiasTrabalhados int     // Número de dias trabalhados
	ValorAdicional  float64 // Valor adicional, se houver
	Descontos       float64 // Descontos aplicáveis, se houver
	Adiantamento    float64 // Valor do adiantamento, se houver
}
type AtualizarApontamentoInput struct {
	FuncionarioID   string  // ID do funcionário
	ObraID          string  // ID da obra
	PeriodoInicio   string  // Formato "YYYY-MM-DD"
	PeriodoFim      string  // Formato "YYYY-MM-DD"
	Diaria          float64 // Valor da diária
	DiasTrabalhados int     // Número de dias trabalhados
	ValorAdicional  float64 // Valor adicional, se houver
	Descontos       float64 // Descontos aplicáveis, se houver
	Adiantamento    float64 // Valor do adiantamento, se houver
	Status          string  // Status do apontamento, ex: "Pendente", "Aprovado", "Rejeitado"
}
