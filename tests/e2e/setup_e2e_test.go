// file: tests/e2e/setup_e2e_test.go
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"
)

// --- Estruturas de Dados (Payloads e Respostas) ---

type UserPayload struct {
	Nome           string `json:"nome"`
	Email          string `json:"email"`
	Senha          string `json:"senha"`
	ConfirmarSenha string `json:"confirmarSenha"`
}

type LoginPayload struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type Funcionario struct {
	ID           string  `json:"id"` // Corrigido para minúsculo para corresponder à resposta da API
	Nome         string  `json:"nome"`
	CPF          string  `json:"cpf"`
	Cargo        string  `json:"cargo"`
	Departamento string  `json:"departamento"`
	ValorDiaria  float64 `json:"valorDiaria"`
}

type Fornecedor struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	CNPJ      string `json:"cnpj"`
	Categoria string `json:"categoria"`
	Contato   string `json:"contato"`
	Email     string `json:"email"`
}

type Material struct {
	ID              string `json:"id"`
	Nome            string `json:"nome"`
	Descricao       string `json:"descricao"`
	UnidadeDeMedida string `json:"unidadeDeMedida"`
	Categoria       string `json:"categoria"`
}

// ... (outras structs se necessário: Obra, Etapa, etc.)

// --- Configuração e Estado Global do Teste ---

const hostname = "http://localhost:8080"

var (
	// client será nosso cliente HTTP autenticado, com o cookie armazenado.
	client *http.Client

	// testData irá armazenar os IDs gerados durante a fase de setup.
	testData struct {
		FuncionarioID string
		FornecedorID  string
		MaterialID1   string
		MaterialID2   string
	}
)

// --- TestMain: O ponto de entrada para o setup e teardown do teste ---

func TestMain(m *testing.M) {
	log.Println(">>> INICIANDO SETUP GLOBAL DOS TESTES E2E...")
	if err := setup(); err != nil {
		log.Fatalf("Falha no setup global: %v", err)
	}

	// Roda todos os outros testes do pacote.
	code := m.Run()

	// Teardown pode ser adicionado aqui para limpar o ambiente, se necessário.
	log.Println(">>> TESTES E2E FINALIZADOS.")
	os.Exit(code)
}

// setup prepara o ambiente, registrando usuário, logando e criando dados base.
func setup() error {
	// Passo 1: Registrar usuário
	adminUser := UserPayload{
		Nome:           "Admin Teste E2E",
		Email:          fmt.Sprintf("admin_e2e_%d@construtora.com", time.Now().UnixNano()),
		Senha:          "senha_forte_123",
		ConfirmarSenha: "senha_forte_123",
	}
	// Usamos um cliente temporário para o registro, pois ainda não estamos autenticados.
	_, err := makeRequest(http.DefaultClient, "POST", "/usuarios/registrar", adminUser, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("falha ao registrar usuário: %w", err)
	}

	// Passo 2: Login para obter o cookie de autenticação
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("falha ao criar cookie jar: %w", err)
	}
	client = &http.Client{Jar: jar}

	loginPayload := LoginPayload{Email: adminUser.Email, Senha: adminUser.Senha}
	_, err = makeRequest(client, "POST", "/usuarios/login", loginPayload, http.StatusOK)
	if err != nil {
		return fmt.Errorf("falha ao fazer login: %w", err)
	}
	log.Println("Login realizado com sucesso, cookie JWT armazenado.")

	// Passo 3: Cadastrar entidades de base (usando o cliente já autenticado)
	err = criarEntidadesDeBase()
	if err != nil {
		return fmt.Errorf("falha ao criar entidades de base: %w", err)
	}

	log.Println("Setup global concluído com sucesso. Entidades de base criadas.")
	return nil
}

func criarEntidadesDeBase() error {
	var funcResp Funcionario
	payloadFunc := map[string]interface{}{"nome": "João da Silva (Eng. Chefe)", "cpf": "111.222.323-44", "cargo": "Engenheiro Civil", "departamento": "Construção Civil", "valorDiaria": 250.00}
	if err := doRequestAndDecode(client, "POST", "/funcionarios", payloadFunc, &funcResp, http.StatusCreated); err != nil {
		return fmt.Errorf("falha ao criar funcionário principal: %w", err)
	}
	testData.FuncionarioID = funcResp.ID

	var fornResp Fornecedor
	// Geramos um CNPJ dinâmico para cada execução do teste
	dynamicCNPJ := fmt.Sprintf("00.123.456/%04d-00", time.Now().UnixNano()%10000)
	payloadForn := map[string]interface{}{
		"nome":      "Casa do Construtor Center",
		"cnpj":      dynamicCNPJ, // Usando o CNPJ dinâmico
		"categoria": "Materiais Básicos",
		"contato":   "Carlos Andrade",                 // Adicionando campo que faltava no teste
		"email":     "comercial@casadoconstrutor.com", // Adicionando campo que faltava no teste
	}
	if err := doRequestAndDecode(client, "POST", "/fornecedores", payloadForn, &fornResp, http.StatusCreated); err != nil {
		return fmt.Errorf("falha ao criar fornecedor principal: %w", err)
	}
	testData.FornecedorID = fornResp.ID

	var mat1Resp Material
	payloadMat1 := map[string]interface{}{"nome": "Cimento Portland CP II 50kg", "unidadeDeMedida": "saco"}
	if err := doRequestAndDecode(client, "POST", "/materiais", payloadMat1, &mat1Resp, http.StatusCreated); err != nil {
		return fmt.Errorf("falha ao criar material 1: %w", err)
	}
	testData.MaterialID1 = mat1Resp.ID

	var mat2Resp Material
	payloadMat2 := map[string]interface{}{"nome": "Vergalhão de Aço CA-50 10mm", "unidadeDeMedida": "barra"}
	if err := doRequestAndDecode(client, "POST", "/materiais", payloadMat2, &mat2Resp, http.StatusCreated); err != nil {
		return fmt.Errorf("falha ao criar material 2: %w", err)
	}
	testData.MaterialID2 = mat2Resp.ID

	return nil
}

// --- Funções Auxiliares (Helpers) ---

// makeRequest cria e executa uma requisição HTTP, verificando o status da resposta.
func makeRequest(c *http.Client, method, path string, payload interface{}, expectedStatus int) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("falha ao fazer marshal do payload: %w", err)
		}
		body = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, hostname+path, body)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler corpo da resposta: %w", err)
	}

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("status inesperado para %s %s. esperado: %d, recebido: %d, corpo: %s",
			method, hostname+path, expectedStatus, resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// doRequestAndDecode é um wrapper sobre makeRequest para requisições que retornam um corpo JSON.
func doRequestAndDecode(c *http.Client, method, path string, payload interface{}, target interface{}, expectedStatus int) error {
	respBody, err := makeRequest(c, method, path, payload, expectedStatus)
	if err != nil {
		return err
	}
	if target != nil {
		if err := json.Unmarshal(respBody, target); err != nil {
			return fmt.Errorf("falha ao fazer unmarshal da resposta para o alvo: %w", err)
		}
	}
	return nil
}
