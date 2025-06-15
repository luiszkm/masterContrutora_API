package dtos

type CadastrarMaterialRequest struct {
	Nome            string `json:"nome"`
	Descricao       string `json:"descricao"`
	UnidadeDeMedida string `json:"unidadeDeMedida"`
	Categoria       string `json:"categoria"`
}

type MaterialResponse struct {
	ID              string `json:"id"`
	Nome            string `json:"nome"`
	Descricao       string `json:"descricao"`
	UnidadeDeMedida string `json:"unidadeDeMedida"`
	Categoria       string `json:"categoria"`
}

type CriarOrcamentoRequest struct {
	Numero       string        `json:"numero"`
	FornecedorID string        `json:"fornecedorId"`
	Itens        []ItemRequest `json:"itens"`
}
type ItemRequest struct {
	MaterialID    string  `json:"materialId"`
	Quantidade    float64 `json:"quantidade"`
	ValorUnitario float64 `json:"valorUnitario"`
}

type AtualizarStatusRequest struct {
	Status string `json:"status"`
}

type AtualizarFornecedorRequest struct {
	Nome      string `json:"nome"`
	CNPJ      string `json:"cnpj"`
	Categoria string `json:"categoria"`
	Contato   string `json:"contato"`
	Email     string `json:"email"`
}
