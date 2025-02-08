package models

import "time"

type PaymentRequest struct {
	OrderID string `json:"order_id"`
	Method  string `json:"method"` // e.g., "card", "wallet"
}

type paymentDetails struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	TransactionId string    `json:"transaction_id"`
	PaymentStatus string    `json:"payment_status"`
	Amount        int       `json:"amount"`
	CreatedDate   time.Time `json:"created_date"`
	UpdatedDate   time.Time `json:"updated_date"`
}

type WebhookData struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	Amount        int    `json:"amount"`
	Status        string `json:"status"` // "success", "failed", etc.
}
