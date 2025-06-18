// file: tests/e2e/funcionario_e2e_test.go
package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestFuncionarioEndpoints(t *testing.T) {
	// O client autenticado já foi preparado pelo setup em 'setup_e2e_test.go'
	// As structs de payload/response também já estão definidas lá.

	var funcionarioID string  // Armazenará o ID do funcionário criado neste teste
	var funcionarioCPF string // Armazenará o CPF para evitar colisões

	// Corresponde à sua requisição: POST /funcionarios (@name CadastrarFuncionarioPrincipal)
	t.Run("Deve Criar um novo funcionário com sucesso", func(t *testing.T) {
		funcionarioCPF = fmt.Sprintf("111.222.313-%02d", time.Now().UnixNano()%100)
		payload := map[string]interface{}{
			"nome":         "João da Silva (Eng. Chefe)",
			"cpf":          funcionarioCPF,
			"cargo":        "Engenheiro Civil",
			"departamento": "Construção Civil",
			"valorDiaria":  100.00,
		}

		var resp Funcionario
		err := doRequestAndDecode(client, "POST", "/funcionarios", payload, &resp, http.StatusCreated)
		if err != nil {
			t.Fatalf("Falha ao criar funcionário: %v", err)
		}

		if resp.ID == "" {
			t.Fatal("Falha ao criar funcionário: ID retornou vazio")
		}
		funcionarioID = resp.ID // Equivalente ao seu: client.global.set("funcionarioId", ...);
	})

	// Corresponde à sua requisição: GET /funcionarios/{id}
	t.Run("Deve Buscar o funcionário recém-criado pelo ID", func(t *testing.T) {
		if funcionarioID == "" {
			t.Skip("Pulando: ID do funcionário não está disponível")
		}

		url := fmt.Sprintf("/funcionarios/%s", funcionarioID)
		var resp Funcionario
		err := doRequestAndDecode(client, "GET", url, nil, &resp, http.StatusOK)
		if err != nil {
			t.Fatalf("Falha ao buscar funcionário pelo ID: %v", err)
		}

		if resp.CPF != funcionarioCPF {
			t.Errorf("CPF do funcionário buscado não corresponde. esperado: %s, recebido: %s", funcionarioCPF, resp.CPF)
		}
	})

	// Corresponde à sua requisição: GET /funcionarios (@name BuscarFuncionarios)
	t.Run("Deve Listar todos os funcionários e encontrar o recém-criado", func(t *testing.T) {
		if funcionarioID == "" {
			t.Skip("Pulando: ID do funcionário não está disponível")
		}

		var resp []Funcionario
		err := doRequestAndDecode(client, "GET", "/funcionarios", nil, &resp, http.StatusOK)
		if err != nil {
			t.Fatalf("Falha ao listar funcionários: %v", err)
		}

		encontrado := false
		for _, f := range resp {
			if f.ID == funcionarioID {
				encontrado = true
				break
			}
		}

		if !encontrado {
			t.Error("Não foi possível encontrar o funcionário recém-criado na lista de funcionários")
		}
	})

	// Corresponde à sua requisição: PUT /funcionarios/{id} (@name Atualizar Funcionário)
	t.Run("Deve Atualizar os dados do funcionário", func(t *testing.T) {
		if funcionarioID == "" {
			t.Skip("Pulando: ID do funcionário não está disponível")
		}

		url := fmt.Sprintf("/funcionarios/%s", funcionarioID)
		payload := map[string]interface{}{
			"nome":   "João da Silva (Eng. Chefe Atualizado)",
			"diaria": 120.00,
		}

		var resp Funcionario
		err := doRequestAndDecode(client, "PUT", url, payload, &resp, http.StatusOK)
		if err != nil {
			t.Fatalf("Falha ao atualizar funcionário: %v", err)
		}

		if resp.Nome != "João da Silva (Eng. Chefe Atualizado)" {
			t.Errorf("Nome não foi atualizado corretamente. esperado: 'João da Silva (Eng. Chefe Atualizado)', recebido: '%s'", resp.Nome)
		}
	})

	// Corresponde à sua requisição: DELETE /funcionarios/{id} (@name DeletarFuncionario)
	t.Run("Deve Deletar o funcionário", func(t *testing.T) {
		if funcionarioID == "" {
			t.Skip("Pulando: ID do funcionário não está disponível")
		}

		url := fmt.Sprintf("/funcionarios/%s", funcionarioID)
		err := doRequestAndDecode(client, "DELETE", url, nil, nil, http.StatusNoContent)
		if err != nil {
			t.Fatalf("Falha ao deletar funcionário: %v", err)
		}
	})

	// Este é um passo extra de verificação que é uma boa prática em testes automatizados.
	t.Run("Deve falhar ao tentar buscar o funcionário deletado", func(t *testing.T) {
		if funcionarioID == "" {
			t.Skip("Pulando: ID do funcionário não está disponível")
		}

		url := fmt.Sprintf("/funcionarios/%s", funcionarioID)
		// Esperamos um 404 Not Found, pois o funcionário não deve mais ser encontrado
		err := doRequestAndDecode(client, "GET", url, nil, nil, http.StatusNotFound)
		if err != nil {
			t.Fatalf("A busca pelo funcionário deletado não retornou o status esperado: %v", err)
		}
	})
}
