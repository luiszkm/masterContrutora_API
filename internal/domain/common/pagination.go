package common

import "math"

// PaginacaoInfo contém os metadados de uma resposta paginada.
type PaginacaoInfo struct {
	TotalItens    int `json:"totalItens"`
	TotalPaginas  int `json:"totalPages"`
	PaginaAtual   int `json:"currentPage"`
	TamanhoPagina int `json:"pageSize"`
}

// RespostaPaginada é uma estrutura genérica para respostas de lista paginada.
type RespostaPaginada[T any] struct {
	Dados     []T           `json:"dados"`
	Paginacao PaginacaoInfo `json:"paginacao"`
}

// NewPaginacaoInfo é um construtor para facilitar a criação dos metadados de paginação.
func NewPaginacaoInfo(totalItens, pagina, tamanhoPagina int) *PaginacaoInfo {
	if tamanhoPagina <= 0 {
		tamanhoPagina = 1 // Evita divisão por zero
	}
	return &PaginacaoInfo{
		TotalItens:    totalItens,
		TamanhoPagina: tamanhoPagina,
		PaginaAtual:   pagina,
		TotalPaginas:  int(math.Ceil(float64(totalItens) / float64(tamanhoPagina))),
	}
}
