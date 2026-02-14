# Network Speedtest Monitor - Frontend

A modern, responsive web application for visualizing and managing network speed test measurements from multiple nodes. Built with React, TypeScript, TailwindCSS, daisyUI, and Apache ECharts.

## Features

- 📊 **Real-time Dashboard** - View aggregated metrics from all nodes
- 📈 **Interactive Charts** - Visualize download/upload speeds, ping, jitter, and packet loss
- 🖥️ **Node Management** - Monitor individual node status and performance
- 🔑 **API Key Management** - Create and manage API keys for nodes
- 🔄 **Auto-refresh** - Automatic data updates every 10 seconds
- 🎨 **Modern UI** - Clean, professional interface with daisyUI components
- 🌓 **Dark Mode** - Light and dark theme with toggle in navbar
- 📱 **Responsive** - Works on desktop and mobile devices

## Tech Stack

- **React 18** - UI framework
- **TypeScript** - Type-safe development
- **Vite** - Fast build tool
- **TailwindCSS + daisyUI** - Styling and UI components
- **Apache ECharts** - Interactive charts
- **Zustand** - State management
- **React Router** - Navigation
- **Axios** - HTTP client
- **date-fns** - Date formatting

## Project Structure

```
src/
├── modules/              # Feature modules
│   ├── auth/            # Authentication
│   ├── dashboard/       # Dashboard page
│   ├── nodes/           # Nodes management
│   ├── api-keys/        # API keys management
│   └── measurements/    # Charts and measurements
├── shared/              # Shared components and utilities
│   ├── components/      # Reusable UI components
│   ├── hooks/           # Custom hooks
│   ├── utils/           # Utility functions
│   └── types/           # Common types
└── core/                # Core application logic
    ├── api/             # API configuration
    ├── config/          # App configuration
    └── router/          # Routing setup
```

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

1. Install dependencies:
```bash
npm install
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Update `.env` with your API server URL:
```bash
VITE_API_URL=http://localhost:8080
```

### Development

Start the development server:
```bash
npm run dev
```

The app will be available at [http://localhost:5173](http://localhost:5173)

### Production Build

Build for production:
```bash
npm run build
```

Preview production build:
```bash
npm run preview
```

## Docker

### Build and Run

```bash
docker-compose up -d
```

The frontend will be available at [http://localhost:3000](http://localhost:3000)

### Using Docker Compose with Backend

See the root `docker-compose.yml` to run the full stack (frontend + backend + database).

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` |
| `VITE_REFRESH_INTERVAL` | Auto-refresh interval (ms) | `10000` |
| `VITE_ENABLE_CHART_ANIMATION` | Enable chart animations | `true` |
| `VITE_DEBUG` | Enable debug mode | `false` |

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Features Guide

### Dashboard
- View summary statistics across all nodes
- Filter data by time range (Last Day, Week, Month)
- Interactive charts for all metrics
- Auto-refreshing data

### Nodes
- List all registered nodes
- View node status (Active, Unreachable, Inactive)
- See last seen timestamp
- Navigate to detailed node view

### Node Details
- View individual node statistics
- Time-filtered charts for specific node
- Historical measurements
- Node metadata

### API Keys
- Create new API keys
- Enable/disable keys
- View key usage statistics
- Delete keys with confirmation

### Dark Mode
- Toggle between light and dark themes
- Theme preference saved in localStorage
- Smooth transitions between themes
- Sun/Moon icon in navbar header

## Code Style

- **TypeScript** - Strongly typed
- **Modular** - Feature-based organization
- **Clean** - Readable, maintainable code
- **Documented** - JSDoc comments on functions
- **Consistent** - ESLint + Prettier

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## License

See [LICENSE](../LICENSE) file in the root directory.

## Contributing

1. Follow the existing code structure
2. Use TypeScript strictly
3. Keep components small and focused
4. Write clean, documented code
5. Test thoroughly before committing
