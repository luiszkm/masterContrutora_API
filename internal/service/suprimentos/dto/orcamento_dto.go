// file: internal/service/suprimentos/dto/orcamento_dto.go
package dto

// CriarOrcamentoInput é o DTO para o caso de uso de criação de orçamento.
type CriarOrcamentoInput struct {
	Numero       string
	FornecedorID string
	Itens        []ItemOrcamentoInput
}

type ItemOrcamentoInput struct {
	MaterialID    string
	Quantidade    float64
	ValorUnitario float64
}
type AtualizarStatusOrcamentoInput struct {
	Status string
}
