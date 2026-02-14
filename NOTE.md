# Network Speed Measurement System - Project Overview

A distributed network monitoring system that measures network speeds from multiple nodes, aggregates data in a central server, and visualizes metrics through an interactive dashboard.

## Project Status
✅ **Planning Phase Complete** - Detailed specifications created for all three components.

## System Architecture

```
┌─────────────────┐
│  Speedtest Node │ ──┐
│   (Go + SQLite) │   │
└─────────────────┘   │
                      │
┌─────────────────┐   │    ┌──────────────────┐      ┌─────────────────┐
│  Speedtest Node │ ──┼───▶│   Data Server    │◀─────│    Frontend     │
│   (Go + SQLite) │   │    │ (Go + PostgreSQL)│      │ (React + Vite)  │
└─────────────────┘   │    └──────────────────┘      └─────────────────┘
                      │
┌─────────────────┐   │
│  Speedtest Node │ ──┘
│   (Go + SQLite) │
└─────────────────┘
```

## Components

### 1. [Speedtest Node](./speedtest-node/README.md) - Data Collector
**Tech Stack**: Go, SQLite, Docker

Lightweight measurement nodes that:
- Run Ookla speedtest every 10 minutes (configurable)
- Store results in local SQLite database
- Send measurements to central server in batches (max 20)
- Send alive signals every minute
- Self-register with unique UUID
- Retry failed tests and queue unsent data
- Keep data for 7 days maximum

**[→ View Detailed Specifications](./speedtest-node/README.md)**

### 2. [Data Server](./data-server/README.md) - Central Collector
**Tech Stack**: Go, PostgreSQL, Docker

Central API server that:
- Accepts measurements from multiple nodes
- Tracks node status (active/unreachable after 2 minutes)
- Manages API keys with bcrypt hashing
- Provides JWT-authenticated admin API
- Aggregates data for chart visualization
- Enforces 1-year retention policy (configurable)
- Handles duplicate measurements (upsert by timestamp)

**[→ View Detailed Specifications](./data-server/README.md)**

### 3. [Frontend](./frontend/README.md) - Dashboard & Management
**Tech Stack**: React, TypeScript, TailwindCSS, Apache ECharts, Docker

Interactive web dashboard that:
- Displays network metrics in real-time charts
- Shows download/upload speeds, ping, jitter, packet loss
- Manages nodes and API keys
- Filters by time range (day/week/month)
- Auto-refreshes every 10 seconds
- Provides per-node detailed views

**[→ View Detailed Specifications](./frontend/README.md)**

## Quick Start

See individual component READMEs for detailed setup instructions.

### Prerequisites
- Docker & Docker Compose
- Ookla Speedtest CLI (for nodes)

### Basic Setup

1. **Start Data Server**
```bash
cd data-server
# Set environment variables
export DB_PASSWORD="your-secure-password"
export ADMIN_USERNAME="admin"
export ADMIN_PASSWORD="your-admin-password"
export JWT_SECRET="your-secret-key"
docker-compose up -d
```

2. **Start Speedtest Node**
```bash
cd speedtest-node
# Set environment variables or use command-line args
export NODE_NAME="home-node"
export SERVER_URL="http://localhost:8080"
export API_KEY="your-api-key"
docker-compose up -d
```

3. **Start Frontend**
```bash
cd frontend
# Set environment variable
export VITE_API_URL="http://localhost:8080"
npm install && npm run dev
# Or use docker-compose up -d
```

## Project Plans

I want to make a small hobby project, that's only goal is to measure network speeds, send them to a central server, that stores it in a database, and than i can draw all sorts of charts with that data.

I need 3 separate applications for it.

## First application: Data collector, measuring node.

This will be a small script, that runs ookla's speedtest (speedtest --format json) in a set cron time. Than stores the output in some sort of storage, like an sqlite db.
It does not have to store anything else, just these data.
Here is how the returned data looks like:
    {"type":"result","timestamp":"2026-02-14T17:46:10Z","ping":{"jitter":0.061,"latency":1.089,"low":1.060,"high":1.198},"download":{"bandwidth":11739031,"bytes":42291720,"elapsed":3602,"latency":{"iqm":9.540,"low":3.194,"high":10.221,"jitter":0.669}},"upload":{"bandwidth":11723027,"bytes":42249312,"elapsed":3604,"latency":{"iqm":128.026,"low":5.304,"high":305.608,"jitter":39.456}},"packetLoss":0,"isp":"Eltenet Eltenet","interface":{"internalIp":"172.31.55.120","name":"eth0","macAddr":"00:15:5D:59:8A:21","isVpn":false,"externalIp":"157.181.192.69"},"server":{"id":28951,"host":"bp-speedtest.zt.hu","port":8080,"name":"ZNET Telekom Zrt.","location":"Budapest","country":"Hungary","ip":"185.232.83.0"},"result":{"id":"be79d55c-6c2d-4a4b-b56b-3f23aa589c0f","url":"https://www.speedtest.net/result/c/be79d55c-6c2d-4a4b-b56b-3f23aa589c0f","persisted":true}}
Store these data for up to configured days. This also has to have a config input for the target data server it will periodically send an alive signal. At every measurement, save that into the db, and also send it to the data server. If successfully sent and got returned 200, mark that row in the db, that it's already sent. If the data could not be sent to the data server, put it in some sort of a queue, and next time it is reachable (known by the alive signal), send it to it. I plan to make this as a docker container, but also runnable just the app.

Configs (not extensive, just the basics):
- Cron time for the measurements.
- Target data server url.
- Target data server api key.
- Number of days to keep the data in the db.
- Node name, so i can identify it in the data server.

Technologies: Golang, SQLite, Docker

## Second app: Data server.

This will be the main part, it runs on a reachable vps, and it collects data from accross the measuring nodes.
I want to be able to have multiple measuring nodes, and all of them send their data to this server. It will have a simple REST API, that accepts the data from the measuring nodes, and stores it in a database.
It also has to have an endpoint for the alive signal, that the measuring nodes will send periodically. This way it can keep track of which measuring nodes are active, and which are not.
I want this to be able to generate api keys for the measuring nodes, so only authorized nodes can send data to it. I only want this to have rest api, no frontend, as the frontend will be a separate app. The api has to have features, like authenticate with the frontend. The frontend authenticated user can see the nodes and their status and data, and also can see a list of active api keys, create new, revoke existing, or just disable/enable existing keys.

Configs:
- Database connection settings.
- Host and port to run the server on.
- Admin username and password for the frontend authentication.

Technologies: Golang, PostgreSQL, Docker

## Third app: Frontend.

This will be a simple web application, that connects to the data server, and shows the data in a nice way. It will have a login system, that authenticates with the data server, and shows the data in a nice way. It will have a dashboard, that shows the status of the measuring nodes, and also shows the data in charts. It will also have a page for managing the api keys and nodes.
I want to make the charts using Apache ECharts, i'm interested in network speeds, ping, jitter, and packet loss, so i want to show these data in a nice way.
I want to be able to click on a node, and see it's details, as well as the measurements in charts, and also have a page where i can see all the measurements in a table for just this node.
But also want to have a dashboard page, where i can see all the measurements accross nodes with charts per data type, and filter them by node, date, etc. (filtering filters all the charts).

Configs:
- Host and port to run the frontend on.
- Data server url for the api calls.

Technologies: React, Vite, TailwindCSS, Apache ECharts, Docker
