// file: internal/handler/http/router/router.go
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/luiszkm/masterCostrutora/internal/authz"
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

	// --- GRUPO ÚNICO PARA TODAS AS ROTAS PROTEGIDAS ---
	r.Group(func(r chi.Router) {
		// Aplicamos o middleware de autenticação UMA VEZ para todo o grupo.
		r.Use(c.JwtService.AuthMiddleware)

		// --- Recursos de Pessoal ---
		r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).Post("/funcionarios", c.PessoalHandler.HandleCadastrarFuncionario)
		r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).Delete("/funcionarios/{id}", c.PessoalHandler.HandleDeletarFuncionario)
		r.With(auth.Authorize(authz.PermissaoPessoalLer)).Get("/funcionarios", c.PessoalHandler.HandleListarFuncionarios)
		r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).Put("/funcionarios/{id}", c.PessoalHandler.HandleAtualizarFuncionario)
		r.With(auth.Authorize(authz.PermissaoPessoalLer)).Get("/funcionarios/{id}", c.PessoalHandler.HandleBuscarFuncionario)
		r.Route("/funcionarios/{funcionarioId}", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoPessoalApontamentoLer)).
				Get("/apontamentos", c.PessoalHandler.HandleListarApontamentosPorFuncionario)
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
		// Apontamentos são um recurso de Pessoal
		r.With(auth.Authorize(authz.PermissaoPessoalApontamentoEscrever)).
			Post("/apontamentos", c.PessoalHandler.HandleCriarApontamento)
		r.Route("/apontamentos/{apontamentoId}", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoPessoalApontamentoAprovar)).
				Patch("/aprovar", c.PessoalHandler.HandleAprovarApontamento)
		})
		r.With(auth.Authorize(authz.PermissaoPessoalApontamentoLer)).
			Get("/apontamentos", c.PessoalHandler.HandleListarApontamentos)
		r.With(auth.Authorize(authz.PermissaoPessoalApontamentoPagar)).
			Patch("/pagar", c.PessoalHandler.HandleRegistrarPagamentoApontamento)
		// --- Recursos de Obras ---
		r.Route("/obras", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoObrasLer)).Get("/", c.ObrasHandler.HandleListarObras)
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/", c.ObrasHandler.HandleCriarObra)

			// Sub-recursos de uma obra específica
			r.Route("/{obraId}", func(r chi.Router) {
				r.With(auth.Authorize(authz.PermissaoObrasLer)).Get("/", c.ObrasHandler.HandleBuscarObra)
				r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Delete("/", c.ObrasHandler.HandleDeletarObra) // NOVA ROTA

				r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/etapas", c.ObrasHandler.HandleAdicionarEtapa)
				r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/alocacoes", c.ObrasHandler.HandleAlocarFuncionario)
			})
		})

		// --- Recursos de Etapas ---
		// Uma etapa pode ser tratada como um recurso de nível superior,
		// pois seu ID já é único.
		r.Route("/etapas/{etapaId}", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Patch("/", c.ObrasHandler.HandleAtualizarEtapaStatus)
			r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).Post("/orcamentos", c.SuprimentosHandler.HandleCriarOrcamento)
		})

		r.With(auth.Authorize(authz.PermissaoSuprimentosEscrever)).
			Patch("/orcamentos/{orcamentoId}", c.SuprimentosHandler.HandleAtualizarOrcamentoStatus)
		// --- Recursos de Financeiro ---
		r.With(auth.Authorize(authz.PermissaoFinanceiroEscrever)).Post("/pagamentos", c.FinanceiroHandler.HandleRegistrarPagamento)
		// r.With(auth.Authorize(authz.PermissaoFinanceiroLer)).Get("/pagamentos", c.FinanceiroHandler.HandleListarPagamentos)
	})

	return r
}
