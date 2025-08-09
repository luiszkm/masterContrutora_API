# Módulo Financeiro

O módulo Financeiro é responsável pelo controle completo das movimentações financeiras da construtora, incluindo receitas de obras, pagamentos a fornecedores, funcionários e controle de fluxo de caixa.

## Visão Geral

### Funcionalidades Principais
- **Contas a Receber**: Gestão de receitas provenientes de obras e serviços
- **Contas a Pagar**: Controle de pagamentos a fornecedores e prestadores de serviços
- **Cronograma de Recebimentos**: Planejamento de receitas por etapas de obra
- **Fluxo de Caixa**: Visão consolidada de entradas e saídas financeiras
- **Integração por Eventos**: Criação automática de contas baseada em orçamentos aprovados

### Arquitetura

```
financeiro/
├── domain/
│   ├── conta_receber.go         # Entidade Contas a Receber
│   ├── conta_pagar.go           # Entidade Contas a Pagar
│   └── registro_pagamento.go    # Entidade Registro de Pagamento
├── service/
│   ├── conta_receber_service.go # Lógica de negócio - Contas a Receber
│   ├── conta_pagar_service.go   # Lógica de negócio - Contas a Pagar
│   └── service.go               # Service principal
├── handler/http/
│   ├── conta_receber_handler.go # Endpoints HTTP - Contas a Receber
│   ├── conta_pagar_handler.go   # Endpoints HTTP - Contas a Pagar
│   └── handler.go               # Handler principal
├── infrastructure/repository/postgres/
│   ├── conta_receber_repository.go
│   ├── conta_pagar_repository.go
│   └── payment_repository.go
└── events/
    └── handler.go               # Manipulador de eventos
```

## Entidades Principais

### 1. ContaReceber

```go
type ContaReceber struct {
    ID                        string
    ObraID                   *string
    CronogramaRecebimentoID  *string
    Cliente                  string
    TipoContaReceber        string  // OBRA, SERVICO, OUTROS
    Descricao               string
    ValorOriginal           float64
    ValorRecebido           float64
    DataVencimento          time.Time
    DataRecebimento         *time.Time
    Status                  string  // PENDENTE, RECEBIDO, VENCIDO, PARCIAL, CANCELADO
    FormaPagamento          *string
    Observacoes             *string
    NumeroDocumento         *string
    CreatedAt               time.Time
    UpdatedAt               time.Time
}
```

**Métodos de Negócio:**
- `ValorSaldo() float64`: Calcula valor restante a receber
- `PercentualRecebido() float64`: Calcula percentual já recebido
- `EstaVencido() bool`: Verifica se a conta está vencida
- `DiasVencimento() int`: Calcula dias de vencimento
- `RegistrarRecebimento(valor, formaPagamento, observacoes)`: Registra recebimento

### 2. ContaPagar

```go
type ContaPagar struct {
    ID              string
    FornecedorID    *string
    ObraID          *string
    OrcamentoID     *string
    FornecedorNome  string
    TipoContaPagar  string  // FORNECEDOR, SERVICO, MATERIAL, OUTROS
    Descricao       string
    ValorOriginal   float64
    ValorPago       float64
    DataVencimento  time.Time
    DataPagamento   *time.Time
    Status          string  // PENDENTE, PAGO, VENCIDO, PARCIAL, CANCELADO
    FormaPagamento  *string
    Observacoes     *string
    NumeroDocumento *string
    NumeroCompraNF  *string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

**Métodos de Negócio:**
- `ValorSaldo() float64`: Calcula valor restante a pagar
- `PercentualPago() float64`: Calcula percentual já pago
- `EstaVencido() bool`: Verifica se a conta está vencida
- `DiasVencimento() int`: Calcula dias de vencimento
- `RegistrarPagamento(valor, formaPagamento, observacoes)`: Registra pagamento
- `MarcarComoVencido()`: Marca conta como vencida

### 3. ParcelaContaPagar

```go
type ParcelaContaPagar struct {
    ID             string
    ContaPagarID   string
    NumeroParcela  int
    ValorParcela   float64
    DataVencimento time.Time
    DataPagamento  *time.Time
    ValorPago      float64
    Status         string
    FormaPagamento *string
    Observacoes    *string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

## APIs Disponíveis

### Contas a Receber

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/contas-receber` | Criar nova conta a receber |
| GET | `/contas-receber` | Listar contas com paginação |
| GET | `/contas-receber/{id}` | Buscar conta por ID |
| POST | `/contas-receber/{id}/recebimentos` | Registrar recebimento |
| GET | `/contas-receber/vencidas` | Listar contas vencidas |
| GET | `/contas-receber/resumo` | Obter resumo financeiro |
| GET | `/obras/{id}/contas-receber` | Listar contas de uma obra |

### Contas a Pagar

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/contas-pagar` | Criar nova conta a pagar |
| POST | `/contas-pagar/orcamentos` | Criar conta a partir de orçamento |
| GET | `/contas-pagar` | Listar contas com paginação |
| GET | `/contas-pagar/{id}` | Buscar conta por ID |
| POST | `/contas-pagar/{id}/pagamentos` | Registrar pagamento |
| GET | `/contas-pagar/vencidas` | Listar contas vencidas |
| GET | `/contas-pagar/resumo` | Obter resumo financeiro |
| GET | `/obras/{id}/contas-pagar` | Listar contas de uma obra |
| GET | `/fornecedores/{id}/contas-pagar` | Listar contas de um fornecedor |

### Cronograma de Recebimentos

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/cronograma-recebimentos` | Criar cronograma individual |
| POST | `/cronograma-recebimentos/lote` | Criar cronograma em lote |
| GET | `/cronograma-recebimentos/{id}` | Buscar cronograma por ID |
| POST | `/cronograma-recebimentos/{id}/recebimentos` | Registrar recebimento |
| GET | `/obras/{id}/cronograma-recebimentos` | Listar cronogramas de uma obra |

## Exemplos de Uso

### Criar Conta a Receber

```http
POST /contas-receber
Content-Type: application/json

{
  "obraId": "obra-uuid",
  "cliente": "Cliente ABC Ltda",
  "tipoContaReceber": "OBRA",
  "descricao": "Primeira parcela da obra",
  "valorOriginal": 50000.00,
  "dataVencimento": "2025-02-15T00:00:00Z",
  "numeroDocumento": "NF-001/2025"
}
```

### Registrar Recebimento

```http
POST /contas-receber/{conta-id}/recebimentos
Content-Type: application/json

{
  "valor": 50000.00,
  "formaPagamento": "PIX",
  "contaBancariaId": "conta-bancaria-uuid",
  "observacoes": "Pagamento integral da primeira parcela"
}
```

### Criar Conta a Pagar

```http
POST /contas-pagar
Content-Type: application/json

{
  "fornecedorId": "fornecedor-uuid",
  "obraId": "obra-uuid",
  "fornecedorNome": "Fornecedor XYZ Ltda",
  "tipoContaPagar": "MATERIAL",
  "descricao": "Materiais de construção - Cimento e areia",
  "valorOriginal": 15000.00,
  "dataVencimento": "2025-03-15T00:00:00Z",
  "numeroDocumento": "NF-456/2025",
  "numeroCompraNf": "COMP-123"
}
```

### Criar Conta a Pagar de Orçamento

```http
POST /contas-pagar/orcamentos
Content-Type: application/json

{
  "orcamentoId": "orcamento-uuid",
  "dataVencimento": "2025-03-15T00:00:00Z",
  "numeroCompraNf": "COMP-456",
  "observacoes": "Conta gerada automaticamente do orçamento aprovado",
  "dividirParcelas": true,
  "quantidadeParcelas": 3
}
```

## Fluxo de Caixa

### Cálculo de Entradas
- **Receitas Reais**: Valor recebido de contas a receber com status `RECEBIDO`
- **Base de Dados**: Tabela `contas_receber` onde `data_recebimento IS NOT NULL`

### Cálculo de Saídas
- **Pagamentos a Fornecedores**: Valor pago de contas a pagar com status `PAGO`
- **Pagamentos a Funcionários**: Registros de pagamento do módulo Pessoal
- **Base de Dados**: 
  - Tabela `contas_pagar` onde `data_pagamento IS NOT NULL`
  - Tabela `registros_pagamento` para funcionários

### API do Fluxo de Caixa

```http
GET /dashboard/fluxo-caixa?dataInicio=2025-01-01&dataFim=2025-12-31
```

**Resposta:**
```json
{
  "totalEntradas": 150000.00,
  "totalSaidas": 85000.00,
  "saldoAtual": 65000.00,
  "fluxoPorPeriodo": [
    {
      "periodo": "2025-01-31T21:00:00-03:00",
      "entradas": 30000.00,
      "saidas": 15000.00,
      "saldoLiquido": 15000.00
    }
  ],
  "tendenciaMensal": "crescimento"
}
```

## Integração por Eventos

### Eventos Publicados

1. **MovimentacaoFinanceiraRegistrada**
   - Publicado quando: Recebimento ou pagamento é registrado
   - Payload: Detalhes da movimentação financeira

### Eventos Consumidos

1. **OrcamentoStatusAtualizado**
   - Quando: Orçamento é aprovado
   - Ação: Cria automaticamente conta a pagar

2. **CronogramaRecebimentoCriado**
   - Quando: Cronograma de recebimento é criado
   - Ação: Cria contas a receber para cada etapa

3. **PagamentoApontamentoRealizado**
   - Quando: Funcionário recebe pagamento
   - Ação: Registra saída no fluxo de caixa

## Regras de Negócio

### Contas a Receber
- Valor recebido não pode exceder valor original
- Status muda automaticamente baseado no valor recebido:
  - `PARCIAL`: 0 < valor_recebido < valor_original
  - `RECEBIDO`: valor_recebido = valor_original
- Contas são marcadas como `VENCIDO` após data de vencimento

### Contas a Pagar
- Valor pago não pode exceder valor original
- Status muda automaticamente baseado no valor pago:
  - `PARCIAL`: 0 < valor_pago < valor_original
  - `PAGO`: valor_pago = valor_original
- Contas são marcadas como `VENCIDO` após data de vencimento
- Suporte a parcelamento com tabela `parcelas_conta_pagar`

### Cronograma de Recebimentos
- Uma obra não pode ter etapas duplicadas (constraint única)
- Valor recebido não pode exceder valor previsto
- Status atualizado automaticamente conforme recebimentos

## Validações

### Dados Obrigatórios
- **ContaReceber**: Cliente, descrição, valor original > 0, data vencimento
- **ContaPagar**: Fornecedor nome, descrição, valor original > 0, data vencimento

### Regras de Validação
- Valores monetários devem ser positivos
- Datas de vencimento devem ser futuras (para novas contas)
- Status devem estar dentro dos valores permitidos
- Referências a outras entidades (obra, fornecedor) devem existir

## Monitoramento e Logs

### Métricas Importantes
- Total de contas a receber pendentes
- Total de contas a pagar vencidas  
- Valor do fluxo de caixa mensal
- Tempo médio de recebimento
- Inadimplência por cliente/obra

### Logs Estruturados
```json
{
  "timestamp": "2025-08-08T13:15:00Z",
  "level": "INFO",
  "service": "ContaReceberService",
  "operation": "RegistrarRecebimento",
  "conta_id": "uuid",
  "valor": 30000.00,
  "status": "RECEBIDO"
}
```

## Próximas Implementações

### Fase 3 - Funcionalidades Avançadas
- Conciliação bancária
- Projeções de fluxo de caixa
- Relatórios financeiros avançados
- Dashboard com gráficos
- Notificações de vencimento
- Integração com bancos (Open Banking)

### Fase 4 - Recursos Empresariais
- Centros de custo por obra
- Orçamento vs. realizado
- Análise de rentabilidade por obra
- Controle de impostos
- Integração contábil
- Auditoria e compliance