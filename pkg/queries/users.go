package queries

const ListUsers = `
    SELECT id, email, password_hash
    FROM users
    ORDER BY email ASC
`

const GetUserByID = `
    SELECT id, email, password_hash
    FROM users
    WHERE id = $1
`

const GetUserByEmail = `
    SELECT id, email, password_hash
    FROM users
    WHERE email = $1
`

const InsertUserReturning = `
    INSERT INTO users (email, password_hash)
    VALUES ($1, $2)
    RETURNING id, email, password_hash
`

const UpdateUserReturning = `
    UPDATE users
    SET email = $1, password_hash = $2
    WHERE id = $3
    RETURNING id, email, password_hash
`

const DeleteUserByID = `
    DELETE FROM users
    WHERE id = $1
`
