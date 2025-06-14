// file: internal/service/obras/dto/etapa_dto.go
package dto

// AdicionarEtapaInput é o DTO para adicionar uma nova etapa a uma obra.
type AdicionarEtapaInput struct {
	Nome               string `json:"nome"`
	DataInicioPrevista string `json:"dataInicioPrevista"` // Formato "YYYY-MM-DD"
	DataFimPrevista    string `json:"dataFimPrevista"`    // Formato "YYYY-MM-DD"
}

type AtualizarStatusEtapaInput struct {
	Status string `json:"status"`
}
