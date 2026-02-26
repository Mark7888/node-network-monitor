-- +goose Up
CREATE TABLE IF NOT EXISTS nodes (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	first_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_alive DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status TEXT NOT NULL DEFAULT 'active',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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
);

CREATE TABLE IF NOT EXISTS failed_measurements (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	node_id TEXT NOT NULL,
	timestamp DATETIME NOT NULL,
	error_message TEXT,
	retry_count INTEGER DEFAULT 0,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS api_keys (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	key_hash TEXT NOT NULL UNIQUE,
	enabled INTEGER NOT NULL DEFAULT 1,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by TEXT,
	last_used DATETIME,
	revoked_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status);
CREATE INDEX IF NOT EXISTS idx_nodes_last_alive ON nodes(last_alive);
CREATE INDEX IF NOT EXISTS idx_nodes_name ON nodes(name);
CREATE INDEX IF NOT EXISTS idx_measurements_node_id ON measurements(node_id);
CREATE INDEX IF NOT EXISTS idx_measurements_timestamp ON measurements(timestamp);
CREATE INDEX IF NOT EXISTS idx_measurements_created_at ON measurements(created_at);
CREATE INDEX IF NOT EXISTS idx_measurements_node_timestamp ON measurements(node_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_failed_node_id ON failed_measurements(node_id);
CREATE INDEX IF NOT EXISTS idx_failed_timestamp ON failed_measurements(timestamp);
CREATE INDEX IF NOT EXISTS idx_api_keys_enabled ON api_keys(enabled);
CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);

-- +goose Down
DROP TABLE IF EXISTS measurements;
DROP TABLE IF EXISTS failed_measurements;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS nodes;
