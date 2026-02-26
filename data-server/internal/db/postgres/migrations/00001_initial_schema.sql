-- +goose Up
CREATE TABLE IF NOT EXISTS nodes (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	first_seen TIMESTAMP NOT NULL DEFAULT NOW(),
	last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
	last_alive TIMESTAMP NOT NULL DEFAULT NOW(),
	status VARCHAR(20) NOT NULL DEFAULT 'active',
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

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
);

CREATE TABLE IF NOT EXISTS failed_measurements (
	id BIGSERIAL PRIMARY KEY,
	node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
	timestamp TIMESTAMP NOT NULL,
	error_message TEXT,
	retry_count INTEGER DEFAULT 0,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS api_keys (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	key_hash VARCHAR(255) NOT NULL UNIQUE,
	enabled BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by VARCHAR(100),
	last_used TIMESTAMP,
	revoked_at TIMESTAMP
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
