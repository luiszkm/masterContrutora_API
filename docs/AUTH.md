# Autenticação e Autorização - Master Construtora

## Visão Geral

O sistema Master Construtora implementa um sistema robusto de autenticação e autorização baseado em:

- **JWT (JSON Web Tokens)** para autenticação
- **RBAC (Role-Based Access Control)** para autorização
- **Cookies httpOnly** para armazenamento seguro de tokens
- **Permissões granulares** para controle fino de acesso

## Arquitetura de Segurança

### Fluxo de Autenticação

```
1. Usuário envia credenciais (email/senha)
2. Sistema valida credenciais
3. Sistema gera JWT com permissões do usuário
4. JWT é enviado em cookie httpOnly
5. Cliente armazena cookie automaticamente
6. Requisições subsequentes incluem cookie automaticamente
7. Middleware valida JWT em cada requisição
```

### Componentes Principais

#### 1. JWTService (`pkg/auth/jwt.go`)
Responsável por:
- Geração de tokens JWT
- Validação de tokens
- Middleware de autenticação
- Extração de claims

#### 2. Sistema de Permissões (`internal/authz/roles.go`)
Define:
- Permissões granulares
- Papéis (roles) e suas permissões
- Mapeamento de permissões por papel

#### 3. Password Hasher (`pkg/security/password.go`)
Gerencia:
- Hash de senhas com bcrypt
- Validação de senhas
- Configurações de segurança

## Implementação da Autenticação

### Registro de Usuário

**Endpoint**: `POST /usuarios/registrar`

```go
type registrarRequest struct {
    Nome           string `json:"nome"`
    Email          string `json:"email"`
    Senha          string `json:"senha"`
    ConfirmarSenha string `json:"confirmarSenha"`
}
```

**Fluxo:**
1. Validação dos dados de entrada
2. Verificação se email já existe
3. Hash da senha com bcrypt
4. Criação do usuário no banco
5. Retorno dos dados do usuário (sem senha)

### Login de Usuário

**Endpoint**: `POST /usuarios/login`

```go
type loginRequest struct {
    Email string `json:"email"`
    Senha string `json:"senha"`
}
```

**Fluxo:**
1. Busca usuário por email
2. Validação da senha com bcrypt
3. Geração do JWT com claims do usuário
4. Configuração do cookie httpOnly
5. Retorno do token e dados do usuário

### Configuração de Cookies

```go
cookie := http.Cookie{
    Name:     "jwt-token",
    Value:    tokenString,
    Expires:  time.Now().Add(time.Hour * 8),
    HttpOnly: true,                    // Impede acesso via JavaScript
    Secure:   isSecure,               // HTTPS em produção
    SameSite: http.SameSiteLaxMode,   // Proteção CSRF
    Path:     "/",                    // Válido para todo o site
}
```

### Estrutura do JWT

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "uuid-do-usuario",
    "email": "usuario@empresa.com",
    "permissions": [
      "obras:ler",
      "obras:escrever",
      "pessoal:ler"
    ],
    "iat": 1640995200,
    "exp": 1641024000
  }
}
```

## Sistema de Autorização

### Permissões Granulares

O sistema define permissões específicas para cada contexto:

```go
const (
    // Obras
    PermissaoObrasLer      = "obras:ler"
    PermissaoObrasEscrever = "obras:escrever"
    
    // Pessoal
    PermissaoPessoalLer                 = "pessoal:ler"
    PermissaoPessoalEscrever            = "pessoal:escrever"
    PermissaoPessoalApontamentoLer      = "pessoal:apontamento:ler"
    PermissaoPessoalApontamentoEscrever = "pessoal:apontamento:escrever"
    PermissaoPessoalApontamentoAprovar  = "pessoal:apontamento:aprovar"
    PermissaoPessoalApontamentoPagar    = "pessoal:apontamento:pagar"
    
    // Suprimentos
    PermissaoSuprimentosLer      = "suprimentos:ler"
    PermissaoSuprimentosEscrever = "suprimentos:escrever"
    
    // Financeiro
    PermissaoFinanceiroLer                = "financeiro:ler"
    PermissaoFinanceiroEscrever           = "financeiro:escrever"
    PermissaoFinanceiroContasReceberLer   = "financeiro:contas_receber:ler"
    PermissaoFinanceiroContasReceberEscrever = "financeiro:contas_receber:escrever"
    PermissaoFinanceiroContasPagarLer     = "financeiro:contas_pagar:ler"
    PermissaoFinanceiroContasPagarEscrever   = "financeiro:contas_pagar:escrever"
    PermissaoFinanceiroCronogramaLer      = "financeiro:cronograma:ler"
    PermissaoFinanceiroCronogramaEscrever = "financeiro:cronograma:escrever"
)
```

### Papéis (Roles)

#### ADMIN
Acesso completo a todas as funcionalidades:
- Todas as permissões do sistema
- Capacidade de gerenciar usuários
- Acesso a todas as operações

#### GERENTE_OBRAS
Gerenciamento completo de obras e recursos:
```go
{
    PermissaoObrasLer,
    PermissaoObrasEscrever,
    PermissaoPessoalLer,
    PermissaoPessoalEscrever,
    PermissaoSuprimentosLer,
    PermissaoSuprimentosEscrever,
    PermissaoFinanceiroLer,
    PermissaoFinanceiroEscrever,
    PermissaoPessoalApontamentoEscrever,
    PermissaoPessoalApontamentoAprovar,
    PermissaoPessoalApontamentoPagar,
}
```

#### VISUALIZADOR
Acesso somente leitura:
```go
{
    PermissaoObrasLer,
    PermissaoPessoalLer,
    PermissaoSuprimentosLer,
    PermissaoFinanceiroLer,
    PermissaoPessoalApontamentoLer,
}
```

### Middleware de Autorização

#### AuthMiddleware
Valida a presença e validade do JWT:

```go
func (j *JWTService) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Extrai token do cookie
        cookie, err := r.Cookie("jwt-token")
        if err != nil {
            // Token não encontrado
            http.Error(w, "Token não encontrado", http.StatusUnauthorized)
            return
        }
        
        // 2. Valida o token
        claims, err := j.ValidateToken(cookie.Value)
        if err != nil {
            // Token inválido
            http.Error(w, "Token inválido", http.StatusUnauthorized)
            return
        }
        
        // 3. Adiciona claims ao contexto
        ctx := context.WithValue(r.Context(), "userClaims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

#### Authorize Function
Verifica permissões específicas:

```go
func Authorize(requiredPermission string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1. Extrai claims do contexto
            claims, ok := r.Context().Value("userClaims").(*Claims)
            if !ok {
                http.Error(w, "Claims não encontradas", http.StatusUnauthorized)
                return
            }
            
            // 2. Verifica se usuário tem a permissão
            if !hasPermission(claims.Permissions, requiredPermission) {
                http.Error(w, "Permissão insuficiente", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### Uso nos Endpoints

```go
// Exemplo de proteção de endpoint
r.With(auth.Authorize(authz.PermissaoObrasEscrever)).
    Post("/obras", obrasHandler.HandleCriarObra)

r.With(auth.Authorize(authz.PermissaoPessoalApontamentoAprovar)).
    Patch("/apontamentos/{id}/aprovar", pessoalHandler.HandleAprovarApontamento)
```

## Implementação de Segurança

### Hash de Senhas

Utiliza bcrypt com custo adequado:

```go
type BcryptHasher struct {
    cost int
}

func (b *BcryptHasher) HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
    return string(hash), err
}

func (b *BcryptHasher) CheckPassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
```

### Configurações de Segurança

#### Desenvolvimento
```env
JWT_SECRET_KEY=chave-desenvolvimento-nao-usar-em-producao
APP_ENV=development
```

#### Produção
```env
JWT_SECRET_KEY=chave-super-secreta-256-bits-minimo
APP_ENV=production
```

### Tempo de Expiração

- **Desenvolvimento**: 8 horas
- **Produção**: 4 horas (recomendado)

## Fluxos de Integração Frontend

### 1. Login

```javascript
// Request
const response = await fetch('/usuarios/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'usuario@empresa.com',
    senha: 'senha123'
  }),
  credentials: 'include' // Importante para cookies
});

// Response
{
  "accessToken": "jwt-token-string",
  "userId": "uuid-do-usuario"
}
```

### 2. Requisições Autenticadas

```javascript
// O cookie é enviado automaticamente
const response = await fetch('/obras', {
  method: 'GET',
  credentials: 'include' // Inclui cookies automaticamente
});
```

### 3. Logout (Frontend)

```javascript
// Remove o cookie (expira imediatamente)
document.cookie = 'jwt-token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
```

### 4. Verificação de Permissões no Frontend

```javascript
// Decode do JWT para verificar permissões (opcional)
function hasPermission(token, permission) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.permissions.includes(permission);
  } catch (e) {
    return false;
  }
}

// Uso
if (hasPermission(token, 'obras:escrever')) {
  // Mostrar botão de criar obra
}
```

## Tratamento de Erros

### Códigos de Resposta

#### 401 Unauthorized
- Token não fornecido
- Token expirado
- Token inválido
- Credenciais incorretas

```json
{
  "erro": {
    "codigo": "TOKEN_INVALIDO",
    "mensagem": "Token JWT inválido ou expirado"
  }
}
```

#### 403 Forbidden
- Usuário autenticado mas sem permissão
- Tentativa de acesso a recurso restrito

```json
{
  "erro": {
    "codigo": "ACESSO_NEGADO", 
    "mensagem": "Usuário não possui permissão para esta operação"
  }
}
```

### Renovação de Token

Atualmente não há renovação automática. Quando o token expira:
1. Frontend recebe 401
2. Redireciona para login
3. Usuário faz login novamente

**Implementação futura**: Refresh tokens para renovação automática.

## Boas Práticas de Segurança

### Servidor

1. **HTTPS obrigatório em produção**
2. **Chave JWT forte (mínimo 256 bits)**
3. **Cookies httpOnly sempre habilitados**
4. **SameSite configurado adequadamente**
5. **Logs de tentativas de acesso**
6. **Rate limiting em endpoints de login**

### Frontend

1. **Nunca armazenar JWT em localStorage**
2. **Usar credentials: 'include' em todas as requisições**
3. **Verificar permissões antes de exibir funcionalidades**
4. **Implementar logout adequado**
5. **Tratar erros 401/403 adequadamente**

### Banco de Dados

1. **Senhas sempre hasheadas**
2. **Não armazenar JWTs no banco**
3. **Logs de tentativas de login**
4. **Índices em campos de busca de usuários**

## Auditoria e Monitoramento

### Logs de Segurança

```go
// Login bem-sucedido
logger.InfoContext(ctx, "usuário logado com sucesso", 
    "email", email, 
    "userId", userID,
    "ip", request.RemoteAddr)

// Tentativa de login falhou
logger.WarnContext(ctx, "tentativa de login falhou", 
    "email", email, 
    "erro", err,
    "ip", request.RemoteAddr)

// Acesso negado
logger.WarnContext(ctx, "acesso negado",
    "userId", userID,
    "endpoint", request.URL.Path,
    "permissaoNecessaria", requiredPermission)
```

### Métricas Importantes

- **Tentativas de login por IP**
- **Tokens expirados vs renovados**
- **Acessos negados por endpoint**
- **Usuários únicos por período**
- **Distribuição de permissões utilizadas**

## Configuração de Desenvolvimento

### Variáveis de Ambiente

```env
# Desenvolvimento
JWT_SECRET_KEY=dev-secret-key-change-in-production
APP_ENV=development
DATABASE_URL=postgres://user:password@localhost:5432/mastercostrutora_db?sslmode=disable
```

### Usuário Padrão

Para desenvolvimento, criar usuário admin:

```json
{
  "nome": "Admin Sistema",
  "email": "admin@construtora.com", 
  "senha": "admin123",
  "permissoes": ["admin"] // Papel ADMIN tem todas as permissões
}
```

## Troubleshooting

### Problemas Comuns

#### "Token não encontrado"
- Verificar se o cookie está sendo enviado
- Verificar configuração de credentials no frontend
- Verificar domínio do cookie

#### "Token inválido"
- Verificar se JWT_SECRET_KEY é a mesma
- Verificar se token não expirou
- Verificar formato do token

#### "Permissão insuficiente"
- Verificar permissões do usuário no banco
- Verificar se endpoint requer permissão correta
- Verificar se middleware está aplicado corretamente

### Comandos Úteis

```bash
# Verificar estrutura do JWT (apenas payload)
echo "TOKEN_JWT" | cut -d. -f2 | base64 -d | jq

# Verificar cookies no navegador
# DevTools > Application > Cookies

# Testar endpoint com curl
curl -X POST http://localhost:8080/usuarios/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@construtora.com","senha":"admin123"}' \
  -c cookies.txt

curl -X GET http://localhost:8080/obras \
  -b cookies.txt
```