# Speedtest Node - Network Measurement Collector

A lightweight network measurement collector that runs Ookla Speedtest periodically, stores results locally in SQLite, and syncs with a central data server.

## 🎯 Features

- **Automatic Testing**: Configurable cron schedule (default: every 10 minutes)
- **Local Storage**: SQLite database for measurements and reliability
- **Batch Sync**: Efficient data transmission (up to 20 measurements per request)
- **Retry Logic**: One retry on test failure, tracking for diagnostics
- **Keepalive Signals**: Regular heartbeat to data server (every 60 seconds)
- **Self-Registration**: Automatic node registration with unique UUID
- **Data Retention**: Configurable local storage (default: 7 days)
- **Graceful Shutdown**: Clean cron job and database closure on SIGINT/SIGTERM
- **Robust Configuration**: POSIX-style flags with environment variable support

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Create `.env` file:**
```bash
NODE_NAME=home-office-node
SERVER_URL=https://speedtest.example.com
API_KEY=your-api-key-here
SPEEDTEST_CRON="*/10 * * * *"
LOG_LEVEL=info
```

2. **Start the node:**
```bash
docker-compose up -d
```

3. **View logs:**
```bash
docker-compose logs -f
```

### Standalone

```bash
# Install Ookla Speedtest CLI
# Docs: https://www.speedtest.net/apps/cli
# Verify: speedtest --version

# Build
go build -o speedtest-node ./cmd/speedtest-node

# Run
./speedtest-node \
  --node-name=home-office \
  --server-url=https://speedtest.example.com \
  --api-key=your-api-key
```

## 📋 Configuration

### Command-Line Flags

```bash
./speedtest-node [flags]

Node Configuration:
  --node-name string         Human-readable node name (default: hostname)
  
Server Configuration:
  --server-url string        Data server URL (HTTPS recommended)
  --api-key string           API key for authentication
  --server-timeout duration  HTTP request timeout (default: 30s)
  --tls-verify               Verify TLS certificates (default: true)
  
Speedtest Configuration:
  --speedtest-cron string    Cron expression for measurements (default: "*/10 * * * *")
  --speedtest-timeout duration Speedtest execution timeout (default: 120s)
  --retry-on-failure         Retry once if speedtest fails (default: true)
  
Sync Configuration:
  --batch-size int           Max measurements per sync request (default: 20)
  --sync-interval duration   Check for unsent data interval (default: 30s)
  --alive-interval duration  Send alive signal interval (default: 60s)
  
Database Configuration:
  --db-path string           SQLite database path (default: "./data/speedtest.db")
  --retention-days int       Keep local data for N days (default: 7)
  
Logging:
  --log-level string         Log level: debug, info, warn, error (default: "info")
  --log-format string        Log format: json or console (default: "json")
  --log-output string        Log file path (default: "./logs/speedtest-node.log")

Other:
  -h, --help                 Show help message
  -v, --version              Print version information
```

### Environment Variables

All flags can be set via environment variables with `SPEEDTEST_` prefix:

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
SPEEDTEST_RETENTION_DAYS=7

# Logging
SPEEDTEST_LOG_LEVEL=info
SPEEDTEST_LOG_FORMAT=json
SPEEDTEST_LOG_OUTPUT=/app/logs/speedtest-node.log
```

**Priority**: Command-line flags > Environment variables > Default values

### Cron Schedule Examples

```bash
"*/10 * * * *"  # Every 10 minutes (default)
"*/15 * * * *"  # Every 15 minutes
"0 * * * *"     # Every hour
"0 */6 * * *"   # Every 6 hours
"0 9,21 * * *"  # Twice daily (9 AM and 9 PM)
```

## 📊 Collected Metrics

Each speedtest measurement includes:

- **Ping**: Jitter, latency, low, high
- **Download**: Bandwidth (bytes/s), total bytes, elapsed time, latency metrics
- **Upload**: Bandwidth (bytes/s), total bytes, elapsed time, latency metrics
- **Network**: Packet loss %, ISP name, internal/external IP, interface info
- **Server**: Test server location, name, country, IP address
- **Result**: Speedtest result ID and URL

## 🗄️ Local Database

SQLite database stores:

### Tables

- **measurements**: Full speedtest results with sync status
- **failed_measurements**: Failed test attempts with error messages
- **config**: Persistent configuration (node UUID)

### Automatic Management

- **Retention**: Old data automatically cleaned up (default: 7 days)
- **Sync Status**: Tracks which measurements have been sent to server
- **Cleanup**: Sent measurements deleted after successful sync

## 🔄 Data Flow

```
┌──────────────┐
│ Cron Trigger │
└──────┬───────┘
       │
       ↓
┌──────────────┐     Success    ┌────────────┐
│   Speedtest  │───────────────►│   SQLite   │
│   Executor   │                 │   Store    │
└──────┬───────┘                 └─────┬──────┘
       │ Failure                        │
       ↓                                │
┌──────────────┐                        │
│   Retry Once │                        │
└──────┬───────┘                        │
       │ Still Failed                   │
       ↓                                │
┌──────────────┐                        │
│ Store Failed │                        │
└──────────────┘                        │
                                        │
       Sync Worker (every 30s)          │
┌──────────────┐◄───────────────────────┘
│ Batch Up to  │
│ 20 Records   │
└──────┬───────┘
       │
       ↓
┌──────────────┐
│ Send to      │
│ Data Server  │
└──────┬───────┘
       │ Success
       ↓
┌──────────────┐
│ Delete from  │
│ Local DB     │
└──────────────┘

Alive Worker (every 60s) ────► Server Keepalive
```

## 🔐 Security

- **API Key**: Bearer token authentication for server communication
- **TLS Verification**: Enforced by default, can be disabled for development
- **Local Data**: SQLite database with restricted file permissions
- **Sensitive Logging**: API keys never logged

## 🐳 Docker Deployment

### Dockerfile Features

- Multi-stage build for minimal image size
- Ookla Speedtest CLI automatically installed
- Alpine Linux base for security and size
- Volume mounts for data and logs persistence
- `network_mode: host` for accurate network measurements

### Volumes

- `./data`: SQLite database persistence
- `./logs`: Application logs

### Network Mode

**Important**: Uses `network_mode: host` to ensure accurate network measurements without Docker NAT overhead.

## 📊 Logging

Structured JSON logs (or console format) include:

```json
{
  "timestamp": "2026-02-14T17:46:10Z",
  "level": "info",
  "message": "Speedtest completed",
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "download_mbps": 93.91,
  "upload_mbps": 93.78,
  "ping_ms": 1.089,
  "jitter_ms": 0.061,
  "packet_loss": 0
}
```

## 🧪 Testing

### Manual Testing

```bash
# Test speedtest CLI
speedtest --format json

# Check database
sqlite3 data/speedtest.db "SELECT * FROM measurements ORDER BY timestamp DESC LIMIT 5;"

# Test server connectivity
curl -H "Authorization: Bearer $API_KEY" $SERVER_URL/api/v1/node/alive
```

### Verify Operation

```bash
# Watch logs for speedtest execution
docker-compose logs -f | grep "Speedtest completed"

# Check sync status
docker-compose logs -f | grep "Measurements synced"

# Verify alive signals
docker-compose logs -f | grep "Alive signal sent"
```

## 🆘 Troubleshooting

**Speedtest not running:**
- Check Speedtest CLI is installed: `speedtest --version`
- Verify cron expression is valid
- Review logs for execution errors
- Check filesystem permissions

**Not syncing with server:**
- Verify server URL is reachable: `curl -I $SERVER_URL`
- Check API key is valid and enabled
- Review network connectivity
- Check TLS certificate verification settings
- Inspect logs for connection errors

**Database growing too large:**
- Verify retention settings (`--retention-days`)
- Check sync is working (measurements should delete after send)
- Manually clean: `DELETE FROM measurements WHERE created_at < datetime('now', '-7 days');`

**High CPU usage:**
- Reduce speedtest frequency (increase cron interval)
- Check for stuck speedtest processes
- Verify timeout settings are appropriate

**Docker network issues:**
- Ensure `network_mode: host` is set in docker-compose.yml
- Check Docker networking configuration
- Verify no firewall blocking outbound connections

## 📈 Performance Considerations

- **Database Size**: ~1,008 records for 7 days at 10-min intervals (~1-2 MB)
- **Memory Usage**: ~20-30 MB typical
- **CPU Usage**: Minimal except during speedtest execution
- **Network Usage**: ~100-500 MB per speedtest (varies by connection speed)
- **Disk I/O**: Minimal, batch operations reduce write load

## 🔧 Build Requirements

- **Go**: 1.24+
- **CGO**: Required for SQLite support
- **GCC**: C compiler for CGO (Linux/macOS) or MinGW (Windows)
- **Ookla Speedtest CLI**: Must be installed and in PATH

```bash
# Build with SQLite support
CGO_ENABLED=1 go build -o speedtest-node ./cmd/speedtest-node

# Verify dependencies
go mod verify
```

## 📚 Additional Documentation

- See `SPECS.md` for detailed technical specifications
- Data server API documentation for endpoint details
- Speedtest CLI documentation: https://www.speedtest.net/apps/cli

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## 📄 License

See LICENSE file in the root directory.

## 🚀 Advanced Configuration

### Custom Speedtest Parameters

The node uses `speedtest --format json` by default. To customize:

1. Modify `internal/speedtest/executor.go`
2. Add flags like `--server-id=<id>` for specific server selection
3. Rebuild and redeploy

### Multiple Nodes on Same Host

Run multiple instances with different configurations:

```bash
# Node 1
./speedtest-node --node-name=node1 --db-path=./data/node1.db --log-output=./logs/node1.log

# Node 2
./speedtest-node --node-name=node2 --db-path=./data/node2.db --log-output=./logs/node2.log
```

### Monitoring Integration

Expose metrics by parsing JSON logs:
- Grep for "Speedtest completed" events
- Extract download_mbps, upload_mbps, ping_ms
- Feed to your monitoring system (Prometheus, Grafana, etc.)
