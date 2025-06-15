// file: internal/handler/http/pessoal/handler.go
package pessoal

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"

	pessoal_service "github.com/luiszkm/masterCostrutora/internal/service/pessoal"
)

type Service interface {
	CadastrarFuncionario(ctx context.Context, nome, cpf, cargo string, salario, diaria float64) (*pessoal.Funcionario, error)
	DeletarFuncionario(ctx context.Context, id string) error
	ListarFuncionarios(ctx context.Context) ([]*pessoal.Funcionario, error)
	AtualizarFuncionario(ctx context.Context, funcionario *pessoal.Funcionario) error // NOVO
	BuscarPorID(ctx context.Context, id string) (*pessoal.Funcionario, error)
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

func (h *Handler) HandleDeletarFuncionario(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "funcionarioId")

	err := h.service.DeletarFuncionario(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FUNCIONARIO_NAO_ENCONTRADO", "Funcionário não encontrado", http.StatusNotFound)
			return
		}
		if errors.Is(err, pessoal_service.ErrFuncionarioAlocado) {
			web.RespondError(w, r, "CONFLITO_REGRA_NEGOCIO", err.Error(), http.StatusConflict) // 409 Conflict
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao deletar funcionário", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao deletar funcionário", http.StatusInternalServerError)
		return
	}
	// Resposta padrão para um DELETE bem-sucedido.
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleListarFuncionarios(w http.ResponseWriter, r *http.Request) {
	funcionarios, err := h.service.ListarFuncionarios(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar funcionários", "erro", err)
		http.Error(w, "Erro ao listar funcionários", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(funcionarios)
}

func (h *Handler) HandleAtualizarFuncionario(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "funcionarioId")

	var req cadastrarFuncionarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	funcionario := &pessoal.Funcionario{
		ID:      id,
		Nome:    req.Nome,
		CPF:     req.CPF,
		Cargo:   req.Cargo,
		Salario: req.Salario,
		Diaria:  req.Diaria,
	}

	if err := h.service.AtualizarFuncionario(r.Context(), funcionario); err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FUNCIONARIO_NAO_ENCONTRADO", "Funcionário não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao atualizar funcionário", "erro", err)
		http.Error(w, "Erro ao atualizar funcionário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleBuscarFuncionario(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "funcionarioId")

	funcionario, err := h.service.BuscarPorID(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FUNCIONARIO_NAO_ENCONTRADO", "Funcionário não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao buscar funcionário", "erro", err)
		http.Error(w, "Erro ao buscar funcionário", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(funcionario)
}
