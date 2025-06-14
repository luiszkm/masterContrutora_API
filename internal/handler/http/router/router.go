// file: internal/handler/http/router/router.go
package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/luiszkm/masterCostrutora/internal/authz"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/identidade"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/obras"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/pessoal"

	"github.com/luiszkm/masterCostrutora/pkg/auth"
)

// Config contém as dependências necessárias para configurar o roteador.
type Config struct {
	JwtService        *auth.JWTService
	IdentidadeHandler *identidade.Handler
	ObrasHandler      *obras.Handler
	PessoalHandler    *pessoal.Handler
}

// New cria e configura um novo roteador chi com todas as rotas da aplicação.
func New(c Config) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- Rotas Públicas ---
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "status: ok")
	})
	r.Post("/usuarios/registrar", c.IdentidadeHandler.HandleRegistrar)
	r.Post("/usuarios/login", c.IdentidadeHandler.HandleLogin)

	// --- Rotas de Pessoal (Protegidas) ---
	// Corrigindo o aninhamento: /funcionarios é uma rota de nível superior.
	r.Group(func(r chi.Router) {
		r.Use(c.JwtService.AuthMiddleware)
		r.Use(auth.Authorize(authz.PermissaoPessoalEscrever))
		r.Post("/funcionarios", c.PessoalHandler.HandleCadastrarFuncionario)
	})

	// --- Rotas de Obras (Protegidas) ---
	r.Route("/obras", func(r chi.Router) {
		r.Use(c.JwtService.AuthMiddleware) // Autenticação para todo o grupo de obras
		r.With(auth.Authorize(authz.PermissaoObrasLer)).Get("/", c.ObrasHandler.HandleListarObras)
		r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/", c.ObrasHandler.HandleCriarObra)

		r.Route("/{obraId}", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoObrasLer)).Get("/", c.ObrasHandler.HandleBuscarObra)
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/etapas", c.ObrasHandler.HandleAdicionarEtapa)
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Patch("/etapas/{etapaId}", c.ObrasHandler.HandleAtualizarEtapaStatus)
			r.With(auth.Authorize(authz.PermissaoObrasEscrever)).Post("/alocacoes", c.ObrasHandler.HandleAlocarFuncionario)

		})
	})

	return r
}
