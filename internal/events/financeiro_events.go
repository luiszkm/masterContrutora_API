package events

import "time"

// Eventos relacionados ao módulo Financeiro

// Nomes dos eventos
const (
	// Evento disparado quando uma conta a receber é criada
	ContaReceberCriada = "financeiro:conta_receber_criada"
	
	// Evento disparado quando uma conta a receber é paga (total ou parcial)
	ContaReceberPaga = "financeiro:conta_receber_paga"
	
	// Evento disparado quando uma conta a receber vence
	ContaReceberVencida = "financeiro:conta_receber_vencida"
	
	// Evento disparado quando uma movimentação financeira é registrada
	MovimentacaoFinanceiraRegistrada = "financeiro:movimentacao_registrada"
)

// ContaReceberCriadaPayload contém dados da conta a receber criada
type ContaReceberCriadaPayload struct {
	ContaReceberID          string     `json:"contaReceberId"`
	ObraID                  *string    `json:"obraId,omitempty"`
	CronogramaRecebimentoID *string    `json:"cronogramaRecebimentoId,omitempty"`
	Cliente                 string     `json:"cliente"`
	TipoContaReceber        string     `json:"tipoContaReceber"`
	Descricao               string     `json:"descricao"`
	ValorOriginal           float64    `json:"valorOriginal"`
	DataVencimento          time.Time  `json:"dataVencimento"`
	NumeroDocumento         *string    `json:"numeroDocumento,omitempty"`
	UsuarioID               string     `json:"usuarioId"`
}

// ContaReceberPagaPayload contém dados do pagamento da conta
type ContaReceberPagaPayload struct {
	ContaReceberID  string     `json:"contaReceberId"`
	ObraID          *string    `json:"obraId,omitempty"`
	Cliente         string     `json:"cliente"`
	ValorRecebido   float64    `json:"valorRecebido"`      // Valor desta operação
	ValorTotalRecebido float64 `json:"valorTotalRecebido"` // Valor total já recebido
	ValorOriginal   float64    `json:"valorOriginal"`
	ValorSaldo      float64    `json:"valorSaldo"`         // Saldo restante
	DataRecebimento time.Time  `json:"dataRecebimento"`
	FormaPagamento  *string    `json:"formaPagamento,omitempty"`
	Status          string     `json:"status"`             // PARCIAL ou RECEBIDO
	ContaBancariaID *string    `json:"contaBancariaId,omitempty"`
	UsuarioID       string     `json:"usuarioId"`
}

// ContaReceberVencidaPayload contém dados da conta vencida
type ContaReceberVencidaPayload struct {
	ContaReceberID   string    `json:"contaReceberId"`
	ObraID           *string   `json:"obraId,omitempty"`
	Cliente          string    `json:"cliente"`
	Descricao        string    `json:"descricao"`
	ValorOriginal    float64   `json:"valorOriginal"`
	ValorSaldo       float64   `json:"valorSaldo"`
	DataVencimento   time.Time `json:"dataVencimento"`
	DiasVencidos     int       `json:"diasVencidos"`
	TipoContaReceber string    `json:"tipoContaReceber"`
}

// MovimentacaoFinanceiraRegistradaPayload contém dados da movimentação
type MovimentacaoFinanceiraRegistradaPayload struct {
	MovimentacaoID       string    `json:"movimentacaoId"`
	ContaBancariaID      string    `json:"contaBancariaId"`
	CategoriaID          *string   `json:"categoriaId,omitempty"`
	TipoMovimentacao     string    `json:"tipoMovimentacao"` // ENTRADA ou SAIDA
	Valor                float64   `json:"valor"`
	DataMovimentacao     time.Time `json:"dataMovimentacao"`
	DataCompetencia      time.Time `json:"dataCompetencia"`
	Descricao            string    `json:"descricao"`
	DocumentoID          *string   `json:"documentoId,omitempty"` // ID do documento origem
	DocumentoTipo        *string   `json:"documentoTipo,omitempty"` // CONTA_RECEBER, ORCAMENTO, APONTAMENTO
	Status               string    `json:"status"`          // PREVISTO, REALIZADO, CONCILIADO
	UsuarioID            string    `json:"usuarioId"`
}