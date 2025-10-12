package models

import "time"

// PricePoint represents a single price measurement for a product
type PricePoint struct {
	ID              string    `json:"id"`
	Price           float64   `json:"price"`
	Timestamp       time.Time `json:"timestamp"`
	Product         Product   `json:"product"`
	CreatedByUser   User      `json:"created_by_user"`
	UpdatedByUser   *User     `json:"updated_by_user,omitempty"`
	CreatedDatetime time.Time `json:"created_datetime"`
	UpdatedDateTime time.Time `json:"updated_datetime"`
}
