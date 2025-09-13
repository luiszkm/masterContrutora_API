# API Documentation - Master Construtora

## Visão Geral

A API Master Construtora é uma REST API construída em Go que gerencia todos os aspectos de uma empresa de construção civil. A API utiliza autenticação JWT via cookies httpOnly e implementa autorização baseada em papéis (RBAC).

**Base URLs**:

- Desenvolvimento: `http://localhost:8081`
- Homologação (Staging): `https://staging.api.masterconstrutora.com`
- Produção: `https://api.masterconstrutora.com`

## Autenticação

### Registro de Usuário

**POST** `/usuarios/registrar`

Registra um novo usuário no sistema.

```json
// Request
{
  "nome": "João Silva",
  "email": "joao@empresa.com",
  "senha": "senha_forte_123",
  "confirmarSenha": "senha_forte_123"
}

// Response (201 Created)
{
  "id": "uuid-do-usuario",
  "nome": "João Silva",
  "email": "joao@empresa.com"
}
```

### Login

**POST** `/usuarios/login`

Autentica um usuário e retorna token JWT em cookie httpOnly.

```json
// Request
{
  "email": "joao@empresa.com",
  "senha": "senha_forte_123"
}

// Response (200 OK)
{
  "accessToken": "jwt-token-string",
  "userId": "uuid-do-usuario"
}
```

**Headers de resposta:**
- `Set-Cookie: jwt-token=<token>; HttpOnly; Secure; SameSite=Lax`

## Módulo Obras

### Listar Obras

**GET** `/obras`

Lista todas as obras com paginação e filtros.

**Permissão**: `obras:ler`

```json
// Query Parameters
?page=1&limit=10&nome=casa&status=Em+Andamento

// Response (200 OK)
{
  "dados": [
    {
      "id": "uuid-obra",
      "nome": "Casa Silva",
      "cliente": "João Silva",
      "endereco": "Rua das Flores, 123",
      "status": "Em Andamento",
      "dataInicio": "2024-01-15T00:00:00Z",
      "dataFim": null
    }
  ],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 5,
    "totalItens": 50,
    "itensPorPagina": 10
  }
}
```

### Criar Obra

**POST** `/obras`

Cria uma nova obra.

**Permissão**: `obras:escrever`

```json
// Request
{
  "nome": "Prédio Comercial ABC",
  "cliente": "Empresa XYZ Ltda",
  "endereco": "Av. Principal, 456",
  "descricao": "Prédio comercial de 10 andares",
  "dataInicio": "2024-02-01",
  "etapasPadrao": [
    "fundacao-id",
    "estrutura-id",
    "acabamento-id"
  ]
}

// Response (201 Created)
{
  "id": "uuid-nova-obra",
  "nome": "Prédio Comercial ABC",
  "cliente": "Empresa XYZ Ltda",
  "endereco": "Av. Principal, 456",
  "descricao": "Prédio comercial de 10 andares",
  "dataInicio": "2024-02-01T00:00:00Z",
  "dataFim": null,
  "status": "Em Planejamento"
}
```

### Buscar Dashboard da Obra

**GET** `/obras/{obraId}/dashboard`

Retorna dashboard completo com métricas calculadas.

**Permissão**: `obras:ler`

```json
// Response (200 OK)
{
  "obraId": "uuid-obra",
  "nomeObra": "Casa Silva",
  "statusObra": "Em Andamento",
  "etapaAtualNome": "Fundação",
  "dataFimPrevistaEtapa": "2024-03-15T00:00:00Z",
  "diasParaPrazoEtapa": 10,
  "percentualConcluido": 35.5,
  "custoTotalRealizado": 85000.00,
  "orcamentoTotalAprovado": 150000.00,
  "balancoFinanceiro": -65000.00,
  "funcionariosAlocados": 8,
  "ultimaAtualizacao": "2024-02-20T14:30:00Z"
}
```

### Buscar Detalhes da Obra

**GET** `/obras/{obraId}`

Retorna detalhes completos da obra.

**Permissão**: `obras:ler`

```json
// Response (200 OK)
{
  "id": "uuid-obra",
  "nome": "Casa Silva",
  "cliente": "João Silva",
  "endereco": "Rua das Flores, 123",
  "descricao": "Casa de 3 quartos",
  "dataInicio": "2024-01-15T00:00:00Z",
  "dataFim": null,
  "status": "Em Andamento"
}
```

### Atualizar Obra

**PUT** `/obras/{obraId}`

Atualiza dados de uma obra existente.

**Permissão**: `obras:escrever`

```json
// Request
{
  "nome": "Casa Silva Ampliada",
  "cliente": "João Silva",
  "endereco": "Rua das Flores, 123",
  "descricao": "Casa de 4 quartos com piscina",
  "dataFim": "2024-06-30",
  "status": "Em Andamento"
}

// Response (204 No Content)
```

### Deletar Obra

**DELETE** `/obras/{obraId}`

Remove uma obra (soft delete).

**Permissão**: `obras:escrever`

```json
// Response (204 No Content)
```

### Adicionar Etapa

**POST** `/obras/{obraId}/etapas`

Adiciona uma etapa à obra.

**Permissão**: `obras:escrever`

```json
// Request
{
  "nome": "Instalações Elétricas",
  "dataInicioPrevista": "2024-03-01",
  "dataFimPrevista": "2024-03-15"
}

// Response (201 Created)
{
  "id": "uuid-etapa",
  "obraId": "uuid-obra",
  "nome": "Instalações Elétricas",
  "dataInicioPrevista": "2024-03-01T00:00:00Z",
  "dataFimPrevista": "2024-03-15T00:00:00Z",
  "status": "Não Iniciada"
}
```

### Listar Etapas da Obra

**GET** `/obras/{obraId}/etapas`

Lista todas as etapas de uma obra específica.

**Permissão**: `obras:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-etapa-1",
    "obraId": "uuid-obra",
    "nome": "Fundação",
    "dataInicioPrevista": "2024-01-15T00:00:00Z",
    "dataFimPrevista": "2024-02-15T00:00:00Z",
    "status": "Concluída"
  },
  {
    "id": "uuid-etapa-2",
    "obraId": "uuid-obra",
    "nome": "Estrutura",
    "dataInicioPrevista": "2024-02-16T00:00:00Z",
    "dataFimPrevista": "2024-04-30T00:00:00Z",
    "status": "Em Andamento"
  }
]
```

### Atualizar Status da Etapa

**PATCH** `/etapas/{etapaId}`

Atualiza o status de uma etapa.

**Permissão**: `obras:escrever`

```json
// Request
{
  "status": "Concluída"
}

// Response (200 OK)
{
  "id": "uuid-etapa",
  "obraId": "uuid-obra",
  "nome": "Fundação",
  "dataInicioPrevista": "2024-01-15T00:00:00Z",
  "dataFimPrevista": "2024-02-15T00:00:00Z",
  "status": "Concluída"
}
```

### Alocar Funcionários

**POST** `/obras/{obraId}/alocacoes`

Aloca funcionários a uma obra.

**Permissão**: `obras:escrever`

```json
// Request
{
  "funcionarioIds": [
    "uuid-funcionario-1",
    "uuid-funcionario-2"
  ],
  "dataInicioAlocacao": "2024-02-01"
}

// Response (201 Created)
[
  {
    "id": "uuid-alocacao-1",
    "obraId": "uuid-obra",
    "funcionarioId": "uuid-funcionario-1",
    "dataInicioAlocacao": "2024-02-01T00:00:00Z",
    "dataFimAlocacao": null
  }
]
```

### Listar Etapas Padrão

**GET** `/etapas-padroes`

Lista etapas padrão disponíveis para criação de obras.

**Permissão**: `obras:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-etapa-padrao-1",
    "nome": "Fundação",
    "descricao": "Escavação e fundação da obra",
    "ordem": 1
  },
  {
    "id": "uuid-etapa-padrao-2",
    "nome": "Estrutura",
    "descricao": "Estrutura de concreto armado",
    "ordem": 2
  }
]
```

## Módulo Pessoal

### Listar Funcionários

**GET** `/funcionarios`

Lista todos os funcionários ativos.

**Permissão**: `pessoal:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-funcionario",
    "nome": "Carlos Santos",
    "cpf": "123.456.789-00",
    "cargo": "Pedreiro",
    "departamento": "Construção",
    "telefone": "(11) 99999-9999",
    "email": "carlos@empresa.com",
    "dataContratacao": "2023-06-15T00:00:00Z",
    "valorDiaria": 180.00,
    "chavePix": "carlos@email.com",
    "status": "Ativo"
  }
]
```

### Cadastrar Funcionário

**POST** `/funcionarios`

Cadastra um novo funcionário.

**Permissão**: `pessoal:escrever`

```json
// Request
{
  "nome": "Maria Oliveira",
  "cpf": "987.654.321-00",
  "cargo": "Eletricista",
  "departamento": "Instalações",
  "telefone": "(11) 88888-8888",
  "email": "maria@empresa.com",
  "diaria": 200.00,
  "chavePix": "maria@email.com"
}

// Response (201 Created)
{
  "id": "uuid-novo-funcionario",
  "nome": "Maria Oliveira",
  "cpf": "987.654.321-00",
  "cargo": "Eletricista",
  "departamento": "Instalações",
  "telefone": "(11) 88888-8888",
  "email": "maria@empresa.com",
  "dataContratacao": "2024-02-20T00:00:00Z",
  "valorDiaria": 200.00,
  "chavePix": "maria@email.com",
  "status": "Ativo"
}
```

### Buscar Funcionário

**GET** `/funcionarios/{funcionarioId}`

Busca detalhes de um funcionário específico.

**Permissão**: `pessoal:ler`

```json
// Response (200 OK)
{
  "id": "uuid-funcionario",
  "nome": "Carlos Santos",
  "cpf": "123.456.789-00",
  "cargo": "Pedreiro",
  "departamento": "Construção",
  "telefone": "(11) 99999-9999",
  "email": "carlos@empresa.com",
  "dataContratacao": "2023-06-15T00:00:00Z",
  "valorDiaria": 180.00,
  "chavePix": "carlos@email.com",
  "status": "Ativo"
}
```

### Atualizar Funcionário

**PUT** `/funcionarios/{funcionarioId}`

Atualiza dados de um funcionário.

**Permissão**: `pessoal:escrever`

```json
// Request
{
  "nome": "Carlos Santos Silva",
  "cargo": "Mestre de Obras",
  "departamento": "Construção",
  "telefone": "(11) 99999-9999",
  "email": "carlos.santos@empresa.com",
  "diaria": 220.00,
  "chavePix": "carlos.santos@email.com"
}

// Response (200 OK)
{
  "id": "uuid-funcionario",
  "nome": "Carlos Santos Silva",
  "cpf": "123.456.789-00",
  "cargo": "Mestre de Obras",
  "departamento": "Construção",
  "telefone": "(11) 99999-9999",
  "email": "carlos.santos@empresa.com",
  "dataContratacao": "2023-06-15T00:00:00Z",
  "valorDiaria": 220.00,
  "chavePix": "carlos.santos@email.com",
  "status": "Ativo"
}
```

### Deletar Funcionário

**DELETE** `/funcionarios/{funcionarioId}`

Remove um funcionário (soft delete).

**Permissão**: `pessoal:escrever`

```json
// Response (204 No Content)
```

### Ativar Funcionário

**PATCH** `/funcionarios/{funcionarioId}/ativar`

Reativa um funcionário desativado.

**Permissão**: `pessoal:apontamento:ler`

```json
// Response (204 No Content)
```

### Criar Apontamento

**POST** `/apontamentos`

Cria um apontamento quinzenal para um funcionário.

**Permissão**: `pessoal:apontamento:escrever`

```json
// Request
{
  "funcionarioId": "uuid-funcionario",
  "obraId": "uuid-obra",
  "periodoInicio": "2024-02-01",
  "periodoFim": "2024-02-15",
  "Diaria": 180.00,
  "DiasTrabalhados": 10,
  "Descontos": 50.00,
  "Adiantamento": 500.00,
  "ValorAdicional": 200.00
}

// Response (201 Created)
{
  "id": "uuid-apontamento",
  "funcionarioId": "uuid-funcionario",
  "obraId": "uuid-obra",
  "periodoInicio": "2024-02-01T00:00:00Z",
  "periodoFim": "2024-02-15T00:00:00Z",
  "diaria": 180.00,
  "diasTrabalhados": 10,
  "adicionais": 200.00,
  "descontos": 50.00,
  "adiantamentos": 500.00,
  "valorTotalCalculado": 1450.00,
  "status": "Em Aberto"
}
```

### Listar Apontamentos

**GET** `/apontamentos`

Lista todos os apontamentos com paginação.

**Permissão**: `pessoal:apontamento:ler`

```json
// Query Parameters
?page=1&limit=10&status=Em+Aberto&funcionarioId=uuid

// Response (200 OK)
{
  "dados": [
    {
      "id": "uuid-apontamento",
      "funcionarioId": "uuid-funcionario",
      "obraId": "uuid-obra",
      "periodoInicio": "2024-02-01T00:00:00Z",
      "periodoFim": "2024-02-15T00:00:00Z",
      "diaria": 180.00,
      "diasTrabalhados": 10,
      "adicionais": 200.00,
      "descontos": 50.00,
      "adiantamentos": 500.00,
      "valorTotalCalculado": 1450.00,
      "status": "Em Aberto",
      "nomeFuncionario": "Carlos Santos"
    }
  ],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 3,
    "totalItens": 25,
    "itensPorPagina": 10
  }
}
```

### Aprovar Apontamento

**PATCH** `/apontamentos/{apontamentoId}/aprovar`

Aprova um apontamento para pagamento.

**Permissão**: `pessoal:apontamento:aprovar`

```json
// Response (200 OK)
{
  "id": "uuid-apontamento",
  "funcionarioId": "uuid-funcionario",
  "obraId": "uuid-obra",
  "periodoInicio": "2024-02-01T00:00:00Z",
  "periodoFim": "2024-02-15T00:00:00Z",
  "diaria": 180.00,
  "diasTrabalhados": 10,
  "adicionais": 200.00,
  "descontos": 50.00,
  "adiantamentos": 500.00,
  "valorTotalCalculado": 1450.00,
  "status": "Aprovado para Pagamento"
}
```

### Cancelar Apontamento

**PATCH** `/apontamentos/{apontamentoId}/cancelar`

Cancela um apontamento existente.

**Permissão**: `pessoal:apontamento:escrever`

```json
// Response (200 OK)
{
  "id": "uuid-apontamento",
  "funcionarioId": "uuid-funcionario",
  "obraId": "uuid-obra",
  "periodoInicio": "2024-02-01T00:00:00Z",
  "periodoFim": "2024-02-15T00:00:00Z",
  "diaria": 180.00,
  "diasTrabalhados": 10,
  "adicionais": 200.00,
  "descontos": 50.00,
  "adiantamentos": 500.00,
  "valorTotalCalculado": 1450.00,
  "status": "Cancelado"
}
```

### Atualizar Apontamento

**PUT** `/funcionarios/apontamentos/{apontamentoId}`

Atualiza dados de um apontamento em aberto.

**Permissão**: `pessoal:apontamento:escrever`

```json
// Request
{
  "funcionarioId": "uuid-funcionario",
  "obraId": "uuid-obra",
  "periodoInicio": "2024-02-01",
  "periodoFim": "2024-02-15",
  "diaria": 185.00,
  "diasTrabalhados": 12,
  "valorAdicional": 250.00,
  "descontos": 30.00,
  "adiantamento": 600.00,
  "status": "Em Aberto"
}

// Response (200 OK)
{
  "id": "uuid-apontamento",
  "funcionarioId": "uuid-funcionario",
  "obraId": "uuid-obra",
  "periodoInicio": "2024-02-01T00:00:00Z",
  "periodoFim": "2024-02-15T00:00:00Z",
  "diaria": 185.00,
  "diasTrabalhados": 12,
  "adicionais": 250.00,
  "descontos": 30.00,
  "adiantamentos": 600.00,
  "valorTotalCalculado": 1840.00,
  "status": "Em Aberto"
}
```

### Listar Funcionários com Apontamentos

**GET** `/funcionarios/apontamentos`

Lista funcionários com dados do último apontamento.

**Permissão**: `pessoal:ler`

```json
// Query Parameters
?status=Ativo&page=1&limit=10

// Response (200 OK)
[
  {
    "id": "uuid-funcionario",
    "nome": "Carlos Santos",
    "cargo": "Pedreiro",
    "departamento": "Construção",
    "dataContratacao": "2023-06-15T00:00:00Z",
    "valorDiaria": 180.00,
    "diasTrabalhados": 10,
    "valorAdicional": 200.00,
    "descontos": 50.00,
    "adiantamento": 500.00,
    "chavePix": "carlos@email.com",
    "statusApontamento": "Em Aberto",
    "apontamentoId": "uuid-apontamento"
  }
]
```

### Listar Apontamentos por Funcionário

**GET** `/funcionarios/{funcionarioId}/apontamentos`

Lista apontamentos de um funcionário específico.

**Permissão**: `pessoal:apontamento:ler`

```json
// Response (200 OK)
{
  "dados": [
    {
      "id": "uuid-apontamento",
      "funcionarioId": "uuid-funcionario",
      "obraId": "uuid-obra",
      "periodoInicio": "2024-02-01T00:00:00Z",
      "periodoFim": "2024-02-15T00:00:00Z",
      "diaria": 180.00,
      "diasTrabalhados": 10,
      "valorTotalCalculado": 1450.00,
      "status": "Pago"
    }
  ],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 2,
    "totalItens": 15,
    "itensPorPagina": 10
  }
}
```

### Replicar Apontamentos

**POST** `/funcionarios/apontamentos/replicar`

Replica apontamentos para a próxima quinzena.

**Permissão**: `pessoal:apontamento:escrever`

```json
// Request
{
  "funcionarioIds": [
    "uuid-funcionario-1",
    "uuid-funcionario-2"
  ]
}

// Response (207 Multi-Status)
{
  "sucessos": [
    {
      "funcionarioId": "uuid-funcionario-1",
      "apontamentoId": "uuid-novo-apontamento",
      "mensagem": "Apontamento replicado com sucesso"
    }
  ],
  "erros": [
    {
      "funcionarioId": "uuid-funcionario-2",
      "erro": "Funcionário já possui apontamento para o próximo período"
    }
  ]
}
```

## Módulo Suprimentos

### Listar Fornecedores

**GET** `/fornecedores`

Lista todos os fornecedores ativos.

**Permissão**: `suprimentos:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-fornecedor",
    "nome": "Fornecedor ABC Ltda",
    "cnpj": "12.345.678/0001-90",
    "contato": "João Fornecedor",
    "email": "contato@fornecedorabc.com",
    "endereco": "Rua dos Fornecedores, 123",
    "status": "Ativo",
    "avaliacao": 4.5,
    "observacoes": "Fornecedor confiável"
  }
]
```

### Cadastrar Fornecedor

**POST** `/fornecedores`

Cadastra um novo fornecedor.

**Permissão**: `suprimentos:escrever`

```json
// Request
{
  "nome": "Materiais XYZ Ltda",
  "cnpj": "98.765.432/0001-10",
  "contato": "Maria Silva",
  "email": "vendas@materiaisxyz.com",
  "endereco": "Av. dos Materiais, 456",
  "observacoes": "Especializada em materiais de acabamento"
}

// Response (201 Created)
{
  "id": "uuid-novo-fornecedor",
  "nome": "Materiais XYZ Ltda",
  "cnpj": "98.765.432/0001-10",
  "contato": "Maria Silva",
  "email": "vendas@materiaisxyz.com",
  "endereco": "Av. dos Materiais, 456",
  "status": "Ativo",
  "avaliacao": null,
  "observacoes": "Especializada em materiais de acabamento"
}
```

### Buscar Fornecedor

**GET** `/fornecedores/{id}`

Busca detalhes de um fornecedor específico.

**Permissão**: `suprimentos:ler`

```json
// Response (200 OK)
{
  "id": "uuid-fornecedor",
  "nome": "Fornecedor ABC Ltda",
  "cnpj": "12.345.678/0001-90",
  "contato": "João Fornecedor",
  "email": "contato@fornecedorabc.com",
  "endereco": "Rua dos Fornecedores, 123",
  "status": "Ativo",
  "avaliacao": 4.5,
  "observacoes": "Fornecedor confiável"
}
```

### Atualizar Fornecedor

**PUT** `/fornecedores/{id}`

Atualiza dados de um fornecedor.

**Permissão**: `suprimentos:escrever`

```json
// Request
{
  "nome": "Fornecedor ABC Materiais Ltda",
  "contato": "João Fornecedor Junior",
  "email": "novo@fornecedorabc.com",
  "endereco": "Rua dos Fornecedores, 456",
  "avaliacao": 4.8,
  "observacoes": "Fornecedor muito confiável, entregas pontuais"
}

// Response (200 OK)
{
  "id": "uuid-fornecedor",
  "nome": "Fornecedor ABC Materiais Ltda",
  "cnpj": "12.345.678/0001-90",
  "contato": "João Fornecedor Junior",
  "email": "novo@fornecedorabc.com",
  "endereco": "Rua dos Fornecedores, 456",
  "status": "Ativo",
  "avaliacao": 4.8,
  "observacoes": "Fornecedor muito confiável, entregas pontuais"
}
```

### Deletar Fornecedor

**DELETE** `/fornecedores/{id}`

Remove um fornecedor (soft delete).

**Permissão**: `suprimentos:escrever`

```json
// Response (204 No Content)
```

### Listar Materiais/Produtos

**GET** `/materiais`

Lista todos os produtos/materiais cadastrados.

**Permissão**: `suprimentos:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-produto",
    "nome": "Cimento CP II",
    "descricao": "Cimento Portland composto",
    "unidadeDeMedida": "saco 50kg",
    "categoria": "Materiais Básicos"
  }
]
```

### Cadastrar Material/Produto

**POST** `/materiais`

Cadastra um novo produto/material.

**Permissão**: `suprimentos:escrever`

```json
// Request
{
  "nome": "Tijolo Cerâmico",
  "descricao": "Tijolo cerâmico 6 furos",
  "unidadeDeMedida": "milheiro",
  "categoria": "Alvenaria"
}

// Response (201 Created)
{
  "id": "uuid-novo-produto",
  "nome": "Tijolo Cerâmico",
  "descricao": "Tijolo cerâmico 6 furos",
  "unidadeDeMedida": "milheiro",
  "categoria": "Alvenaria"
}
```

### Listar Categorias

**GET** `/categorias`

Lista todas as categorias de produtos com suporte opcional a paginação.

**Permissão**: `suprimentos:ler`

```json
// Query Parameters (opcional para paginação)
?page=1&limit=10

// Response (200 OK) - Sem paginação
[
  {
    "id": "uuid-categoria",
    "nome": "Materiais Básicos"
  },
  {
    "id": "uuid-categoria-2",
    "nome": "Acabamento"
  }
]

// Response (200 OK) - Com paginação
{
  "dados": [
    {
      "id": "uuid-categoria",
      "nome": "Materiais Básicos"
    },
    {
      "id": "uuid-categoria-2",
      "nome": "Acabamento"
    }
  ],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 2,
    "totalItens": 15,
    "itensPorPagina": 10
  }
}
```

### Criar Categoria

**POST** `/categorias`

Cria uma nova categoria de produtos.

**Permissão**: `suprimentos:escrever`

```json
// Request
{
  "nome": "Ferramentas"
}

// Response (201 Created)
{
  "id": "uuid-nova-categoria",
  "nome": "Ferramentas"
}
```

### Buscar Categoria

**GET** `/categorias/{categoriaId}`

Busca detalhes de uma categoria específica.

**Permissão**: `suprimentos:ler`

```json
// Response (200 OK)
{
  "id": "uuid-categoria",
  "nome": "Materiais Básicos"
}
```

### Atualizar Categoria

**PUT** `/categorias/{categoriaId}`

Atualiza dados de uma categoria existente.

**Permissão**: `suprimentos:escrever`

```json
// Request
{
  "nome": "Materiais Básicos e Estruturais"
}

// Response (200 OK)
{
  "id": "uuid-categoria",
  "nome": "Materiais Básicos e Estruturais"
}
```

### Deletar Categoria

**DELETE** `/categorias/{categoriaId}`

Remove uma categoria do sistema.

**Permissão**: `suprimentos:escrever`

```json
// Response (204 No Content)

// Response (409 Conflict) - Se categoria estiver em uso
{
  "erro": {
    "codigo": "CONFLITO",
    "mensagem": "A categoria está em uso e não pode ser deletada"
  }
}
```

### Criar Orçamento

**POST** `/etapas/{etapaId}/orcamentos`

Cria um orçamento para uma etapa específica.

**Permissão**: `suprimentos:escrever`

```json
// Request
{
  "fornecedorId": "uuid-fornecedor",
  "itens": [
    {
      "nomeProduto": "Cimento CP II",
      "unidadeDeMedida": "saco 50kg",
      "categoria": "Materiais Básicos",
      "quantidade": 20,
      "valorUnitario": 35.50
    },
    {
      "nomeProduto": "Areia Média",
      "unidadeDeMedida": "m³",
      "categoria": "Agregados",
      "quantidade": 5,
      "valorUnitario": 45.00
    }
  ]
}

// Response (201 Created)
{
  "id": "uuid-orcamento",
  "numero": "ORC-2024-001",
  "etapaId": "uuid-etapa",
  "fornecedorId": "uuid-fornecedor",
  "valorTotal": 935.00,
  "status": "Em Aberto",
  "dataEmissao": "2024-02-20T14:30:00Z",
  "itens": [
    {
      "nomeProduto": "Cimento CP II",
      "unidadeDeMedida": "saco 50kg",
      "categoria": "Materiais Básicos",
      "quantidade": 20,
      "valorUnitario": 35.50
    }
  ]
}
```

### Listar Orçamentos

**GET** `/orcamentos`

Lista todos os orçamentos com filtros e paginação.

**Permissão**: `suprimentos:ler`

```json
// Query Parameters
?page=1&limit=10&status=Em+Aberto&fornecedorId=uuid&obraId=uuid

// Response (200 OK)
{
  "dados": [
    {
      "id": "uuid-orcamento",
      "numero": "ORC-2024-001",
      "valorTotal": 935.00,
      "status": "Em Aberto",
      "dataEmissao": "2024-02-20T14:30:00Z",
      "obraId": "uuid-obra",
      "obraNome": "Casa Silva",
      "fornecedorId": "uuid-fornecedor",
      "fornecedorNome": "Fornecedor ABC Ltda",
      "itensCount": 2,
      "categorias": ["Materiais Básicos", "Acabamento"]
    }
  ],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 5,
    "totalItens": 50,
    "itensPorPagina": 10
  }
}
```

### Buscar Orçamento Detalhado

**GET** `/orcamentos/{orcamentoId}`

Busca detalhes completos de um orçamento.

**Permissão**: `suprimentos:ler`

```json
// Response (200 OK)
{
  "id": "uuid-orcamento",
  "numero": "ORC-2024-001",
  "valorTotal": 935.00,
  "status": "Em Aberto",
  "dataEmissao": "2024-02-20T14:30:00Z",
  "observacoes": "Orçamento para primeira etapa",
  "condicoesPagamento": "30 dias",
  "obra": {
    "id": "uuid-obra",
    "nome": "Casa Silva"
  },
  "etapa": {
    "id": "uuid-etapa",
    "nome": "Fundação"
  },
  "fornecedor": {
    "id": "uuid-fornecedor",
    "nome": "Fornecedor ABC Ltda"
  },
  "itens": [
    {
      "ProdutoNome": "Cimento CP II",
      "UnidadeDeMedida": "saco 50kg",
      "Categoria": "Materiais Básicos",
      "Quantidade": 20,
      "ValorUnitario": 35.50
    }
  ]
}
```

### Atualizar Orçamento

**PUT** `/orcamentos/{orcamentoId}`

Atualiza dados de um orçamento existente.

**Permissão**: `suprimentos:escrever`

```json
// Request
{
  "fornecedorId": "uuid-fornecedor",
  "etapaId": "uuid-etapa",
  "observacoes": "Orçamento revisado com desconto",
  "condicoesPagamento": "45 dias",
  "itens": [
    {
      "nomeProduto": "Cimento CP II",
      "unidadeDeMedida": "saco 50kg",
      "categoria": "Materiais Básicos",
      "quantidade": 25,
      "valorUnitario": 33.00
    }
  ]
}

// Response (200 OK)
{
  "id": "uuid-orcamento",
  "numero": "ORC-2024-001",
  "valorTotal": 825.00,
  "status": "Em Aberto",
  "dataEmissao": "2024-02-20T14:30:00Z",
  "observacoes": "Orçamento revisado com desconto",
  "condicoesPagamento": "45 dias"
}
```

### Atualizar Status do Orçamento

**PATCH** `/orcamentos/{orcamentoId}/status`

Atualiza apenas o status de um orçamento.

**Permissão**: `suprimentos:escrever`

```json
// Request
{
  "status": "Aprovado"
}

// Response (200 OK)
{
  "id": "uuid-orcamento",
  "numero": "ORC-2024-001",
  "valorTotal": 825.00,
  "status": "Aprovado",
  "dataEmissao": "2024-02-20T14:30:00Z",
  "dataAprovacao": "2024-02-21T10:15:00Z"
}
```

## Módulo Financeiro

### Contas a Receber

#### Criar Conta a Receber

**POST** `/contas-receber`

Cria uma nova conta a receber.

**Permissão**: `financeiro:escrever`

```json
// Request
{
  "obraId": "uuid-obra",
  "descricao": "Pagamento da 1ª etapa - Fundação",
  "valorTotal": 15000.00,
  "dataVencimento": "2024-03-15",
  "parcela": 1,
  "totalParcelas": 4
}

// Response (201 Created)
{
  "id": "uuid-conta-receber",
  "obraId": "uuid-obra",
  "descricao": "Pagamento da 1ª etapa - Fundação",
  "valorTotal": 15000.00,
  "valorRecebido": 0.00,
  "valorPendente": 15000.00,
  "dataVencimento": "2024-03-15T00:00:00Z",
  "status": "Pendente",
  "parcela": 1,
  "totalParcelas": 4,
  "dataCriacao": "2024-02-20T14:30:00Z"
}
```

#### Listar Contas a Receber

**GET** `/contas-receber`

Lista todas as contas a receber com paginação.

**Permissão**: `financeiro:ler`

```json
// Query Parameters
?page=1&limit=10&status=Pendente&obraId=uuid

// Response (200 OK)
{
  "dados": [
    {
      "id": "uuid-conta-receber",
      "obraId": "uuid-obra",
      "obraNome": "Casa Silva",
      "descricao": "Pagamento da 1ª etapa",
      "valorTotal": 15000.00,
      "valorRecebido": 0.00,
      "valorPendente": 15000.00,
      "dataVencimento": "2024-03-15T00:00:00Z",
      "status": "Pendente",
      "parcela": 1,
      "totalParcelas": 4
    }
  ],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 3,
    "totalItens": 25,
    "itensPorPagina": 10
  }
}
```

#### Buscar Conta a Receber

**GET** `/contas-receber/{contaId}`

Busca detalhes de uma conta a receber específica.

**Permissão**: `financeiro:ler`

```json
// Response (200 OK)
{
  "id": "uuid-conta-receber",
  "obraId": "uuid-obra",
  "obraNome": "Casa Silva",
  "descricao": "Pagamento da 1ª etapa - Fundação",
  "valorTotal": 15000.00,
  "valorRecebido": 7500.00,
  "valorPendente": 7500.00,
  "dataVencimento": "2024-03-15T00:00:00Z",
  "status": "Parcialmente Recebido",
  "parcela": 1,
  "totalParcelas": 4,
  "recebimentos": [
    {
      "valor": 7500.00,
      "dataRecebimento": "2024-03-10T10:00:00Z",
      "observacoes": "Recebimento parcial via PIX"
    }
  ]
}
```

#### Registrar Recebimento

**POST** `/contas-receber/{contaId}/recebimentos`

Registra um recebimento para uma conta a receber.

**Permissão**: `financeiro:escrever`

```json
// Request
{
  "valor": 7500.00,
  "dataRecebimento": "2024-03-10",
  "observacoes": "Recebimento via PIX"
}

// Response (201 Created)
{
  "id": "uuid-conta-receber",
  "valorTotal": 15000.00,
  "valorRecebido": 7500.00,
  "valorPendente": 7500.00,
  "status": "Parcialmente Recebido"
}
```

#### Listar Contas Vencidas

**GET** `/contas-receber/vencidas`

Lista contas a receber vencidas.

**Permissão**: `financeiro:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-conta-receber",
    "obraNome": "Casa Silva",
    "descricao": "Pagamento da 1ª etapa",
    "valorPendente": 15000.00,
    "dataVencimento": "2024-02-15T00:00:00Z",
    "diasVencido": 5
  }
]
```

#### Obter Resumo

**GET** `/contas-receber/resumo`

Obter resumo das contas a receber.

**Permissão**: `financeiro:ler`

```json
// Response (200 OK)
{
  "totalPendente": 45000.00,
  "totalRecebido": 25000.00,
  "totalVencido": 10000.00,
  "quantidadePendente": 8,
  "quantidadeRecebida": 5,
  "quantidadeVencida": 2
}
```

### Contas a Pagar

#### Criar Conta a Pagar

**POST** `/contas-pagar`

Cria uma nova conta a pagar.

**Permissão**: `financeiro:escrever`

```json
// Request
{
  "fornecedorId": "uuid-fornecedor",
  "obraId": "uuid-obra",
  "descricao": "Materiais para fundação",
  "valorTotal": 8500.00,
  "dataVencimento": "2024-03-20",
  "tipo": "FORNECEDOR"
}

// Response (201 Created)
{
  "id": "uuid-conta-pagar",
  "fornecedorId": "uuid-fornecedor",
  "obraId": "uuid-obra",
  "descricao": "Materiais para fundação",
  "valorTotal": 8500.00,
  "valorPago": 0.00,
  "valorPendente": 8500.00,
  "dataVencimento": "2024-03-20T00:00:00Z",
  "status": "Pendente",
  "tipo": "FORNECEDOR"
}
```

#### Listar Contas a Pagar

**GET** `/contas-pagar`

Lista todas as contas a pagar com paginação.

**Permissão**: `financeiro:ler`

```json
// Query Parameters
?page=1&limit=10&status=Pendente&fornecedorId=uuid

// Response (200 OK)
{
  "dados": [
    {
      "id": "uuid-conta-pagar",
      "fornecedorId": "uuid-fornecedor",
      "fornecedorNome": "Materiais ABC Ltda",
      "obraId": "uuid-obra",
      "obraNome": "Casa Silva",
      "descricao": "Materiais para fundação",
      "valorTotal": 8500.00,
      "valorPago": 0.00,
      "valorPendente": 8500.00,
      "dataVencimento": "2024-03-20T00:00:00Z",
      "status": "Pendente",
      "tipo": "FORNECEDOR"
    }
  ],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 2,
    "totalItens": 15,
    "itensPorPagina": 10
  }
}
```

#### Buscar Conta a Pagar

**GET** `/contas-pagar/{contaId}`

Busca detalhes de uma conta a pagar específica.

**Permissão**: `financeiro:ler`

```json
// Response (200 OK)
{
  "id": "uuid-conta-pagar",
  "fornecedorId": "uuid-fornecedor",
  "fornecedorNome": "Materiais ABC Ltda",
  "obraId": "uuid-obra",
  "obraNome": "Casa Silva",
  "descricao": "Materiais para fundação",
  "valorTotal": 8500.00,
  "valorPago": 4000.00,
  "valorPendente": 4500.00,
  "dataVencimento": "2024-03-20T00:00:00Z",
  "status": "Parcialmente Pago",
  "tipo": "FORNECEDOR",
  "pagamentos": [
    {
      "valor": 4000.00,
      "dataPagamento": "2024-03-15T14:00:00Z",
      "observacoes": "Pagamento parcial"
    }
  ]
}
```

#### Registrar Pagamento

**POST** `/contas-pagar/{contaId}/pagamentos`

Registra um pagamento para uma conta a pagar.

**Permissão**: `financeiro:escrever`

```json
// Request
{
  "valor": 4000.00,
  "dataPagamento": "2024-03-15",
  "observacoes": "Pagamento parcial"
}

// Response (201 Created)
{
  "id": "uuid-conta-pagar",
  "valorTotal": 8500.00,
  "valorPago": 4000.00,
  "valorPendente": 4500.00,
  "status": "Parcialmente Pago"
}
```

#### Criar Conta de Orçamento

**POST** `/contas-pagar/orcamentos`

Cria conta a pagar automaticamente a partir de orçamento aprovado.

**Permissão**: `financeiro:escrever`

```json
// Request
{
  "orcamentoId": "uuid-orcamento"
}

// Response (201 Created)
{
  "id": "uuid-conta-pagar",
  "orcamentoId": "uuid-orcamento",
  "fornecedorId": "uuid-fornecedor",
  "obraId": "uuid-obra",
  "descricao": "Conta gerada automaticamente do orçamento ORC-2024-001",
  "valorTotal": 8500.00,
  "status": "Pendente",
  "tipo": "FORNECEDOR"
}
```

#### Listar Contas Vencidas

**GET** `/contas-pagar/vencidas`

Lista contas a pagar vencidas.

**Permissão**: `financeiro:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-conta-pagar",
    "fornecedorNome": "Materiais ABC Ltda",
    "descricao": "Materiais para fundação",
    "valorPendente": 8500.00,
    "dataVencimento": "2024-02-15T00:00:00Z",
    "diasVencido": 10
  }
]
```

#### Obter Resumo

**GET** `/contas-pagar/resumo`

Obter resumo das contas a pagar.

**Permissão**: `financeiro:ler`

```json
// Response (200 OK)
{
  "totalPendente": 35000.00,
  "totalPago": 18000.00,
  "totalVencido": 12000.00,
  "quantidadePendente": 6,
  "quantidadePaga": 8,
  "quantidadeVencida": 3
}
```

### Consultas por Relacionamentos

#### Listar Contas a Receber por Obra

**GET** `/obras/{obraId}/contas-receber`

Lista contas a receber de uma obra específica.

**Permissão**: `financeiro:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-conta-receber",
    "descricao": "Pagamento da 1ª etapa",
    "valorTotal": 15000.00,
    "valorPendente": 7500.00,
    "status": "Parcialmente Recebido",
    "dataVencimento": "2024-03-15T00:00:00Z"
  }
]
```

#### Listar Contas a Pagar por Obra

**GET** `/obras/{obraId}/contas-pagar`

Lista contas a pagar de uma obra específica.

**Permissão**: `financeiro:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-conta-pagar",
    "fornecedorNome": "Materiais ABC Ltda",
    "descricao": "Materiais para fundação",
    "valorTotal": 8500.00,
    "valorPendente": 4500.00,
    "status": "Parcialmente Pago",
    "dataVencimento": "2024-03-20T00:00:00Z"
  }
]
```

#### Listar Contas a Pagar por Fornecedor

**GET** `/fornecedores/{fornecedorId}/contas-pagar`

Lista contas a pagar de um fornecedor específico.

**Permissão**: `financeiro:ler`

```json
// Response (200 OK)
[
  {
    "id": "uuid-conta-pagar",
    "obraNome": "Casa Silva",
    "descricao": "Materiais para fundação",
    "valorTotal": 8500.00,
    "valorPendente": 4500.00,
    "status": "Parcialmente Pago",
    "dataVencimento": "2024-03-20T00:00:00Z"
  }
]
```

### Cronograma de Recebimentos

#### Criar Cronograma

**POST** `/cronograma-recebimentos`

Cria um cronograma de recebimento.

**Permissão**: `obras:escrever`

```json
// Request
{
  "obraId": "uuid-obra",
  "etapaId": "uuid-etapa",
  "descricao": "Pagamento da fundação",
  "valorPrevisto": 15000.00,
  "dataVencimentoPrevista": "2024-03-15",
  "parcela": 1,
  "totalParcelas": 4
}

// Response (201 Created)
{
  "id": "uuid-cronograma",
  "obraId": "uuid-obra",
  "etapaId": "uuid-etapa",
  "descricao": "Pagamento da fundação",
  "valorPrevisto": 15000.00,
  "valorRecebido": 0.00,
  "dataVencimentoPrevista": "2024-03-15T00:00:00Z",
  "status": "Pendente",
  "parcela": 1,
  "totalParcelas": 4
}
```

#### Criar Cronogramas em Lote

**POST** `/cronograma-recebimentos/lote`

Cria múltiplos cronogramas de uma vez.

**Permissão**: `obras:escrever`

```json
// Request
{
  "obraId": "uuid-obra",
  "cronogramas": [
    {
      "etapaId": "uuid-etapa-1",
      "descricao": "Pagamento da fundação",
      "valorPrevisto": 15000.00,
      "dataVencimentoPrevista": "2024-03-15",
      "parcela": 1
    },
    {
      "etapaId": "uuid-etapa-2",
      "descricao": "Pagamento da estrutura",
      "valorPrevisto": 25000.00,
      "dataVencimentoPrevista": "2024-04-15",
      "parcela": 2
    }
  ]
}

// Response (201 Created)
[
  {
    "id": "uuid-cronograma-1",
    "obraId": "uuid-obra",
    "descricao": "Pagamento da fundação",
    "valorPrevisto": 15000.00,
    "status": "Pendente"
  }
]
```

#### Buscar Cronograma

**GET** `/cronograma-recebimentos/{cronogramaId}`

Busca detalhes de um cronograma específico.

**Permissão**: `obras:ler`

```json
// Response (200 OK)
{
  "id": "uuid-cronograma",
  "obraId": "uuid-obra",
  "obraNome": "Casa Silva",
  "etapaId": "uuid-etapa",
  "etapaNome": "Fundação",
  "descricao": "Pagamento da fundação",
  "valorPrevisto": 15000.00,
  "valorRecebido": 7500.00,
  "dataVencimentoPrevista": "2024-03-15T00:00:00Z",
  "status": "Parcialmente Recebido",
  "recebimentos": [
    {
      "valor": 7500.00,
      "dataRecebimento": "2024-03-10T10:00:00Z",
      "observacoes": "Recebimento parcial"
    }
  ]
}
```

#### Registrar Recebimento

**POST** `/cronograma-recebimentos/{cronogramaId}/recebimentos`

Registra um recebimento para um cronograma.

**Permissão**: `obras:escrever`

```json
// Request
{
  "valor": 7500.00,
  "dataRecebimento": "2024-03-10",
  "observacoes": "Recebimento parcial via PIX"
}

// Response (201 Created)
{
  "id": "uuid-cronograma",
  "valorPrevisto": 15000.00,
  "valorRecebido": 7500.00,
  "status": "Parcialmente Recebido"
}
```

#### Listar Cronogramas por Obra

**GET** `/obras/{obraId}/cronograma-recebimentos`

Lista cronogramas de uma obra específica.

**Permissão**: `obras:ler`

```json
// Query Parameters
?page=1&limit=10

// Response (200 OK) - Com paginação
{
  "dados": [
    {
      "id": "uuid-cronograma",
      "etapaId": "uuid-etapa",
      "etapaNome": "Fundação",
      "descricao": "Pagamento da fundação",
      "valorPrevisto": 15000.00,
      "valorRecebido": 7500.00,
      "dataVencimentoPrevista": "2024-03-15T00:00:00Z",
      "status": "Parcialmente Recebido"
    }
  ],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 2,
    "totalItens": 12,
    "itensPorPagina": 10
  }
}

// Response (200 OK) - Sem paginação (backward compatibility)
[
  {
    "id": "uuid-cronograma",
    "etapaId": "uuid-etapa",
    "etapaNome": "Fundação",
    "valorPrevisto": 15000.00,
    "valorRecebido": 7500.00,
    "status": "Parcialmente Recebido"
  }
]
```

### Registros de Pagamento (Funcionários)

#### Registrar Pagamento

**POST** `/pagamentos`

Registra um pagamento individual.

**Permissão**: `financeiro:escrever`

```json
// Request
{
  "funcionarioId": "uuid-funcionario",
  "obraId": "uuid-obra",
  "periodoReferencia": "Fevereiro/2024",
  "valorCalculado": 1450.00,
  "contaBancariaId": "uuid-conta-bancaria"
}

// Response (201 Created)
{
  "id": "uuid-pagamento",
  "funcionarioId": "uuid-funcionario",
  "obraId": "uuid-obra",
  "periodoReferencia": "Fevereiro/2024",
  "valorCalculado": 1450.00,
  "dataDeEfetivacao": "2024-02-21T14:30:00Z",
  "contaBancariaId": "uuid-conta-bancaria"
}
```

#### Registrar Pagamentos em Lote

**POST** `/pagamentos/lote`

Registra múltiplos pagamentos de funcionários.

**Permissão**: `financeiro:escrever`

```json
// Request
{
  "apontamentoIds": [
    "uuid-apontamento-1",
    "uuid-apontamento-2"
  ],
  "contaBancariaId": "uuid-conta-bancaria",
  "dataDeEfetivacao": "2024-02-21"
}

// Response (207 Multi-Status)
{
  "sucessos": [
    {
      "apontamentoId": "uuid-apontamento-1",
      "funcionarioId": "uuid-funcionario-1",
      "valor": 1450.00,
      "pagamentoId": "uuid-pagamento-1"
    }
  ],
  "erros": [
    {
      "apontamentoId": "uuid-apontamento-2",
      "erro": "Apontamento não encontrado ou não aprovado"
    }
  ]
}
```

#### Listar Pagamentos

**GET** `/pagamentos`

Lista todos os pagamentos de funcionários.

**Permissão**: `financeiro:ler`

```json
// Query Parameters
?page=1&limit=10&funcionarioId=uuid&obraId=uuid

// Response (200 OK)
{
  "dados": [
    {
      "id": "uuid-pagamento",
      "funcionarioId": "uuid-funcionario",
      "funcionarioNome": "Carlos Santos",
      "obraId": "uuid-obra",
      "obraNome": "Casa Silva",
      "periodoReferencia": "Fevereiro/2024",
      "valorCalculado": 1450.00,
      "dataDeEfetivacao": "2024-02-21T14:30:00Z"
    }
  ],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 3,
    "totalItens": 25,
    "itensPorPagina": 10
  }
}
```

## Módulo Dashboard

### Dashboard Completo

**GET** `/dashboard/`

Retorna dashboard completo com todas as seções.

**Autenticação**: Não requerida (debug)

```json
// Response (200 OK)
{
  "financeiro": {
    "totalReceber": 150000.00,
    "totalPagar": 85000.00,
    "fluxoCaixa": 65000.00,
    "contasVencidas": 5
  },
  "obras": {
    "totalObras": 8,
    "obrasAtivas": 5,
    "obrasFinalizadas": 3,
    "faturamentoTotal": 750000.00
  },
  "funcionarios": {
    "totalFuncionarios": 15,
    "funcionariosAtivos": 12,
    "totalApontamentos": 45,
    "folhaPagamento": 28000.00
  },
  "fornecedores": {
    "totalFornecedores": 25,
    "fornecedoresAtivos": 20,
    "orcamentosAbertos": 8,
    "valorTotalOrcamentos": 95000.00
  }
}
```

### Dashboard Financeiro

**GET** `/dashboard/financeiro`

Retorna métricas financeiras consolidadas.

**Autenticação**: Não requerida (debug)

```json
// Response (200 OK)
{
  "contasReceber": {
    "total": 150000.00,
    "pendente": 120000.00,
    "vencido": 30000.00,
    "quantidade": 12
  },
  "contasPagar": {
    "total": 85000.00,
    "pendente": 65000.00,
    "vencido": 20000.00,
    "quantidade": 8
  },
  "fluxoCaixa": {
    "saldo": 65000.00,
    "entradas": 180000.00,
    "saidas": 115000.00,
    "projecao": 85000.00
  }
}
```

### Dashboard de Obras

**GET** `/dashboard/obras`

Retorna métricas de obras.

**Autenticação**: Não requerida (debug)

```json
// Response (200 OK)
{
  "resumo": {
    "total": 8,
    "ativas": 5,
    "finalizadas": 3,
    "planejamento": 2
  },
  "financeiro": {
    "faturamentoTotal": 750000.00,
    "valorRecebido": 450000.00,
    "valorPendente": 300000.00
  },
  "cronograma": {
    "noPrazo": 3,
    "atrasadas": 2,
    "concluidas": 3
  }
}
```

### Dashboard de Funcionários

**GET** `/dashboard/funcionarios`

Retorna métricas de funcionários.

**Autenticação**: Não requerida (debug)

```json
// Response (200 OK)
{
  "resumo": {
    "total": 15,
    "ativos": 12,
    "inativos": 3
  },
  "apontamentos": {
    "total": 45,
    "aprovados": 38,
    "pendentes": 7,
    "valorTotal": 28000.00
  },
  "produtividade": {
    "horasTotal": 1800,
    "mediaDiaria": 7.5,
    "funcionarioMaisAtivo": "Carlos Santos"
  }
}
```

### Dashboard de Fornecedores

**GET** `/dashboard/fornecedores`

Retorna métricas de fornecedores.

**Autenticação**: Não requerida (debug)

```json
// Response (200 OK)
{
  "resumo": {
    "total": 25,
    "ativos": 20,
    "inativos": 5
  },
  "orcamentos": {
    "total": 15,
    "aprovados": 7,
    "pendentes": 8,
    "valorTotal": 95000.00
  },
  "avaliacoes": {
    "mediaGeral": 4.2,
    "melhorAvaliado": "Materiais ABC Ltda",
    "piorAvaliado": "Fornecedor XYZ"
  }
}
```

### Fluxo de Caixa

**GET** `/dashboard/fluxo-caixa`

Retorna fluxo de caixa consolidado.

**Autenticação**: Não requerida (debug)

```json
// Response (200 OK)
{
  "saldoAtual": 65000.00,
  "entradas": {
    "mes": 45000.00,
    "trimestre": 120000.00,
    "ano": 380000.00
  },
  "saidas": {
    "mes": 32000.00,
    "trimestre": 85000.00,
    "ano": 280000.00
  },
  "projecoes": {
    "proximoMes": 75000.00,
    "proximoTrimestre": 95000.00
  },
  "principais": {
    "maiorEntrada": {
      "descricao": "Pagamento Casa Silva - 2ª etapa",
      "valor": 25000.00
    },
    "maiorSaida": {
      "descricao": "Folha de pagamento quinzenal",
      "valor": 15000.00
    }
  }
}
```

### Dashboard por Seção

**GET** `/dashboard/{secao}`

Retorna dados de uma seção específica do dashboard.

**Parâmetros**: `secao` pode ser `financeiro`, `obras`, `funcionarios` ou `fornecedores`

**Autenticação**: Não requerida (debug)

```json
// Response (200 OK)
// Retorna os mesmos dados das rotas específicas acima
```

### Parâmetros de Cache

**GET** `/dashboard/cache-info`

Retorna informações sobre cache do dashboard.

**Autenticação**: Não requerida (debug)

```json
// Response (200 OK)
{
  "cacheAtivo": true,
  "tempoExpiracao": "5m",
  "ultimaAtualizacao": "2024-02-20T14:30:00Z",
  "proximaAtualizacao": "2024-02-20T14:35:00Z",
  "hits": 245,
  "misses": 12
}
```

## Health Check

### Verificar Status da API

**GET** `/health`

Verifica se a API está funcionando.

**Autenticação**: Não requerida

```json
// Response (200 OK)
{
  "status": "ok"
}
```

## Códigos de Erro Padrão

### Códigos HTTP
- **200**: OK - Sucesso
- **201**: Created - Recurso criado com sucesso
- **204**: No Content - Sucesso sem conteúdo de resposta
- **400**: Bad Request - Dados inválidos
- **401**: Unauthorized - Não autenticado
- **403**: Forbidden - Sem permissão
- **404**: Not Found - Recurso não encontrado
- **409**: Conflict - Conflito de regra de negócio
- **500**: Internal Server Error - Erro interno do servidor

### Formato de Erro Padrão

```json
{
  "erro": {
    "codigo": "CODIGO_ERRO",
    "mensagem": "Descrição do erro",
    "detalhes": "Informações adicionais (opcional)"
  }
}
```

### Códigos de Erro Específicos

#### Autenticação
- `PAYLOAD_INVALIDO`: Dados enviados são inválidos
- `SENHAS_NAO_CONFEREM`: Senhas não coincidem no registro
- `CREDENCIAIS_INVALIDAS`: Email ou senha incorretos
- `TOKEN_INVALIDO`: Token JWT inválido ou expirado

#### Autorização
- `ACESSO_NEGADO`: Usuário sem permissão para a operação
- `PERMISSAO_INSUFICIENTE`: Permissão específica necessária

#### Recursos
- `RECURSO_NAO_ENCONTRADO`: Recurso solicitado não existe
- `FUNCIONARIO_NAO_ENCONTRADO`: Funcionário específico não encontrado
- `OBRA_NAO_ENCONTRADA`: Obra específica não encontrada
- `ETAPA_NAO_ENCONTRADA`: Etapa específica não encontrada

#### Regras de Negócio
- `CONFLITO_REGRA_NEGOCIO`: Violação de regra de negócio
- `FUNCIONARIO_ALOCADO`: Não é possível deletar funcionário alocado
- `ETAPA_EM_ANDAMENTO`: Etapa não pode ser deletada enquanto em andamento

#### Sistema
- `ERRO_INTERNO`: Erro interno do servidor
- `ERRO_BANCO_DADOS`: Erro na operação do banco de dados
- `METODO_NAO_PERMITIDO`: Método HTTP não suportado

## Paginação

A API utiliza paginação baseada em offset para listagens:

### Parâmetros de Query
- `page`: Número da página (padrão: 1)
- `limit`: Itens por página (padrão: 10, máximo: 100)

### Resposta Paginada
```json
{
  "dados": [...],
  "paginacao": {
    "paginaAtual": 1,
    "totalPaginas": 5,
    "totalItens": 50,
    "itensPorPagina": 10
  }
}
```

## Filtros

Muitos endpoints suportam filtros via query parameters:

### Filtros Comuns
- `nome`: Filtro por nome (busca parcial)
- `status`: Filtro por status específico
- `dataInicio`: Filtro por data de início (formato: YYYY-MM-DD)
- `dataFim`: Filtro por data de fim (formato: YYYY-MM-DD)

### Exemplo de URL com Filtros
```
GET /obras?page=2&limit=20&nome=casa&status=Em+Andamento&dataInicio=2024-01-01
```