-- Create price_points table
CREATE TABLE IF NOT EXISTS price_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    price NUMERIC(12,2) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_by_user_id UUID NOT NULL REFERENCES users(id),
    updated_by_user_id UUID REFERENCES users(id),
    created_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW()
);