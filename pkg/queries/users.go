package queries

const ListUsers = `
    SELECT id, email, password_hash, preferred_currency, decimal_places, thousand_separator, currency_symbol_placement
    FROM users
    ORDER BY email ASC
`

const GetUserByID = `
    SELECT id, email, password_hash, preferred_currency, decimal_places, thousand_separator, currency_symbol_placement
    FROM users
    WHERE id = $1
`

const GetUserByEmail = `
    SELECT id, email, password_hash, preferred_currency, decimal_places, thousand_separator, currency_symbol_placement
    FROM users
    WHERE email = $1
`

const InsertUserReturning = `
    INSERT INTO users (email, password_hash, preferred_currency, decimal_places, thousand_separator, currency_symbol_placement)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, email, password_hash, preferred_currency, decimal_places, thousand_separator, currency_symbol_placement
`

const UpdateUserReturning = `
    UPDATE users
    SET email = $1, password_hash = $2,
        preferred_currency = $3, decimal_places = $4, thousand_separator = $5, currency_symbol_placement = $6
    WHERE id = $7
    RETURNING id, email, password_hash, preferred_currency, decimal_places, thousand_separator, currency_symbol_placement
`

const DeleteUserByID = `
    DELETE FROM users
    WHERE id = $1
`
