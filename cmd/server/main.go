// file: cmd/server/main.go
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	// --- Importações Internas Padronizadas ---

	"github.com/luiszkm/masterCostrutora/internal/events"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/router"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	"github.com/luiszkm/masterCostrutora/pkg/auth"
	"github.com/luiszkm/masterCostrutora/pkg/security"

	// Usaremos um único nome 'postgres' para o pacote de repositório para clareza

	financeiro_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/financeiro"
	identidade_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/identidade"
	obras_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/obras"
	pessoal_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/pessoal"
	suprimentos_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/suprimentos"

	financeiro_service "github.com/luiszkm/masterCostrutora/internal/service/financeiro"
	identidade_service "github.com/luiszkm/masterCostrutora/internal/service/identidade"
	obras_service "github.com/luiszkm/masterCostrutora/internal/service/obras"
	pessoal_service "github.com/luiszkm/masterCostrutora/internal/service/pessoal"
	suprimentos_service "github.com/luiszkm/masterCostrutora/internal/service/suprimentos"

	obras_events "github.com/luiszkm/masterCostrutora/internal/service/obras/events"
)

func main() {
	// 1. Configuração do Logger Estruturado (Correto)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	logger.Info("iniciando o sistema Master Construtora")

	// 2. Carregamento de Configurações (Correto)
	if err := godotenv.Load(); err != nil {
		logger.Warn("arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
	dbURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if dbURL == "" || jwtSecret == "" {
		logger.Error("as variáveis de ambiente DATABASE_URL e JWT_SECRET_KEY são obrigatórias")
		os.Exit(1)
	}

	// 3. Inicialização de Plataforma (Correto)
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Error("não foi possível conectar ao banco de dados", "erro", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	logger.Info("conexão com o PostgreSQL estabelecida com sucesso")

	jwtService := auth.NewJWTService(jwtSecret)
	passwordHasher := security.NewBcryptHasher()
	eventBus := bus.NovoEventBus(logger.With("component", "EventBus"))

	// --- 4. Injeção de Dependências (Wiring) - SEÇÃO CORRIGIDA ---

	// Repositórios Concretos
	usuarioRepo := postgres.NewUsuarioRepository(dbpool, logger)
	obraRepo := postgres.NovaObraRepository(dbpool, logger)
	etapaRepo := postgres.NovoEtapaRepository(dbpool, logger)
	alocacaoRepo := postgres.NovoAlocacaoRepository(dbpool, logger)
	funcionarioRepo := postgres.NovoFuncionarioRepository(dbpool, logger)
	fornecedorRepo := postgres.NovoFornecedorRepository(dbpool, logger)
	materialRepo := postgres.NovoMaterialRepository(dbpool, logger)
	orcamentoRepo := postgres.NovoOrcamentoRepository(dbpool, logger)
	financeiroRepo := postgres.NovoRegistroPagamentoRepository(dbpool, logger)
	apontamentoRepo := postgres.NovoApontamentoRepository(dbpool, logger)
	categoriaRepo := postgres.NovoCategoriaRepository(dbpool, logger)

	// Serviços
	identidadeSvc := identidade_service.NovoServico(usuarioRepo, passwordHasher, jwtService, logger)
	pessoalSvc := pessoal_service.NovoServico(
		funcionarioRepo, // Satisafaz pessoal.FuncionarioRepository
		apontamentoRepo, // A dependência que estava faltando
		alocacaoRepo,    // Satisafaz pessoal.AlocacaoFinder
		obraRepo,        // Satisafaz pessoal.ObraFinder
		eventBus,        // Satisafaz pessoal.EventPublisher
		funcionarioRepo,
		logger,
	)
	financeiroSvc := financeiro_service.NovoServico(financeiroRepo, funcionarioRepo, obraRepo, logger)
	obraSvc := obras_service.NovoServico(
		obraRepo,
		etapaRepo,
		alocacaoRepo,
		funcionarioRepo, // PessoalFinder implementado por FuncionarioRepository
		obraRepo,
		logger,
	)

	suprimentosSvc := suprimentos_service.NovoServico(
		fornecedorRepo,
		materialRepo,
		orcamentoRepo,
		categoriaRepo,
		etapaRepo,      // EtapaRepository implementa a interface EtapaFinder
		fornecedorRepo, // FornecedorRepository implementa a interface FornecedorFinder
		materialRepo,   // MaterialRepository implementa a interface MaterialFinder
		eventBus,
		logger,
	)

	// Handlers HTTP
	identidadeHandler := identidade_handler.NovoIdentidadeHandler(identidadeSvc, logger)
	pessoalHandler := pessoal_handler.NovoPessoalHandler(pessoalSvc, logger)
	obraHandler := obras_handler.NovoObrasHandler(obraSvc, logger)
	finaceiroHandler := financeiro_handler.NovoFinanceiroHandler(financeiroSvc, logger)
	// CORREÇÃO: Usando a variável com nome correto 'suprimentosSvc'.
	suprimentosHandler := suprimentos_handler.NovoSuprimentosHandler(suprimentosSvc, logger)

	// 4. Configuração do Event Bus e Manipuladores de Eventos (Correto)
	obrasEventHandler := obras_events.NovoObrasEventHandler(logger)
	eventBus.Subscrever(events.OrcamentoStatusAtualizado, obrasEventHandler.HandleOrcamentoStatusAtualizado)

	// 5. Configuração do Servidor HTTP e Roteamento (Correto)
	routerCfg := router.Config{
		JwtService:         jwtService,
		IdentidadeHandler:  identidadeHandler,
		ObrasHandler:       obraHandler,
		PessoalHandler:     pessoalHandler,
		SuprimentosHandler: suprimentosHandler,
		FinanceiroHandler:  finaceiroHandler,
	}
	r := router.New(routerCfg)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	logger.Info("servidor escutando na porta :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("não foi possível iniciar o servidor: %v", err)
	}
}
