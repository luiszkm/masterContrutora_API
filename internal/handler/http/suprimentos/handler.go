// file: internal/handler/http/suprimentos/handler.go
package suprimentos

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
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
	AtualizarFornecedor(ctx context.Context, id string, input dto.AtualizarFornecedorInput) (*suprimentos.Fornecedor, error)
	DeletarFornecedor(ctx context.Context, id string) error
	BuscarPorID(ctx context.Context, id string) (*suprimentos.Fornecedor, error)
	CriarCategoria(ctx context.Context, nome string) (*suprimentos.Categoria, error)
	ListarCategorias(ctx context.Context) ([]*suprimentos.Categoria, error)
	ListarOrcamentos(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*suprimentos.Orcamento], error)
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

func (h *Handler) HandleAtualizarFornecedor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.RespondError(w, r, "ID_FORNECEDOR_INVALIDO", "O ID do fornecedor não pode ser vazio", http.StatusBadRequest)
		return
	}

	var req handler_dto.AtualizarFornecedorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	input := dto.AtualizarFornecedorInput{
		Nome:      req.Nome,
		CNPJ:      req.CNPJ,
		Categoria: req.Categoria,
		Contato:   req.Contato,
		Email:     req.Email,
	}

	fornecedor, err := h.service.AtualizarFornecedor(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FORNECEDOR_NAO_ENCONTRADO", "Fornecedor não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao atualizar fornecedor", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao atualizar fornecedor", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, fornecedor, http.StatusOK)
}

func (h *Handler) HandleDeletarFornecedor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.RespondError(w, r, "ID_FORNECEDOR_INVALIDO", "O ID do fornecedor não pode ser vazio", http.StatusBadRequest)
		return
	}

	if err := h.service.DeletarFornecedor(r.Context(), id); err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FORNECEDOR_NAO_ENCONTRADO", "Fornecedor não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao deletar fornecedor", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao deletar fornecedor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (h *Handler) HandleBuscarFornecedor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		web.RespondError(w, r, "ID_FORNECEDOR_INVALIDO", "O ID do fornecedor não pode ser vazio", http.StatusBadRequest)
		return
	}

	fornecedor, err := h.service.BuscarPorID(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FORNECEDOR_NAO_ENCONTRADO", "Fornecedor não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao buscar fornecedor", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao buscar fornecedor", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, fornecedor, http.StatusOK)
}

func (h *Handler) HandleCriarCategoria(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nome string `json:"nome"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}
	if req.Nome == "" {
		web.RespondError(w, r, "NOME_CATEGORIA_OBRIGATORIO", "O nome da categoria é obrigatório", http.StatusBadRequest)
		return
	}
	_, err := h.service.CriarCategoria(r.Context(), req.Nome)
	if err != nil {
		if errors.Is(err, postgres.ErrCategoriaJaExiste) {
			web.RespondError(w, r, "CATEGORIA_JA_EXISTE", "Categoria já existe", http.StatusConflict)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao criar categoria", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao criar categoria", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, map[string]string{"message": "Categoria criada com sucesso"}, http.StatusCreated)
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
