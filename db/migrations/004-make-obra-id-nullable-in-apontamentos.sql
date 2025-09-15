-- Migration: Make obra_id nullable in apontamentos_quinzenais table
-- This allows creating apontamentos without associating them to a specific obra

ALTER TABLE apontamentos_quinzenais
ALTER COLUMN obra_id DROP NOT NULL;

-- Add a comment to document this change
COMMENT ON COLUMN apontamentos_quinzenais.obra_id IS 'ID da obra associada ao apontamento (opcional - permite apontamentos gerais não vinculados a projetos específicos)';