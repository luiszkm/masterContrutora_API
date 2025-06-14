// file: internal/service/obras/dto/dashboard_dto.go
package dto

import "time"

// ObraDashboard é o nosso Modelo de Leitura (ViewModel) para o painel.
// Conforme definido na documentação arquitetural.
type ObraDashboard struct {
	ObraID                 string     `json:"obraId"`
	NomeObra               string     `json:"nomeObra"`
	StatusObra             string     `json:"statusObra"`
	EtapaAtualNome         *string    `json:"etapaAtualNome"`       // Ponteiro para aceitar NULL
	DataFimPrevistaEtapa   *time.Time `json:"dataFimPrevistaEtapa"` // Ponteiro para aceitar NULL
	DiasParaPrazoEtapa     *int       `json:"diasParaPrazoEtapa"`   // Ponteiro para aceitar NULL
	PercentualConcluido    float64    `json:"percentualConcluido"`
	CustoTotalRealizado    float64    `json:"custoTotalRealizado"`
	OrcamentoTotalAprovado float64    `json:"orcamentoTotalAprovado"`
	BalancoFinanceiro      float64    `json:"balancoFinanceiro"`
	FuncionariosAlocados   int        `json:"funcionariosAlocados"`
	UltimaAtualizacao      time.Time  `json:"ultimaAtualizacao"`
}
