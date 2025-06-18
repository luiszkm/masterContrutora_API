-- file: db/init/01-init.sql
-- Script de inicialização completo e corrigido para o banco de dados.
-- Todas as tabelas necessárias para o sistema estão definidas aqui.

-- Tabelas sem dependências externas diretas
CREATE TABLE IF NOT EXISTS usuarios (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha_hash TEXT NOT NULL,
    permissoes TEXT[] NOT NULL,
    ativo BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS obras (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cliente VARCHAR(255) NOT NULL,
    endereco TEXT NOT NULL,
    data_inicio DATE NOT NULL,
    data_fim DATE,
    status VARCHAR(50) NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL -- Corrigido: vírgula final removida
);

-- Tabela de funcionários V2 (a versão antiga é descartada)
-- O DROP garante que, ao reiniciar o ambiente, a estrutura mais nova seja usada.
DROP TABLE IF EXISTS funcionarios CASCADE;
CREATE TABLE IF NOT EXISTS funcionarios (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) NOT NULL UNIQUE,
    telefone VARCHAR(20) NOT NULL DEFAULT '',
    cargo VARCHAR(100) NOT NULL,
    departamento VARCHAR(100) NOT NULL DEFAULT '',
    data_contratacao DATE NOT NULL,
    valor_diaria NUMERIC(10, 2) NOT NULL,
    chave_pix VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL,
    desligamento_data TIMESTAMPTZ,
    motivo_desligamento TEXT NOT NULL DEFAULT '',
    observacoes TEXT NOT NULL DEFAULT '',
    avaliação_desempenho TEXT NOT NULL DEFAULT '',
    avatar_url VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fornecedores (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cnpj VARCHAR(18) NOT NULL UNIQUE,
    categoria VARCHAR(100) NOT NULL,
    contato VARCHAR(255),
    email VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS materiais (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    unidade_de_medida VARCHAR(20) NOT NULL,
    categoria VARCHAR(100) NOT NULL
);


-- Tabelas com dependências (chaves estrangeiras)

CREATE TABLE IF NOT EXISTS etapas (
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    data_inicio_prevista DATE,
    data_fim_prevista DATE,
    status VARCHAR(50) NOT NULL
);

-- Recriando a tabela 'alocacoes' para garantir que ela exista e referencie a nova tabela de funcionários.
CREATE TABLE IF NOT EXISTS alocacoes (
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras(id) ON DELETE CASCADE,
    funcionario_id UUID NOT NULL REFERENCES funcionarios(id) ON DELETE CASCADE,
    data_inicio_alocacao DATE NOT NULL,
    data_fim_alocacao DATE
);

CREATE TABLE IF NOT EXISTS apontamentos_quinzenais (
    id UUID PRIMARY KEY,
    funcionario_id UUID NOT NULL REFERENCES funcionarios(id),
    obra_id UUID NOT NULL REFERENCES obras(id),
    periodo_inicio DATE NOT NULL,
    periodo_fim DATE NOT NULL,
    dias_trabalhados INT NOT NULL DEFAULT 0,
    adicionais NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    descontos NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    adiantamentos NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    valor_total_calculado NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(funcionario_id, periodo_inicio, periodo_fim) -- Garante que não haja apontamentos duplicados.
);

CREATE TABLE IF NOT EXISTS orcamentos (
    id UUID PRIMARY KEY,
    numero VARCHAR(50) NOT NULL UNIQUE,
    etapa_id UUID NOT NULL REFERENCES etapas(id),
    fornecedor_id UUID NOT NULL REFERENCES fornecedores(id),
    valor_total NUMERIC(15, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    data_emissao TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS orcamento_itens (
    id UUID PRIMARY KEY,
    orcamento_id UUID NOT NULL REFERENCES orcamentos(id) ON DELETE CASCADE,
    material_id UUID NOT NULL REFERENCES materiais(id),
    quantidade NUMERIC(10, 2) NOT NULL,
    valor_unitario NUMERIC(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS registros_pagamento (
    id UUID PRIMARY KEY,
    funcionario_id UUID NOT NULL REFERENCES funcionarios(id),
    obra_id UUID NOT NULL REFERENCES obras(id),
    periodo_referencia VARCHAR(100) NOT NULL,
    valor_calculado NUMERIC(10, 2) NOT NULL,
    data_de_efetivacao TIMESTAMPTZ NOT NULL,
    conta_bancaria_id UUID NOT NULL
);

-- Adiciona índices para otimizar as buscas (JOINs).
CREATE INDEX IF NOT EXISTS idx_etapas_obra_id ON etapas(obra_id);
CREATE INDEX IF NOT EXISTS idx_orcamentos_etapa_id ON orcamentos(etapa_id);
CREATE INDEX IF NOT EXISTS idx_alocacoes_obra_id ON alocacoes(obra_id);
CREATE INDEX IF NOT EXISTS idx_alocacoes_funcionario_id ON alocacoes(funcionario_id);
CREATE INDEX IF NOT EXISTS idx_materiais_categoria ON materiais(categoria);
CREATE INDEX IF NOT EXISTS idx_apontamentos_funcionario_id ON apontamentos_quinzenais(funcionario_id);
CREATE INDEX IF NOT EXISTS idx_pagamentos_funcionario_id ON registros_pagamento(funcionario_id);
CREATE INDEX IF NOT EXISTS idx_pagamentos_obra_id ON registros_pagamento(obra_id);