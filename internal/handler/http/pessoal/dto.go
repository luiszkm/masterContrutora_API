package pessoal

type registrarPagamentoRequest struct {
	ContaBancariaID string `json:"contaBancariaId"`
}

type cadastrarFuncionarioRequest struct {
	Nome         string  `json:"nome"`
	CPF          string  `json:"cpf"`
	Cargo        string  `json:"cargo"`
	Departamento string  `json:"departamento"` // Adicionando o campo Departamento
	Diaria       float64 `json:"diaria"`       // Adicionando o campo Diaria
	ChavePix     string  `json:"chavePix"`
	Observacoes  string  `json:"observacoes"`
	Telefone     string  `json:"telefone"`
}
type atualizarFuncionarioRequest struct {
	Nome                *string  `json:"nome,omitempty"`
	CPF                 *string  `json:"cpf,omitempty"`
	Cargo               *string  `json:"cargo,omitempty"`
	Departamento        *string  `json:"departamento,omitempty"`
	ValorDiaria         *float64 `json:"valorDiaria,omitempty"`
	ChavePix            *string  `json:"chavePix,omitempty"`
	Status              *string  `json:"status,omitempty"`
	Telefone            *string  `json:"telefone,omitempty"`
	MotivoDesligamento  *string  `json:"motivoDesligamento,omitempty"`
	DataContratacao     *string  `json:"dataContratacao,omitempty"`
	DesligamentoData    *string  `json:"desligamentoData,omitempty"`
	Observacoes         *string  `json:"observacoes,omitempty"`
	AvaliacaoDesempenho *string  `json:"avaliacaoDesempenho,omitempty"`
	Diaria              *float64 `json:"diaria,omitempty"` // Adicionando o campo Diaria
	Email               *string  `json:"email,omitempty"`
}
