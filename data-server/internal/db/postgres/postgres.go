package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"mark7888/speedtest-data-server/internal/config"
	"mark7888/speedtest-data-server/internal/logger"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// PostgresDB implements the Database interface for PostgreSQL
type PostgresDB struct {
	db      *sql.DB
	builder sq.StatementBuilderType

	// Ping cache to prevent DDoS via health checks
	pingMutex    sync.RWMutex
	lastPingTime time.Time
	lastPingErr  error
}

// New creates a new PostgreSQL database connection
func New(cfg *config.Config) (*PostgresDB, error) {
	dsn := cfg.GetDSN()

	logger.Log.Info("Connecting to PostgreSQL database",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("database", cfg.Database.Name),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL database: %w", err)
	}

	// PostgreSQL connection pool settings
	db.SetMaxOpenConns(cfg.Database.MaxConnections)
	db.SetMaxIdleConns(cfg.Database.MaxIdle)
	db.SetConnMaxLifetime(cfg.Database.ConnectionLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Log.Info("PostgreSQL connection established")

	// Create statement builder with PostgreSQL placeholder format ($1, $2, ...)
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &PostgresDB{
		db:      db,
		builder: builder,
	}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	logger.Log.Info("Closing PostgreSQL connection")
	return p.db.Close()
}

// Ping checks if the database connection is alive
func (p *PostgresDB) Ping() error {
	ctx, cancel := withTimeout()
	defer cancel()
	return p.db.PingContext(ctx)
}

// SafePing checks if the database connection is alive with caching to prevent DDoS
// Only performs an actual ping if the last ping was more than 5 seconds ago
func (p *PostgresDB) SafePing() error {
	// First, try to read the cached result
	p.pingMutex.RLock()
	if time.Since(p.lastPingTime) < 5*time.Second {
		err := p.lastPingErr
		p.pingMutex.RUnlock()
		return err
	}
	p.pingMutex.RUnlock()

	// Cache is stale, acquire write lock and ping
	p.pingMutex.Lock()
	defer p.pingMutex.Unlock()

	// Double-check in case another goroutine already updated the cache
	if time.Since(p.lastPingTime) < 5*time.Second {
		return p.lastPingErr
	}

	// Perform the actual ping
	ctx, cancel := withTimeout()
	defer cancel()
	err := p.db.PingContext(ctx)

	// Update cache
	p.lastPingTime = time.Now()
	p.lastPingErr = err

	return err
}

// withTimeout creates a context with a default timeout
func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
