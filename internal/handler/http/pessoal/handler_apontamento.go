package pessoal

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	pessoal_service "github.com/luiszkm/masterCostrutora/internal/service/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/service/pessoal/dto"
)

func (h *Handler) HandleCriarApontamento(w http.ResponseWriter, r *http.Request) {
	var req dto.CriarApontamentoInput // reusando o DTO do serviço por simplicidade
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", err.Error(), http.StatusBadRequest)
		return
	}

	apontamento, err := h.service.CriarApontamento(r.Context(), req)
	if err != nil {
		// TODO: Tratar erros específicos, como 404
		h.logger.ErrorContext(r.Context(), "falha ao criar apontamento", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", err.Error(), http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, apontamento, http.StatusCreated)
}
func (h *Handler) HandleAprovarApontamento(w http.ResponseWriter, r *http.Request) {
	apontamentoID := chi.URLParam(r, "apontamentoId")
	// Não há corpo na requisição, a ação é implícita pelo endpoint.

	apontamento, err := h.service.AprovarApontamento(r.Context(), apontamentoID)
	if err != nil {
		// Trata erros específicos, como 404 ou 409 (conflito de regra de negócio)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "APONTAMENTO_NAO_ENCONTRADO", "Apontamento não encontrado", http.StatusNotFound)
			return
		}
		// Se o erro veio da regra de negócio do domínio (ex: tentar aprovar algo já pago)
		web.RespondError(w, r, "REGRA_NEGOCIO_VIOLADA", err.Error(), http.StatusConflict)
		return
	}

	web.Respond(w, r, apontamento, http.StatusOK)
}

func (h *Handler) HandleRegistrarPagamentoApontamento(w http.ResponseWriter, r *http.Request) {
	apontamentoID := chi.URLParam(r, "apontamentoId")

	var req registrarPagamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}
	if req.ContaBancariaID == "" {
		web.RespondError(w, r, "DADOS_OBRIGATORIOS", "O campo contaBancariaId é obrigatório", http.StatusBadRequest)
		return
	}

	apontamento, err := h.service.RegistrarPagamentoApontamento(r.Context(), apontamentoID, req.ContaBancariaID)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "APONTAMENTO_NAO_ENCONTRADO", "Apontamento não encontrado", http.StatusNotFound)
			return
		}
		// Trata o erro de regra de negócio vindo do método .RegistrarPagamento() do agregado.
		if errors.Is(err, pessoal_service.ErrFuncionarioAlocado) || err.Error() == "só é possível pagar um apontamento que está 'Aprovado para Pagamento'" {
			web.RespondError(w, r, "REGRA_NEGOCIO_VIOLADA", err.Error(), http.StatusConflict) // 409 Conflict
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao registrar pagamento de apontamento", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao registrar pagamento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, apontamento, http.StatusOK)
}

func (h *Handler) HandleListarApontamentos(w http.ResponseWriter, r *http.Request) {
	// Lógica para extrair filtros da URL (status, page, pageSize)
	filtros := web.ParseFiltros(r) // Reutilizando a função de parsing de filtros

	resposta, err := h.service.ListarApontamentos(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar apontamentos", "erro", err)
		web.RespondError(w, r, "ERRO_LISTAR_APONTAMENTOS", "Erro ao listar apontamentos", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, resposta, http.StatusOK)
}

func (h *Handler) HandleListarApontamentosPorFuncionario(w http.ResponseWriter, r *http.Request) {
	funcionarioID := chi.URLParam(r, "funcionarioId")
	// Lógica para extrair filtros da URL
	filtros := web.ParseFiltros(r) // Reutilizando a função de parsing de filtros

	resposta, err := h.service.ListarApontamentosPorFuncionario(r.Context(), funcionarioID, filtros)
	if err != nil { /* ... tratamento de erro */
	}

	web.Respond(w, r, resposta, http.StatusOK)
}

func (h *Handler) HandleListarComUltimoApontamento(w http.ResponseWriter, r *http.Request) {
	// 1. Reutiliza nossa função helper para extrair os filtros da requisição.
	filtros := web.ParseFiltros(r)

	// 2. Chama o método de serviço correspondente.
	respostaPaginada, _, err := h.service.ListarComUltimoApontamento(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar funcionários com apontamentos", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar funcionários com apontamentos", http.StatusInternalServerError)
		return
	}
	log.Println("Resposta paginada:", respostaPaginada)

	web.Respond(w, r, respostaPaginada, http.StatusOK)
}

func (h *Handler) HandleAtualizarApontamento(w http.ResponseWriter, r *http.Request) {
	apontamentoID := chi.URLParam(r, "apontamentoId")
	if _, err := uuid.Parse(apontamentoID); err != nil {
		web.RespondError(w, r, "ID_INVALIDO", "O ID do apontamento não é um UUID válido", http.StatusBadRequest)
		return
	}

	var req dto.AtualizarApontamentoInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	apontamento, err := h.service.AtualizarApontamento(r.Context(), apontamentoID, req)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "APONTAMENTO_NAO_ENCONTRADO", "Apontamento não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao atualizar apontamento", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao atualizar apontamento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, apontamento, http.StatusOK)

}
