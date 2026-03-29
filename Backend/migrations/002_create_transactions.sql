-- Migration: 002_create_transactions.sql
-- Description: Create transactions table for indexed blockchain events

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    tx_hash VARCHAR(66) UNIQUE NOT NULL,
    event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('deposit', 'withdrawal', 'transfer')),
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    amount VARCHAR(78) NOT NULL,
    block_number BIGINT NOT NULL,
    block_timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transactions_from ON transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions(block_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_block ON transactions(block_number);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(event_type);
