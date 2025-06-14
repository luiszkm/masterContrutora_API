// file: cmd/server/main.go
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	// Importações organizadas por função

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	// Importações internas do projeto

	obras_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/obras"
	"github.com/luiszkm/masterCostrutora/internal/handler/http/router"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	obras_repository "github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	obras_service "github.com/luiszkm/masterCostrutora/internal/service/obras"
	"github.com/luiszkm/masterCostrutora/pkg/auth"
	"github.com/luiszkm/masterCostrutora/pkg/security"

	identidade_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/identidade"
	identidade_service "github.com/luiszkm/masterCostrutora/internal/service/identidade"

	pessoal_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/pessoal"
	pessoal_service "github.com/luiszkm/masterCostrutora/internal/service/pessoal"
)

func main() {
	// 1. Configuração do Logger Estruturado (ADR-008)
	// Logs serão em JSON para facilitar a futura integração com plataformas de logging.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("iniciando o sistema Master Construtora")

	// 2. Carregamento de Configurações (ex: .env)
	if err := godotenv.Load(); err != nil {
		logger.Warn("arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
	dbURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if dbURL == "" || jwtSecret == "" {
		logger.Error("as variáveis de ambiente DATABASE_URL e JWT_SECRET_KEY são obrigatórias")
		os.Exit(1)
	}

	// 3. Inicialização do Banco de Dados (ADR-002)
	// Usamos um pool de conexões para performance e resiliência.
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Error("não foi possível conectar ao banco de dados", "erro", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	logger.Info("conexão com o PostgreSQL estabelecida com sucesso")

	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET_KEY"))
	passwordHasher := security.NewBcryptHasher()

	// 4. Injeção de Dependências (Wiring)
	// Aqui conectamos as implementações concretas com as interfaces.

	// Repositório
	obraRepo := obras_repository.NovaObraRepository(dbpool, logger)
	etapaRepo := postgres.NovoEtapaRepository(dbpool, logger)
	usuarioRepo := postgres.NewUsuarioRepository(dbpool, logger)
	funcionarioRepo := postgres.NovoFuncionarioRepository(dbpool, logger)
	alocacaoRepo := postgres.NovoAlocacaoRepository(dbpool, logger)

	// Serviço
	obraSvc := obras_service.NovoServico(obraRepo, etapaRepo, alocacaoRepo, funcionarioRepo, obraRepo, logger)
	identidadeSvc := identidade_service.NovoServico(usuarioRepo, passwordHasher, jwtService, logger.With("component", "IdentidadeService"))
	pessoalSvc := pessoal_service.NovoServico(funcionarioRepo, logger.With("component", "PessoalService"))

	// Handler HTTP
	obraHandler := obras_handler.NovoObrasHandler(obraSvc, logger.With("component", "ObrasHandler"))
	identidadeHandler := identidade_handler.NovoIdentidadeHandler(identidadeSvc, logger.With("component", "IdentidadeHandler"))
	pessoalHandler := pessoal_handler.NovoPessoalHandler(pessoalSvc, logger.With("component", "PessoalHandler"))

	// 5. Configuração do Servidor HTTP e Roteamento
	routerCfg := router.Config{
		JwtService:        jwtService,
		IdentidadeHandler: identidadeHandler,
		ObrasHandler:      obraHandler,
		PessoalHandler:    pessoalHandler,
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
