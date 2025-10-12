-- Rollback core profile attributes from users table
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_jpy_decimal_places,
    DROP CONSTRAINT IF EXISTS users_symbol_placement_valid,
    DROP CONSTRAINT IF EXISTS users_thousand_separator_valid,
    DROP CONSTRAINT IF EXISTS users_decimal_places_range,
    DROP CONSTRAINT IF EXISTS users_preferred_currency_format,
    DROP COLUMN IF EXISTS currency_symbol_placement,
    DROP COLUMN IF EXISTS thousand_separator,
    DROP COLUMN IF EXISTS decimal_places,
    DROP COLUMN IF EXISTS preferred_currency;