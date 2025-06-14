// file: internal/handler/http/suprimentos/handler.go
package suprimentos

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	handler_dto "github.com/luiszkm/masterCostrutora/internal/handler/http/suprimentos/dtos"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

type Service interface {
	CadastrarFornecedor(ctx context.Context, input dto.CadastrarFornecedorInput) (*suprimentos.Fornecedor, error)
	ListarFornecedores(ctx context.Context) ([]*suprimentos.Fornecedor, error)
	CadastrarMaterial(ctx context.Context, input dto.CadastrarMaterialInput) (*suprimentos.Material, error)
	ListarMateriais(ctx context.Context) ([]*suprimentos.Material, error)
	CriarOrcamento(ctx context.Context, etapaID string, input dto.CriarOrcamentoInput) (*suprimentos.Orcamento, error)
	AtualizarStatusOrcamento(ctx context.Context, orcamentoID string, input dto.AtualizarStatusOrcamentoInput) (*suprimentos.Orcamento, error)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NovoSuprimentosHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

type cadastrarRequest struct {
	Nome      string `json:"nome"`
	CNPJ      string `json:"cnpj"`
	Categoria string `json:"categoria"`
	Contato   string `json:"contato"`
	Email     string `json:"email"`
}

func (h *Handler) HandleCadastrarFornecedor(w http.ResponseWriter, r *http.Request) {
	var req cadastrarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	input := dto.CadastrarFornecedorInput{
		Nome:      req.Nome,
		CNPJ:      req.CNPJ,
		Categoria: req.Categoria,
		Contato:   req.Contato,
		Email:     req.Email,
	}

	f, err := h.service.CadastrarFornecedor(r.Context(), input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao cadastrar fornecedor", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao cadastrar fornecedor", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, f, http.StatusCreated)
}

func (h *Handler) HandleListarFornecedores(w http.ResponseWriter, r *http.Request) {
	fornecedores, err := h.service.ListarFornecedores(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar fornecedores", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar fornecedores", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, fornecedores, http.StatusOK)
}

// HandleCadastrarMaterial trata a requisição para criar um novo material.
func (h *Handler) HandleCadastrarMaterial(w http.ResponseWriter, r *http.Request) {
	var req handler_dto.CadastrarMaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload da requisição é inválido", http.StatusBadRequest)
		return
	}

	// TODO: Adicionar validação para os campos da requisição.

	input := dto.CadastrarMaterialInput{
		Nome:            req.Nome,
		Descricao:       req.Descricao,
		UnidadeDeMedida: req.UnidadeDeMedida,
		Categoria:       req.Categoria,
	}

	material, err := h.service.CadastrarMaterial(r.Context(), input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao cadastrar material", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Não foi possível cadastrar o material", http.StatusInternalServerError)
		return
	}

	// Converte o agregado de domínio para um DTO de resposta para não expor o modelo interno.
	resp := handler_dto.MaterialResponse{
		ID:              material.ID,
		Nome:            material.Nome,
		Descricao:       material.Descricao,
		UnidadeDeMedida: material.UnidadeDeMedida,
		Categoria:       material.Categoria,
	}

	web.Respond(w, r, resp, http.StatusCreated)
}

// HandleListarMateriais trata a requisição para listar todos os materiais.
func (h *Handler) HandleListarMateriais(w http.ResponseWriter, r *http.Request) {
	materiais, err := h.service.ListarMateriais(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar materiais", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar os materiais", http.StatusInternalServerError)
		return
	}

	// Converte a lista de agregados de domínio para uma lista de DTOs de resposta.
	resp := make([]handler_dto.MaterialResponse, len(materiais))
	for i, m := range materiais {
		resp[i] = handler_dto.MaterialResponse{
			ID:              m.ID,
			Nome:            m.Nome,
			Descricao:       m.Descricao,
			UnidadeDeMedida: m.UnidadeDeMedida,
			Categoria:       m.Categoria,
		}
	}

	web.Respond(w, r, resp, http.StatusOK)
}

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
