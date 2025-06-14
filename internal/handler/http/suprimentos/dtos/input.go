package dtos

type CadastrarMaterialRequest struct {
	Nome            string `json:"nome"`
	Descricao       string `json:"descricao"`
	UnidadeDeMedida string `json:"unidadeDeMedida"`
	Categoria       string `json:"categoria"`
}

type MaterialResponse struct {
	ID              string `json:"id"`
	Nome            string `json:"nome"`
	Descricao       string `json:"descricao"`
	UnidadeDeMedida string `json:"unidadeDeMedida"`
	Categoria       string `json:"categoria"`
}
