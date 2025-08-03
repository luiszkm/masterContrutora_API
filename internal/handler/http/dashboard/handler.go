package dashboard

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/internal/service/dashboard/dto"
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
	service Service
	logger  *slog.Logger
}

// NovoDashboardHandler cria um novo handler para o dashboard
func NovoDashboardHandler(s Service, l *slog.Logger) *Handler {
	return &Handler{
		service: s,
		logger:  l,
	}
}

// HandleObterDashboardCompleto retorna o dashboard completo
func (h *Handler) HandleObterDashboardCompleto(w http.ResponseWriter, r *http.Request) {
	parametros := h.parseParametrosDashboard(r)
	
	dashboard, err := h.service.ObterDashboardCompleto(r.Context(), parametros)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard completo", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dashboard", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, dashboard, http.StatusOK)
}

// HandleObterDashboardFinanceiro retorna apenas a seção financeira
func (h *Handler) HandleObterDashboardFinanceiro(w http.ResponseWriter, r *http.Request) {
	dataInicio, dataFim := h.parsePeriodo(r)
	
	financeiro, err := h.service.ObterDashboardFinanceiro(r.Context(), dataInicio, dataFim)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard financeiro", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dados financeiros", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, financeiro, http.StatusOK)
}

// HandleObterDashboardObras retorna apenas a seção de obras
func (h *Handler) HandleObterDashboardObras(w http.ResponseWriter, r *http.Request) {
	obras, err := h.service.ObterDashboardObras(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard obras", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dados de obras", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, obras, http.StatusOK)
}

// HandleObterDashboardFuncionarios retorna apenas a seção de funcionários
func (h *Handler) HandleObterDashboardFuncionarios(w http.ResponseWriter, r *http.Request) {
	dataInicio, dataFim := h.parsePeriodo(r)
	
	funcionarios, err := h.service.ObterDashboardFuncionarios(r.Context(), dataInicio, dataFim)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard funcionários", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dados de funcionários", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, funcionarios, http.StatusOK)
}

// HandleObterDashboardFornecedores retorna apenas a seção de fornecedores
func (h *Handler) HandleObterDashboardFornecedores(w http.ResponseWriter, r *http.Request) {
	dataInicio, dataFim := h.parsePeriodo(r)
	
	fornecedores, err := h.service.ObterDashboardFornecedores(r.Context(), dataInicio, dataFim)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao obter dashboard fornecedores", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter dados de fornecedores", http.StatusInternalServerError)
		return
	}

	web.Respond(w, r, fornecedores, http.StatusOK)
}

// HandleObterFluxoCaixa retorna apenas dados de fluxo de caixa
func (h *Handler) HandleObterFluxoCaixa(w http.ResponseWriter, r *http.Request) {
	dataInicio, dataFim := h.parsePeriodo(r)
	
	fluxo, err := h.service.ObterFluxoCaixa(r.Context(), dataInicio, dataFim)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "falha ao obter fluxo de caixa", "erro", err)
		web.RespondError(w, r, "ERRO_INTERNO", "Erro ao obter fluxo de caixa", http.StatusInternalServerError)
		return
	}

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