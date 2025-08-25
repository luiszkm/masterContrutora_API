# API Pessoal - Documentação Completa

## Visão Geral

A API Pessoal gerencia funcionários e apontamentos quinzenais de trabalho na aplicação Master Construtora. Esta documentação detalha todas as rotas, payloads de entrada, respostas e códigos de status.

**Importante**: O pagamento de apontamentos é realizado através do módulo financeiro. Quando um apontamento é aprovado, uma conta a pagar é criada automaticamente. O status do apontamento é atualizado para "PAGO" quando a conta correspondente é quitada no módulo financeiro.

**Formatos de Data**:
- Campos de entrada (períodos): `"YYYY-MM-DD"` (ex: "2025-01-15")
- Campos de resposta (timestamps): `"YYYY-MM-DDTHH:mm:ssZ"` (ex: "2025-01-15T10:30:00Z")

## Estruturas de Dados

### Funcionário

```json
{
  "id": "uuid",
  "nome": "string",
  "cpf": "string", 
  "telefone": "string",
  "cargo": "string",
  "email": "string (opcional)",
  "departamento": "string",
  "dataContratacao": "2025-01-01T00:00:00Z",
  "valorDiaria": 150.00,
  "chavePix": "string",
  "status": "ATIVO|INATIVO|DESLIGADO",
  "desligamentoData": "2025-01-01T00:00:00Z (opcional)",
  "motivoDesligamento": "string (opcional)",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z",
  "diaria": 150.00,
  "avaliacaoDesempenho": "string (opcional)",
  "observacoes": "string (opcional)"
}
```

### Apontamento Quinzenal

```json
{
  "id": "uuid",
  "funcionarioId": "uuid",
  "obraId": "uuid",
  "periodoInicio": "2025-01-01T00:00:00Z",
  "periodoFim": "2025-01-15T00:00:00Z",
  "diaria": 150.00,
  "diasTrabalhados": 15,
  "adicionais": 100.00,
  "descontos": 50.00,
  "adiantamentos": 200.00,
  "valorTotalCalculado": 2100.00,
  "status": "EM_ABERTO|APROVADO_PARA_PAGAMENTO|PAGO|CANCELADO",
  "createdAt": "2025-01-01T00:00:00Z",
  "updatedAt": "2025-01-01T00:00:00Z",
  "funcionarioNome": "Nome do Funcionário"
}
```

## Rotas - Funcionários

### 1. Cadastrar Funcionário

**Endpoint:** `POST /funcionarios`  
**Permissão:** `PESSOAL_ESCREVER`  
**Autenticação:** Obrigatória

#### Payload de Entrada:
```json
{
  "nome": "João Silva",
  "cpf": "12345678901",
  "cargo": "Pedreiro",
  "departamento": "Construção",
  "diaria": 150.00,
  "chavePix": "joao@email.com",
  "observacoes": "Funcionário experiente",
  "telefone": "(11) 99999-9999"
}
```

#### Resposta de Sucesso (201):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "nome": "João Silva",
  "cpf": "12345678901",
  "telefone": "(11) 99999-9999",
  "cargo": "Pedreiro",
  "email": null,
  "departamento": "Construção",
  "dataContratacao": "2025-01-24T10:30:00Z",
  "valorDiaria": 150.00,
  "chavePix": "joao@email.com",
  "status": "ATIVO",
  "desligamentoData": null,
  "motivoDesligamento": "",
  "created_at": "2025-01-24T10:30:00Z",
  "updated_at": "2025-01-24T10:30:00Z",
  "diaria": 150.00,
  "avaliacaoDesempenho": "",
  "observacoes": "Funcionário experiente"
}
```

#### Erros Possíveis:
- **400 Bad Request:** Payload inválido
- **500 Internal Server Error:** Erro interno do servidor

---

### 2. Listar Funcionários

**Endpoint:** `GET /funcionarios`  
**Permissão:** `PESSOAL_LER`  
**Autenticação:** Obrigatória

#### Resposta de Sucesso (200):
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "nome": "João Silva",
    "cpf": "12345678901",
    "telefone": "(11) 99999-9999",
    "cargo": "Pedreiro",
    "email": null,
    "departamento": "Construção",
    "dataContratacao": "2025-01-24T10:30:00Z",
    "valorDiaria": 150.00,
    "chavePix": "joao@email.com",
    "status": "ATIVO",
    "desligamentoData": null,
    "motivoDesligamento": "",
    "created_at": "2025-01-24T10:30:00Z",
    "updated_at": "2025-01-24T10:30:00Z",
    "diaria": 150.00,
    "avaliacaoDesempenho": "",
    "observacoes": "Funcionário experiente"
  }
]
```

#### Erros Possíveis:
- **500 Internal Server Error:** Erro ao listar funcionários

---

### 3. Buscar Funcionário por ID

**Endpoint:** `GET /funcionarios/{funcionarioId}`  
**Permissão:** `PESSOAL_LER`  
**Autenticação:** Obrigatória

#### Parâmetros de URL:
- `funcionarioId`: UUID do funcionário

#### Resposta de Sucesso (200):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "nome": "João Silva",
  "cpf": "12345678901",
  "telefone": "(11) 99999-9999",
  "cargo": "Pedreiro",
  "email": null,
  "departamento": "Construção",
  "dataContratacao": "2025-01-24T10:30:00Z",
  "valorDiaria": 150.00,
  "chavePix": "joao@email.com",
  "status": "ATIVO",
  "desligamentoData": null,
  "motivoDesligamento": "",
  "created_at": "2025-01-24T10:30:00Z",
  "updated_at": "2025-01-24T10:30:00Z",
  "diaria": 150.00,
  "avaliacaoDesempenho": "",
  "observacoes": "Funcionário experiente"
}
```

#### Erros Possíveis:
- **404 Not Found:** Funcionário não encontrado
- **500 Internal Server Error:** Erro ao buscar funcionário

---

### 4. Atualizar Funcionário

**Endpoint:** `PUT /funcionarios/{funcionarioId}`  
**Permissão:** `PESSOAL_ESCREVER`  
**Autenticação:** Obrigatória

#### Parâmetros de URL:
- `funcionarioId`: UUID do funcionário

#### Payload de Entrada (todos os campos opcionais):
```json
{
  "nome": "João Santos Silva",
  "cpf": "12345678901",
  "cargo": "Mestre de Obras",
  "departamento": "Construção",
  "valorDiaria": 180.00,
  "chavePix": "joao.santos@email.com",
  "status": "ATIVO",
  "telefone": "(11) 88888-8888",
  "motivoDesligamento": null,
  "dataContratacao": "2025-01-24",
  "desligamentoData": null,
  "observacoes": "Promovido a mestre de obras",
  "avaliacaoDesempenho": "Excelente",
  "email": "joao.santos@email.com"
}
```

#### Resposta de Sucesso (200):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "nome": "João Santos Silva",
  "cpf": "12345678901",
  "telefone": "(11) 88888-8888",
  "cargo": "Mestre de Obras",
  "email": "joao.santos@email.com",
  "departamento": "Construção",
  "dataContratacao": "2025-01-24T00:00:00Z",
  "valorDiaria": 180.00,
  "chavePix": "joao.santos@email.com",
  "status": "ATIVO",
  "desligamentoData": null,
  "motivoDesligamento": "",
  "created_at": "2025-01-24T10:30:00Z",
  "updated_at": "2025-01-24T11:30:00Z",
  "diaria": 180.00,
  "avaliacaoDesempenho": "Excelente",
  "observacoes": "Promovido a mestre de obras"
}
```

#### Erros Possíveis:
- **400 Bad Request:** ID inválido ou payload inválido
- **404 Not Found:** Funcionário não encontrado
- **500 Internal Server Error:** Erro interno

---

### 5. Ativar Funcionário

**Endpoint:** `PATCH /funcionarios/{funcionarioId}/ativar`  
**Permissão:** `PESSOAL_APONTAMENTO_LER`  
**Autenticação:** Obrigatória

#### Parâmetros de URL:
- `funcionarioId`: UUID do funcionário

#### Resposta de Sucesso (204):
```
Status: 204 No Content
```

#### Erros Possíveis:
- **404 Not Found:** Funcionário não encontrado
- **500 Internal Server Error:** Erro ao buscar funcionário

---

### 6. Deletar Funcionário

**Endpoint:** `DELETE /funcionarios/{funcionarioId}`  
**Permissão:** `PESSOAL_ESCREVER`  
**Autenticação:** Obrigatória

#### Parâmetros de URL:
- `funcionarioId`: UUID do funcionário

#### Resposta de Sucesso (204):
```
Status: 204 No Content
```

#### Erros Possíveis:
- **404 Not Found:** Funcionário não encontrado
- **409 Conflict:** Funcionário possui alocações ativas (não pode ser deletado)
- **500 Internal Server Error:** Erro interno

---

### 7. Listar Funcionários com Último Apontamento

**Endpoint:** `GET /funcionarios/apontamentos`  
**Permissão:** `PESSOAL_LER`  
**Autenticação:** Obrigatória

#### Parâmetros de Query (opcionais):
- `page`: Número da página (padrão: 1)
- `pageSize`: Itens por página (padrão: 10)
- `status`: Filtro por status

#### Resposta de Sucesso (200):
```json
[
  {
    "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
    "funcionarioNome": "João Silva",
    "ultimoApontamento": {
      "id": "456e7890-e89b-12d3-a456-426614174001",
      "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
      "obraId": "789e0123-e89b-12d3-a456-426614174002",
      "periodoInicio": "2025-01-01T00:00:00Z",
      "periodoFim": "2025-01-15T00:00:00Z",
      "diaria": 150.00,
      "diasTrabalhados": 15,
      "adicionais": 100.00,
      "descontos": 50.00,
      "adiantamentos": 200.00,
      "valorTotalCalculado": 2100.00,
      "status": "EM_ABERTO",
      "createdAt": "2025-01-01T00:00:00Z",
      "updatedAt": "2025-01-01T00:00:00Z",
      "funcionarioNome": "João Silva"
    }
  }
]
```

#### Erros Possíveis:
- **500 Internal Server Error:** Erro interno

---

### 8. Listar Apontamentos por Funcionário

**Endpoint:** `GET /funcionarios/{funcionarioId}/apontamentos`  
**Permissão:** `PESSOAL_APONTAMENTO_LER`  
**Autenticação:** Obrigatória

#### Parâmetros de URL:
- `funcionarioId`: UUID do funcionário

#### Parâmetros de Query (opcionais):
- `page`: Número da página (padrão: 1)
- `pageSize`: Itens por página (padrão: 10)
- `status`: Filtro por status

#### Resposta de Sucesso (200):
```json
{
  "data": [
    {
      "id": "456e7890-e89b-12d3-a456-426614174001",
      "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
      "obraId": "789e0123-e89b-12d3-a456-426614174002",
      "periodoInicio": "2025-01-01T00:00:00Z",
      "periodoFim": "2025-01-15T00:00:00Z",
      "diaria": 150.00,
      "diasTrabalhados": 15,
      "adicionais": 100.00,
      "descontos": 50.00,
      "adiantamentos": 200.00,
      "valorTotalCalculado": 2100.00,
      "status": "EM_ABERTO",
      "createdAt": "2025-01-01T00:00:00Z",
      "updatedAt": "2025-01-01T00:00:00Z",
      "funcionarioNome": "João Silva"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "total": 1,
    "totalPages": 1
  }
}
```

---

## Rotas - Apontamentos

### 1. Criar Apontamento

**Endpoint:** `POST /apontamentos`  
**Permissão:** `PESSOAL_APONTAMENTO_ESCREVER`  
**Autenticação:** Obrigatória

#### Payload de Entrada:
```json
{
  "FuncionarioID": "123e4567-e89b-12d3-a456-426614174000",
  "ObraID": "789e0123-e89b-12d3-a456-426614174002",
  "PeriodoInicio": "2025-01-01",
  "PeriodoFim": "2025-01-15",
  "Diaria": 150.00,
  "DiasTrabalhados": 15,
  "ValorAdicional": 100.00,
  "Descontos": 50.00,
  "Adiantamento": 200.00
}
```

#### Resposta de Sucesso (201):
```json
{
  "id": "456e7890-e89b-12d3-a456-426614174001",
  "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
  "obraId": "789e0123-e89b-12d3-a456-426614174002",
  "periodoInicio": "2025-01-01T00:00:00Z",
  "periodoFim": "2025-01-15T00:00:00Z",
  "diaria": 150.00,
  "diasTrabalhados": 15,
  "adicionais": 100.00,
  "descontos": 50.00,
  "adiantamentos": 200.00,
  "valorTotalCalculado": 2100.00,
  "status": "EM_ABERTO",
  "createdAt": "2025-01-24T10:30:00Z",
  "updatedAt": "2025-01-24T10:30:00Z",
  "funcionarioNome": "João Silva"
}
```

#### Erros Possíveis:
- **400 Bad Request:** Payload inválido
- **500 Internal Server Error:** Erro interno

---

### 2. Listar Apontamentos

**Endpoint:** `GET /apontamentos`  
**Permissão:** `PESSOAL_APONTAMENTO_LER`  
**Autenticação:** Obrigatória

#### Parâmetros de Query (opcionais):
- `page`: Número da página (padrão: 1)
- `pageSize`: Itens por página (padrão: 10)
- `status`: Filtro por status

#### Resposta de Sucesso (200):
```json
{
  "data": [
    {
      "id": "456e7890-e89b-12d3-a456-426614174001",
      "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
      "obraId": "789e0123-e89b-12d3-a456-426614174002",
      "periodoInicio": "2025-01-01T00:00:00Z",
      "periodoFim": "2025-01-15T00:00:00Z",
      "diaria": 150.00,
      "diasTrabalhados": 15,
      "adicionais": 100.00,
      "descontos": 50.00,
      "adiantamentos": 200.00,
      "valorTotalCalculado": 2100.00,
      "status": "EM_ABERTO",
      "createdAt": "2025-01-01T00:00:00Z",
      "updatedAt": "2025-01-01T00:00:00Z",
      "funcionarioNome": "João Silva"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "total": 1,
    "totalPages": 1
  }
}
```

#### Erros Possíveis:
- **500 Internal Server Error:** Erro ao listar apontamentos

---

### 3. Atualizar Apontamento

**Endpoint:** `PUT /funcionarios/apontamentos/{apontamentoId}`  
**Permissão:** `PESSOAL_ESCREVER`  
**Autenticação:** Obrigatória

#### Parâmetros de URL:
- `apontamentoId`: UUID do apontamento

#### Payload de Entrada:
```json
{
  "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
  "obraId": "789e0123-e89b-12d3-a456-426614174002",
  "periodoInicio": "2025-01-01",
  "periodoFim": "2025-01-15",
  "diaria": 160.00,
  "diasTrabalhados": 14,
  "valorAdicional": 120.00,
  "descontos": 30.00,
  "adiantamento": 150.00,
  "status": "EM_ABERTO"
}
```

#### Resposta de Sucesso (200):
```json
{
  "id": "456e7890-e89b-12d3-a456-426614174001",
  "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
  "obraId": "789e0123-e89b-12d3-a456-426614174002",
  "periodoInicio": "2025-01-01T00:00:00Z",
  "periodoFim": "2025-01-15T00:00:00Z",
  "diaria": 160.00,
  "diasTrabalhados": 14,
  "adicionais": 120.00,
  "descontos": 30.00,
  "adiantamentos": 150.00,
  "valorTotalCalculado": 2180.00,
  "status": "EM_ABERTO",
  "createdAt": "2025-01-01T00:00:00Z",
  "updatedAt": "2025-01-24T11:30:00Z",
  "funcionarioNome": "João Silva"
}
```

#### Erros Possíveis:
- **400 Bad Request:** ID inválido ou payload inválido
- **404 Not Found:** Apontamento não encontrado
- **500 Internal Server Error:** Erro interno

---

### 4. Aprovar Apontamento

**Endpoint:** `PATCH /apontamentos/{apontamentoId}/aprovar`  
**Permissão:** `PESSOAL_APONTAMENTO_APROVAR`  
**Autenticação:** Obrigatória

#### Parâmetros de URL:
- `apontamentoId`: UUID do apontamento

#### Resposta de Sucesso (200):
```json
{
  "id": "456e7890-e89b-12d3-a456-426614174001",
  "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
  "obraId": "789e0123-e89b-12d3-a456-426614174002",
  "periodoInicio": "2025-01-01T00:00:00Z",
  "periodoFim": "2025-01-15T00:00:00Z",
  "diaria": 150.00,
  "diasTrabalhados": 15,
  "adicionais": 100.00,
  "descontos": 50.00,
  "adiantamentos": 200.00,
  "valorTotalCalculado": 2100.00,
  "status": "APROVADO_PARA_PAGAMENTO",
  "createdAt": "2025-01-01T00:00:00Z",
  "updatedAt": "2025-01-24T12:00:00Z",
  "funcionarioNome": "João Silva"
}
```

#### Erros Possíveis:
- **404 Not Found:** Apontamento não encontrado
- **409 Conflict:** Regra de negócio violada (ex: apontamento já foi pago)

---

### 5. Cancelar Apontamento

**Endpoint:** `PATCH /apontamentos/{apontamentoId}/cancelar`  
**Permissão:** `PESSOAL_APONTAMENTO_ESCREVER`  
**Autenticação:** Obrigatória

#### Parâmetros de URL:
- `apontamentoId`: UUID do apontamento

#### Payload de Entrada:
```json
{
  "motivoCancelamento": "Funcionário solicitou cancelamento do apontamento"
}
```

#### Resposta de Sucesso (200):
```json
{
  "id": "456e7890-e89b-12d3-a456-426614174001",
  "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
  "obraId": "789e0123-e89b-12d3-a456-426614174002",
  "periodoInicio": "2025-01-01T00:00:00Z",
  "periodoFim": "2025-01-15T00:00:00Z",
  "diaria": 150.00,
  "diasTrabalhados": 15,
  "adicionais": 100.00,
  "descontos": 50.00,
  "adiantamentos": 200.00,
  "valorTotalCalculado": 2100.00,
  "status": "CANCELADO",
  "createdAt": "2025-01-01T00:00:00Z",
  "updatedAt": "2025-01-24T14:30:00Z",
  "funcionarioNome": "João Silva"
}
```

#### Erros Possíveis:
- **404 Not Found:** Apontamento não encontrado
- **409 Conflict:** Regra de negócio violada (ex: apontamento já foi pago)
- **500 Internal Server Error:** Erro interno

---

### 6. Replicar Apontamentos para Próxima Quinzena

**Endpoint:** `POST /funcionarios/apontamentos/replicar`  
**Permissão:** `PESSOAL_APONTAMENTO_ESCREVER`  
**Autenticação:** Obrigatória

#### Payload de Entrada:
```json
{
  "funcionarioIds": [
    "123e4567-e89b-12d3-a456-426614174000",
    "456e7890-e89b-12d3-a456-426614174001"
  ]
}
```

#### Resposta de Sucesso (207 Multi-Status):
```json
{
  "resumo": {
    "totalSolicitado": 2,
    "totalSucesso": 1,
    "totalFalha": 1
  },
  "sucessos": [
    {
      "funcionarioId": "123e4567-e89b-12d3-a456-426614174000",
      "novoApontamentoId": "789e0123-e89b-12d3-a456-426614174003"
    }
  ],
  "falhas": [
    {
      "funcionarioId": "456e7890-e89b-12d3-a456-426614174001",
      "motivo": "Funcionário não possui apontamento na quinzena atual"
    }
  ]
}
```

#### Erros Possíveis:
- **400 Bad Request:** Payload inválido ou lista vazia
- **500 Internal Server Error:** Falha inesperada na replicação

---

## Status de Apontamentos

Os apontamentos seguem um ciclo de vida com os seguintes status:

- **EM_ABERTO**: Apontamento criado, pode ser editado e cancelado
- **APROVADO_PARA_PAGAMENTO**: Apontamento aprovado, conta a pagar criada automaticamente no módulo financeiro, pode ser cancelado
- **PAGO**: Status atualizado automaticamente quando a conta correspondente é quitada no módulo financeiro, não pode mais ser alterado ou cancelado
- **CANCELADO**: Apontamento cancelado, conta a pagar correspondente é cancelada automaticamente no módulo financeiro, não pode mais ser alterado

## Status de Funcionários

Os funcionários possuem os seguintes status possíveis:

- **ATIVO**: Funcionário ativo na empresa
- **INATIVO**: Funcionário temporariamente inativo
- **DESLIGADO**: Funcionário desligado da empresa

## Permissões Necessárias

- **PESSOAL_LER**: Visualizar funcionários e apontamentos
- **PESSOAL_ESCREVER**: Criar e editar funcionários
- **PESSOAL_APONTAMENTO_LER**: Visualizar apontamentos
- **PESSOAL_APONTAMENTO_ESCREVER**: Criar e editar apontamentos
- **PESSOAL_APONTAMENTO_APROVAR**: Aprovar apontamentos

## Regras de Negócio

1. **Funcionários alocados não podem ser deletados** - retorna erro 409
2. **Apontamentos pagos não podem ser editados ou cancelados** - retorna erro 409
3. **Apontamentos cancelados não podem ser editados** - retorna erro 409
4. **Só é possível aprovar apontamentos em aberto**
5. **Só é possível cancelar apontamentos em aberto ou aprovados para pagamento**
6. **Aprovação de apontamento cria conta a pagar automaticamente** no módulo financeiro
7. **Cancelamento de apontamento cancela automaticamente a conta a pagar** no módulo financeiro
8. **Status PAGO é atualizado automaticamente** quando conta é quitada no módulo financeiro
9. **Replicação considera apenas apontamentos da quinzena atual**
10. **Valor total é calculado automaticamente**: (diasTrabalhados × diaria) + adicionais - descontos - adiantamentos

## Integração com Módulo Financeiro

O módulo Pessoal integra-se automaticamente com o módulo Financeiro através de eventos:

### Fluxo de Aprovação
1. **Apontamento Aprovado** → Publica evento `pessoal:apontamento_aprovado`
2. **Módulo Financeiro** → Escuta o evento e cria uma **Conta a Pagar** automaticamente
3. **Conta a Pagar** → Contém informações do funcionário, valor calculado e data de vencimento
4. **Pagamento da Conta** → Quando conta é paga no módulo financeiro, status do apontamento muda para `PAGO`

### Fluxo de Cancelamento
1. **Apontamento Cancelado** → Publica evento `pessoal:apontamento_cancelado`
2. **Módulo Financeiro** → Escuta o evento e **cancela automaticamente** a conta a pagar correspondente
3. **Conta a Pagar** → Status alterado para `CANCELADO`, não pode mais ser paga
4. **Regras**: Só cancela contas pendentes (sem pagamentos realizados)

### Dados Criados Automaticamente na Conta a Pagar
- **Fornecedor**: Nome do funcionário
- **Tipo**: `FUNCIONARIO`
- **Descrição**: "Pagamento de funcionário - [Nome] ([Período])"
- **Valor**: Valor total calculado do apontamento
- **Vencimento**: Data prevista para pagamento (7 dias após aprovação)
- **Número Documento**: ID do apontamento como referência

### Eventos Relacionados
- `pessoal:apontamento_aprovado`: Disparado quando apontamento é aprovado
- `pessoal:apontamento_cancelado`: Disparado quando apontamento é cancelado
- `financeiro:conta_pagar_criada`: Disparado quando conta é criada automaticamente
- `financeiro:conta_pagar_paga`: Quando conta é paga, status do apontamento muda para PAGO
- `financeiro:conta_pagar_cancelada`: Disparado quando conta é cancelada automaticamente

## Códigos de Erro Padronizados

- **PAYLOAD_INVALIDO**: Dados da requisição inválidos
- **ID_INVALIDO**: UUID fornecido não é válido
- **FUNCIONARIO_NAO_ENCONTRADO**: Funcionário não existe
- **APONTAMENTO_NAO_ENCONTRADO**: Apontamento não existe
- **CONFLITO_REGRA_NEGOCIO**: Violação de regra de negócio
- **REGRA_NEGOCIO_VIOLADA**: Operação não permitida pelo estado atual
- **DADOS_OBRIGATORIOS**: Campos obrigatórios não informados
- **ERRO_INTERNO**: Erro interno do servidor