package repository

import (
	"database/sql"
	"price-story/pkg/models"
	"price-story/pkg/queries"
)

// PriceStoryRepository handles database operations for price stories
type PriceStoryRepository struct {
	db *sql.DB
}

// NewPriceStoryRepository creates a new price story repository
func NewPriceStoryRepository(db *sql.DB) *PriceStoryRepository {
	return &PriceStoryRepository{db: db}
}

// ListPriceStories retrieves all price stories
func (r *PriceStoryRepository) ListPriceStories() ([]models.PriceStory, error) {
	rows, err := r.db.Query(queries.ListPriceStories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []models.PriceStory
	for rows.Next() {
		var ps models.PriceStory
		var puuID, puuEmail, puuHash sql.NullString
		var uuID, uuEmail, uuHash sql.NullString
		if err := rows.Scan(
			&ps.ID,
			&ps.CreatedDatetime,
			&ps.UpdatedDateTime,
			&ps.Product.ID,
			&ps.Product.ProductName,
			&ps.Product.ProductUrl,
			&ps.Product.ProductImageUrl,
			&ps.Product.ProductDescription,
			&ps.Product.CreatedDatetime,
			&ps.Product.UpdatedDateTime,
			&ps.Product.CreatedByUser.ID,
			&ps.Product.CreatedByUser.Email,
			&ps.Product.CreatedByUser.PasswordHash,
			&puuID,
			&puuEmail,
			&puuHash,
			&ps.CreatedByUser.ID,
			&ps.CreatedByUser.Email,
			&ps.CreatedByUser.PasswordHash,
			&uuID,
			&uuEmail,
			&uuHash,
		); err != nil {
			return nil, err
		}
		if puuID.Valid {
			ps.Product.UpdatedByUser = &models.User{ID: puuID.String, Email: puuEmail.String, PasswordHash: puuHash.String}
		}
		if uuID.Valid {
			ps.UpdatedByUser = &models.User{ID: uuID.String, Email: uuEmail.String, PasswordHash: uuHash.String}
		}
		stories = append(stories, ps)
	}
	return stories, nil
}

// GetPriceStory retrieves a price story by ID
func (r *PriceStoryRepository) GetPriceStory(id string) (models.PriceStory, error) {
	var ps models.PriceStory
	var puuID, puuEmail, puuHash sql.NullString
	var uuID, uuEmail, uuHash sql.NullString
	err := r.db.QueryRow(queries.GetPriceStoryByID, id).Scan(
		&ps.ID,
		&ps.CreatedDatetime,
		&ps.UpdatedDateTime,
		&ps.Product.ID,
		&ps.Product.ProductName,
		&ps.Product.ProductUrl,
		&ps.Product.ProductImageUrl,
		&ps.Product.ProductDescription,
		&ps.Product.CreatedDatetime,
		&ps.Product.UpdatedDateTime,
		&ps.Product.CreatedByUser.ID,
		&ps.Product.CreatedByUser.Email,
		&ps.Product.CreatedByUser.PasswordHash,
		&puuID,
		&puuEmail,
		&puuHash,
		&ps.CreatedByUser.ID,
		&ps.CreatedByUser.Email,
		&ps.CreatedByUser.PasswordHash,
		&uuID,
		&uuEmail,
		&uuHash,
	)
	if err != nil {
		return ps, err
	}
	if puuID.Valid {
		ps.Product.UpdatedByUser = &models.User{ID: puuID.String, Email: puuEmail.String, PasswordHash: puuHash.String}
	}
	if uuID.Valid {
		ps.UpdatedByUser = &models.User{ID: uuID.String, Email: uuEmail.String, PasswordHash: uuHash.String}
	}
	return ps, nil
}

// CreatePriceStory creates a new price story
func (r *PriceStoryRepository) CreatePriceStory(ps models.PriceStory) (models.PriceStory, error) {
	var updatedByID interface{}
	if ps.UpdatedByUser != nil {
		updatedByID = ps.UpdatedByUser.ID
	} else {
		updatedByID = nil
	}
	var puuID, puuEmail, puuHash sql.NullString
	var uuID, uuEmail, uuHash sql.NullString
	err := r.db.QueryRow(queries.InsertPriceStoryWithUsersCTE,
		ps.Product.ID,
		ps.CreatedByUser.ID,
		updatedByID,
	).Scan(
		&ps.ID,
		&ps.CreatedDatetime,
		&ps.UpdatedDateTime,
		&ps.Product.ID,
		&ps.Product.ProductName,
		&ps.Product.ProductUrl,
		&ps.Product.ProductImageUrl,
		&ps.Product.ProductDescription,
		&ps.Product.CreatedDatetime,
		&ps.Product.UpdatedDateTime,
		&ps.Product.CreatedByUser.ID,
		&ps.Product.CreatedByUser.Email,
		&ps.Product.CreatedByUser.PasswordHash,
		&puuID,
		&puuEmail,
		&puuHash,
		&ps.CreatedByUser.ID,
		&ps.CreatedByUser.Email,
		&ps.CreatedByUser.PasswordHash,
		&uuID,
		&uuEmail,
		&uuHash,
	)
	if err != nil {
		return ps, err
	}
	if puuID.Valid {
		ps.Product.UpdatedByUser = &models.User{ID: puuID.String, Email: puuEmail.String, PasswordHash: puuHash.String}
	}
	if uuID.Valid {
		ps.UpdatedByUser = &models.User{ID: uuID.String, Email: uuEmail.String, PasswordHash: uuHash.String}
	}
	return ps, nil
}

// UpdatePriceStory updates an existing price story by ID
func (r *PriceStoryRepository) UpdatePriceStory(id string, ps models.PriceStory) (models.PriceStory, error) {
	ps.ID = id
	var updatedByID interface{}
	if ps.UpdatedByUser != nil {
		updatedByID = ps.UpdatedByUser.ID
	} else {
		updatedByID = nil
	}
	var puuID, puuEmail, puuHash sql.NullString
	var uuID, uuEmail, uuHash sql.NullString
	err := r.db.QueryRow(queries.UpdatePriceStoryWithUsersCTE,
		ps.Product.ID,
		updatedByID,
		id,
	).Scan(
		&ps.ID,
		&ps.CreatedDatetime,
		&ps.UpdatedDateTime,
		&ps.Product.ID,
		&ps.Product.ProductName,
		&ps.Product.ProductUrl,
		&ps.Product.ProductImageUrl,
		&ps.Product.ProductDescription,
		&ps.Product.CreatedDatetime,
		&ps.Product.UpdatedDateTime,
		&ps.Product.CreatedByUser.ID,
		&ps.Product.CreatedByUser.Email,
		&ps.Product.CreatedByUser.PasswordHash,
		&puuID,
		&puuEmail,
		&puuHash,
		&ps.CreatedByUser.ID,
		&ps.CreatedByUser.Email,
		&ps.CreatedByUser.PasswordHash,
		&uuID,
		&uuEmail,
		&uuHash,
	)
	if err != nil {
		return ps, err
	}
	if puuID.Valid {
		ps.Product.UpdatedByUser = &models.User{ID: puuID.String, Email: puuEmail.String, PasswordHash: puuHash.String}
	}
	if uuID.Valid {
		ps.UpdatedByUser = &models.User{ID: uuID.String, Email: uuEmail.String, PasswordHash: uuHash.String}
	}
	return ps, nil
}

// DeletePriceStory deletes a price story by ID
func (r *PriceStoryRepository) DeletePriceStory(id string) (bool, error) {
	res, err := r.db.Exec(queries.DeletePriceStoryByID, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
