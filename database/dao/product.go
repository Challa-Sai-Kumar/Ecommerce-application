package dao

import (
	"database/sql"
	"ecommerce/database"
	"ecommerce/models"
	"fmt"
)

// CreateProduct inserts a new product into the database
func CreateProduct(product *models.Product) error {
	query := `INSERT INTO products (id, name, description, price, stock, category, created_date, updated_date)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, product.ID, product.Name, product.Description, product.Price, product.Stock, product.Category, product.CreatedDate, product.UpdatedDate)
	return err
}

// GetProducts retrieves all products
func GetProducts() ([]models.Product, error) {
	query := `SELECT id, name, description, price, stock, category FROM products`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func ReserveStock(tx *sql.Tx, qunatity int, ID string) error {
	var availableStock int
	var reservedStock int
	query := "SELECT stock, reserved_stock FROM products WHERE id = ? FOR UPDATE"
	row := tx.QueryRow(query, ID)
	err := row.Scan(&availableStock, &reservedStock)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no product found with id : %s", ID)
		}
		return err
	}

	if availableStock-reservedStock < qunatity {
		return fmt.Errorf("sufficient stock not available for quantity : %d", qunatity)
	}

	query = "UPDATE products SET reserved_stock = reserved_stock + ? WHERE id = ?"
	_, err = tx.Exec(query, qunatity, ID)
	return err

}

func DeductStockForOrder(tx *sql.Tx, orderID string) error {
	query := `
        UPDATE products p
		JOIN order_items oi ON p.id = oi.product_id
		SET p.stock = p.stock - oi.quantity, p.reserved_stock = p.reserved_stock - oi.quantity
		WHERE oi.order_id = ? AND p.stock >= oi.quantity AND p.reserved_stock >= oi.quantity
 		`
	result, err := tx.Exec(query, orderID)
	if err != nil {
		return fmt.Errorf("failed to deduct stock: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient stock or no matching records for order %s", orderID)
	}

	return nil
}

func RestoreReservedStock(tx *sql.Tx, orderID string) error {
	query := `
        UPDATE products p
        JOIN order_items oi ON p.id = oi.product_id
        SET p.reserved_stock = p.reserved_stock - oi.quantity
        WHERE oi.order_id = ? AND p.reserved_stock >= oi.quantity
    `
	result, err := tx.Exec(query, orderID)
	if err != nil {
		return fmt.Errorf("failed to deduct stock: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient reserved stock or no matching records for order %s", orderID)
	}

	return nil
}

// func UpdateInventory(item *models.ProductDetails) error {
// 	query := "UPDATE products SET  "
// }
