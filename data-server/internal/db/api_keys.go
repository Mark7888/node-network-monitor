package db

import (
	"database/sql"
	"fmt"
	"time"

	"mark7888/speedtest-data-server/internal/auth"
	"mark7888/speedtest-data-server/pkg/models"

	"github.com/google/uuid"
)

// CreateAPIKey creates a new API key
func (db *DB) CreateAPIKey(name, plainKey, createdBy string) (*models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	// Hash the key
	keyHash, err := auth.HashAPIKey(plainKey)
	if err != nil {
		return nil, fmt.Errorf("failed to hash API key: %w", err)
	}

	apiKey := &models.APIKey{
		ID:        uuid.New(),
		Name:      name,
		KeyHash:   keyHash,
		Enabled:   true,
		CreatedAt: time.Now().UTC(),
	}

	if createdBy != "" {
		apiKey.CreatedBy = &createdBy
	}

	query := `
		INSERT INTO api_keys (id, name, key_hash, enabled, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = db.ExecContext(ctx, query,
		apiKey.ID,
		apiKey.Name,
		apiKey.KeyHash,
		apiKey.Enabled,
		apiKey.CreatedAt,
		apiKey.CreatedBy,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return apiKey, nil
}

// GetAPIKeyByID retrieves an API key by ID
func (db *DB) GetAPIKeyByID(id uuid.UUID) (*models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query := `
		SELECT id, name, key_hash, enabled, created_at, created_by, last_used, revoked_at
		FROM api_keys
		WHERE id = $1
	`

	var apiKey models.APIKey
	err := db.QueryRowContext(ctx, query, id).Scan(
		&apiKey.ID,
		&apiKey.Name,
		&apiKey.KeyHash,
		&apiKey.Enabled,
		&apiKey.CreatedAt,
		&apiKey.CreatedBy,
		&apiKey.LastUsed,
		&apiKey.RevokedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API key not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return &apiKey, nil
}

// GetAllAPIKeys retrieves all API keys
func (db *DB) GetAllAPIKeys() ([]models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query := `
		SELECT id, name, key_hash, enabled, created_at, created_by, last_used, revoked_at
		FROM api_keys
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		err := rows.Scan(
			&apiKey.ID,
			&apiKey.Name,
			&apiKey.KeyHash,
			&apiKey.Enabled,
			&apiKey.CreatedAt,
			&apiKey.CreatedBy,
			&apiKey.LastUsed,
			&apiKey.RevokedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

// GetEnabledAPIKeys retrieves all enabled API keys
func (db *DB) GetEnabledAPIKeys() ([]models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query := `
		SELECT id, name, key_hash, enabled, created_at, created_by, last_used, revoked_at
		FROM api_keys
		WHERE enabled = true AND revoked_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query enabled API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		err := rows.Scan(
			&apiKey.ID,
			&apiKey.Name,
			&apiKey.KeyHash,
			&apiKey.Enabled,
			&apiKey.CreatedAt,
			&apiKey.CreatedBy,
			&apiKey.LastUsed,
			&apiKey.RevokedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

// UpdateAPIKeyEnabled updates the enabled status of an API key
func (db *DB) UpdateAPIKeyEnabled(id uuid.UUID, enabled bool) error {
	ctx, cancel := withTimeout()
	defer cancel()

	query := `
		UPDATE api_keys
		SET enabled = $1
		WHERE id = $2
	`

	result, err := db.ExecContext(ctx, query, enabled, id)
	if err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// DeleteAPIKey deletes an API key
func (db *DB) DeleteAPIKey(id uuid.UUID) error {
	ctx, cancel := withTimeout()
	defer cancel()

	query := "DELETE FROM api_keys WHERE id = $1"

	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// UpdateAPIKeyLastUsed updates the last_used timestamp of an API key
func (db *DB) UpdateAPIKeyLastUsed(id uuid.UUID) error {
	ctx, cancel := withTimeout()
	defer cancel()

	nowSQL := db.getNowSQL()
	query := fmt.Sprintf(`
		UPDATE api_keys
		SET last_used = %s
		WHERE id = $1
	`, nowSQL)

	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update API key last_used: %w", err)
	}

	return nil
}

// VerifyAPIKey verifies an API key and returns the key details if valid
func (db *DB) VerifyAPIKey(plainKey string) (*models.APIKey, error) {
	// Get all enabled keys
	apiKeys, err := db.GetEnabledAPIKeys()
	if err != nil {
		return nil, err
	}

	// Check each key
	for _, apiKey := range apiKeys {
		if auth.VerifyAPIKey(plainKey, apiKey.KeyHash) {
			// Update last used timestamp
			_ = db.UpdateAPIKeyLastUsed(apiKey.ID)
			return &apiKey, nil
		}
	}

	return nil, fmt.Errorf("invalid API key")
}
