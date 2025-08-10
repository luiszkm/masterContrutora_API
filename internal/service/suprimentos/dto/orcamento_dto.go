// file: internal/service/suprimentos/dto/orcamento_dto.go
package dto

import "time"

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

type OrcamentoListItemDTO struct {
	ID             string    `json:"id" db:"id"`
	Numero         string    `json:"numero" db:"numero"`
	ValorTotal     float64   `json:"valorTotal" db:"valor_total"`
	Status         string    `json:"status" db:"status"`
	DataEmissao    time.Time `json:"dataEmissao" db:"data_emissao"`
	ObraID         string    `json:"obraId" db:"obra_id"`
	ObraNome       string    `json:"obraNome" db:"obra_nome"`
	FornecedorID   string    `json:"fornecedorId" db:"fornecedor_id"`
	FornecedorNome string    `json:"fornecedorNome" db:"fornecedor_nome"` // CAMPO ADICIONADO
	ItensCount     int       `json:"itensCount" db:"itens_count"`         // NOVO CAMPO
	Categorias     []string  `json:"categorias" db:"categorias"`          // ARRAY DE CATEGORIAS
}

type OrcamentoDetalhadoDTO struct {
	ID                 string                      `json:"id" db:"id"`
	Numero             string                      `json:"numero" db:"numero"`
	ValorTotal         float64                     `json:"valorTotal" db:"valor_total"`
	Status             string                      `json:"status" db:"status"`
	DataEmissao        time.Time                   `json:"dataEmissao" db:"data_emissao"`
	Observacoes        *string                     `json:"observacoes,omitempty" db:"observacoes"`
	CondicoesPagamento *string                     `json:"condicoesPagamento,omitempty" db:"condicoes_pagamento"`
	Obra               ObraInfoDTO                 `json:"obra" db:"obra"`
	Etapa              EtapaInfoDTO                `json:"etapa" db:"etapa"`
	Fornecedor         FornecedorInfoDTO           `json:"fornecedor" db:"fornecedor"`
	Itens              []ItemOrcamentoDetalhadoDTO `json:"itens" db:"itens"`
}

// Sub-DTOs para manter a resposta organizada
type ObraInfoDTO struct {
	ID   string `json:"id" db:"id"`
	Nome string `json:"nome" db:"nome"`
}
type EtapaInfoDTO struct {
	ID   string `json:"id" db:"id"`
	Nome string `json:"nome" db:"nome"`
}
type FornecedorInfoDTO struct {
	ID   string `json:"id" db:"id"`
	Nome string `json:"nome" db:"nome"`
}

type ItemOrcamentoDetalhadoDTO struct {
	NomeProduto     string `json:"ProdutoNome" db:"produto_nome"`
	UnidadeDeMedida string `json:"UnidadeDeMedida" db:"unidade_de_medida"`
	Categoria       string `json:"Categoria" db:"categoria"`
	Quantidade      float64
	ValorUnitario   float64
}

type AtualizarOrcamentoInput struct {
	FornecedorID       string
	EtapaID            string
	Observacoes        *string
	CondicoesPagamento *string
	Itens              []ItemOrcamentoInput
}

// CompararOrcamentosRequest é o DTO de entrada para comparar orçamentos por categoria
type CompararOrcamentosRequest struct {
	Categoria string `json:"categoria"` // Nome da categoria (ex: "Cimento", "Aço")
}

// CompararOrcamentosResponse é o DTO de resposta para comparação de orçamentos
type CompararOrcamentosResponse struct {
	Categoria  string                `json:"categoria"`
	Orcamentos []OrcamentoComparacao `json:"orcamentos"`
}

// OrcamentoComparacao representa um orçamento na comparação
type OrcamentoComparacao struct {
	ID             string    `json:"id"`
	Numero         string    `json:"numero"`
	FornecedorNome string    `json:"fornecedorNome"`
	ValorTotal     float64   `json:"valorTotal"`
	Status         string    `json:"status"`
	DataEmissao    time.Time `json:"dataEmissao"`
	ItensCategoria int       `json:"itensCategoria"` // Quantidade de itens da categoria específica
}
