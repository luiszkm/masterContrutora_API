// file: internal/handler/http/obras.go
package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/luiszkm/masterCostrutora/internal/domain/obras"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
)

// ObrasService define a interface que o handler espera do serviço.
// Isso permite testar o handler com um serviço mock.
type ObrasService interface {
	CriarNovaObra(ctx context.Context, input dto.CriarNovaObraInput) (*obras.Obra, error)
	BuscarDashboard(ctx context.Context, id string) (*dto.ObraDashboard, error)
	AdicionarEtapa(ctx context.Context, obraID string, input dto.AdicionarEtapaInput) (*obras.Etapa, error)
	AtualizarStatusEtapa(ctx context.Context, etapaID string, input dto.AtualizarStatusEtapaInput) (*obras.Etapa, error)
}

// ObrasHandler gerencia as requisições HTTP para o contexto de Obras.
type ObrasHandler struct {
	service ObrasService
	logger  *slog.Logger
}

// NovoObrasHandler cria um novo handler para obras.
func NovoObrasHandler(s ObrasService, l *slog.Logger) *ObrasHandler {
	return &ObrasHandler{
		service: s,
		logger:  l,
	}
}

// HandleCriarObra trata a criação de uma nova obra.
func (h *ObrasHandler) HandleCriarObra(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var input dto.CriarNovaObraInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("falha ao decodificar payload", "erro", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	obra, err := h.service.CriarNovaObra(r.Context(), input)
	if err != nil {
		// Aqui poderíamos ter uma lógica mais granular para mapear
		// erros de serviço para status HTTP (ex: 400, 409, etc).
		h.logger.ErrorContext(r.Context(), "falha ao criar obra", "erro", err)
		http.Error(w, "Erro interno ao processar sua requisição", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(obra); err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao encodificar resposta", "erro", err)
	}
}

func (h *ObrasHandler) HandleAdicionarEtapa(w http.ResponseWriter, r *http.Request) {
	obraID := chi.URLParam(r, "obraId")
	if obraID == "" {
		http.Error(w, "ID da obra na URL é obrigatório", http.StatusBadRequest)
		return
	}

	var input dto.AdicionarEtapaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("falha ao decodificar payload da etapa", "erro", err)
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	etapa, err := h.service.AdicionarEtapa(r.Context(), obraID, input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao adicionar etapa", "erro", err)
		http.Error(w, "Erro interno ao processar sua requisição", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(etapa)
}

func (h *ObrasHandler) HandleBuscarObra(w http.ResponseWriter, r *http.Request) {
	// A forma correta e segura de pegar o parâmetro da URL com chi
	id := chi.URLParam(r, "obraId")
	if id == "" {
		http.Error(w, "ID da obra é obrigatório", http.StatusBadRequest)
		return
	}

	dashboard, err := h.service.BuscarDashboard(r.Context(), id)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao buscar dashboard da obra", "erro", err, "obra_id", id)
		http.Error(w, "Erro interno ao buscar a obra", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dashboard)
}

func (h *ObrasHandler) HandleAtualizarEtapaStatus(w http.ResponseWriter, r *http.Request) {
	etapaID := chi.URLParam(r, "etapaId")

	var input dto.AtualizarStatusEtapaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	etapa, err := h.service.AtualizarStatusEtapa(r.Context(), etapaID, input)
	if err != nil {
		if errors.Is(err, postgres.ErrNaoEncontrado) {
			http.Error(w, "Etapa não encontrada", http.StatusNotFound)
			return
		}
		h.logger.ErrorContext(r.Context(), "falha ao atualizar etapa", "erro", err)
		http.Error(w, "Erro interno ao processar sua requisição", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(etapa)
}
