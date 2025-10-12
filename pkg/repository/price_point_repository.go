package repository

import (
	"database/sql"
	"price-story/pkg/models"
	"price-story/pkg/queries"
)

// PricePointRepository handles database operations for price points
type PricePointRepository struct {
	db *sql.DB
}

// NewPricePointRepository creates a new price point repository
func NewPricePointRepository(db *sql.DB) *PricePointRepository {
	return &PricePointRepository{db: db}
}

// ListPricePoints retrieves all price points with nested product and users
func (r *PricePointRepository) ListPricePoints() ([]models.PricePoint, error) {
	rows, err := r.db.Query(queries.ListPricePoints)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var points []models.PricePoint
	for rows.Next() {
		var pp models.PricePoint
		var puuID, puuEmail, puuHash sql.NullString
		var uuID, uuEmail, uuHash sql.NullString
		if err := rows.Scan(
			&pp.ID,
			&pp.Price,
			&pp.Timestamp,
			&pp.CreatedDatetime,
			&pp.UpdatedDateTime,
			&pp.Product.ID,
			&pp.Product.ProductName,
			&pp.Product.ProductUrl,
			&pp.Product.ProductImageUrl,
			&pp.Product.ProductDescription,
			&pp.Product.CreatedDatetime,
			&pp.Product.UpdatedDateTime,
			&pp.Product.CreatedByUser.ID,
			&pp.Product.CreatedByUser.Email,
			&pp.Product.CreatedByUser.PasswordHash,
			&puuID,
			&puuEmail,
			&puuHash,
			&pp.CreatedByUser.ID,
			&pp.CreatedByUser.Email,
			&pp.CreatedByUser.PasswordHash,
			&uuID,
			&uuEmail,
			&uuHash,
		); err != nil {
			return nil, err
		}
		if puuID.Valid {
			pp.Product.UpdatedByUser = &models.User{ID: puuID.String, Email: puuEmail.String, PasswordHash: puuHash.String}
		}
		if uuID.Valid {
			pp.UpdatedByUser = &models.User{ID: uuID.String, Email: uuEmail.String, PasswordHash: uuHash.String}
		}
		points = append(points, pp)
	}
	return points, nil
}

// GetPricePoint retrieves a price point by ID with nested product and users
func (r *PricePointRepository) GetPricePoint(id string) (models.PricePoint, error) {
	var pp models.PricePoint
	var puuID, puuEmail, puuHash sql.NullString
	var uuID, uuEmail, uuHash sql.NullString
	err := r.db.QueryRow(queries.GetPricePointByID, id).Scan(
		&pp.ID,
		&pp.Price,
		&pp.Timestamp,
		&pp.CreatedDatetime,
		&pp.UpdatedDateTime,
		&pp.Product.ID,
		&pp.Product.ProductName,
		&pp.Product.ProductUrl,
		&pp.Product.ProductImageUrl,
		&pp.Product.ProductDescription,
		&pp.Product.CreatedDatetime,
		&pp.Product.UpdatedDateTime,
		&pp.Product.CreatedByUser.ID,
		&pp.Product.CreatedByUser.Email,
		&pp.Product.CreatedByUser.PasswordHash,
		&puuID,
		&puuEmail,
		&puuHash,
		&pp.CreatedByUser.ID,
		&pp.CreatedByUser.Email,
		&pp.CreatedByUser.PasswordHash,
		&uuID,
		&uuEmail,
		&uuHash,
	)
	if err != nil {
		return pp, err
	}
	if puuID.Valid {
		pp.Product.UpdatedByUser = &models.User{ID: puuID.String, Email: puuEmail.String, PasswordHash: puuHash.String}
	}
	if uuID.Valid {
		pp.UpdatedByUser = &models.User{ID: uuID.String, Email: uuEmail.String, PasswordHash: uuHash.String}
	}
	return pp, nil
}

// CreatePricePoint creates a new price point and returns hydrated nested entities
func (r *PricePointRepository) CreatePricePoint(pp models.PricePoint) (models.PricePoint, error) {
	var updatedByID interface{}
	if pp.UpdatedByUser != nil {
		updatedByID = pp.UpdatedByUser.ID
	} else {
		updatedByID = nil
	}
	var puuID, puuEmail, puuHash sql.NullString
	var uuID, uuEmail, uuHash sql.NullString
	err := r.db.QueryRow(queries.InsertPricePointWithUsersCTE,
		pp.Price,
		pp.Timestamp,
		pp.Product.ID,
		pp.CreatedByUser.ID,
		updatedByID,
	).Scan(
		&pp.ID,
		&pp.Price,
		&pp.Timestamp,
		&pp.CreatedDatetime,
		&pp.UpdatedDateTime,
		&pp.Product.ID,
		&pp.Product.ProductName,
		&pp.Product.ProductUrl,
		&pp.Product.ProductImageUrl,
		&pp.Product.ProductDescription,
		&pp.Product.CreatedDatetime,
		&pp.Product.UpdatedDateTime,
		&pp.Product.CreatedByUser.ID,
		&pp.Product.CreatedByUser.Email,
		&pp.Product.CreatedByUser.PasswordHash,
		&puuID,
		&puuEmail,
		&puuHash,
		&pp.CreatedByUser.ID,
		&pp.CreatedByUser.Email,
		&pp.CreatedByUser.PasswordHash,
		&uuID,
		&uuEmail,
		&uuHash,
	)
	if err != nil {
		return pp, err
	}
	if puuID.Valid {
		pp.Product.UpdatedByUser = &models.User{ID: puuID.String, Email: puuEmail.String, PasswordHash: puuHash.String}
	}
	if uuID.Valid {
		pp.UpdatedByUser = &models.User{ID: uuID.String, Email: uuEmail.String, PasswordHash: uuHash.String}
	}
	return pp, nil
}

// UpdatePricePoint updates an existing price point by ID and returns hydrated nested entities
func (r *PricePointRepository) UpdatePricePoint(id string, pp models.PricePoint) (models.PricePoint, error) {
	pp.ID = id
	var updatedByID interface{}
	if pp.UpdatedByUser != nil {
		updatedByID = pp.UpdatedByUser.ID
	} else {
		updatedByID = nil
	}
	var puuID, puuEmail, puuHash sql.NullString
	var uuID, uuEmail, uuHash sql.NullString
	err := r.db.QueryRow(queries.UpdatePricePointWithUsersCTE,
		pp.Price,
		pp.Timestamp,
		pp.Product.ID,
		updatedByID,
		id,
	).Scan(
		&pp.ID,
		&pp.Price,
		&pp.Timestamp,
		&pp.CreatedDatetime,
		&pp.UpdatedDateTime,
		&pp.Product.ID,
		&pp.Product.ProductName,
		&pp.Product.ProductUrl,
		&pp.Product.ProductImageUrl,
		&pp.Product.ProductDescription,
		&pp.Product.CreatedDatetime,
		&pp.Product.UpdatedDateTime,
		&pp.Product.CreatedByUser.ID,
		&pp.Product.CreatedByUser.Email,
		&pp.Product.CreatedByUser.PasswordHash,
		&puuID,
		&puuEmail,
		&puuHash,
		&pp.CreatedByUser.ID,
		&pp.CreatedByUser.Email,
		&pp.CreatedByUser.PasswordHash,
		&uuID,
		&uuEmail,
		&uuHash,
	)
	if err != nil {
		return pp, err
	}
	if puuID.Valid {
		pp.Product.UpdatedByUser = &models.User{ID: puuID.String, Email: puuEmail.String, PasswordHash: puuHash.String}
	}
	if uuID.Valid {
		pp.UpdatedByUser = &models.User{ID: uuID.String, Email: uuEmail.String, PasswordHash: uuHash.String}
	}
	return pp, nil
}

// DeletePricePoint deletes a price point by ID
func (r *PricePointRepository) DeletePricePoint(id string) (bool, error) {
	res, err := r.db.Exec(queries.DeletePricePointByID, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
