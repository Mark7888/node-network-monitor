# Network Speed Measurement System

A distributed network monitoring system that collects internet speed measurements from multiple nodes, stores data centrally, and provides an interactive web dashboard for visualization and analysis.

## 🎯 Overview

This system enables continuous monitoring of network performance across multiple locations. Deploy lightweight measurement nodes anywhere, collect data in a central server, and visualize trends through a web interface.

**Perfect for:**
- Home network monitoring across different locations
- ISP performance tracking and accountability
- Network diagnostics and troubleshooting
- Historical bandwidth analysis
- Multi-site network comparison

## ✨ Key Features

- **Automated Measurements**: Configurable cron-based speedtest execution (default: every 10 minutes)
- **Multi-Node Support**: Deploy unlimited measurement nodes with automatic registration
- **Comprehensive Metrics**: Download/upload speeds, ping latency, jitter, packet loss, ISP info
- **Local Resilience**: Nodes store data locally and sync when connected
- **Secure Authentication**: API keys for nodes, JWT tokens for dashboard access
- **Interactive Dashboard**: Real-time charts with historical data visualization
- **Self-Hosted**: Complete Docker deployment, no external dependencies
- **Data Retention**: Configurable retention policies for both nodes and server

## 🏗️ Architecture

The system consists of three main components that work together:

```
┌─────────────────────────────────────────────┐
│             Speedtest Nodes                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  Node 1  │  │  Node 2  │  │  Node N  │   │
│  │  Go +    │  │  Go +    │  │  Go +    │   │
│  │  SQLite  │  │  SQLite  │  │  SQLite  │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼─────────────┼─────────────┼─────────┘
        │             │             │
        └─────────────┼─────────────┘
                      │ HTTPS + API Key
                      ▼
        ┌─────────────────────────┐
        │     Data Server         │
        │   Go + Gin Framework    │
        │   PostgreSQL / SQLite   │
        └────────────┬────────────┘
                     │ REST API + JWT
                     ▼
        ┌─────────────────────────┐
        │    Frontend Dashboard   │
        │   React + TypeScript    │
        │   ECharts + TailwindCSS │
        └─────────────────────────┘
```

## 📦 Components

### 1. [Speedtest Node](./speedtest-node/)
Lightweight measurement collectors that run Ookla Speedtest CLI on a schedule.

- **Technology**: Go 1.24+, SQLite3, Cron scheduler
- **Features**: Automated testing, local storage, batch sync, retry logic
- **Deployment**: Docker or standalone binary
- **Documentation**: [speedtest-node/README.md](./speedtest-node/README.md)

### 2. [Data Server](./data-server/)
Central API server that collects measurements and provides data access.

- **Technology**: Go 1.24+, Gin framework, PostgreSQL/SQLite
- **Features**: REST API, node management, API key system, JWT auth
- **Deployment**: Docker Compose with PostgreSQL or standalone
- **Documentation**: [data-server/README.md](./data-server/README.md)

### 3. [Frontend Dashboard](./frontend/)
Web interface for visualizing measurements and managing the system.

- **Technology**: React 18+, TypeScript, Vite, TailwindCSS, Apache ECharts
- **Features**: Interactive charts, node management, API key admin, time range filters
- **Deployment**: Docker with Nginx or development server
- **Documentation**: [frontend/README.md](./frontend/README.md) *(coming soon)*

## � Quick Start

### Prerequisites
- Docker & Docker Compose (recommended)
- Or: Go 1.24+, PostgreSQL 15+, Node.js 18+ (for manual setup)
- Ookla Speedtest CLI (for nodes)

### 1. Deploy the Data Server

```bash
cd data-server
cp .env.example .env  # Edit with your passwords and secrets
docker-compose up -d
```

The server will start at `http://localhost:8080`. See [data-server/README.md](./data-server/README.md) for detailed configuration.

### 2. Create an API Key

First, login to get a JWT token:
```bash
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'
```

Then create an API key for your node:
```bash
curl -X POST http://localhost:8080/api/v1/admin/api-keys \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My First Node"}'
```

Save the returned API key securely.

### 3. Deploy Speedtest Nodes

```bash
cd speedtest-node
export NODE_NAME="home-node"
export SERVER_URL="http://localhost:8080"
export API_KEY="your-api-key-here"
docker-compose up -d
```

Deploy additional nodes by repeating with different `NODE_NAME` values. See [speedtest-node/README.md](./speedtest-node/README.md) for more options.

### 4. Start the Dashboard

```bash
cd frontend
export VITE_API_URL="http://localhost:8080"
npm install && npm run dev
# Or use Docker: docker-compose up -d
```

Access the dashboard at `http://localhost:3000` and login with your admin credentials.

## 📖 Documentation

Each component has comprehensive documentation in its directory:

- **[Speedtest Node](./speedtest-node/README.md)** - Node setup, configuration, troubleshooting
- **[Data Server](./data-server/README.md)** - API endpoints, database schema, deployment
- **[Frontend](./frontend/README.md)** - Dashboard features, customization *(coming soon)*

## 🔒 Security Notes

- Change default admin credentials immediately
- Use HTTPS/TLS in production deployments
- Store API keys securely, never commit to version control
- Enable PostgreSQL SSL connections for production
- Use strong passwords and JWT secrets
- Keep Docker images and dependencies updated

See component READMEs for detailed security configurations.

## 📊 How It Works

1. **Measurement**: Nodes run Ookla Speedtest on a cron schedule (default: every 10 minutes)
2. **Local Storage**: Results stored in local SQLite database for resilience
3. **Sync**: Nodes batch-send measurements to the data server (up to 20 at a time)
4. **Keepalive**: Nodes send heartbeat signals every 60 seconds
5. **Status Tracking**: Server monitors node health based on heartbeat signals
6. **Visualization**: Dashboard fetches aggregated data and displays interactive charts
7. **Retention**: Old data automatically cleaned up based on retention policies

## 🛠️ Development

```bash
# Clone the repository
git clone https://github.com/yourusername/network-measure-app.git
cd network-measure-app

# Each component can be developed independently
cd speedtest-node && go run ./cmd/speedtest-node
cd data-server && go run ./cmd/data-server
cd frontend && npm run dev
```

## 📄 License

This project is licensed under the terms specified in the LICENSE file.

## 🙏 Acknowledgments

- [Ookla Speedtest CLI](https://www.speedtest.net/apps/cli) for network measurements
- [Apache ECharts](https://echarts.apache.org/) for powerful data visualization
- Open source community for the excellent tools and libraries

---

**Built for network enthusiasts who want to track their internet performance** 📈
