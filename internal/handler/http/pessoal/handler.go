// file: internal/handler/http/pessoal/handler.go
package pessoal

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"

	pessoal_service "github.com/luiszkm/masterCostrutora/internal/service/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/service/pessoal/dto"
	pessoal_dto "github.com/luiszkm/masterCostrutora/internal/service/pessoal/dto"
)

type Service interface {
	CadastrarFuncionario(ctx context.Context, nome, cpf, cargo, departamento string, diaria float64) (*pessoal.Funcionario, error)
	DeletarFuncionario(ctx context.Context, id string) error
	ListarFuncionarios(ctx context.Context) ([]*pessoal.Funcionario, error)
	AtualizarFuncionario(ctx context.Context, id string, input dto.AtualizarFuncionarioInput) error // NOVO
	BuscarPorID(ctx context.Context, id string) (*pessoal.Funcionario, error)
	CriarApontamento(ctx context.Context, input dto.CriarApontamentoInput) (*pessoal.ApontamentoQuinzenal, error)
	AprovarApontamento(ctx context.Context, apontamentoID string) (*pessoal.ApontamentoQuinzenal, error)
	RegistrarPagamentoApontamento(ctx context.Context, apontamentoID string, contaPagamentoID string) (*pessoal.ApontamentoQuinzenal, error) // NOVO

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

	f, err := h.service.CadastrarFuncionario(r.Context(), req.Nome, req.CPF, req.Cargo, req.Departamento, req.Diaria)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao cadastrar funcionário", "erro", err)
		web.RespondError(w, r, "ERRO_CADASTRAR_FUNCIONARIO", "Erro ao cadastrar funcionário", http.StatusInternalServerError)
		return
	}

	h.logger.InfoContext(r.Context(), "Funcionário cadastrado com sucesso", "funcionarioId", f.ID)
	web.Respond(w, r, f, http.StatusCreated)
}

func (h *Handler) HandleDeletarFuncionario(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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
	id := chi.URLParam(r, "id")

	var req atualizarFuncionarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao decodificar payload", "erro", err)
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}
	input := pessoal_dto.AtualizarFuncionarioInput{
		Nome:                req.Nome,
		CPF:                 req.CPF,
		Cargo:               req.Cargo,
		Departamento:        req.Departamento,
		ValorDiaria:         req.Diaria,
		ChavePix:            req.ChavePix,
		Status:              req.Status,
		Telefone:            req.Telefone,
		MotivoDesligamento:  req.MotivoDesligamento,
		DataContratacao:     req.DataContratacao,
		DesligamentoData:    req.DesligamentoData,
		Observacoes:         req.Observacoes,
		AvaliacaoDesempenho: req.AvaliacaoDesempenho,
	}

	if err := h.service.AtualizarFuncionario(r.Context(), id, input); err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "FUNCIONARIO_NAO_ENCONTRADO", "Funcionário não encontrado", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao atualizar funcionário", "erro", err)
		web.RespondError(w, r, "ERRO_ATUALIZAR_FUNCIONARIO", "Erro ao atualizar funcionário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleBuscarFuncionario(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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
