package repository

import (
    "database/sql"
    "regexp"
    "testing"

    sqlmock "github.com/DATA-DOG/go-sqlmock"
    "price-story/pkg/models"
    "price-story/pkg/queries"
)

func TestEnsureUserByEmail_ReturnsExisting(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("sqlmock.New error: %v", err)
    }
    defer db.Close()

    repo := NewUserRepository(db)
    email := "test@example.com"
    mock.ExpectQuery(regexp.QuoteMeta(queries.GetUserByEmail)).
        WithArgs(email).
        WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).
            AddRow("u1", email, "hash"))

    u, err := repo.EnsureUserByEmail(email, "irrelevant")
    if err != nil {
        t.Fatalf("EnsureUserByEmail error: %v", err)
    }
    if u.ID != "u1" || u.Email != email || u.PasswordHash != "hash" {
        t.Fatalf("unexpected user: %+v", u)
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}

func TestEnsureUserByEmail_CreatesWhenMissing(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("sqlmock.New error: %v", err)
    }
    defer db.Close()

    repo := NewUserRepository(db)
    email := "new@example.com"
    pass := "hash"
    mock.ExpectQuery(regexp.QuoteMeta(queries.GetUserByEmail)).
        WithArgs(email).
        WillReturnError(sql.ErrNoRows)

    mock.ExpectQuery(regexp.QuoteMeta(queries.InsertUserReturning)).
        WithArgs(email, pass).
        WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).
            AddRow("u2", email, pass))

    u, err := repo.EnsureUserByEmail(email, pass)
    if err != nil {
        t.Fatalf("EnsureUserByEmail error: %v", err)
    }
    if u.ID != "u2" || u.Email != email || u.PasswordHash != pass {
        t.Fatalf("unexpected user: %+v", u)
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Fatalf("unmet expectations: %v", err)
    }
}

func TestUpdateUser_ReturnsUpdated(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("sqlmock.New error: %v", err)
    }
    defer db.Close()

    repo := NewUserRepository(db)
    id := "u3"
    u := models.User{Email: "e@x.com", PasswordHash: "h"}

    mock.ExpectQuery(regexp.QuoteMeta(queries.UpdateUserReturning)).
        WithArgs(u.Email, u.PasswordHash, id).
        WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}).
            AddRow(id, u.Email, u.PasswordHash))

    got, err := repo.UpdateUser(id, u)
    if err != nil {
        t.Fatalf("UpdateUser error: %v", err)
    }
    if got.ID != id || got.Email != u.Email || got.PasswordHash != u.PasswordHash {
        t.Fatalf("unexpected updated user: %+v", got)
    }
}