package models

import "time"

// PriceStory captures the story of price changes for a product
type PriceStory struct {
	ID              string    `json:"id"`
	Product         Product   `json:"product"`
	CreatedByUser   User      `json:"created_by_user"`
	UpdatedByUser   *User     `json:"updated_by_user,omitempty"`
	CreatedDatetime time.Time `json:"created_datetime"`
	UpdatedDateTime time.Time `json:"updated_datetime"`
}
