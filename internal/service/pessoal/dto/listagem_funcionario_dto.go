package dto

import "time"

type ListagemFuncionarioDTO struct {
	ID                  string    `json:"id"`
	Nome                string    `json:"nome"`
	Cargo               string    `json:"cargo"`
	Departamento        *string   `json:"departamento"`
	DataContratacao     time.Time `json:"dataContratacao"`
	Diaria              float64   `json:"valorDiaria"`
	DiasTrabalhados     *int      `json:"diasTrabalhados"`
	ValorAdicional      *float64  `json:"valorAdicional"`
	Descontos           *float64  `json:"descontos"`
	Adiantamento        *float64  `json:"adiantamento"`
	ChavePix            *string   `json:"chavePix"`
	Avaliacao           *string   `json:"avaliacao"` // Nota: Campo novo, não populado ainda.
	StatusApontamento   *string   `json:"statusApontamento"`
	ApontamentoId       *string   `json:"apontamentoId"` // ID do apontamento quinzenal, se aplicável
	Observacoes         *string   `json:"observacoes"`
	AvaliacaoDesempenho *string   `json:"avaliacaoDesempenho"` // Nota: Campo novo, não populado ainda.
}

type ApontamentoDTO struct {
	// Campos do ApontamentoQuinzenal original
	ID                  string    `json:"id"`
	FuncionarioID       string    `json:"funcionarioId"`
	ObraID              string    `json:"obraId"`
	PeriodoInicio       time.Time `json:"periodoInicio"`
	PeriodoFim          time.Time `json:"periodoFim"`
	Diaria              float64   `json:"diaria"`
	DiasTrabalhados     int       `json:"diasTrabalhados"`
	Adicionais          float64   `json:"adicionais"`
	Descontos           float64   `json:"descontos"`
	Adiantamentos       float64   `json:"adiantamentos"`
	ValorTotalCalculado float64   `json:"valorTotalCalculado"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	NomeFuncionario     string    `json:"nomeFuncionario"`
}
