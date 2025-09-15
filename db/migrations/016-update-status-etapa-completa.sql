-- Migração 016: Atualizar status das etapas de 'Concluída' para 'Completa'
-- Data: 2025-01-15
-- Descrição: Padroniza status das etapas para usar 'Completa' ao invés de 'Concluída'

-- Atualizar todas as etapas com status 'Concluída' para 'Completa'
UPDATE etapas
SET status = 'Completa'
WHERE status = 'Concluída';

-- Verificar quantos registros foram atualizados (comentário para log)
-- SELECT COUNT(*) FROM etapas WHERE status = 'Completa';

-- Comentário explicativo
COMMENT ON COLUMN etapas.status IS 'Status da etapa: Pendente, Em Andamento, Completa';