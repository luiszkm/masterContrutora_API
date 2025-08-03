package dashboard

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/service/dashboard/dto"
	"github.com/luiszkm/masterCostrutora/pkg/auth"
	"github.com/luiszkm/masterCostrutora/pkg/logging"
)

// Service define a interface que o handler espera do serviço de dashboard
type Service interface {
	ObterDashboardCompleto(ctx context.Context, parametros dto.ParametrosDashboardDTO) (*dto.DashboardGeralDTO, error)
	ObterDashboardFinanceiro(ctx context.Context, dataInicio, dataFim time.Time) (*dto.DashboardFinanceiroDTO, error)
	ObterDashboardObras(ctx context.Context) (*dto.DashboardObrasDTO, error)
	ObterDashboardFuncionarios(ctx context.Context, dataInicio, dataFim time.Time) (*dto.DashboardFuncionariosDTO, error)
	ObterDashboardFornecedores(ctx context.Context, dataInicio, dataFim time.Time) (*dto.DashboardFornecedoresDTO, error)
	ObterFluxoCaixa(ctx context.Context, dataInicio, dataFim time.Time) (*dto.FluxoCaixaResumoDTO, error)
}

// Handler gerencia as requisições HTTP para o dashboard
type Handler struct {
	service    Service
	logger     *slog.Logger
	dashLogger *logging.DashboardLogger
	jwtService *auth.JWTService
}

// NovoDashboardHandler cria um novo handler para o dashboard
func NovoDashboardHandler(s Service, l *slog.Logger, dashLogger *logging.DashboardLogger, jwtService *auth.JWTService) *Handler {
	return &Handler{
		service:    s,
		logger:     l,
		dashLogger: dashLogger,
		jwtService: jwtService,
	}
}

// HandleObterDashboardCompleto retorna o dashboard completo
func (h *Handler) HandleObterDashboardCompleto(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	userID := h.extractUserID(r)
	
	// Log acesso
	h.dashLogger.LogDashboardAuth(r.Context(), "geral", userID, "dashboard_read", true, map[string]interface{}{
		"endpoint": r.URL.Path,
		"method": r.Method,
	})
	
	parametros := h.parseParametrosDashboard(r)
	
	// Log parâmetros recebidos
	h.dashLogger.LogDashboardData(r.Context(), "geral", "parametros_request", 1, false, map[string]interface{}{
		"parametros": parametros,
		"userAgent": r.UserAgent(),
		"remoteAddr": r.RemoteAddr,
	})
	
	dashboard, err := h.service.ObterDashboardCompleto(r.Context(), parametros)
	if err != nil {
		// Log erro detalhado
		h.dashLogger.LogDashboardError(r.Context(), "geral", "HandleObterDashboardCompleto", err, map[string]interface{}{
			"parametros": parametros,
			"userID": userID,
			"duration": time.Since(startTime).String(),
		})
		
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard completo", "erro", err, "parametros", parametros)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dashboard", http.StatusInternalServerError)
		return
	}

	// Log sucesso
	h.dashLogger.LogDashboardData(r.Context(), "geral", "dashboard_response", 1, dashboard == nil, map[string]interface{}{
		"duration": time.Since(startTime).String(),
		"responseSize": "large", // Poderia calcular tamanho real
	})

	web.Respond(w, r, dashboard, http.StatusOK)
}

// HandleObterDashboardFinanceiro retorna apenas a seção financeira
func (h *Handler) HandleObterDashboardFinanceiro(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	userID := h.extractUserID(r)
	
	// Log acesso
	h.dashLogger.LogDashboardAuth(r.Context(), "financeiro", userID, "dashboard_read", true, map[string]interface{}{
		"endpoint": r.URL.Path,
		"method": r.Method,
	})
	
	dataInicio, dataFim := h.parsePeriodo(r)
	
	financeiro, err := h.service.ObterDashboardFinanceiro(r.Context(), dataInicio, dataFim)
	if err != nil {
		h.dashLogger.LogDashboardError(r.Context(), "financeiro", "HandleObterDashboardFinanceiro", err, map[string]interface{}{
			"dataInicio": dataInicio,
			"dataFim": dataFim,
			"userID": userID,
			"duration": time.Since(startTime).String(),
		})
		
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard financeiro", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dados financeiros", http.StatusInternalServerError)
		return
	}

	// Log sucesso
	h.dashLogger.LogDashboardData(r.Context(), "financeiro", "dashboard_response", 1, financeiro == nil, map[string]interface{}{
		"duration": time.Since(startTime).String(),
	})

	web.Respond(w, r, financeiro, http.StatusOK)
}

// HandleObterDashboardObras retorna apenas a seção de obras
func (h *Handler) HandleObterDashboardObras(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	userID := h.extractUserID(r)
	
	// Log acesso
	h.dashLogger.LogDashboardAuth(r.Context(), "obras", userID, "dashboard_read", true, map[string]interface{}{
		"endpoint": r.URL.Path,
		"method": r.Method,
	})
	
	obras, err := h.service.ObterDashboardObras(r.Context())
	if err != nil {
		h.dashLogger.LogDashboardError(r.Context(), "obras", "HandleObterDashboardObras", err, map[string]interface{}{
			"userID": userID,
			"duration": time.Since(startTime).String(),
		})
		
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard obras", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dados de obras", http.StatusInternalServerError)
		return
	}

	// Log sucesso
	h.dashLogger.LogDashboardData(r.Context(), "obras", "dashboard_response", 1, obras == nil, map[string]interface{}{
		"duration": time.Since(startTime).String(),
	})

	web.Respond(w, r, obras, http.StatusOK)
}

// HandleObterDashboardFuncionarios retorna apenas a seção de funcionários
func (h *Handler) HandleObterDashboardFuncionarios(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	userID := h.extractUserID(r)
	
	// Log acesso
	h.dashLogger.LogDashboardAuth(r.Context(), "funcionarios", userID, "dashboard_read", true, map[string]interface{}{
		"endpoint": r.URL.Path,
		"method": r.Method,
	})
	
	dataInicio, dataFim := h.parsePeriodo(r)
	
	funcionarios, err := h.service.ObterDashboardFuncionarios(r.Context(), dataInicio, dataFim)
	if err != nil {
		h.dashLogger.LogDashboardError(r.Context(), "funcionarios", "HandleObterDashboardFuncionarios", err, map[string]interface{}{
			"dataInicio": dataInicio,
			"dataFim": dataFim,
			"userID": userID,
			"duration": time.Since(startTime).String(),
		})
		
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard funcionários", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dados de funcionários", http.StatusInternalServerError)
		return
	}

	// Log sucesso
	h.dashLogger.LogDashboardData(r.Context(), "funcionarios", "dashboard_response", 1, funcionarios == nil, map[string]interface{}{
		"duration": time.Since(startTime).String(),
	})

	web.Respond(w, r, funcionarios, http.StatusOK)
}

// HandleObterDashboardFornecedores retorna apenas a seção de fornecedores
func (h *Handler) HandleObterDashboardFornecedores(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	userID := h.extractUserID(r)
	
	// Log acesso
	h.dashLogger.LogDashboardAuth(r.Context(), "fornecedores", userID, "dashboard_read", true, map[string]interface{}{
		"endpoint": r.URL.Path,
		"method": r.Method,
	})
	
	dataInicio, dataFim := h.parsePeriodo(r)
	
	fornecedores, err := h.service.ObterDashboardFornecedores(r.Context(), dataInicio, dataFim)
	if err != nil {
		h.dashLogger.LogDashboardError(r.Context(), "fornecedores", "HandleObterDashboardFornecedores", err, map[string]interface{}{
			"dataInicio": dataInicio,
			"dataFim": dataFim,
			"userID": userID,
			"duration": time.Since(startTime).String(),
		})
		
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard fornecedores", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dados de fornecedores", http.StatusInternalServerError)
		return
	}

	// Log sucesso
	h.dashLogger.LogDashboardData(r.Context(), "fornecedores", "dashboard_response", 1, fornecedores == nil, map[string]interface{}{
		"duration": time.Since(startTime).String(),
	})

	web.Respond(w, r, fornecedores, http.StatusOK)
}

// HandleObterFluxoCaixa retorna apenas dados de fluxo de caixa
func (h *Handler) HandleObterFluxoCaixa(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	userID := h.extractUserID(r)
	
	// Log acesso
	h.dashLogger.LogDashboardAuth(r.Context(), "fluxo-caixa", userID, "dashboard_read", true, map[string]interface{}{
		"endpoint": r.URL.Path,
		"method": r.Method,
	})
	
	dataInicio, dataFim := h.parsePeriodo(r)
	
	fluxo, err := h.service.ObterFluxoCaixa(r.Context(), dataInicio, dataFim)
	if err != nil {
		h.dashLogger.LogDashboardError(r.Context(), "fluxo-caixa", "HandleObterFluxoCaixa", err, map[string]interface{}{
			"dataInicio": dataInicio,
			"dataFim": dataFim,
			"userID": userID,
			"duration": time.Since(startTime).String(),
		})
		
		h.logger.ErrorContext(r.Context(), "falha ao obter fluxo de caixa", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter fluxo de caixa", http.StatusInternalServerError)
		return
	}

	// Log sucesso
	h.dashLogger.LogDashboardData(r.Context(), "fluxo-caixa", "dashboard_response", 1, fluxo == nil, map[string]interface{}{
		"duration": time.Since(startTime).String(),
	})

	web.Respond(w, r, fluxo, http.StatusOK)
}

// parseParametrosDashboard extrai os parâmetros da query string
func (h *Handler) parseParametrosDashboard(r *http.Request) dto.ParametrosDashboardDTO {
	parametros := dto.ParametrosDashboardDTO{}
	
	// Parse período
	dataInicio, dataFim := h.parsePeriodo(r)
	parametros.DataInicio = &dataInicio
	parametros.DataFim = &dataFim
	
	// Parse seções específicas
	if secoes := r.URL.Query()["secoes"]; len(secoes) > 0 {
		parametros.Secoes = secoes
	}
	
	// Parse obras específicas
	if obraIDs := r.URL.Query()["obraIds"]; len(obraIDs) > 0 {
		parametros.ObraIDs = obraIDs
	}
	
	// Parse fornecedores específicos
	if fornecedorIDs := r.URL.Query()["fornecedorIds"]; len(fornecedorIDs) > 0 {
		parametros.FornecedorIDs = fornecedorIDs
	}
	
	// Parse incluir inativos
	if incluirInativosStr := r.URL.Query().Get("incluirInativos"); incluirInativosStr != "" {
		if incluirInativos, err := strconv.ParseBool(incluirInativosStr); err == nil {
			parametros.IncluirInativos = incluirInativos
		}
	}
	
	return parametros
}

// parsePeriodo extrai as datas de início e fim dos parâmetros da query
func (h *Handler) parsePeriodo(r *http.Request) (time.Time, time.Time) {
	// Período padrão: últimos 6 meses
	dataFim := time.Now()
	dataInicio := dataFim.AddDate(0, -6, 0)
	
	// Parse data de início
	if dataInicioStr := r.URL.Query().Get("dataInicio"); dataInicioStr != "" {
		if parsed, err := time.Parse("2006-01-02", dataInicioStr); err == nil {
			dataInicio = parsed
		}
	}
	
	// Parse data de fim
	if dataFimStr := r.URL.Query().Get("dataFim"); dataFimStr != "" {
		if parsed, err := time.Parse("2006-01-02", dataFimStr); err == nil {
			dataFim = parsed
		}
	}
	
	return dataInicio, dataFim
}

// HandleObterDashboardPorSecao permite obter dados de uma seção específica via URL
func (h *Handler) HandleObterDashboardPorSecao(w http.ResponseWriter, r *http.Request) {
	secao := chi.URLParam(r, "secao")
	
	switch secao {
	case "financeiro":
		h.HandleObterDashboardFinanceiro(w, r)
	case "obras":
		h.HandleObterDashboardObras(w, r)
	case "funcionarios":
		h.HandleObterDashboardFuncionarios(w, r)
	case "fornecedores":
		h.HandleObterDashboardFornecedores(w, r)
	case "fluxo-caixa":
		h.HandleObterFluxoCaixa(w, r)
	default:
		web.RespondError(w, r, "SECAO_INVALIDA", "Seção de dashboard não encontrada", http.StatusNotFound)
		return
	}
}

// HandleObterParametrosCache retorna parâmetros para controle de cache do frontend
func (h *Handler) HandleObterParametrosCache(w http.ResponseWriter, r *http.Request) {
	// Retorna informações que ajudam o frontend a gerenciar cache
	cacheInfo := map[string]interface{}{
		"ultimaAtualizacao": time.Now(),
		"ttlRecomendado":    300, // 5 minutos em segundos
		"secoesDisponiveis": []string{"financeiro", "obras", "funcionarios", "fornecedores"},
		"versao":           "1.0",
	}
	
	web.Respond(w, r, cacheInfo, http.StatusOK)
}

// extractUserID extrai o ID do usuário do contexto ou do JWT
func (h *Handler) extractUserID(r *http.Request) string {
	// Primeiro, tentar obter do contexto (se foi definido pelo middleware de autenticação)
	if userID := r.Context().Value(auth.UserContextKey); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	
	// Se não encontrou no contexto, tentar extrair do JWT diretamente
	if h.jwtService != nil {
		return h.extractUserIDFromJWT(r)
	}
	
	return ""
}

// extractUserIDFromJWT extrai o ID do usuário do token JWT
func (h *Handler) extractUserIDFromJWT(r *http.Request) string {
	// Tentar extrair do header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Tentar extrair do cookie se não há header
		cookie, err := r.Cookie("token")
		if err != nil {
			return ""
		}
		authHeader = "Bearer " + cookie.Value
	}
	
	// Formato: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	
	token := parts[1]
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		return ""
	}
	
	// Extrair user ID do claims
	if sub, ok := claims["sub"].(string); ok {
		return sub
	}
	
	return ""
}