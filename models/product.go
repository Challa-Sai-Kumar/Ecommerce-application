package models

import (
	"time"
)

// Product structure
type Product struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	Category    string    `json:"category" db:"category"`
	CreatedDate time.Time `json:"created_date" db:"created_at"`
	UpdatedDate time.Time `json:"updated_date" db:"updated_at"`
}
