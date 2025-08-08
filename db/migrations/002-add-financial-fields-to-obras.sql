-- Migração para adicionar campos financeiros na tabela obras
-- Data: 2025-08-08
-- Descrição: Adiciona campos para controle financeiro das obras

-- Adicionar campos financeiros na tabela obras
ALTER TABLE obras 
ADD COLUMN IF NOT EXISTS valor_contrato_total NUMERIC(15, 2) DEFAULT 0.00 NOT NULL,
ADD COLUMN IF NOT EXISTS valor_recebido NUMERIC(15, 2) DEFAULT 0.00 NOT NULL,
ADD COLUMN IF NOT EXISTS tipo_cobranca VARCHAR(20) DEFAULT 'VISTA' NOT NULL 
    CHECK (tipo_cobranca IN ('VISTA', 'PARCELADO', 'ETAPAS')),
ADD COLUMN IF NOT EXISTS data_assinatura_contrato DATE DEFAULT NULL;

-- Criar tabela para cronograma de recebimentos
CREATE TABLE IF NOT EXISTS cronograma_recebimentos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    obra_id UUID NOT NULL REFERENCES obras(id) ON DELETE CASCADE,
    numero_etapa INT NOT NULL,
    descricao_etapa VARCHAR(500) NOT NULL,
    valor_previsto NUMERIC(15, 2) NOT NULL CHECK (valor_previsto > 0),
    data_vencimento DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE' 
        CHECK (status IN ('PENDENTE', 'RECEBIDO', 'VENCIDO', 'PARCIAL')),
    data_recebimento TIMESTAMPTZ DEFAULT NULL,
    valor_recebido NUMERIC(15, 2) DEFAULT 0.00 NOT NULL CHECK (valor_recebido >= 0),
    observacoes_recebimento TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    UNIQUE (obra_id, numero_etapa), -- Uma obra não pode ter etapas duplicadas
    CHECK (valor_recebido <= valor_previsto) -- Valor recebido não pode exceder previsto
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_cronograma_recebimentos_obra_id ON cronograma_recebimentos(obra_id);
CREATE INDEX IF NOT EXISTS idx_cronograma_recebimentos_data_vencimento ON cronograma_recebimentos(data_vencimento);
CREATE INDEX IF NOT EXISTS idx_cronograma_recebimentos_status ON cronograma_recebimentos(status);

-- Comentários para documentação
COMMENT ON COLUMN obras.valor_contrato_total IS 'Valor total do contrato da obra';
COMMENT ON COLUMN obras.valor_recebido IS 'Valor total já recebido da obra';
COMMENT ON COLUMN obras.tipo_cobranca IS 'Tipo de cobrança: VISTA, PARCELADO ou ETAPAS';
COMMENT ON COLUMN obras.data_assinatura_contrato IS 'Data de assinatura do contrato';

COMMENT ON TABLE cronograma_recebimentos IS 'Cronograma de recebimentos por etapa das obras';
COMMENT ON COLUMN cronograma_recebimentos.numero_etapa IS 'Número sequencial da etapa de recebimento';
COMMENT ON COLUMN cronograma_recebimentos.valor_previsto IS 'Valor previsto para recebimento nesta etapa';
COMMENT ON COLUMN cronograma_recebimentos.valor_recebido IS 'Valor já recebido nesta etapa';

-- Criar tabela para contas a receber
CREATE TABLE IF NOT EXISTS contas_receber (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    obra_id UUID REFERENCES obras(id) ON DELETE SET NULL, -- Pode ser NULL para contas não relacionadas a obras
    cronograma_recebimento_id UUID REFERENCES cronograma_recebimentos(id) ON DELETE SET NULL,
    cliente VARCHAR(255) NOT NULL,
    tipo_conta_receber VARCHAR(20) NOT NULL DEFAULT 'OBRA'
        CHECK (tipo_conta_receber IN ('OBRA', 'SERVICO', 'OUTROS')),
    descricao TEXT NOT NULL,
    valor_original NUMERIC(15, 2) NOT NULL CHECK (valor_original > 0),
    valor_recebido NUMERIC(15, 2) DEFAULT 0.00 NOT NULL CHECK (valor_recebido >= 0),
    data_vencimento DATE NOT NULL,
    data_recebimento TIMESTAMPTZ DEFAULT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE'
        CHECK (status IN ('PENDENTE', 'RECEBIDO', 'VENCIDO', 'PARCIAL', 'CANCELADO')),
    forma_pagamento VARCHAR(50) DEFAULT NULL,
    observacoes TEXT DEFAULT NULL,
    numero_documento VARCHAR(100) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CHECK (valor_recebido <= valor_original)
);

-- Índices para performance das contas a receber
CREATE INDEX IF NOT EXISTS idx_contas_receber_obra_id ON contas_receber(obra_id);
CREATE INDEX IF NOT EXISTS idx_contas_receber_cliente ON contas_receber(cliente);
CREATE INDEX IF NOT EXISTS idx_contas_receber_data_vencimento ON contas_receber(data_vencimento);
CREATE INDEX IF NOT EXISTS idx_contas_receber_status ON contas_receber(status);
CREATE INDEX IF NOT EXISTS idx_contas_receber_cronograma ON contas_receber(cronograma_recebimento_id);

-- Comentários para contas a receber
COMMENT ON TABLE contas_receber IS 'Contas a receber do sistema';
COMMENT ON COLUMN contas_receber.tipo_conta_receber IS 'Tipo da conta: OBRA, SERVICO ou OUTROS';
COMMENT ON COLUMN contas_receber.valor_original IS 'Valor original da conta a receber';
COMMENT ON COLUMN contas_receber.valor_recebido IS 'Valor já recebido desta conta';

-- Criar tabela para contas a pagar
CREATE TABLE IF NOT EXISTS contas_pagar (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fornecedor_id UUID REFERENCES fornecedores(id) ON DELETE SET NULL, -- Pode ser NULL para fornecedores não cadastrados
    obra_id UUID REFERENCES obras(id) ON DELETE SET NULL, -- Pode ser NULL para contas não relacionadas a obras
    orcamento_id UUID REFERENCES orcamentos(id) ON DELETE SET NULL, -- Referência ao orçamento que originou
    fornecedor_nome VARCHAR(255) NOT NULL,
    tipo_conta_pagar VARCHAR(20) NOT NULL DEFAULT 'FORNECEDOR'
        CHECK (tipo_conta_pagar IN ('FORNECEDOR', 'SERVICO', 'MATERIAL', 'OUTROS')),
    descricao TEXT NOT NULL,
    valor_original NUMERIC(15, 2) NOT NULL CHECK (valor_original > 0),
    valor_pago NUMERIC(15, 2) DEFAULT 0.00 NOT NULL CHECK (valor_pago >= 0),
    data_vencimento DATE NOT NULL,
    data_pagamento TIMESTAMPTZ DEFAULT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE'
        CHECK (status IN ('PENDENTE', 'PAGO', 'VENCIDO', 'PARCIAL', 'CANCELADO')),
    forma_pagamento VARCHAR(50) DEFAULT NULL,
    observacoes TEXT DEFAULT NULL,
    numero_documento VARCHAR(100) DEFAULT NULL,
    numero_compra_nf VARCHAR(100) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CHECK (valor_pago <= valor_original)
);

-- Criar tabela para parcelas de contas a pagar (para casos de parcelamento)
CREATE TABLE IF NOT EXISTS parcelas_conta_pagar (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conta_pagar_id UUID NOT NULL REFERENCES contas_pagar(id) ON DELETE CASCADE,
    numero_parcela INT NOT NULL CHECK (numero_parcela > 0),
    valor_parcela NUMERIC(15, 2) NOT NULL CHECK (valor_parcela > 0),
    data_vencimento DATE NOT NULL,
    data_pagamento TIMESTAMPTZ DEFAULT NULL,
    valor_pago NUMERIC(15, 2) DEFAULT 0.00 NOT NULL CHECK (valor_pago >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDENTE'
        CHECK (status IN ('PENDENTE', 'PAGO', 'VENCIDO', 'PARCIAL', 'CANCELADO')),
    forma_pagamento VARCHAR(50) DEFAULT NULL,
    observacoes TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CHECK (valor_pago <= valor_parcela),
    UNIQUE (conta_pagar_id, numero_parcela) -- Uma conta não pode ter parcelas duplicadas
);

-- Índices para performance das contas a pagar
CREATE INDEX IF NOT EXISTS idx_contas_pagar_fornecedor_id ON contas_pagar(fornecedor_id);
CREATE INDEX IF NOT EXISTS idx_contas_pagar_obra_id ON contas_pagar(obra_id);
CREATE INDEX IF NOT EXISTS idx_contas_pagar_orcamento_id ON contas_pagar(orcamento_id);
CREATE INDEX IF NOT EXISTS idx_contas_pagar_fornecedor_nome ON contas_pagar(fornecedor_nome);
CREATE INDEX IF NOT EXISTS idx_contas_pagar_data_vencimento ON contas_pagar(data_vencimento);
CREATE INDEX IF NOT EXISTS idx_contas_pagar_status ON contas_pagar(status);

-- Índices para parcelas
CREATE INDEX IF NOT EXISTS idx_parcelas_conta_pagar_conta_id ON parcelas_conta_pagar(conta_pagar_id);
CREATE INDEX IF NOT EXISTS idx_parcelas_conta_pagar_data_vencimento ON parcelas_conta_pagar(data_vencimento);
CREATE INDEX IF NOT EXISTS idx_parcelas_conta_pagar_status ON parcelas_conta_pagar(status);

-- Comentários para contas a pagar
COMMENT ON TABLE contas_pagar IS 'Contas a pagar do sistema (fornecedores, serviços, etc.)';
COMMENT ON COLUMN contas_pagar.tipo_conta_pagar IS 'Tipo da conta: FORNECEDOR, SERVICO, MATERIAL ou OUTROS';
COMMENT ON COLUMN contas_pagar.valor_original IS 'Valor original da conta a pagar';
COMMENT ON COLUMN contas_pagar.valor_pago IS 'Valor já pago desta conta';
COMMENT ON COLUMN contas_pagar.numero_compra_nf IS 'Número da compra ou nota fiscal';

COMMENT ON TABLE parcelas_conta_pagar IS 'Parcelas de contas a pagar parceladas';
COMMENT ON COLUMN parcelas_conta_pagar.numero_parcela IS 'Número sequencial da parcela (1, 2, 3...)';
COMMENT ON COLUMN parcelas_conta_pagar.valor_parcela IS 'Valor desta parcela específica';