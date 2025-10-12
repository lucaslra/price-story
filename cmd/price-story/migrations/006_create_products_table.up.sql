-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_name TEXT NOT NULL,
    product_url TEXT NOT NULL,
    created_by_user_id UUID NOT NULL REFERENCES users(id),
    updated_by_user_id UUID REFERENCES users(id),
    created_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW()
);