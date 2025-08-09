-- Migration to add soft delete support to produtos and orcamentos tables

-- Add deleted_at column to produtos table
ALTER TABLE produtos ADD COLUMN deleted_at TIMESTAMPTZ DEFAULT NULL;

-- Add deleted_at column to orcamentos table  
ALTER TABLE orcamentos ADD COLUMN deleted_at TIMESTAMPTZ DEFAULT NULL;

-- Create indexes for better performance on soft delete queries
CREATE INDEX IF NOT EXISTS idx_produtos_deleted_at ON produtos(deleted_at);
CREATE INDEX IF NOT EXISTS idx_orcamentos_deleted_at ON orcamentos(deleted_at);