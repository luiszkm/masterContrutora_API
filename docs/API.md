# API Documentation - Master Construtora

## Visão Geral

A API Master Construtora é uma REST API construída em Go que gerencia todos os aspectos de uma empresa de construção civil. A API utiliza autenticação JWT via cookies httpOnly e implementa autorização baseada em papéis (RBAC).

**Base URLs**:

- Desenvolvimento: `http://localhost:8080`
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

### Pagar Apontamento

**PATCH** `/apontamentos/{apontamentoId}/pagar`

Registra o pagamento de um apontamento aprovado.

**Permissão**: `pessoal:apontamento:pagar`

```json
// Request
{
  "contaBancariaId": "uuid-conta-bancaria",
  "apontamentoId": ["uuid-apontamento"]
}

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
  "status": "Pago"
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

Lista todas as categorias de produtos.

**Permissão**: `suprimentos:ler`

```json
// Response (200 OK)
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
      "itensCount": 2
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

### Registrar Pagamento

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

### Registrar Pagamentos em Lote

**POST** `/pagamentos/lote`

Registra múltiplos pagamentos de uma vez.

**Permissão**: `financeiro:escrever`

```json
// Request
{
  "pagamentos": [
    {
      "funcionarioId": "uuid-funcionario-1",
      "obraId": "uuid-obra",
      "periodoReferencia": "Fevereiro/2024",
      "valorCalculado": 1450.00,
      "contaBancariaId": "uuid-conta-bancaria"
    },
    {
      "funcionarioId": "uuid-funcionario-2",
      "obraId": "uuid-obra",
      "periodoReferencia": "Fevereiro/2024",
      "valorCalculado": 1800.00,
      "contaBancariaId": "uuid-conta-bancaria"
    }
  ]
}

// Response (207 Multi-Status)
{
  "sucessos": [
    {
      "funcionarioId": "uuid-funcionario-1",
      "pagamentoId": "uuid-pagamento-1",
      "valor": 1450.00
    }
  ],
  "erros": [
    {
      "funcionarioId": "uuid-funcionario-2",
      "erro": "Funcionário não encontrado"
    }
  ]
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