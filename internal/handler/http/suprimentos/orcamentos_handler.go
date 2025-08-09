package suprimentos

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/suprimentos/dtos"
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
			NomeProduto:     item.NomeProduto,
			UnidadeDeMedida: item.UnidadeDeMedida,
			Categoria:       item.Categoria,
			Quantidade:      item.Quantidade,
			ValorUnitario:   item.ValorUnitario,
		}
	}
	input := dto.CriarOrcamentoInput{
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
	// A lógica de parsing dos filtros permanece a mesma.
	filtros := web.ParseFiltros(r)

	// A chamada ao serviço agora retorna uma única variável 'resposta'.
	resposta, err := h.service.ListarOrcamentos(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar orçamentos", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar orçamentos", http.StatusInternalServerError)
		return
	}

	// A resposta do serviço já está no formato correto para a API.
	// Simplesmente a passamos para a função de resposta.
	web.Respond(w, r, resposta, http.StatusOK)
}
func (h *Handler) HandleBuscarOrcamentoPorID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "orcamentoId")

	orcamento, err := h.service.BuscarOrcamentoPorID(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Orçamento não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao buscar orçamento por id", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Falha ao buscar orçamento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, orcamento, http.StatusOK)
}
func (h *Handler) HandleAtualizarOrcamento(w http.ResponseWriter, r *http.Request) {
	orcamentoID := chi.URLParam(r, "orcamentoId")
	if orcamentoID == "" {
		web.RespondError(w, r, "ID_ORCAMENTO_OBRIGATORIO", "O ID do orçamento é obrigatório.", http.StatusBadRequest)
		return
	}

	// 1. Decodifica o payload para o DTO de requisição do handler.
	var req dtos.AtualizarOrcamentoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	// 2. Mapeia os itens do DTO de requisição para o DTO do serviço.
	itensInput := make([]dto.ItemOrcamentoInput, len(req.Itens))
	for i, item := range req.Itens {
		itensInput[i] = dto.ItemOrcamentoInput{
			NomeProduto:     item.NomeProduto,
			UnidadeDeMedida: item.UnidadeDeMedida,
			Categoria:       item.Categoria,
			Quantidade:      item.Quantidade,
			ValorUnitario:   item.ValorUnitario,
		}
	}

	// 3. Mapeia o resto dos campos do DTO de requisição para o DTO de serviço.
	// ESTA É A ETAPA CRÍTICA QUE ESTAVA EM FALTA.
	input := dto.AtualizarOrcamentoInput{
		FornecedorID:       req.FornecedorID,
		EtapaID:            req.EtapaID,
		Observacoes:        req.Observacoes,
		CondicoesPagamento: req.CondicoesPagamento,
		Itens:              itensInput,
	}

	// 4. Chama o serviço com o DTO completamente preenchido.
	orcamentoAtualizado, err := h.service.AtualizarOrcamento(r.Context(), orcamentoID, input)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Orçamento, etapa ou fornecedor não encontrado.", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao atualizar orçamento", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao atualizar orçamento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, orcamentoAtualizado, http.StatusOK)
}
