-- Migration: Remove valor_diaria column from funcionarios table
-- Justification: Diaria values should only exist in apontamentos, not in funcionarios
-- This allows for flexible daily rates per work period/project

-- Remove the valor_diaria column from funcionarios table
ALTER TABLE funcionarios DROP COLUMN IF EXISTS valor_diaria;