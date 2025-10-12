-- Remove optional fields from products table
ALTER TABLE products DROP COLUMN IF EXISTS product_image_url;
ALTER TABLE products DROP COLUMN IF EXISTS product_description;