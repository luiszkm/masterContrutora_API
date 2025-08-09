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

// ContaPagarService define a interface para o service de conta a pagar
type ContaPagarService interface {
	CriarConta(ctx context.Context, input dto.CriarContaPagarInput) (*dto.ContaPagarOutput, error)
	CriarContaDeOrcamento(ctx context.Context, input dto.CriarContaPagarDeOrcamentoInput, orcamento interface{}) (*dto.ContaPagarOutput, error)
	RegistrarPagamento(ctx context.Context, contaID string, input dto.RegistrarPagamentoContaPagarInput) (*dto.ContaPagarOutput, error)
	BuscarPorID(ctx context.Context, id string) (*dto.ContaPagarOutput, error)
	ListarPorObraID(ctx context.Context, obraID string) ([]*dto.ContaPagarOutput, error)
	ListarPorFornecedorID(ctx context.Context, fornecedorID string) ([]*dto.ContaPagarOutput, error)
	ListarVencidas(ctx context.Context) ([]*dto.ContaPagarOutput, error)
	Listar(ctx context.Context, filtros common.ListarFiltros) (*common.RespostaPaginada[*dto.ContaPagarOutput], error)
	ObterResumo(ctx context.Context, filtros dto.FiltrosContaPagarInput) (*dto.ResumoContasPagarOutput, error)
}

// ContaPagarHandler gerencia as rotas de contas a pagar
type ContaPagarHandler struct {
	service ContaPagarService
	logger  *slog.Logger
}

func NovoContaPagarHandler(service ContaPagarService, logger *slog.Logger) *ContaPagarHandler {
	return &ContaPagarHandler{
		service: service,
		logger:  logger.With("handler", "conta_pagar"),
	}
}

// HandleCriarConta cria uma nova conta a pagar
func (h *ContaPagarHandler) HandleCriarConta(w http.ResponseWriter, r *http.Request) {
	var input dto.CriarContaPagarInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	conta, err := h.service.CriarConta(r.Context(), input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao criar conta a pagar", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao criar conta a pagar", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, conta, http.StatusCreated)
}

// HandleCriarContaDeOrcamento cria uma conta a pagar a partir de um orçamento
func (h *ContaPagarHandler) HandleCriarContaDeOrcamento(w http.ResponseWriter, r *http.Request) {
	var input dto.CriarContaPagarDeOrcamentoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	conta, err := h.service.CriarContaDeOrcamento(r.Context(), input, nil) // TODO: buscar orçamento
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao criar conta de orçamento", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao criar conta de orçamento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, conta, http.StatusCreated)
}

// HandleRegistrarPagamento registra um pagamento
func (h *ContaPagarHandler) HandleRegistrarPagamento(w http.ResponseWriter, r *http.Request) {
	contaID := chi.URLParam(r, "contaId")
	if contaID == "" {
		web.RespondError(w, r, "PARAMETRO_OBRIGATORIO", "contaId é obrigatório", http.StatusBadRequest)
		return
	}

	var input dto.RegistrarPagamentoContaPagarInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
		return
	}

	conta, err := h.service.RegistrarPagamento(r.Context(), contaID, input)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao registrar pagamento", 
			"conta_id", contaID, "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao registrar pagamento", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, conta, http.StatusOK)
}

// HandleBuscarConta busca uma conta por ID
func (h *ContaPagarHandler) HandleBuscarConta(w http.ResponseWriter, r *http.Request) {
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
func (h *ContaPagarHandler) HandleListarContasPorObra(w http.ResponseWriter, r *http.Request) {
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

// HandleListarContasPorFornecedor lista contas de um fornecedor
func (h *ContaPagarHandler) HandleListarContasPorFornecedor(w http.ResponseWriter, r *http.Request) {
	fornecedorID := chi.URLParam(r, "fornecedorId")
	if fornecedorID == "" {
		web.RespondError(w, r, "PARAMETRO_OBRIGATORIO", "fornecedorId é obrigatório", http.StatusBadRequest)
		return
	}

	contas, err := h.service.ListarPorFornecedorID(r.Context(), fornecedorID)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar contas por fornecedor", 
			"fornecedor_id", fornecedorID, "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar contas", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, contas, http.StatusOK)
}

// HandleListarContasVencidas lista contas vencidas
func (h *ContaPagarHandler) HandleListarContasVencidas(w http.ResponseWriter, r *http.Request) {
	contas, err := h.service.ListarVencidas(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar contas vencidas", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar contas vencidas", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, contas, http.StatusOK)
}

// HandleListarContas lista contas com filtros e paginação
func (h *ContaPagarHandler) HandleListarContas(w http.ResponseWriter, r *http.Request) {
	filtros := web.ParseFiltros(r)
	
	contas, err := h.service.Listar(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao listar contas", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao listar contas", http.StatusInternalServerError)
		return
	}
	
	web.Respond(w, r, contas, http.StatusOK)
}

// HandleObterResumo obtém resumo das contas a pagar
func (h *ContaPagarHandler) HandleObterResumo(w http.ResponseWriter, r *http.Request) {
	// Por enquanto, filtros vazios - pode ser expandido para aceitar query params
	filtros := dto.FiltrosContaPagarInput{}
	
	resumo, err := h.service.ObterResumo(r.Context(), filtros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao obter resumo de contas", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter resumo", http.StatusInternalServerError)
		return
	}
	
	web.Respond(w, r, resumo, http.StatusOK)
}

// SetupContaPagarRoutes configura as rotas de contas a pagar
func SetupContaPagarRoutes(r chi.Router, handler *ContaPagarHandler) {
	r.Route("/contas-pagar", func(r chi.Router) {
		// CRUD básico
		r.Post("/", handler.HandleCriarConta)                    // Criar conta
		r.Get("/", handler.HandleListarContas)                   // Listar com paginação e filtros
		r.Get("/{contaId}", handler.HandleBuscarConta)           // Buscar por ID
		
		// Ações específicas
		r.Post("/{contaId}/pagamentos", handler.HandleRegistrarPagamento) // Registrar pagamento
		r.Post("/orcamentos", handler.HandleCriarContaDeOrcamento)        // Criar de orçamento
		
		// Relatórios e consultas
		r.Get("/vencidas", handler.HandleListarContasVencidas)   // Listar vencidas
		r.Get("/resumo", handler.HandleObterResumo)              // Resumo geral
	})

	// Rotas específicas por entidade relacionada
	r.Get("/obras/{obraId}/contas-pagar", handler.HandleListarContasPorObra)
	r.Get("/fornecedores/{fornecedorId}/contas-pagar", handler.HandleListarContasPorFornecedor)
}