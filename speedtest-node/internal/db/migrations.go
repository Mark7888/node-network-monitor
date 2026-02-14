package db

const (
	// Schema migrations
	createMeasurementsTable = `
	CREATE TABLE IF NOT EXISTS measurements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		
		-- Ping data
		ping_jitter REAL,
		ping_latency REAL,
		ping_low REAL,
		ping_high REAL,
		
		-- Download data
		download_bandwidth INTEGER,
		download_bytes INTEGER,
		download_elapsed INTEGER,
		download_latency_iqm REAL,
		download_latency_low REAL,
		download_latency_high REAL,
		download_latency_jitter REAL,
		
		-- Upload data
		upload_bandwidth INTEGER,
		upload_bytes INTEGER,
		upload_elapsed INTEGER,
		upload_latency_iqm REAL,
		upload_latency_low REAL,
		upload_latency_high REAL,
		upload_latency_jitter REAL,
		
		-- Network info
		packet_loss REAL,
		isp TEXT,
		interface_internal_ip TEXT,
		interface_name TEXT,
		interface_mac TEXT,
		interface_is_vpn BOOLEAN,
		interface_external_ip TEXT,
		
		-- Server info
		server_id INTEGER,
		server_host TEXT,
		server_port INTEGER,
		server_name TEXT,
		server_location TEXT,
		server_country TEXT,
		server_ip TEXT,
		
		-- Result info
		result_id TEXT,
		result_url TEXT,
		
		-- Sync status
		sent BOOLEAN DEFAULT 0,
		sent_at DATETIME,
		
		UNIQUE(timestamp)
	);`

	createMeasurementsIndexes = `
	CREATE INDEX IF NOT EXISTS idx_measurements_timestamp ON measurements(timestamp);
	CREATE INDEX IF NOT EXISTS idx_measurements_sent ON measurements(sent);
	CREATE INDEX IF NOT EXISTS idx_measurements_created_at ON measurements(created_at);`

	createFailedMeasurementsTable = `
	CREATE TABLE IF NOT EXISTS failed_measurements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		error_message TEXT,
		retry_count INTEGER DEFAULT 0,
		sent BOOLEAN DEFAULT 0,
		sent_at DATETIME
	);`

	createFailedMeasurementsIndexes = `
	CREATE INDEX IF NOT EXISTS idx_failed_timestamp ON failed_measurements(timestamp);
	CREATE INDEX IF NOT EXISTS idx_failed_sent ON failed_measurements(sent);`

	createConfigTable = `
	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);`
)

// runMigrations executes all database migrations
func (db *DB) runMigrations() error {
	migrations := []string{
		createMeasurementsTable,
		createMeasurementsIndexes,
		createFailedMeasurementsTable,
		createFailedMeasurementsIndexes,
		createConfigTable,
	}

	for _, migration := range migrations {
		if _, err := db.conn.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}
