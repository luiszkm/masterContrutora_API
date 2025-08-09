// file: internal/service/pessoal/dto/replicacao_dto.go
package dto

// ReplicarApontamentosInput é o DTO para o comando de replicação.
type ReplicarApontamentosInput struct {
	FuncionarioIDs []string `json:"funcionarioIds"`
}

// ResultadoReplicacao é a estrutura completa da resposta 207 Multi-Status.
type ResultadoReplicacao struct {
	Resumo   ResumoReplicacao `json:"resumo"`
	Sucessos []DetalheSucesso `json:"sucessos"`
	Falhas   []DetalheFalha   `json:"falhas"`
}

// ResumoReplicacao contém os totais da operação.
type ResumoReplicacao struct {
	TotalSolicitado int `json:"totalSolicitado"`
	TotalSucesso    int `json:"totalSucesso"`
	TotalFalha      int `json:"totalFalha"`
}

// DetalheSucesso informa o resultado de uma replicação bem-sucedida.
type DetalheSucesso struct {
	FuncionarioID     string `json:"funcionarioId"`
	NovoApontamentoID string `json:"novoApontamentoId"`
}

// DetalheFalha informa por que a replicação falhou para um funcionário.
type DetalheFalha struct {
	FuncionarioID string `json:"funcionarioId"`
	Motivo        string `json:"motivo"`
}
