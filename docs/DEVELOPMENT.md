# Guia de Desenvolvimento - Master Construtora

## Pré-requisitos

### Software Necessário

- **Go**: Versão 1.23.0 ou superior
- **Docker**: Para execução do banco de dados
- **Docker Compose**: Para orquestração dos containers
- **Git**: Para controle de versão
- **Make**: Para execução de comandos (opcional)

### Ferramentas Recomendadas

- **VS Code** ou **GoLand**: IDEs com ótimo suporte para Go
- **TablePlus** ou **pgAdmin**: Clientes PostgreSQL para visualização do banco
- **Postman** ou **REST Client** (VS Code): Para testes da API
- **golang-migrate**: Para criação de migrations (opcional)

### Extensões VS Code Recomendadas

```json
{
  "recommendations": [
    "golang.go",
    "humao.rest-client",
    "ms-vscode.vscode-json",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode"
  ]
}
```

## Configuração do Ambiente

### 1. Clone do Repositório

```bash
git clone <url-do-repositorio>
cd masterContrutora
```

### 2. Configuração das Variáveis de Ambiente

Copie o arquivo de exemplo e configure as variáveis:

```bash
cp .env.example .env
```

Edite o arquivo `.env`:

```env
# URL de conexão para o banco de dados PostgreSQL que roda via Docker
DATABASE_URL="postgres://user:password@localhost:5432/mastercostrutora_db?sslmode=disable"

# Chave secreta para assinar os tokens JWT. Use um valor forte.
JWT_SECRET_KEY="SUA_CHAVE_SECRETA_E_LONGA_AQUI"

# Define o ambiente da aplicação (usado para habilitar/desabilitar features de debug)
APP_ENV="development"
```

**⚠️ Importante**: Em produção, use uma chave JWT forte com pelo menos 256 bits (32 caracteres).

### 3. Instalação das Dependências

```bash
# Download das dependências
go mod download

# Verificar se tudo está correto
go mod verify
```

### 4. Inicialização do Banco de Dados

```bash
# Iniciar o container PostgreSQL
docker-compose up -d

# Verificar se o container está rodando
docker-compose ps

# Ver logs do banco (opcional)
docker-compose logs db
```

O banco será inicializado automaticamente com o schema definido em `db/init/01-init.sql`.

### 5. Execução da Aplicação

```bash
# Executar o servidor
go run ./cmd/server/main.go

# Ou usando o Makefile
make up    # Inicia o banco
# Em outro terminal:
go run ./cmd/server/main.go
```

A aplicação estará disponível em `http://localhost:8080`.

## Estrutura do Projeto

### Organização de Diretórios

```
masterContrutora/
├── cmd/                              # Pontos de entrada
│   ├── server/main.go               # Servidor principal
│   └── seeder/main.go               # Populador de dados
├── internal/                        # Código interno da aplicação
│   ├── domain/                      # Entidades de domínio
│   │   ├── common/                  # Utilitários compartilhados
│   │   ├── identidade/              # Entidades de identidade
│   │   ├── obras/                   # Entidades de obras
│   │   ├── pessoal/                 # Entidades de pessoal
│   │   ├── suprimentos/             # Entidades de suprimentos
│   │   └── financeiro/              # Entidades financeiras
│   ├── service/                     # Lógica de negócio
│   │   ├── identidade/              # Serviços de identidade
│   │   ├── obras/                   # Serviços de obras
│   │   ├── pessoal/                 # Serviços de pessoal
│   │   ├── suprimentos/             # Serviços de suprimentos
│   │   └── financeiro/              # Serviços financeiros
│   ├── handler/                     # Controladores HTTP
│   │   ├── http/                    # Handlers HTTP
│   │   └── web/                     # Utilitários web
│   ├── infrastructure/              # Infraestrutura
│   │   └── repository/postgres/     # Repositórios PostgreSQL
│   ├── platform/                    # Componentes de plataforma
│   │   └── bus/                     # Event Bus
│   ├── events/                      # Definições de eventos
│   └── authz/                       # Sistema de autorização
├── pkg/                             # Pacotes reutilizáveis
│   ├── auth/                        # Utilitários de autenticação
│   ├── security/                    # Utilitários de segurança
│   └── storage/                     # Utilitários de armazenamento
├── db/                              # Scripts de banco
│   └── init/                        # Scripts de inicialização
├── docs/                            # Documentação
├── restclient/                      # Arquivos de teste HTTP
└── docker-compose.yml              # Configuração Docker
```

### Convenções de Nomenclatura

#### Pacotes
- **Minúsculas**: `identidade`, `obras`, `pessoal`
- **Singular**: `service`, `handler`, `repository`

#### Arquivos
- **Snake_case** para arquivos: `usuario_repository.go`
- **CamelCase** para tipos: `UsuarioRepository`

#### Banco de Dados
- **Snake_case** para tabelas e colunas: `usuarios`, `data_contratacao`
- **Plural** para nomes de tabelas: `funcionarios`, `obras`

## Comandos de Desenvolvimento

### Makefile

O projeto inclui um Makefile com comandos úteis:

```bash
# Rodar testes
make test

# Iniciar banco de dados
make up

# Parar banco de dados
make down

# Resetar banco de dados (apaga todos os dados)
make down-v
```

### Comandos Go Úteis

```bash
# Rodar testes verbosamente
go test -v ./...

# Rodar testes com coverage
go test -cover ./...

# Verificar código
go vet ./...

# Formatar código
go fmt ./...

# Atualizar dependências
go mod tidy

# Verificar dependências não utilizadas
go mod why <package>

# Build da aplicação
go build -o bin/server ./cmd/server

# Executar com flags específicas
go run ./cmd/server/main.go -debug

# Ver todas as dependências
go list -m all
```

## Desenvolvimento de Features

### Fluxo de Desenvolvimento

1. **Criar branch**: `git checkout -b feature/nova-funcionalidade`
2. **Implementar**: Seguir arquitetura limpa
3. **Testar**: Escrever testes unitários
4. **Documentar**: Atualizar documentação
5. **Pull Request**: Solicitar revisão

### Estrutura de uma Feature

#### 1. Entidade de Domínio (`internal/domain/`)

```go
// internal/domain/exemplo/entidade.go
package exemplo

import "time"

type MinhaEntidade struct {
    ID        string    `json:"id"`
    Nome      string    `json:"nome"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// Validações e métodos de negócio
func (e *MinhaEntidade) Validar() error {
    if e.Nome == "" {
        return errors.New("nome é obrigatório")
    }
    return nil
}
```

#### 2. Repository Interface (`internal/domain/`)

```go
// internal/domain/exemplo/repository.go
package exemplo

import "context"

type Repository interface {
    Salvar(ctx context.Context, entidade *MinhaEntidade) error
    BuscarPorID(ctx context.Context, id string) (*MinhaEntidade, error)
    Listar(ctx context.Context) ([]*MinhaEntidade, error)
    Deletar(ctx context.Context, id string) error
}
```

#### 3. Service (`internal/service/`)

```go
// internal/service/exemplo/service.go
package exemplo

import (
    "context"
    "github.com/luiszkm/masterCostrutora/internal/domain/exemplo"
)

type Service struct {
    repo   exemplo.Repository
    logger *slog.Logger
}

func NovoServico(repo exemplo.Repository, logger *slog.Logger) *Service {
    return &Service{
        repo:   repo,
        logger: logger,
    }
}

func (s *Service) Criar(ctx context.Context, input CriarInput) (*exemplo.MinhaEntidade, error) {
    entidade := &exemplo.MinhaEntidade{
        ID:   uuid.New().String(),
        Nome: input.Nome,
        // ... outros campos
    }
    
    if err := entidade.Validar(); err != nil {
        return nil, err
    }
    
    err := s.repo.Salvar(ctx, entidade)
    if err != nil {
        s.logger.Error("erro ao salvar entidade", "erro", err)
        return nil, err
    }
    
    return entidade, nil
}
```

#### 4. Repository Implementation (`internal/infrastructure/`)

```go
// internal/infrastructure/repository/postgres/exemplo_repository.go
package postgres

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/luiszkm/masterCostrutora/internal/domain/exemplo"
)

type exemploRepository struct {
    db     *pgxpool.Pool
    logger *slog.Logger
}

func NovoExemploRepository(db *pgxpool.Pool, logger *slog.Logger) exemplo.Repository {
    return &exemploRepository{
        db:     db,
        logger: logger,
    }
}

func (r *exemploRepository) Salvar(ctx context.Context, e *exemplo.MinhaEntidade) error {
    query := `
        INSERT INTO minha_tabela (id, nome, status, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())`
    
    _, err := r.db.Exec(ctx, query, e.ID, e.Nome, e.Status)
    return err
}
```

#### 5. HTTP Handler (`internal/handler/http/`)

```go
// internal/handler/http/exemplo/handler.go
package exemplo

import (
    "encoding/json"
    "net/http"
    "github.com/luiszkm/masterCostrutora/internal/handler/web"
)

type Handler struct {
    service Service
    logger  *slog.Logger
}

func NovoHandler(service Service, logger *slog.Logger) *Handler {
    return &Handler{
        service: service,
        logger:  logger,
    }
}

func (h *Handler) HandleCriar(w http.ResponseWriter, r *http.Request) {
    var input CriarInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        web.RespondError(w, r, "PAYLOAD_INVALIDO", "Payload inválido", http.StatusBadRequest)
        return
    }
    
    entidade, err := h.service.Criar(r.Context(), input)
    if err != nil {
        h.logger.Error("erro ao criar entidade", "erro", err)
        web.RespondError(w, r, "ERRO_CRIAR", "Erro ao criar", http.StatusInternalServerError)
        return
    }
    
    web.Respond(w, r, entidade, http.StatusCreated)
}
```

### Adicionando Rotas

No arquivo `internal/handler/http/router/router.go`:

```go
// Adicionar no grupo de rotas protegidas
r.With(auth.Authorize(authz.PermissaoExemploLer)).
    Get("/exemplo", exemploHandler.HandleListar)

r.With(auth.Authorize(authz.PermissaoExemploEscrever)).
    Post("/exemplo", exemploHandler.HandleCriar)
```

## Testes

### Estrutura de Testes

```
internal/
├── service/
│   ├── exemplo/
│   │   ├── service.go
│   │   └── service_test.go
├── infrastructure/
│   ├── repository/
│   │   └── postgres/
│   │       ├── exemplo_repository.go
│   │       └── exemplo_repository_test.go
```

### Teste Unitário de Service

```go
// internal/service/exemplo/service_test.go
package exemplo

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock do repository
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Salvar(ctx context.Context, e *exemplo.MinhaEntidade) error {
    args := m.Called(ctx, e)
    return args.Error(0)
}

func TestService_Criar(t *testing.T) {
    // Arrange
    mockRepo := new(MockRepository)
    service := NovoServico(mockRepo, slog.Default())
    
    input := CriarInput{Nome: "Teste"}
    
    mockRepo.On("Salvar", mock.Anything, mock.AnythingOfType("*exemplo.MinhaEntidade")).
        Return(nil)
    
    // Act
    resultado, err := service.Criar(context.Background(), input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, resultado)
    assert.Equal(t, "Teste", resultado.Nome)
    mockRepo.AssertExpectations(t)
}
```

### Teste de Integração (Repository)

```go
// internal/infrastructure/repository/postgres/exemplo_repository_test.go
package postgres

import (
    "context"
    "testing"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestExemploRepository_Salvar(t *testing.T) {
    // Configurar banco de teste
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NovoExemploRepository(db, slog.Default())
    
    entidade := &exemplo.MinhaEntidade{
        ID:   "test-id",
        Nome: "Teste",
    }
    
    // Act
    err := repo.Salvar(context.Background(), entidade)
    
    // Assert
    require.NoError(t, err)
    
    // Verificar se foi salvo
    encontrada, err := repo.BuscarPorID(context.Background(), "test-id")
    require.NoError(t, err)
    assert.Equal(t, "Teste", encontrada.Nome)
}
```

### Executar Testes

```bash
# Todos os testes
go test ./...

# Testes específicos
go test ./internal/service/exemplo/

# Com coverage
go test -cover ./internal/service/exemplo/

# Testes específicos com verbosidade
go test -v ./internal/service/exemplo/ -run TestService_Criar

# Testes de integração (requer banco)
go test -tags=integration ./internal/infrastructure/repository/postgres/
```

## Debugging

### VS Code Launch Configuration

Crie `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server",
            "env": {
                "DATABASE_URL": "postgres://user:password@localhost:5432/mastercostrutora_db?sslmode=disable",
                "JWT_SECRET_KEY": "dev-secret-key",
                "APP_ENV": "development"
            },
            "console": "integratedTerminal"
        }
    ]
}
```

### Logging

O sistema usa `slog` para logging estruturado:

```go
// Logs de diferentes níveis
logger.Debug("informação de debug", "usuario", userID)
logger.Info("operação realizada", "id", entityID)
logger.Warn("possível problema", "erro", err)
logger.Error("erro crítico", "erro", err, "contexto", context)

// Log com contexto
logger.With("component", "UserService").Info("usuário criado", "id", userID)
```

### Profiling

Para análise de performance:

```go
import _ "net/http/pprof"

// No main.go, adicionar:
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

Acessar `http://localhost:6060/debug/pprof/` para visualizar profiles.

## Banco de Dados

### Conexão Manual

```bash
# Conectar ao banco
docker exec -it masterconstrutora-db psql -U user -d mastercostrutora_db

# Ou usando cliente externo
psql -h localhost -p 5432 -U user -d mastercostrutora_db
```

### Queries Úteis

```sql
-- Ver todas as tabelas
\dt

-- Descrever uma tabela
\d usuarios

-- Ver dados de uma tabela
SELECT * FROM usuarios LIMIT 5;

-- Ver conexões ativas
SELECT count(*) FROM pg_stat_activity WHERE state = 'active';

-- Ver tamanho das tabelas
SELECT schemaname, tablename, 
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public' 
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### Backup e Restore

```bash
# Backup
docker exec masterconstrutora-db pg_dump -U user mastercostrutora_db > backup.sql

# Restore
docker exec -i masterconstrutora-db psql -U user mastercostrutora_db < backup.sql

# Reset completo do banco
make down-v
make up
```

## Migrations (Futuras)

Para quando implementarmos migrations automáticas:

```bash
# Instalar golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Criar nova migration
migrate create -ext sql -dir db/migrations -seq add_new_table

# Aplicar migrations
migrate -path db/migrations -database "postgres://user:password@localhost:5432/mastercostrutora_db?sslmode=disable" up

# Reverter última migration
migrate -path db/migrations -database "postgres://user:password@localhost:5432/mastercostrutora_db?sslmode=disable" down 1
```

## Boas Práticas

### 1. Estrutura de Commits

```bash
# Formato: tipo(escopo): descrição

feat(obras): adicionar endpoint para listar etapas
fix(auth): corrigir validação de token expirado
docs(api): atualizar documentação de autenticação
refactor(pessoal): extrair lógica de cálculo para service
test(suprimentos): adicionar testes para orçamentos
```

### 2. Code Review

Pontos importantes para revisão:
- **Arquitetura**: Seguir Clean Architecture
- **Testes**: Cobertura adequada
- **Logs**: Logs estruturados e úteis
- **Erros**: Tratamento adequado de erros
- **Segurança**: Validação de entrada, autorização
- **Performance**: Queries otimizadas, índices

### 3. Documentação

Manter atualizado:
- **README.md**: Instruções básicas
- **API.md**: Documentação da API
- **Comentários**: Código complexo deve ser comentado
- **Testes**: Exemplos de uso nos testes

### 4. Configuração de Produção

Diferenças entre desenvolvimento e produção:
- **Logs**: JSON em produção, texto em desenvolvimento
- **CORS**: Configurar origins específicos
- **HTTPS**: Obrigatório em produção
- **Timeouts**: Configurar timeouts apropriados
- **Rate Limiting**: Implementar limitação de requisições

## Troubleshooting

### Problemas Comuns

#### 1. Erro de Conexão com Banco
```bash
# Verificar se o container está rodando
docker-compose ps

# Ver logs do banco
docker-compose logs db

# Reiniciar banco
docker-compose restart db
```

#### 2. Porta 8080 já em uso
```bash
# Verificar que processo está usando a porta
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Matar processo (se necessário)
kill -9 <PID>
```

#### 3. Módulos Go não encontrados
```bash
# Limpar cache e redownload
go clean -modcache
go mod download
```

#### 4. Testes falhando
```bash
# Rodar testes com mais verbosidade
go test -v ./internal/service/exemplo/

# Verificar se banco de teste está configurado
# Usar tags para separar testes de integração
go test -tags=integration ./...
```

### Logs de Debug

Para debugging avançado:

```go
// No main.go, para ambiente de desenvolvimento
if os.Getenv("APP_ENV") == "development" {
    logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
}
```

### Performance

```bash
# Verificar goroutines
curl http://localhost:6060/debug/pprof/goroutine?debug=1

# Verificar memória
curl http://localhost:6060/debug/pprof/heap?debug=1

# CPU profiling
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```