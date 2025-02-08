package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// Secret key shared between your server and the webhook provider.
var secretKey = []byte("your-secret-key")

// VerifyWebhook verifies the authenticity of the webhook request.
func VerifyWebhook(r *http.Request) error {
	// Step 1: Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read request body")
	}
	defer r.Body.Close()

	// Step 2: Get the signature from headers
	signature := r.Header.Get("X-Signature")
	if signature == "" {
		return errors.New("missing signature header")
	}

	// Step 3: Generate HMAC-SHA256 of the body using the secret key
	hash := hmac.New(sha256.New, secretKey)
	hash.Write(body)
	expectedSignature := hex.EncodeToString(hash.Sum(nil))

	// Step 4: Compare the signature with the expected value
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return errors.New("invalid signature")
	}

	// Step 5: (Optional) Validate the timestamp
	timestamp := r.Header.Get("X-Timestamp")
	if timestamp != "" {
		reqTime, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return errors.New("invalid timestamp format")
		}

		if time.Since(reqTime) > 5*time.Minute {
			return errors.New("webhook request expired")
		}
	}

	return nil
}
