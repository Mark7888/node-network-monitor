package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mark7888/speedtest-data-server/internal/config"
	"mark7888/speedtest-data-server/internal/logger"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// DB wraps the database connection
type DB struct {
	*sql.DB
	dbType string // "postgres" or "sqlite"
}

// New creates a new database connection
func New(cfg *config.Config) (*DB, error) {
	var db *sql.DB
	var err error

	if cfg.Database.Type == "sqlite" {
		// Create directory if it doesn't exist
		dbPath := cfg.Database.Path
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}

		logger.Log.Info("Connecting to SQLite database", zap.String("path", dbPath))

		db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite database: %w", err)
		}

		// SQLite specific settings
		db.SetMaxOpenConns(1) // SQLite works best with single connection
		db.SetMaxIdleConns(1)
	} else {
		// PostgreSQL
		dsn := cfg.GetDSN()

		logger.Log.Info("Connecting to PostgreSQL database",
			zap.String("host", cfg.Database.Host),
			zap.Int("port", cfg.Database.Port),
			zap.String("database", cfg.Database.Name),
		)

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open PostgreSQL database: %w", err)
		}

		// PostgreSQL connection pool settings
		db.SetMaxOpenConns(cfg.Database.MaxConnections)
		db.SetMaxIdleConns(cfg.Database.MaxIdle)
		db.SetConnMaxLifetime(cfg.Database.ConnectionLifetime)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Log.Info("Database connection established", zap.String("type", cfg.Database.Type))

	return &DB{DB: db, dbType: cfg.Database.Type}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	logger.Log.Info("Closing database connection")
	return db.DB.Close()
}

// Ping checks if the database connection is alive
func (db *DB) Ping() error {
	ctx, cancel := withTimeout()
	defer cancel()
	return db.DB.PingContext(ctx)
}

// withTimeout creates a context with a default timeout
func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
