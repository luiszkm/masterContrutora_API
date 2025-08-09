# Módulo Obras

O módulo Obras é responsável pela gestão completa de projetos de construção, incluindo controle financeiro, cronogramas de recebimento, alocação de recursos e acompanhamento de progresso.

## Visão Geral

### Funcionalidades Principais
- **Gestão de Obras**: CRUD completo de projetos de construção
- **Controle Financeiro**: Valores contratuais, recebimentos e saldos
- **Cronogramas de Recebimento**: Planejamento de receitas por etapas
- **Gestão de Etapas**: Controle de progresso e marcos da obra
- **Alocação de Recursos**: Designação de funcionários para obras
- **Integração Financeira**: Comunicação automática com módulo Financeiro

### Arquitetura

```
obras/
├── domain/
│   ├── obra.go                    # Entidade principal da Obra
│   ├── cronograma_recebimento.go  # Cronograma de recebimentos
│   ├── etapa.go                   # Etapas da obra
│   ├── etapa_padrao.go           # Templates de etapas
│   └── alocacao.go               # Alocação de funcionários
├── service/
│   ├── obras_service.go           # Lógica de negócio principal
│   ├── cronograma_service.go      # Gestão de cronogramas
│   └── dto/                       # Data Transfer Objects
├── handler/http/
│   ├── obras_handler.go           # Endpoints HTTP principais
│   ├── cronograma_handler.go      # Endpoints de cronogramas
│   └── etapas_handler.go         # Endpoints de etapas
├── infrastructure/repository/postgres/
│   ├── obras_repository.go
│   ├── cronograma_recebimento_repository.go
│   ├── etapas_repository.go
│   └── alocacao_repository.go
└── events/
    └── handler.go                # Manipulador de eventos
```

## Entidades Principais

### 1. Obra

```go
type Obra struct {
    ID                     string
    Nome                   string
    Cliente                string
    Endereco               string
    Descricao              string
    DataInicio             time.Time
    DataFim                *time.Time
    Status                 string    // "Em Andamento", "Concluída", "Pausada", "Cancelada"
    
    // Campos Financeiros
    ValorContratoTotal     float64   // Valor total do contrato
    ValorRecebido          float64   // Valor já recebido
    TipoCobranca          string    // "VISTA", "PARCELADO", "ETAPAS"
    DataAssinaturaContrato *time.Time // Data da assinatura do contrato
    
    CreatedAt              time.Time
    UpdatedAt              time.Time
}
```

**Métodos de Negócio:**
- `ValorSaldo() float64`: Calcula valor restante a receber
- `PercentualRecebido() float64`: Calcula percentual já recebido do contrato
- `RegistrarRecebimento(valor, observacoes)`: Registra recebimento parcial
- `PodeIniciar() bool`: Verifica se obra pode ser iniciada
- `PodeConcluir() bool`: Verifica se obra pode ser concluída

### 2. CronogramaRecebimento

```go
type CronogramaRecebimento struct {
    ID                     string
    ObraID                 string
    NumeroEtapa           int
    DescricaoEtapa        string
    ValorPrevisto         float64
    DataVencimento        time.Time
    Status                string    // "PENDENTE", "RECEBIDO", "VENCIDO", "PARCIAL"
    DataRecebimento       *time.Time
    ValorRecebido         float64
    ObservacoesRecebimento *string
    CreatedAt             time.Time
    UpdatedAt             time.Time
}
```

**Métodos de Negócio:**
- `ValorSaldo() float64`: Valor restante a receber nesta etapa
- `PercentualRecebido() float64`: Percentual recebido da etapa
- `EstaVencido() bool`: Verifica se etapa está vencida
- `PodeMarcarComoRecebido() bool`: Verifica se pode ser marcada como recebida

### 3. Etapa

```go
type Etapa struct {
    ID          string
    ObraID      string
    Nome        string
    Descricao   string
    DataInicio  time.Time
    DataFim     *time.Time
    Status      string    // "Pendente", "Em Andamento", "Concluída"
    Ordem       int       // Ordem sequencial na obra
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 4. Alocacao

```go
type Alocacao struct {
    ID           string
    ObraID       string
    FuncionarioID string
    DataInicio   time.Time
    DataFim      *time.Time
    Status       string    // "Ativa", "Finalizada"
    Observacoes  *string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

## APIs Disponíveis

### Obras

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/obras` | Criar nova obra |
| GET | `/obras` | Listar obras com filtros |
| GET | `/obras/{id}` | Buscar obra por ID |
| PUT | `/obras/{id}` | Atualizar obra |
| DELETE | `/obras/{id}` | Excluir obra (soft delete) |
| POST | `/obras/{id}/recebimentos` | Registrar recebimento |
| GET | `/obras/{id}/financeiro` | Obter resumo financeiro da obra |

### Cronogramas de Recebimento

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/cronograma-recebimentos` | Criar cronograma individual |
| POST | `/cronograma-recebimentos/lote` | Criar múltiplos cronogramas |
| GET | `/obras/{id}/cronograma-recebimentos` | Listar cronogramas da obra |
| GET | `/cronograma-recebimentos/{id}` | Buscar cronograma por ID |
| PUT | `/cronograma-recebimentos/{id}` | Atualizar cronograma |
| POST | `/cronograma-recebimentos/{id}/recebimentos` | Registrar recebimento |

### Etapas

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/obras/{id}/etapas` | Criar etapa na obra |
| GET | `/obras/{id}/etapas` | Listar etapas da obra |
| PUT | `/etapas/{id}` | Atualizar etapa |
| DELETE | `/etapas/{id}` | Excluir etapa |
| PATCH | `/etapas/{id}/status` | Alterar status da etapa |

### Alocações

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/obras/{id}/alocacoes` | Alocar funcionário à obra |
| GET | `/obras/{id}/alocacoes` | Listar alocações da obra |
| PUT | `/alocacoes/{id}` | Atualizar alocação |
| DELETE | `/alocacoes/{id}` | Finalizar alocação |

## Exemplos de Uso

### Criar Obra com Controle Financeiro

```http
POST /obras
Content-Type: application/json

{
  "nome": "Casa Residencial - João Silva",
  "cliente": "João Silva",
  "endereco": "Rua das Flores, 123 - Centro",
  "dataInicio": "2025-02-01",
  "descricao": "Casa de 150m² com 3 quartos",
  "valorContratoTotal": 200000.00,
  "tipoCobranca": "ETAPAS",
  "dataAssinaturaContrato": "2025-01-15T00:00:00Z"
}
```

### Criar Cronograma de Recebimento em Lote

```http
POST /cronograma-recebimentos/lote
Content-Type: application/json

{
  "obraId": "obra-uuid",
  "substituirExistente": true,
  "cronogramas": [
    {
      "numeroEtapa": 1,
      "descricaoEtapa": "Fundação e estrutura",
      "valorPrevisto": 80000.00,
      "dataVencimento": "2025-03-15T00:00:00Z"
    },
    {
      "numeroEtapa": 2,
      "descricaoEtapa": "Alvenaria e cobertura",
      "valorPrevisto": 70000.00,
      "dataVencimento": "2025-05-15T00:00:00Z"
    },
    {
      "numeroEtapa": 3,
      "descricaoEtapa": "Acabamentos finais",
      "valorPrevisto": 50000.00,
      "dataVencimento": "2025-07-15T00:00:00Z"
    }
  ]
}
```

### Registrar Recebimento de Etapa

```http
POST /cronograma-recebimentos/{cronograma-id}/recebimentos
Content-Type: application/json

{
  "valor": 80000.00,
  "observacoes": "Recebimento integral da primeira etapa - PIX"
}
```

### Alocar Funcionário à Obra

```http
POST /obras/{obra-id}/alocacoes
Content-Type: application/json

{
  "funcionarioId": "funcionario-uuid",
  "dataInicio": "2025-02-01T08:00:00Z",
  "observacoes": "Responsável pela supervisão geral da obra"
}
```

## Integrações e Eventos

### Eventos Publicados

1. **CronogramaRecebimentoCriado**
   - Publicado quando: Cronograma de recebimento é criado
   - Payload: Detalhes do cronograma e da obra
   - Consumido por: Módulo Financeiro (cria contas a receber)

2. **RecebimentoRealizado**
   - Publicado quando: Recebimento é registrado
   - Payload: Valor recebido e detalhes do pagamento
   - Consumido por: Módulo Financeiro (atualiza fluxo de caixa)

### Eventos Consumidos

1. **OrcamentoStatusAtualizado**
   - Quando: Orçamento de materiais é aprovado
   - Ação: Atualiza métricas financeiras da obra

## Regras de Negócio

### Obras
- Valor do contrato deve ser positivo
- Data de início não pode ser no passado (para novas obras)
- Cliente é obrigatório
- Status deve ser válido: "Em Andamento", "Concluída", "Pausada", "Cancelada"
- Tipo de cobrança deve ser: "VISTA", "PARCELADO", "ETAPAS"

### Cronogramas de Recebimento
- Uma obra não pode ter etapas com números duplicados
- Valor previsto deve ser positivo
- Valor recebido não pode exceder valor previsto
- Status muda automaticamente baseado nos recebimentos:
  - `PARCIAL`: 0 < valor_recebido < valor_previsto
  - `RECEBIDO`: valor_recebido = valor_previsto
- Etapas são marcadas como `VENCIDO` após data de vencimento

### Etapas
- Ordem deve ser sequencial e única por obra
- Status deve seguir fluxo lógico: Pendente → Em Andamento → Concluída
- Data de fim deve ser posterior à data de início

### Alocações
- Funcionário não pode estar alocado em duas obras no mesmo período
- Data de fim deve ser posterior à data de início
- Status válidos: "Ativa", "Finalizada"

## Dashboard e Métricas

### Métricas por Obra
```http
GET /obras/{obra-id}/financeiro
```

**Resposta:**
```json
{
  "obra": {
    "id": "obra-uuid",
    "nome": "Casa João Silva",
    "valorContrato": 200000.00,
    "valorRecebido": 80000.00,
    "valorSaldo": 120000.00,
    "percentualRecebido": 40.0
  },
  "cronogramas": [
    {
      "etapa": 1,
      "descricao": "Fundação",
      "valorPrevisto": 80000.00,
      "valorRecebido": 80000.00,
      "status": "RECEBIDO",
      "dataVencimento": "2025-03-15T00:00:00Z"
    }
  ],
  "proximosVencimentos": [
    {
      "etapa": 2,
      "valorPrevisto": 70000.00,
      "dataVencimento": "2025-05-15T00:00:00Z",
      "diasParaVencimento": 45
    }
  ]
}
```

### Métricas Gerais do Dashboard
- Total de obras ativas
- Valor total em contratos
- Valor total recebido
- Valor em aberto
- Obras por status
- Cronograma de recebimentos futuros

## Validações e Constraints

### Validações de Entrada
- **Nome da obra**: Obrigatório, máximo 200 caracteres
- **Cliente**: Obrigatório, máximo 100 caracteres
- **Valor do contrato**: Positivo, máximo 2 casas decimais
- **Datas**: Formato ISO 8601, validações de lógica temporal

### Constraints de Banco
- Chaves primárias UUID
- Relacionamentos com foreign keys
- Índices para performance:
  - `idx_obras_cliente` - Busca por cliente
  - `idx_obras_status` - Filtro por status
  - `idx_cronograma_obra_etapa` - Unique constraint (obra_id, numero_etapa)
  - `idx_alocacao_funcionario_data` - Verificação de sobreposição

## Monitoramento e Logs

### Métricas Importantes
- Número de obras criadas por período
- Tempo médio de conclusão de obras
- Valor médio de contratos
- Taxa de recebimento (valor recebido / valor contratado)
- Obras com atraso no cronograma

### Logs Estruturados
```json
{
  "timestamp": "2025-08-08T14:00:00Z",
  "level": "INFO",
  "service": "ObrasService",
  "operation": "CriarObra",
  "obra_id": "uuid",
  "cliente": "João Silva",
  "valor_contrato": 200000.00
}
```

## Relatórios Disponíveis

### Relatório Financeiro por Obra
- Valor contratado vs. valor recebido
- Cronograma de recebimentos (planejado vs. real)
- Análise de atrasos nos recebimentos
- Projeção de fluxo de caixa

### Relatório de Progresso
- Status das etapas por obra
- Tempo gasto vs. planejado
- Alocação de recursos
- Marcos importantes

### Relatório Gerencial
- Obras mais rentáveis
- Clientes com maior volume
- Performance por período
- Indicadores de eficiência

## Próximas Funcionalidades

### Curto Prazo
- Fotos de progresso da obra
- Notificações de vencimento de cronograma
- Templates de obras por tipo
- Integração com GPS para localização

### Médio Prazo
- Gestão de documentos (contratos, alvarás)
- Cronograma físico-financeiro
- Controle de qualidade
- Avaliação de fornecedores

### Longo Prazo
- Análise preditiva de custos
- Integração com drones para monitoramento
- BIM (Building Information Modeling)
- Sustentabilidade e certificações verdes