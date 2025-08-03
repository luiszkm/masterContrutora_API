package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestWithAuth testa endpoints do dashboard com autenticação
func TestWithAuth(t *testing.T) {
	fmt.Println("=== Teste Dashboard com Autenticação ===")

	// 1. Fazer login primeiro
	token, err := doLogin()
	if err != nil {
		t.Skipf("❌ Não foi possível fazer login: %v", err)
		return
	}

	fmt.Printf("✅ Login realizado, token: %s...\n", token[:20])

	// 2. Testar endpoints com token
	endpoints := []struct {
		name string
		url  string
		desc string
	}{
		{
			name: "Dashboard_Geral_Auth",
			url:  "http://localhost:8080/dashboard",
			desc: "Dashboard geral com auth",
		},
		{
			name: "Dashboard_Financeiro_Auth",
			url:  "http://localhost:8080/dashboard/financeiro",
			desc: "Financeiro com auth",
		},
		{
			name: "Dashboard_Obras_Auth",
			url:  "http://localhost:8080/dashboard/obras",
			desc: "Obras com auth",
		},
		{
			name: "Dashboard_Funcionarios_Auth",
			url:  "http://localhost:8080/dashboard/funcionarios",
			desc: "Funcionários com auth",
		},
		{
			name: "Dashboard_Fornecedores_Auth",
			url:  "http://localhost:8080/dashboard/fornecedores",
			desc: "Fornecedores com auth",
		},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			fmt.Printf("\n=== %s ===\n", endpoint.desc)

			// Criar request com Authorization header
			req, err := http.NewRequest("GET", endpoint.url, nil)
			if err != nil {
				t.Errorf("Erro ao criar request: %v", err)
				return
			}

			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			// Executar request
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("❌ Erro na requisição: %v\n", err)
				t.Errorf("Erro de conectividade: %v", err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Status: %d %s\n", resp.StatusCode, resp.Status)

			// Ler resposta
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("❌ Erro ao ler resposta: %v\n", err)
				return
			}

			// Analisar resposta
			if resp.StatusCode == 500 {
				fmt.Printf("❌ ERRO 500 DETECTADO!\n")
				fmt.Printf("Resposta completa:\n%s\n", string(body))
				
				// Tentar parse como JSON
				var errorData map[string]interface{}
				if json.Unmarshal(body, &errorData) == nil {
					fmt.Printf("\nErro estruturado:\n")
					for k, v := range errorData {
						fmt.Printf("  %s: %v\n", k, v)
					}
				}
				
				t.Errorf("❌ Erro 500 em %s", endpoint.url)
			} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				fmt.Printf("✅ Sucesso!\n")
				
				// Verificar se é JSON válido
				var responseData map[string]interface{}
				if json.Unmarshal(body, &responseData) == nil {
					fmt.Printf("Campos principais: %v\n", getTopLevelKeys(responseData))
				} else {
					fmt.Printf("⚠️ Resposta não é JSON: %s\n", string(body)[:100])
				}
			} else {
				fmt.Printf("⚠️ Status inesperado: %d\n", resp.StatusCode)
				fmt.Printf("Resposta: %s\n", string(body))
			}
		})
	}
}

// doLogin faz login e retorna o token JWT
func doLogin() (string, error) {
	// Dados de login (usando dados de exemplo)
	loginData := map[string]string{
		"email": "admin@mastercontrutora.com",
		"senha": "admin123",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	// Fazer request de login (rota correta)
	resp, err := http.Post("http://localhost:8080/usuarios/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login falhou: %d - %s", resp.StatusCode, string(body))
	}

	// Parse da resposta
	var loginResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return "", err
	}

	// Extrair token
	token, ok := loginResponse["token"].(string)
	if !ok {
		return "", fmt.Errorf("token não encontrado na resposta")
	}

	return token, nil
}

// getTopLevelKeys retorna as chaves principais do JSON
func getTopLevelKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}