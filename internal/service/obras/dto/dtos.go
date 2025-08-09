package dto

import "time"

// CriarNovaObraInput representa os dados necessários para criar uma obra.
type CriarNovaObraInput struct {
	Nome       string `json:"nome"`
	Cliente    string `json:"cliente"`
	Endereco   string `json:"endereco"`
	DataInicio string `json:"dataInicio"` // Espera-se "YYYY-MM-DD"
	DataFim    string `json:"dataFim"`    // Espera-se "YYYY-MM-DD"
	Descricao  string `json:"descricao"`
	
	// Campos financeiros (opcionais na criação)
	ValorContratoTotal     *float64   `json:"valorContratoTotal,omitempty"`
	TipoCobranca           *string    `json:"tipoCobranca,omitempty"` // "VISTA", "PARCELADO", "ETAPAS"
	DataAssinaturaContrato *time.Time `json:"dataAssinaturaContrato,omitempty"`
}

type AtualizarObraInput struct {
	Nome       string `json:"nome"`
	Cliente    string `json:"cliente"`
	Endereco   string `json:"endereco"`
	DataInicio string `json:"dataInicio"` // Espera-se "YYYY-MM-DD"
	DataFim    string `json:"dataFim"`    // Espera-se "YYYY-MM-DD"
	Descricao  string `json:"descricao"`
	Status     string `json:"status"`
	
	// Campos financeiros
	ValorContratoTotal     *float64   `json:"valorContratoTotal,omitempty"`
	TipoCobranca           *string    `json:"tipoCobranca,omitempty"`
	DataAssinaturaContrato *time.Time `json:"dataAssinaturaContrato,omitempty"`
}

// AtualizarValoresContratoInput permite atualizar apenas valores financeiros
type AtualizarValoresContratoInput struct {
	ValorContratoTotal     float64    `json:"valorContratoTotal" validate:"required,gt=0"`
	TipoCobranca           string     `json:"tipoCobranca" validate:"required,oneof=VISTA PARCELADO ETAPAS"`
	DataAssinaturaContrato *time.Time `json:"dataAssinaturaContrato,omitempty"`
}

// CriarEtapaPadraoInput representa os dados necessários para criar uma etapa padrão.
type CriarEtapaPadraoInput struct {
	Nome      string  `json:"nome"`
	Descricao *string `json:"descricao,omitempty"`
	Ordem     int     `json:"ordem"`
}

// AtualizarEtapaPadraoInput representa os dados necessários para atualizar uma etapa padrão.
type AtualizarEtapaPadraoInput struct {
	Nome      string  `json:"nome"`
	Descricao *string `json:"descricao,omitempty"`
	Ordem     int     `json:"ordem"`
}
