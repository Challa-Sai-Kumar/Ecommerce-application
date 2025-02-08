package handlers

import (
	"ecommerce/database"
	"ecommerce/database/dao"
	"ecommerce/models"
	"ecommerce/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// CreateOrder handles creating an order
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	order.ID = utils.NewID()
	order.CreatedDate = time.Now()

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("unable to start transaction, err : %s", err)
		http.Error(w, "Unable to process order", http.StatusInternalServerError)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("recovered from panic: %v", r)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	items, err := dao.GetCartItems(tx, order.UserID)
	if err != nil {
		log.Printf("unable to fetch order price, err : %s", err)
		tx.Rollback()
		http.Error(w, "unable to fetch order price", http.StatusInternalServerError)
		return
	}
	if len(items) == 0 {
		log.Printf("your cart is empty")
		tx.Rollback()
		http.Error(w, "your cart is empty", http.StatusBadRequest)
	}

	for _, item := range items {
		err := dao.ReserveStock(tx, item.Quantity, item.ID)
		if err != nil {
			log.Printf("unable to reserve stock, err : %s", err)
			tx.Rollback()
			http.Error(w, "unable to reserve stock", http.StatusInternalServerError)
			return
		}
	}

	var totalPrice float32
	for _, item := range items {
		totalPrice += float32(item.Price * float32(item.Quantity))
	}

	order.TotalPrice = totalPrice
	order.Status = "pending"

	if err := dao.CreateOrder(tx, &order); err != nil {
		log.Printf("unable to creat order, err : %s", err)
		tx.Rollback()
		http.Error(w, "Unable to create order", http.StatusInternalServerError)
		return
	}

	for _, item := range items {
		err := dao.UpdateOrderItems(tx, order.ID, item)
		if err != nil {
			log.Printf("unable to update order items, err : %s", err)
			tx.Rollback()
			http.Error(w, "unable to update order items", http.StatusInternalServerError)
			return
		}
	}

	// Step 6: Delete cart items after order is placed
	err = dao.DeleteCartItems(tx, order.UserID)
	if err != nil {
		log.Printf("unable to delete cart items, err : %s", err)
		tx.Rollback()
		http.Error(w, "unable to delete cart items", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("unable to commit orders, err : %s", err)
		tx.Rollback()
		http.Error(w, "unable to commit orders", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"order_id": order.ID,
		"status":   "pending",
		"message":  "Order created successfully. Proceed to payment.",
	}
	json.NewEncoder(w).Encode(response)
}

// GetOrders handles fetching orders for a user
func GetOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	orders, err := dao.GetOrders(userID)
	if err != nil {
		http.Error(w, "Unable to fetch orders", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
