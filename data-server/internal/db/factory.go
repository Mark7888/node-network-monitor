package db

import (
	"fmt"

	"mark7888/speedtest-data-server/internal/config"
	"mark7888/speedtest-data-server/internal/db/postgres"
	"mark7888/speedtest-data-server/internal/db/sqlite"
)

// New creates a new database connection based on configuration
// Returns a Database interface that can be either PostgreSQL or SQLite
func New(cfg *config.Config) (Database, error) {
	switch cfg.Database.Type {
	case "postgres":
		return postgres.New(cfg)
	case "sqlite":
		return sqlite.New(cfg)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}
}
