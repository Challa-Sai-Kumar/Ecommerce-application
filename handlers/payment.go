package handlers

import (
	"ecommerce/database"
	"ecommerce/database/dao"
	"ecommerce/kafka"
	"ecommerce/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Payment struct {
	producer *kafka.Producer
}

func NewPayment(producer *kafka.Producer) *Payment {
	return &Payment{producer: producer}
}

func InitiatePayment(w http.ResponseWriter, r *http.Request) {
	var paymentRequest models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate the order
	order, err := dao.GetOrderByID(paymentRequest.OrderID)
	if err != nil || order.Status != "pending" {
		http.Error(w, "Invalid or non-pending order", http.StatusBadRequest)
		return
	}

	// Simulate payment link generation (replace with actual payment gateway integration)
	paymentLink := fmt.Sprintf("https://paymentgateway.com/pay?order_id=%s&amount=%.2f",
		order.ID, order.TotalPrice)

	// Optionally, save the payment initiation details in the database
	// err = dao.SavePaymentDetails(order.ID, paymentRequest.Method, "pending")
	// if err != nil {
	// 	http.Error(w, "Unable to initiate payment", http.StatusInternalServerError)
	// 	return
	// }

	response := map[string]string{
		"payment_link": paymentLink,
		"status":       "pending",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (p *Payment) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	var webhookData models.WebhookData
	if err := json.NewDecoder(r.Body).Decode(&webhookData); err != nil {
		log.Printf("invalid payload, err : %s", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// // Verify the webhook request (e.g., HMAC signature)
	// if err := utils.VerifyWebhook(r); err != nil { // Implement VerifyWebhook for your gateway
	// 	log.Printf("invalid webhook request, err : %s", err)
	// 	http.Error(w, "Unauthorized request", http.StatusUnauthorized)
	// 	return
	// }

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Unable to start transaction, err: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	// Update order and payment status
	if webhookData.Status == "success" {
		if err := dao.UpdateOrderStatus(tx, webhookData.OrderID, "completed"); err != nil {
			tx.Rollback()
			log.Printf("Failed to update order status, err: %v", err)
			http.Error(w, "Failed to update order status", http.StatusInternalServerError)
			return
		}

		if err := dao.DeductStockForOrder(tx, webhookData.OrderID); err != nil {
			tx.Rollback()
			log.Printf("Failed to deduct stock, err: %v", err)
			http.Error(w, "Failed to deduct stock", http.StatusInternalServerError)
			return
		}

		p.producer.PublishOrderStatus(webhookData.OrderID, "order placed")
		// kafka.PublishOrderStatus(webhookData.OrderID, "order placed")
		// sendOrderStatusDetails(webhookData.OrderID, "order placed")

	} else if webhookData.Status == "failed" {
		if err := dao.UpdateOrderStatus(tx, webhookData.OrderID, "canceled"); err != nil {
			tx.Rollback()
			log.Printf("Failed to update order status, err: %v", err)
			http.Error(w, "Failed to update order status", http.StatusInternalServerError)
			return
		}

		if err := dao.RestoreReservedStock(tx, webhookData.OrderID); err != nil {
			tx.Rollback()
			log.Printf("Failed to restore stock, err: %v", err)
			http.Error(w, "Failed to restore stock", http.StatusInternalServerError)
			return
		}

		// sendOrderStatusDetails(webhookData.OrderID, "order failed")
		p.producer.PublishOrderStatus(webhookData.OrderID, "order failed")

	}

	if err := tx.Commit(); err != nil {
		log.Printf("Transaction commit failed, err: %v", err)
		http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook processed successfully"))
}

// func (p *Payment) sendOrderStatusDetails(webhookData.OrderID, "order place")
