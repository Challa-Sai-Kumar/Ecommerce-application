package dao

import (
	"database/sql"
	"ecommerce/database"
	"ecommerce/models"
)

// CreateCart creates a cart for the user
func CreateCart(cart *models.Cart) error {
	query := `INSERT INTO carts (user_id, created_at, updated_at)
              VALUES (?, NOW(), NOW())`
	_, err := database.DB.Exec(query, cart.UserID)
	return err
}

// AddProductToCart adds a product to the user's cart
func AddProductToCart(cart *models.Cart) error {
	query := `INSERT INTO carts (id, user_id, product_id, quantity, created_date) 
			VALUES (?,?,?,?,?)`
	_, err := database.DB.Exec(query, cart.ID, cart.UserID, cart.ProductID, cart.Quantity, cart.CreatedDate)
	return err
}

// GetCart retrieves the user's cart
func GetCartItems(executor database.QueryExecutor, userID string) ([]*models.ProductDetails, error) {
	query := "SELECT product_id, quantity, name, price " +
		"FROM carts " +
		"INNER JOIN products ON carts.product_id=products.id " +
		"WHERE carts.user_id=?"
	rows, err := executor.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.ProductDetails
	for rows.Next() {
		var item models.ProductDetails
		err = rows.Scan(&item.ID, &item.Quantity, &item.Name, &item.Price)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

func DeleteCartItems(tx *sql.Tx, userID string) error {
	query := `DELETE FROM carts WHERE user_id = ?`
	_, err := tx.Exec(query, userID)
	return err
}
