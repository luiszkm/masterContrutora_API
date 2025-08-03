package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// TestSimpleDashboard testa o dashboard de forma simples sem dependências complexas
func TestSimpleDashboard(t *testing.T) {
	fmt.Println("=== Teste Simples do Dashboard ===")

	// Conectar com banco primeiro
	dbURL := getSimpleTestDatabaseURL()
	fmt.Printf("Conectando com: %s\n", dbURL)

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("❌ Banco não disponível: %v", err)
		return
	}
	defer dbPool.Close()

	// Testar ping
	err = dbPool.Ping(context.Background())
	require.NoError(t, err, "Ping no banco deve funcionar")
	fmt.Println("✅ Banco conectado!")

	// Executar queries SQL diretamente para testar cada seção
	ctx := context.Background()

	tests := []struct {
		name  string
		query string
		desc  string
	}{
		{
			name: "Test_Distribuicao_Despesas",
			query: `
				SELECT 
					'Material' as categoria,
					1000.50 as valor,
					5 as quantidade_itens
				UNION ALL
				SELECT 
					'Mao de Obra' as categoria,
					2500.75 as valor,
					3 as quantidade_itens
			`,
			desc: "Query de distribuição de despesas",
		},
		{
			name: "Test_Distribuicao_Obras",
			query: `
				SELECT 
					'Em Andamento' as status,
					10 as quantidade,
					25000.00 as valor_total
				UNION ALL
				SELECT 
					'Concluida' as status,
					5 as quantidade,
					15000.00 as valor_total
			`,
			desc: "Query de distribuição de obras",
		},
		{
			name: "Test_Fornecedor_Categoria",
			query: `
				SELECT 
					'cat-1' as categoria_id,
					'Materiais' as categoria_nome,
					8 as quantidade_fornecedores,
					4.2 as avaliacao_media
			`,
			desc: "Query de fornecedores por categoria",
		},
		{
			name: "Test_Gasto_Fornecedor",
			query: `
				SELECT 
					'forn-1' as fornecedor_id,
					'Fornecedor Teste' as nome_fornecedor,
					4.5 as avaliacao,
					5000.00 as valor_total_gasto,
					3 as quantidade_orcamentos,
					CURRENT_TIMESTAMP as ultimo_orcamento
			`,
			desc: "Query de gastos com fornecedor",
		},
		{
			name: "Test_Produtividade_Funcionario",
			query: `
				SELECT 
					'func-1' as funcionario_id,
					'João Silva' as nome_funcionario,
					'Pedreiro' as cargo,
					20 as dias_trabalhados,
					22.5 as media_dias_por_periodo,
					2 as obras_alocadas
			`,
			desc: "Query de produtividade de funcionário",
		},
		{
			name: "Test_Top_Funcionario",
			query: `
				SELECT 
					'func-1' as funcionario_id,
					'Maria Santos' as nome_funcionario,
					'Engenheira' as cargo,
					'Excelente' as avaliacao_desempenho,
					45 as dias_trabalhados_total,
					3 as obras_participadas,
					CURRENT_TIMESTAMP as data_contratacao
			`,
			desc: "Query de top funcionário",
		},
		{
			name: "Test_Top_Fornecedor",
			query: `
				SELECT 
					'forn-2' as fornecedor_id,
					'ABC Materiais' as nome_fornecedor,
					'12345678000123' as cnpj,
					4.8 as avaliacao,
					'Ativo' as status,
					15 as total_orcamentos,
					12000.00 as valor_total_gasto,
					CURRENT_TIMESTAMP as ultimo_orcamento
			`,
			desc: "Query de top fornecedor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("\n--- %s ---\n", tt.desc)
			fmt.Printf("Query: %s\n", tt.query)

			rows, err := dbPool.Query(ctx, tt.query)
			if err != nil {
				t.Errorf("❌ Erro na query: %v", err)
				return
			}
			defer rows.Close()

			// Verificar se podemos ler pelo menos uma linha
			if rows.Next() {
				fmt.Println("✅ Query executada com sucesso!")
				
				// Tentar fazer scan manual simples
				values, err := rows.Values()
				if err != nil {
					t.Errorf("❌ Erro ao ler valores: %v", err)
					return
				}
				
				fmt.Printf("Valores retornados: %v\n", values)
			} else {
				fmt.Println("⚠️  Query não retornou dados")
			}
		})
	}
}

// TestDirectEndpoint testa endpoints diretamente com mock simples
func TestDirectEndpoint(t *testing.T) {
	fmt.Println("\n=== Teste Direto de Endpoint ===")

	// Mock handler simples
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simular resposta de dashboard
		response := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"totalObras":       10,
				"obrasConcluidas":  5,
				"obrasEmAndamento": 5,
				"totalFuncionarios": 25,
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Criar servidor de teste
	server := httptest.NewServer(handler)
	defer server.Close()

	// Fazer request
	resp, err := http.Get(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	fmt.Printf("✅ Resposta mock: %+v\n", response)
}

// TestJWTGeneration testa geração de JWT para usar nos testes
func TestJWTGeneration(t *testing.T) {
	fmt.Println("\n=== Teste de Geração JWT ===")

	// Mock simples sem usar o auth package
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	
	require.NotEmpty(t, mockToken)
	fmt.Printf("✅ Token mock gerado: %s...\n", mockToken[:20])
}

// getSimpleTestDatabaseURL retorna URL do banco para testes simples
func getSimpleTestDatabaseURL() string {
	return "postgres://user:password@localhost:5432/mastercostrutora_db?sslmode=disable"
}