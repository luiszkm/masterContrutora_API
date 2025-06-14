// file: internal/handler/web/web.go
package web

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse é a estrutura padronizada para erros da API, conforme a documentação.
type ErrorResponse struct {
	Codigo   string `json:"codigo"`
	Mensagem string `json:"mensagem"`
}

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
