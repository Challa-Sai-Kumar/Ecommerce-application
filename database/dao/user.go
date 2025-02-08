package dao

import (
	"database/sql"
	"ecommerce/database"
	"ecommerce/models"
	"fmt"
	"log"
)

// CreateUser inserts a new user into the database
func CreateUser(user *models.User) error {
	query := `INSERT INTO users (id,first_name, last_name, email, password, created_date, updated_date) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.CreatedDate, user.UpdatedDate)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// GetUserByID retrieves a user by ID
func GetUserByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, first_name, last_name, email, password FROM users WHERE id = ?`

	// Use QueryRow to fetch a single row
	row := database.DB.QueryRow(query, id)

	// Scan the row into the user struct
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id %s", id)
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, first_name, last_name, email, password FROM users WHERE email = ?`

	// Use QueryRow to fetch a single row
	row := database.DB.QueryRow(query, email)

	// Scan the row into the user struct
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email %s", email)
		}
		return nil, err
	}

	return &user, nil
}
