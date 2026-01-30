package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"-" validate:"required"`
	Password     string    `json:"-" validate:"required_without=ID"` // Only required for registration
	FirstName    string    `json:"first_name" validate:"required"`
	LastName     string    `json:"last_name,omitempty"`
	Verified     bool      `json:"verified"`
	Addresses    []Address `json:"addresses,omitempty"`
	CreatedAt    time.Time `json:"created_at" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at" validate:"required"`
	Phone        string    `json:"phone" validate:"required"`
}
