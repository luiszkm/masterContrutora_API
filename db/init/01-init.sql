-- file: db/init/01-init.sql
-- Script de inicialização V4.0 Final
-- Define a estrutura completa e correta do banco de dados, com todas as refatorações aplicadas.

-- Tabelas sem dependências externas (entidades principais)
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
    nome VARCHAR(255) NOT NULL UNIQUE, -- <<--- RESTRIÇÃO UNIQUE ADICIONADA AQUI
    cliente VARCHAR(255) NOT NULL,
    endereco TEXT NOT NULL,
    descricao TEXT,
    data_inicio DATE NOT NULL,
    data_fim DATE,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS funcionarios (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) NOT NULL UNIQUE,
    telefone VARCHAR(20) NOT NULL DEFAULT '',
    cargo VARCHAR(100) NOT NULL,
    departamento VARCHAR(100) NOT NULL DEFAULT '',
    email VARCHAR(100) NOT NULL DEFAULT '',
    data_contratacao DATE NOT NULL,
    valor_diaria NUMERIC(10, 2) NOT NULL,
    chave_pix VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL,
    desligamento_data TIMESTAMPTZ,
    motivo_desligamento TEXT NOT NULL DEFAULT '',
    observacoes TEXT NOT NULL DEFAULT '',
    avaliacao_desempenho TEXT NOT NULL DEFAULT '',
    avatar_url VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categorias (
    id UUID PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS fornecedores (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cnpj VARCHAR(18) NOT NULL UNIQUE,
    contato VARCHAR(255),
    email VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    endereco TEXT,
    avaliacao NUMERIC(3, 1),
    observacoes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS produtos (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL UNIQUE, -- Nome do produto agora é único
    descricao TEXT,
    unidade_de_medida VARCHAR(50) NOT NULL,
    categoria VARCHAR(100), -- Categoria como string, pode ser vinculada a tabela 'categorias' no futuro se necessário
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS etapas_padrao (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL UNIQUE,
    descricao TEXT,
    ordem INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Tabelas com dependências (chaves estrangeiras)
CREATE TABLE IF NOT EXISTS etapas (
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    data_inicio_prevista DATE,
    data_fim_prevista DATE,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS alocacoes (
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras(id) ON DELETE CASCADE,
    funcionario_id UUID NOT NULL REFERENCES funcionarios(id) ON DELETE RESTRICT, -- Evita deletar funcionário alocado
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
    diaria NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    descontos NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    adiantamentos NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    valor_total_calculado NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(funcionario_id, periodo_inicio, periodo_fim)
);

CREATE TABLE IF NOT EXISTS orcamentos (
    id UUID PRIMARY KEY,
    numero VARCHAR(50) NOT NULL UNIQUE,
    etapa_id UUID NOT NULL REFERENCES etapas(id),
    fornecedor_id UUID NOT NULL REFERENCES fornecedores(id),
    valor_total NUMERIC(15, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    data_emissao TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    data_aprovacao TIMESTAMPTZ,
    observacoes TEXT,
    condicoes_pagamento VARCHAR(255),
    criado_por_usuario_id UUID REFERENCES usuarios(id), -- Permitir nulo caso o usuário seja deletado
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orcamento_itens (
    id UUID PRIMARY KEY,
    orcamento_id UUID NOT NULL REFERENCES orcamentos(id) ON DELETE CASCADE,
    produto_id UUID NOT NULL REFERENCES produtos(id), -- Corrigido para referenciar 'produtos'
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

CREATE TABLE IF NOT EXISTS fornecedor_categorias (
    fornecedor_id UUID NOT NULL REFERENCES fornecedores(id) ON DELETE CASCADE,
    categoria_id UUID NOT NULL REFERENCES categorias(id) ON DELETE CASCADE,
    PRIMARY KEY (fornecedor_id, categoria_id)
);


-- Adiciona índices para otimizar as buscas (JOINs).
CREATE INDEX IF NOT EXISTS idx_etapas_obra_id ON etapas(obra_id);
CREATE INDEX IF NOT EXISTS idx_orcamentos_etapa_id ON orcamentos(etapa_id);
CREATE INDEX IF NOT EXISTS idx_orcamentos_fornecedor_id ON orcamentos(fornecedor_id);
CREATE INDEX IF NOT EXISTS idx_alocacoes_obra_id ON alocacoes(obra_id);
CREATE INDEX IF NOT EXISTS idx_alocacoes_funcionario_id ON alocacoes(funcionario_id);
CREATE INDEX IF NOT EXISTS idx_produtos_categoria ON produtos(categoria);
CREATE INDEX IF NOT EXISTS idx_apontamentos_funcionario_id ON apontamentos_quinzenais(funcionario_id);
CREATE INDEX IF NOT EXISTS idx_pagamentos_funcionario_id ON registros_pagamento(funcionario_id);
CREATE INDEX IF NOT EXISTS idx_fornecedor_categorias_categoria_id ON fornecedor_categorias(categoria_id);