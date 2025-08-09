// file: internal/service/financeiro/dto/pagamento_lote_dto.go
package dto

// RegistrarPagamentoEmLoteInput é o DTO para o comando de pagamento em lote.
type RegistrarPagamentoEmLoteInput struct {
	ApontamentoIDs   []string `json:"apontamentoIds"`
	ContaBancariaID  string   `json:"contaBancariaId"`
	DataDeEfetivacao string   `json:"dataDeEfetivacao"` // Formato "YYYY-MM-DD"
}

// ResultadoExecucaoLote é a estrutura de resposta para o 207 Multi-Status.
type ResultadoExecucaoLote struct {
	Resumo   ResumoExecucao   `json:"resumo"`
	Sucessos []DetalheSucesso `json:"sucessos"`
	Falhas   []DetalheFalha   `json:"falhas"`
}

type ResumoExecucao struct {
	TotalSolicitado int `json:"totalSolicitado"`
	TotalSucesso    int `json:"totalSucesso"`
	TotalFalha      int `json:"totalFalha"`
}

type DetalheSucesso struct {
	ApontamentoID       string `json:"apontamentoId"`
	RegistroPagamentoID string `json:"registroPagamentoId"`
}

type DetalheFalha struct {
	ApontamentoID string `json:"apontamentoId"`
	Motivo        string `json:"motivo"`
}
