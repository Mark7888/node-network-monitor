# Speedtest Node - Network Measurement Collector

## Overview
A lightweight network measurement collector that runs Ookla Speedtest periodically, stores results locally, and syncs with a central data server. Designed to run as a Docker container or standalone application.

## Technologies & Dependencies

### Core
- **Language**: Go 1.24+
- **Database**: SQLite3
- **Container**: Docker
- **External Tool**: Ookla Speedtest CLI

### Go Packages
```
github.com/mattn/go-sqlite3          # SQLite driver
github.com/robfig/cron/v3            # Cron scheduler
github.com/google/uuid               # UUID generation
github.com/spf13/pflag               # POSIX-style command-line flags
github.com/spf13/viper               # Configuration management
go.uber.org/zap                      # Structured logging
```

## Architecture

### Components
1. **Speedtest Executor**: Runs `speedtest --format json` on schedule
2. **Local Storage**: SQLite database for measurements and failures
3. **Sync Manager**: Sends data to server in batches
4. **Alive Signal Sender**: Keepalive heartbeat every 60 seconds
5. **Retry Handler**: Manages failed measurements and unsent data

### Data Flow
```
Cron Trigger → Run Speedtest → Store in SQLite → Queue for Sync
                     ↓ (if fails)
                Retry Once → Store as Failed
                     
Background Worker → Batch Unsent (max 20) → Send to Server → Remove on Success
Background Worker → Send Alive Signal (every 60s)
```

## Database Schema

### Table: `measurements`
```sql
CREATE TABLE measurements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME NOT NULL,           -- UTC timestamp from speedtest
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- Ping data
    ping_jitter REAL,
    ping_latency REAL,
    ping_low REAL,
    ping_high REAL,
    
    -- Download data
    download_bandwidth INTEGER,            -- bytes per second
    download_bytes INTEGER,
    download_elapsed INTEGER,              -- milliseconds
    download_latency_iqm REAL,
    download_latency_low REAL,
    download_latency_high REAL,
    download_latency_jitter REAL,
    
    -- Upload data
    upload_bandwidth INTEGER,              -- bytes per second
    upload_bytes INTEGER,
    upload_elapsed INTEGER,                -- milliseconds
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
    result_id TEXT,                        -- Speedtest result ID
    result_url TEXT,
    
    -- Sync status
    sent BOOLEAN DEFAULT 0,
    sent_at DATETIME,
    
    UNIQUE(timestamp)                      -- Prevent duplicates
);

CREATE INDEX idx_measurements_timestamp ON measurements(timestamp);
CREATE INDEX idx_measurements_sent ON measurements(sent);
CREATE INDEX idx_measurements_created_at ON measurements(created_at);
```

### Table: `failed_measurements`
```sql
CREATE TABLE failed_measurements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME NOT NULL,           -- UTC when test was attempted
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    sent BOOLEAN DEFAULT 0,
    sent_at DATETIME
);

CREATE INDEX idx_failed_timestamp ON failed_measurements(timestamp);
CREATE INDEX idx_failed_sent ON failed_measurements(sent);
```

### Table: `config`
```sql
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Stores: node_id (UUID generated on first run)
```

## Configuration

### Command-Line Interface

The application uses **pflag** and **viper** for robust configuration management with proper POSIX-style flags and automatic help generation. Command-line arguments take precedence over environment variables.

#### Getting Help

```bash
./speedtest-node --help
# or
./speedtest-node -h
```

#### Version Information

```bash
./speedtest-node --version
# or
./speedtest-node -v
```

#### Available Flags

```bash
./speedtest-node [flags]

Flags:
  --node-name string
        Human-readable node name (default: hostname)
        
  --server-url string
        Data server URL (HTTPS recommended)
        
  --api-key string
        API key for authentication
        
  --server-timeout duration
        HTTP request timeout (default: 30s)
        
  --tls-verify
        Verify TLS certificates (default: true)
        
  --speedtest-cron string
        Cron expression for measurements (default: "*/10 * * * *")
        
  --speedtest-timeout duration
        Speedtest execution timeout (default: 120s)
        
  --retry-on-failure
        Retry once if speedtest fails (default: true)
        
  --batch-size int
        Max measurements per sync request (default: 20)
        
  --sync-interval duration
        Check for unsent data interval (default: 30s)
        
  --alive-interval duration
        Send alive signal interval (default: 60s)
        
  --db-path string
        SQLite database path (default: "./data/speedtest.db")
        
  --retention-days int
        Keep local data for N days (default: 7)
        
  --log-level string
        Log level: debug, info, warn, error (default: "info")
        
  --log-format string
        Log format: json or console (default: "json")
        
  --log-output string
        Log file path (default: "./logs/speedtest-node.log")
        
  -h, --help
        Show help message
        
  -v, --version
        Print version information
```

### Environment Variables

All configuration flags can be set via environment variables with the `SPEEDTEST_` prefix. Environment variables use uppercase with underscores, replacing hyphens in flag names.

**Priority:** Command-line flags > Environment variables > Default values

```bash
# Node configuration
SPEEDTEST_NODE_NAME=home-office-node

# Server configuration
SPEEDTEST_SERVER_URL=https://speedtest.example.com
SPEEDTEST_SERVER_API_KEY=your-api-key-here
SPEEDTEST_SERVER_TIMEOUT=30s
SPEEDTEST_TLS_VERIFY=true

# Speedtest configuration
SPEEDTEST_CRON="*/10 * * * *"
SPEEDTEST_TIMEOUT=120s
SPEEDTEST_RETRY_ON_FAILURE=true

# Sync configuration
SPEEDTEST_BATCH_SIZE=20
SPEEDTEST_SYNC_INTERVAL=30s
SPEEDTEST_ALIVE_INTERVAL=60s

# Database configuration
SPEEDTEST_DB_PATH=/data/speedtest.db

# Retention
SPEEDTEST_RETENTION_DAYS=7

# Logging
SPEEDTEST_LOG_LEVEL=info
SPEEDTEST_LOG_FORMAT=json
SPEEDTEST_LOG_OUTPUT=/app/logs/speedtest-node.log
```

### Configuration Examples

**Using command-line flags:**
```bash
./speedtest-node \
  --node-name=home-office \
  --server-url=https://speedtest.example.com \
  --api-key=your-api-key \
  --log-level=debug
```

**Using environment variables:**
```bash
export SPEEDTEST_NODE_NAME=home-office
export SPEEDTEST_SERVER_URL=https://speedtest.example.com
export SPEEDTEST_SERVER_API_KEY=your-api-key
export SPEEDTEST_LOG_LEVEL=debug
./speedtest-node
```

**Mixed (flags override env vars):**
```bash
export SPEEDTEST_NODE_NAME=default-name
./speedtest-node --node-name=override-name  # Uses "override-name"
```

## API Communication

### Endpoints Used on Data Server

#### 1. Register/Alive Signal
**POST** `/api/v1/node/alive`

Headers:
```
Authorization: Bearer <API_KEY>
Content-Type: application/json
```

Request:
```json
{
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "node_name": "home-office-node",
  "timestamp": "2026-02-14T17:46:10Z"
}
```

Response (200):
```json
{
  "status": "ok",
  "server_time": "2026-02-14T17:46:10Z"
}
```

#### 2. Send Measurements
**POST** `/api/v1/measurements`

Headers:
```
Authorization: Bearer <API_KEY>
Content-Type: application/json
```

Request (batch of measurements):
```json
{
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "node_name": "home-office-node",
  "measurements": [
    {
      "timestamp": "2026-02-14T17:46:10Z",
      "ping": {
        "jitter": 0.061,
        "latency": 1.089,
        "low": 1.060,
        "high": 1.198
      },
      "download": {
        "bandwidth": 11739031,
        "bytes": 42291720,
        "elapsed": 3602,
        "latency": {
          "iqm": 9.540,
          "low": 3.194,
          "high": 10.221,
          "jitter": 0.669
        }
      },
      "upload": {
        "bandwidth": 11723027,
        "bytes": 42249312,
        "elapsed": 3604,
        "latency": {
          "iqm": 128.026,
          "low": 5.304,
          "high": 305.608,
          "jitter": 39.456
        }
      },
      "packet_loss": 0,
      "isp": "Eltenet Eltenet",
      "interface": {
        "internal_ip": "172.31.55.120",
        "name": "eth0",
        "mac_addr": "00:15:5D:59:8A:21",
        "is_vpn": false,
        "external_ip": "157.181.192.69"
      },
      "server": {
        "id": 28951,
        "host": "bp-speedtest.zt.hu",
        "port": 8080,
        "name": "ZNET Telekom Zrt.",
        "location": "Budapest",
        "country": "Hungary",
        "ip": "185.232.83.0"
      },
      "result": {
        "id": "be79d55c-6c2d-4a4b-b56b-3f23aa589c0f",
        "url": "https://www.speedtest.net/result/c/be79d55c-6c2d-4a4b-b56b-3f23aa589c0f"
      }
    }
  ]
}
```

Response (200):
```json
{
  "status": "ok",
  "received": 1,
  "failed": 0
}
```

#### 3. Send Failed Measurements
**POST** `/api/v1/measurements/failed`

Headers:
```
Authorization: Bearer <API_KEY>
Content-Type: application/json
```

Request:
```json
{
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "node_name": "home-office-node",
  "failed_tests": [
    {
      "timestamp": "2026-02-14T17:30:00Z",
      "error_message": "speedtest process timeout",
      "retry_count": 1
    }
  ]
}
```

Response (200):
```json
{
  "status": "ok",
  "received": 1
}
```

## Application Structure

```
speedtest-node/
├── cmd/
│   └── speedtest-node/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration from args/env
│   ├── db/
│   │   ├── sqlite.go                  # Database connection
│   │   ├── measurements.go            # Measurement operations
│   │   └── migrations.go              # Database migrations
│   ├── speedtest/
│   │   ├── executor.go                # Speedtest execution
│   │   └── parser.go                  # JSON parsing
│   ├── sync/
│   │   ├── client.go                  # HTTP client for server
│   │   ├── sender.go                  # Batch sender
│   │   └── alive.go                   # Alive signal sender
│   ├── scheduler/
│   │   └── scheduler.go               # Cron scheduling
│   └── logger/
│       └── logger.go                  # Logging setup
├── pkg/
│   └── models/
│       └── measurement.go             # Data models
├── data/                               # SQLite database (gitignored)
├── logs/                               # Log files (gitignored)
├── Dockerfile                          # Docker build
├── docker-compose.yml                 # Docker Compose setup
├── .dockerignore
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Key Implementation Details

### 1. Node ID Generation
On first startup, generate a UUID and store in `config` table:
```go
// Check if node_id exists
nodeID := db.GetConfig("node_id")
if nodeID == "" {
    nodeID = uuid.New().String()
    db.SetConfig("node_id", nodeID)
}
```

### 2. Speedtest Execution with Retry
```go
func RunSpeedtest() (*Measurement, error) {
    result, err := executeSpeedtest()
    if err != nil {
        // Retry once
        time.Sleep(5 * time.Second)
        result, err = executeSpeedtest()
        if err != nil {
            // Store as failed
            db.InsertFailedMeasurement(time.Now(), err.Error(), 1)
            return nil, err
        }
    }
    return result, nil
}
```

### 3. Batch Sync Logic
```go
func SyncMeasurements() {
    unsent := db.GetUnsentMeasurements(20) // Max 20
    if len(unsent) == 0 {
        return
    }
    
    err := client.SendMeasurements(unsent)
    if err != nil {
        logger.Error("Failed to sync", zap.Error(err))
        return
    }
    
    // Delete successfully sent measurements
    db.DeleteMeasurements(unsent)
}
```

### 4. Data Cleanup
```go
func CleanupOldData() {
    retentionDate := time.Now().AddDate(0, 0, -7) // 7 days ago
    db.DeleteMeasurementsBefore(retentionDate)
    db.DeleteFailedMeasurementsBefore(retentionDate)
}
```

### 5. Graceful Shutdown
Handle SIGINT/SIGTERM to stop cron jobs and close DB properly.

## Docker Setup

### Dockerfile
```dockerfile
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o speedtest-node ./cmd/speedtest-node

FROM alpine:latest

# Install speedtest CLI
RUN apk add --no-cache ca-certificates curl && \
    curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.alpine.sh | sh && \
    apk add speedtest

WORKDIR /app
COPY --from=builder /app/speedtest-node .

RUN mkdir -p /app/data /app/logs

VOLUME ["/app/data", "/app/logs"]

CMD ["./speedtest-node"]
```

### docker-compose.yml
```yaml
version: '3.8'

services:
  speedtest-node:
    build: .
    container_name: speedtest-node
    restart: unless-stopped
    environment:
      - NODE_NAME=${NODE_NAME:-home-node}
      - SERVER_URL=${SERVER_URL}
      - API_KEY=${API_KEY}
      - SPEEDTEST_CRON=${CRON:-*/10 * * * *}
      - LOG_LEVEL=${LOG_LEVEL:-info}
      - DB_PATH=/app/data/speedtest.db
      - LOG_OUTPUT=/app/logs/speedtest-node.log
    volumes:
      - ./data:/app/data
      - ./logs:/app/logs
    network_mode: host  # Required for accurate network measurements
```

### .env file
```bash
NODE_NAME=home-office-node
SERVER_URL=https://speedtest.example.com
API_KEY=your-api-key-here
CRON=*/10 * * * *
```

## Running the Application

### Standalone
```bash
# Install dependencies
go mod download

# Build
go build -o speedtest-node ./cmd/speedtest-node

# Run
./speedtest-node
```

### Docker
```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## Logging

All logs are structured JSON (or console format) with fields:
- `timestamp`: ISO8601 timestamp
- `level`: debug, info, warn, error
- `message`: Log message
- `node_id`: Node UUID
- `context`: Additional context fields

Example log entry:
```json
{
  "timestamp": "2026-02-14T17:46:10Z",
  "level": "info",
  "message": "Speedtest completed",
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "download_mbps": 93.91,
  "upload_mbps": 93.78,
  "ping_ms": 1.089
}
```

## Error Handling

### Speedtest Failures
- Timeout: Retry once, then log as failed
- Parse error: Log error and continue
- CLI not found: Fatal error, exit application

### Server Communication Failures
- Connection timeout: Queue for retry
- 401 Unauthorized: Log error, continue (invalid API key)
- 5xx Server Error: Queue for retry
- Network unreachable: Queue for retry

### Database Failures
- Failed to open DB: Fatal error
- Failed to insert: Log error, continue
- Disk full: Log critical error

## Health Monitoring

The application should expose health status via:
- Exit code (for Docker health checks)
- Log messages for critical failures
- Metrics (optional): measurements taken, sync success/failure

## Security Considerations

1. **API Key Storage**: Store in config file with restricted permissions (0600)
2. **TLS Verification**: Always verify server certificates in production
3. **Data Privacy**: Speedtest data includes external IP addresses
4. **File Permissions**: Restrict access to database and logs
5. **Sensitive Logging**: Never log API keys

## Performance Considerations

- **Database Size**: With measurements every 10 minutes for 7 days: ~1,008 records
- **Disk Usage**: ~1-2 MB for SQLite database
- **Memory**: ~20-30 MB typical usage
- **CPU**: Minimal except during speedtest execution
- **Network**: ~100-500 MB per speedtest (varies by connection speed)

## Testing

### Manual Testing
1. Run speedtest manually: `speedtest --format json`
2. Check database: `sqlite3 data/speedtest.db "SELECT * FROM measurements;"`
3. Test server connectivity: `curl -H "Authorization: Bearer $API_KEY" $SERVER_URL/api/v1/node/alive`

### Integration Testing
- Mock speedtest CLI output
- Mock server endpoints
- Test retry logic
- Test data cleanup

## Troubleshooting

**Problem**: Speedtest not running
- Check if `speedtest` CLI is installed: `speedtest --version`
- Check cron expression is valid
- Check logs for errors

**Problem**: Not syncing with server
- Verify server URL is reachable
- Check API key is valid
- Check network connectivity
- Review logs for connection errors

**Problem**: Database growing too large
- Check retention settings
- Manually run cleanup: Delete old records
- Check if sync is working (should delete after send)

## Future Enhancements

- Metrics endpoint (Prometheus format)
- Web UI for local status
- Multiple server support (failover)
- Configurable speedtest parameters (server selection)
- Bandwidth throttling for speedtests
