// file: internal/handler/http/pessoal/handler.go
package pessoal

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
)

type Service interface {
	CadastrarFuncionario(ctx context.Context, nome, cpf, cargo string, salario, diaria float64) (*pessoal.Funcionario, error)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NovoPessoalHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

type cadastrarFuncionarioRequest struct {
	Nome    string  `json:"nome"`
	CPF     string  `json:"cpf"`
	Cargo   string  `json:"cargo"`
	Salario float64 `json:"salario"`
	Diaria  float64 `json:"diaria"`
}

func (h *Handler) HandleCadastrarFuncionario(w http.ResponseWriter, r *http.Request) {
	var req cadastrarFuncionarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	f, err := h.service.CadastrarFuncionario(r.Context(), req.Nome, req.CPF, req.Cargo, req.Salario, req.Diaria)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao cadastrar funcionário", "erro", err)
		http.Error(w, "Erro ao cadastrar funcionário", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}
