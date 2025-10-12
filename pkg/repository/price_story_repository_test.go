package repository

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"price-story/pkg/queries"
)

func TestDeletePriceStory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := NewPriceStoryRepository(db)
	id := "ps1"
	mock.ExpectExec(regexp.QuoteMeta(queries.DeletePriceStoryByID)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ok, err := repo.DeletePriceStory(id)
	if err != nil {
		t.Fatalf("DeletePriceStory error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete ok=true")
	}
}
