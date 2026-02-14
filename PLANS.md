# Network Speed Measurement System

A distributed network monitoring system that measures internet speeds from multiple nodes, aggregates data in a central server, and visualizes metrics through an interactive dashboard.

## 🎯 Project Goals

- Monitor network performance across multiple locations
- Store historical speedtest data for trend analysis
- Visualize download/upload speeds, latency, jitter, and packet loss
- Self-hosting friendly with Docker support
- Lightweight and efficient for long-term deployment

## 📊 Features

### Multi-Node Monitoring
- Deploy measurement nodes at different locations
- Automatic node registration with unique IDs
- Real-time node status tracking (active/unreachable/inactive)
- Independent operation with local data storage

### Comprehensive Metrics
- **Download/Upload Speed**: Bandwidth measurements in Mbps
- **Ping Latency**: Round-trip time to test servers
- **Jitter**: Latency variation over time
- **Packet Loss**: Percentage of lost packets
- **ISP Information**: Network provider details
- **Server Details**: Test server location and specs

### Data Management
- Configurable measurement frequency (default: every 10 minutes)
- Local data retention: 7 days on nodes
- Server data retention: 1 year (configurable)
- Automatic cleanup of old data
- Batch data transmission (up to 20 measurements)
- Retry mechanism for failed transmissions

### Secure API
- Bcrypt-hashed API keys for node authentication
- JWT-based authentication for admin access
- HTTPS/TLS support
- Per-key usage tracking

### Interactive Dashboard
- Real-time charts with 10-second auto-refresh
- Time range filters (day/week/month)
- Per-node detailed views
- Multi-node comparison charts
- API key management interface
- Summary statistics and analytics

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Internet / Network                       │
└─────────────────────────────────────────────────────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
    ┌───────────┐      ┌───────────┐      ┌───────────┐
    │   Node 1  │      │   Node 2  │      │   Node 3  │
    │           │      │           │      │           │
    │  Speedtest│      │  Speedtest│      │  Speedtest│
    │  + SQLite │      │  + SQLite │      │  + SQLite │
    └─────┬─────┘      └─────┬─────┘      └─────┬─────┘
          │                  │                  │
          │  Measurements    │  Measurements    │
          │  Alive Signals   │  Alive Signals   │
          │                  │                  │
          └──────────────────┼──────────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │   Data Server  │
                    │                │
                    │  REST API      │
                    │  + PostgreSQL  │
                    └────────┬───────┘
                             │
                             │ REST API
                             │
                             ▼
                    ┌────────────────┐
                    │    Frontend    │
                    │                │
                    │  React + Vite  │
                    │  + Charts      │
                    └────────────────┘
```

## 🛠️ Technology Stack

### Speedtest Node
- **Language**: Go 1.24+
- **Database**: SQLite3
- **Scheduler**: Cron
- **External**: Ookla Speedtest CLI
- **Container**: Docker

### Data Server
- **Language**: Go 1.24+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL 15+
- **Authentication**: JWT + bcrypt
- **Container**: Docker

### Frontend
- **Framework**: React 18+
- **Build Tool**: Vite 5+
- **Language**: TypeScript
- **Styling**: TailwindCSS 3+
- **Charts**: Apache ECharts 5+
- **State**: Zustand
- **Container**: Docker + Nginx

## 📁 Project Structure

```
network-measure-app/
├── speedtest-node/           # Measurement node application
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── README.md             # ← Detailed node specs
│
├── data-server/              # Central API server
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── migrations/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── README.md             # ← Detailed server specs
│
├── frontend/                 # Web dashboard
│   ├── src/
│   ├── public/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── README.md             # ← Detailed frontend specs
│
├── README.md                 # This file
└── NOTE.md                   # Original planning notes
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- (For nodes) Ookla Speedtest CLI

### 1. Start the Data Server

```bash
cd data-server

# Set up environment variables
export DB_PASSWORD="your-secure-db-password"
export ADMIN_USERNAME="admin"
export ADMIN_PASSWORD="your-secure-admin-password"
export JWT_SECRET="your-secret-jwt-key-change-in-production"

# Start with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f
```

The server will be available at `http://localhost:8080`

### 2. Generate an API Key

```bash
# Access the server and create an API key
# You can do this via the frontend once it's running
# Or use curl to create one directly:
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'

# Use the returned token to create an API key
curl -X POST http://localhost:8080/api/v1/admin/api-keys \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"My First Node"}'
```

### 3. Deploy Speedtest Node(s)

```bash
cd speedtest-node

# Run with command-line arguments
./speedtest-node \
  --node-name="home-node" \
  --server-url="http://localhost:8080" \
  --api-key="your-api-key-here"

# Or use environment variables
export NODE_NAME="home-node"
export SERVER_URL="http://localhost:8080"
export API_KEY="your-api-key-here"
./speedtest-node

# Or use Docker Compose
docker-compose up -d

# Verify it's running
docker-compose logs -f
```

Deploy multiple nodes by copying the directory and using different node names via arguments or environment variables.

### 4. Start the Frontend

```bash
cd frontend

# Set environment variable for API URL
export VITE_API_URL="http://localhost:8080"

# Install dependencies and start
npm install
npm run dev

# Or use Docker
docker-compose up -d
```

The dashboard will be available at `http://localhost:3000`

### 5. Login and Explore

1. Navigate to `http://localhost:3000`
2. Login with your admin credentials
3. View the dashboard to see measurements coming in
4. Manage nodes and API keys from the admin panel

## 📖 Documentation

Each component has comprehensive documentation:

- **[Speedtest Node Documentation](./speedtest-node/README.md)**
  - Configuration options
  - Database schema
  - API communication
  - Docker deployment
  - Troubleshooting

- **[Data Server Documentation](./data-server/README.md)**
  - API endpoints
  - Database schema
  - Authentication & security
  - Docker deployment
  - Performance tuning

- **[Frontend Documentation](./frontend/README.md)**
  - Component structure
  - State management
  - Chart configuration
  - Docker deployment
  - Customization

## 🔒 Security Best Practices

1. **Change Default Credentials**: Update admin username/password in data server config
2. **Use HTTPS**: Enable TLS/SSL for production deployments
3. **Secure API Keys**: Store API keys securely, never commit to version control
4. **Database Security**: Use strong passwords, enable SSL connections
5. **Network Isolation**: Use Docker networks to isolate components
6. **Regular Updates**: Keep dependencies and Docker images updated
7. **Backup Data**: Regular backups of PostgreSQL database

## 📊 Data Flow

### Measurement Collection
1. Node runs speedtest at configured interval (default: every 10 minutes)
2. Result stored in local SQLite database
3. Measurement queued for transmission to server

### Data Transmission
1. Background worker checks for unsent data every 30 seconds
2. Up to 20 measurements sent in batch to server
3. On success (HTTP 200), local data can be deleted after 7 days
4. On failure, data retained and retried later

### Alive Signals
1. Every node sends alive signal every 60 seconds
2. Server updates node's `last_alive` timestamp
3. Nodes missing 2+ signals (2 minutes) marked as "unreachable"
4. Nodes with no signal for 24+ hours marked as "inactive"

### Dashboard Updates
1. Frontend polls server API every 10 seconds
2. Aggregated data fetched based on selected time range
3. Charts update with latest measurements
4. Node status indicators reflect real-time state

## 🔧 Configuration

All components use **command-line arguments** and/or **environment variables** for configuration. Command-line arguments take precedence over environment variables.

### Node Configuration
```bash
# Command-line arguments
./speedtest-node \
  --node-name="home-office" \
  --server-url="https://speedtest.example.com" \
  --api-key="sk_live_..." \
  --speedtest-cron="*/10 * * * *" \
  --retention-days=7

# Environment variables
export NODE_NAME="home-office"
export SERVER_URL="https://speedtest.example.com"
export API_KEY="sk_live_..."
export SPEEDTEST_CRON="*/10 * * * *"
export RETENTION_DAYS=7
```

### Server Configuration
```bash
# Command-line arguments
./data-server \
  --host="0.0.0.0" \
  --port=8080 \
  --db-host="postgres" \
  --db-password="secure-password" \
  --admin-username="admin" \
  --admin-password="admin-password" \
  --jwt-secret="your-secret-key" \
  --retention-measurements=365

# Environment variables
export HOST="0.0.0.0"
export PORT=8080
export DB_HOST="postgres"
export DB_PASSWORD="secure-password"
export ADMIN_USERNAME="admin"
export ADMIN_PASSWORD="admin-password"
export JWT_SECRET="your-secret-key"
export RETENTION_MEASUREMENTS=365
```

### Frontend Configuration
```bash
VITE_API_URL=https://speedtest-api.example.com
VITE_REFRESH_INTERVAL=10000  # 10 seconds
```

## 📈 Monitoring & Maintenance

### Health Checks
- Node: Check Docker container status and logs
- Server: Access `/health` endpoint
- Frontend: Access the dashboard

### Log Locations
- Node: `./speedtest-node/logs/`
- Server: `./data-server/logs/`
- Frontend: Browser console + Nginx logs

### Database Maintenance
```bash
# Backup PostgreSQL
docker exec speedtest-db pg_dump -U speedtest speedtest > backup.sql

# Restore
docker exec -i speedtest-db psql -U speedtest speedtest < backup.sql

# Check database size
docker exec speedtest-db psql -U speedtest -d speedtest \
  -c "SELECT pg_size_pretty(pg_database_size('speedtest'));"
```

## 🐛 Troubleshooting

### Node Can't Connect to Server
- Verify server URL in node config
- Check API key is valid and enabled
- Test network connectivity: `curl https://your-server.com/health`
- Review node logs for specific errors

### No Data in Dashboard
- Verify nodes are running and sending data
- Check server logs for incoming measurements
- Confirm database contains data: `SELECT COUNT(*) FROM measurements;`
- Verify frontend API URL configuration

### Charts Not Updating
- Check browser console for errors
- Verify frontend can reach server API
- Check CORS configuration on server
- Ensure JWT token is valid

## 🚧 Future Enhancements

### Planned Features
- [ ] Multi-user support with role-based access
- [ ] WebSocket for real-time updates (replace polling)
- [ ] Alert system for node failures or performance issues
- [ ] Email/webhook notifications
- [ ] Data export (CSV, JSON, PDF reports)
- [ ] Grafana integration
- [ ] Prometheus metrics export
- [ ] Mobile app
- [ ] Advanced analytics and anomaly detection
- [ ] Custom speedtest parameters
- [ ] Network topology visualization

### Nice to Have
- [ ] Dark mode for frontend
- [ ] Multi-language support
- [ ] Speed test scheduling per node
- [ ] Bandwidth usage calculator
- [ ] Historical comparison tools
- [ ] Performance trends and predictions
- [ ] ISP comparison features

## 📄 License

This is a personal hobby project. Use at your own risk.

## 🤝 Contributing

This is currently a personal project, but suggestions and improvements are welcome!

## 📞 Support

For detailed documentation, refer to individual component READMEs:
- [Speedtest Node](./speedtest-node/README.md)
- [Data Server](./data-server/README.md)
- [Frontend](./frontend/README.md)

## 🙏 Acknowledgments

- Ookla for the Speedtest CLI
- Apache ECharts for charting library
- Open source community for the various libraries and tools used

---

**Made with ❤️ for better network monitoring**
