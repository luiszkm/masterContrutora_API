package pessoal


type cadastrarFuncionarioRequest struct {
	Nome         string `json:"nome"`
	CPF          string `json:"cpf"`
	Cargo        string `json:"cargo"`
	Departamento string `json:"departamento"`
	ChavePix     string `json:"chavePix"`
	Observacoes  string `json:"observacoes"`
	Telefone     string `json:"telefone"`
}
type atualizarFuncionarioRequest struct {
	Nome                *string `json:"nome,omitempty"`
	CPF                 *string `json:"cpf,omitempty"`
	Cargo               *string `json:"cargo,omitempty"`
	Departamento        *string `json:"departamento,omitempty"`
	ChavePix            *string `json:"chavePix,omitempty"`
	Status              *string `json:"status,omitempty"`
	Telefone            *string `json:"telefone,omitempty"`
	MotivoDesligamento  *string `json:"motivoDesligamento,omitempty"`
	DataContratacao     *string `json:"dataContratacao,omitempty"`
	DesligamentoData    *string `json:"desligamentoData,omitempty"`
	Observacoes         *string `json:"observacoes,omitempty"`
	AvaliacaoDesempenho *string `json:"avaliacaoDesempenho,omitempty"`
	Email               *string `json:"email,omitempty"`
}
