package dto

type CriarServicoInput struct {
	Nome string `json:"nome"`
}

type AtualizarServicoInput struct {
	Nome string `json:"nome"`
}