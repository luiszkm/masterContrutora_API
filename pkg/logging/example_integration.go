package logging

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ExampleIntegration demonstra como integrar o sistema de logging
func ExampleIntegration() {
	// 1. Criar o logger da aplicação
	appLogger, err := NewAppLogger(Config{
		Service:      "master-construtora",
		LogDirectory: "./logs",
		LogLevel:     slog.LevelDebug,
	})
	if err != nil {
		panic(err)
	}
	defer appLogger.Close()

	// 2. Criar logger específico do dashboard
	_ = NewDashboardLogger(appLogger) // dashLogger := NewDashboardLogger(appLogger)

	// 3. Configurar router com middlewares de logging
	r := chi.NewRouter()
	
	// Middlewares de logging - ORDEM IMPORTANTE
	// NOTA: Importe github.com/luiszkm/masterCostrutora/pkg/middleware no seu código
	// r.Use(middleware.RequestID)                          // 1. Request ID primeiro
	// r.Use(middleware.RealIP)                            // 2. IP real
	// r.Use(middleware.RequestLogger(appLogger, jwtService)) // 3. Log de requisições
	// r.Use(middleware.ErrorRecovery(appLogger))          // 4. Recovery de panics
	// r.Use(middleware.ErrorLogger(appLogger))            // 5. Log de erros HTTP

	// 4. Exemplo de uso em handlers
	r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		// O middleware já faz o logging automático
		// Você só precisa usar o logger para erros específicos da aplicação
		
		// Exemplo de log de erro de negócio:
		// dashLogger.LogDashboardError(r.Context(), "financeiro", "calcular_fluxo", err, nil)
		
		w.WriteHeader(http.StatusOK)
	})

	// 5. Inicializar servidor
	// http.ListenAndServe(":8080", r)
}

// ExampleServiceIntegration demonstra como integrar no serviço
func ExampleServiceIntegration() {
	// No seu main.go ou onde você configura os serviços:
	
	// 1. Criar logger
	appLogger, err := NewAppLogger(Config{
		Service:      "master-construtora",
		LogDirectory: "./logs",
		LogLevel:     slog.LevelInfo,
	})
	if err != nil {
		panic(err)
	}
	
	// 2. Criar logger do dashboard
	_ = NewDashboardLogger(appLogger) // dashLogger := NewDashboardLogger(appLogger)
	
	// 3. Usar na criação dos serviços
	// dashboardService := dashboard.NovoServicoDashboard(querier, logger, dashLogger)
	
	// 4. Criar funções de logging para repositórios
	// NOTA: Importe github.com/luiszkm/masterCostrutora/pkg/middleware no seu código
	// dbErrorLogger := middleware.DatabaseErrorLogger(appLogger)
	// serviceErrorLogger := middleware.ServiceErrorLogger(appLogger)
	
	// 5. Usar nos repositórios
	// _ = dbErrorLogger // repo.SetErrorLogger(dbErrorLogger)
	// _ = serviceErrorLogger // service.SetErrorLogger(serviceErrorLogger)
}