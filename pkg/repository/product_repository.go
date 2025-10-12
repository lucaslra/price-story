package repository

import (
	"database/sql"
	"price-story/pkg/models"
	"price-story/pkg/queries"
)

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// ListProducts retrieves all products with nested created/updated users
func (r *ProductRepository) ListProducts() ([]models.Product, error) {
	rows, err := r.db.Query(queries.ListProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var updID, updEmail, updHash sql.NullString
		if err := rows.Scan(
			&p.ID,
			&p.ProductName,
			&p.ProductUrl,
			&p.ProductImageUrl,
			&p.ProductDescription,
			&p.CreatedDatetime,
			&p.UpdatedDateTime,
			&p.CreatedByUser.ID,
			&p.CreatedByUser.Email,
			&p.CreatedByUser.PasswordHash,
			&updID,
			&updEmail,
			&updHash,
		); err != nil {
			return nil, err
		}
		if updID.Valid {
			p.UpdatedByUser = &models.User{ID: updID.String, Email: updEmail.String, PasswordHash: updHash.String}
		} else {
			p.UpdatedByUser = nil
		}
		products = append(products, p)
	}
	return products, nil
}

// GetProduct retrieves a product by ID with nested created/updated users
func (r *ProductRepository) GetProduct(id string) (models.Product, error) {
	var p models.Product
	var updID, updEmail, updHash sql.NullString
	err := r.db.QueryRow(queries.GetProductByID, id).Scan(
		&p.ID,
		&p.ProductName,
		&p.ProductUrl,
		&p.ProductImageUrl,
		&p.ProductDescription,
		&p.CreatedDatetime,
		&p.UpdatedDateTime,
		&p.CreatedByUser.ID,
		&p.CreatedByUser.Email,
		&p.CreatedByUser.PasswordHash,
		&updID,
		&updEmail,
		&updHash,
	)
	if err != nil {
		return p, err
	}
	if updID.Valid {
		p.UpdatedByUser = &models.User{ID: updID.String, Email: updEmail.String, PasswordHash: updHash.String}
	} else {
		p.UpdatedByUser = nil
	}
	return p, err
}

// CreateProduct creates a new product
func (r *ProductRepository) CreateProduct(p models.Product) (models.Product, error) {
	var updatedByID interface{}
	if p.UpdatedByUser != nil {
		updatedByID = p.UpdatedByUser.ID
	} else {
		updatedByID = nil
	}
	var err error
	if updatedByID == nil {
		err = r.db.QueryRow(queries.InsertProductWithoutUpdatedCTE,
			p.ProductName,
			p.ProductUrl,
			p.ProductImageUrl,
			p.ProductDescription,
			p.CreatedByUser.ID,
		).Scan(
			&p.ID,
			&p.ProductName,
			&p.ProductUrl,
			&p.ProductImageUrl,
			&p.ProductDescription,
			&p.CreatedDatetime,
			&p.UpdatedDateTime,
			&p.CreatedByUser.ID,
			&p.CreatedByUser.Email,
			&p.CreatedByUser.PasswordHash,
			new(sql.NullString), // uu.id (ignored here)
			new(sql.NullString), // uu.email (ignored here)
			new(sql.NullString), // uu.password_hash (ignored here)
		)
	} else {
		err = r.db.QueryRow(queries.InsertProductWithUsersCTE,
			p.ProductName,
			p.ProductUrl,
			p.ProductImageUrl,
			p.ProductDescription,
			p.CreatedByUser.ID,
			updatedByID,
		).Scan(
			&p.ID,
			&p.ProductName,
			&p.ProductUrl,
			&p.ProductImageUrl,
			&p.ProductDescription,
			&p.CreatedDatetime,
			&p.UpdatedDateTime,
			&p.CreatedByUser.ID,
			&p.CreatedByUser.Email,
			&p.CreatedByUser.PasswordHash,
			// Updated user fields may be NULL; scan into sql.NullString
			new(sql.NullString), // uu.id (ignored here)
			new(sql.NullString), // uu.email (ignored here)
			new(sql.NullString), // uu.password_hash (ignored here)
		)
	}
	if err != nil {
		return p, err
	}
	// Hydrate updated user if present in the row
	var updID sql.NullString
	if err := r.db.QueryRow(queries.SelectProductUpdatedByUserID, p.ID).Scan(&updID); err == nil && updID.Valid {
		var uu models.User
		if err := r.db.QueryRow(queries.GetUserByID, updID.String).Scan(&uu.ID, &uu.Email, &uu.PasswordHash); err == nil {
			p.UpdatedByUser = &uu
		} else {
			p.UpdatedByUser = nil
		}
	} else {
		p.UpdatedByUser = nil
	}
	return p, err
}

// UpdateProduct updates an existing product by ID
func (r *ProductRepository) UpdateProduct(id string, p models.Product) (models.Product, error) {
	p.ID = id
	var updatedByID interface{}
	if p.UpdatedByUser != nil {
		updatedByID = p.UpdatedByUser.ID
	} else {
		updatedByID = nil
	}
	err := r.db.QueryRow(queries.UpdateProductWithUsersCTE,
		p.ProductName,
		p.ProductUrl,
		p.ProductImageUrl,
		p.ProductDescription,
		updatedByID,
		id,
	).Scan(
		&p.ID,
		&p.ProductName,
		&p.ProductUrl,
		&p.ProductImageUrl,
		&p.ProductDescription,
		&p.CreatedDatetime,
		&p.UpdatedDateTime,
		&p.CreatedByUser.ID,
		&p.CreatedByUser.Email,
		&p.CreatedByUser.PasswordHash,
		new(sql.NullString), // uu.id (ignored)
		new(sql.NullString), // uu.email (ignored)
		new(sql.NullString), // uu.password_hash (ignored)
	)
	if err != nil {
		return p, err
	}
	// Hydrate updated user using the stored column if present
	var updID sql.NullString
	if err := r.db.QueryRow(queries.SelectProductUpdatedByUserID, p.ID).Scan(&updID); err == nil && updID.Valid {
		var uu models.User
		if err := r.db.QueryRow(queries.GetUserByID, updID.String).Scan(&uu.ID, &uu.Email, &uu.PasswordHash); err == nil {
			p.UpdatedByUser = &uu
		} else {
			p.UpdatedByUser = nil
		}
	} else {
		p.UpdatedByUser = nil
	}
	return p, err
}

// DeleteProduct deletes a product by ID
func (r *ProductRepository) DeleteProduct(id string) (bool, error) {
	res, err := r.db.Exec(queries.DeleteProductByID, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
