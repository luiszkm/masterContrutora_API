// file: internal/handler/http/obras.go
package obras

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
)

// ObrasService define a interface que o handler espera do serviço.
// Isso permite testar o handler com um serviço mock.
type Service interface {
	CriarNovaObra(ctx context.Context, input dto.CriarNovaObraInput) (*obras.Obra, error)
	BuscarDashboard(ctx context.Context, id string) (*dto.ObraDashboard, error)
	AdicionarEtapa(ctx context.Context, obraID string, input dto.AdicionarEtapaInput) (*obras.Etapa, error)
	AtualizarStatusEtapa(ctx context.Context, etapaID string, input dto.AtualizarStatusEtapaInput) (*obras.Etapa, error)
	AlocarFuncionario(ctx context.Context, obraID string, input dto.AlocarFuncionarioInput) (*obras.Alocacao, error)
	ListarObras(ctx context.Context) ([]*dto.ObraListItemDTO, error)
}

// ObrasHandler gerencia as requisições HTTP para o contexto de Obras.
type Handler struct {
	service Service
	logger  *slog.Logger
}

// NovoObrasHandler cria um novo handler para obras.
func NovoObrasHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{
		service: s,
		logger:  l,
	}
}

// HandleCriarObra trata a criação de uma nova obra.
func (h *Handler) HandleCriarObra(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var input dto.CriarNovaObraInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("falha ao decodificar payload", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Payload inválido", http.StatusBadRequest)
		return
	}

	obra, err := h.service.CriarNovaObra(r.Context(), input)
	if err != nil {
		// Aqui poderíamos ter uma lógica mais granular para mapear
		// erros de serviço para status HTTP (ex: 400, 409, etc).
		h.logger.ErrorContext(r.Context(), "falha ao criar obra", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro interno ao processar sua requisição", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(obra); err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao encodificar resposta", "erro", err)
	}
}

func (h *Handler) HandleAdicionarEtapa(w http.ResponseWriter, r *http.Request) {
	obraID := chi.URLParam(r, "obraId")
	if obraID == "" {
		http.Error(w, "ID da obra na URL é obrigatório", http.StatusBadRequest)
		return
	}

	var input dto.AdicionarEtapaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("falha ao decodificar payload da etapa", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Payload inválido", http.StatusBadRequest)
		return
	}

	etapa, err := h.service.AdicionarEtapa(r.Context(), obraID, input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao adicionar etapa", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro interno ao processar sua requisição", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, etapa, http.StatusCreated)
}

func (h *Handler) HandleBuscarObra(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "obraId")

	id, err := uuid.Parse(idStr)
	if err != nil {
		web.RespondError(w, r, "ID_INVALIDO", "O ID da obra fornecido na URL não é um UUID válido", http.StatusBadRequest)
		return
	}

	dashboard, err := h.service.BuscarDashboard(r.Context(), id.String())
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "OBRA_NAO_ENCONTRADA", "Obra não encontrada", http.StatusNotFound)
			return
		}

		h.logger.ErrorContext(r.Context(), "falha ao buscar dashboard da obra", "erro", err, "obra_id", id)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro interno ao buscar a obra", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, dashboard, http.StatusOK)
}

func (h *Handler) HandleAtualizarEtapaStatus(w http.ResponseWriter, r *http.Request) {
	etapaID := chi.URLParam(r, "etapaId")

	var input dto.AtualizarStatusEtapaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	etapa, err := h.service.AtualizarStatusEtapa(r.Context(), etapaID, input)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			web.RespondError(w, r, "ETAPA_NAO_ENCONTRADA", "Etapa não encontrada", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao atualizar etapa", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro interno ao processar sua requisição", http.StatusInternalServerError)
		return
	}
	web.Respond(w, r, etapa, http.StatusOK)
}

func (h *Handler) HandleAlocarFuncionario(w http.ResponseWriter, r *http.Request) {
	obraID := chi.URLParam(r, "obraId")

	var input dto.AlocarFuncionarioInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	alocacao, err := h.service.AlocarFuncionario(r.Context(), obraID, input)
	if err != nil {
		// TODO: Tratar erros específicos (obra ou funcionário não encontrado) com 404
		h.logger.ErrorContext(r.Context(), "falha ao alocar funcionário", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao alocar funcionário", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, alocacao, http.StatusCreated)
}

func (h *Handler) HandleListarObras(w http.ResponseWriter, r *http.Request) {
	obras, err := h.service.ListarObras(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar obras", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar obras", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, obras, http.StatusOK)
}
