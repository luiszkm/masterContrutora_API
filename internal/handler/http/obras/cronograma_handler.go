package obras

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/service/obras/dto"
)

// CronogramaService define a interface para o service de cronograma
type CronogramaService interface {
	CriarCronograma(ctx context.Context, input dto.CriarCronogramaRecebimentoInput) (*dto.CronogramaRecebimentoOutput, error)
	CriarCronogramaEmLote(ctx context.Context, input dto.CriarCronogramaEmLoteInput) ([]*dto.CronogramaRecebimentoOutput, error)
	RegistrarRecebimento(ctx context.Context, cronogramaID string, input dto.RegistrarRecebimentoInput) (*dto.CronogramaRecebimentoOutput, error)
	ListarPorObraID(ctx context.Context, obraID string) ([]*dto.CronogramaRecebimentoOutput, error)
	BuscarPorID(ctx context.Context, id string) (*dto.CronogramaRecebimentoOutput, error)
}

// CronogramaHandler gerencia as rotas de cronograma de recebimento
type CronogramaHandler struct {
	service CronogramaService
	logger  *slog.Logger
}

func NovoCronogramaHandler(service CronogramaService, logger *slog.Logger) *CronogramaHandler {
	return &CronogramaHandler{
		service: service,
		logger:  logger.With("handler", "cronograma"),
	}
}

// HandleCriarCronograma cria um cronograma de recebimento
func (h *CronogramaHandler) HandleCriarCronograma(w http.ResponseWriter, r *http.Request) {
	var input dto.CriarCronogramaRecebimentoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	cronograma, err := h.service.CriarCronograma(r.Context(), input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao criar cronograma", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao criar cronograma", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, cronograma, http.StatusCreated)
}

// HandleCriarCronogramaEmLote cria múltiplos cronogramas
func (h *CronogramaHandler) HandleCriarCronogramaEmLote(w http.ResponseWriter, r *http.Request) {
	var input dto.CriarCronogramaEmLoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	cronogramas, err := h.service.CriarCronogramaEmLote(r.Context(), input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao criar cronogramas em lote", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao criar cronogramas", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, cronogramas, http.StatusCreated)
}

// HandleRegistrarRecebimento registra um recebimento
func (h *CronogramaHandler) HandleRegistrarRecebimento(w http.ResponseWriter, r *http.Request) {
	cronogramaID := chi.URLParam(r, "cronogramaId")
	if cronogramaID == "" {
		web.RespondError(w, r, "PARAMETRO_OBRIGATORIO", "cronogramaId é obrigatório", http.StatusBadRequest)
		return
	}

	var input dto.RegistrarRecebimentoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	cronograma, err := h.service.RegistrarRecebimento(r.Context(), cronogramaID, input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao registrar recebimento", 
			"cronograma_id", cronogramaID, "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao registrar recebimento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, cronograma, http.StatusOK)
}

// HandleListarCronogramasPorObra lista cronogramas de uma obra
func (h *CronogramaHandler) HandleListarCronogramasPorObra(w http.ResponseWriter, r *http.Request) {
	obraID := chi.URLParam(r, "obraId")
	if obraID == "" {
		web.RespondError(w, r, "PARAMETRO_OBRIGATORIO", "obraId é obrigatório", http.StatusBadRequest)
		return
	}

	cronogramas, err := h.service.ListarPorObraID(r.Context(), obraID)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar cronogramas", "obra_id", obraID, "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar cronogramas", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, cronogramas, http.StatusOK)
}

// HandleBuscarCronograma busca um cronograma por ID
func (h *CronogramaHandler) HandleBuscarCronograma(w http.ResponseWriter, r *http.Request) {
	cronogramaID := chi.URLParam(r, "cronogramaId")
	if cronogramaID == "" {
		web.RespondError(w, r, "PARAMETRO_OBRIGATORIO", "cronogramaId é obrigatório", http.StatusBadRequest)
		return
	}

	cronograma, err := h.service.BuscarPorID(r.Context(), cronogramaID)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao buscar cronograma", "cronograma_id", cronogramaID, "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Cronograma não encontrado", http.StatusNotFound)
		return
	}

	web.Respond(w, r, cronograma, http.StatusOK)
}

// SetupCronogramaRoutes configura as rotas do cronograma
func SetupCronogramaRoutes(r chi.Router, handler *CronogramaHandler) {
	r.Route("/cronograma-recebimentos", func(r chi.Router) {
		// Criar cronograma individual
		r.Post("/", handler.HandleCriarCronograma)
		
		// Criar cronogramas em lote  
		r.Post("/lote", handler.HandleCriarCronogramaEmLote)
		
		// Buscar cronograma específico
		r.Get("/{cronogramaId}", handler.HandleBuscarCronograma)
		
		// Registrar recebimento
		r.Post("/{cronogramaId}/recebimentos", handler.HandleRegistrarRecebimento)
	})

	// Listar cronogramas de uma obra específica
	r.Get("/obras/{obraId}/cronograma-recebimentos", handler.HandleListarCronogramasPorObra)
}