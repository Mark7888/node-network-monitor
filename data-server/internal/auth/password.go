package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost is the cost factor for bcrypt hashing
	BcryptCost = 12
	// APIKeyLength is the length of generated API keys in bytes
	APIKeyLength = 32
	// APIKeyPrefix is the prefix for API keys
	APIKeyPrefix = "sk_live_"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateAPIKey generates a new random API key
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, APIKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}

	key := APIKeyPrefix + base64.URLEncoding.EncodeToString(bytes)
	return key, nil
}

// HashAPIKey hashes an API key using bcrypt
func HashAPIKey(key string) (string, error) {
	return HashPassword(key)
}

// VerifyAPIKey verifies an API key against a hash
func VerifyAPIKey(key, hash string) bool {
	return VerifyPassword(key, hash)
}
