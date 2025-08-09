package dto

// ObraListItemDTO representa os dados de uma única obra em uma lista.
type ObraListItemDTO struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	Cliente  string `json:"cliente"`
	Status   string `json:"status"`
	Etapa    string `json:"etapa"`
	Evolucao string `json:"evolucao"` // Evolução da obra, como percentual de conclusão
}
