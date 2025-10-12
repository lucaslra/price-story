package repository

import (
	"database/sql"
	"price-story/pkg/models"
	"price-story/pkg/queries"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ListUsers retrieves all users
func (r *UserRepository) ListUsers() ([]models.User, error) {
	rows, err := r.db.Query(queries.ListUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUser retrieves a user by ID
func (r *UserRepository) GetUser(id string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(queries.GetUserByID, id).Scan(&u.ID, &u.Email, &u.PasswordHash)
	return u, err
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(queries.GetUserByEmail, email).Scan(&u.ID, &u.Email, &u.PasswordHash)
	return u, err
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(u models.User) (models.User, error) {
	// Return created record
	err := r.db.QueryRow(queries.InsertUserReturning, u.Email, u.PasswordHash).Scan(&u.ID, &u.Email, &u.PasswordHash)
	return u, err
}

// EnsureUserByEmail returns the user if exists, otherwise creates it
func (r *UserRepository) EnsureUserByEmail(email, passwordHash string) (models.User, error) {
	u, err := r.GetUserByEmail(email)
	if err == nil {
		return u, nil
	}
	if err == sql.ErrNoRows {
		return r.CreateUser(models.User{Email: email, PasswordHash: passwordHash})
	}
	return models.User{}, err
}

// UpdateUser updates an existing user by ID
func (r *UserRepository) UpdateUser(id string, u models.User) (models.User, error) {
	u.ID = id
	err := r.db.QueryRow(queries.UpdateUserReturning, u.Email, u.PasswordHash, id).Scan(&u.ID, &u.Email, &u.PasswordHash)
	return u, err
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id string) (bool, error) {
	res, err := r.db.Exec(queries.DeleteUserByID, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
