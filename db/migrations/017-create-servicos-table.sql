-- Migration: Create servicos table
-- Date: 2025-09-20

CREATE TABLE IF NOT EXISTS servicos (
    id UUID PRIMARY KEY,
    nome VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index for performance on nome field
CREATE INDEX IF NOT EXISTS idx_servicos_nome ON servicos(nome);

-- Insert some initial data
INSERT INTO servicos (id, nome, created_at, updated_at) VALUES
    (gen_random_uuid(), 'Instalação Elétrica', NOW(), NOW()),
    (gen_random_uuid(), 'Instalação Hidráulica', NOW(), NOW()),
    (gen_random_uuid(), 'Pintura', NOW(), NOW()),
    (gen_random_uuid(), 'Reboco', NOW(), NOW()),
    (gen_random_uuid(), 'Alvenaria', NOW(), NOW())
ON CONFLICT (nome) DO NOTHING;