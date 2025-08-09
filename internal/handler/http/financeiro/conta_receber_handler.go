package financeiro

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/service/financeiro/dto"
)

// ContaReceberService define a interface para o service de conta a receber
type ContaReceberService interface {
	CriarConta(ctx context.Context, input dto.CriarContaReceberInput) (*dto.ContaReceberOutput, error)
	RegistrarRecebimento(ctx context.Context, contaID string, input dto.RegistrarRecebimentoContaInput) (*dto.ContaReceberOutput, error)
	BuscarPorID(ctx context.Context, id string) (*dto.ContaReceberOutput, error)
	ListarPorObraID(ctx context.Context, obraID string) ([]*dto.ContaReceberOutput, error)
	ListarVencidas(ctx context.Context) ([]*dto.ContaReceberOutput, error)
	Listar(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*dto.ContaReceberOutput], error)
	ObterResumo(ctx context.Context, filtros dto.FiltrosContaReceberInput) (*dto.ResumoContasReceberOutput, error)
}

// ContaReceberHandler gerencia as rotas de contas a receber
type ContaReceberHandler struct {
	service ContaReceberService
	logger  *slog.Logger
}

func NovoContaReceberHandler(service ContaReceberService, logger *slog.Logger) *ContaReceberHandler {
	return &ContaReceberHandler{
		service: service,
		logger:  logger.With("handler", "conta_receber"),
	}
}

// HandleCriarConta cria uma nova conta a receber
func (h *ContaReceberHandler) HandleCriarConta(w http.ResponseWriter, r *http.Request) {
	var input dto.CriarContaReceberInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	conta, err := h.service.CriarConta(r.Context(), input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao criar conta a receber", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao criar conta a receber", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, conta, http.StatusCreated)
}

// HandleRegistrarRecebimento registra um recebimento
func (h *ContaReceberHandler) HandleRegistrarRecebimento(w http.ResponseWriter, r *http.Request) {
	contaID := chi.URLParam(r, "contaId")
	if contaID == "" {
		web.RespondError(w, r, "PARAMETRO_OBRIGATORIO", "contaId é obrigatório", http.StatusBadRequest)
		return
	}

	var input dto.RegistrarRecebimentoContaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	conta, err := h.service.RegistrarRecebimento(r.Context(), contaID, input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao registrar recebimento", 
			"conta_id", contaID, "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao registrar recebimento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, conta, http.StatusOK)
}

// HandleBuscarConta busca uma conta por ID
func (h *ContaReceberHandler) HandleBuscarConta(w http.ResponseWriter, r *http.Request) {
	contaID := chi.URLParam(r, "contaId")
	if contaID == "" {
		web.RespondError(w, r, "PARAMETRO_OBRIGATORIO", "contaId é obrigatório", http.StatusBadRequest)
		return
	}

	conta, err := h.service.BuscarPorID(r.Context(), contaID)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao buscar conta", "conta_id", contaID, "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Conta não encontrada", http.StatusNotFound)
		return
	}

	web.Respond(w, r, conta, http.StatusOK)
}

// HandleListarContasPorObra lista contas de uma obra
func (h *ContaReceberHandler) HandleListarContasPorObra(w http.ResponseWriter, r *http.Request) {
	obraID := chi.URLParam(r, "obraId")
	if obraID == "" {
		web.RespondError(w, r, "PARAMETRO_OBRIGATORIO", "obraId é obrigatório", http.StatusBadRequest)
		return
	}

	contas, err := h.service.ListarPorObraID(r.Context(), obraID)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar contas por obra", "obra_id", obraID, "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar contas", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, contas, http.StatusOK)
}

// HandleListarContasVencidas lista contas vencidas
func (h *ContaReceberHandler) HandleListarContasVencidas(w http.ResponseWriter, r *http.Request) {
	contas, err := h.service.ListarVencidas(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar contas vencidas", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar contas vencidas", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, contas, http.StatusOK)
}

// HandleListarContas lista contas com filtros e paginação
func (h *ContaReceberHandler) HandleListarContas(w http.ResponseWriter, r *http.Request) {
	filtros := web.ParseFiltros(r)
	
	contas, err := h.service.Listar(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar contas", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar contas", http.StatusInternalServerError)
		return
	}
	
	web.Respond(w, r, contas, http.StatusOK)
}

// HandleObterResumo obtém resumo das contas a receber
func (h *ContaReceberHandler) HandleObterResumo(w http.ResponseWriter, r *http.Request) {
	// Por enquanto, filtros vazios - pode ser expandido para aceitar query params
	filtros := dto.FiltrosContaReceberInput{}
	
	resumo, err := h.service.ObterResumo(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao obter resumo de contas", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter resumo", http.StatusInternalServerError)
		return
	}
	
	web.Respond(w, r, resumo, http.StatusOK)
}

// SetupContaReceberRoutes configura as rotas de contas a receber
func SetupContaReceberRoutes(r chi.Router, handler *ContaReceberHandler) {
	r.Route("/contas-receber", func(r chi.Router) {
		// CRUD básico
		r.Post("/", handler.HandleCriarConta)                    // Criar conta
		r.Get("/", handler.HandleListarContas)                   // Listar com paginação e filtros
		r.Get("/{contaId}", handler.HandleBuscarConta)           // Buscar por ID
		
		// Ações específicas
		r.Post("/{contaId}/recebimentos", handler.HandleRegistrarRecebimento) // Registrar recebimento
		
		// Relatórios e consultas
		r.Get("/vencidas", handler.HandleListarContasVencidas)   // Listar vencidas
		r.Get("/resumo", handler.HandleObterResumo)              // Resumo geral
	})

	// Listar contas de uma obra específica
	r.Get("/obras/{obraId}/contas-receber", handler.HandleListarContasPorObra)
}