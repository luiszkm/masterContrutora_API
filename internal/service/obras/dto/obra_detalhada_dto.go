// file: internal/service/obras/dto/obra_detalhada_dto.go
package dto

import "time"

// ObraDetalhadaDTO é a estrutura completa de resposta para o endpoint de detalhes da obra.
type ObraDetalhadaDTO struct {
	ID           string                  `json:"id"`
	Nome         string                  `json:"nome"`
	Cliente      string                  `json:"cliente"`
	Endereco     string                  `json:"endereco"`
	DataInicio   time.Time               `json:"dataInicio"`
	DataFim      *time.Time              `json:"dataFim,omitempty"`
	Descricao    *string                 `json:"descricao,omitempty"` // Adiciona descrição opcional
	Status       string                  `json:"status"`
	Etapas       []EtapaDTO              `json:"etapas"`
	Funcionarios []FuncionarioAlocadoDTO `json:"funcionarios"`
	Fornecedores []FornecedorDTO         `json:"fornecedores"`
	Orcamentos   []OrcamentoDTO          `json:"orcamentos"`
	Produtos     []ProdutoDto            `json:"produtos"`
}

// EtapaDTO representa uma etapa dentro da resposta detalhada.
type EtapaDTO struct {
	ID                 string    `json:"id"`
	Nome               string    `json:"nome"`
	DataInicioPrevista time.Time `json:"dataInicioPrevista"`
	DataFimPrevista    time.Time `json:"dataFimPrevista"`
	Status             string    `json:"status"`
}

// FuncionarioAlocadoDTO representa um funcionário alocado na obra.
type FuncionarioAlocadoDTO struct {
	FuncionarioID      string    `json:"funcionarioId"`
	NomeFuncionario    string    `json:"nomeFuncionario"`
	DataInicioAlocacao time.Time `json:"dataInicioAlocacao"`
}

type FornecedorDTO struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}

type OrcamentoDTO struct {
	ID         string  `json:"id"`
	Numero     string  `json:"numero"`
	ValorTotal float64 `json:"valorTotal"` // ADICIONADO: Expõe o valor total
	Status     string  `json:"status"`     // ADICIONADO: Expõe o status
}

type ProdutoDto struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
}
