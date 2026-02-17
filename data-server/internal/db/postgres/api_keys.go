package postgres

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
func (p *PostgresDB) CreateAPIKey(name, plainKey, createdBy string) (*models.APIKey, error) {
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

	query, args, err := p.builder.
		Insert("api_keys").
		Columns("id", "name", "key_hash", "enabled", "created_at", "created_by").
		Values(apiKey.ID, apiKey.Name, apiKey.KeyHash, apiKey.Enabled, apiKey.CreatedAt, apiKey.CreatedBy).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return apiKey, nil
}

// GetAPIKeyByID retrieves an API key by ID
func (p *PostgresDB) GetAPIKeyByID(id uuid.UUID) (*models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := p.builder.
		Select("id", "name", "key_hash", "enabled", "created_at", "created_by", "last_used", "revoked_at").
		From("api_keys").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var apiKey models.APIKey
	err = p.db.QueryRowContext(ctx, query, args...).Scan(
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
func (p *PostgresDB) GetAllAPIKeys() ([]models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := p.builder.
		Select("id", "name", "key_hash", "enabled", "created_at", "created_by", "last_used", "revoked_at").
		From("api_keys").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := p.db.QueryContext(ctx, query, args...)
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
func (p *PostgresDB) GetEnabledAPIKeys() ([]models.APIKey, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := p.builder.
		Select("id", "name", "key_hash", "enabled", "created_at", "created_by", "last_used", "revoked_at").
		From("api_keys").
		Where(sq.Eq{"enabled": true}).
		Where("revoked_at IS NULL").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := p.db.QueryContext(ctx, query, args...)
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
func (p *PostgresDB) UpdateAPIKeyEnabled(id uuid.UUID, enabled bool) error {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := p.builder.
		Update("api_keys").
		Set("enabled", enabled).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := p.db.ExecContext(ctx, query, args...)
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
func (p *PostgresDB) DeleteAPIKey(id uuid.UUID) error {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := p.builder.
		Delete("api_keys").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	result, err := p.db.ExecContext(ctx, query, args...)
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
func (p *PostgresDB) UpdateAPIKeyLastUsed(id uuid.UUID) error {
	ctx, cancel := withTimeout()
	defer cancel()

	query, args, err := p.builder.
		Update("api_keys").
		Set("last_used", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update API key last_used: %w", err)
	}

	return nil
}

// VerifyAPIKey verifies an API key and returns the key details if valid
func (p *PostgresDB) VerifyAPIKey(plainKey string) (*models.APIKey, error) {
	// Get all enabled keys
	apiKeys, err := p.GetEnabledAPIKeys()
	if err != nil {
		return nil, err
	}

	// Check each key
	for _, apiKey := range apiKeys {
		if auth.VerifyAPIKey(plainKey, apiKey.KeyHash) {
			// Update last used timestamp
			_ = p.UpdateAPIKeyLastUsed(apiKey.ID)
			return &apiKey, nil
		}
	}

	return nil, fmt.Errorf("invalid API key")
}
