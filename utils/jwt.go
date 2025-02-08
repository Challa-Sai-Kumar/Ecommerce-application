package utils

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret is the secret key for signing the token. Store it securely!
var JWTSecret = []byte("your-very-secure-secret")

// GenerateJWT generates a JWT for a given user ID
func GenerateJWT(userID string) (string, error) {
	// Define token claims
	claims := jwt.MapClaims{
		"user_id": userID,                           // Include user-specific data
		"exp":     time.Now().Add(time.Hour).Unix(), // Set expiration time
		"iat":     time.Now().Unix(),                // Issued at time
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	signedToken, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateJWT validates the JWT token by parsing and checking its signature and claims.
func ValidateJWT(tokenStr string) bool {
	// Parse the JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token's signing method is HMAC (HS256 in this case)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return JWTSecret, nil
	})

	// If there's an error or the token is invalid, return false
	if err != nil {
		log.Println("Error parsing token:", err)
		return false
	}

	// Check if the token is valid (including expiration and other claims)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if the token is expired
		expirationTime, ok := claims["exp"].(float64)
		if !ok {
			log.Println("Token does not have an expiration time.")
			return false
		}

		// Check if the token has expired
		if time.Now().Unix() > int64(expirationTime) {
			log.Println("Token has expired.")
			return false
		}

		return true
	}

	log.Println("Invalid token claims.")
	return false
}
