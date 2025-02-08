package models

import (
	"time"
)

// Cart structure
type Cart struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	ProductID   string    `json:"product_id" db:"product_id"` // store list of product IDs in cart
	Quantity    int       `json:"quantity" db:"quantity"`
	CreatedDate time.Time `json:"created_date" db:"created_at"`
}

type CartItems struct {
	UserID   string            `json:"userID"`
	Products []*ProductDetails `json:"products" db:"products"`
}

type ProductDetails struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
}
