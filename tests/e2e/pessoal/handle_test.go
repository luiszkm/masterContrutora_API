// file: internal/handler/http/pessoal/handler_test.go
package pessoal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luiszkm/masterCostrutora/internal/authz"
	"github.com/luiszkm/masterCostrutora/internal/domain/common"
	"github.com/luiszkm/masterCostrutora/internal/domain/pessoal"
	identidade_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/identidade"
	pessoal_handler "github.com/luiszkm/masterCostrutora/internal/handler/http/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/infrastructure/repository/postgres"
	"github.com/luiszkm/masterCostrutora/internal/platform/bus"
	identidade_service "github.com/luiszkm/masterCostrutora/internal/service/identidade"
	pessoal_service "github.com/luiszkm/masterCostrutora/internal/service/pessoal"
	"github.com/luiszkm/masterCostrutora/internal/service/pessoal/dto"
	"github.com/luiszkm/masterCostrutora/pkg/auth"
	"github.com/luiszkm/masterCostrutora/pkg/security"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestApplication inicializa uma aplicação completa para o teste.
// setupTestApplication agora configura todas as rotas de Pessoal e Identidade.
func setupTestApplication(t *testing.T) (*chi.Mux, *pgxpool.Pool) {
	dbURL := "postgres://user:password@localhost:5432/mastercostrutora_db?sslmode=disable"
	if dbURL == "" {
		t.Skip("DATABASE_URL não definida, pulando teste de integração")
	}

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)

	cleanupDatabase(t, dbpool)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	// Injeção de dependência completa
	// Plataforma
	passwordHasher := security.NewBcryptHasher()
	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET_KEY"))
	eventBus := bus.NovoEventBus(logger.With("component", "EventBus"))

	// Repositórios
	usuarioRepo := postgres.NewUsuarioRepository(dbpool, logger)
	funcionarioRepo := postgres.NovoFuncionarioRepository(dbpool, logger)
	apontamentoRepo := postgres.NovoApontamentoRepository(dbpool, logger)
	alocacaoRepo := postgres.NovoAlocacaoRepository(dbpool, logger)
	obraRepo := postgres.NovaObraRepository(dbpool, logger)

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
	// Handlers
	pessoalHandler := pessoal_handler.NovoPessoalHandler(pessoalSvc, logger)

	identidadeHandler := identidade_handler.NovoIdentidadeHandler(identidadeSvc, logger)

	// Roteador
	r := chi.NewRouter()
	r.Post("/usuarios/registrar", identidadeHandler.HandleRegistrar)
	r.Post("/usuarios/login", identidadeHandler.HandleLogin)

	// Grupo de rotas protegidas
	r.Group(func(r chi.Router) {
		r.Use(jwtService.AuthMiddleware)

		// Rotas de Funcionários
		r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).Post("/funcionarios", pessoalHandler.HandleCadastrarFuncionario)
		r.With(auth.Authorize(authz.PermissaoPessoalLer)).Get("/funcionarios", pessoalHandler.HandleListarFuncionarios)
		r.With(auth.Authorize(authz.PermissaoPessoalLer)).Get("/funcionarios/{id}", pessoalHandler.HandleBuscarFuncionario)
		r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).Put("/funcionarios/{id}", pessoalHandler.HandleAtualizarFuncionario)
		r.With(auth.Authorize(authz.PermissaoPessoalEscrever)).Delete("/funcionarios/{id}", pessoalHandler.HandleDeletarFuncionario)
		r.With(auth.Authorize(authz.PermissaoPessoalApontamentoLer)).Get("/funcionarios/{id}/apontamentos", pessoalHandler.HandleListarApontamentosPorFuncionario)

		r.With(auth.Authorize(authz.PermissaoPessoalApontamentoEscrever)).Post("/apontamentos", pessoalHandler.HandleCriarApontamento)
		r.Route("/apontamentos/{apontamentoId}", func(r chi.Router) {
			r.With(auth.Authorize(authz.PermissaoPessoalApontamentoAprovar)).Patch("/aprovar", pessoalHandler.HandleAprovarApontamento)
			r.With(auth.Authorize(authz.PermissaoPessoalApontamentoPagar)).Patch("/pagar", pessoalHandler.HandleRegistrarPagamentoApontamento)
		})
		r.With(auth.Authorize(authz.PermissaoPessoalApontamentoLer)).Get("/apontamentos", pessoalHandler.HandleListarApontamentos)

	})

	return r, dbpool
}
func cleanupDatabase(t *testing.T, dbpool *pgxpool.Pool) {
	query := `
		DELETE FROM registros_pagamento; DELETE FROM orcamento_itens; DELETE FROM orcamentos;
		DELETE FROM apontamentos_quinzenais; DELETE FROM alocacoes; DELETE FROM etapas;
		DELETE FROM materiais; DELETE FROM fornecedores; DELETE FROM usuarios;
		DELETE FROM funcionarios; DELETE FROM obras;
	`
	_, err := dbpool.Exec(context.Background(), query)
	require.NoError(t, err, "A limpeza do banco de dados não deve falhar")
}
func loginParaTeste(t *testing.T, server *httptest.Server) *http.Cookie {
	// Registrar
	userPayload := []byte(`{"nome": "User Teste", "email": "teste@email.com", "senha": "123", "confirmarSenha": "123"}`)
	reqReg, _ := http.NewRequest(http.MethodPost, server.URL+"/usuarios/registrar", bytes.NewBuffer(userPayload))
	respReg, _ := server.Client().Do(reqReg)
	require.Equal(t, http.StatusCreated, respReg.StatusCode)

	// Login
	loginPayload := []byte(`{"email": "teste@email.com", "senha": "123"}`)
	reqLogin, _ := http.NewRequest(http.MethodPost, server.URL+"/usuarios/login", bytes.NewBuffer(loginPayload))
	respLogin, _ := server.Client().Do(reqLogin)
	require.Equal(t, http.StatusOK, respLogin.StatusCode)

	// Extrai o cookie da resposta do login
	cookies := respLogin.Cookies()
	require.NotEmpty(t, cookies, "O login deve retornar um cookie de autenticação")
	return cookies[0]
}

func TestFuncionarioHandlers_Integration(t *testing.T) {
	// Arrange
	router, _ := setupTestApplication(t)
	server := httptest.NewServer(router)
	defer server.Close()
	authCookie := loginParaTeste(t, server)
	// --- Teste de Criação ---
	var funcionarioID string
	t.Run("Deve criar um funcionário com sucesso", func(t *testing.T) {
		payload := []byte(`{"nome": "Clarice Lispector", "cpf": "999.888.777-66", "cargo": "Escritora", "salario": 7000.00, "valor_diaria": 60.00, "data_contratacao": "2025-01-01"}`)
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/funcionarios", bytes.NewBuffer(payload))
		req.AddCookie(authCookie)

		resp, _ := server.Client().Do(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var funcionarioRetornado pessoal.Funcionario
		json.NewDecoder(resp.Body).Decode(&funcionarioRetornado)
		assert.Equal(t, "Clarice Lispector", funcionarioRetornado.Nome)
		funcionarioID = funcionarioRetornado.ID // Salva o ID para os próximos testes
	})

	// --- Teste de Busca por ID ---
	t.Run("Deve buscar o funcionário criado pelo ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/funcionarios/"+funcionarioID, nil)
		req.AddCookie(authCookie)

		resp, _ := server.Client().Do(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var funcionarioRetornado pessoal.Funcionario
		json.NewDecoder(resp.Body).Decode(&funcionarioRetornado)
		assert.Equal(t, "Clarice Lispector", funcionarioRetornado.Nome)
	})

	// --- Teste de Listagem ---
	t.Run("Deve listar todos os funcionários", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/funcionarios", nil)
		req.AddCookie(authCookie)

		resp, _ := server.Client().Do(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var funcionarios []pessoal.Funcionario
		json.NewDecoder(resp.Body).Decode(&funcionarios)
		assert.Len(t, funcionarios, 1)
	})
	// --- Teste de Atualização ---
	t.Run("Deve atualizar um funcionário com sucesso", func(t *testing.T) {
		// Arrange
		updatePayload := []byte(`{
			"nome": "Clarice Lispector (Atualizado)",
			"cpf": "999.888.777-66",
			"cargo": "Romancista",
			"salario": 7500.00,
			"diaria": 80.00
		}`)
		req, _ := http.NewRequest(http.MethodPut, server.URL+"/funcionarios/"+funcionarioID, bytes.NewBuffer(updatePayload))
		req.AddCookie(authCookie)

		// Act
		resp, err := server.Client().Do(req)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

	})
	// --- Teste de Deleção ---
	t.Run("Deve deletar um funcionário com sucesso", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, server.URL+"/funcionarios/"+funcionarioID, nil)
		req.AddCookie(authCookie)

		resp, _ := server.Client().Do(req)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verifica que o funcionário não aparece mais na listagem
		reqList, _ := http.NewRequest(http.MethodGet, server.URL+"/funcionarios", nil)
		respList, _ := server.Client().Do(reqList)
		var funcionarios []pessoal.Funcionario
		json.NewDecoder(respList.Body).Decode(&funcionarios)
		assert.Len(t, funcionarios, 0)
	})
}

func TestApontamentoHandlers_Integration(t *testing.T) {
	// Arrange
	router, dbpool := setupTestApplication(t)
	server := httptest.NewServer(router)
	defer server.Close()
	ctx := context.Background()

	// Autentica para obter o cookie
	authCookie := loginParaTeste(t, server)

	// Prepara os dados necessários para os testes
	funcionarioID := uuid.NewString()
	_, err := dbpool.Exec(ctx, `INSERT INTO funcionarios (id, nome, cpf, cargo, data_contratacao, valor_diaria, status) VALUES ($1, 'Func Teste', '123', 'Cargo', '2025-01-01', 100, 'Ativo')`, funcionarioID)
	require.NoError(t, err)
	obraID := uuid.NewString()
	_, err = dbpool.Exec(ctx, `INSERT INTO obras (id, nome, cliente, endereco, data_inicio, status) VALUES ($1, 'Obra Teste', 'Cliente', 'End', '2025-01-01', 'Em Andamento')`, obraID)
	require.NoError(t, err)
	etapaID := uuid.NewString()
	_, err = dbpool.Exec(ctx, `INSERT INTO etapas (id, obra_id, nome, status) VALUES ($1, $2, 'Fundação', 'Em Andamento')`, etapaID, obraID)
	require.NoError(t, err)

	// --- Teste de Criação de Apontamento ---
	var apontamentoCriado pessoal.ApontamentoQuinzenal
	t.Run("Deve criar um apontamento com sucesso", func(t *testing.T) {
		payload, _ := json.Marshal(dto.CriarApontamentoInput{
			FuncionarioID: funcionarioID,
			ObraID:        obraID,
			PeriodoInicio: "2025-06-01",
			PeriodoFim:    "2025-06-15",
		})
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/apontamentos", bytes.NewBuffer(payload))
		req.AddCookie(authCookie)
		resp, _ := server.Client().Do(req)

		require.Equal(t, http.StatusCreated, resp.StatusCode)
		json.NewDecoder(resp.Body).Decode(&apontamentoCriado)
		assert.Equal(t, pessoal.StatusApontamentoEmAberto, apontamentoCriado.Status)
	})

	// --- Teste de Aprovação de Apontamento ---
	t.Run("Deve aprovar um apontamento com sucesso", func(t *testing.T) {
		require.NotEmpty(t, apontamentoCriado.ID)
		url := fmt.Sprintf("%s/apontamentos/%s/aprovar", server.URL, apontamentoCriado.ID)
		req, _ := http.NewRequest(http.MethodPatch, url, nil)
		req.AddCookie(authCookie)
		resp, _ := server.Client().Do(req)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		var apontamentoAprovado pessoal.ApontamentoQuinzenal
		json.NewDecoder(resp.Body).Decode(&apontamentoAprovado)
		assert.Equal(t, pessoal.StatusApontamentoAprovadoParaPagamento, apontamentoAprovado.Status)
	})

	// --- Teste de Pagamento de Apontamento ---
	t.Run("Deve registrar o pagamento de um apontamento com sucesso", func(t *testing.T) {
		require.NotEmpty(t, apontamentoCriado.ID)
		payload := []byte(`{"contaBancariaId": "f47ac10b-58cc-4372-a567-0e02b2c3d479"}`)
		url := fmt.Sprintf("%s/apontamentos/%s/pagar", server.URL, apontamentoCriado.ID)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
		req.AddCookie(authCookie)
		resp, _ := server.Client().Do(req)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		var apontamentoPago pessoal.ApontamentoQuinzenal
		json.NewDecoder(resp.Body).Decode(&apontamentoPago)
		assert.Equal(t, pessoal.StatusApontamentoPago, apontamentoPago.Status)
	})
}

// file: internal/handler/http/pessoal/handler_test.go
// ...

func TestListarApontamentos_Integration(t *testing.T) {
	// Arrange (Preparação)
	router, dbpool := setupTestApplication(t)
	server := httptest.NewServer(router)
	defer server.Close()
	ctx := context.Background()
	authCookie := loginParaTeste(t, server)

	// Cria dados de teste: 2 funcionários, 1 obra, 3 apontamentos
	funcA_ID := uuid.NewString()
	funcB_ID := uuid.NewString()
	obraID := uuid.NewString()
	_, err := dbpool.Exec(ctx, `INSERT INTO funcionarios (id, nome, cpf, cargo, data_contratacao, valor_diaria, status) VALUES ($1, 'Func A', '111', 'Cargo A', '2025-01-01', 100, 'Ativo'), ($2, 'Func B', '222', 'Cargo B', '2025-01-01', 100, 'Ativo')`, funcA_ID, funcB_ID)
	require.NoError(t, err)
	_, err = dbpool.Exec(ctx, `INSERT INTO obras (id, nome, cliente, endereco, data_inicio, status) VALUES ($1, 'Obra Teste Listagem', 'Cliente', 'End', '2025-01-01', 'Em Andamento')`, obraID)
	require.NoError(t, err)

	dbpool.Exec(ctx, `INSERT INTO apontamentos_quinzenais (id, funcionario_id, obra_id, periodo_inicio, periodo_fim, status, created_at, updated_at) VALUES ($1, $2, $3, '2025-01-01', '2025-01-15', 'EM_ABERTO', NOW(), NOW())`, uuid.NewString(), funcA_ID, obraID)
	dbpool.Exec(ctx, `INSERT INTO apontamentos_quinzenais (id, funcionario_id, obra_id, periodo_inicio, periodo_fim, status, created_at, updated_at) VALUES ($1, $2, $3, '2025-01-16', '2025-01-31', 'PAGO', NOW(), NOW())`, uuid.NewString(), funcA_ID, obraID)
	dbpool.Exec(ctx, `INSERT INTO apontamentos_quinzenais (id, funcionario_id, obra_id, periodo_inicio, periodo_fim, status, created_at, updated_at) VALUES ($1, $2, $3, '2025-02-01', '2025-02-15', 'EM_ABERTO', NOW(), NOW())`, uuid.NewString(), funcB_ID, obraID)

	t.Run("Deve listar todos os 3 apontamentos sem filtro", func(t *testing.T) {
		// Act
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/apontamentos", nil)
		req.AddCookie(authCookie)
		resp, _ := server.Client().Do(req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var resposta common.RespostaPaginada[pessoal.ApontamentoQuinzenal]
		json.NewDecoder(resp.Body).Decode(&resposta)
		assert.Equal(t, 3, resposta.Paginacao.TotalItens)
		assert.Len(t, resposta.Dados, 3)
	})

	t.Run("Deve listar apenas 2 apontamentos com status EM_ABERTO", func(t *testing.T) {
		// Act
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/apontamentos?status=EM_ABERTO", nil)
		req.AddCookie(authCookie)
		resp, _ := server.Client().Do(req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var resposta common.RespostaPaginada[pessoal.ApontamentoQuinzenal]
		json.NewDecoder(resp.Body).Decode(&resposta)
		assert.Equal(t, 2, resposta.Paginacao.TotalItens)
		assert.Len(t, resposta.Dados, 2)
	})

	t.Run("Deve listar os 2 apontamentos do Funcionario A", func(t *testing.T) {
		// Act
		url := fmt.Sprintf("%s/funcionarios/%s/apontamentos", server.URL, funcA_ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.AddCookie(authCookie)
		resp, _ := server.Client().Do(req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var resposta common.RespostaPaginada[pessoal.ApontamentoQuinzenal]
		json.NewDecoder(resp.Body).Decode(&resposta)
		assert.Equal(t, 2, resposta.Paginacao.TotalItens)
	})

	t.Run("Deve listar apenas 1 apontamento PAGO do Funcionario A", func(t *testing.T) {
		// Act
		url := fmt.Sprintf("%s/funcionarios/%s/apontamentos?status=PAGO", server.URL, funcA_ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.AddCookie(authCookie)
		resp, _ := server.Client().Do(req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var resposta common.RespostaPaginada[pessoal.ApontamentoQuinzenal]
		json.NewDecoder(resp.Body).Decode(&resposta)
		assert.Equal(t, 1, resposta.Paginacao.TotalItens)
		assert.Equal(t, "PAGO", resposta.Dados[0].Status)
	})
}
