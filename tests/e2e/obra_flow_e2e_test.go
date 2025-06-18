// file: tests/e2e/obra_flow_e2e_test.go
package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// ObraDashboard é a struct para decodificar a resposta do endpoint de dashboard.
type ObraDashboard struct {
	OrcamentoTotalAprovado float64 `json:"orcamentoTotalAprovado"`
	FuncionariosAlocados   int     `json:"funcionariosAlocados"`
}

func TestObraCompletaFlow(t *testing.T) {
	// O setup já foi executado pelo TestMain.
	// O 'client' já está autenticado e 'testData' está populado.

	var obraID string
	var etapaID string
	var orcamentoID string
	var valorTotalOrcamento float64

	t.Run("Deve criar uma obra principal", func(t *testing.T) {
		type ObraResponse struct {
			ID string `json:"id"`
		}
		var resp ObraResponse

		payload := map[string]interface{}{
			"nome":       "Obra Residencial Alphaville",
			"cliente":    "Família Silva",
			"endereco":   "Av. dos Testes, 456",
			"dataInicio": time.Now().Format("2006-01-02"),
		}

		err := doRequestAndDecode(client, "POST", "/obras", payload, &resp, http.StatusCreated)
		if err != nil {
			t.Fatalf("Falha ao criar obra: %v", err)
		}
		if resp.ID == "" {
			t.Fatal("Falha ao criar obra: ID retornou vazio")
		}
		obraID = resp.ID
	})

	t.Run("Deve adicionar uma etapa à obra", func(t *testing.T) {
		if obraID == "" {
			t.Skip("Pulando: ID da obra não está disponível")
		}

		type EtapaResponse struct {
			ID string `json:"id"`
		}
		var resp EtapaResponse

		url := fmt.Sprintf("/obras/%s/etapas", obraID)
		payload := map[string]interface{}{
			"nome":               "Fundações e Estrutura",
			"dataInicioPrevista": time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			"dataFimPrevista":    time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
		}

		err := doRequestAndDecode(client, "POST", url, payload, &resp, http.StatusCreated)
		if err != nil {
			t.Fatalf("Falha ao adicionar etapa: %v", err)
		}
		if resp.ID == "" {
			t.Fatal("Falha ao adicionar etapa: ID da etapa retornou vazio")
		}
		etapaID = resp.ID
	})

	t.Run("Deve alocar funcionário na obra", func(t *testing.T) {
		if obraID == "" || testData.FuncionarioID == "" {
			t.Skip("Pulando: ID da obra ou do funcionário não está disponível")
		}

		url := fmt.Sprintf("/obras/%s/alocacoes", obraID)
		payload := map[string]interface{}{
			"funcionarioId":      testData.FuncionarioID,
			"dataInicioAlocacao": time.Now().Format("2006-01-02"),
		}
		// Para esta requisição, não precisamos decodificar a resposta, apenas verificar o status.
		err := doRequestAndDecode(client, "POST", url, payload, nil, http.StatusCreated)
		if err != nil {
			t.Fatalf("Falha ao alocar funcionário: %v", err)
		}
	})

	t.Run("Deve criar e aprovar um orçamento para a etapa", func(t *testing.T) {
		if etapaID == "" {
			t.Skip("Pulando: ID da etapa não está disponível")
		}

		// 1. Criar o orçamento
		type OrcamentoResponse struct {
			ID string `json:"id"`
		}
		var resp OrcamentoResponse

		urlCriar := fmt.Sprintf("/etapas/%s/orcamentos", etapaID)
		valorUnitarioCimento := 30.50
		valorUnitarioAco := 60.00
		quantidadeCimento := 100.0
		quantidadeAco := 50.0
		valorTotalOrcamento = (valorUnitarioCimento * quantidadeCimento) + (valorUnitarioAco * quantidadeAco)

		payloadCriar := map[string]interface{}{
			"numero":       fmt.Sprintf("ORC-E2E-%d", time.Now().UnixNano()),
			"fornecedorId": testData.FornecedorID,
			"itens": []map[string]interface{}{
				{"materialId": testData.MaterialID1, "quantidade": quantidadeCimento, "valorUnitario": valorUnitarioCimento},
				{"materialId": testData.MaterialID2, "quantidade": quantidadeAco, "valorUnitario": valorUnitarioAco},
			},
		}
		err := doRequestAndDecode(client, "POST", urlCriar, payloadCriar, &resp, http.StatusCreated)
		if err != nil {
			t.Fatalf("Falha ao criar orçamento: %v", err)
		}
		orcamentoID = resp.ID

		// 2. Aprovar o orçamento
		urlAprovar := fmt.Sprintf("/orcamentos/%s", orcamentoID)
		payloadAprovar := map[string]string{"status": "Aprovado"}
		err = doRequestAndDecode(client, "PATCH", urlAprovar, payloadAprovar, nil, http.StatusOK)
		if err != nil {
			t.Fatalf("Falha ao aprovar orçamento: %v", err)
		}
	})

	t.Run("Deve buscar o dashboard da obra com os dados atualizados", func(t *testing.T) {
		if obraID == "" {
			t.Skip("Pulando: ID da obra não está disponível")
		}

		var dashboard ObraDashboard
		url := fmt.Sprintf("/obras/%s", obraID)

		// Pode levar um instante para o evento do orçamento ser processado e o dashboard ser atualizado.
		// Em um cenário real, poderíamos usar uma estratégia de retentativa aqui.
		// Por simplicidade, vamos apenas esperar um pouco.
		time.Sleep(200 * time.Millisecond)

		err := doRequestAndDecode(client, "GET", url, nil, &dashboard, http.StatusOK)
		if err != nil {
			t.Fatalf("Falha ao buscar dashboard da obra: %v", err)
		}

		// Verificações finais
		if dashboard.FuncionariosAlocados != 1 {
			t.Errorf("Número de funcionários alocados inesperado. esperado: 1, recebido: %d", dashboard.FuncionariosAlocados)
		}

		if dashboard.OrcamentoTotalAprovado != valorTotalOrcamento {
			t.Errorf("Valor do orçamento total aprovado inesperado. esperado: %f, recebido: %f", valorTotalOrcamento, dashboard.OrcamentoTotalAprovado)
		}
	})
}
