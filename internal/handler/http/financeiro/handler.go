// file: internal/handler/http/financeiro/handler.go
package financeiro

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/luiszkm/masterCostrutora/internal/domain/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	financeiro_service "github.com/luiszkm/masterCostrutora/internal/service/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/service/financeiro/dto"
)

type Service interface {
	RegistrarPagamento(ctx context.Context, input dto.RegistrarPagamentoInput) (*financeiro.RegistroDePagamento, error)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NovoFinanceiroHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

type registrarPagamentoRequest struct {
	FuncionarioID     string  `json:"funcionarioId"`
	ObraID            string  `json:"obraId"`
	PeriodoReferencia string  `json:"periodoReferencia"`
	ValorCalculado    float64 `json:"valorCalculado"`
	ContaBancariaID   string  `json:"contaBancariaId"`
}

func (h *Handler) HandleRegistrarPagamento(w http.ResponseWriter, r *http.Request) {
	var req registrarPagamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	input := dto.RegistrarPagamentoInput{
		FuncionarioID:     req.FuncionarioID,
		ObraID:            req.ObraID,
		PeriodoReferencia: req.PeriodoReferencia,
		ValorCalculado:    req.ValorCalculado,
		ContaBancariaID:   req.ContaBancariaID,
	}

	pagamento, err := h.service.RegistrarPagamento(r.Context(), input)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "RECURSO_NAO_ENCONTRADO", "Funcionário ou Obra não encontrado(a)", http.StatusNotFound)
			return
		}
		if errors.Is(err, financeiro_service.ErrFuncionarioInativo) {
			web.RespondError(w, r, "REGRA_NEGOCIO_VIOLADA", err.Error(), http.StatusUnprocessableEntity) // 422
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao registrar pagamento", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao registrar pagamento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, pagamento, http.StatusCreated)
}
