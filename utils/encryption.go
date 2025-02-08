package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"
)

func Encrypt(data, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatalf("Error creating cipher: %v", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("Error creating GCM: %v", err)
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(encryptedData, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatalf("Error creating cipher: %v", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("Error creating GCM: %v", err)
		return "", err
	}

	data, err := hex.DecodeString(encryptedData)
	if err != nil {
		log.Fatalf("Error decoding hex: %v", err)
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalf("Error decrypting: %v", err)
	}

	return string(plaintext), nil
}
