// file: internal/handler/http/pessoal/handler.go
package pessoal

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"

	pessoal_service "github.com/luiszkm/masterCostrutora/internal/service/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/service/pessoal/dto"
)

type Service interface {
	CadastrarFuncionario(ctx context.Context, nome, cpf, cargo, departamento, telefone, chavePix string, diaria float64) (*pessoal.Funcionario, error)
	DeletarFuncionario(ctx context.Context, id string) error
	ListarFuncionarios(ctx context.Context) ([]*pessoal.Funcionario, error)
	AtualizarFuncionario(ctx context.Context, id string, input dto.AtualizarFuncionarioInput) (*pessoal.Funcionario, error)
	BuscarPorID(ctx context.Context, id string) (*pessoal.Funcionario, error)
	CriarApontamento(ctx context.Context, input dto.CriarApontamentoInput) (*pessoal.ApontamentoQuinzenal, error)
	AprovarApontamento(ctx context.Context, apontamentoID string) (*pessoal.ApontamentoQuinzenal, error)
	RegistrarPagamentoApontamento(ctx context.Context, apontamentoID string, contaPagamentoID string) (*pessoal.ApontamentoQuinzenal, error) // NOVO
	ListarApontamentos(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*pessoal.ApontamentoQuinzenal], error)
	ListarApontamentosPorFuncionario(ctx context.Context, funcionarioID string, filtros common.ListarFiltros) (*common.RespostaPaginada[*pessoal.ApontamentoQuinzenal], error)
	ListarComUltimoApontamento(ctx context.Context, filtros common.ListarFiltros) ([]*dto.ListagemFuncionarioDTO, *common.PaginacaoInfo, error)
	AtualizarApontamento(ctx context.Context, id string, input dto.AtualizarApontamentoInput) (*pessoal.ApontamentoQuinzenal, error)
	AtivarFuncionario(ctx context.Context, id string) error
	ReplicarParaProximaQuinzena(ctx context.Context, input dto.ReplicarApontamentosInput) (*dto.ResultadoReplicacao, error)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NovoPessoalHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{service: s, logger: l}
}

func (h *Handler) HandleCadastrarFuncionario(w http.ResponseWriter, r *http.Request) {
	var req cadastrarFuncionarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao decodificar payload", "erro", err)
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	f, err := h.service.CadastrarFuncionario(r.Context(), req.Nome, req.CPF, req.Cargo, req.Departamento, req.Telefone, req.ChavePix, req.Diaria)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao cadastrar funcionário", "erro", err)
		web.RespondError(w, r, "ERRO_CADASTRAR_FUNCIONARIO", "Erro ao cadastrar funcionário", http.StatusInternalServerError)
		return
	}

	h.logger.InfoContext(r.Context(), "Funcionário cadastrado com sucesso", "funcionarioId", f.ID)
	web.Respond(w, r, f, http.StatusCreated)
}

func (h *Handler) HandleDeletarFuncionario(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "funcionarioId")

	err := h.service.DeletarFuncionario(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FUNCIONARIO_NAO_ENCONTRADO", "Funcionário não encontrado", http.StatusNotFound)
			return
		}
		if errors.Is(err, pessoal_service.ErrFuncionarioAlocado) {
			web.RespondError(w, r, "CONFLITO_REGRA_NEGOCIO", err.Error(), http.StatusConflict) // 409 Conflict
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao deletar funcionário", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao deletar funcionário", http.StatusInternalServerError)
		return
	}
	// Resposta padrão para um DELETE bem-sucedido.
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleListarFuncionarios(w http.ResponseWriter, r *http.Request) {
	funcionarios, err := h.service.ListarFuncionarios(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar funcionários", "erro", err)
		web.RespondError(w, r, "ERRO_LISTAR_FUNCIONARIOS", "Erro ao listar funcionários", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, funcionarios, http.StatusOK)
}

func (h *Handler) HandleAtualizarFuncionario(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "funcionarioId")
	if _, err := uuid.Parse(id); err != nil {
		web.RespondError(w, r, "ID_INVALIDO", "O ID do funcionário não é um UUID válido", http.StatusBadRequest)
		return
	}

	var req atualizarFuncionarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	// Mapeia o request do handler para o DTO do serviço.
	input := dto.AtualizarFuncionarioInput{
		Nome:                req.Nome,
		CPF:                 req.CPF,
		Telefone:            req.Telefone,
		Cargo:               req.Cargo,
		Departamento:        req.Departamento,
		ValorDiaria:         req.ValorDiaria,
		ChavePix:            req.ChavePix,
		Status:              req.Status,
		DesligamentoData:    req.DesligamentoData,
		MotivoDesligamento:  req.MotivoDesligamento,
		DataContratacao:     req.DataContratacao,
		Observacoes:         req.Observacoes,
		AvaliacaoDesempenho: req.AvaliacaoDesempenho,
		Email:               req.Email,
	}

	funcionario, err := h.service.AtualizarFuncionario(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FUNCIONARIO_NAO_ENCONTRADO", "Funcionário não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao atualizar funcionário", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao atualizar funcionário", http.StatusInternalServerError)
		return
	}
	// Retorna o objeto completo atualizado.
	web.Respond(w, r, funcionario, http.StatusOK)
}

func (h *Handler) HandleBuscarFuncionario(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "funcionarioId")

	funcionario, err := h.service.BuscarPorID(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FUNCIONARIO_NAO_ENCONTRADO", "Funcionário não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao buscar funcionário", "erro", err)
		web.RespondError(w, r, "ERRO_BUSCAR_FUNCIONARIO", "Erro ao buscar funcionário", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, funcionario, http.StatusOK)
}

func (h *Handler) HandleAtivarFuncionario(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "funcionarioId")

	err := h.service.AtivarFuncionario(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FUNCIONARIO_NAO_ENCONTRADO", "Funcionário não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao buscar funcionário", "erro", err)
		web.RespondError(w, r, "ERRO_BUSCAR_FUNCIONARIO", "Erro ao buscar funcionário", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, nil, http.StatusNoContent)
}
