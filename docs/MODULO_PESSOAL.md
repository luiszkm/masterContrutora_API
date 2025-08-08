# Módulo Pessoal

O módulo Pessoal é responsável pela gestão completa de recursos humanos da construtora, incluindo cadastro de funcionários, controle de apontamentos de horas, aprovação de pagamentos e integração com o sistema financeiro.

## Visão Geral

### Funcionalidades Principais
- **Gestão de Funcionários**: CRUD completo de colaboradores
- **Apontamentos Quinzenais**: Registro de horas trabalhadas por período
- **Controle de Pagamentos**: Aprovação e processamento de pagamentos
- **Alocação de Recursos**: Designação de funcionários para obras
- **Integração Financeira**: Comunicação automática com módulo Financeiro

### Arquitetura

```
pessoal/
├── domain/
│   ├── funcionario.go             # Entidade Funcionário
│   └── apontamento_quinzenal.go   # Apontamentos de horas
├── service/
│   ├── pessoal_service.go         # Lógica de negócio principal
│   └── dto/                       # Data Transfer Objects
├── handler/http/
│   └── pessoal_handler.go         # Endpoints HTTP
├── infrastructure/repository/postgres/
│   ├── funcionario_repository.go
│   └── apontamento_repository.go
└── events/
    └── handler.go                 # Manipulador de eventos
```

## Entidades Principais

### 1. Funcionario

```go
type Funcionario struct {
    ID          string
    Nome        string
    Email       string
    Telefone    string
    CPF         string
    Cargo       string
    Salario     float64
    DataAdmissao time.Time
    DataDemissao *time.Time
    Status      string    // "Ativo", "Inativo", "Férias", "Licença"
    Endereco    string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Métodos de Negócio:**
- `EstaAtivo() bool`: Verifica se funcionário está ativo
- `PodeReceberPagamento() bool`: Verifica se pode receber pagamentos
- `TempoServico() time.Duration`: Calcula tempo de serviço
- `SalarioProporcional(dias int) float64`: Calcula salário proporcional

### 2. ApontamentoQuinzenal

```go
type ApontamentoQuinzenal struct {
    ID            string
    FuncionarioID string
    ObraID        *string
    DataInicio    time.Time
    DataFim       time.Time
    HorasTrabalhadas int
    ValorHora     float64
    ValorTotal    float64
    ValorCalculado float64
    Status        string    // "Pendente", "Aprovado", "Rejeitado", "Pago"
    Observacoes   *string
    DataAprovacao *time.Time
    DataPagamento *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

**Métodos de Negócio:**
- `CalcularValorTotal()`: Calcula valor total baseado em horas × valor por hora
- `PodeSerAprovado() bool`: Verifica se apontamento pode ser aprovado
- `PodeSerPago() bool`: Verifica se apontamento pode ser pago
- `EstaVencido() bool`: Verifica se apontamento está atrasado
- `Aprovar(observacoes)`: Aprova o apontamento
- `Rejeitar(motivo)`: Rejeita o apontamento
- `RegistrarPagamento(valor, observacoes)`: Registra pagamento

## APIs Disponíveis

### Funcionários

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/funcionarios` | Criar novo funcionário |
| GET | `/funcionarios` | Listar funcionários com filtros |
| GET | `/funcionarios/{id}` | Buscar funcionário por ID |
| PUT | `/funcionarios/{id}` | Atualizar funcionário |
| DELETE | `/funcionarios/{id}` | Inativar funcionário (soft delete) |
| PATCH | `/funcionarios/{id}/status` | Alterar status do funcionário |

### Apontamentos

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/apontamentos` | Criar novo apontamento |
| GET | `/apontamentos` | Listar apontamentos com filtros |
| GET | `/apontamentos/{id}` | Buscar apontamento por ID |
| PUT | `/apontamentos/{id}` | Atualizar apontamento |
| PATCH | `/apontamentos/{id}/aprovar` | Aprovar apontamento |
| PATCH | `/apontamentos/{id}/rejeitar` | Rejeitar apontamento |
| POST | `/apontamentos/{id}/pagamento` | Registrar pagamento |
| GET | `/funcionarios/{id}/apontamentos` | Apontamentos por funcionário |
| GET | `/obras/{id}/apontamentos` | Apontamentos por obra |

## Exemplos de Uso

### Criar Funcionário

```http
POST /funcionarios
Content-Type: application/json

{
  "nome": "João Carlos Silva",
  "email": "joao.silva@construtora.com",
  "telefone": "(11) 99999-9999",
  "cpf": "123.456.789-00",
  "cargo": "Pedreiro",
  "salario": 3500.00,
  "dataAdmissao": "2025-01-15T00:00:00Z",
  "endereco": "Rua das Palmeiras, 456 - Vila Nova"
}
```

### Criar Apontamento Quinzenal

```http
POST /apontamentos
Content-Type: application/json

{
  "funcionarioId": "funcionario-uuid",
  "obraId": "obra-uuid",
  "dataInicio": "2025-08-01T00:00:00Z",
  "dataFim": "2025-08-15T00:00:00Z",
  "horasTrabalhadas": 120,
  "valorHora": 25.00,
  "observacoes": "Trabalho na fundação da casa"
}
```

### Aprovar Apontamento

```http
PATCH /apontamentos/{apontamento-id}/aprovar
Content-Type: application/json

{
  "observacoes": "Apontamento aprovado após verificação em obra"
}
```

### Registrar Pagamento

```http
POST /apontamentos/{apontamento-id}/pagamento
Content-Type: application/json

{
  "valor": 3000.00,
  "observacoes": "Pagamento via PIX - Quinzena 01-15/08/2025",
  "dataPagamento": "2025-08-20T00:00:00Z"
}
```

### Listar Apontamentos com Filtros

```http
GET /apontamentos?status=Aprovado&funcionarioId=func-uuid&dataInicio=2025-08-01&dataFim=2025-08-31
```

## Fluxo de Trabalho

### 1. Ciclo do Apontamento

```
1. Funcionário registra horas trabalhadas
2. Sistema calcula valor total automaticamente
3. Apontamento fica com status "Pendente"
4. Supervisor/RH aprova ou rejeita
5. Se aprovado, RH pode processar pagamento
6. Sistema registra pagamento e atualiza status
7. Evento é enviado para módulo Financeiro
```

### 2. Estados do Apontamento

```
Pendente → Aprovado → Pago
    ↓
Rejeitado (fim do fluxo)
```

## Integrações e Eventos

### Eventos Publicados

1. **PagamentoApontamentoRealizado**
   - Publicado quando: Pagamento de apontamento é processado
   - Payload: Detalhes do pagamento e funcionário
   - Consumido por: Módulo Financeiro (registra saída de caixa)

### Eventos Consumidos

Atualmente o módulo não consome eventos, mas pode ser expandido para:
- **FuncionarioAlocado**: Quando funcionário é alocado a uma obra
- **ObraConcluida**: Para finalizar apontamentos pendentes

## Regras de Negócio

### Funcionários
- CPF deve ser único no sistema
- Email deve ser único e válido
- Salário deve ser positivo
- Status válidos: "Ativo", "Inativo", "Férias", "Licença"
- Data de demissão deve ser posterior à admissão
- Funcionários inativos não podem receber novos apontamentos

### Apontamentos
- Horas trabalhadas devem ser positivas e não exceder 200h por quinzena
- Valor por hora deve ser positivo
- Data fim deve ser posterior à data início
- Período não pode exceder 15 dias
- Funcionário deve estar ativo no período do apontamento
- Não pode haver sobreposição de períodos para o mesmo funcionário
- Valor calculado é automático: horas × valor_hora
- Apontamento só pode ser pago se estiver aprovado

### Aprovações e Pagamentos
- Apenas usuários com permissão podem aprovar apontamentos
- Apenas usuários com permissão podem processar pagamentos
- Apontamentos rejeitados não podem ser pagos
- Valor do pagamento pode ser diferente do valor calculado (ajustes)

## Validações e Constraints

### Validações de Entrada
- **CPF**: Formato válido (XXX.XXX.XXX-XX) e único
- **Email**: Formato válido e único
- **Telefone**: Formato brasileiro válido
- **Salário**: Positivo, máximo 2 casas decimais
- **Horas trabalhadas**: Entre 1 e 200 horas

### Constraints de Banco
- Unique constraint em CPF
- Unique constraint em email
- Check constraint em salário (> 0)
- Check constraint em horas trabalhadas (> 0)
- Foreign keys para obras e funcionários
- Índices para performance:
  - `idx_funcionarios_cpf` - Busca por CPF
  - `idx_funcionarios_status` - Filtro por status
  - `idx_apontamentos_funcionario_periodo` - Busca por funcionário e período
  - `idx_apontamentos_status` - Filtro por status

## Relatórios e Dashboards

### Relatório de Folha de Pagamento

```http
GET /relatorios/folha-pagamento?dataInicio=2025-08-01&dataFim=2025-08-31
```

**Resposta:**
```json
{
  "periodo": {
    "dataInicio": "2025-08-01",
    "dataFim": "2025-08-31"
  },
  "resumo": {
    "totalFuncionarios": 25,
    "totalApontamentos": 50,
    "valorTotalAprovado": 150000.00,
    "valorTotalPago": 145000.00,
    "valorPendente": 5000.00
  },
  "funcionarios": [
    {
      "id": "func-uuid",
      "nome": "João Silva",
      "cargo": "Pedreiro",
      "totalHoras": 160,
      "valorCalculado": 4000.00,
      "valorPago": 4000.00,
      "status": "Pago"
    }
  ]
}
```

### Dashboard de RH

- Funcionários ativos vs inativos
- Distribuição por cargo
- Média salarial por cargo
- Taxa de aprovação de apontamentos
- Tempo médio de aprovação
- Valor total pago por mês

### Relatório Individual do Funcionário

```http
GET /funcionarios/{funcionario-id}/relatorio?ano=2025
```

- Histórico de apontamentos
- Total de horas trabalhadas
- Valor total recebido
- Obras em que trabalhou
- Performance mensal

## Métricas e Monitoramento

### Métricas Importantes
- Taxa de aprovação de apontamentos
- Tempo médio entre criação e aprovação
- Tempo médio entre aprovação e pagamento
- Valor médio por hora por cargo
- Funcionários mais produtivos (horas/mês)

### Logs Estruturados
```json
{
  "timestamp": "2025-08-08T14:30:00Z",
  "level": "INFO",
  "service": "PessoalService",
  "operation": "AprovarApontamento",
  "apontamento_id": "uuid",
  "funcionario_id": "uuid",
  "valor_aprovado": 3000.00,
  "aprovado_por": "supervisor-uuid"
}
```

### Alertas
- Apontamentos pendentes há mais de 3 dias
- Pagamentos aprovados há mais de 5 dias sem processamento
- Funcionários sem apontamentos há mais de 30 dias
- Valores por hora discrepantes da média do cargo

## Integrações Futuras

### Ponto Eletrônico
- Integração com sistemas de ponto
- Importação automática de horas trabalhadas
- Validação cruzada de dados

### Sistema de Pagamento
- Integração com bancos para pagamentos automáticos
- Geração de arquivos CNAB
- Conciliação bancária

### Aplicativo Mobile
- App para funcionários registrarem apontamentos
- Aprovação via mobile para supervisores
- Consulta de histórico de pagamentos

## Configurações e Parametrizações

### Regras de Negócio Configuráveis
```json
{
  "horasMaximasPorQuinzena": 200,
  "diasMaximosAprovacao": 3,
  "diasMaximosPagamento": 5,
  "valorMinimoHora": 15.00,
  "cargos": [
    "Pedreiro",
    "Servente",
    "Mestre de Obras",
    "Encarregado",
    "Supervisor"
  ]
}
```

### Permissões Específicas
- `pessoal:funcionario:criar`: Criar funcionários
- `pessoal:funcionario:editar`: Editar dados de funcionários
- `pessoal:apontamento:aprovar`: Aprovar apontamentos
- `pessoal:apontamento:pagar`: Processar pagamentos
- `pessoal:relatorio:visualizar`: Visualizar relatórios

## Troubleshooting

### Problemas Comuns

#### "Funcionário já possui apontamento no período"
- Verificar se não há sobreposição de datas
- Verificar se apontamento anterior foi finalizado

#### "Valor por hora muito baixo/alto"
- Verificar configurações de valor mínimo
- Comparar com média do cargo

#### "Apontamento não pode ser aprovado"
- Verificar se funcionário está ativo
- Verificar se período é válido
- Verificar permissões do usuário

### Comandos Úteis
```sql
-- Verificar apontamentos pendentes
SELECT * FROM apontamentos_quinzenais 
WHERE status = 'Pendente' 
AND created_at < NOW() - INTERVAL '3 days';

-- Calcular valor total por funcionário no mês
SELECT f.nome, SUM(aq.valor_calculado) as total
FROM funcionarios f
JOIN apontamentos_quinzenals aq ON f.id = aq.funcionario_id
WHERE aq.data_inicio >= '2025-08-01'
AND aq.data_fim <= '2025-08-31'
GROUP BY f.id, f.nome;
```