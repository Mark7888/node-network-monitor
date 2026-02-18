package sqlite

import (
	"fmt"

	"mark7888/speedtest-data-server/internal/logger"
)

// Migrate runs database migrations for SQLite
func (s *SQLiteDB) Migrate() error {
	logger.Log.Info("Running SQLite database migrations")

	// Create nodes table
	_, err := s.db.Exec(createNodesTable)
	if err != nil {
		return fmt.Errorf("failed to create nodes table: %w", err)
	}

	// Create measurements table
	_, err = s.db.Exec(createMeasurementsTable)
	if err != nil {
		return fmt.Errorf("failed to create measurements table: %w", err)
	}

	// Create failed_measurements table
	_, err = s.db.Exec(createFailedMeasurementsTable)
	if err != nil {
		return fmt.Errorf("failed to create failed_measurements table: %w", err)
	}

	// Create api_keys table
	_, err = s.db.Exec(createAPIKeysTable)
	if err != nil {
		return fmt.Errorf("failed to create api_keys table: %w", err)
	}

	// Add missing columns to existing tables
	if err := s.addMissingColumns(); err != nil {
		return fmt.Errorf("failed to add missing columns: %w", err)
	}

	// Create indexes
	if err := s.createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	logger.Log.Info("SQLite migrations completed successfully")
	return nil
}

// addMissingColumns adds any missing columns to existing tables
// This function is called during migration to ensure existing databases
// have all the columns from newer schema versions
func (s *SQLiteDB) addMissingColumns() error {
	// Check if archived column exists
	var archivedExists int
	err := s.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('nodes') WHERE name='archived'").Scan(&archivedExists)
	if err != nil {
		return fmt.Errorf("failed to check for archived column: %w", err)
	}
	if archivedExists == 0 {
		_, err = s.db.Exec(`ALTER TABLE nodes ADD COLUMN archived INTEGER NOT NULL DEFAULT 0`)
		if err != nil {
			return fmt.Errorf("failed to add archived column: %w", err)
		}
	}

	// Check if favorite column exists
	var favoriteExists int
	err = s.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('nodes') WHERE name='favorite'").Scan(&favoriteExists)
	if err != nil {
		return fmt.Errorf("failed to check for favorite column: %w", err)
	}
	if favoriteExists == 0 {
		_, err = s.db.Exec(`ALTER TABLE nodes ADD COLUMN favorite INTEGER NOT NULL DEFAULT 0`)
		if err != nil {
			return fmt.Errorf("failed to add favorite column: %w", err)
		}
	}

	return nil
}

const createNodesTable = `
CREATE TABLE IF NOT EXISTS nodes (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	first_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_alive DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status TEXT NOT NULL DEFAULT 'active',
	archived INTEGER NOT NULL DEFAULT 0,
	favorite INTEGER NOT NULL DEFAULT 0,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
`

const createMeasurementsTable = `
CREATE TABLE IF NOT EXISTS measurements (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	node_id TEXT NOT NULL,
	timestamp DATETIME NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	ping_jitter REAL,
	ping_latency REAL,
	ping_low REAL,
	ping_high REAL,
	
	download_bandwidth INTEGER,
	download_bytes INTEGER,
	download_elapsed INTEGER,
	download_latency_iqm REAL,
	download_latency_low REAL,
	download_latency_high REAL,
	download_latency_jitter REAL,
	
	upload_bandwidth INTEGER,
	upload_bytes INTEGER,
	upload_elapsed INTEGER,
	upload_latency_iqm REAL,
	upload_latency_low REAL,
	upload_latency_high REAL,
	upload_latency_jitter REAL,
	
	packet_loss REAL,
	isp TEXT,
	interface_internal_ip TEXT,
	interface_name TEXT,
	interface_mac TEXT,
	interface_is_vpn INTEGER,
	interface_external_ip TEXT,
	
	server_id INTEGER,
	server_host TEXT,
	server_port INTEGER,
	server_name TEXT,
	server_location TEXT,
	server_country TEXT,
	server_ip TEXT,
	
	result_id TEXT,
	result_url TEXT,
	
	UNIQUE(node_id, timestamp),
	FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
)
`

const createFailedMeasurementsTable = `
CREATE TABLE IF NOT EXISTS failed_measurements (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	node_id TEXT NOT NULL,
	timestamp DATETIME NOT NULL,
	error_message TEXT,
	retry_count INTEGER DEFAULT 0,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
)
`

const createAPIKeysTable = `
CREATE TABLE IF NOT EXISTS api_keys (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	key_hash TEXT NOT NULL UNIQUE,
	enabled INTEGER NOT NULL DEFAULT 1,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by TEXT,
	last_used DATETIME,
	revoked_at DATETIME
)
`

// createIndexes creates indexes for all tables
func (s *SQLiteDB) createIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status)",
		"CREATE INDEX IF NOT EXISTS idx_nodes_last_alive ON nodes(last_alive)",
		"CREATE INDEX IF NOT EXISTS idx_nodes_name ON nodes(name)",
		"CREATE INDEX IF NOT EXISTS idx_nodes_archived ON nodes(archived)",
		"CREATE INDEX IF NOT EXISTS idx_nodes_favorite ON nodes(favorite)",
		"CREATE INDEX IF NOT EXISTS idx_measurements_node_id ON measurements(node_id)",
		"CREATE INDEX IF NOT EXISTS idx_measurements_timestamp ON measurements(timestamp)",
		"CREATE INDEX IF NOT EXISTS idx_measurements_created_at ON measurements(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_measurements_node_timestamp ON measurements(node_id, timestamp DESC)",
		"CREATE INDEX IF NOT EXISTS idx_failed_node_id ON failed_measurements(node_id)",
		"CREATE INDEX IF NOT EXISTS idx_failed_timestamp ON failed_measurements(timestamp)",
		"CREATE INDEX IF NOT EXISTS idx_api_keys_enabled ON api_keys(enabled)",
		"CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash)",
	}

	for _, indexSQL := range indexes {
		if _, err := s.db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}
