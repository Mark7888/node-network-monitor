# Data Server - Central Network Measurement Collector

## Overview
A REST API server that collects network measurement data from multiple speedtest nodes, manages API keys, tracks node status, and provides data access for the frontend dashboard. Built to run on a VPS with PostgreSQL for reliable data storage.

## Technologies & Dependencies

### Core
- **Language**: Go 1.24+
- **Database**: PostgreSQL 15+ (production) or SQLite 3+ (development/testing)
- **Container**: Docker & Docker Compose
- **Web Framework**: Gin (HTTP router)

### Go Packages
```
github.com/gin-gonic/gin              # HTTP web framework
github.com/lib/pq                     # PostgreSQL driver
github.com/mattn/go-sqlite3           # SQLite driver (CGO required)
github.com/google/uuid                # UUID generation
go.uber.org/zap                       # Structured logging
golang.org/x/crypto/bcrypt            # Password & API key hashing
github.com/golang-jwt/jwt/v5          # JWT tokens for frontend auth
github.com/gin-contrib/cors           # CORS middleware
```

### Build Requirements
- **CGO**: Required for SQLite support. Ensure `CGO_ENABLED=1` when building with SQLite.
- **GCC/Clang**: C compiler required for CGO (typically pre-installed on Linux/macOS, use MinGW on Windows)
- To build without SQLite support: `CGO_ENABLED=0 go build` (PostgreSQL only)

## Architecture

### Components
1. **Node API**: Receives measurements and alive signals from nodes
2. **Admin API**: Manages API keys and provides data access for frontend
3. **Authentication**: JWT-based auth for frontend
4. **Node Status Tracker**: Monitors node health based on alive signals
5. **Data Storage**: PostgreSQL with normalized tables
6. **Cleanup Service**: Removes old data based on retention policy

### Data Flow
```
Speedtest Nodes → [Node API] → Database (PostgreSQL/SQLite)
                       ↓
                   Alive Tracker → Node Status Updates
                   
Frontend → [Auth] → JWT Token → [Admin API] → Database (PostgreSQL/SQLite)
```

### Database Backend Options

**PostgreSQL (Production)**
- Recommended for production deployments
- Better performance with large datasets
- Full-featured with advanced SQL capabilities
- Concurrent write operations

**SQLite (Development/Testing)**
- Perfect for local development without Docker
- Ideal for testing and CI/CD pipelines
- Single file database, easy backup
- No separate database server required
- Use `--db-type=sqlite` to enable

## Database Schema

### Table: `nodes`
```sql
CREATE TABLE nodes (
    id UUID PRIMARY KEY,                          -- Node UUID (from client)
    name VARCHAR(255) NOT NULL,                   -- Human-readable name
    first_seen TIMESTAMP NOT NULL DEFAULT NOW(),  -- UTC timestamp
    last_seen TIMESTAMP NOT NULL DEFAULT NOW(),   -- UTC timestamp  
    last_alive TIMESTAMP NOT NULL DEFAULT NOW(),  -- UTC timestamp
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, unreachable, inactive
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_nodes_status ON nodes(status);
CREATE INDEX idx_nodes_last_alive ON nodes(last_alive);
CREATE INDEX idx_nodes_name ON nodes(name);
```

### Table: `measurements`
```sql
CREATE TABLE measurements (
    id BIGSERIAL PRIMARY KEY,
    node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    timestamp TIMESTAMP NOT NULL,                 -- UTC from speedtest
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),  -- UTC when received
    
    -- Ping metrics
    ping_jitter DOUBLE PRECISION,
    ping_latency DOUBLE PRECISION,
    ping_low DOUBLE PRECISION,
    ping_high DOUBLE PRECISION,
    
    -- Download metrics
    download_bandwidth BIGINT,                    -- bytes per second
    download_bytes BIGINT,
    download_elapsed INTEGER,                     -- milliseconds
    download_latency_iqm DOUBLE PRECISION,
    download_latency_low DOUBLE PRECISION,
    download_latency_high DOUBLE PRECISION,
    download_latency_jitter DOUBLE PRECISION,
    
    -- Upload metrics
    upload_bandwidth BIGINT,                      -- bytes per second
    upload_bytes BIGINT,
    upload_elapsed INTEGER,                       -- milliseconds
    upload_latency_iqm DOUBLE PRECISION,
    upload_latency_low DOUBLE PRECISION,
    upload_latency_high DOUBLE PRECISION,
    upload_latency_jitter DOUBLE PRECISION,
    
    -- Network info
    packet_loss DOUBLE PRECISION,
    isp VARCHAR(255),
    interface_internal_ip VARCHAR(45),
    interface_name VARCHAR(100),
    interface_mac VARCHAR(17),
    interface_is_vpn BOOLEAN,
    interface_external_ip VARCHAR(45),
    
    -- Server info
    server_id INTEGER,
    server_host VARCHAR(255),
    server_port INTEGER,
    server_name VARCHAR(255),
    server_location VARCHAR(255),
    server_country VARCHAR(100),
    server_ip VARCHAR(45),
    
    -- Result info
    result_id VARCHAR(255),
    result_url TEXT,
    
    UNIQUE(node_id, timestamp)                    -- Prevent duplicates
);

CREATE INDEX idx_measurements_node_id ON measurements(node_id);
CREATE INDEX idx_measurements_timestamp ON measurements(timestamp);
CREATE INDEX idx_measurements_created_at ON measurements(created_at);
CREATE INDEX idx_measurements_node_timestamp ON measurements(node_id, timestamp DESC);
```

### Table: `failed_measurements`
```sql
CREATE TABLE failed_measurements (
    id BIGSERIAL PRIMARY KEY,
    node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    timestamp TIMESTAMP NOT NULL,                 -- UTC when test failed
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_failed_node_id ON failed_measurements(node_id);
CREATE INDEX idx_failed_timestamp ON failed_measurements(timestamp);
```

### Table: `api_keys`
```sql
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,                   -- Human-readable name
    key_hash VARCHAR(255) NOT NULL UNIQUE,        -- bcrypt hash
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by VARCHAR(100),                      -- Admin username
    last_used TIMESTAMP,
    revoked_at TIMESTAMP,
    
    UNIQUE(key_hash)
);

CREATE INDEX idx_api_keys_enabled ON api_keys(enabled);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
```

### Table: `admin_sessions` (optional, for session tracking)
```sql
CREATE TABLE admin_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_activity TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_token ON admin_sessions(token_hash);
CREATE INDEX idx_sessions_expires ON admin_sessions(expires_at);
```

## Configuration

### Command-Line Arguments

All configuration is done via command-line arguments or environment variables. Arguments take precedence over environment variables.

```bash
./data-server [options]

Options:
  --host string
        Server host address (default: "0.0.0.0")
        
  --port int
        Server port (default: 8080)
        
  --mode string
        Server mode: debug, release (default: "release")
        
  --tls-enabled
        Enable HTTPS/TLS (default: false)
        
  --tls-cert string
        Path to TLS certificate file
        
  --tls-key string
        Path to TLS key file
        
  --db-type string
        Database type: postgres, sqlite (default: "postgres")
        
  --db-path string
        SQLite database file path (default: "./data/speedtest.db")
        
  --db-host string
        PostgreSQL host (default: "localhost")
        
  --db-port int
        PostgreSQL port (default: 5432)
        
  --db-user string
        PostgreSQL user (default: "speedtest")
        
  --db-password string
        PostgreSQL password
        
  --db-name string
        PostgreSQL database name (default: "speedtest")
        
  --db-sslmode string
        PostgreSQL SSL mode: disable, require, verify-full (default: "require")
        
  --db-max-connections int
        Max database connections (default: 25)
        
  --db-max-idle int
        Max idle connections (default: 5)
        
  --db-connection-lifetime duration
        Connection lifetime (default: 5m)
        
  --admin-username string
        Admin username (default: "admin")
        
  --admin-password string
        Admin password (required)
        
  --jwt-secret string
        JWT secret key (required)
        
  --jwt-expiry duration
        JWT token expiration (default: 24h)
        
  --alive-timeout duration
        Node alive signal timeout (default: 2m)
        
  --inactive-timeout duration
        Node inactive timeout (default: 24h)
        
  --status-check-interval duration
        Node status check interval (default: 30s)
        
  --retention-measurements int
        Keep measurements for N days (default: 365)
        
  --retention-failed int
        Keep failed measurements for N days (default: 90)
        
  --cleanup-interval duration
        Data cleanup interval (default: 24h)
        
  --rate-limit int
        Requests per minute per API key (default: 100)
        
  --timeout duration
        API request timeout (default: 30s)
        
  --log-level string
        Log level: debug, info, warn, error (default: "info")
        
  --log-format string
        Log format: json or console (default: "json")
        
  --log-output string
        Log file path (default: "./logs/data-server.log")
```

### Environment Variables

Environment variables can be used instead of arguments. Arguments override environment variables if both are present.
```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database Configuration
DB_TYPE=postgres  # or sqlite

# PostgreSQL (when DB_TYPE=postgres)
DB_HOST=postgres
DB_PORT=5432
DB_USER=speedtest
DB_PASSWORD=your-secure-password
DB_NAME=speedtest
DB_SSLMODE=require

# SQLite (when DB_TYPE=sqlite)
DB_PATH=./data/speedtest.db

# Authentication
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-secure-password
JWT_SECRET=your-secret-key

# Logging
LOG_LEVEL=info
```

## API Endpoints

### Node API (Used by Speedtest Nodes)

All node endpoints require `Authorization: Bearer <API_KEY>` header.

#### 1. Node Alive Signal / Registration
**POST** `/api/v1/node/alive`

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
  "server_time": "2026-02-14T17:46:10Z",
  "node_registered": true
}
```

Logic:
- If node doesn't exist, create new entry (self-registration)
- Update `last_alive` and `status` to `active`
- Update `last_seen` timestamp

#### 2. Submit Measurements
**POST** `/api/v1/measurements`

Request:
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
      "isp": "Eltenet",
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
        "name": "ZNET Telekom",
        "location": "Budapest",
        "country": "Hungary",
        "ip": "185.232.83.0"
      },
      "result": {
        "id": "be79d55c-6c2d-4a4b-b56b-3f23aa589c0f",
        "url": "https://www.speedtest.net/result/c/..."
      }
    }
  ]
}
```

Response (200):
```json
{
  "status": "ok",
  "received": 20,
  "inserted": 18,
  "updated": 2,
  "failed": 0
}
```

Response (400):
```json
{
  "error": "Invalid measurement data",
  "details": "Missing required field: timestamp"
}
```

Logic:
- Validate node exists (or register if needed)
- Parse each measurement
- Insert or update (UPSERT on duplicate timestamp+node_id)
- Return count of successful operations

#### 3. Submit Failed Measurements
**POST** `/api/v1/measurements/failed`

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

### Admin API (Used by Frontend)

All admin endpoints require `Authorization: Bearer <JWT_TOKEN>` header (except login).

#### 1. Admin Login
**POST** `/api/v1/admin/login`

Request:
```json
{
  "username": "admin",
  "password": "your-password"
}
```

Response (200):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2026-02-15T17:46:10Z",
  "username": "admin"
}
```

Response (401):
```json
{
  "error": "Invalid credentials"
}
```

#### 2. Refresh Token
**POST** `/api/v1/admin/refresh`

Headers:
```
Authorization: Bearer <JWT_TOKEN>
```

Response (200):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2026-02-15T17:46:10Z"
}
```

#### 3. List All Nodes
**GET** `/api/v1/admin/nodes`

Query params:
- `status` (optional): Filter by status (active, unreachable, inactive)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)

Response (200):
```json
{
  "nodes": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "home-office-node",
      "status": "active",
      "first_seen": "2026-01-01T10:00:00Z",
      "last_seen": "2026-02-14T17:46:10Z",
      "last_alive": "2026-02-14T17:46:00Z",
      "measurement_count": 4320,
      "latest_measurement": {
        "timestamp": "2026-02-14T17:40:00Z",
        "download_mbps": 93.91,
        "upload_mbps": 93.78,
        "ping_ms": 1.089
      }
    }
  ],
  "total": 5,
  "page": 1,
  "limit": 50
}
```

#### 4. Get Node Details
**GET** `/api/v1/admin/nodes/:id`

Response (200):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "home-office-node",
  "status": "active",
  "first_seen": "2026-01-01T10:00:00Z",
  "last_seen": "2026-02-14T17:46:10Z",
  "last_alive": "2026-02-14T17:46:00Z",
  "measurement_count": 4320,
  "failed_test_count": 12,
  "statistics": {
    "avg_download_mbps": 95.2,
    "avg_upload_mbps": 89.7,
    "avg_ping_ms": 1.2,
    "avg_jitter_ms": 0.08,
    "avg_packet_loss": 0.01
  }
}
```

#### 5. Get Node Measurements
**GET** `/api/v1/admin/nodes/:id/measurements`

Query params:
- `from` (optional): Start timestamp (ISO8601)
- `to` (optional): End timestamp (ISO8601)
- `limit` (optional): Max results (default: 1000)
- `page` (optional): Page number (default: 1)

Response (200):
```json
{
  "measurements": [
    {
      "id": 12345,
      "timestamp": "2026-02-14T17:40:00Z",
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
      "isp": "Eltenet",
      "interface": {
        "internal_ip": "172.31.55.120",
        "external_ip": "157.181.192.69"
      },
      "server": {
        "name": "ZNET Telekom",
        "location": "Budapest",
        "country": "Hungary"
      }
    }
  ],
  "total": 4320,
  "page": 1,
  "limit": 1000
}
```

#### 6. Get Aggregated Measurements (for Charts)
**GET** `/api/v1/admin/measurements/aggregate`

Query params:
- `node_ids` (optional): Comma-separated node IDs (all if omitted)
- `from` (required): Start timestamp (ISO8601)
- `to` (required): End timestamp (ISO8601)
- `interval` (required): Aggregation interval (5m, 15m, 1h, 6h, 1d)

Response (200):
```json
{
  "data": [
    {
      "timestamp": "2026-02-14T17:00:00Z",
      "node_id": "550e8400-e29b-41d4-a716-446655440000",
      "node_name": "home-office-node",
      "avg_download_mbps": 95.2,
      "avg_upload_mbps": 89.7,
      "avg_ping_ms": 1.2,
      "avg_jitter_ms": 0.08,
      "avg_packet_loss": 0.01,
      "min_download_mbps": 85.0,
      "max_download_mbps": 99.5,
      "sample_count": 6
    }
  ],
  "interval": "1h",
  "total_samples": 144
}
```

#### 7. List API Keys
**GET** `/api/v1/admin/api-keys`

Response (200):
```json
{
  "api_keys": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Production Node 1",
      "enabled": true,
      "created_at": "2026-01-01T10:00:00Z",
      "created_by": "admin",
      "last_used": "2026-02-14T17:46:00Z"
    }
  ],
  "total": 5
}
```

#### 8. Create API Key
**POST** `/api/v1/admin/api-keys`

Request:
```json
{
  "name": "New Measurement Node"
}
```

Response (201):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "New Measurement Node",
  "key": "sk_live_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "enabled": true,
  "created_at": "2026-02-14T17:46:10Z",
  "warning": "Save this key securely. It won't be shown again."
}
```

Note: The plain API key is only returned once during creation.

#### 9. Update API Key
**PATCH** `/api/v1/admin/api-keys/:id`

Request:
```json
{
  "enabled": false
}
```

Response (200):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Production Node 1",
  "enabled": false,
  "updated_at": "2026-02-14T17:46:10Z"
}
```

#### 10. Delete API Key
**DELETE** `/api/v1/admin/api-keys/:id`

Response (204): No content

#### 11. Get Dashboard Summary
**GET** `/api/v1/admin/dashboard`

Response (200):
```json
{
  "total_nodes": 5,
  "active_nodes": 4,
  "unreachable_nodes": 1,
  "total_measurements": 21600,
  "measurements_last_24h": 576,
  "last_measurement": "2026-02-14T17:46:00Z",
  "average_stats_24h": {
    "download_mbps": 92.5,
    "upload_mbps": 88.3,
    "ping_ms": 1.3,
    "jitter_ms": 0.09,
    "packet_loss": 0.02
  }
}
```

## Application Structure

```
data-server/
├── cmd/
│   └── data-server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration from args/env
│   ├── db/
│   │   ├── postgres.go                # Database connection
│   │   ├── migrations.go              # Database migrations
│   │   ├── nodes.go                   # Node operations
│   │   ├── measurements.go            # Measurement operations
│   │   └── api_keys.go                # API key operations
│   ├── api/
│   │   ├── router.go                  # API routes setup
│   │   ├── middleware/
│   │   │   ├── auth.go                # Authentication middleware
│   │   │   ├── cors.go                # CORS middleware
│   │   │   └── rate_limit.go          # Rate limiting
│   │   ├── handlers/
│   │   │   ├── node.go                # Node endpoints
│   │   │   ├── measurement.go         # Measurement endpoints
│   │   │   ├── admin.go               # Admin endpoints
│   │   │   └── api_keys.go            # API key management
│   │   └── validators/
│   │       └── validators.go          # Input validation
│   ├── auth/
│   │   ├── jwt.go                     # JWT token handling
│   │   └── password.go                # Password hashing
│   ├── services/
│   │   ├── node_tracker.go            # Node status monitoring
│   │   ├── cleanup.go                 # Data retention cleanup
│   │   └── aggregation.go             # Data aggregation for charts
│   └── logger/
│       └── logger.go                  # Logging setup
├── pkg/
│   └── models/
│       ├── node.go                    # Node models
│       ├── measurement.go             # Measurement models
│       ├── api_key.go                 # API key models
│       └── auth.go                    # Auth models
├── migrations/                         # SQL migration files
│   ├── 001_initial_schema.up.sql
│   └── 001_initial_schema.down.sql
├── Dockerfile                          # Docker build
├── docker-compose.yml                 # Docker Compose with PostgreSQL
├── .dockerignore
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Key Implementation Details

### 1. API Key Authentication
```go
func AuthenticateAPIKey(keyString string) (*APIKey, error) {
    // Extract key from "Bearer <key>" header
    // Query all enabled keys
    // Compare using bcrypt
    // Update last_used timestamp
    // Return key details if valid
}
```

### 2. JWT Authentication
```go
func GenerateJWT(username string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
        "iat": time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}
```

### 3. Node Status Tracking
```go
// Background goroutine
func MonitorNodeStatus() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        // Find nodes with last_alive > 2 minutes ago
        // Update status to "unreachable"
        
        // Find nodes with last_alive > 24 hours ago
        // Update status to "inactive"
    }
}
```

### 4. Measurement Upsert
```go
// Handle duplicate timestamps per node
INSERT INTO measurements (node_id, timestamp, ...)
VALUES ($1, $2, ...)
ON CONFLICT (node_id, timestamp)
DO UPDATE SET
    ping_jitter = EXCLUDED.ping_jitter,
    ...
    updated_at = NOW()
```

### 5. Data Cleanup
```go
func CleanupOldData() {
    retentionDate := time.Now().AddDate(0, 0, -365)
    db.Exec("DELETE FROM measurements WHERE timestamp < $1", retentionDate)
    
    failedRetentionDate := time.Now().AddDate(0, 0, -90)
    db.Exec("DELETE FROM failed_measurements WHERE timestamp < $1", failedRetentionDate)
}
```

### 6. Data Aggregation
```go
// For charts - aggregate by time interval
SELECT 
    date_trunc('hour', timestamp) as time_bucket,
    node_id,
    AVG(download_bandwidth) as avg_download,
    AVG(upload_bandwidth) as avg_upload,
    AVG(ping_latency) as avg_ping,
    COUNT(*) as sample_count
FROM measurements
WHERE timestamp BETWEEN $1 AND $2
GROUP BY time_bucket, node_id
ORDER BY time_bucket
```

## Docker Setup

### Dockerfile
```dockerfile
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o data-server ./cmd/data-server

FROM alpine:latest

RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app
COPY --from=builder /app/data-server .
COPY migrations ./migrations

RUN mkdir -p /app/logs

VOLUME ["/app/logs"]

EXPOSE 8080

CMD ["./data-server"]
```

### docker-compose.yml
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: speedtest-db
    restart: unless-stopped
    environment:
      POSTGRES_DB: speedtest
      POSTGRES_USER: speedtest
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U speedtest"]
      interval: 10s
      timeout: 5s
      retries: 5

  data-server:
    build: .
    container_name: speedtest-server
    restart: unless-stopped
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - HOST=0.0.0.0
      - PORT=8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=speedtest
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=speedtest
      - DB_SSLMODE=disable
      - ADMIN_USERNAME=${ADMIN_USERNAME}
      - ADMIN_PASSWORD=${ADMIN_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
      - LOG_LEVEL=info
      - LOG_OUTPUT=/app/logs/data-server.log
    ports:
      - "8080:8080"
    volumes:
      - ./logs:/app/logs

volumes:
  postgres_data:
```

### .env file
```bash
DB_PASSWORD=your-secure-db-password
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-secure-admin-password
JWT_SECRET=your-secret-jwt-key-change-in-production
```

## Database Migrations

Use a migration tool or manual SQL scripts:

### Initial Schema Migration
```sql
-- migrations/001_initial_schema.up.sql
-- (Include all CREATE TABLE statements from schema section)
```

### Running Migrations
```bash
# Manual approach
psql -h localhost -U speedtest -d speedtest -f migrations/001_initial_schema.up.sql

# Or use golang-migrate
migrate -path migrations -database "postgres://user:pass@localhost:5432/speedtest?sslmode=disable" up
```

## Running the Application

### Standalone (Development)

**With SQLite (Easiest for local development):**
```bash
# Install dependencies
go mod download

# Build (CGO required for SQLite)
go build -o data-server ./cmd/data-server

# Run with SQLite
./data-server --db-type=sqlite --db-path=./data/speedtest.db \
  --admin-password=admin123 --jwt-secret=dev-secret-key

# Or use environment variables
export DB_TYPE=sqlite
export DB_PATH=./data/speedtest.db
export ADMIN_PASSWORD=admin123
export JWT_SECRET=dev-secret-key
./data-server
```

**With PostgreSQL:**
```bash
# Install dependencies
go mod download

# Setup PostgreSQL
createdb speedtest
psql speedtest < migrations/001_initial_schema.up.sql

# Build
go build -o data-server ./cmd/data-server

# Run
./data-server --db-type=postgres --db-password=yourpassword \
  --admin-password=admin123 --jwt-secret=dev-secret-key
```

### Docker (Production)
```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f data-server

# Stop
docker-compose down

# Stop and remove data
docker-compose down -v
```

## Security Considerations

### API Key Security
- Store hashed with bcrypt (cost 12+)
- Generate with cryptographically secure random (32+ bytes)
- Use prefix (e.g., `sk_live_`) for identification
- Never log full keys

### Admin Password
- Hash with bcrypt
- Enforce strong password policy (min length, complexity)
- Allow password change via config reload

### JWT Tokens
- Use strong secret (256+ bits)
- Short expiry (24 hours)
- Include refresh mechanism
- Validate on every protected endpoint

### Database
- Use connection pooling
- Enable SSL/TLS in production
- Restrict network access
- Regular backups

### TLS/HTTPS
- Use valid certificates (Let's Encrypt)
- TLS 1.2+ only
- Strong cipher suites

### Rate Limiting
- Per API key: 100 req/min
- Per IP for admin: 20 req/min
- Exponential backoff on failures

## Performance Optimization

### Database Indexes
- Index all foreign keys
- Composite index on (node_id, timestamp) for queries
- Index on status, last_alive for monitoring

### Connection Pooling
- Max 25 connections
- Idle connection timeout: 5 minutes
- Connection lifetime: 1 hour

### Query Optimization
- Use prepared statements
- Batch inserts for measurements
- Use aggregation queries for charts
- Add LIMIT to prevent large result sets

### Caching (Future)
- Cache dashboard summary (1 minute TTL)
- Cache node list (30 seconds TTL)
- Use Redis for session storage

## Monitoring & Logging

### Structured Logging
```json
{
  "timestamp": "2026-02-14T17:46:10Z",
  "level": "info",
  "message": "Measurement received",
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "measurement_count": 20,
  "duration_ms": 45
}
```

### Metrics to Track
- API request count by endpoint
- Request duration
- Database query duration
- Active connections
- Error rates
- Node status distribution

### Health Endpoint
**GET** `/health`
```json
{
  "status": "healthy",
  "database": "connected",
  "uptime_seconds": 86400,
  "version": "1.0.0"
}
```

## Testing

### Unit Tests
- Test authentication logic
- Test validation functions
- Test data aggregation queries

### Integration Tests
- Test full API endpoints with test database
- Mock speedtest node requests
- Test JWT flow

### Load Testing
- Simulate multiple nodes sending data
- Test batch insert performance
- Measure query response times

## Troubleshooting

**Problem**: Database connection failed (PostgreSQL)
- Check PostgreSQL is running
- Verify credentials
- Check network connectivity
- Review SSL/TLS settings

**Problem**: SQLite build errors (`undefined: sqlite3.SQLiteDriver`)
- Ensure `CGO_ENABLED=1` when building: `CGO_ENABLED=1 go build`
- Install GCC/Clang compiler (Linux/macOS) or MinGW (Windows)
- Run `go mod download` to ensure dependencies are fetched

**Problem**: SQLite database locked
- SQLite uses single-writer concurrency
- Ensure `--db-max-connections=1` when using SQLite (default behavior)
- Check no other process has the database file open

**Problem**: SQLite permission denied
- Ensure the directory for database file exists and is writable
- Default: `./data/speedtest.db` requires `./data/` directory
- Check file permissions on database file

**Problem**: API key authentication failing
- Verify key format (Bearer token)
- Check key is enabled
- Verify key hash in database

**Problem**: JWT token invalid
- Check token expiry
- Verify JWT secret matches
- Check token format

**Problem**: Nodes showing as unreachable
- Check alive signal frequency
- Review node_tracker service logs
- Verify timeout settings

## Future Enhancements

- Multi-tenancy support (multiple organizations)
- Webhooks for alerts (node down, poor performance)
- Prometheus metrics export
- Grafana integration
- Email notifications
- Data export (CSV, JSON)
- Advanced filtering and search
- Performance trends and anomaly detection
