package dto

import "time"

type ListagemFuncionarioDTO struct {
	ID                string    `json:"id"`
	Nome              string    `json:"nome"`
	Cargo             string    `json:"cargo"`
	Departamento      *string   `json:"departamento"`
	DataContratacao   time.Time `json:"dataContratacao"`
	Diaria            float64   `json:"valorDiaria"`
	DiasTrabalhados   *int      `json:"diasTrabalhados"`
	ValorAdicional    *float64  `json:"valorAdicional"`
	Descontos         *float64  `json:"descontos"`
	Adiantamento      *float64  `json:"adiantamento"`
	ChavePix          *string   `json:"chavePix"`
	Avaliacao         *string   `json:"avaliacao"` // Nota: Campo novo, não populado ainda.
	StatusApontamento *string   `json:"statusApontamento"`
	ApontamentoId     *string   `json:"apontamentoId"` // ID do apontamento quinzenal, se aplicável
}
