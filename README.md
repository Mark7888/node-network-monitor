# Network Speed Measurement System

This is a simple system to monitor your internet speed over time. It runs speedtests automatically, stores the results, and shows you nice charts in a web dashboard.

## What does it do?

The system has a few parts that work together:
- **Speedtest nodes** that run speed tests every 10 minutes (configurable)
- A **data server** that collects and stores all the measurements
- A **web dashboard** where you can see charts and stats
- Everything runs in Docker containers, so it's easy to set up

It's useful if you want to track your ISP's performance, see when your internet is slow, or just keep an eye on your connection quality.

## Getting Started

### 1. Install Docker

If you don't have Docker and Docker Compose installed yet:

- Install Docker: https://docs.docker.com/engine/install/
- Install Docker Compose: https://docs.docker.com/compose/install/

Make sure both are working by running `docker --version` and `docker compose version`

### 2. Configure the application

Copy the example environment file and edit it with your settings:

```bash
cp .env.example .env
nano .env  # or use your favorite editor
```

You need to set these values in the `.env` file:

```bash
# Database password (choose something secure)
POSTGRES_PASSWORD=your_secure_password_here
DB_PASSWORD=your_secure_password_here

# Admin credentials for the dashboard
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_admin_password

# JWT secret (use a random string)
JWT_SECRET=your_random_secret_key

# Node configuration
SPEEDTEST_NODE_NAME=my-home-node
SPEEDTEST_SERVER_API_KEY=temporary_key  # We'll generate a real one after starting
```

The other settings have sensible defaults, but you can adjust them if needed.

### 3. Run everything

Start all the services with Docker Compose:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database (port 5432)
- Data server API (port 8080)
- Web dashboard (port 3000)
- A speedtest node

Wait a minute for everything to start up, then check if it's running:

```bash
docker-compose ps
```

### 4. Generate an API key for the node

The speedtest node needs a proper API key to send data to the server. First, get a login token:

```bash
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_admin_password"}'
```

Copy the token from the response, then create an API key:

```bash
curl -X POST http://localhost:8080/api/v1/admin/api-keys \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"name":"my-home-node"}'
```

Copy the API key from the response, add it to your `.env` file under `SPEEDTEST_SERVER_API_KEY`, and restart the node:

```bash
docker-compose restart network-monitor-node
```

### 5. Open the dashboard

Go to http://localhost:3000 in your browser and log in with the admin credentials you set up. You should start seeing speedtest results after a few minutes.

## Managing the system

**View logs:**
```bash
docker-compose logs -f  # all services
docker-compose logs -f network-monitor-node  # just the speedtest node
```

**Stop everything:**
```bash
docker-compose down
```

**Update to latest version:**
```bash
docker-compose pull
docker-compose up -d
```

**Add more nodes:**
You can run additional nodes on other computers or networks by using the individual docker-compose files in the `speedtest-node/` directory. Each node needs its own API key.

## Troubleshooting

**Can't connect to the dashboard?**
- Make sure all containers are running: `docker-compose ps`
- Check the logs: `docker-compose logs`

**No data showing up?**
- Check if the node is running: `docker-compose logs network-monitor-node`
- Make sure you've set a valid API key in the `.env` file
- The first measurement takes about 10 minutes to appear

**Want to change settings?**
- Edit the `.env` file
- Run `docker-compose up -d` to apply changes

## More Information

Each component has its own README with more details:
- [data-server/README.md](./data-server/README.md) - API documentation and advanced config
- [speedtest-node/README.md](./speedtest-node/README.md) - Node setup and troubleshooting
- [frontend/README.md](./frontend/README.md) - Dashboard features

## License

See the LICENSE file for details.
