package models

import "time"

// Product represents a tracked product
type Product struct {
	ID                 string    `json:"id"`
	ProductName        string    `json:"product_name"`
	ProductUrl         string    `json:"product_url"`
	ProductImageUrl    string    `json:"product_image_url"`
	ProductDescription string    `json:"product_description"`
	CreatedByUser      User      `json:"created_by_user"`
	UpdatedByUser      *User     `json:"updated_by_user,omitempty"`
	CreatedDatetime    time.Time `json:"created_datetime"`
	UpdatedDateTime    time.Time `json:"updated_datetime"`
}
