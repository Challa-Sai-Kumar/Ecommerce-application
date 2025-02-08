package models

import (
	"time"
)

// Order structure
type Order struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	TotalPrice  float32   `json:"total_price" db:"total_price"`
	Status      string    `json:"status" db:"status"` // Example: "pending", "shipped", "delivered"
	CreatedDate time.Time `json:"created_date" db:"created_at"`
}

type OrderDetails struct {
	ID         string  `json:"id" db:"id"`
	UserID     string  `json:"user_id" db:"user_id"`
	TotalPrice float32 `json:"total_price" db:"total_price"`
	Status     string  `json:"status" db:"status"` // Example: "pending", "shipped", "delivered"
	Username   string
	Email      string
}
