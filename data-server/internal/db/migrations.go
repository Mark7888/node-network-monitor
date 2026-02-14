package db

import (
	"fmt"

	"mark7888/speedtest-data-server/internal/logger"

	"go.uber.org/zap"
)

// Migrate runs database migrations
func (db *DB) Migrate() error {
	logger.Log.Info("Running database migrations")

	logger.Log.Info("Detected database type", zap.String("type", db.dbType))

	// Create nodes table
	_, err := db.Exec(getSQLForNodes(db.dbType))
	if err != nil {
		return fmt.Errorf("failed to create nodes table: %w", err)
	}

	// Create measurements table
	_, err = db.Exec(getSQLForMeasurements(db.dbType))
	if err != nil {
		return fmt.Errorf("failed to create measurements table: %w", err)
	}

	// Create failed_measurements table
	_, err = db.Exec(getSQLForFailedMeasurements(db.dbType))
	if err != nil {
		return fmt.Errorf("failed to create failed_measurements table: %w", err)
	}

	// Create api_keys table
	_, err = db.Exec(getSQLForAPIKeys(db.dbType))
	if err != nil {
		return fmt.Errorf("failed to create api_keys table: %w", err)
	}

	// Create indexes
	if err := db.createIndexes(db.dbType); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	logger.Log.Info("Database migrations completed successfully")
	return nil
}
