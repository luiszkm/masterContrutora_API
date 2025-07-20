// file: internal/service/suprimentos/dto/orcamento_dto.go
package dto

// CriarOrcamentoInput é o DTO para o caso de uso de criação de orçamento.
type CriarOrcamentoInput struct {
	FornecedorID string
	Itens        []ItemOrcamentoInput
}
type ItemOrcamentoInput struct {
	NomeProduto     string
	UnidadeDeMedida string
	Categoria       string
	Quantidade      float64
	ValorUnitario   float64
}
type AtualizarStatusOrcamentoInput struct {
	Status string
}
