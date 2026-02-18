package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"mark7888/speedtest-data-server/internal/config"
	"mark7888/speedtest-data-server/internal/logger"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// SQLiteDB implements the Database interface for SQLite
type SQLiteDB struct {
	db      *sql.DB
	builder sq.StatementBuilderType

	// Ping cache to prevent DDoS via health checks
	pingMutex    sync.RWMutex
	lastPingTime time.Time
	lastPingErr  error
}

// New creates a new SQLite database connection
func New(cfg *config.Config) (*SQLiteDB, error) {
	// Create directory if it doesn't exist
	dbPath := cfg.Database.Path
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	logger.Log.Info("Connecting to SQLite database", zap.String("path", dbPath))

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// SQLite specific settings
	db.SetMaxOpenConns(1) // SQLite works best with single connection
	db.SetMaxIdleConns(1)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Log.Info("SQLite connection established")

	// Create statement builder with ? placeholder format
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Question)

	return &SQLiteDB{
		db:      db,
		builder: builder,
	}, nil
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	logger.Log.Info("Closing SQLite connection")
	return s.db.Close()
}

// Ping checks if the database connection is alive
func (s *SQLiteDB) Ping() error {
	ctx, cancel := withTimeout()
	defer cancel()
	return s.db.PingContext(ctx)
}

// SafePing checks if the database connection is alive with caching to prevent DDoS
// Only performs an actual ping if the last ping was more than 5 seconds ago
func (s *SQLiteDB) SafePing() error {
	// First, try to read the cached result
	s.pingMutex.RLock()
	if time.Since(s.lastPingTime) < 5*time.Second {
		err := s.lastPingErr
		s.pingMutex.RUnlock()
		return err
	}
	s.pingMutex.RUnlock()

	// Cache is stale, acquire write lock and ping
	s.pingMutex.Lock()
	defer s.pingMutex.Unlock()

	// Double-check in case another goroutine already updated the cache
	if time.Since(s.lastPingTime) < 5*time.Second {
		return s.lastPingErr
	}

	// Perform the actual ping
	ctx, cancel := withTimeout()
	defer cancel()
	err := s.db.PingContext(ctx)

	// Update cache
	s.lastPingTime = time.Now()
	s.lastPingErr = err

	return err
}

// withTimeout creates a context with a default timeout
func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
