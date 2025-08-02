# Guia de Integração Frontend - Master Construtora

## Visão Geral

Este guia fornece todas as informações necessárias para desenvolver um frontend que se integre com a API Master Construtora. A API é projetada para ser consumida por aplicações React, Vue, Angular ou qualquer framework que suporte requisições HTTP.

## Configuração Inicial

### Base URL
```javascript
const API_BASE_URL = 'http://localhost:8080'; // Desenvolvimento
// const API_BASE_URL = 'https://api.construtora.com'; // Produção
```

### Headers Padrão
```javascript
const defaultHeaders = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
};
```

### Configuração de Cookies
**IMPORTANTE**: Sempre incluir `credentials: 'include'` para enviar cookies de autenticação:

```javascript
const apiRequest = async (url, options = {}) => {
  return fetch(`${API_BASE_URL}${url}`, {
    ...options,
    credentials: 'include', // Essencial para autenticação
    headers: {
      ...defaultHeaders,
      ...options.headers
    }
  });
};
```

## Sistema de Autenticação

### 1. Login

```javascript
const login = async (email, senha) => {
  try {
    const response = await apiRequest('/usuarios/login', {
      method: 'POST',
      body: JSON.stringify({ email, senha })
    });
    
    if (!response.ok) {
      throw new Error('Credenciais inválidas');
    }
    
    const data = await response.json();
    
    // Cookie é automaticamente armazenado pelo navegador
    // Opcionalmente, armazenar dados do usuário
    localStorage.setItem('user', JSON.stringify({
      userId: data.userId,
      token: data.accessToken // Para verificação de permissões
    }));
    
    return data;
  } catch (error) {
    console.error('Erro no login:', error);
    throw error;
  }
};
```

### 2. Verificação de Autenticação

```javascript
const checkAuth = async () => {
  try {
    const response = await apiRequest('/health', {
      method: 'GET'
    });
    
    // Se retornar 401, usuário não está autenticado
    return response.ok;
  } catch (error) {
    return false;
  }
};
```

### 3. Logout

```javascript
const logout = () => {
  // Remove dados locais
  localStorage.removeItem('user');
  
  // Remove cookie (expira imediatamente)
  document.cookie = 'jwt-token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; SameSite=Lax';
  
  // Redirecionar para login
  window.location.href = '/login';
};
```

### 4. Verificação de Permissões

```javascript
// Função para decodificar JWT (apenas payload, não validação)
const decodeJWT = (token) => {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch (error) {
    return null;
  }
};

// Verificar se usuário tem permissão específica
const hasPermission = (permission) => {
  const user = JSON.parse(localStorage.getItem('user') || '{}');
  if (!user.token) return false;
  
  const payload = decodeJWT(user.token);
  return payload?.permissions?.includes(permission) || false;
};

// Exemplo de uso
if (hasPermission('obras:escrever')) {
  // Mostrar botão de criar obra
}
```

## Interceptador HTTP

### Configuração com Axios

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // Importante para cookies
  headers: defaultHeaders
});

// Interceptador de resposta para tratar erros de auth
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expirado ou inválido
      logout();
      return Promise.reject(error);
    }
    
    if (error.response?.status === 403) {
      // Sem permissão
      console.error('Acesso negado:', error.response.data);
      // Mostrar mensagem de erro ou redirecionar
    }
    
    return Promise.reject(error);
  }
);

export default api;
```

### Configuração com Fetch

```javascript
const apiRequest = async (url, options = {}) => {
  try {
    const response = await fetch(`${API_BASE_URL}${url}`, {
      credentials: 'include',
      headers: defaultHeaders,
      ...options,
      headers: { ...defaultHeaders, ...options.headers }
    });
    
    // Tratar erros de autenticação
    if (response.status === 401) {
      logout();
      throw new Error('Sessão expirada');
    }
    
    if (response.status === 403) {
      throw new Error('Acesso negado');
    }
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.erro?.mensagem || 'Erro na requisição');
    }
    
    return response;
  } catch (error) {
    console.error('Erro na API:', error);
    throw error;
  }
};
```

## Modelos de Dados (TypeScript)

### Entidades Principais

```typescript
// Obras
interface Obra {
  id: string;
  nome: string;
  cliente: string;
  endereco: string;
  descricao?: string;
  dataInicio: string;
  dataFim?: string;
  status: 'Em Planejamento' | 'Em Andamento' | 'Concluída' | 'Cancelada';
}

interface ObraDashboard {
  obraId: string;
  nomeObra: string;
  statusObra: string;
  etapaAtualNome?: string;
  dataFimPrevistaEtapa?: string;
  diasParaPrazoEtapa?: number;
  percentualConcluido: number;
  custoTotalRealizado: number;
  orcamentoTotalAprovado: number;
  balancoFinanceiro: number;
  funcionariosAlocados: number;
  ultimaAtualizacao: string;
}

interface Etapa {
  id: string;
  obraId: string;
  nome: string;
  dataInicioPrevista?: string;
  dataFimPrevista?: string;
  status: 'Não Iniciada' | 'Em Andamento' | 'Concluída';
}

// Pessoal
interface Funcionario {
  id: string;
  nome: string;
  cpf: string;
  cargo: string;
  departamento: string;
  telefone: string;
  email: string;
  dataContratacao: string;
  valorDiaria: number;
  chavePix: string;
  status: 'Ativo' | 'Inativo' | 'Desligado';
}

interface ApontamentoQuinzenal {
  id: string;
  funcionarioId: string;
  obraId: string;
  periodoInicio: string;
  periodoFim: string;
  diaria: number;
  diasTrabalhados: number;
  adicionais: number;
  descontos: number;
  adiantamentos: number;
  valorTotalCalculado: number;
  status: 'Em Aberto' | 'Aprovado para Pagamento' | 'Pago';
  createdAt: string;
  updatedAt: string;
  nomeFuncionario?: string;
}

// Suprimentos
interface Fornecedor {
  id: string;
  nome: string;
  cnpj: string;
  contato?: string;
  email?: string;
  endereco?: string;
  status: 'Ativo' | 'Inativo';
  avaliacao?: number;
  observacoes?: string;
}

interface Produto {
  id: string;
  nome: string;
  descricao?: string;
  unidadeDeMedida: string;
  categoria: string;
}

interface Orcamento {
  id: string;
  numero: string;
  etapaId: string;
  fornecedorId: string;
  valorTotal: number;
  status: 'Em Aberto' | 'Aprovado' | 'Rejeitado' | 'Pago';
  dataEmissao: string;
  dataAprovacao?: string;
  observacoes?: string;
  condicoesPagamento?: string;
  itens: ItemOrcamento[];
}

interface ItemOrcamento {
  ProdutoNome: string;
  UnidadeDeMedida: string;
  Categoria: string;
  Quantidade: number;
  ValorUnitario: number;
}

// Resposta paginada
interface RespostaPaginada<T> {
  dados: T[];
  paginacao: {
    paginaAtual: number;
    totalPaginas: number;
    totalItens: number;
    itensPorPagina: number;
  };
}

// Erro padrão da API
interface ErroAPI {
  erro: {
    codigo: string;
    mensagem: string;
    detalhes?: string;
  };
}
```

### DTOs de Input

```typescript
// Obras
interface CriarObraInput {
  nome: string;
  cliente: string;
  endereco: string;
  descricao?: string;
  dataInicio: string;
  etapasPadrao?: string[];
}

interface AtualizarObraInput {
  nome?: string;
  cliente?: string;
  endereco?: string;
  descricao?: string;
  dataFim?: string;
  status?: string;
}

// Pessoal
interface CadastrarFuncionarioInput {
  nome: string;
  cpf: string;
  cargo: string;
  departamento: string;
  telefone: string;
  email: string;
  diaria: number;
  chavePix: string;
}

interface CriarApontamentoInput {
  funcionarioId: string;
  obraId: string;
  periodoInicio: string;
  periodoFim: string;
  Diaria: number;
  DiasTrabalhados: number;
  Descontos: number;
  Adiantamento: number;
  ValorAdicional: number;
}

// Suprimentos
interface CadastrarFornecedorInput {
  nome: string;
  cnpj: string;
  contato?: string;
  email?: string;
  endereco?: string;
  observacoes?: string;
}

interface CriarOrcamentoInput {
  fornecedorId: string;
  itens: {
    nomeProduto: string;
    unidadeDeMedida: string;
    categoria: string;
    quantidade: number;
    valorUnitario: number;
  }[];
}
```

## Serviços da API

### Serviço de Obras

```typescript
class ObrasService {
  // Listar obras
  static async listar(filtros: {
    page?: number;
    limit?: number;
    nome?: string;
    status?: string;
  } = {}): Promise<RespostaPaginada<Obra>> {
    const params = new URLSearchParams();
    Object.entries(filtros).forEach(([key, value]) => {
      if (value) params.append(key, value.toString());
    });
    
    const response = await apiRequest(`/obras?${params}`);
    return response.json();
  }
  
  // Criar obra
  static async criar(data: CriarObraInput): Promise<Obra> {
    const response = await apiRequest('/obras', {
      method: 'POST',
      body: JSON.stringify(data)
    });
    return response.json();
  }
  
  // Buscar dashboard
  static async buscarDashboard(obraId: string): Promise<ObraDashboard> {
    const response = await apiRequest(`/obras/${obraId}/dashboard`);
    return response.json();
  }
  
  // Buscar detalhes
  static async buscarDetalhes(obraId: string): Promise<Obra> {
    const response = await apiRequest(`/obras/${obraId}`);
    return response.json();
  }
  
  // Atualizar obra
  static async atualizar(obraId: string, data: AtualizarObraInput): Promise<void> {
    await apiRequest(`/obras/${obraId}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    });
  }
  
  // Deletar obra
  static async deletar(obraId: string): Promise<void> {
    await apiRequest(`/obras/${obraId}`, {
      method: 'DELETE'
    });
  }
  
  // Listar etapas da obra
  static async listarEtapas(obraId: string): Promise<Etapa[]> {
    const response = await apiRequest(`/obras/${obraId}/etapas`);
    return response.json();
  }
  
  // Adicionar etapa
  static async adicionarEtapa(obraId: string, etapa: {
    nome: string;
    dataInicioPrevista?: string;
    dataFimPrevista?: string;
  }): Promise<Etapa> {
    const response = await apiRequest(`/obras/${obraId}/etapas`, {
      method: 'POST',
      body: JSON.stringify(etapa)
    });
    return response.json();
  }
  
  // Alocar funcionários
  static async alocarFuncionarios(obraId: string, data: {
    funcionarioIds: string[];
    dataInicioAlocacao: string;
  }): Promise<any[]> {
    const response = await apiRequest(`/obras/${obraId}/alocacoes`, {
      method: 'POST',
      body: JSON.stringify(data)
    });
    return response.json();
  }
}
```

### Serviço de Pessoal

```typescript
class PessoalService {
  // Listar funcionários
  static async listarFuncionarios(): Promise<Funcionario[]> {
    const response = await apiRequest('/funcionarios');
    return response.json();
  }
  
  // Cadastrar funcionário
  static async cadastrar(data: CadastrarFuncionarioInput): Promise<Funcionario> {
    const response = await apiRequest('/funcionarios', {
      method: 'POST',
      body: JSON.stringify(data)
    });
    return response.json();
  }
  
  // Buscar funcionário
  static async buscar(funcionarioId: string): Promise<Funcionario> {
    const response = await apiRequest(`/funcionarios/${funcionarioId}`);
    return response.json();
  }
  
  // Listar apontamentos
  static async listarApontamentos(filtros: {
    page?: number;
    limit?: number;
    status?: string;
    funcionarioId?: string;
  } = {}): Promise<RespostaPaginada<ApontamentoQuinzenal>> {
    const params = new URLSearchParams();
    Object.entries(filtros).forEach(([key, value]) => {
      if (value) params.append(key, value.toString());
    });
    
    const response = await apiRequest(`/apontamentos?${params}`);
    return response.json();
  }
  
  // Criar apontamento
  static async criarApontamento(data: CriarApontamentoInput): Promise<ApontamentoQuinzenal> {
    const response = await apiRequest('/apontamentos', {
      method: 'POST',
      body: JSON.stringify(data)
    });
    return response.json();
  }
  
  // Aprovar apontamento
  static async aprovarApontamento(apontamentoId: string): Promise<ApontamentoQuinzenal> {
    const response = await apiRequest(`/apontamentos/${apontamentoId}/aprovar`, {
      method: 'PATCH'
    });
    return response.json();
  }
  
  // Pagar apontamento
  static async pagarApontamento(apontamentoId: string, contaBancariaId: string): Promise<ApontamentoQuinzenal> {
    const response = await apiRequest(`/apontamentos/${apontamentoId}/pagar`, {
      method: 'PATCH',
      body: JSON.stringify({
        contaBancariaId,
        apontamentoId: [apontamentoId]
      })
    });
    return response.json();
  }
}
```

### Serviço de Suprimentos

```typescript
class SuprimentosService {
  // Listar fornecedores
  static async listarFornecedores(): Promise<Fornecedor[]> {
    const response = await apiRequest('/fornecedores');
    return response.json();
  }
  
  // Listar orçamentos
  static async listarOrcamentos(filtros: {
    page?: number;
    limit?: number;
    status?: string;
    fornecedorId?: string;
    obraId?: string;
  } = {}): Promise<RespostaPaginada<any>> {
    const params = new URLSearchParams();
    Object.entries(filtros).forEach(([key, value]) => {
      if (value) params.append(key, value.toString());
    });
    
    const response = await apiRequest(`/orcamentos?${params}`);
    return response.json();
  }
  
  // Buscar orçamento detalhado
  static async buscarOrcamento(orcamentoId: string): Promise<Orcamento> {
    const response = await apiRequest(`/orcamentos/${orcamentoId}`);
    return response.json();
  }
  
  // Criar orçamento
  static async criarOrcamento(etapaId: string, data: CriarOrcamentoInput): Promise<Orcamento> {
    const response = await apiRequest(`/etapas/${etapaId}/orcamentos`, {
      method: 'POST',
      body: JSON.stringify(data)
    });
    return response.json();
  }
  
  // Atualizar status do orçamento
  static async atualizarStatus(orcamentoId: string, status: string): Promise<Orcamento> {
    const response = await apiRequest(`/orcamentos/${orcamentoId}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status })
    });
    return response.json();
  }
}
```

## Componentes React de Exemplo

### Hook de Autenticação

```typescript
import { useState, useEffect, createContext, useContext } from 'react';

interface AuthContextType {
  user: any;
  login: (email: string, senha: string) => Promise<void>;
  logout: () => void;
  hasPermission: (permission: string) => boolean;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState(null);
  
  useEffect(() => {
    // Verificar se há usuário logado
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      setUser(JSON.parse(savedUser));
    }
  }, []);
  
  const login = async (email: string, senha: string) => {
    const data = await loginAPI(email, senha);
    setUser(data);
  };
  
  const logout = () => {
    localStorage.removeItem('user');
    document.cookie = 'jwt-token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    setUser(null);
  };
  
  const hasPermission = (permission: string) => {
    if (!user?.token) return false;
    const payload = decodeJWT(user.token);
    return payload?.permissions?.includes(permission) || false;
  };
  
  return (
    <AuthContext.Provider value={{
      user,
      login,
      logout,
      hasPermission,
      isAuthenticated: !!user
    }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth deve ser usado dentro de AuthProvider');
  }
  return context;
};
```

### Componente de Login

```typescript
import React, { useState } from 'react';
import { useAuth } from './AuthContext';

const LoginForm: React.FC = () => {
  const [email, setEmail] = useState('');
  const [senha, setSenha] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  const { login } = useAuth();
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    
    try {
      await login(email, senha);
      // Redirecionar após login bem-sucedido
    } catch (err) {
      setError('Credenciais inválidas');
    } finally {
      setLoading(false);
    }
  };
  
  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Email:</label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
      </div>
      
      <div>
        <label>Senha:</label>
        <input
          type="password"
          value={senha}
          onChange={(e) => setSenha(e.target.value)}
          required
        />
      </div>
      
      {error && <div className="error">{error}</div>}
      
      <button type="submit" disabled={loading}>
        {loading ? 'Entrando...' : 'Entrar'}
      </button>
    </form>
  );
};
```

### Componente Protegido

```typescript
import React from 'react';
import { useAuth } from './AuthContext';

interface ProtectedComponentProps {
  permission: string;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

const ProtectedComponent: React.FC<ProtectedComponentProps> = ({
  permission,
  children,
  fallback = <div>Você não tem permissão para ver este conteúdo.</div>
}) => {
  const { hasPermission } = useAuth();
  
  if (!hasPermission(permission)) {
    return <>{fallback}</>;
  }
  
  return <>{children}</>;
};

// Uso
const ObrasList: React.FC = () => {
  return (
    <div>
      <h1>Lista de Obras</h1>
      
      <ProtectedComponent permission="obras:escrever">
        <button>Criar Nova Obra</button>
      </ProtectedComponent>
      
      {/* Lista de obras aqui */}
    </div>
  );
};
```

## Tratamento de Erros

### Hook para Tratamento de Erros

```typescript
import { useState } from 'react';

interface ApiError {
  codigo: string;
  mensagem: string;
  detalhes?: string;
}

export const useApiError = () => {
  const [error, setError] = useState<ApiError | null>(null);
  
  const handleError = (err: any) => {
    if (err.response?.data?.erro) {
      setError(err.response.data.erro);
    } else {
      setError({
        codigo: 'ERRO_DESCONHECIDO',
        mensagem: err.message || 'Erro desconhecido'
      });
    }
  };
  
  const clearError = () => setError(null);
  
  return { error, handleError, clearError };
};
```

### Componente de Notificação

```typescript
import React from 'react';

interface NotificationProps {
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
  onClose: () => void;
}

const Notification: React.FC<NotificationProps> = ({ type, message, onClose }) => {
  return (
    <div className={`notification notification-${type}`}>
      <span>{message}</span>
      <button onClick={onClose}>×</button>
    </div>
  );
};
```

## Paginação

### Hook de Paginação

```typescript
import { useState, useEffect } from 'react';

interface UsePaginationProps<T> {
  fetchFunction: (page: number, limit: number, filters?: any) => Promise<RespostaPaginada<T>>;
  initialLimit?: number;
  filters?: any;
}

export const usePagination = <T>({
  fetchFunction,
  initialLimit = 10,
  filters = {}
}: UsePaginationProps<T>) => {
  const [data, setData] = useState<T[]>([]);
  const [pagination, setPagination] = useState({
    paginaAtual: 1,
    totalPaginas: 0,
    totalItens: 0,
    itensPorPagina: initialLimit
  });
  const [loading, setLoading] = useState(false);
  
  const fetchData = async (page: number = 1) => {
    setLoading(true);
    try {
      const response = await fetchFunction(page, pagination.itensPorPagina, filters);
      setData(response.dados);
      setPagination(response.paginacao);
    } catch (error) {
      console.error('Erro ao buscar dados:', error);
    } finally {
      setLoading(false);
    }
  };
  
  useEffect(() => {
    fetchData(1);
  }, [filters]);
  
  const goToPage = (page: number) => {
    fetchData(page);
  };
  
  const nextPage = () => {
    if (pagination.paginaAtual < pagination.totalPaginas) {
      fetchData(pagination.paginaAtual + 1);
    }
  };
  
  const prevPage = () => {
    if (pagination.paginaAtual > 1) {
      fetchData(pagination.paginaAtual - 1);
    }
  };
  
  return {
    data,
    pagination,
    loading,
    goToPage,
    nextPage,
    prevPage,
    refresh: () => fetchData(pagination.paginaAtual)
  };
};
```

## Validação de Formulários

### Esquemas de Validação (usando Zod)

```typescript
import { z } from 'zod';

export const obraSchema = z.object({
  nome: z.string().min(1, 'Nome é obrigatório').max(255, 'Nome muito longo'),
  cliente: z.string().min(1, 'Cliente é obrigatório').max(255, 'Nome muito longo'),
  endereco: z.string().min(1, 'Endereço é obrigatório'),
  descricao: z.string().optional(),
  dataInicio: z.string().refine((date) => !isNaN(Date.parse(date)), {
    message: 'Data de início inválida'
  })
});

export const funcionarioSchema = z.object({
  nome: z.string().min(1, 'Nome é obrigatório').max(255, 'Nome muito longo'),
  cpf: z.string().regex(/^\d{3}\.\d{3}\.\d{3}-\d{2}$/, 'CPF inválido'),
  cargo: z.string().min(1, 'Cargo é obrigatório'),
  departamento: z.string().min(1, 'Departamento é obrigatório'),
  telefone: z.string().min(1, 'Telefone é obrigatório'),
  email: z.string().email('Email inválido'),
  diaria: z.number().positive('Diária deve ser positiva'),
  chavePix: z.string().min(1, 'Chave PIX é obrigatória')
});

export type ObraFormData = z.infer<typeof obraSchema>;
export type FuncionarioFormData = z.infer<typeof funcionarioSchema>;
```

## Configuração de Build

### Variáveis de Ambiente

```bash
# .env.development
REACT_APP_API_URL=http://localhost:8080
REACT_APP_ENV=development

# .env.production
REACT_APP_API_URL=https://api.construtora.com
REACT_APP_ENV=production
```

### Configuração do Proxy (Create React App)

```json
// package.json
{
  "proxy": "http://localhost:8080"
}
```

## Testes

### Teste de Integração com API

```typescript
import { render, screen, waitFor } from '@testing-library/react';
import { rest } from 'msw';
import { setupServer } from 'msw/node';
import ObrasListComponent from './ObrasListComponent';

// Mock da API
const server = setupServer(
  rest.get('/obras', (req, res, ctx) => {
    return res(ctx.json({
      dados: [
        {
          id: '1',
          nome: 'Obra Teste',
          cliente: 'Cliente Teste',
          status: 'Em Andamento'
        }
      ],
      paginacao: {
        paginaAtual: 1,
        totalPaginas: 1,
        totalItens: 1,
        itensPorPagina: 10
      }
    }));
  })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

test('deve listar obras', async () => {
  render(<ObrasListComponent />);
  
  await waitFor(() => {
    expect(screen.getByText('Obra Teste')).toBeInTheDocument();
  });
});
```

## Considerações de Performance

### 1. Debounce em Filtros
```typescript
import { useMemo } from 'react';
import { debounce } from 'lodash';

const SearchComponent = () => {
  const debouncedSearch = useMemo(
    () => debounce((term: string) => {
      // Fazer busca na API
    }, 300),
    []
  );
  
  return (
    <input
      onChange={(e) => debouncedSearch(e.target.value)}
      placeholder="Buscar..."
    />
  );
};
```

### 2. Cache de Dados
```typescript
// Usando React Query
import { useQuery } from 'react-query';

const useObras = (filtros: any) => {
  return useQuery(
    ['obras', filtros],
    () => ObrasService.listar(filtros),
    {
      staleTime: 5 * 60 * 1000, // 5 minutos
      cacheTime: 10 * 60 * 1000 // 10 minutos
    }
  );
};
```

### 3. Lazy Loading
```typescript
import { lazy, Suspense } from 'react';

const ObrasModule = lazy(() => import('./modules/Obras'));
const PessoalModule = lazy(() => import('./modules/Pessoal'));

const App = () => {
  return (
    <Suspense fallback={<div>Carregando...</div>}>
      <Routes>
        <Route path="/obras/*" element={<ObrasModule />} />
        <Route path="/pessoal/*" element={<PessoalModule />} />
      </Routes>
    </Suspense>
  );
};
```

## Troubleshooting

### Problemas Comuns

1. **CORS Error**: Verificar se `credentials: 'include'` está sendo usado
2. **Cookie não enviado**: Verificar configuração de SameSite e Secure
3. **401 em todas as requisições**: Verificar se login foi feito corretamente
4. **403 em operações específicas**: Verificar permissões do usuário

### Debug de Requisições

```javascript
// Adicionar logs para debug
const apiRequest = async (url, options = {}) => {
  console.log('API Request:', { url, options });
  
  const response = await fetch(`${API_BASE_URL}${url}`, {
    credentials: 'include',
    ...options
  });
  
  console.log('API Response:', { 
    status: response.status, 
    headers: Object.fromEntries(response.headers.entries())
  });
  
  return response;
};
```