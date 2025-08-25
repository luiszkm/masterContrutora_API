-- Migração 014: Adicionar tipo FUNCIONARIO para contas a pagar
-- Data: 2025-08-24
-- Descrição: Adiciona o tipo "FUNCIONARIO" às contas a pagar para suportar 
--            pagamentos de funcionários criados automaticamente quando 
--            apontamentos são aprovados.

-- Remove a constraint existente
ALTER TABLE contas_pagar 
DROP CONSTRAINT contas_pagar_tipo_conta_pagar_check;

-- Adiciona nova constraint incluindo FUNCIONARIO
ALTER TABLE contas_pagar 
ADD CONSTRAINT contas_pagar_tipo_conta_pagar_check 
CHECK (tipo_conta_pagar::text = ANY (ARRAY[
    'FORNECEDOR'::character varying, 
    'SERVICO'::character varying, 
    'MATERIAL'::character varying, 
    'FUNCIONARIO'::character varying, 
    'OUTROS'::character varying
]::text[]));

-- Comentário para documentar a alteração
COMMENT ON CONSTRAINT contas_pagar_tipo_conta_pagar_check ON contas_pagar 
IS 'Constraint atualizada para incluir FUNCIONARIO - suporte a pagamentos de funcionários via apontamentos aprovados';