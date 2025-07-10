package dto

// CriarNovaObraInput representa os dados necessários para criar uma obra.
type CriarNovaObraInput struct {
	Nome       string `json:"nome"`
	Cliente    string `json:"cliente"`
	Endereco   string `json:"endereco"`
	DataInicio string `json:"dataInicio"` // Espera-se "YYYY-MM-DD"
	DataFim    string `json:"dataFim"`    // Espera-se "YYYY-MM-DD"
	Descricao  string `json:"descricao"`
}

type AtualizarObraInput struct {
	Nome       string `json:"nome"`
	Cliente    string `json:"cliente"`
	Endereco   string `json:"endereco"`
	DataInicio string `json:"dataInicio"` // Espera-se "YYYY-MM-DD"
	DataFim    string `json:"dataFim"`    // Espera-se "YYYY-MM-DD"
	Descricao  string `json:"descricao"`
	Status     string `json:"status"`
}
