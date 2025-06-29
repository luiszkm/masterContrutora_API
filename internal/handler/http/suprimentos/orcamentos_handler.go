package suprimentos

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	handler_dto "github.com/luiszkm/masterCostrutora/internal/handler/http/suprimentos/dtos"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

func (h *Handler) HandleCriarOrcamento(w http.ResponseWriter, r *http.Request) {
	etapaIDStr := chi.URLParam(r, "etapaId")
	if _, err := uuid.Parse(etapaIDStr); err != nil {
		web.RespondError(w, r, "ID_ETAPA_INVALIDO", "O ID da etapa fornecido na URL não é um UUID válido", http.StatusBadRequest)
		return
	}

	var req handler_dto.CriarOrcamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload da requisição é inválido", http.StatusBadRequest)
		return
	}
	// TODO: Adicionar validação mais robusta para os campos da requisição.
	if req.FornecedorID == "" || len(req.Itens) == 0 {
		web.RespondError(w, r, "DADOS_OBRIGATORIOS", "fornecedorId e pelo menos um item são obrigatórios", http.StatusBadRequest)
		return
	}

	// Converte o DTO do handler para o DTO do serviço
	itensInput := make([]dto.ItemOrcamentoInput, len(req.Itens))
	for i, item := range req.Itens {
		itensInput[i] = dto.ItemOrcamentoInput{
			MaterialID:    item.MaterialID,
			Quantidade:    item.Quantidade,
			ValorUnitario: item.ValorUnitario,
		}
	}
	input := dto.CriarOrcamentoInput{
		Numero:       req.Numero,
		FornecedorID: req.FornecedorID,
		Itens:        itensInput,
	}

	orcamento, err := h.service.CriarOrcamento(r.Context(), etapaIDStr, input)
	if err != nil {
		// Se o erro for de recurso não encontrado (etapa, fornecedor, material), retorna 404.
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "RECURSO_NAO_ENCONTRADO", err.Error(), http.StatusNotFound)
			return
		}
		// Para qualquer outro erro, retorna 500.
		h.logger.ErrorContext(r.Context(), "falha ao criar orçamento", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Não foi possível criar o orçamento", http.StatusInternalServerError)
		return
	}

	// Usa nosso helper para responder com sucesso.
	web.Respond(w, r, orcamento, http.StatusCreated)
}
func (h *Handler) HandleAtualizarOrcamentoStatus(w http.ResponseWriter, r *http.Request) {
	orcamentoID := chi.URLParam(r, "orcamentoId")

	var req handler_dto.AtualizarStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	input := dto.AtualizarStatusOrcamentoInput{Status: req.Status}

	orcamento, err := h.service.AtualizarStatusOrcamento(r.Context(), orcamentoID, input)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "ORCAMENTO_NAO_ENCONTRADO", "Orçamento não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao atualizar orçamento", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao atualizar orçamento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, orcamento, http.StatusOK)
}

func (h *Handler) HandleListarOrcamentos(w http.ResponseWriter, r *http.Request) {
	filtros := web.ParseFiltros(r) // Reutiliza nosso helper!

	resposta, err := h.service.ListarOrcamentos(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar orçamentos", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar orçamentos", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, resposta, http.StatusOK)
}
