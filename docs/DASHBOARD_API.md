# Dashboard API - Documentação Completa

## Visão Geral

A API de Dashboard fornece uma visão abrangente dos dados da construtora, incluindo informações financeiras, de obras, funcionários e fornecedores. Todos os endpoints são protegidos por autenticação e autorização baseada em permissões.

## Base URL
```
http://localhost:8080/dashboard
```

> **Nota**: Durante desenvolvimento, o servidor pode rodar em portas alternativas (ex: 8081). Verifique os logs de inicialização para a porta correta.

## Autenticação
Todos os endpoints requerem autenticação via JWT:
```
Authorization: Bearer <token>
```

## Parâmetros Comuns

### Filtros de Período
- `dataInicio` (opcional): Data de início no formato YYYY-MM-DD (padrão: 6 meses atrás)
- `dataFim` (opcional): Data de fim no formato YYYY-MM-DD (padrão: hoje)

### Filtros Avançados
- `secoes` (opcional): Lista de seções separadas por vírgula (`financeiro,obras,funcionarios,fornecedores`)
- `obraIds` (opcional): Lista de IDs de obras específicas
- `fornecedorIds` (opcional): Lista de IDs de fornecedores específicos
- `incluirInativos` (opcional): Incluir registros inativos (padrão: false)

---

## 1. Dashboard Completo

### Endpoint
```http
GET /dashboard
```

### Permissões
- Requer: `PermissaoObrasLer` (permissão básica)

### Exemplo de Requisição
```http
GET /dashboard?dataInicio=2024-01-01&dataFim=2024-12-31&secoes=financeiro,obras
Authorization: Bearer <token>
```

### Payload de Resposta
```json
{
  "resumoGeral": {
    "totalObras": 15,
    "obrasEmAndamento": 8,
    "totalFuncionarios": 45,
    "funcionariosAtivos": 42,
    "totalFornecedores": 28,
    "fornecedoresAtivos": 25,
    "saldoFinanceiroAtual": 125000.50,
    "totalInvestido": 850000.00,
    "progressoMedioObras": 67.5,
    "obrasEmAtraso": 2,
    "percentualAtraso": 25.0
  },
  "alertas": {
    "obrasComAtraso": [
      "Residencial Vila Verde",
      "Edifício São Paulo"
    ],
    "fornecedoresInativos": [
      "Materiais XYZ Ltda"
    ],
    "funcionariosSemApontamento": [
      "João Silva",
      "Maria Santos"
    ],
    "orcamentosPendentes": 5,
    "pagamentosPendentes": 3
  },
  "financeiro": {
    "fluxoCaixa": {
      "totalEntradas": 450000.00,
      "totalSaidas": 380000.00,
      "saldoAtual": 70000.00,
      "fluxoPorPeriodo": [
        {
          "periodo": "2024-01-01T00:00:00Z",
          "entradas": 75000.00,
          "saidas": 65000.00,
          "saldoLiquido": 10000.00
        },
        {
          "periodo": "2024-02-01T00:00:00Z",
          "entradas": 80000.00,
          "saidas": 70000.00,
          "saldoLiquido": 10000.00
        }
      ],
      "tendenciaMensal": "crescente"
    },
    "distribuicaoDespesas": {
      "totalGasto": 380000.00,
      "distribuicao": [
        {
          "categoria": "Mão de Obra",
          "valor": 180000.00,
          "percentual": 47.37,
          "quantidadeItens": 156
        },
        {
          "categoria": "Materiais de Construção",
          "valor": 120000.00,
          "percentual": 31.58,
          "quantidadeItens": 89
        },
        {
          "categoria": "Equipamentos",
          "valor": 80000.00,
          "percentual": 21.05,
          "quantidadeItens": 45
        }
      ],
      "maiorCategoria": "Mão de Obra",
      "valorMaiorCategoria": 180000.00
    },
    "ultimaAtualizacao": "2024-08-03T14:30:00Z"
  },
  "obras": {
    "progresso": {
      "progressoMedio": 67.5,
      "obrasEmAndamento": 8,
      "obrasConcluidas": 5,
      "totalObras": 15,
      "progressoPorObra": [
        {
          "obraId": "uuid-obra-1",
          "nomeObra": "Residencial Jardim das Flores",
          "percentualConcluido": 85.0,
          "etapasConcluidas": 8,
          "etapasTotal": 10,
          "status": "Em Andamento",
          "dataInicio": "2024-01-15T00:00:00Z",
          "dataFimPrevista": "2024-09-30T00:00:00Z"
        }
      ]
    },
    "distribuicao": {
      "totalObras": 15,
      "distribuicaoPorStatus": [
        {
          "status": "Em Andamento",
          "quantidade": 8,
          "percentual": 53.33,
          "valorTotal": 650000.00
        },
        {
          "status": "Concluída",
          "quantidade": 5,
          "percentual": 33.33,
          "valorTotal": 450000.00
        },
        {
          "status": "Pausada",
          "quantidade": 2,
          "percentual": 13.33,
          "valorTotal": 180000.00
        }
      ],
      "statusMaisComum": "Em Andamento"
    },
    "tendencias": {
      "obrasEmAtraso": 2,
      "obrasNoPrazo": 6,
      "percentualAtraso": 25.0,
      "tendenciaMensal": [
        {
          "periodo": "2024-01-01T00:00:00Z",
          "obrasIniciadas": 3,
          "obrasConcluidas": 1,
          "obrasEmAtraso": 0
        },
        {
          "periodo": "2024-02-01T00:00:00Z",
          "obrasIniciadas": 2,
          "obrasConcluidas": 2,
          "obrasEmAtraso": 1
        }
      ],
      "previsaoConclusaoMes": 3,
      "tendenciaGeral": "melhorando"
    },
    "ultimaAtualizacao": "2024-08-03T14:30:00Z"
  },
  "funcionarios": {
    "produtividade": {
      "mediaGeralProdutividade": 12.5,
      "totalFuncionarios": 45,
      "funcionariosAtivos": 42,
      "produtividadePorFuncionario": [
        {
          "funcionarioId": "uuid-func-1",
          "nomeFuncionario": "Carlos Pereira",
          "cargo": "Pedreiro",
          "diasTrabalhados": 180,
          "mediaDiasPorPeriodo": 14.5,
          "indiceProdutividade": 96.67,
          "obrasAlocadas": 3
        }
      ],
      "top5Produtivos": [
        {
          "funcionarioId": "uuid-func-1",
          "nomeFuncionario": "Carlos Pereira",
          "cargo": "Pedreiro",
          "diasTrabalhados": 180,
          "mediaDiasPorPeriodo": 14.5,
          "indiceProdutividade": 96.67,
          "obrasAlocadas": 3
        }
      ]
    },
    "custosMaoObra": {
      "custoTotalMaoObra": 180000.00,
      "custoMedioFuncionario": 4285.71,
      "custoMedioObra": 22500.00,
      "custosPorFuncionario": [
        {
          "funcionarioId": "uuid-func-1",
          "nomeFuncionario": "Carlos Pereira",
          "cargo": "Pedreiro",
          "custoTotal": 8500.00,
          "custoMedio": 850.00,
          "valorDiaria": 120.00,
          "periodosTrabalho": 10
        }
      ],
      "custosPorObra": [
        {
          "obraId": "uuid-obra-1",
          "nomeObra": "Residencial Jardim das Flores",
          "custoTotal": 35000.00,
          "custoMedio": 3500.00,
          "numFuncionarios": 8
        }
      ]
    },
    "topFuncionarios": {
      "top5Funcionarios": [
        {
          "funcionarioId": "uuid-func-1",
          "nomeFuncionario": "Carlos Pereira",
          "cargo": "Pedreiro",
          "avaliacaoDesempenho": "Excelente",
          "notaAvaliacao": 9.5,
          "diasTrabalhadosTotal": 180,
          "obrasParticipadas": 3,
          "dataContratacao": "2023-01-15T00:00:00Z"
        }
      ],
      "criterioAvaliacao": "Tempo de empresa + Produtividade"
    },
    "ultimaAtualizacao": "2024-08-03T14:30:00Z"
  },
  "fornecedores": {
    "fornecedoresPorCategoria": {
      "totalFornecedores": 28,
      "totalCategorias": 8,
      "distribuicaoPorCategoria": [
        {
          "categoriaId": "uuid-cat-1",
          "categoriaNome": "Materiais de Construção",
          "quantidadeFornecedores": 12,
          "percentual": 42.86,
          "avaliacaoMedia": 4.2
        },
        {
          "categoriaId": "uuid-cat-2",
          "categoriaNome": "Equipamentos",
          "quantidadeFornecedores": 8,
          "percentual": 28.57,
          "avaliacaoMedia": 4.5
        }
      ],
      "categoriaMaisPopular": "Materiais de Construção",
      "categoriaComMelhorAvaliacao": "Equipamentos"
    },
    "topFornecedores": {
      "top5Fornecedores": [
        {
          "fornecedorId": "uuid-forn-1",
          "nomeFornecedor": "Construtora ABC Ltda",
          "cnpj": "12.345.678/0001-90",
          "avaliacao": 4.8,
          "status": "Ativo",
          "totalOrcamentos": 15,
          "valorTotalGasto": 85000.00,
          "ultimoOrcamento": "2024-07-15T10:00:00Z",
          "categorias": ["Materiais de Construção", "Equipamentos"]
        }
      ],
      "criterioAvaliacao": "Avaliação + Volume de negócios",
      "avaliacaoMedia": 4.3,
      "fornecedoresAtivos": 25
    },
    "gastosFornecedores": {
      "totalGastoFornecedores": 320000.00,
      "gastoMedioFornecedor": 11428.57,
      "top10Gastos": [
        {
          "fornecedorId": "uuid-forn-1",
          "nomeFornecedor": "Construtora ABC Ltda",
          "avaliacao": 4.8,
          "valorTotalGasto": 85000.00,
          "quantidadeOrcamentos": 15,
          "ticketMedio": 5666.67,
          "ultimoOrcamento": "2024-07-15T10:00:00Z",
          "percentual": 26.56
        }
      ],
      "fornecedorMaiorGasto": "Construtora ABC Ltda",
      "valorMaiorGasto": 85000.00
    },
    "estatisticasGerais": {
      "totalFornecedores": 28,
      "fornecedoresAtivos": 25,
      "fornecedoresInativos": 3,
      "avaliacaoMediaGeral": 4.2,
      "tempoMedioContrato": 456
    },
    "ultimaAtualizacao": "2024-08-03T14:30:00Z"
  },
  "ultimaAtualizacao": "2024-08-03T14:30:00Z",
  "versaoCache": "1.0"
}
```

---

## 2. Dashboard Financeiro

### Endpoint
```http
GET /dashboard/financeiro
```

### Permissões
- Requer: `PermissaoFinanceiroLer`

### Payload de Resposta
```json
{
  "fluxoCaixa": {
    "totalEntradas": 450000.00,
    "totalSaidas": 380000.00,
    "saldoAtual": 70000.00,
    "fluxoPorPeriodo": [
      {
        "periodo": "2024-01-01T00:00:00Z",
        "entradas": 75000.00,
        "saidas": 65000.00,
        "saldoLiquido": 10000.00
      },
      {
        "periodo": "2024-02-01T00:00:00Z",
        "entradas": 80000.00,
        "saidas": 70000.00,
        "saldoLiquido": 10000.00
      }
    ],
    "tendenciaMensal": "crescente"
  },
  "distribuicaoDespesas": {
    "totalGasto": 380000.00,
    "distribuicao": [
      {
        "categoria": "Mão de Obra",
        "valor": 180000.00,
        "percentual": 47.37,
        "quantidadeItens": 156
      },
      {
        "categoria": "Materiais de Construção",
        "valor": 120000.00,
        "percentual": 31.58,
        "quantidadeItens": 89
      }
    ],
    "maiorCategoria": "Mão de Obra",
    "valorMaiorCategoria": 180000.00
  },
  "ultimaAtualizacao": "2024-08-03T14:30:00Z"
}
```

---

## 3. Dashboard Obras

### Endpoint
```http
GET /dashboard/obras
```

### Permissões
- Requer: `PermissaoObrasLer`

### Payload de Resposta
```json
{
  "progresso": {
    "progressoMedio": 67.5,
    "obrasEmAndamento": 8,
    "obrasConcluidas": 5,
    "totalObras": 15,
    "progressoPorObra": [
      {
        "obraId": "uuid-obra-1",
        "nomeObra": "Residencial Jardim das Flores",
        "percentualConcluido": 85.0,
        "etapasConcluidas": 8,
        "etapasTotal": 10,
        "status": "Em Andamento",
        "dataInicio": "2024-01-15T00:00:00Z",
        "dataFimPrevista": "2024-09-30T00:00:00Z"
      }
    ]
  },
  "distribuicao": {
    "totalObras": 15,
    "distribuicaoPorStatus": [
      {
        "status": "Em Andamento",
        "quantidade": 8,
        "percentual": 53.33,
        "valorTotal": 650000.00
      },
      {
        "status": "Concluída",
        "quantidade": 5,
        "percentual": 33.33,
        "valorTotal": 450000.00
      }
    ],
    "statusMaisComum": "Em Andamento"
  },
  "tendencias": {
    "obrasEmAtraso": 2,
    "obrasNoPrazo": 6,
    "percentualAtraso": 25.0,
    "tendenciaMensal": [
      {
        "periodo": "2024-01-01T00:00:00Z",
        "obrasIniciadas": 3,
        "obrasConcluidas": 1,
        "obrasEmAtraso": 0
      }
    ],
    "previsaoConclusaoMes": 3,
    "tendenciaGeral": "melhorando"
  },
  "ultimaAtualizacao": "2024-08-03T14:30:00Z"
}
```

---

## 4. Dashboard Funcionários

### Endpoint
```http
GET /dashboard/funcionarios
```

### Permissões
- Requer: `PermissaoPessoalLer`

### Payload de Resposta
```json
{
  "produtividade": {
    "mediaGeralProdutividade": 12.5,
    "totalFuncionarios": 45,
    "funcionariosAtivos": 42,
    "produtividadePorFuncionario": [
      {
        "funcionarioId": "uuid-func-1",
        "nomeFuncionario": "Carlos Pereira",
        "cargo": "Pedreiro",
        "diasTrabalhados": 180,
        "mediaDiasPorPeriodo": 14.5,
        "indiceProdutividade": 96.67,
        "obrasAlocadas": 3
      }
    ],
    "top5Produtivos": [
      {
        "funcionarioId": "uuid-func-1",
        "nomeFuncionario": "Carlos Pereira",
        "cargo": "Pedreiro",
        "diasTrabalhados": 180,
        "mediaDiasPorPeriodo": 14.5,
        "indiceProdutividade": 96.67,
        "obrasAlocadas": 3
      }
    ]
  },
  "custosMaoObra": {
    "custoTotalMaoObra": 180000.00,
    "custoMedioFuncionario": 4285.71,
    "custoMedioObra": 22500.00,
    "custosPorFuncionario": [
      {
        "funcionarioId": "uuid-func-1",
        "nomeFuncionario": "Carlos Pereira",
        "cargo": "Pedreiro",
        "custoTotal": 8500.00,
        "custoMedio": 850.00,
        "valorDiaria": 120.00,
        "periodosTrabalho": 10
      }
    ],
    "custosPorObra": [
      {
        "obraId": "uuid-obra-1",
        "nomeObra": "Residencial Jardim das Flores",
        "custoTotal": 35000.00,
        "custoMedio": 3500.00,
        "numFuncionarios": 8
      }
    ]
  },
  "topFuncionarios": {
    "top5Funcionarios": [
      {
        "funcionarioId": "uuid-func-1",
        "nomeFuncionario": "Carlos Pereira",
        "cargo": "Pedreiro",
        "avaliacaoDesempenho": "Excelente",
        "notaAvaliacao": 9.5,
        "diasTrabalhadosTotal": 180,
        "obrasParticipadas": 3,
        "dataContratacao": "2023-01-15T00:00:00Z"
      }
    ],
    "criterioAvaliacao": "Tempo de empresa + Produtividade"
  },
  "ultimaAtualizacao": "2024-08-03T14:30:00Z"
}
```

---

## 5. Dashboard Fornecedores

### Endpoint
```http
GET /dashboard/fornecedores
```

### Permissões
- Requer: `PermissaoSuprimentosLer`

### Payload de Resposta
```json
{
  "fornecedoresPorCategoria": {
    "totalFornecedores": 28,
    "totalCategorias": 8,
    "distribuicaoPorCategoria": [
      {
        "categoriaId": "uuid-cat-1",
        "categoriaNome": "Materiais de Construção",
        "quantidadeFornecedores": 12,
        "percentual": 42.86,
        "avaliacaoMedia": 4.2
      }
    ],
    "categoriaMaisPopular": "Materiais de Construção",
    "categoriaComMelhorAvaliacao": "Equipamentos"
  },
  "topFornecedores": {
    "top5Fornecedores": [
      {
        "fornecedorId": "uuid-forn-1",
        "nomeFornecedor": "Construtora ABC Ltda",
        "cnpj": "12.345.678/0001-90",
        "avaliacao": 4.8,
        "status": "Ativo",
        "totalOrcamentos": 15,
        "valorTotalGasto": 85000.00,
        "ultimoOrcamento": "2024-07-15T10:00:00Z",
        "categorias": ["Materiais de Construção", "Equipamentos"]
      }
    ],
    "criterioAvaliacao": "Avaliação + Volume de negócios",
    "avaliacaoMedia": 4.3,
    "fornecedoresAtivos": 25
  },
  "gastosFornecedores": {
    "totalGastoFornecedores": 320000.00,
    "gastoMedioFornecedor": 11428.57,
    "top10Gastos": [
      {
        "fornecedorId": "uuid-forn-1",
        "nomeFornecedor": "Construtora ABC Ltda",
        "avaliacao": 4.8,
        "valorTotalGasto": 85000.00,
        "quantidadeOrcamentos": 15,
        "ticketMedio": 5666.67,
        "ultimoOrcamento": "2024-07-15T10:00:00Z",
        "percentual": 26.56
      }
    ],
    "fornecedorMaiorGasto": "Construtora ABC Ltda",
    "valorMaiorGasto": 85000.00
  },
  "estatisticasGerais": {
    "totalFornecedores": 28,
    "fornecedoresAtivos": 25,
    "fornecedoresInativos": 3,
    "avaliacaoMediaGeral": 4.2,
    "tempoMedioContrato": 456
  },
  "ultimaAtualizacao": "2024-08-03T14:30:00Z"
}
```

---

## 6. Fluxo de Caixa Detalhado

### Endpoint
```http
GET /dashboard/fluxo-caixa
```

### Permissões
- Requer: `PermissaoFinanceiroLer`

### Payload de Resposta
```json
{
  "totalEntradas": 450000.00,
  "totalSaidas": 380000.00,
  "saldoAtual": 70000.00,
  "fluxoPorPeriodo": [
    {
      "periodo": "2024-01-01T00:00:00Z",
      "entradas": 75000.00,
      "saidas": 65000.00,
      "saldoLiquido": 10000.00
    },
    {
      "periodo": "2024-02-01T00:00:00Z",
      "entradas": 80000.00,
      "saidas": 70000.00,
      "saldoLiquido": 10000.00
    },
    {
      "periodo": "2024-03-01T00:00:00Z",
      "entradas": 85000.00,
      "saidas": 75000.00,
      "saldoLiquido": 10000.00
    }
  ],
  "tendenciaMensal": "crescente"
}
```

---

## 7. Informações de Cache

### Endpoint
```http
GET /dashboard/cache-info
```

### Permissões
- Público (sem autenticação requerida)

### Payload de Resposta
```json
{
  "ultimaAtualizacao": "2024-08-03T14:30:00Z",
  "ttlRecomendado": 300,
  "secoesDisponiveis": [
    "financeiro",
    "obras", 
    "funcionarios",
    "fornecedores"
  ],
  "versao": "1.0"
}
```

---

## Códigos de Status HTTP

### Sucesso
- `200 OK` - Dados retornados com sucesso

### Erros do Cliente
- `400 Bad Request` - Parâmetros inválidos
- `401 Unauthorized` - Token de autenticação ausente ou inválido
- `403 Forbidden` - Usuário não tem permissão para acessar o recurso
- `404 Not Found` - Seção de dashboard não encontrada

### Erros do Servidor
- `500 Internal Server Error` - Erro interno do servidor

---

## Exemplos de Uso

### Frontend React
```javascript
// Obter dashboard completo
const response = await fetch('/dashboard', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
const dashboard = await response.json();

// Obter apenas seção financeira
const financeiroResponse = await fetch('/dashboard/financeiro?dataInicio=2024-01-01&dataFim=2024-12-31', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
const financeiro = await financeiroResponse.json();
```

### Cache no Frontend
```javascript
// Obter informações de cache
const cacheInfo = await fetch('/dashboard/cache-info').then(r => r.json());

// Implementar cache baseado no TTL recomendado
const CACHE_TTL = cacheInfo.ttlRecomendado * 1000; // Converter para ms
```

---

## Considerações de Performance

1. **Cache**: Recomenda-se implementar cache no frontend com TTL de 5 minutos
2. **Paginação**: Alguns endpoints podem retornar muitos dados - considere filtros
3. **Período**: Limitar consultas a períodos razoáveis (máximo 1 ano)
4. **Seções**: Use o parâmetro `secoes` para obter apenas dados necessários

---

## Monitoramento e Logs

Todos os endpoints geram logs estruturados com:
- Tempo de resposta
- Usuário que fez a requisição
- Parâmetros utilizados
- Erros encontrados

Os logs podem ser utilizados para:
- Análise de performance
- Auditoria de acesso
- Detecção de problemas

### Logs de Exemplo
```json
{
  "time": "2025-08-03T19:05:12.803780209-03:00",
  "level": "INFO",
  "msg": "Dashboard service call completed: geral.ObterDashboardCompleto",
  "extra": {
    "component": "dashboard",
    "dashboardSection": "geral",
    "duration": "77.395215ms",
    "serviceMethod": "ObterDashboardCompleto",
    "success": true
  }
}
```

---

## Correções Aplicadas

### Histórico de Fixes (Agosto 2025)

1. **✅ Correção SQL Ambiguidade**: Resolvido erro `column reference "status" is ambiguous` em `ObterDistribuicaoObras`
2. **✅ Correção NULL Scan**: Aplicado `COALESCE` para tratar valores NULL em campos de avaliação
3. **✅ Correção Função EXTRACT**: Substituída sintaxe inválida por cast de data/intervalo
4. **✅ Logging Estruturado**: Implementado sistema de logs detalhado para auditoria e debug

### Status Atual
- ✅ **Funcional**: Todos os endpoints funcionando corretamente
- ✅ **Performance**: Queries otimizadas com tratamento de NULL
- ✅ **Logging**: Sistema completo de auditoria e monitoramento
- ✅ **Testes**: Validado com dados reais do sistema