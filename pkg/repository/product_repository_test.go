package repository

import (
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"price-story/pkg/models"
	"price-story/pkg/queries"
)

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer func() { _ = db.Close() }()

	repo := NewProductRepository(db)
	id := "p1"
	mock.ExpectExec(regexp.QuoteMeta(queries.DeleteProductByID)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ok, err := repo.DeleteProduct(id)
	if err != nil {
		t.Fatalf("DeleteProduct error: %v", err)
	}
	if !ok {
		t.Fatalf("expected delete ok=true")
	}

	// zero rows affected
	id2 := "p2"
	mock.ExpectExec(regexp.QuoteMeta(queries.DeleteProductByID)).
		WithArgs(id2).
		WillReturnResult(sqlmock.NewResult(0, 0))

	ok, err = repo.DeleteProduct(id2)
	if err != nil {
		t.Fatalf("DeleteProduct error: %v", err)
	}
	if ok {
		t.Fatalf("expected delete ok=false")
	}
}

func TestCreateProduct_WithoutUpdatedUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer func() { _ = db.Close() }()

	repo := NewProductRepository(db)
	p := models.Product{
		ProductName:        "Widget",
		ProductUrl:         "https://example.com/widget",
		ProductImageUrl:    "https://example.com/widget.png",
		ProductDescription: "desc",
		CreatedByUser:      models.User{ID: "u1", Email: "u1@example.com", PasswordHash: "h"},
	}

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "product_name", "product_url", "product_image_url", "product_description",
		"created_datetime", "updated_datetime",
		"uc.id", "uc.email", "uc.password_hash",
		"uu.id", "uu.email", "uu.password_hash",
	}).AddRow(
		"p123", p.ProductName, p.ProductUrl, p.ProductImageUrl, p.ProductDescription,
		now, now,
		p.CreatedByUser.ID, p.CreatedByUser.Email, p.CreatedByUser.PasswordHash,
		nil, nil, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(queries.InsertProductWithoutUpdatedCTE)).
		WithArgs(p.ProductName, p.ProductUrl, p.ProductImageUrl, p.ProductDescription, p.CreatedByUser.ID).
		WillReturnRows(rows)

	// select updated by user id returns NULL
	mock.ExpectQuery(regexp.QuoteMeta(queries.SelectProductUpdatedByUserID)).
		WithArgs("p123").
		WillReturnRows(sqlmock.NewRows([]string{"updated_by_user_id"}).AddRow(nil))

	got, err := repo.CreateProduct(p)
	if err != nil {
		t.Fatalf("CreateProduct error: %v", err)
	}
	if got.ID != "p123" || got.UpdatedByUser != nil {
		t.Fatalf("unexpected created product: %+v", got)
	}
}
