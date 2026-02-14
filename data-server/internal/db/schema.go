package db

import (
	"fmt"
)

// getDBType detects the database type
func (db *DB) getDBType() (string, error) {
	// Try PostgreSQL specific query
	var version string
	err := db.QueryRow("SELECT version()").Scan(&version)
	if err == nil {
		return "postgres", nil
	}

	// Try SQLite specific query
	err = db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err == nil {
		return "sqlite", nil
	}

	return "", fmt.Errorf("unknown database type")
}

// getNowSQL returns the SQL function for current timestamp based on database type
func (db *DB) getNowSQL() string {
	if db.dbType == "sqlite" {
		return "CURRENT_TIMESTAMP"
	}
	return "NOW()"
}

// getSQLForNodes returns the CREATE TABLE statement for nodes
func getSQLForNodes(dbType string) string {
	if dbType == "sqlite" {
		return `
			CREATE TABLE IF NOT EXISTS nodes (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				first_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				last_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				last_alive DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				status TEXT NOT NULL DEFAULT 'active',
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
			)
		`
	}
	// PostgreSQL
	return `
		CREATE TABLE IF NOT EXISTS nodes (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			first_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			last_alive TIMESTAMP NOT NULL DEFAULT NOW(),
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
}

// getSQLForMeasurements returns the CREATE TABLE statement for measurements
func getSQLForMeasurements(dbType string) string {
	if dbType == "sqlite" {
		return `
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
	}
	// PostgreSQL
	return `
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
}

// getSQLForFailedMeasurements returns the CREATE TABLE statement for failed_measurements
func getSQLForFailedMeasurements(dbType string) string {
	if dbType == "sqlite" {
		return `
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
	}
	// PostgreSQL
	return `
		CREATE TABLE IF NOT EXISTS failed_measurements (
			id BIGSERIAL PRIMARY KEY,
			node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
			timestamp TIMESTAMP NOT NULL,
			error_message TEXT,
			retry_count INTEGER DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
}

// getSQLForAPIKeys returns the CREATE TABLE statement for api_keys
func getSQLForAPIKeys(dbType string) string {
	if dbType == "sqlite" {
		return `
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
	}
	// PostgreSQL
	return `
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
}

// createIndexes creates indexes for all tables
func (db *DB) createIndexes(dbType string) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status)",
		"CREATE INDEX IF NOT EXISTS idx_nodes_last_alive ON nodes(last_alive)",
		"CREATE INDEX IF NOT EXISTS idx_nodes_name ON nodes(name)",
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
		if _, err := db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}
