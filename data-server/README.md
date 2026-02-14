# Data Server - Network Measurement Collector

A REST API server that collects network measurement data from multiple speedtest nodes, manages API keys, tracks node status, and provides data access for the frontend dashboard.

## 🎯 Features

- **Node Management**: Self-registration, status tracking (active/unreachable/inactive), keepalive monitoring
- **Data Collection**: Batch measurement ingestion with upsert support, failed test tracking
- **Authentication**: Bcrypt-secured API keys for nodes, JWT tokens for admin access
- **Admin API**: Complete data access for frontend dashboard, aggregated statistics
- **Data Retention**: Configurable cleanup policies (default: 365 days measurements, 90 days failures)
- **Database Support**: PostgreSQL (production) or SQLite (development/testing)
- **Rate Limiting**: Configurable per-key and per-IP rate limits
- **TLS/HTTPS**: Optional TLS support for secure communication

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Create `.env` file:**
```bash
DB_PASSWORD=your-secure-db-password
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-secure-admin-password
JWT_SECRET=your-secret-jwt-key-change-in-production
```

2. **Start the server:**
```bash
docker-compose up -d
```

The server will be available at `http://localhost:8080`

3. **View logs:**
```bash
docker-compose logs -f data-server
```

### Standalone (Development)

#### Option 1: SQLite (Easiest for local development)

```bash
# Build (CGO required for SQLite)
CGO_ENABLED=1 go build -o data-server ./cmd/data-server

# Run
./data-server \
  --db-type=sqlite \
  --db-path=./data/speedtest.db \
  --admin-password=admin123 \
  --jwt-secret=dev-secret-key
```

#### Option 2: PostgreSQL

```bash
# Setup database
createdb speedtest
psql speedtest < migrations/001_initial_schema.up.sql

# Build
go build -o data-server ./cmd/data-server

# Run
./data-server \
  --db-type=postgres \
  --db-password=yourpassword \
  --admin-password=admin123 \
  --jwt-secret=dev-secret-key
```

## 📋 Configuration

### Command-Line Flags

```bash
./data-server [options]

Core Options:
  --host string              Server host address (default: "0.0.0.0")
  --port int                 Server port (default: 8080)
  --mode string              Server mode: debug, release (default: "release")
  
Database:
  --db-type string           Database type: postgres, sqlite (default: "postgres")
  --db-path string           SQLite database file path (default: "./data/speedtest.db")
  --db-host string           PostgreSQL host (default: "localhost")
  --db-port int              PostgreSQL port (default: 5432)
  --db-user string           PostgreSQL user (default: "speedtest")
  --db-password string       PostgreSQL password
  --db-name string           PostgreSQL database name (default: "speedtest")
  --db-sslmode string        PostgreSQL SSL mode (default: "require")
  
Authentication:
  --admin-username string    Admin username (default: "admin")
  --admin-password string    Admin password (required)
  --jwt-secret string        JWT secret key (required)
  --jwt-expiry duration      JWT token expiration (default: 24h)
  
Node Management:
  --alive-timeout duration   Node alive signal timeout (default: 2m)
  --inactive-timeout duration Node inactive timeout (default: 24h)
  
Data Retention:
  --retention-measurements int Keep measurements for N days (default: 365)
  --retention-failed int       Keep failed measurements for N days (default: 90)
  --cleanup-interval duration  Data cleanup interval (default: 24h)
  
Logging:
  --log-level string         Log level: debug, info, warn, error (default: "info")
  --log-format string        Log format: json or console (default: "json")
  --log-output string        Log file path (default: "./logs/data-server.log")
```

### Environment Variables

All flags can be set via environment variables with `SERVER_` or `DB_` prefix:

```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_TYPE=postgres
DB_HOST=postgres
DB_PASSWORD=your-secure-password
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-secure-password
JWT_SECRET=your-secret-key
LOG_LEVEL=info
```

## 📡 API Endpoints

### Node API (Authentication: Bearer API Key)

- **POST** `/api/v1/node/alive` - Node keepalive/registration
- **POST** `/api/v1/measurements` - Submit measurements (batch)
- **POST** `/api/v1/measurements/failed` - Submit failed tests

### Admin API (Authentication: Bearer JWT Token)

- **POST** `/api/v1/admin/login` - Admin authentication
- **POST** `/api/v1/admin/refresh` - Refresh JWT token
- **GET** `/api/v1/admin/nodes` - List all nodes
- **GET** `/api/v1/admin/nodes/:id` - Get node details
- **GET** `/api/v1/admin/nodes/:id/measurements` - Get node measurements
- **GET** `/api/v1/admin/measurements/aggregate` - Aggregated data for charts
- **GET** `/api/v1/admin/api-keys` - List API keys
- **POST** `/api/v1/admin/api-keys` - Create API key
- **PATCH** `/api/v1/admin/api-keys/:id` - Update API key
- **DELETE** `/api/v1/admin/api-keys/:id` - Delete API key
- **GET** `/api/v1/admin/dashboard` - Dashboard summary

### Health Check

- **GET** `/health` - Server health status

## 🔐 Security

- **API Keys**: Bcrypt-hashed, 32+ byte cryptographically secure random
- **Admin Password**: Bcrypt-hashed with strong policy enforcement
- **JWT Tokens**: HS256 signing, configurable expiry (default 24h)
- **Rate Limiting**: 100 req/min per API key, 20 req/min per IP for admin
- **TLS Support**: Optional HTTPS with certificate verification
- **Database**: PostgreSQL SSL mode, connection pooling with limits

## 🗄️ Database Schema

### Tables

- **nodes**: Node registration and status tracking
- **measurements**: Speedtest results with full metrics
- **failed_measurements**: Failed test attempts and errors
- **api_keys**: API key management with usage tracking
- **admin_sessions**: Optional session tracking

### Indexes

Optimized for common queries:
- Node status and last_alive lookups
- Measurement queries by node_id and timestamp
- API key authentication

## 🏗️ Architecture

```
┌─────────────────┐
│ Speedtest Nodes │
└────────┬────────┘
         │ POST /api/v1/measurements
         │ POST /api/v1/node/alive
         ↓
┌─────────────────┐         ┌──────────────┐
│   Data Server   │◄────────┤  PostgreSQL  │
│  (Gin Router)   │         │   or SQLite  │
└────────┬────────┘         └──────────────┘
         │
         │ JWT Auth
         ↓
┌─────────────────┐
│  Frontend       │
│  Dashboard      │
└─────────────────┘
```

## 🔧 Build Requirements

- **Go**: 1.24+
- **CGO**: Required for SQLite support only
  - Linux/macOS: GCC or Clang (usually pre-installed)
  - Windows: MinGW or TDM-GCC
  - Set `CGO_ENABLED=1` when building with SQLite
- **PostgreSQL**: 15+ (production) or SQLite 3+ (development)

**Build without SQLite** (PostgreSQL only):
```bash
CGO_ENABLED=0 go build -o data-server ./cmd/data-server
```

## 📊 Monitoring

### Structured Logging

JSON logs with contextual information:
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

### Background Services

- **Node Status Tracker**: Monitors alive signals, updates node status (every 30s)
- **Data Cleanup**: Removes old measurements per retention policy (every 24h)

## 🐳 Docker Deployment

The `docker-compose.yml` includes:
- PostgreSQL 15 with health checks
- Data server with automatic migrations
- Volume persistence for database and logs
- Environment-based configuration

**Ports:**
- `8080`: API server
- `5432`: PostgreSQL (exposed for debugging)

**Volumes:**
- `postgres_data`: Database persistence
- `./logs`: Application logs

## 🧪 Testing

```bash
# Run tests
go test ./...

# Integration tests with test database
go test -tags=integration ./...

# Check health endpoint
curl http://localhost:8080/health
```

## 📚 Additional Documentation

- See `SPECS.md` for detailed technical specifications
- API documentation available at `/api/v1/docs` (when enabled)
- Database schema migrations in `migrations/` directory

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## 📄 License

See LICENSE file in the root directory.

## 🆘 Troubleshooting

**PostgreSQL connection failed:**
- Verify PostgreSQL is running: `pg_isready`
- Check credentials and database exists
- Review SSL/TLS settings

**SQLite build errors:**
- Ensure `CGO_ENABLED=1`
- Install C compiler (gcc, clang, or MinGW)
- Run `go mod download`

**API authentication failing:**
- Verify Bearer token format
- Check API key is enabled
- Review rate limit settings
