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

func (h *Handler) HandleCriarCategoria(w http.ResponseWriter, r *http.Request) {
	var input dto.CriarCategoriaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	categoria, err := h.service.CriarCategoria(r.Context(), input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao criar categoria", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Categoria não encontrada", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao criar categoria", http.StatusInternalServerError)
	}
	web.Respond(w, r, categoria, http.StatusCreated)
}

func (h *Handler) HandleListarCategorias(w http.ResponseWriter, r *http.Request) {
	categorias, err := h.service.ListarCategorias(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar categorias", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar categorias", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, categorias, http.StatusOK)
}

func (h *Handler) HandleBuscarCategoria(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "categoriaId")
	categoria, err := h.service.BuscarCategoria(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao buscar categoria", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Categoria não encontrada", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao buscar categoria", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, categoria, http.StatusOK)
}

func (h *Handler) HandleAtualizarCategoria(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "categoriaId")
	var input dto.AtualizarCategoriaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	categoria, err := h.service.AtualizarCategoria(r.Context(), id, input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao atualizar categoria", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Categoria não encontrada", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao atualizar categoria", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, categoria, http.StatusOK)
}

func (h *Handler) HandleDeletarCategoria(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "categoriaId")
	err := h.service.DeletarCategoria(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao deletar categoria", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Categoria não encontrada", http.StatusNotFound)
			return
		}
		// Tratar erro de chave estrangeira
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			web.RespondError(w, r, "CONFLITO", "A categoria está em uso e não pode ser deletada", http.StatusConflict)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao deletar categoria", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
