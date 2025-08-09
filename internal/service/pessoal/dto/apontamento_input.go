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
	FuncionarioID   string  `json:"funcionarioId"`   // ID do funcionário (não alterável)
	ObraID          string  `json:"obraId"`          // ID da obra (opcional - mantém atual se vazio)
	PeriodoInicio   string  `json:"periodoInicio"`   // Formato "YYYY-MM-DD" (obrigatório)
	PeriodoFim      string  `json:"periodoFim"`      // Formato "YYYY-MM-DD" (obrigatório)
	Diaria          float64 `json:"diaria"`          // Valor da diária (obrigatório)
	DiasTrabalhados int     `json:"diasTrabalhados"` // Número de dias trabalhados (obrigatório)
	ValorAdicional  float64 `json:"valorAdicional"`  // Valor adicional, se houver (opcional)
	Descontos       float64 `json:"descontos"`       // Descontos aplicáveis, se houver (opcional)
	Adiantamento    float64 `json:"adiantamento"`    // Valor do adiantamento, se houver (opcional)
	Status          string  `json:"status"`          // Status do apontamento (não alterável via PUT)
}
