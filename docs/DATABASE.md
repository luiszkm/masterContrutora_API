# Documentação do Banco de Dados - Master Construtora

## Visão Geral

O sistema Master Construtora utiliza **PostgreSQL 16** como banco de dados principal. O schema foi projetado para suportar a arquitetura modular da aplicação, com tabelas organizadas por contextos de negócio e relacionamentos bem definidos.

## Características Gerais

### Convenções
- **Chaves Primárias**: Todas as tabelas usam UUIDs como chave primária
- **Timestamps**: Campos `created_at` e `updated_at` com timezone (TIMESTAMPTZ)
- **Soft Delete**: Implementado via campo `deleted_at` onde aplicável
- **Naming**: Nomes em português para melhor compreensão do domínio
- **Índices**: Otimizados para consultas frequentes e JOINs

### Tipos de Dados Padrão
- **IDs**: `UUID`
- **Dinheiro**: `NUMERIC(10,2)` ou `NUMERIC(15,2)` para valores grandes
- **Texto Longo**: `TEXT`
- **Texto Curto**: `VARCHAR(n)` com limite apropriado
- **Status**: `VARCHAR(50)` para enums de status
- **Datas**: `DATE` para datas simples, `TIMESTAMPTZ` para timestamps completos

## Schema de Tabelas

### 1. Contexto Identidade

#### usuarios
Gerenciamento de usuários e autenticação.

```sql
CREATE TABLE usuarios (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha_hash TEXT NOT NULL,
    permissoes TEXT[] NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE
);
```

**Campos:**
- `id`: Identificador único do usuário
- `nome`: Nome completo do usuário
- `email`: Email único para login
- `senha_hash`: Hash da senha (bcrypt)
- `permissoes`: Array de permissões do usuário
- `ativo`: Status ativo/inativo do usuário

**Índices:**
- Primary Key em `id`
- Unique em `email`

### 2. Contexto Obras

#### obras
Entidade principal para projetos de construção.

```sql
CREATE TABLE obras (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL UNIQUE,
    cliente VARCHAR(255) NOT NULL,
    endereco TEXT NOT NULL,
    descricao TEXT,
    data_inicio DATE NOT NULL,
    data_fim DATE,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
```

**Campos:**
- `id`: Identificador único da obra
- `nome`: Nome único da obra
- `cliente`: Nome do cliente/contratante
- `endereco`: Endereço completo da obra
- `descricao`: Descrição detalhada (opcional)
- `data_inicio`: Data de início planejada/real
- `data_fim`: Data de conclusão (opcional)
- `status`: Status atual (`Em Planejamento`, `Em Andamento`, `Concluída`, `Cancelada`)
- `deleted_at`: Timestamp de exclusão lógica

**Índices:**
- Primary Key em `id`
- Unique em `nome`

#### etapas
Fases/etapas dentro de uma obra.

```sql
CREATE TABLE etapas (
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    data_inicio_prevista DATE,
    data_fim_prevista DATE,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Campos:**
- `id`: Identificador único da etapa
- `obra_id`: Referência à obra (FK)
- `nome`: Nome da etapa
- `data_inicio_prevista`: Data prevista para início
- `data_fim_prevista`: Data prevista para conclusão
- `status`: Status da etapa (`Não Iniciada`, `Em Andamento`, `Concluída`)

**Índices:**
- Primary Key em `id`
- Foreign Key em `obra_id`
- Index em `obra_id` para otimizar JOINs

#### etapas_padrao
Templates de etapas para reutilização.

```sql
CREATE TABLE etapas_padrao (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL UNIQUE,
    descricao TEXT,
    ordem INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Campos:**
- `id`: Identificador único da etapa padrão
- `nome`: Nome único da etapa padrão
- `descricao`: Descrição da etapa
- `ordem`: Ordem sugerida na sequência de etapas

#### alocacoes
Alocação de funcionários em obras.

```sql
CREATE TABLE alocacoes (
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras(id) ON DELETE CASCADE,
    funcionario_id UUID NOT NULL REFERENCES funcionarios(id) ON DELETE RESTRICT,
    data_inicio_alocacao DATE NOT NULL,
    data_fim_alocacao DATE
);
```

**Campos:**
- `id`: Identificador único da alocação
- `obra_id`: Referência à obra (FK)
- `funcionario_id`: Referência ao funcionário (FK)
- `data_inicio_alocacao`: Data de início da alocação
- `data_fim_alocacao`: Data de fim da alocação (opcional)

**Índices:**
- Primary Key em `id`
- Index em `obra_id`
- Index em `funcionario_id`

### 3. Contexto Pessoal

#### funcionarios
Cadastro de funcionários.

```sql
CREATE TABLE funcionarios (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) NOT NULL UNIQUE,
    telefone VARCHAR(20) NOT NULL DEFAULT '',
    cargo VARCHAR(100) NOT NULL,
    departamento VARCHAR(100) NOT NULL DEFAULT '',
    email VARCHAR(100) NOT NULL DEFAULT '',
    data_contratacao DATE NOT NULL,
    valor_diaria NUMERIC(10, 2) NOT NULL,
    chave_pix VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL,
    desligamento_data TIMESTAMPTZ,
    motivo_desligamento TEXT NOT NULL DEFAULT '',
    observacoes TEXT NOT NULL DEFAULT '',
    avaliacao_desempenho TEXT NOT NULL DEFAULT '',
    avatar_url VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Campos:**
- `id`: Identificador único do funcionário
- `nome`: Nome completo
- `cpf`: CPF único do funcionário
- `telefone`: Telefone de contato
- `cargo`: Cargo/função do funcionário
- `departamento`: Departamento de trabalho
- `email`: Email do funcionário
- `data_contratacao`: Data de contratação
- `valor_diaria`: Valor da diária de trabalho
- `chave_pix`: Chave PIX para pagamentos
- `status`: Status atual (`Ativo`, `Inativo`, `Desligado`)
- `desligamento_data`: Data de desligamento (se aplicável)
- `motivo_desligamento`: Motivo do desligamento
- `observacoes`: Observações gerais
- `avaliacao_desempenho`: Avaliação de desempenho
- `avatar_url`: URL da foto do funcionário

**Índices:**
- Primary Key em `id`
- Unique em `cpf`

#### apontamentos_quinzenais
Registro de horas trabalhadas por quinzena.

```sql
CREATE TABLE apontamentos_quinzenais (
    id UUID PRIMARY KEY,
    funcionario_id UUID NOT NULL REFERENCES funcionarios(id),
    obra_id UUID NOT NULL REFERENCES obras(id),
    periodo_inicio DATE NOT NULL,
    periodo_fim DATE NOT NULL,
    dias_trabalhados INT NOT NULL DEFAULT 0,
    adicionais NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    diaria NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    descontos NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    adiantamentos NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    valor_total_calculado NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(funcionario_id, periodo_inicio, periodo_fim)
);
```

**Campos:**
- `id`: Identificador único do apontamento
- `funcionario_id`: Referência ao funcionário (FK)
- `obra_id`: Referência à obra (FK)
- `periodo_inicio`: Data de início da quinzena
- `periodo_fim`: Data de fim da quinzena
- `dias_trabalhados`: Número de dias trabalhados
- `adicionais`: Valores adicionais (horas extras, etc.)
- `diaria`: Valor da diária utilizada
- `descontos`: Descontos aplicados
- `adiantamentos`: Adiantamentos dados
- `valor_total_calculado`: Valor total calculado para pagamento
- `status`: Status (`Em Aberto`, `Aprovado para Pagamento`, `Pago`)

**Índices:**
- Primary Key em `id`
- Unique em `(funcionario_id, periodo_inicio, periodo_fim)`
- Index em `funcionario_id`

### 4. Contexto Suprimentos

#### fornecedores
Cadastro de fornecedores.

```sql
CREATE TABLE fornecedores (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cnpj VARCHAR(18) NOT NULL UNIQUE,
    contato VARCHAR(255),
    email VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    endereco TEXT,
    avaliacao NUMERIC(3, 1),
    observacoes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
```

**Campos:**
- `id`: Identificador único do fornecedor
- `nome`: Nome/razão social
- `cnpj`: CNPJ único do fornecedor
- `contato`: Nome da pessoa de contato
- `email`: Email de contato
- `status`: Status (`Ativo`, `Inativo`)
- `endereco`: Endereço completo
- `avaliacao`: Nota de avaliação (0.0 a 5.0)
- `observacoes`: Observações sobre o fornecedor
- `deleted_at`: Timestamp de exclusão lógica

**Índices:**
- Primary Key em `id`
- Unique em `cnpj`

#### produtos
Catálogo de produtos/materiais.

```sql
CREATE TABLE produtos (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL UNIQUE,
    descricao TEXT,
    unidade_de_medida VARCHAR(50) NOT NULL,
    categoria VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Campos:**
- `id`: Identificador único do produto
- `nome`: Nome único do produto
- `descricao`: Descrição detalhada
- `unidade_de_medida`: Unidade de medida (kg, m³, unidade, etc.)
- `categoria`: Categoria do produto

**Índices:**
- Primary Key em `id`
- Unique em `nome`
- Index em `categoria`

#### categorias
Categorias de produtos.

```sql
CREATE TABLE categorias (
    id UUID PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Campos:**
- `id`: Identificador único da categoria
- `nome`: Nome único da categoria

#### orcamentos
Orçamentos de fornecedores para etapas.

```sql
CREATE TABLE orcamentos (
    id UUID PRIMARY KEY,
    numero VARCHAR(50) NOT NULL UNIQUE,
    etapa_id UUID NOT NULL REFERENCES etapas(id),
    fornecedor_id UUID NOT NULL REFERENCES fornecedores(id),
    valor_total NUMERIC(15, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    data_emissao TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    data_aprovacao TIMESTAMPTZ,
    observacoes TEXT,
    condicoes_pagamento VARCHAR(255),
    criado_por_usuario_id UUID REFERENCES usuarios(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Campos:**
- `id`: Identificador único do orçamento
- `numero`: Número único do orçamento (auto-gerado)
- `etapa_id`: Referência à etapa (FK)
- `fornecedor_id`: Referência ao fornecedor (FK)
- `valor_total`: Valor total do orçamento
- `status`: Status (`Em Aberto`, `Aprovado`, `Rejeitado`, `Pago`)
- `data_emissao`: Data de emissão do orçamento
- `data_aprovacao`: Data de aprovação (se aplicável)
- `observacoes`: Observações sobre o orçamento
- `condicoes_pagamento`: Condições de pagamento acordadas
- `criado_por_usuario_id`: Usuário que criou o orçamento

**Índices:**
- Primary Key em `id`
- Unique em `numero`
- Index em `etapa_id`
- Index em `fornecedor_id`

#### orcamento_itens
Itens individuais de cada orçamento.

```sql
CREATE TABLE orcamento_itens (
    id UUID PRIMARY KEY,
    orcamento_id UUID NOT NULL REFERENCES orcamentos(id) ON DELETE CASCADE,
    produto_id UUID NOT NULL REFERENCES produtos(id),
    quantidade NUMERIC(10, 2) NOT NULL,
    valor_unitario NUMERIC(10, 2) NOT NULL
);
```

**Campos:**
- `id`: Identificador único do item
- `orcamento_id`: Referência ao orçamento (FK)
- `produto_id`: Referência ao produto (FK)
- `quantidade`: Quantidade solicitada
- `valor_unitario`: Valor unitário do produto

#### fornecedor_categorias
Relacionamento many-to-many entre fornecedores e categorias.

```sql
CREATE TABLE fornecedor_categorias (
    fornecedor_id UUID NOT NULL REFERENCES fornecedores(id) ON DELETE CASCADE,
    categoria_id UUID NOT NULL REFERENCES categorias(id) ON DELETE CASCADE,
    PRIMARY KEY (fornecedor_id, categoria_id)
);
```

### 5. Contexto Financeiro

#### registros_pagamento
Registro de pagamentos realizados.

```sql
CREATE TABLE registros_pagamento (
    id UUID PRIMARY KEY,
    funcionario_id UUID NOT NULL REFERENCES funcionarios(id),
    obra_id UUID NOT NULL REFERENCES obras(id),
    periodo_referencia VARCHAR(100) NOT NULL,
    valor_calculado NUMERIC(10, 2) NOT NULL,
    data_de_efetivacao TIMESTAMPTZ NOT NULL,
    conta_bancaria_id UUID NOT NULL
);
```

**Campos:**
- `id`: Identificador único do pagamento
- `funcionario_id`: Referência ao funcionário (FK)
- `obra_id`: Referência à obra (FK)
- `periodo_referencia`: Período de referência do pagamento
- `valor_calculado`: Valor calculado para pagamento
- `data_de_efetivacao`: Data em que o pagamento foi efetivado
- `conta_bancaria_id`: Identificador da conta bancária

**Índices:**
- Primary Key em `id`
- Index em `funcionario_id`

## Relacionamentos

### Diagrama de Relacionamentos Principais

```
usuarios
├── orcamentos (criado_por_usuario_id)

obras
├── etapas (obra_id)
├── alocacoes (obra_id)
├── apontamentos_quinzenais (obra_id)
└── registros_pagamento (obra_id)

funcionarios
├── alocacoes (funcionario_id)
├── apontamentos_quinzenais (funcionario_id)
└── registros_pagamento (funcionario_id)

fornecedores
├── orcamentos (fornecedor_id)
└── fornecedor_categorias (fornecedor_id)

etapas
└── orcamentos (etapa_id)

produtos
└── orcamento_itens (produto_id)

orcamentos
└── orcamento_itens (orcamento_id)

categorias
└── fornecedor_categorias (categoria_id)
```

### Tipos de Relacionamento

#### One-to-Many (1:N)
- `obras` → `etapas`
- `obras` → `alocacoes`
- `funcionarios` → `alocacoes`
- `funcionarios` → `apontamentos_quinzenais`
- `etapas` → `orcamentos`
- `fornecedores` → `orcamentos`
- `orcamentos` → `orcamento_itens`

#### Many-to-Many (N:N)
- `fornecedores` ↔ `categorias` (via `fornecedor_categorias`)

#### Referências Opcionais
- `usuarios` → `orcamentos` (criado_por_usuario_id - pode ser NULL)

## Estratégias de Performance

### Índices Implementados

```sql
-- Índices para otimizar JOINs frequentes
CREATE INDEX idx_etapas_obra_id ON etapas(obra_id);
CREATE INDEX idx_orcamentos_etapa_id ON orcamentos(etapa_id);
CREATE INDEX idx_orcamentos_fornecedor_id ON orcamentos(fornecedor_id);
CREATE INDEX idx_alocacoes_obra_id ON alocacoes(obra_id);
CREATE INDEX idx_alocacoes_funcionario_id ON alocacoes(funcionario_id);
CREATE INDEX idx_produtos_categoria ON produtos(categoria);
CREATE INDEX idx_apontamentos_funcionario_id ON apontamentos_quinzenais(funcionario_id);
CREATE INDEX idx_pagamentos_funcionario_id ON registros_pagamento(funcionario_id);
CREATE INDEX idx_fornecedor_categorias_categoria_id ON fornecedor_categorias(categoria_id);
```

### Consultas Otimizadas

#### Dashboard de Obra
Query complexa que calcula métricas em tempo real:
- Percentual de conclusão baseado em etapas
- Custos realizados vs orçados
- Funcionários alocados ativos

#### Listagem de Funcionários com Apontamentos
Query que traz dados do último apontamento de cada funcionário para otimizar a interface.

## Política de Backup

### Estratégia Recomendada
1. **Backup Completo**: Diário, fora do horário comercial
2. **Backup Incremental**: A cada 4 horas durante horário comercial
3. **Retenção**: 30 dias para backups diários, 7 dias para incrementais
4. **Teste de Restore**: Mensal em ambiente de teste

### Scripts de Backup

```bash
# Backup completo
pg_dump -h localhost -U user -d mastercostrutora_db -f backup_$(date +%Y%m%d).sql

# Backup com compressão
pg_dump -h localhost -U user -d mastercostrutora_db | gzip > backup_$(date +%Y%m%d).sql.gz
```

## Migrações

### Versionamento do Schema
- Scripts de migração organizados em ordem cronológica
- Cada migration inclui script UP e DOWN
- Uso de ferramentas como `golang-migrate`

### Estrutura de Migrations
```
migrations/
├── 001_initial_schema.up.sql
├── 001_initial_schema.down.sql
├── 002_add_etapas_padrao.up.sql
├── 002_add_etapas_padrao.down.sql
```

## Monitoramento

### Métricas Importantes
- **Conexões ativas**: Monitorar pool de conexões
- **Slow queries**: Queries > 1 segundo
- **Deadlocks**: Bloqueios de transação
- **Tamanho do banco**: Crescimento das tabelas
- **Índices não utilizados**: Otimização de performance

### Queries de Monitoramento

```sql
-- Conexões ativas
SELECT count(*) FROM pg_stat_activity WHERE state = 'active';

-- Queries lentas
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;

-- Tamanho das tabelas
SELECT schemaname, tablename, 
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public' 
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## Considerações de Segurança

### Controle de Acesso
- **Usuário da aplicação**: Apenas permissões necessárias
- **Usuário administrativo**: Para manutenção e backup
- **Conexões criptografadas**: SSL/TLS obrigatório em produção

### Auditoria
- **Logs de conexão**: Registrar tentativas de conexão
- **Logs de modificação**: Triggers para tabelas críticas
- **Retenção de logs**: Mínimo 90 dias

### Backup de Segurança
- **Criptografia**: Backups sempre criptografados
- **Armazenamento externo**: Cópias em local seguro
- **Teste de integridade**: Verificação periódica dos backups