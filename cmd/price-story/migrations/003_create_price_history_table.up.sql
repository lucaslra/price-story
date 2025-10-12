-- Create legacy price_history table (for symmetric migrations)
-- Note: This table is dropped later by migration 010
CREATE TABLE IF NOT EXISTS price_history (
    id UUID PRIMARY KEY,
    product_id UUID,
    price NUMERIC(12,2) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
);