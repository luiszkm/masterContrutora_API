# Documentação Master Construtora

Esta pasta contém toda a documentação técnica do projeto Master Construtora, um sistema de gestão para empresas de construção civil.

## 📚 Índice da Documentação

### 🏗️ [ARCHITECTURE.md](./ARCHITECTURE.md)
**Arquitetura do Sistema**
- Visão geral da arquitetura modular monolítica
- Padrões implementados (Clean Architecture, CQRS, Event-Driven)
- Estrutura de módulos e bounded contexts
- Tecnologias utilizadas e decisões arquiteturais

### 🔌 [API.md](./API.md)
**Documentação da API REST**
- Todos os endpoints organizados por módulo
- Exemplos de request/response para cada operação
- Códigos de erro e tratamento
- Guia completo para integração com frontend

### 🗄️ [DATABASE.md](./DATABASE.md)
**Schema e Banco de Dados**
- Schema completo do PostgreSQL
- Relacionamentos entre tabelas
- Índices e otimizações de performance
- Estratégias de backup e migração

### 🔐 [AUTH.md](./AUTH.md)
**Autenticação e Autorização**
- Sistema JWT com cookies httpOnly
- Permissões granulares e papéis (RBAC)
- Fluxos de autenticação e middleware
- Boas práticas de segurança

### 📡 [EVENTS.md](./EVENTS.md)
**Sistema de Eventos**
- Event Bus interno e comunicação entre módulos
- Eventos disponíveis e seus payloads
- Implementação de handlers
- Padrões event-driven

### 💻 [FRONTEND.md](./FRONTEND.md)
**Guia de Integração Frontend**
- Configuração e consumo da API
- Modelos de dados TypeScript
- Componentes React de exemplo
- Tratamento de autenticação e erros

### 🛠️ [DEVELOPMENT.md](./DEVELOPMENT.md)
**Guia de Desenvolvimento**
- Setup do ambiente local
- Estrutura do projeto e convenções
- Fluxo de desenvolvimento de features
- Testes e debugging

### 🚀 [DEPLOYMENT.md](./DEPLOYMENT.md)
**Guia de Deploy**
- Deploy local, staging e produção
- Configurações Docker e Kubernetes
- AWS, GCP e outras clouds
- CI/CD, monitoramento e backup

## 🚀 Quick Start

Para começar rapidamente:

1. **Setup Local**: Siga [DEVELOPMENT.md](./DEVELOPMENT.md#configuração-do-ambiente)
2. **Entender a API**: Consulte [API.md](./API.md#autenticação) 
3. **Integrar Frontend**: Use [FRONTEND.md](./FRONTEND.md#configuração-inicial)

## 📋 Funcionalidades Principais

### 🏗️ **Obras**
- CRUD completo de obras e etapas
- Dashboard com métricas em tempo real
- Alocação de funcionários
- Acompanhamento de progresso

### 👥 **Pessoal**
- Gestão de funcionários
- Apontamentos quinzenais
- Aprovação e pagamento
- Controle de horas trabalhadas

### 📦 **Suprimentos**
- Cadastro de fornecedores e produtos
- Criação e gestão de orçamentos
- Controle de categorias
- Aprovação de compras

### 💰 **Financeiro**
- Registro de pagamentos
- Controle de custos por obra
- Integração com outros módulos

### 🔐 **Identidade**
- Autenticação JWT
- Controle de permissões granular
- Gestão de usuários

## 🏗️ Arquitetura Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Frontend     │    │    API REST     │    │   PostgreSQL    │
│   (React/Vue)   │◄──►│   (Go/Chi)      │◄──►│   (Database)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                       ┌──────┴──────┐
                       │  Event Bus  │
                       │ (Internal)  │
                       └─────────────┘
```

**Principais características:**
- **Modular Monolith** com bounded contexts
- **Clean Architecture** com dependências apontando para o core
- **CQRS** para separação de leitura e escrita
- **Event-Driven** para comunicação entre módulos

## 🛠️ Stack Tecnológica

### Backend
- **Go 1.23+** - Linguagem principal
- **Chi Router** - Roteamento HTTP
- **PostgreSQL 16** - Banco de dados
- **pgx/v5** - Driver PostgreSQL
- **JWT** - Autenticação
- **Docker** - Containerização

### Arquitetura
- **Clean Architecture** - Organização do código
- **Repository Pattern** - Acesso a dados
- **Event Bus** - Comunicação interna
- **RBAC** - Controle de acesso

## 📖 Convenções de Documentação

### Estrutura dos Documentos
- **Visão Geral**: Introdução e contexto
- **Implementação**: Detalhes técnicos
- **Exemplos**: Código e casos de uso
- **Troubleshooting**: Problemas comuns

### Formato de Código
- Exemplos em Go, JavaScript/TypeScript, SQL
- Comandos de terminal com descrição
- Configurações em YAML/JSON

### Referencias Cruzadas
- Links entre documentos relacionados
- Referências a linhas de código específicas
- Menções a endpoints e entidades

## 🤝 Como Contribuir

### Para a Documentação
1. Identificar gaps ou informações desatualizadas
2. Criar/atualizar documentos seguindo o padrão
3. Testar exemplos de código
4. Solicitar review da equipe

### Para o Código
1. Seguir [DEVELOPMENT.md](./DEVELOPMENT.md#fluxo-de-desenvolvimento)
2. Atualizar documentação relevante
3. Incluir testes para novas features
4. Documentar mudanças na API

## 📞 Suporte

### Documentação
- 📁 **Issues**: Reportar problemas na documentação
- 💬 **Discussões**: Sugerir melhorias
- 📝 **Pull Requests**: Contribuir com atualizações

### Desenvolvimento
- 🐛 **Bugs**: Reportar problemas técnicos
- ✨ **Features**: Sugerir novas funcionalidades
- ❓ **Dúvidas**: Fazer perguntas técnicas

## 📈 Roadmap da Documentação

### ✅ Concluído
- Documentação completa da API
- Guias de setup e desenvolvimento
- Arquitetura e padrões
- Integração frontend

### 🔄 Em Andamento
- Exemplos práticos de uso
- Vídeos tutoriais
- Troubleshooting avançado

### 📋 Planejado
- Guia de performance
- Documentação de testes
- Guia de contribuição
- Glossário técnico

---

**Última atualização**: Agosto 2025  
**Versão da documentação**: 1.0  
**Versão da aplicação**: 1.0