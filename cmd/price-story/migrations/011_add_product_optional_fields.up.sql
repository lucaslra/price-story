-- Add optional fields to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS product_image_url TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS product_description TEXT;