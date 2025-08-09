# Módulo Suprimentos

O módulo Suprimentos é responsável pela gestão completa de fornecedores, produtos, materiais e orçamentos da construtora, incluindo controle de estoque, processo de cotação e integração com o sistema financeiro para criação automática de contas a pagar.

## Visão Geral

### Funcionalidades Principais
- **Gestão de Fornecedores**: Cadastro e avaliação de fornecedores
- **Catálogo de Produtos**: Materiais e serviços organizados por categoria
- **Sistema de Orçamentos**: Cotações e aprovações de compras
- **Controle de Categorias**: Organização hierárquica de produtos
- **Integração Financeira**: Criação automática de contas a pagar

### Arquitetura

```
suprimentos/
├── domain/
│   ├── fornecedor.go              # Entidade Fornecedor
│   ├── produto.go                 # Produtos e materiais
│   ├── orcamento.go               # Orçamentos de compra
│   ├── orcamento_item.go          # Itens do orçamento
│   └── categoria.go               # Categorias de produtos
├── service/
│   ├── suprimentos_service.go     # Lógica de negócio principal
│   └── dto/                       # Data Transfer Objects
├── handler/http/
│   └── suprimentos_handler.go     # Endpoints HTTP
├── infrastructure/repository/postgres/
│   ├── fornecedor_repository.go
│   ├── produto_repository.go
│   ├── orcamento_repository.go
│   └── categoria_repository.go
└── events/
    └── handler.go                 # Manipulador de eventos
```

## Entidades Principais

### 1. Fornecedor

```go
type Fornecedor struct {
    ID          string
    Nome        string
    Email       string
    Telefone    string
    CNPJ        string
    Endereco    string
    Contato     string    // Pessoa de contato
    Avaliacao   float64   // Nota de 0 a 5
    Status      string    // "Ativo", "Inativo", "Bloqueado"
    Observacoes *string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Métodos de Negócio:**
- `EstaAtivo() bool`: Verifica se fornecedor está ativo
- `PodeFornecerPara(obraID string) bool`: Verifica se pode fornecer para obra
- `AtualizarAvaliacao(nova float64)`: Atualiza nota do fornecedor
- `ValidarCNPJ() bool`: Valida formato do CNPJ

### 2. Produto

```go
type Produto struct {
    ID          string
    Nome        string
    Descricao   string
    Unidade     string    // "m²", "kg", "unidade", "m³", etc.
    CategoriaID string
    PrecoMedio  float64   // Preço médio baseado em cotações
    Status      string    // "Ativo", "Inativo", "Descontinuado"
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Métodos de Negócio:**
- `EstaDisponivel() bool`: Verifica se produto está disponível
- `AtualizarPrecoMedio(preco float64)`: Atualiza preço baseado em cotações
- `ObterHistoricoPrecos() []PrecoHistorico`: Histórico de preços

### 3. Orcamento

```go
type Orcamento struct {
    ID              string
    ObraID          *string
    FornecedorID    string
    Status          string    // "Pendente", "Aprovado", "Rejeitado", "Cancelado"
    ValorTotal      float64
    DataCriacao     time.Time
    DataAprovacao   *time.Time
    DataValidade    time.Time
    Observacoes     *string
    UsuarioAprovador *string
    MotivoRejeicao  *string
    CreatedAt       time.Time
    UpdatedAt       time.Time
    
    // Relacionamentos
    Itens []OrcamentoItem `json:"itens,omitempty"`
}
```

**Métodos de Negócio:**
- `PodeSerAprovado() bool`: Verifica se orçamento pode ser aprovado
- `EstaVencido() bool`: Verifica se orçamento está vencido
- `CalcularValorTotal()`: Recalcula valor total baseado nos itens
- `Aprovar(usuarioID, observacoes)`: Aprova orçamento
- `Rejeitar(motivo)`: Rejeita orçamento

### 4. OrcamentoItem

```go
type OrcamentoItem struct {
    ID          string
    OrcamentoID string
    ProdutoID   string
    Quantidade  float64
    PrecoUnitario float64
    PrecoTotal  float64
    Observacoes *string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Métodos de Negócio:**
- `CalcularPrecoTotal()`: Calcula preço total (quantidade × preço unitário)
- `ValidarQuantidade() bool`: Valida se quantidade é positiva

### 5. Categoria

```go
type Categoria struct {
    ID        string
    Nome      string
    Descricao string
    Cor       string    // Código hex para UI (#FF5733)
    Status    string    // "Ativa", "Inativa"
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## APIs Disponíveis

### Fornecedores

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/fornecedores` | Criar novo fornecedor |
| GET | `/fornecedores` | Listar fornecedores com filtros |
| GET | `/fornecedores/{id}` | Buscar fornecedor por ID |
| PUT | `/fornecedores/{id}` | Atualizar fornecedor |
| DELETE | `/fornecedores/{id}` | Inativar fornecedor |
| PATCH | `/fornecedores/{id}/avaliacao` | Atualizar avaliação |
| GET | `/fornecedores/{id}/orcamentos` | Orçamentos do fornecedor |

### Produtos

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/produtos` | Criar novo produto |
| GET | `/produtos` | Listar produtos com filtros |
| GET | `/produtos/{id}` | Buscar produto por ID |
| PUT | `/produtos/{id}` | Atualizar produto |
| DELETE | `/produtos/{id}` | Inativar produto |
| GET | `/produtos/{id}/historico-precos` | Histórico de preços |
| GET | `/categorias/{id}/produtos` | Produtos por categoria |

### Orçamentos

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/orcamentos` | Criar novo orçamento |
| GET | `/orcamentos` | Listar orçamentos com filtros |
| GET | `/orcamentos/{id}` | Buscar orçamento por ID |
| PUT | `/orcamentos/{id}` | Atualizar orçamento |
| PATCH | `/orcamentos/{id}/aprovar` | Aprovar orçamento |
| PATCH | `/orcamentos/{id}/rejeitar` | Rejeitar orçamento |
| POST | `/orcamentos/{id}/itens` | Adicionar item ao orçamento |
| PUT | `/orcamentos/{id}/itens/{item-id}` | Atualizar item |
| DELETE | `/orcamentos/{id}/itens/{item-id}` | Remover item |

### Categorias

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/categorias` | Criar nova categoria |
| GET | `/categorias` | Listar categorias ativas |
| GET | `/categorias/{id}` | Buscar categoria por ID |
| PUT | `/categorias/{id}` | Atualizar categoria |
| DELETE | `/categorias/{id}` | Inativar categoria |

## Exemplos de Uso

### Criar Fornecedor

```http
POST /fornecedores
Content-Type: application/json

{
  "nome": "Materiais de Construção ABC Ltda",
  "email": "vendas@materialsabc.com.br",
  "telefone": "(11) 3333-4444",
  "cnpj": "12.345.678/0001-99",
  "endereco": "Rua dos Fornecedores, 789 - Distrito Industrial",
  "contato": "Carlos Vendas",
  "observacoes": "Fornecedor especializado em cimento e agregados"
}
```

### Criar Produto

```http
POST /produtos
Content-Type: application/json

{
  "nome": "Cimento Portland CP II-E 50kg",
  "descricao": "Cimento Portland composto com escória para uso geral",
  "unidade": "saco",
  "categoriaId": "categoria-cimento-uuid",
  "precoMedio": 28.50
}
```

### Criar Orçamento com Itens

```http
POST /orcamentos
Content-Type: application/json

{
  "obraId": "obra-uuid",
  "fornecedorId": "fornecedor-uuid",
  "dataValidade": "2025-09-15T00:00:00Z",
  "observacoes": "Orçamento para materiais da fundação",
  "itens": [
    {
      "produtoId": "produto-cimento-uuid",
      "quantidade": 50,
      "precoUnitario": 28.50,
      "observacoes": "Entrega em até 5 dias úteis"
    },
    {
      "produtoId": "produto-areia-uuid", 
      "quantidade": 15,
      "precoUnitario": 45.00,
      "observacoes": "Areia média lavada"
    }
  ]
}
```

### Aprovar Orçamento

```http
PATCH /orcamentos/{orcamento-id}/aprovar
Content-Type: application/json

{
  "observacoes": "Orçamento aprovado - melhores preços do mercado"
}
```

### Listar Orçamentos com Filtros

```http
GET /orcamentos?status=Pendente&fornecedorId=forn-uuid&dataInicio=2025-08-01&dataFim=2025-08-31
```

## Fluxo de Trabalho

### 1. Processo de Orçamentação

```
1. Criar orçamento com itens
2. Sistema calcula valor total automaticamente
3. Orçamento fica com status "Pendente"
4. Gerente/Comprador analisa e aprova ou rejeita
5. Se aprovado, sistema cria conta a pagar automaticamente
6. Evento é enviado para módulo Financeiro
7. Preços médios dos produtos são atualizados
```

### 2. Estados do Orçamento

```
Pendente → Aprovado → [Integração com Financeiro]
    ↓
Rejeitado (fim do fluxo)
    ↓
Cancelado (pode ser reaberto)
```

## Integrações e Eventos

### Eventos Publicados

1. **OrcamentoStatusAtualizado**
   - Publicado quando: Status do orçamento muda para "Aprovado"
   - Payload: Detalhes do orçamento, fornecedor e valor total
   - Consumido por: Módulo Financeiro (cria conta a pagar automaticamente)

### Eventos Consumidos

Atualmente o módulo não consome eventos, mas pode ser expandido para:
- **ObraIniciada**: Para solicitar orçamentos automáticos
- **EstoqueMinimo**: Para disparar processo de recompra

## Regras de Negócio

### Fornecedores
- CNPJ deve ser único no sistema
- Email deve ser válido e único
- Avaliação deve estar entre 0 e 5
- Status válidos: "Ativo", "Inativo", "Bloqueado"
- Fornecedores inativos não podem ter novos orçamentos

### Produtos
- Nome deve ser único por categoria
- Preço médio deve ser positivo
- Unidade deve ser válida (m², kg, unidade, etc.)
- Produtos inativos não podem ser cotados

### Orçamentos
- Data de validade deve ser futura
- Valor total é calculado automaticamente (soma dos itens)
- Orçamentos vencidos não podem ser aprovados
- Quantidade nos itens deve ser positiva
- Preço unitário deve ser positivo
- Apenas usuários autorizados podem aprovar

### Cotações e Preços
- Preço médio é atualizado quando orçamentos são aprovados
- Histórico de preços é mantido para análise
- Variações acima de 20% geram alertas

## Validações e Constraints

### Validações de Entrada
- **CNPJ**: Formato válido (XX.XXX.XXX/XXXX-XX) e único
- **Email**: Formato válido
- **Avaliação**: Entre 0.0 e 5.0
- **Quantidade**: Positiva
- **Preço**: Positivo, máximo 2 casas decimais

### Constraints de Banco
- Unique constraint em CNPJ de fornecedores
- Check constraint em avaliação (>= 0 AND <= 5)
- Check constraint em quantidades e preços (> 0)
- Foreign keys para produtos, fornecedores e categorias
- Índices para performance:
  - `idx_fornecedores_cnpj` - Busca por CNPJ
  - `idx_fornecedores_status` - Filtro por status
  - `idx_orcamentos_fornecedor_data` - Orçamentos por fornecedor e período
  - `idx_orcamentos_status` - Filtro por status

## Dashboard e Relatórios

### Dashboard de Suprimentos

```http
GET /dashboard/suprimentos
```

**Resposta:**
```json
{
  "resumoOrcamentos": {
    "totalPendentes": 15,
    "totalAprovadosMes": 45,
    "valorTotalAprovado": 250000.00,
    "tempoMedioAprovacao": 2.5
  },
  "fornecedores": {
    "totalAtivos": 120,
    "avaliacaoMedia": 4.2,
    "novosCadastrosMes": 5
  },
  "produtos": {
    "totalAtivos": 450,
    "categorias": 25,
    "variacao_precos": [
      {
        "produto": "Cimento CP II",
        "variacao": "+5.2%",
        "periodo": "30 dias"
      }
    ]
  }
}
```

### Relatório de Orçamentos por Período

```http
GET /relatorios/orcamentos?dataInicio=2025-08-01&dataFim=2025-08-31
```

- Orçamentos criados vs. aprovados
- Valor total por fornecedor
- Tempo médio de aprovação
- Produtos mais cotados
- Análise de preços por categoria

### Relatório de Performance de Fornecedores

```http
GET /relatorios/fornecedores-performance?ano=2025
```

- Ranking por avaliação
- Volume de negócios por fornecedor
- Tempo de resposta médio
- Taxa de aprovação de orçamentos
- Cumprimento de prazos

## Análises e Métricas

### Análise de Preços
- Variação de preços por produto
- Comparativo entre fornecedores
- Tendências de mercado
- Alertas de variações significativas

### Performance de Fornecedores
- Tempo médio de resposta
- Taxa de aprovação
- Qualidade dos produtos entregues
- Cumprimento de prazos

### Análise de Compras
- Produtos mais comprados
- Categorias com maior gasto
- Sazonalidade das compras
- Oportunidades de economia

## Integrações Futuras

### E-procurement
- Integração com plataformas de compras
- Catálogos eletrônicos de fornecedores
- Leilões reversos automatizados

### Sistema de Qualidade
- Avaliação de produtos recebidos
- Controle de não conformidades
- Certificações de fornecedores

### Análise Preditiva
- Previsão de demanda de materiais
- Otimização de estoques
- Análise de risco de fornecedores

## Configurações

### Parâmetros do Sistema
```json
{
  "diasValidadeOrcamento": 30,
  "percentualAlertaVariacaoPreco": 20,
  "avaliacaoMinimaFornecedor": 2.0,
  "limiteMinimoCompra": 1000.00,
  "unidadesPadroes": [
    "unidade", "kg", "m²", "m³", "m", "litro", "saco", "caixa"
  ]
}
```

### Permissões Específicas
- `suprimentos:fornecedor:criar`: Cadastrar fornecedores
- `suprimentos:orcamento:aprovar`: Aprovar orçamentos
- `suprimentos:relatorio:visualizar`: Visualizar relatórios
- `suprimentos:preco:alterar`: Alterar preços de produtos

## Monitoramento e Logs

### Métricas Importantes
- Tempo médio de aprovação de orçamentos
- Taxa de aprovação por fornecedor
- Variação de preços por categoria
- Volume de compras por mês
- Fornecedores mais utilizados

### Logs Estruturados
```json
{
  "timestamp": "2025-08-08T15:00:00Z",
  "level": "INFO",
  "service": "SuprimentosService",
  "operation": "AprovarOrcamento",
  "orcamento_id": "uuid",
  "fornecedor_id": "uuid",
  "valor_total": 50000.00,
  "aprovado_por": "gerente-uuid"
}
```

### Alertas
- Orçamentos pendentes há mais de 5 dias
- Variações de preço acima de 20%
- Fornecedores com avaliação baixa
- Orçamentos próximos do vencimento

## Troubleshooting

### Problemas Comuns

#### "CNPJ já cadastrado"
- Verificar se fornecedor não existe
- Verificar formato do CNPJ

#### "Orçamento não pode ser aprovado"
- Verificar se está vencido
- Verificar permissões do usuário
- Verificar se fornecedor está ativo

#### "Preço muito acima da média"
- Verificar histórico de preços do produto
- Analisar variações de mercado
- Confirmar se preço está correto

### Comandos Úteis
```sql
-- Orçamentos pendentes há mais de 5 dias
SELECT * FROM orcamentos 
WHERE status = 'Pendente' 
AND created_at < NOW() - INTERVAL '5 days';

-- Fornecedores com melhor avaliação
SELECT nome, avaliacao, COUNT(*) as total_orcamentos
FROM fornecedores f
JOIN orcamentos o ON f.id = o.fornecedor_id
WHERE f.status = 'Ativo'
GROUP BY f.id, nome, avaliacao
ORDER BY avaliacao DESC;
```