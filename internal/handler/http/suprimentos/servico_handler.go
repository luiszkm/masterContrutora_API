package suprimentos

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

func (h *Handler) HandleCriarServico(w http.ResponseWriter, r *http.Request) {
	var input dto.CriarServicoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	servico, err := h.service.CriarServico(r.Context(), input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao criar servico", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Servico não encontrado", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao criar servico", http.StatusInternalServerError)
	}
	web.Respond(w, r, servico, http.StatusCreated)
}

func (h *Handler) HandleListarServicos(w http.ResponseWriter, r *http.Request) {
	filtros := web.ParseFiltros(r)

	if filtros.Pagina == 0 && filtros.TamanhoPagina == 0 {
		servicos, err := h.service.ListarServicos(r.Context())
		if err != nil {
			h.logger.ErrorContext(r.Context(), "falha ao listar servicos", "erro", err)
			web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar servicos", http.StatusInternalServerError)
			return
		}
		web.Respond(w, r, servicos, http.StatusOK)
		return
	}

	resposta, err := h.service.ListarServicosPaginado(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar servicos paginados", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar servicos", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, resposta, http.StatusOK)
}

func (h *Handler) HandleBuscarServico(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "servicoId")
	servico, err := h.service.BuscarServico(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao buscar servico", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Servico não encontrado", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao buscar servico", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, servico, http.StatusOK)
}

func (h *Handler) HandleAtualizarServico(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "servicoId")
	var input dto.AtualizarServicoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	servico, err := h.service.AtualizarServico(r.Context(), id, input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao atualizar servico", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Servico não encontrado", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao atualizar servico", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, servico, http.StatusOK)
}

func (h *Handler) HandleDeletarServico(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "servicoId")
	err := h.service.DeletarServico(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao deletar servico", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Servico não encontrado", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			web.RespondError(w, r, "CONFLITO", "O servico está em uso e não pode ser deletado", http.StatusConflict)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao deletar servico", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}