package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Notify(w http.ResponseWriter, r *http.Request) {
	var details struct {
		OrderID string `json:"order_id"`
		Amount  int    `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := SimulateWebhook(details.OrderID, details.Amount)
	if err != nil {
		http.Error(w, "Payment failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "payment successful"})
}

func SimulateWebhook(orderID string, amount int) error {
	webhookURL := "http://localhost:8080/payments/webhook"

	payload := map[string]interface{}{
		"order_id":       orderID,
		"transaction_id": "txn_67890",
		"amount":         amount,
		"status":         "success",
	}

	payloadBytes, _ := json.Marshal(payload)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return err
	}

	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY3MDA1NjAsImlhdCI6MTczNjY5Njk2MCwidXNlcl9pZCI6ImEwNjJlM2Q0ZGFmYjRhNjZiZWFiYTY5ZGRmYTRiNGJkIn0.Cjn5HtlndZR8fCPDWAuzsi9Dhr0e3WVtqPl744SCkvw"
	// Add headers (e.g., for authentication)
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token)) // Replace with the actual token or key

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending webhook: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("Webhook response status: %v", resp.Status)
	return nil
}
