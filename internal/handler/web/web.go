// file: internal/handler/web/web.go
package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/luiszkm/masterCostrutora/internal/domain/common"
)

// ErrorResponse é a estrutura padronizada para erros da API, conforme a documentação.
type ErrorResponse struct {
	Codigo   string `json:"codigo"`
	Mensagem string `json:"mensagem"`
}

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// Respond converte um payload Go para JSON e o escreve na resposta HTTP.
// Esta função centraliza a lógica de resposta de sucesso.
func Respond(w http.ResponseWriter, r *http.Request, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// Se a codificação falhar, logamos o erro e enviamos um erro HTTP genérico.
			// Em um sistema em produção, usaríamos um logger estruturado aqui.
			http.Error(w, "could not encode response to json", http.StatusInternalServerError)
		}
	}
}

// RespondError envia uma resposta de erro JSON padronizada.
// Esta função centraliza toda a lógica de formatação de erro.
func RespondError(w http.ResponseWriter, r *http.Request, codigo string, mensagem string, statusCode int) {
	errResponse := ErrorResponse{
		Codigo:   codigo,
		Mensagem: mensagem,
	}
	Respond(w, r, errResponse, statusCode)
}

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

func ParseFiltros(r *http.Request) common.ListarFiltros {
	q := r.URL.Query()

	// Parse do status (simples, pois é uma string)
	status := q.Get("status")

	// Parse da página com valor padrão
	pagina, err := strconv.Atoi(q.Get("page"))
	if err != nil || pagina < 1 {
		pagina = defaultPage
	}

	// Parse do tamanho da página com valor padrão e limite máximo
	tamanhoPagina, err := strconv.Atoi(q.Get("pageSize"))

	if err != nil || tamanhoPagina < 1 {
		tamanhoPagina = defaultPageSize
	}
	if tamanhoPagina > maxPageSize {
		tamanhoPagina = maxPageSize
	}

	fornecedorID := q.Get("fornecedorId")
	obraID := q.Get("obraId")
	return common.ListarFiltros{
		Status:        status,
		Pagina:        pagina,
		FornecedorID:  fornecedorID,
		ObraID:        obraID,
		TamanhoPagina: tamanhoPagina,
	}
}
