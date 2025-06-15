-- file: db/init/init.sql

CREATE TABLE
IF NOT EXISTS obras
(
    id UUID PRIMARY KEY,
    nome VARCHAR
(255) NOT NULL,
    cliente VARCHAR
(255) NOT NULL,
    endereco TEXT NOT NULL,
    data_inicio DATE NOT NULL,
    data_fim DATE,
    status VARCHAR
(50) NOT NULL
        deleted_at  TIMESTAMPTZ DEFAULT NULL, -- Para soft delete

);

CREATE TABLE
IF NOT EXISTS etapas
(
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras
(id) ON
DELETE CASCADE,
    nome VARCHAR(255)
NOT NULL,
    data_inicio_prevista DATE,
    data_fim_prevista DATE,
    status VARCHAR
(50) NOT NULL
);

CREATE TABLE
IF NOT EXISTS alocacoes
(
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras
(id) ON
DELETE CASCADE,
    funcionario_id UUID
NOT NULL, -- Futuramente referenciará a tabela de funcionários
    data_inicio_alocacao DATE NOT NULL,
    data_fim_alocacao DATE
);

CREATE TABLE
IF NOT EXISTS usuarios
(
    id UUID PRIMARY KEY,
    nome VARCHAR
(255) NOT NULL,
    email VARCHAR
(255) NOT NULL UNIQUE,
    senha_hash TEXT NOT NULL,
    permissoes TEXT[] NOT NULL, -- Array de strings para as permissões
    ativo BOOLEAN NOT NULL DEFAULT TRUE
);

-- file: migrations/000005_create_funcionarios_table.up.sql
CREATE TABLE
IF NOT EXISTS funcionarios
(
    id UUID PRIMARY KEY,
    nome VARCHAR
(255) NOT NULL,
    cpf VARCHAR
(14) NOT NULL UNIQUE, -- Ex: 123.456.789-00
    cargo VARCHAR
(100) NOT NULL,
    data_contratacao DATE NOT NULL,
    salario NUMERIC
(10, 2) NOT NULL,
    diaria NUMERIC
(10, 2) NOT NULL,
    status VARCHAR
(50) NOT NULL,
 deleted_at TIMESTAMPT DEFAULT NULL
);

CREATE TABLE
IF NOT EXISTS fornecedores
(
    id UUID PRIMARY KEY,
    nome VARCHAR
(255) NOT NULL,
    cnpj VARCHAR
(18) NOT NULL UNIQUE, -- Ex: 00.000.000/0001-00
    categoria VARCHAR
(100) NOT NULL,
    contato VARCHAR
(255),
    email VARCHAR
(255),
    status VARCHAR
(50) NOT NULL
);

CREATE TABLE
IF NOT EXISTS materiais
(
    id UUID PRIMARY KEY,
    nome VARCHAR
(255) NOT NULL,
    descricao TEXT,
    unidade_de_medida VARCHAR
(20) NOT NULL,
    categoria VARCHAR
(100) NOT NULL
);
-- Tabela principal do orçamento
CREATE TABLE
IF NOT EXISTS orcamentos
(
    id UUID PRIMARY KEY,
    numero VARCHAR
(50) NOT NULL UNIQUE,
    etapa_id UUID NOT NULL REFERENCES etapas
(id),
    fornecedor_id UUID NOT NULL REFERENCES fornecedores
(id),
    valor_total NUMERIC
(15, 2) NOT NULL,
    status VARCHAR
(50) NOT NULL,
    data_emissao TIMESTAMPTZ NOT NULL
);

-- Tabela para os itens do orçamento
CREATE TABLE
IF NOT EXISTS orcamento_itens
(
    id UUID PRIMARY KEY,
    orcamento_id UUID NOT NULL REFERENCES orcamentos
(id) ON
DELETE CASCADE,
    material_id UUID
NOT NULL REFERENCES materiais
(id),
    quantidade NUMERIC
(10, 2) NOT NULL,
    valor_unitario NUMERIC
(10, 2) NOT NULL
);


-- Adiciona índices nas colunas de chave estrangeira para otimizar as buscas (JOINs).
CREATE INDEX
IF NOT EXISTS idx_etapas_obra_id ON etapas
(obra_id);
CREATE INDEX
IF NOT EXISTS idx_orcamentos_etapa_id ON orcamentos
(etapa_id);

CREATE INDEX
IF NOT EXISTS idx_alocacoes_obra_id ON alocacoes
(obra_id);
CREATE INDEX
IF NOT EXISTS idx_materiais_categoria ON materiais
(categoria);
