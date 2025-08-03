# 🐛 Ambiente de Debug para Go - Master Construtora

## 📋 Configurações Criadas

### VS Code Debug (`.vscode/launch.json`)
- **Debug Server**: Debug do servidor principal (`cmd/server`)
- **Debug Seeder**: Debug do seeder (`cmd/seeder`) 
- **Debug Tests**: Debug de testes com verbosidade
- **Debug Current File**: Debug do arquivo atual
- **Attach to Process**: Anexar a processo em execução

### Configurações VS Code (`.vscode/settings.json`)
- Delve configurado com API v2
- Build/lint automático ao salvar
- Organização automática de imports
- Suporte completo ao Go Language Server

### Tasks (`.vscode/tasks.json`)
- Build do servidor
- Execução de testes
- Start/Stop do banco de dados
- Debug manual com Delve

## 🚀 Como Usar

### 1. Debug via VS Code (Recomendado)
1. Abra o VS Code no projeto
2. Vá para Run and Debug (Ctrl+Shift+D)
3. Selecione "Debug Server" 
4. Aperte F5 ou clique em "Start Debugging"

### 2. Debug Manual com Delve
```bash
# Terminal 1: Iniciar servidor em modo debug
dlv debug ./cmd/server --headless --listen=:2345 --api-version=2

# Terminal 2: Conectar ao debugger
dlv connect :2345
```

### 3. Debug de Testes
```bash
# Debug de todos os testes
dlv test ./...

# Debug de teste específico
dlv test ./internal/service/dashboard -- -test.run TestDashboardService
```

## 🎯 Breakpoints e Debug

### Comandos Delve Úteis
- `break main.main` - Breakpoint na função main
- `break arquivo.go:123` - Breakpoint na linha 123
- `continue` ou `c` - Continuar execução
- `next` ou `n` - Próxima linha
- `step` ou `s` - Entrar na função
- `locals` - Variáveis locais
- `print variavel` - Imprimir variável
- `goroutines` - Listar goroutines
- `quit` - Sair do debugger

### Breakpoints Sugeridos
- `cmd/server/main.go:43` - Início da aplicação
- `internal/handler/http/router/router.go` - Setup de rotas
- `internal/service/*/service.go` - Lógica de negócio
- `internal/infrastructure/repository/postgres/*.go` - Queries de banco

## 🔧 Preparação do Ambiente

### 1. Banco de Dados
```bash
# Iniciar PostgreSQL
docker-compose up -d

# Verificar se está rodando
docker-compose ps
```

### 2. Variáveis de Ambiente
Certifique-se que o arquivo `.env` existe com:
```env
DATABASE_URL=postgres://user:password@localhost:5432/master_construtora?sslmode=disable
JWT_SECRET_KEY=seu-jwt-secret-aqui
```

### 3. Dependências
```bash
# Instalar/atualizar dependências
go mod tidy

# Verificar se Delve está instalado
dlv version
```

## 🏗️ Estrutura de Debug

### Pontos de Entrada Principais
- `cmd/server/main.go:43` - Início da aplicação
- `internal/handler/http/router/router.go` - Configuração de rotas
- `internal/service/*/service.go` - Serviços de negócio

### Módulos para Debug
- **Identidade**: Autenticação e autorização
- **Obras**: Gestão de projetos e etapas
- **Pessoal**: Funcionários e apontamentos
- **Suprimentos**: Fornecedores e orçamentos
- **Financeiro**: Pagamentos e registros
- **Dashboard**: Métricas e relatórios

## 💡 Dicas de Debug

1. **Use logs estruturados**: O projeto já tem logging configurado
2. **Monitore goroutines**: Use `goroutines` para identificar vazamentos
3. **Debug de banco**: Coloque breakpoints nos repositórios
4. **Teste endpoints**: Use os arquivos `.http` na pasta `restclient/`
5. **Variables watch**: Configure watches no VS Code para variáveis importantes

## 🆘 Troubleshooting

### Delve não funciona
```bash
# Reinstalar Delve
go install github.com/go-delve/delve/cmd/dlv@latest
```

### Build falha
```bash
# Limpar cache e rebuildar
go clean -cache
go mod tidy
go build ./cmd/server
```

### Debug não para nos breakpoints
- Certifique-se que está compilando com símbolos de debug
- Verifique se o arquivo `.env` está presente
- Confirme que o banco está rodando