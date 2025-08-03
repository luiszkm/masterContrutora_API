package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// TestRealEndpoint testa o endpoint real do dashboard
func TestRealEndpoint(t *testing.T) {
	fmt.Println("=== Teste Real do Endpoint Dashboard ===")

	// Lista de endpoints para testar
	endpoints := []struct {
		name string
		url  string
		desc string
	}{
		{
			name: "Dashboard_Geral",
			url:  "http://localhost:8080/dashboard",
			desc: "Dashboard geral completo",
		},
		{
			name: "Dashboard_Financeiro",
			url:  "http://localhost:8080/dashboard/financeiro",
			desc: "Seção financeira",
		},
		{
			name: "Dashboard_Obras",
			url:  "http://localhost:8080/dashboard/obras",
			desc: "Seção obras",
		},
		{
			name: "Dashboard_Funcionarios",
			url:  "http://localhost:8080/dashboard/funcionarios",
			desc: "Seção funcionários",
		},
		{
			name: "Dashboard_Fornecedores",
			url:  "http://localhost:8080/dashboard/fornecedores",
			desc: "Seção fornecedores",
		},
		{
			name: "Dashboard_FluxoCaixa",
			url:  "http://localhost:8080/dashboard/fluxo-caixa",
			desc: "Fluxo de caixa",
		},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			fmt.Printf("\n=== Testando: %s ===\n", endpoint.desc)
			fmt.Printf("URL: %s\n", endpoint.url)

			// Fazer request simples
			resp, err := http.Get(endpoint.url)
			if err != nil {
				fmt.Printf("❌ Erro ao fazer request: %v\n", err)
				t.Errorf("Erro de conectividade: %v", err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Status: %d %s\n", resp.StatusCode, resp.Status)

			// Ler corpo da resposta
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("❌ Erro ao ler corpo: %v\n", err)
				return
			}

			// Se é erro 500, mostrar detalhes
			if resp.StatusCode == 500 {
				fmt.Printf("❌ ERRO 500 ENCONTRADO!\n")
				fmt.Printf("Corpo da resposta:\n%s\n", string(body))
				
				// Tentar fazer parse do JSON de erro
				var errorResponse map[string]interface{}
				if err := json.Unmarshal(body, &errorResponse); err == nil {
					fmt.Printf("Erro estruturado: %+v\n", errorResponse)
				}
				
				t.Errorf("Erro 500 em %s", endpoint.url)
				return
			}

			// Se é sucesso, verificar se resposta é JSON válido
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					fmt.Printf("⚠️ Resposta não é JSON válido: %v\n", err)
					fmt.Printf("Corpo: %s\n", string(body)[:200]) // Primeiros 200 chars
				} else {
					fmt.Printf("✅ Sucesso! Campos na resposta: %v\n", getJSONKeys(response))
				}
			} else {
				fmt.Printf("⚠️ Status não-sucesso: %d\n", resp.StatusCode)
				fmt.Printf("Corpo: %s\n", string(body))
			}
		})
	}
}

// TestHealthCheck testa se servidor está respondendo
func TestHealthCheck(t *testing.T) {
	fmt.Println("=== Teste de Health Check ===")

	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Skipf("❌ Servidor não está rodando: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Health check status: %d\n", resp.StatusCode)

	if resp.StatusCode == 200 {
		fmt.Println("✅ Servidor está rodando!")
	} else {
		fmt.Printf("⚠️ Health check retornou: %d\n", resp.StatusCode)
	}
}

// Função auxiliar para extrair chaves do JSON
func getJSONKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}