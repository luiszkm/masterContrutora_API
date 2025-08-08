// file: internal/handler/http/router/router.go
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/luiszkm/masterCostrutora/internal/authz"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/dashboard"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/financeiro"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/identidade"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/obras"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/suprimentos"
	"github.com/luiszkm/masterCostrutora/internal/handler/web"
	"github.com/luiszkm/masterCostrutora/pkg/auth"
)

type Config struct {
	JwtService         *auth.JWTService
	IdentidadeHandler  *identidade.Handler
	ObrasHandler       *obras.Handler
	PessoalHandler     *pessoal.Handler
	SuprimentosHandler *suprimentos.Handler
	FinanceiroHandler  *financeiro.Handler
	ContaReceberHandler *financeiro.ContaReceberHandler
	ContaPagarHandler  *financeiro.ContaPagarHandler
	DashboardHandler   *dashboard.Handler
}

func New(c Config) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares globais aplicados a todas as rotas
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- ROTAS PÚBLICAS ---
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		web.Respond(w, r, map[string]string{"status": "ok"}, http.StatusOK)
	})
	r.Route("/usuarios", func(r chi.Router) {
		r.Post("/registrar", c.IdentidadeHandler.HandleRegistrar)
		r.Post("/login", c.IdentidadeHandler.HandleLogin)
	})

	// --- DASHBOARD PÚBLICO PARA DEBUG ---
	r.Route("/dashboard", func(r chi.Router) {
		// Dashboard completo - sem autenticação para debug
		r.Get("/", c.DashboardHandler.HandleObterDashboardCompleto)
		
		// Seções específicas do dashboard - sem autenticação para debug
		r.Get("/financeiro", c.DashboardHandler.HandleObterDashboardFinanceiro)
		r.Get("/obras", c.DashboardHandler.HandleObterDashboardObras)
		r.Get("/funcionarios", c.DashboardHandler.HandleObterDashboardFuncionarios)
		r.Get("/fornecedores", c.DashboardHandler.HandleObterDashboardFornecedores)
		r.Get("/fluxo-caixa", c.DashboardHandler.HandleObterFluxoCaixa)
		r.Get("/{secao}", c.DashboardHandler.HandleObterDashboardPorSecao)
		r.Get("/cache-info", c.DashboardHandler.HandleObterParametrosCache)
	})

	// --- GRUPO ÚNICO PARA TODAS AS ROTAS PROTEGIDAS ---
	r.Group(func(r chi.Router) {
		// Aplicamos o middleware de autenticação UMA VEZ para todo o grupo.
		r.Use(c.JwtService.AuthMiddleware)

		// --- Recursos de Pessoal ---
		r.Route("/funcionarios", func(r chi.Router) {
			// Rotas que operam na coleção de funcionários: /funcionarios
			r.Use(c.JwtService.AuthMiddleware) // Aplica autenticação para todo o grupo

			r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).
				Post("/", c.PessoalHandler.HandleCadastrarFuncionario)

			r.With(auth.Authorize(authz.PermissaoPessoalLer)).
				Get("/", c.PessoalHandler.HandleListarFuncionarios)

			r.With(auth.Authorize(authz.PermissaoPessoalLer)).
				Get("/apontamentos", c.PessoalHandler.HandleListarComUltimoApontamento)

			r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).
				Put("/apontamentos/{apontamentoId}", c.PessoalHandler.HandleAtualizarApontamento)

			// r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).
			// 	Patch("/apontamentos/{apontamentoId}/pagar", c.PessoalHandler.HandleRegistrarPagamentoApontamento)

			r.With(auth.Authorize(authz.PermissaoPessoalApontamentoEscrever)).
				Post("/apontamentos/replicar", c.PessoalHandler.HandleReplicarApontamentos)

			// Sub-rotas que operam em um funcionário específico: /funcionarios/{funcionarioId}
			r.Route("/{funcionarioId}", func(r chi.Router) {
				r.With(auth.Authorize(authz.PermissaoPessoalLer)).
					Get("/", c.PessoalHandler.HandleBuscarFuncionario)

				r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).
					Put("/", c.PessoalHandler.HandleAtualizarFuncionario)

				r.With(auth.Authorize(authz.PermissaoPessoalApontamentoLer)).
					Patch("/ativar", c.PessoalHandler.HandleAtivarFuncionario)

				r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).
					Delete("/", c.PessoalHandler.HandleDeletarFuncionario)

				// Rota aninhada para listar os apontamentos deste funcionário
				r.With(auth.Authorize(authz.PermissaoPessoalApontamentoLer)).
					Get("/apontamentos", c.PessoalHandler.HandleListarApontamentosPorFuncionario)

			})
		})

		r.With(auth.Authorize(authz.PermissaoPessoalApontamentoEscrever)).
			Post("/apontamentos", c.PessoalHandler.HandleCriarApontamento)
		r.With(auth.Authorize(authz.PermissaoPessoalApontamentoLer)).
			Get("/apontamentos", c.PessoalHandler.HandleListarApontamentos)
		r.Route("/apontamentos/{apontamentoId}", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoPessoalApontamentoAprovar)).
				Patch("/aprovar", c.PessoalHandler.HandleAprovarApontamento)
			r.With(auth.Authorize(authz.PermissaoPessoalApontamentoPagar)).
				Patch("/pagar", c.PessoalHandler.HandleRegistrarPagamentoApontamento)
		})
		// --- Recursos de Suprimentos ---
		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Post("/fornecedores", c.SuprimentosHandler.HandleCadastrarFornecedor)
		r.With(auth.Authorize(authz.PermissaoSuprimentosLer)).Get("/fornecedores", c.SuprimentosHandler.HandleListarFornecedores)
		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Put("/fornecedores/{id}", c.SuprimentosHandler.HandleAtualizarFornecedor)
		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Delete("/fornecedores/{id}", c.SuprimentosHandler.HandleDeletarFornecedor)
		r.With(auth.Authorize(authz.PermissaoSuprimentosLer)).Get("/fornecedores/{id}", c.SuprimentosHandler.HandleBuscarFornecedor)
		// --- Recursos de Materiais ---
		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Post("/materiais", c.SuprimentosHandler.HandleCadastrarMaterial)
		r.With(auth.Authorize(authz.PermissaoSuprimentosLer)).Get("/materiais", c.SuprimentosHandler.HandleListarMateriais)
		r.With(auth.Authorize(authz.PermissaoSuprimentosLer)).Get("/materiais/{materialId}", c.SuprimentosHandler.HandleBuscarMaterial)
		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Put("/materiais/{materialId}", c.SuprimentosHandler.HandleAtualizarMaterial)
		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Delete("/materiais/{materialId}", c.SuprimentosHandler.HandleDeletarMaterial)
		// categorias do material
		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Post("/categorias", c.SuprimentosHandler.HandleCriarCategoria)
		r.With(auth.Authorize(authz.PermissaoSuprimentosLer)).Get("/categorias", c.SuprimentosHandler.HandleListarCategorias)
		r.With(auth.Authorize(authz.PermissaoSuprimentosLer)).Get("/categorias/{categoriaId}", c.SuprimentosHandler.HandleBuscarCategoria)
		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Put("/categorias/{categoriaId}", c.SuprimentosHandler.HandleAtualizarCategoria)
		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Delete("/categorias/{categoriaId}", c.SuprimentosHandler.HandleDeletarCategoria)

		// --- Recursos de Obras ---
		r.Route("/obras", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoObrasLer)).Get("/", c.ObrasHandler.HandleListarObras)
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/", c.ObrasHandler.HandleCriarObra)

			// Sub-recursos de uma obra específica
			r.Route("/{obraId}", func(r chi.Router) {
				r.With(auth.Authorize(authz.PermissaoObrasLer)).Get("/dashboard", c.ObrasHandler.HandleBuscarObra)
				r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Delete("/", c.ObrasHandler.HandleDeletarObra) // NOVA ROTA
				r.With(auth.Authorize(authz.PermissaoObrasLer)).Get("/", c.ObrasHandler.HandleBuscarObraPorID)
				r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Put("/", c.ObrasHandler.HandleAtualizarObra)
				r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/etapas", c.ObrasHandler.HandleAdicionarEtapa)
				r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/alocacoes", c.ObrasHandler.HandleAlocarFuncionario)
				r.With(auth.Authorize(authz.PermissaoObrasLer)).
					Get("/etapas", c.ObrasHandler.HandleListarEtapasPorObra)

			})
		})

		// --- Recursos de Etapas ---
		// Uma etapa pode ser tratada como um recurso de nível superior,
		// pois seu ID já é único.
		r.Route("/etapas/{etapaId}", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Patch("/", c.ObrasHandler.HandleAtualizarEtapaStatus)
			r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Post("/orcamentos", c.SuprimentosHandler.HandleCriarOrcamento)
		})

		r.Route("/orcamentos", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoSuprimentosLer)).
				Get("/", c.SuprimentosHandler.HandleListarOrcamentos)

			r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).
				Put("/{orcamentoId}", c.SuprimentosHandler.HandleAtualizarOrcamento)

			r.With(auth.Authorize(authz.PermissaoSuprimentosLer)).
				Get("/{orcamentoId}", c.SuprimentosHandler.HandleBuscarOrcamentoPorID)

			r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).
				Patch("/{orcamentoId}/status", c.SuprimentosHandler.HandleAtualizarOrcamentoStatus)
			
			r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).
				Delete("/{orcamentoId}", c.SuprimentosHandler.HandleDeletarOrcamento)
		})

		// --- Recursos de Financeiro ---
		r.With(auth.Authorize(authz.PermissaoFinanceiroEscrever)).
			Post("/pagamentos", c.FinanceiroHandler.HandleRegistrarPagamento)
		r.With(auth.Authorize(authz.PermissaoFinanceiroEscrever)).
			Post("/pagamentos/lote", c.FinanceiroHandler.HandleRegistrarPagamentosEmLote)
		r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).Get("/pagamentos", c.FinanceiroHandler.HandleListarPagamentos)

		// --- Contas a Receber ---
		r.Route("/contas-receber", func(r chi.Router) {
			// CRUD básico
			r.With(auth.Authorize(authz.PermissaoFinanceiroEscrever)).
				Post("/", c.ContaReceberHandler.HandleCriarConta)
			r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
				Get("/", c.ContaReceberHandler.HandleListarContas)
			r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
				Get("/{contaId}", c.ContaReceberHandler.HandleBuscarConta)
			
			// Ações específicas
			r.With(auth.Authorize(authz.PermissaoFinanceiroEscrever)).
				Post("/{contaId}/recebimentos", c.ContaReceberHandler.HandleRegistrarRecebimento)
			
			// Relatórios e consultas
			r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
				Get("/vencidas", c.ContaReceberHandler.HandleListarContasVencidas)
			r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
				Get("/resumo", c.ContaReceberHandler.HandleObterResumo)
		})

		// --- Contas a Pagar ---
		r.Route("/contas-pagar", func(r chi.Router) {
			// CRUD básico
			r.With(auth.Authorize(authz.PermissaoFinanceiroEscrever)).
				Post("/", c.ContaPagarHandler.HandleCriarConta)
			r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
				Get("/", c.ContaPagarHandler.HandleListarContas)
			r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
				Get("/{contaId}", c.ContaPagarHandler.HandleBuscarConta)
			
			// Ações específicas
			r.With(auth.Authorize(authz.PermissaoFinanceiroEscrever)).
				Post("/{contaId}/pagamentos", c.ContaPagarHandler.HandleRegistrarPagamento)
			r.With(auth.Authorize(authz.PermissaoFinanceiroEscrever)).
				Post("/orcamentos", c.ContaPagarHandler.HandleCriarContaDeOrcamento)
			
			// Relatórios e consultas
			r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
				Get("/vencidas", c.ContaPagarHandler.HandleListarContasVencidas)
			r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
				Get("/resumo", c.ContaPagarHandler.HandleObterResumo)
		})

		// Rotas específicas por entidade relacionada
		r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
			Get("/obras/{obraId}/contas-receber", c.ContaReceberHandler.HandleListarContasPorObra)
		r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
			Get("/obras/{obraId}/contas-pagar", c.ContaPagarHandler.HandleListarContasPorObra)
		r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
			Get("/fornecedores/{fornecedorId}/contas-pagar", c.ContaPagarHandler.HandleListarContasPorFornecedor)

		r.Route("/etapas-padroes", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoObrasLer)).Get("/", c.ObrasHandler.HandleListarEtapasPadrao)
			r.With(auth.Authorize(authz.PermissaoObrasLer)).Get("/{etapaId}", c.ObrasHandler.HandleBuscarEtapaPadrao)
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/", c.ObrasHandler.HandleCriarEtapaPadrao)
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Put("/{etapaId}", c.ObrasHandler.HandleAtualizarEtapaPadrao)
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Delete("/{etapaId}", c.ObrasHandler.HandleDeletarEtapaPadrao)
		})

		// --- Recursos de Dashboard (COMENTADO PARA DEBUG) ---
		// r.Route("/dashboard", func(r chi.Router) {
		// 	// Dashboard completo - requer permissão de leitura geral
		// 	r.With(auth.Authorize(authz.PermissaoObrasLer)).
		// 		Get("/", c.DashboardHandler.HandleObterDashboardCompleto)

		// 	// Seções específicas do dashboard
		// 	r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
		// 		Get("/financeiro", c.DashboardHandler.HandleObterDashboardFinanceiro)
			
		// 	r.With(auth.Authorize(authz.PermissaoObrasLer)).
		// 		Get("/obras", c.DashboardHandler.HandleObterDashboardObras)
			
		// 	r.With(auth.Authorize(authz.PermissaoPessoalLer)).
		// 		Get("/funcionarios", c.DashboardHandler.HandleObterDashboardFuncionarios)
			
		// 	r.With(auth.Authorize(authz.PermissaoSuprimentosLer)).
		// 		Get("/fornecedores", c.DashboardHandler.HandleObterDashboardFornecedores)

		// 	// Endpoints específicos
		// 	r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).
		// 		Get("/fluxo-caixa", c.DashboardHandler.HandleObterFluxoCaixa)

		// 	// Endpoint genérico para seções por URL
		// 	r.With(auth.Authorize(authz.PermissaoObrasLer)).
		// 		Get("/{secao}", c.DashboardHandler.HandleObterDashboardPorSecao)

		// 	// Parâmetros de cache
		// 	r.Get("/cache-info", c.DashboardHandler.HandleObterParametrosCache)
		// })
	})

	return r
}
