// file: internal/handler/http/suprimentos/handler.go
package suprimentos

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/luiszkm/masterCostrutora/internal/domain/suprimentos"
	materialDto "github.com/luiszkm/masterCostrutora/internal/handler/http/suprimentos/dtos"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/service/suprimentos/dto"
)

type Service interface {
	CadastrarFornecedor(ctx context.Context, input dto.CadastrarFornecedorInput) (*suprimentos.Fornecedor, error)
	ListarFornecedores(ctx context.Context) ([]*suprimentos.Fornecedor, error)
	CadastrarMaterial(ctx context.Context, input dto.CadastrarMaterialInput) (*suprimentos.Material, error)
	ListarMateriais(ctx context.Context) ([]*suprimentos.Material, error)
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
	var req materialDto.CadastrarMaterialRequest
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
	resp := materialDto.MaterialResponse{
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
	resp := make([]materialDto.MaterialResponse, len(materiais))
	for i, m := range materiais {
		resp[i] = materialDto.MaterialResponse{
			ID:              m.ID,
			Nome:            m.Nome,
			Descricao:       m.Descricao,
			UnidadeDeMedida: m.UnidadeDeMedida,
			Categoria:       m.Categoria,
		}
	}

	web.Respond(w, r, resp, http.StatusOK)
}
