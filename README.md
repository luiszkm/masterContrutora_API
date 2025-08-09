# API Master Construtora

![Go Version](https://img.shields.io/badge/go-1.21%2B-blue.svg) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![PostgreSQL](https://img.shields.io/badge/postgresql-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)

API de backend para o sistema de gestão de construções "Master Construtora", projetada para ser robusta, escalável e de fácil manutenção.

## Visão Geral da Arquitetura

Este projeto foi construído seguindo princípios modernos de arquitetura de software para garantir a qualidade e a longevidade do código.

* [cite_start]**Monólito Modular**: A aplicação é um único binário implementável, mas seu código-fonte é rigorosamente organizado em Módulos (Bounded Contexts), como `Obras`, `Pessoal`, `Suprimentos` e `Financeiro`. 
* **Arquitetura Limpa (Clean Architecture)**: As dependências sempre apontam para o centro. A lógica de negócio e as entidades de domínio (`internal/domain`) são o núcleo e não possuem dependências de camadas externas como banco de dados ou frameworks web.
* **CQRS (Command Query Responsibility Segregation)**: Separamos a responsabilidade de escrita (Commands) da de leitura (Queries). [cite_start]Isso nos permite ter operações de escrita consistentes e, ao mesmo tempo, criar modelos de leitura otimizados e complexos, como o `ObraDashboard`. 
* **Comunicação Orientada a Eventos**: Módulos se comunicam de forma assíncrona através de um Event Bus interno. [cite_start]Isso desacopla os contextos, permitindo que um evento em um módulo (ex: `Orçamento Atualizado`) dispare ações em outros módulos sem que haja um acoplamento direto. 

## Funcionalidades Implementadas

* [cite_start]**Autenticação & Autorização**: Sistema de segurança completo com login via JWT (em cookies `httpOnly`) e autorização granular baseada em papéis e permissões (RBAC). 
* [cite_start]**Contexto de Identidade**: Registro e login de usuários. 
* [cite_start]**Contexto de Pessoal**: Cadastro, listagem e exclusão lógica de funcionários. 
* **Contexto de Obras**:
    * CRUD completo (Create, Read, Update, Delete) para Obras e Etapas.
    * Alocação de funcionários a obras.
    * Dashboard de leitura (`GET /obras/{id}`) com dados calculados em tempo real (percentual de conclusão, funcionários alocados, balanço financeiro).
* **Contexto de Suprimentos**:
    * CRUD completo para Fornecedores e Materiais.
    * Criação e atualização de status de Orçamentos com múltiplos itens.
* **Contexto Financeiro**:
    * [cite_start]Registro de pagamentos a funcionários, associados a uma obra.
* **Dashboard Executivo**:
    * API completa de dashboard com dados agregados de todas as seções
    * Métricas em tempo real de obras, funcionários, fornecedores e finanças
    * Sistema de alertas para obras em atraso e pagamentos pendentes
    * Logging estruturado para auditoria e monitoramento de performance 

## Tecnologias Utilizadas

* **Linguagem**: Go (1.21+)
* [cite_start]**Banco de Dados**: PostgreSQL 
* **Contêineres**: Docker & Docker Compose
* **Roteamento HTTP**: `go-chi/chi/v5`
* **Acesso ao Banco de Dados**: `jackc/pgx/v5`
* [cite_start]**Logging**: `log/slog` (para logs estruturados em JSON) 
* **Autenticação**: `golang-jwt/jwt/v5`
* **Migrations de Banco**: `golang-migrate/migrate` (usado para criar os scripts `up` e `down`)

## Pré-requisitos

Antes de começar, garanta que você tenha os seguintes softwares instalados:
* Go (versão 1.21 ou superior)
* Docker e Docker Compose
* A CLI do `golang-migrate` (se precisar criar novas migrations):
    ```sh
    go install -tags 'postgres' [github.com/golang-migrate/migrate/v4/cmd/migrate@latest](https://github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
    ```

## Configuração do Ambiente

1.  **Clone o Repositório**
    ```sh
    git clone [URL_DO_SEU_REPOSITORIO]
    cd masterCostrutora
    ```

2.  **Crie o Arquivo de Ambiente**
    Copie o arquivo de exemplo `.env.example` para um novo arquivo chamado `.env`.
    ```sh
    cp .env.example .env
    ```
    Edite o arquivo `.env` e preencha as variáveis, especialmente a `JWT_SECRET_KEY` com um valor longo e aleatório.

    **`.env` (Exemplo):**
    ```env
    # URL de conexão para o banco de dados PostgreSQL que roda via Docker
    DATABASE_URL="postgres://user:password@localhost:5432/mastercostrutora_db?sslmode=disable"

    # Chave secreta para assinar os tokens JWT. Use um valor forte.
    JWT_SECRET_KEY="SUA_CHAVE_SECRETA_E_LONGA_AQUI"

    # Define o ambiente da aplicação (usado para habilitar/desabilitar features de debug)
    APP_ENV="development"
    ```

## Como Executar a Aplicação

Este projeto foi desenhado para ser executado com o banco de dados em um contêiner Docker e a aplicação Go rodando localmente na sua máquina para facilitar a depuração.

1.  **Inicie o Banco de Dados**
    Este comando irá iniciar o contêiner do PostgreSQL em background. Na primeira vez, ele executará os scripts em `/db/init` para criar todas as tabelas.
    ```sh
    docker-compose up -d
    ```
    Para parar o banco, use `docker-compose down`. Para resetar os dados, use `docker-compose down -v`.

2.  **Execute a API Go**
    Em outro terminal, na raiz do projeto, execute a aplicação. Ela lerá o arquivo `.env` e se conectará ao banco de dados no Docker.
    ```sh
    go run ./cmd/server/main.go
    ```
    Você deverá ver a mensagem de log `servidor escutando na porta :8080`.

## Documentação da API

A documentação prática e executável da API está no arquivo `requests.http`, que pode ser usado com a extensão **REST Client** do Visual Studio Code.

Ele contém exemplos para todos os endpoints, incluindo o fluxo completo de autenticação e captura de variáveis para testes encadeados.

### Documentação Específica

- **[Dashboard API](docs/DASHBOARD_API.md)**: Documentação completa da API de dashboard com exemplos de payloads e casos de uso
- **[Arquivo de Testes](restclient/dashboard.http)**: Testes HTTP executáveis para todos os endpoints do dashboard

## Registro de Decisões Arquiteturais (ADRs)

As decisões de arquitetura mais importantes tomadas durante o desenvolvimento deste projeto estão documentadas como **ADRs (Architectural Decision Records)** na pasta `/docs/adr`. [cite_start]Eles explicam o "porquê" por trás das nossas escolhas técnicas.