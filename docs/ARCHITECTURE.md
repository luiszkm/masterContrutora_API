# Arquitetura do Sistema Master Construtora

## Visão Geral

O Master Construtora é um sistema de gestão para empresas de construção civil, construído como um **Monólito Modular** seguindo princípios de **Clean Architecture**. O sistema é organizado em módulos independentes (Bounded Contexts) que representam diferentes domínios de negócio.

## Princípios Arquiteturais

### 1. Monólito Modular
- **Única aplicação deployável**: Um binário Go que contém todos os módulos
- **Módulos bem definidos**: Cada contexto de negócio é isolado
- **Comunicação controlada**: Módulos se comunicam via interfaces bem definidas
- **Facilita evolução**: Possibilidade futura de extração para microserviços

### 2. Clean Architecture
- **Dependências apontam para dentro**: Camadas externas dependem das internas
- **Domínio isolado**: Regras de negócio não dependem de frameworks
- **Testabilidade**: Cada camada pode ser testada independentemente
- **Flexibilidade**: Fácil substituição de componentes externos

### 3. CQRS (Command Query Responsibility Segregation)
- **Separação de responsabilidades**: Commands (escrita) e Queries (leitura) separados
- **Modelos otimizados**: DTOs específicos para diferentes necessidades
- **Performance**: Queries otimizadas para casos específicos (ex: Dashboard)

### 4. Event-Driven Architecture
- **Comunicação assíncrona**: Módulos se comunicam via eventos
- **Baixo acoplamento**: Módulos não precisam conhecer uns aos outros diretamente
- **Extensibilidade**: Novos handlers podem ser adicionados facilmente

## Estrutura de Módulos (Bounded Contexts)

### 1. Identidade (Identity)
**Responsabilidade**: Autenticação e autorização de usuários

- **Entidades**: Usuario
- **Casos de Uso**: Registro, Login, Gerenciamento de Permissões
- **Tecnologias**: JWT, bcrypt

### 2. Obras (Construction Projects)
**Responsabilidade**: Gestão de projetos de construção

- **Entidades**: Obra, Etapa, EtapaPadrao, Alocacao, CronogramaRecebimento
- **Casos de Uso**: 
  - CRUD de obras com controle financeiro
  - Gestão de etapas de construção
  - Alocação de funcionários
  - Cronogramas de recebimento por etapas
  - Controle de valores contratuais
- **Funcionalidades Especiais**: Dashboard com métricas calculadas e financeiras
- **Integração**: Publica eventos de cronogramas criados para o módulo Financeiro

### 3. Pessoal (Personnel)
**Responsabilidade**: Gestão de recursos humanos

- **Entidades**: Funcionario, ApontamentoQuinzenal
- **Casos de Uso**: Cadastro de funcionários, Apontamentos de horas, Aprovação de pagamentos
- **Integração**: Comunica com Financeiro via eventos

### 4. Suprimentos (Supplies)
**Responsabilidade**: Gestão de fornecedores e materiais

- **Entidades**: Fornecedor, Produto, Orcamento, Categoria
- **Casos de Uso**: Gestão de fornecedores, Controle de estoque, Orçamentos
- **Integração**: Envia eventos para atualização de custos

### 5. Financeiro (Financial)
**Responsabilidade**: Controle financeiro completo da construtora

- **Entidades**: ContaReceber, ContaPagar, ParcelaContaPagar, CronogramaRecebimento, RegistroPagamento
- **Casos de Uso**: 
  - Gestão de contas a receber (receitas de obras)
  - Gestão de contas a pagar (fornecedores e serviços)
  - Cronogramas de recebimento por etapas
  - Fluxo de caixa consolidado
  - Processamento de pagamentos
- **Integração**: 
  - Recebe eventos de orçamentos aprovados (Suprimentos) → cria contas a pagar
  - Recebe eventos de pagamentos (Pessoal) → registra saídas
  - Publica eventos de movimentações financeiras

## Estrutura de Diretórios

```
masterContrutora/
├── cmd/                              # Pontos de entrada da aplicação
│   ├── server/main.go               # Servidor principal
│   └── seeder/main.go               # Populador de dados
├── internal/                        # Código interno da aplicação
│   ├── domain/                      # Camada de Domínio (Clean Architecture)
│   │   ├── common/                  # Utilitários compartilhados
│   │   ├── identidade/              # Entidades do domínio Identidade
│   │   ├── obras/                   # Entidades do domínio Obras
│   │   ├── pessoal/                 # Entidades do domínio Pessoal
│   │   ├── suprimentos/             # Entidades do domínio Suprimentos
│   │   └── financeiro/              # Entidades do domínio Financeiro
│   ├── service/                     # Camada de Aplicação (Use Cases)
│   │   ├── identidade/              # Serviços de Identidade
│   │   ├── obras/                   # Serviços de Obras
│   │   ├── pessoal/                 # Serviços de Pessoal
│   │   ├── suprimentos/             # Serviços de Suprimentos
│   │   └── financeiro/              # Serviços Financeiros
│   ├── handler/                     # Camada de Interface (Controllers)
│   │   ├── http/                    # Handlers HTTP
│   │   └── web/                     # Utilitários web
│   ├── infrastructure/              # Camada de Infraestrutura
│   │   └── repository/postgres/     # Implementações PostgreSQL
│   ├── platform/                    # Componentes de plataforma
│   │   └── bus/                     # Event Bus interno
│   ├── events/                      # Definições de eventos
│   └── authz/                       # Sistema de autorização
├── pkg/                             # Pacotes reutilizáveis
│   ├── auth/                        # Utilitários de autenticação
│   ├── security/                    # Utilitários de segurança
│   └── storage/                     # Utilitários de armazenamento
├── db/                              # Scripts de banco de dados
│   └── init/                        # Scripts de inicialização
└── docs/                            # Documentação
```

## Fluxo de Dados

### 1. Requisição HTTP
```
HTTP Request → Router → Middleware (Auth) → Handler → Service → Repository → Database
```

### 2. Comunicação entre Módulos
```
Módulo A → Event Bus → Event Handler → Módulo B
```

### 3. Exemplo: Aprovação de Orçamento
```
1. Usuário aprova orçamento (Suprimentos)
2. Evento "OrcamentoStatusAtualizado" é publicado
3. Handler no módulo Obras atualiza métricas financeiras
4. Dashboard é automaticamente atualizado
```

## Padrões Implementados

### Repository Pattern
- **Interfaces no domínio**: Contratos definidos na camada de domínio
- **Implementações na infraestrutura**: PostgreSQL como implementação
- **Testabilidade**: Fácil mock para testes unitários

### Dependency Injection
- **Configuração central**: Todas as dependências configuradas no main.go
- **Interfaces**: Dependências via interfaces, não implementações concretas
- **Flexibilidade**: Fácil substituição de implementações

### DTO Pattern
- **Separação de responsabilidades**: DTOs específicos para cada caso de uso
- **Versionamento**: Facilita evolução da API sem quebrar contratos
- **Validação**: Cada DTO tem suas próprias regras de validação

## Tecnologias Core

### Backend
- **Go 1.23+**: Linguagem principal
- **go-chi/chi**: Router HTTP
- **pgx/v5**: Driver PostgreSQL
- **golang-jwt/jwt**: Autenticação JWT
- **slog**: Logging estruturado

### Banco de Dados
- **PostgreSQL 16**: Banco principal
- **UUIDs**: Chaves primárias
- **Soft Delete**: Exclusão lógica
- **Índices otimizados**: Performance em consultas

### Infraestrutura
- **Docker**: Containerização
- **Docker Compose**: Orquestração local

## Principais Benefícios da Arquitetura

### 1. Manutenibilidade
- Código organizado por domínio
- Responsabilidades bem definidas
- Baixo acoplamento entre módulos

### 2. Testabilidade
- Interfaces bem definidas
- Dependências injetáveis
- Cada camada testável independentemente

### 3. Escalabilidade
- Módulos podem ser otimizados independentemente
- Possibilidade futura de extração para microserviços
- Event-driven permite scaling horizontal

### 4. Segurança
- Autenticação centralizada
- Autorização granular
- Validação em múltiplas camadas

### 5. Performance
- Queries otimizadas via CQRS
- Índices apropriados no banco
- Caching em pontos estratégicos

## Considerações para Evolução

### Microserviços
A arquitetura atual facilita uma eventual migração para microserviços:
- Cada módulo já tem boundaries bem definidos
- Comunicação via eventos já está implementada
- Interfaces de domínio facilitam extração de serviços

### Observabilidade
Pontos para melhoria:
- Métricas de performance
- Distributed tracing
- Health checks mais robustos

### Caching
Oportunidades de melhoria:
- Cache de queries frequentes
- Cache de sessões de usuário
- Cache de dados de dashboard