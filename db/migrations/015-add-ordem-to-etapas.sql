-- Migração 015: Adicionar campo ordem na tabela etapas e atualizar dados existentes
-- Data: 2025-01-15
-- Descrição: Adiciona campo ordem para ordenação sequencial das etapas

-- Adicionar campo ordem na tabela etapas
ALTER TABLE etapas ADD COLUMN ordem INTEGER DEFAULT 0;

-- Atualizar etapas existentes com ordem baseada no nome
-- Ordem definida: Fundações(1), Estrutura(2), Alvenaria(3), Instalações(4), Pintura(5), Acabamentos(6)
UPDATE etapas SET ordem = 1 WHERE LOWER(nome) LIKE '%fundaç%' OR LOWER(nome) LIKE '%fundac%';
UPDATE etapas SET ordem = 2 WHERE LOWER(nome) LIKE '%estrutura%' OR LOWER(nome) LIKE '%concreto%';
UPDATE etapas SET ordem = 3 WHERE LOWER(nome) LIKE '%alvenaria%' OR LOWER(nome) LIKE '%tijolo%' OR LOWER(nome) LIKE '%parede%';
UPDATE etapas SET ordem = 4 WHERE LOWER(nome) LIKE '%instalaç%' OR LOWER(nome) LIKE '%instalac%' OR LOWER(nome) LIKE '%elétric%' OR LOWER(nome) LIKE '%hidráulic%';
UPDATE etapas SET ordem = 5 WHERE LOWER(nome) LIKE '%pintura%' OR LOWER(nome) LIKE '%tinta%';
UPDATE etapas SET ordem = 6 WHERE LOWER(nome) LIKE '%acabamento%' OR LOWER(nome) LIKE '%finaliz%';

-- Para etapas que não se encaixam nos padrões acima, manter ordem 0 ou definir manualmente
UPDATE etapas SET ordem = 7 WHERE ordem = 0;

-- Criar índice para otimizar consultas ordenadas por obra e ordem
CREATE INDEX IF NOT EXISTS idx_etapas_obra_ordem ON etapas(obra_id, ordem);

-- Comentário na tabela
COMMENT ON COLUMN etapas.ordem IS 'Ordem sequencial da etapa para ordenação (1=Fundações, 2=Estrutura, 3=Alvenaria, 4=Instalações, 5=Pintura, 6=Acabamentos)';