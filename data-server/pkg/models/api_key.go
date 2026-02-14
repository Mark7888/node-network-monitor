package models

import (
	"time"

	"github.com/google/uuid"
)

// APIKey represents an API key for node authentication
type APIKey struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	KeyHash   string     `json:"-" db:"key_hash"` // Never expose the hash
	Enabled   bool       `json:"enabled" db:"enabled"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	CreatedBy *string    `json:"created_by,omitempty" db:"created_by"`
	LastUsed  *time.Time `json:"last_used,omitempty" db:"last_used"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

// CreateAPIKeyRequest represents a request to create a new API key
type CreateAPIKeyRequest struct {
	Name string `json:"name" binding:"required,min=3,max=255"`
}

// CreateAPIKeyResponse represents the response when creating an API key
type CreateAPIKeyResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"` // Plain key, only shown once
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	Warning   string    `json:"warning"`
}

// UpdateAPIKeyRequest represents a request to update an API key
type UpdateAPIKeyRequest struct {
	Enabled *bool `json:"enabled"`
}

// ListAPIKeysResponse represents the response when listing API keys
type ListAPIKeysResponse struct {
	APIKeys []APIKey `json:"api_keys"`
	Total   int      `json:"total"`
}
