package models

import "time"

type Address struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id" validate:"required"`
	Title      string    `json:"title,omitempty"`
	IsDefault  bool      `json:"is_default"`
	Line1      string    `json:"line_1" validate:"required"`
	Line2      string    `json:"line_2,omitempty"`
	PostalCode string    `json:"postal_code" validate:"required"`
	CreatedAt  time.Time `json:"created_at" validate:"required"`
	UpdatedAt  time.Time `json:"updated_at" validate:"required"`
}
