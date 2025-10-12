package repository

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"price-story/pkg/queries"
)

func TestDeletePricePoint(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := NewPricePointRepository(db)
	id := "pp1"
	mock.ExpectExec(regexp.QuoteMeta(queries.DeletePricePointByID)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ok, err := repo.DeletePricePoint(id)
	if err != nil {
		t.Fatalf("DeletePricePoint error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete ok=true")
	}
}
