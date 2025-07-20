// file: internal/domain/suprimentos/orcamento.go
package suprimentos

import (
	"time"
)

type Status string

const (
	StatusEmAberto  Status = "Em Aberto"
	StatusAprovado  Status = "Aprovado"
	StatusRejeitado Status = "Rejeitado"
	StatusPago      Status = "Pago"
)

type Orcamento struct {
	ID                 string          `json:"id" db:"id"`
	Numero             string          `json:"numero" db:"numero"`
	EtapaID            string          `json:"etapaId" db:"etapa_id"`
	FornecedorID       string          `json:"fornecedorId" db:"fornecedor_id"`
	Itens              []ItemOrcamento `json:"itens"` // Não tem tag db porque é preenchido separadamente
	ValorTotal         float64         `json:"valorTotal" db:"valor_total"`
	Status             string          `json:"status" db:"status"`
	DataEmissao        time.Time       `json:"dataEmissao" db:"data_emissao"`
	DataAprovacao      *time.Time      `json:"dataAprovacao,omitempty" db:"data_aprovacao"`
	CondicoesPagamento *string         `json:"condicoesPagamento,omitempty" db:"condicoes_pagamento"` // CORRIGIDO
	Observacoes        *string         `json:"observacoes,omitempty" db:"observacoes"`
}

type ItemOrcamento struct {
	ID                 string  `json:"id" db:"id"`
	OrcamentoID        string  `json:"orcamentoId" db:"orcamento_id"`
	ProdutoID          string  `json:"produtoId" db:"produto_id"`
	Quantidade         float64 `json:"quantidade" db:"quantidade"`
	ValorUnitario      float64 `json:"valorUnitario" db:"valor_unitario"`
	Categoria          string  `json:"categoria" db:"categoria"`
	UnidadeDeutoMedida string  `json:"unidadeDeMedida" db:"unidade_de_medida"`
	NomeProd           string  `json:"nomeProduto" db:"produto_nome"`
}
