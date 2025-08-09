package dtos

type CadastrarMaterialRequest struct {
	Nome            string  `json:"nome"`
	Descricao       *string `json:"descricao"`
	UnidadeDeMedida string  `json:"unidadeDeMedida"`
	Categoria       string  `json:"categoria"`
}

type MaterialResponse struct {
	ID              string  `json:"id"`
	Nome            string  `json:"nome"`
	Descricao       *string `json:"descricao"`
	UnidadeDeMedida string  `json:"unidadeDeMedida"`
	Categoria       string  `json:"categoria"`
}

type CriarOrcamentoRequest struct {
	Numero       string        `json:"numero"`
	FornecedorID string        `json:"fornecedorId"`
	Itens        []ItemRequest `json:"itens"`
}
type ItemRequest struct {
	NomeProduto     string  `json:"nomeProduto"`
	UnidadeDeMedida string  `json:"unidadeDeMedida"`
	Categoria       string  `json:"categoria"`
	Quantidade      float64 `json:"quantidade"`
	ValorUnitario   float64 `json:"valorUnitario"`
}
type AtualizarStatusRequest struct {
	Status string `json:"status"`
}

type AtualizarFornecedorRequest struct {
	Nome         *string   `json:"nome,omitempty"`
	CNPJ         *string   `json:"cnpj,omitempty"`
	CategoriaIDs *[]string `json:"categoriaIds,omitempty"`
	Contato      *string   `json:"contato,omitempty"`
	Email        *string   `json:"email,omitempty"`
	Status       *string   `json:"status,omitempty"`
	Endereco     *string   `json:"endereco,omitempty"`
	Avaliacao    *float64  `json:"avaliacao,omitempty"`
	Observacoes  *string   `json:"observacoes,omitempty"`
}

type CadastrarFornecedorRequest struct {
	Nome         string   `json:"nome"`
	CNPJ         string   `json:"cnpj"`
	CategoriaIDs []string `json:"categoriaIds"`
	Contato      string   `json:"contato"`
	Email        string   `json:"email"`
	Endereco     *string  `json:"endereco,omitempty"`
	Observacoes  *string  `json:"observacoes,omitempty"`
	Avaliacao    *float64 `json:"avaliacao,omitempty"`
}

type AtualizarOrcamentoRequest struct {
	FornecedorID       string        `json:"fornecedorId"`
	EtapaID            string        `json:"etapaId"`
	Observacoes        *string       `json:"observacoes"`
	CondicoesPagamento *string       `json:"condicoesPagamento"`
	Itens              []ItemRequest `json:"itens"`
}
