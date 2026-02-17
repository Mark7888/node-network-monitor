package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"mark7888/speedtest-data-server/internal/auth"
	"mark7888/speedtest-data-server/pkg/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// CreateAPIKey creates a new API key
func (s *SQLiteDB) CreateAPIKey(name, plainKey, createdBy string) (*models.APIKey, error) {
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

	query, args, err := s.builder.
		Insert("api_keys").
		Columns("id", "name", "key_hash", "enabled", "created_at", "created_by").
		Values(apiKey.ID.String(), apiKey.Name, apiKey.KeyHash, 1, apiKey.CreatedAt, apiKey.CreatedBy).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return apiKey, nil
}

// GetAPIKeyByID retrieves an API key by ID
func (s *SQLiteDB) GetAPIKeyByID(id uuid.UUID) (*models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := s.builder.
		Select("id", "name", "key_hash", "enabled", "created_at", "created_by", "last_used", "revoked_at").
		From("api_keys").
		Where(sq.Eq{"id": id.String()}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var apiKey models.APIKey
	var idStr string
	var enabled int

	err = s.db.QueryRowContext(ctx, query, args...).Scan(
		&idStr,
		&apiKey.Name,
		&apiKey.KeyHash,
		&enabled,
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

	apiKey.ID, _ = uuid.Parse(idStr)
	apiKey.Enabled = enabled == 1

	return &apiKey, nil
}

// GetAllAPIKeys retrieves all API keys
func (s *SQLiteDB) GetAllAPIKeys() ([]models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := s.builder.
		Select("id", "name", "key_hash", "enabled", "created_at", "created_by", "last_used", "revoked_at").
		From("api_keys").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		var idStr string
		var enabled int

		err := rows.Scan(
			&idStr,
			&apiKey.Name,
			&apiKey.KeyHash,
			&enabled,
			&apiKey.CreatedAt,
			&apiKey.CreatedBy,
			&apiKey.LastUsed,
			&apiKey.RevokedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}

		apiKey.ID, _ = uuid.Parse(idStr)
		apiKey.Enabled = enabled == 1
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

// GetEnabledAPIKeys retrieves all enabled API keys
func (s *SQLiteDB) GetEnabledAPIKeys() ([]models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := s.builder.
		Select("id", "name", "key_hash", "enabled", "created_at", "created_by", "last_used", "revoked_at").
		From("api_keys").
		Where(sq.Eq{"enabled": 1}).
		Where("revoked_at IS NULL").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query enabled API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		var idStr string
		var enabled int

		err := rows.Scan(
			&idStr,
			&apiKey.Name,
			&apiKey.KeyHash,
			&enabled,
			&apiKey.CreatedAt,
			&apiKey.CreatedBy,
			&apiKey.LastUsed,
			&apiKey.RevokedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}

		apiKey.ID, _ = uuid.Parse(idStr)
		apiKey.Enabled = enabled == 1
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

// UpdateAPIKeyEnabled updates the enabled status of an API key
func (s *SQLiteDB) UpdateAPIKeyEnabled(id uuid.UUID, enabled bool) error {
	ctx, cancel := withTimeout()
	defer cancel()

	enabledInt := 0
	if enabled {
		enabledInt = 1
	}

	query, args, err := s.builder.
		Update("api_keys").
		Set("enabled", enabledInt).
		Where(sq.Eq{"id": id.String()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
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
func (s *SQLiteDB) DeleteAPIKey(id uuid.UUID) error {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := s.builder.
		Delete("api_keys").
		Where(sq.Eq{"id": id.String()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
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
func (s *SQLiteDB) UpdateAPIKeyLastUsed(id uuid.UUID) error {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := s.builder.
		Update("api_keys").
		Set("last_used", time.Now().UTC()).
		Where(sq.Eq{"id": id.String()}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update API key last_used: %w", err)
	}

	return nil
}

// VerifyAPIKey verifies an API key and returns the key details if valid
func (s *SQLiteDB) VerifyAPIKey(plainKey string) (*models.APIKey, error) {
	// Get all enabled keys
	apiKeys, err := s.GetEnabledAPIKeys()
	if err != nil {
		return nil, err
	}

	// Check each key
	for _, apiKey := range apiKeys {
		if auth.VerifyAPIKey(plainKey, apiKey.KeyHash) {
			// Update last used timestamp
			_ = s.UpdateAPIKeyLastUsed(apiKey.ID)
			return &apiKey, nil
		}
	}

	return nil, fmt.Errorf("invalid API key")
}
