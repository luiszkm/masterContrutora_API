package events

import "time"

// Eventos relacionados ao módulo de Obras

// Nomes dos eventos
const (
	// Evento disparado quando uma obra tem seu contrato definido/atualizado
	ObraContratoDefinido = "obra:contrato_definido"
	
	// Evento disparado quando um cronograma de recebimento é criado
	CronogramaRecebimentoCriado = "cronograma:recebimento_criado"
	
	// Evento disparado quando uma etapa de recebimento vence
	EtapaRecebimentoVencida = "cronograma:etapa_vencida"
	
	// Evento disparado quando um recebimento é realizado
	RecebimentoRealizado = "obra:recebimento_realizado"
)

// ObraContratoDefinidoPayload contém os dados do evento de contrato definido
type ObraContratoDefinidoPayload struct {
	ObraID                 string     `json:"obraId"`
	ObraNome               string     `json:"obraNome"`
	Cliente                string     `json:"cliente"`
	ValorContratoTotal     float64    `json:"valorContratoTotal"`
	TipoCobranca           string     `json:"tipoCobranca"`
	DataAssinaturaContrato *time.Time `json:"dataAssinaturaContrato,omitempty"`
	UsuarioID              string     `json:"usuarioId"` // Quem definiu o contrato
}

// CronogramaRecebimentoCriadoPayload contém dados do cronograma criado
type CronogramaRecebimentoCriadoPayload struct {
	ObraID             string                     `json:"obraId"`
	ObraNome           string                     `json:"obraNome"`
	Cliente            string                     `json:"cliente"`
	CronogramasIds     []string                   `json:"cronogramasIds"`
	ValorTotalPrevisto float64                    `json:"valorTotalPrevisto"`
	QuantidadeEtapas   int                        `json:"quantidadeEtapas"`
	PrimeiroVencimento time.Time                  `json:"primeiroVencimento"`
	UsuarioID          string                     `json:"usuarioId"`
}

// EtapaRecebimentoVencidaPayload contém dados da etapa vencida
type EtapaRecebimentoVencidaPayload struct {
	CronogramaRecebimentoID string    `json:"cronogramaRecebimentoId"`
	ObraID                  string    `json:"obraId"`
	ObraNome                string    `json:"obraNome"`
	Cliente                 string    `json:"cliente"`
	NumeroEtapa             int       `json:"numeroEtapa"`
	DescricaoEtapa          string    `json:"descricaoEtapa"`
	ValorPrevisto           float64   `json:"valorPrevisto"`
	ValorSaldo              float64   `json:"valorSaldo"` // Valor ainda não recebido
	DataVencimento          time.Time `json:"dataVencimento"`
	DiasVencidos            int       `json:"diasVencidos"`
}

// RecebimentoRealizadoPayload contém dados do recebimento realizado
type RecebimentoRealizadoPayload struct {
	CronogramaRecebimentoID *string    `json:"cronogramaRecebimentoId,omitempty"` // Pode ser null se não for de cronograma
	ObraID                  *string    `json:"obraId,omitempty"`
	ObraNome                *string    `json:"obraNome,omitempty"`
	Cliente                 string     `json:"cliente"`
	ValorRecebido           float64    `json:"valorRecebido"`
	DataRecebimento         time.Time  `json:"dataRecebimento"`
	FormaPagamento          *string    `json:"formaPagamento,omitempty"`
	Descricao               string     `json:"descricao"`
	ContaBancariaID         *string    `json:"contaBancariaId,omitempty"` // Onde foi depositado
	UsuarioID               string     `json:"usuarioId"` // Quem registrou o recebimento
}