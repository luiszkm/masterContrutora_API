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
	ID            string
	Numero        string // Um número de identificação amigável, ex: "ORC-2025-001"
	EtapaID       string
	FornecedorID  string
	Itens         []ItemOrcamento
	ValorTotal    float64
	Status        string // Ex: "Em Aberto", "Aprovado", "Rejeitado", "Pago"
	DataEmissao   time.Time
	DataAprovacao time.Time
	Observacoes   string // Observações adicionais sobre o orçamento
}

type ItemOrcamento struct {
	ID            string
	OrcamentoID   string
	ProdutoID     string
	Quantidade    float64
	ValorUnitario float64
}
