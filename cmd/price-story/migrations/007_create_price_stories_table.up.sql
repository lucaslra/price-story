-- Create price_stories table
CREATE TABLE IF NOT EXISTS price_stories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_by_user_id UUID NOT NULL REFERENCES users(id),
    updated_by_user_id UUID REFERENCES users(id),
    created_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW()
);