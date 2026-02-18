package postgres

import (
	"fmt"

	"mark7888/speedtest-data-server/internal/logger"
)

// Migrate runs database migrations for PostgreSQL
func (p *PostgresDB) Migrate() error {
	logger.Log.Info("Running PostgreSQL database migrations")

	// Create nodes table
	_, err := p.db.Exec(createNodesTable)
	if err != nil {
		return fmt.Errorf("failed to create nodes table: %w", err)
	}

	// Create measurements table
	_, err = p.db.Exec(createMeasurementsTable)
	if err != nil {
		return fmt.Errorf("failed to create measurements table: %w", err)
	}

	// Create failed_measurements table
	_, err = p.db.Exec(createFailedMeasurementsTable)
	if err != nil {
		return fmt.Errorf("failed to create failed_measurements table: %w", err)
	}

	// Create api_keys table
	_, err = p.db.Exec(createAPIKeysTable)
	if err != nil {
		return fmt.Errorf("failed to create api_keys table: %w", err)
	}

	// Add missing columns to existing tables
	if err := p.addMissingColumns(); err != nil {
		return fmt.Errorf("failed to add missing columns: %w", err)
	}

	// Create indexes
	if err := p.createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	logger.Log.Info("PostgreSQL migrations completed successfully")
	return nil
}

// addMissingColumns adds any missing columns to existing tables
// This function is called during migration to ensure existing databases
// have all the columns from newer schema versions
func (p *PostgresDB) addMissingColumns() error {
	// Add archived column to nodes table
	_, err := p.db.Exec(`ALTER TABLE nodes ADD COLUMN IF NOT EXISTS archived BOOLEAN NOT NULL DEFAULT false`)
	if err != nil {
		return fmt.Errorf("failed to add archived column: %w", err)
	}

	// Add favorite column to nodes table
	_, err = p.db.Exec(`ALTER TABLE nodes ADD COLUMN IF NOT EXISTS favorite BOOLEAN NOT NULL DEFAULT false`)
	if err != nil {
		return fmt.Errorf("failed to add favorite column: %w", err)
	}

	return nil
}

const createNodesTable = `
CREATE TABLE IF NOT EXISTS nodes (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	first_seen TIMESTAMP NOT NULL DEFAULT NOW(),
	last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
	last_alive TIMESTAMP NOT NULL DEFAULT NOW(),
	status VARCHAR(20) NOT NULL DEFAULT 'active',
	archived BOOLEAN NOT NULL DEFAULT false,
	favorite BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
)
`

const createMeasurementsTable = `
CREATE TABLE IF NOT EXISTS measurements (
	id BIGSERIAL PRIMARY KEY,
	node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
	timestamp TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	
	ping_jitter DOUBLE PRECISION,
	ping_latency DOUBLE PRECISION,
	ping_low DOUBLE PRECISION,
	ping_high DOUBLE PRECISION,
	
	download_bandwidth BIGINT,
	download_bytes BIGINT,
	download_elapsed INTEGER,
	download_latency_iqm DOUBLE PRECISION,
	download_latency_low DOUBLE PRECISION,
	download_latency_high DOUBLE PRECISION,
	download_latency_jitter DOUBLE PRECISION,
	
	upload_bandwidth BIGINT,
	upload_bytes BIGINT,
	upload_elapsed INTEGER,
	upload_latency_iqm DOUBLE PRECISION,
	upload_latency_low DOUBLE PRECISION,
	upload_latency_high DOUBLE PRECISION,
	upload_latency_jitter DOUBLE PRECISION,
	
	packet_loss DOUBLE PRECISION,
	isp VARCHAR(255),
	interface_internal_ip VARCHAR(45),
	interface_name VARCHAR(100),
	interface_mac VARCHAR(17),
	interface_is_vpn BOOLEAN,
	interface_external_ip VARCHAR(45),
	
	server_id INTEGER,
	server_host VARCHAR(255),
	server_port INTEGER,
	server_name VARCHAR(255),
	server_location VARCHAR(255),
	server_country VARCHAR(100),
	server_ip VARCHAR(45),
	
	result_id VARCHAR(255),
	result_url TEXT,
	
	UNIQUE(node_id, timestamp)
)
`

const createFailedMeasurementsTable = `
CREATE TABLE IF NOT EXISTS failed_measurements (
	id BIGSERIAL PRIMARY KEY,
	node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
	timestamp TIMESTAMP NOT NULL,
	error_message TEXT,
	retry_count INTEGER DEFAULT 0,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
)
`

const createAPIKeysTable = `
CREATE TABLE IF NOT EXISTS api_keys (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	key_hash VARCHAR(255) NOT NULL UNIQUE,
	enabled BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by VARCHAR(100),
	last_used TIMESTAMP,
	revoked_at TIMESTAMP
)
`

// createIndexes creates indexes for all tables
func (p *PostgresDB) createIndexes() error {
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
		if _, err := p.db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}
