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
	CadastrarMaterial(ctx context.Context, input dto.CadastrarProdutoInput) (*suprimentos.Produto, error)
	ListarMateriais(ctx context.Context) ([]*suprimentos.Produto, error)
	BuscarMaterialPorID(ctx context.Context, id string) (*suprimentos.Produto, error)
	AtualizarMaterial(ctx context.Context, id string, input dto.CadastrarProdutoInput) (*suprimentos.Produto, error)
	DeletarMaterial(ctx context.Context, id string) error
	DeletarOrcamento(ctx context.Context, id string) error
	CriarOrcamento(ctx context.Context, etapaID string, input dto.CriarOrcamentoInput) (*suprimentos.Orcamento, error)
	AtualizarStatusOrcamento(ctx context.Context, orcamentoID string, input dto.AtualizarStatusOrcamentoInput) (*suprimentos.Orcamento, error)
	AtualizarFornecedor(ctx context.Context, id string, input dto.AtualizarFornecedorInput) (*suprimentos.Fornecedor, error)
	DeletarFornecedor(ctx context.Context, id string) error
	BuscarPorID(ctx context.Context, id string) (*suprimentos.Fornecedor, error)
	CriarCategoria(ctx context.Context, input dto.CriarCategoriaInput) (*suprimentos.Categoria, error)
	ListarCategorias(ctx context.Context) ([]*suprimentos.Categoria, error)
	BuscarCategoria(ctx context.Context, id string) (*suprimentos.Categoria, error)
	AtualizarCategoria(ctx context.Context, id string, input dto.AtualizarCategoriaInput) (*suprimentos.Categoria, error)
	DeletarCategoria(ctx context.Context, id string) error
	ListarOrcamentos(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*dto.OrcamentoListItemDTO], error)
	BuscarOrcamentoPorID(ctx context.Context, id string) (*dto.OrcamentoDetalhadoDTO, error)
	AtualizarOrcamento(ctx context.Context, orcamentoID string, input dto.AtualizarOrcamentoInput) (*dto.OrcamentoDetalhadoDTO, error)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NovoSuprimentosHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

func (h *Handler) HandleCadastrarFornecedor(w http.ResponseWriter, r *http.Request) {
	var req handler_dto.CadastrarFornecedorRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	input := dto.CadastrarFornecedorInput{
		Nome:         req.Nome,
		CNPJ:         req.CNPJ,
		CategoriaIDs: req.CategoriaIDs,
		Contato:      req.Contato,
		Email:        req.Email,
		Endereco:     req.Endereco,
		Avaliacao:    req.Avaliacao,
		Observacoes:  req.Observacoes,
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

	input := dto.CadastrarProdutoInput{
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
		Nome:         req.Nome,
		CNPJ:         req.CNPJ,
		CategoriaIDs: req.CategoriaIDs,
		Contato:      req.Contato,
		Email:        req.Email,
		Status:       req.Status,
		Endereco:     req.Endereco,
		Avaliacao:    req.Avaliacao,
		Observacoes:  req.Observacoes,
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

func (h *Handler) HandleBuscarMaterial(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "materialId")
	material, err := h.service.BuscarMaterialPorID(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao buscar material", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Material não encontrado", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao buscar material", http.StatusInternalServerError)
		return
	}
	
	// Converte para o DTO de resposta
	resp := handler_dto.MaterialResponse{
		ID:              material.ID,
		Nome:            material.Nome,
		Descricao:       material.Descricao,
		UnidadeDeMedida: material.UnidadeDeMedida,
		Categoria:       material.Categoria,
	}
	web.Respond(w, r, resp, http.StatusOK)
}

func (h *Handler) HandleAtualizarMaterial(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "materialId")
	var req handler_dto.CadastrarMaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	input := dto.CadastrarProdutoInput{
		Nome:            req.Nome,
		Descricao:       req.Descricao,
		UnidadeDeMedida: req.UnidadeDeMedida,
		Categoria:       req.Categoria,
	}

	material, err := h.service.AtualizarMaterial(r.Context(), id, input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao atualizar material", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Material não encontrado", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao atualizar material", http.StatusInternalServerError)
		return
	}

	// Converte para o DTO de resposta
	resp := handler_dto.MaterialResponse{
		ID:              material.ID,
		Nome:            material.Nome,
		Descricao:       material.Descricao,
		UnidadeDeMedida: material.UnidadeDeMedida,
		Categoria:       material.Categoria,
	}
	web.Respond(w, r, resp, http.StatusOK)
}

func (h *Handler) HandleDeletarMaterial(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "materialId")
	err := h.service.DeletarMaterial(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao deletar material (soft delete)", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Material não encontrado", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao deletar material", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleDeletarOrcamento(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "orcamentoId")
	err := h.service.DeletarOrcamento(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao deletar orçamento (soft delete)", "erro", err)
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "NAO_ENCONTRADO", "Orçamento não encontrado", http.StatusNotFound)
			return
		}
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao deletar orçamento", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
