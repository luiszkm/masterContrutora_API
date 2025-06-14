-- file: db/init/init.sql

CREATE TABLE IF NOT EXISTS obras (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cliente VARCHAR(255) NOT NULL,
    endereco TEXT NOT NULL,
    data_inicio DATE NOT NULL,
    data_fim DATE,
    status VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS etapas (
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras(id) ON DELETE CASCADE,
    nome VARCHAR(255) NOT NULL,
    data_inicio_prevista DATE,
    data_fim_prevista DATE,
    status VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS alocacoes (
    id UUID PRIMARY KEY,
    obra_id UUID NOT NULL REFERENCES obras(id) ON DELETE CASCADE,
    funcionario_id UUID NOT NULL, -- Futuramente referenciará a tabela de funcionários
    data_inicio_alocacao DATE NOT NULL,
    data_fim_alocacao DATE
);

CREATE TABLE IF NOT EXISTS usuarios (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha_hash TEXT NOT NULL,
    permissoes TEXT[] NOT NULL, -- Array de strings para as permissões
    ativo BOOLEAN NOT NULL DEFAULT TRUE
);

-- Adiciona índices nas colunas de chave estrangeira para otimizar as buscas (JOINs).
CREATE INDEX IF NOT EXISTS idx_etapas_obra_id ON etapas(obra_id);
CREATE INDEX IF NOT EXISTS idx_alocacoes_obra_id ON alocacoes(obra_id);