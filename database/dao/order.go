package dao

import (
	"database/sql"
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/utils"
	"fmt"
)

// CreateOrder creates a new order
func CreateOrder(tx *sql.Tx, order *models.Order) error {
	query := `INSERT INTO orders (id, user_id, total_price, status, created_date)
              VALUES (?, ?, ?, ?, ?)`
	_, err := tx.Exec(query, order.ID, order.UserID, order.TotalPrice, order.Status, order.CreatedDate)
	return err
}

func UpdateOrderItems(tx *sql.Tx, orderID string, item *models.ProductDetails) error {
	query := "INSERT INTO order_items (id, order_id, product_id, product_price, quantity) " +
		"VALUES (?,?,?,?,?)"
	_, err := tx.Exec(query, utils.NewID(), orderID, item.ID, item.Price, item.Quantity)
	return err
}

func GetOrderByID(ID string) (*models.Order, error) {
	query := "SELECT id, status, total_price FROM orders WHERE id = ?"

	var order models.Order
	err := database.DB.QueryRow(query, ID).Scan(&order.ID, &order.Status, &order.TotalPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no orders found with id : %s", ID)
		}
		return nil, err
	}
	return &order, nil
}

func UpdateOrderStatus(tx *sql.Tx, ID, status string) error {
	query := "UPDATE orders SET status = ? WHERE id = ?"
	_, err := tx.Exec(query, status, ID)
	return err
}

// GetOrders retrieves all orders for a specific user
func GetOrders(userID string) ([]models.Order, error) {
	query := `SELECT * FROM orders WHERE user_id = ?`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, order.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func GetOrderDetails(orderID string) (*models.OrderDetails, error) {
	query := "SELECT orders.id, orders.total_price, orders.status, users.id, CONCAT(users.first_name, ' ', users.last_name) AS username, users.email " +
		"FROM orders " +
		"INNER JOIN users ON users.id = orders.user_id " +
		"WHERE orders.id = ?"

	var orderDetails models.OrderDetails
	err := database.DB.QueryRow(query, orderID).Scan(&orderDetails.ID, &orderDetails.TotalPrice, &orderDetails.Status, &orderDetails.UserID, &orderDetails.Username, &orderDetails.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no orders details found with id : %s", orderID)
		}
		return nil, err
	}
	return &orderDetails, nil
}
