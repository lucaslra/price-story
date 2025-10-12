-- Add core profile attributes to users table with validation
ALTER TABLE users
    ADD COLUMN preferred_currency CHAR(3) NOT NULL DEFAULT 'USD',
    ADD COLUMN decimal_places INTEGER NOT NULL DEFAULT 2,
    ADD COLUMN thousand_separator TEXT NOT NULL DEFAULT ',',
    ADD COLUMN currency_symbol_placement TEXT NOT NULL DEFAULT 'before';

-- Constraints for currency code format and number formatting rules
ALTER TABLE users
    ADD CONSTRAINT users_preferred_currency_format CHECK (preferred_currency ~ '^[A-Z]{3}$'),
    ADD CONSTRAINT users_decimal_places_range CHECK (decimal_places BETWEEN 0 AND 4),
    ADD CONSTRAINT users_thousand_separator_valid CHECK (thousand_separator IN (',', '.', ' ')),
    ADD CONSTRAINT users_symbol_placement_valid CHECK (currency_symbol_placement IN ('before', 'after')),
    ADD CONSTRAINT users_jpy_decimal_places CHECK (preferred_currency <> 'JPY' OR decimal_places = 0);