package sqlite

import (
	"context"
	"embed"
	"fmt"

	"mark7888/speedtest-data-server/internal/logger"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// latestMigrationVersion is the highest version number in the migrations
// directory. Increment this constant whenever a new migration file is added.
const latestMigrationVersion int64 = 3

// Migrate runs database migrations using goose.
//
// For databases created before goose was introduced the function detects the
// old schema and fast-forwards the version table so that already-applied
// changes are not executed again.
func (s *SQLiteDB) Migrate() error {
	logger.Log.Info("Running SQLite database migrations")

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	goose.SetBaseFS(migrationsFS)

	if err := s.bootstrapGooseIfNeeded(); err != nil {
		return fmt.Errorf("failed to bootstrap goose version table: %w", err)
	}

	if err := goose.Up(s.db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Log.Info("SQLite migrations completed successfully")
	return nil
}

// bootstrapGooseIfNeeded detects a pre-goose database (nodes table present but
// no goose_db_version table) and inserts versioned records for all existing
// migrations so that goose treats them as already applied. This avoids ALTER
// TABLE failures on columns that already exist.
func (s *SQLiteDB) bootstrapGooseIfNeeded() error {
	ctx := context.Background()

	// Check whether the nodes table already exists.
	var tableCount int
	err := s.db.QueryRowContext(ctx,
		"SELECT count(*) FROM sqlite_master WHERE type='table' AND name='nodes'",
	).Scan(&tableCount)
	if err != nil {
		return fmt.Errorf("failed to check for nodes table: %w", err)
	}
	if tableCount == 0 {
		// Fresh database – let goose run all migrations normally.
		return nil
	}

	// Check whether the goose version table already exists.
	var gooseTableCount int
	err = s.db.QueryRowContext(ctx,
		"SELECT count(*) FROM sqlite_master WHERE type='table' AND name='goose_db_version'",
	).Scan(&gooseTableCount)
	if err != nil {
		return fmt.Errorf("failed to check for goose_db_version table: %w", err)
	}
	if gooseTableCount > 0 {
		// Goose is already initialized – nothing to do.
		return nil
	}

	// Pre-goose database: create the version table and mark all known
	// migrations as applied without executing them.
	logger.Log.Info("Pre-goose SQLite database detected: fast-forwarding migration history",
		zap.Int64("target_version", latestMigrationVersion),
	)

	_, err = s.db.ExecContext(ctx, `
		CREATE TABLE goose_db_version (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			version_id INTEGER   NOT NULL,
			is_applied INTEGER   NOT NULL,
			tstamp     TIMESTAMP DEFAULT (datetime('now'))
		)`)
	if err != nil {
		return fmt.Errorf("failed to create goose_db_version table: %w", err)
	}

	for v := int64(1); v <= latestMigrationVersion; v++ {
		if _, err = s.db.ExecContext(ctx,
			`INSERT INTO goose_db_version (version_id, is_applied) VALUES (?, 1)`, v,
		); err != nil {
			return fmt.Errorf("failed to record migration %d as applied: %w", v, err)
		}
	}

	logger.Log.Info("Migration history fast-forwarded successfully",
		zap.Int64("version", latestMigrationVersion),
	)
	return nil
}
