package models

import (
	"time"
)

// User structure
type User struct {
	ID          string    `json:"id" db:"id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"password" db:"password"`
	Role        string    `json:"role,omitempty" db:"role"` // Example: "customer", "admin"
	CreatedDate time.Time `json:"created_date,omitempty" db:"created_date"`
	UpdatedDate time.Time `json:"updated_date,omitempty" db:"updated_date"`
}
