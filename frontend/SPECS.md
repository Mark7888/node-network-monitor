# Frontend - Network Measurement Dashboard

## Overview
A modern web application for visualizing network measurement data from multiple speedtest nodes. Features real-time updates, interactive charts, node management, and API key administration. Built with React, TailwindCSS, and Apache ECharts.

## Technologies & Dependencies

### Core
- **Framework**: React 18+
- **Build Tool**: Vite 5+
- **Language**: TypeScript
- **Styling**: TailwindCSS 3+
- **Charts**: Apache ECharts 5+
- **Container**: Docker

### Key Packages
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.22.0",
    "echarts": "^5.5.0",
    "echarts-for-react": "^3.0.2",
    "axios": "^1.6.5",
    "date-fns": "^3.3.1",
    "zustand": "^4.5.0",
    "react-hot-toast": "^2.4.1",
    "@headlessui/react": "^1.7.18",
    "@heroicons/react": "^2.1.1"
  },
  "devDependencies": {
    "@types/react": "^18.2.55",
    "@types/react-dom": "^18.2.19",
    "@vitejs/plugin-react": "^4.2.1",
    "tailwindcss": "^3.4.1",
    "postcss": "^8.4.35",
    "autoprefixer": "^10.4.17",
    "typescript": "^5.3.3",
    "vite": "^5.1.0"
  }
}
```

## Architecture

### Components Structure
```
src/
├── components/
│   ├── layout/
│   │   ├── Header.tsx              # Top navigation bar
│   │   ├── Sidebar.tsx             # Side navigation
│   │   └── Layout.tsx              # Main layout wrapper
│   ├── charts/
│   │   ├── DownloadChart.tsx       # Download speed chart
│   │   ├── UploadChart.tsx         # Upload speed chart
│   │   ├── PingChart.tsx           # Ping/latency chart
│   │   ├── JitterChart.tsx         # Jitter chart
│   │   └── PacketLossChart.tsx     # Packet loss chart
│   ├── nodes/
│   │   ├── NodeCard.tsx            # Node status card
│   │   ├── NodeList.tsx            # List of nodes
│   │   ├── NodeDetails.tsx         # Detailed node view
│   │   └── NodeStatusBadge.tsx     # Status indicator
│   ├── api-keys/
│   │   ├── APIKeyList.tsx          # List of API keys
│   │   ├── APIKeyCard.tsx          # API key item
│   │   ├── CreateAPIKeyModal.tsx   # Create new key modal
│   │   └── APIKeyDialog.tsx        # Show generated key
│   ├── auth/
│   │   ├── LoginForm.tsx           # Login form
│   │   └── ProtectedRoute.tsx      # Auth guard
│   ├── common/
│   │   ├── Button.tsx              # Reusable button
│   │   ├── Input.tsx               # Reusable input
│   │   ├── Card.tsx                # Reusable card
│   │   ├── Badge.tsx               # Status badges
│   │   ├── Modal.tsx               # Modal dialog
│   │   ├── Table.tsx               # Data table
│   │   ├── Spinner.tsx             # Loading spinner
│   │   └── TimeRangeFilter.tsx     # Time filter component
│   └── dashboard/
│       ├── DashboardSummary.tsx    # Summary cards
│       └── AllNodesCharts.tsx      # Combined charts view
```

### State Management
- **Zustand** for global state (auth, nodes, measurements)
- **React Query** (optional) for server state caching
- Local component state for UI

### Routing
```
/ - Redirect to /dashboard
/login - Login page
/dashboard - Main dashboard (all nodes overview)
/nodes - Node list
/nodes/:id - Single node details
/api-keys - API key management
```

## Pages & Features

### 1. Login Page
**Route**: `/login`

**Features**:
- Username/password form
- Error handling
- Redirect to dashboard on success
- Remember token in localStorage

**UI**:
```
┌──────────────────────────────────┐
│     Network Speedtest Monitor    │
│                                   │
│  ┌─────────────────────────────┐ │
│  │ Username                     │ │
│  │ [________________]           │ │
│  │                              │ │
│  │ Password                     │ │
│  │ [________________]           │ │
│  │                              │ │
│  │       [  Login  ]            │ │
│  └─────────────────────────────┘ │
└──────────────────────────────────┘
```

### 2. Dashboard Page
**Route**: `/dashboard`

**Features**:
- Summary statistics (total nodes, active nodes, avg speeds)
- Time range filter (Last Day / Last Week / Last Month)
- Real-time updates (refresh every 10 seconds)
- Charts for all nodes combined:
  - Download speeds over time
  - Upload speeds over time
  - Ping/latency over time
  - Jitter over time
  - Packet loss over time
- Node filter (multi-select)

**Layout**:
```
┌─────────────────────────────────────────────────────────┐
│ [☰] Network Speedtest Monitor          [Admin] [Logout] │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Dashboard                                               │
│                                                          │
│  [Last Day] [Last Week] [Last Month]  [All Nodes ▾]     │
│                                                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐  │
│  │  5       │ │  4       │ │  95.2    │ │  1.2      │  │
│  │  Nodes   │ │  Active  │ │  Mbps ↓  │ │  ms ping  │  │
│  └──────────┘ └──────────┘ └──────────┘ └───────────┘  │
│                                                          │
│  Download Speed (Mbps)                                   │
│  ┌─────────────────────────────────────────────────┐    │
│  │        [Chart with multiple node lines]         │    │
│  └─────────────────────────────────────────────────┘    │
│                                                          │
│  Upload Speed (Mbps)                                     │
│  ┌─────────────────────────────────────────────────┐    │
│  │        [Chart with multiple node lines]         │    │
│  └─────────────────────────────────────────────────┘    │
│                                                          │
│  Ping / Latency (ms)                                     │
│  ┌─────────────────────────────────────────────────┐    │
│  │        [Chart with multiple node lines]         │    │
│  └─────────────────────────────────────────────────┘    │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 3. Nodes List Page
**Route**: `/nodes`

**Features**:
- List of all nodes with status
- Quick stats per node
- Click to view details
- Status indicators (active, unreachable, inactive)
- Last seen timestamp

**Layout**:
```
┌─────────────────────────────────────────────────────────┐
│ [☰] Network Speedtest Monitor                           │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Nodes                                                   │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │ ● home-office-node              [Active]           │ │
│  │   Last seen: 2 minutes ago                         │ │
│  │   Avg: ↓ 95.2 Mbps  ↑ 89.7 Mbps  Ping: 1.2 ms     │ │
│  │   [View Details →]                                  │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │ ● office-main                   [Active]           │ │
│  │   Last seen: 1 minute ago                          │ │
│  │   Avg: ↓ 120.5 Mbps ↑ 95.3 Mbps  Ping: 0.8 ms     │ │
│  │   [View Details →]                                  │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │ ◌ remote-site-1              [Unreachable]        │ │
│  │   Last seen: 5 minutes ago                         │ │
│  │   Avg: ↓ 45.2 Mbps  ↑ 12.5 Mbps  Ping: 25.3 ms    │ │
│  │   [View Details →]                                  │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 4. Node Details Page
**Route**: `/nodes/:id`

**Features**:
- Node information (name, ID, status, first seen, last seen)
- Time range filter
- Charts specific to this node:
  - Download speed
  - Upload speed
  - Ping/latency
  - Jitter
  - Packet loss
- Measurements table (paginated)
- Export data button (future)

**Layout**:
```
┌─────────────────────────────────────────────────────────┐
│ [☰] Network Speedtest Monitor                           │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ← Back to Nodes                                         │
│                                                          │
│  home-office-node                         [Active]      │
│  ID: 550e8400-e29b-41d4-a716-446655440000               │
│  First seen: Jan 1, 2026  |  Last seen: 2 min ago       │
│                                                          │
│  [Last Day] [Last Week] [Last Month]                    │
│                                                          │
│  Statistics                                              │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐  │
│  │  95.2    │ │  89.7    │ │  1.2     │ │  0.08     │  │
│  │  Mbps ↓  │ │  Mbps ↑  │ │  ms ping │ │  ms jitter│  │
│  └──────────┘ └──────────┘ └──────────┘ └───────────┘  │
│                                                          │
│  Download Speed (Mbps)                                   │
│  ┌─────────────────────────────────────────────────┐    │
│  │        [Line chart with data points]            │    │
│  └─────────────────────────────────────────────────┘    │
│                                                          │
│  Upload Speed (Mbps)                                     │
│  ┌─────────────────────────────────────────────────┐    │
│  │        [Line chart with data points]            │    │
│  └─────────────────────────────────────────────────┘    │
│                                                          │
│  Measurements History                                    │
│  ┌─────────────────────────────────────────────────┐    │
│  │ Timestamp       | Download | Upload | Ping      │    │
│  │ Feb 14, 17:40   | 95.2 Mbps| 89.7  | 1.2 ms    │    │
│  │ Feb 14, 17:30   | 93.8 Mbps| 88.2  | 1.3 ms    │    │
│  │ ...                                              │    │
│  │                           [1] [2] [3] ... [50] │    │
│  └─────────────────────────────────────────────────┘    │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 5. API Keys Management Page
**Route**: `/api-keys`

**Features**:
- List all API keys
- Show name, status (enabled/disabled), created date, last used
- Create new key button
- Enable/disable toggle
- Delete key button (with confirmation)
- Show generated key only once in modal

**Layout**:
```
┌─────────────────────────────────────────────────────────┐
│ [☰] Network Speedtest Monitor                           │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  API Keys                                                │
│                                            [+ New Key]   │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │ Production Node 1                      [Enabled]   │ │
│  │ Created: Jan 1, 2026  |  Last used: 2 min ago      │ │
│  │ [Disable] [Delete]                                  │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │ Test Node                             [Disabled]   │ │
│  │ Created: Jan 15, 2026  |  Last used: Never         │ │
│  │ [Enable] [Delete]                                   │ │
│  └────────────────────────────────────────────────────┘ │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

**Create API Key Modal**:
```
┌──────────────────────────────┐
│ Create New API Key     [✕]   │
├──────────────────────────────┤
│                              │
│ Name:                        │
│ [____________________]       │
│                              │
│        [Cancel] [Create]     │
└──────────────────────────────┘
```

**Show Generated Key Modal**:
```
┌──────────────────────────────────────────┐
│ API Key Created Successfully       [✕]   │
├──────────────────────────────────────────┤
│                                          │
│ ⚠️  Save this key securely!              │
│ It won't be shown again.                 │
│                                          │
│ sk_live_a1b2c3d4e5f6g7h8i9j0...         │
│ [Copy to Clipboard]                      │
│                                          │
│                           [Close]        │
└──────────────────────────────────────────┘
```

## Configuration

### Config File: `.env`
```bash
# API Server URL
VITE_API_URL=https://speedtest-api.example.com

# Refresh interval (milliseconds)
VITE_REFRESH_INTERVAL=10000

# Chart animation
VITE_ENABLE_CHART_ANIMATION=true

# Debug mode
VITE_DEBUG=false
```

### Environment-specific configs
- `.env.development` - Development settings
- `.env.production` - Production settings

## API Client

### Axios Instance
```typescript
// src/lib/api.ts
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
```

### API Service Functions
```typescript
// src/services/authService.ts
export const login = async (username: string, password: string) => {
  const response = await api.post('/api/v1/admin/login', { username, password });
  return response.data;
};

// src/services/nodeService.ts
export const getNodes = async () => {
  const response = await api.get('/api/v1/admin/nodes');
  return response.data;
};

export const getNodeDetails = async (nodeId: string) => {
  const response = await api.get(`/api/v1/admin/nodes/${nodeId}`);
  return response.data;
};

export const getNodeMeasurements = async (
  nodeId: string, 
  params: { from?: string; to?: string; limit?: number }
) => {
  const response = await api.get(`/api/v1/admin/nodes/${nodeId}/measurements`, { params });
  return response.data;
};

// src/services/measurementService.ts
export const getAggregatedMeasurements = async (params: {
  node_ids?: string[];
  from: string;
  to: string;
  interval: string;
}) => {
  const response = await api.get('/api/v1/admin/measurements/aggregate', { params });
  return response.data;
};

// src/services/apiKeyService.ts
export const getAPIKeys = async () => {
  const response = await api.get('/api/v1/admin/api-keys');
  return response.data;
};

export const createAPIKey = async (name: string) => {
  const response = await api.post('/api/v1/admin/api-keys', { name });
  return response.data;
};

export const updateAPIKey = async (id: string, enabled: boolean) => {
  const response = await api.patch(`/api/v1/admin/api-keys/${id}`, { enabled });
  return response.data;
};

export const deleteAPIKey = async (id: string) => {
  await api.delete(`/api/v1/admin/api-keys/${id}`);
};

// src/services/dashboardService.ts
export const getDashboardSummary = async () => {
  const response = await api.get('/api/v1/admin/dashboard');
  return response.data;
};
```

## State Management (Zustand)

### Auth Store
```typescript
// src/store/authStore.ts
import { create } from 'zustand';

interface AuthState {
  token: string | null;
  username: string | null;
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  username: localStorage.getItem('username'),
  isAuthenticated: !!localStorage.getItem('token'),
  
  login: async (username, password) => {
    const data = await authService.login(username, password);
    localStorage.setItem('token', data.token);
    localStorage.setItem('username', data.username);
    set({ token: data.token, username: data.username, isAuthenticated: true });
  },
  
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    set({ token: null, username: null, isAuthenticated: false });
  },
  
  checkAuth: () => {
    const token = localStorage.getItem('token');
    set({ isAuthenticated: !!token });
  },
}));
```

### Nodes Store
```typescript
// src/store/nodesStore.ts
import { create } from 'zustand';

interface NodesState {
  nodes: Node[];
  loading: boolean;
  error: string | null;
  fetchNodes: () => Promise<void>;
}

export const useNodesStore = create<NodesState>((set) => ({
  nodes: [],
  loading: false,
  error: null,
  
  fetchNodes: async () => {
    set({ loading: true, error: null });
    try {
      const data = await nodeService.getNodes();
      set({ nodes: data.nodes, loading: false });
    } catch (error) {
      set({ error: error.message, loading: false });
    }
  },
}));
```

## Chart Configuration

### ECharts Options Example
```typescript
// src/components/charts/DownloadChart.tsx
import ReactECharts from 'echarts-for-react';

const DownloadChart = ({ data, timeRange }) => {
  const option = {
    title: {
      text: 'Download Speed',
      textStyle: { color: '#374151', fontSize: 16, fontWeight: 600 }
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params) => {
        let tooltip = `${params[0].axisValue}<br/>`;
        params.forEach(param => {
          tooltip += `${param.marker} ${param.seriesName}: ${param.value} Mbps<br/>`;
        });
        return tooltip;
      }
    },
    legend: {
      data: data.map(d => d.node_name),
      bottom: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'time',
      boundaryGap: false,
      axisLabel: {
        formatter: (value) => format(new Date(value), 'HH:mm')
      }
    },
    yAxis: {
      type: 'value',
      name: 'Mbps',
      axisLabel: {
        formatter: '{value}'
      }
    },
    series: data.map(node => ({
      name: node.node_name,
      type: 'line',
      smooth: true,
      data: node.measurements.map(m => [
        m.timestamp,
        (m.download_bandwidth / 1000000 * 8).toFixed(2) // Convert to Mbps
      ]),
      itemStyle: {
        color: getNodeColor(node.node_id)
      }
    }))
  };

  return <ReactECharts option={option} style={{ height: 400 }} />;
};
```

### Chart Data Transformation
```typescript
// Convert bandwidth (bytes/sec) to Mbps
const bytesToMbps = (bytes: number): number => {
  return (bytes / 1000000) * 8;
};

// Format timestamp for display
const formatTimestamp = (timestamp: string, range: string): string => {
  const date = new Date(timestamp);
  if (range === 'day') return format(date, 'HH:mm');
  if (range === 'week') return format(date, 'MMM dd HH:mm');
  return format(date, 'MMM dd');
};

// Aggregate data for display
const aggregateData = (measurements, interval) => {
  // Group by time intervals
  // Calculate averages, min, max
  // Return formatted data for charts
};
```

## Real-time Updates

### Auto-refresh Hook
```typescript
// src/hooks/useAutoRefresh.ts
import { useEffect, useRef } from 'react';

export const useAutoRefresh = (callback: () => void, interval: number) => {
  const savedCallback = useRef(callback);

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    const tick = () => savedCallback.current();
    const id = setInterval(tick, interval);
    return () => clearInterval(id);
  }, [interval]);
};

// Usage in Dashboard component
const Dashboard = () => {
  const fetchData = async () => {
    await fetchDashboardSummary();
    await fetchAggregatedMeasurements();
  };

  useAutoRefresh(fetchData, 10000); // Refresh every 10 seconds

  // ...
};
```

## Time Range Filter

### Time Range Component
```typescript
// src/components/common/TimeRangeFilter.tsx
const TimeRangeFilter = ({ value, onChange }) => {
  const ranges = [
    { label: 'Last Day', value: 'day', hours: 24 },
    { label: 'Last Week', value: 'week', hours: 168 },
    { label: 'Last Month', value: 'month', hours: 720 },
  ];

  return (
    <div className="flex gap-2">
      {ranges.map(range => (
        <button
          key={range.value}
          onClick={() => onChange(range.value)}
          className={`px-4 py-2 rounded ${
            value === range.value
              ? 'bg-blue-600 text-white'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
          }`}
        >
          {range.label}
        </button>
      ))}
    </div>
  );
};

// Calculate time range
const getTimeRange = (range: string) => {
  const to = new Date();
  const from = new Date();
  
  switch (range) {
    case 'day':
      from.setHours(from.getHours() - 24);
      break;
    case 'week':
      from.setDate(from.getDate() - 7);
      break;
    case 'month':
      from.setMonth(from.getMonth() - 1);
      break;
  }
  
  return {
    from: from.toISOString(),
    to: to.toISOString(),
  };
};
```

## Styling with TailwindCSS

### Tailwind Config
```javascript
// tailwind.config.js
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          // ... more shades
          900: '#1e3a8a',
        },
      },
    },
  },
  plugins: [],
}
```

### Common Styles
```typescript
// Status badge colors
const statusColors = {
  active: 'bg-green-100 text-green-800',
  unreachable: 'bg-yellow-100 text-yellow-800',
  inactive: 'bg-gray-100 text-gray-800',
};

// Button variants
const buttonVariants = {
  primary: 'bg-blue-600 hover:bg-blue-700 text-white',
  secondary: 'bg-gray-200 hover:bg-gray-300 text-gray-700',
  danger: 'bg-red-600 hover:bg-red-700 text-white',
};
```

## Docker Setup

### Dockerfile
```dockerfile
FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### nginx.conf
```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy (optional - if serving from same domain)
    location /api/ {
        proxy_pass http://data-server:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
}
```

### docker-compose.yml
```yaml
version: '3.8'

services:
  frontend:
    build: .
    container_name: speedtest-frontend
    restart: unless-stopped
    ports:
      - "3000:80"
    environment:
      - VITE_API_URL=${API_URL:-http://localhost:8080}
    depends_on:
      - data-server
```

### .env file
```bash
API_URL=http://localhost:8080
REFRESH_INTERVAL=10000
```

## Project Structure

```
frontend/
├── public/
│   └── favicon.ico
├── src/
│   ├── components/         # React components
│   ├── pages/             # Page components
│   ├── services/          # API services
│   ├── store/             # Zustand stores
│   ├── hooks/             # Custom hooks
│   ├── lib/               # Utilities
│   ├── types/             # TypeScript types
│   ├── styles/            # Global styles
│   ├── App.tsx            # Main app component
│   ├── main.tsx           # Entry point
│   └── router.tsx         # Route definitions
├── .env
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── index.html
├── nginx.conf
├── package.json
├── postcss.config.js
├── tailwind.config.js
├── tsconfig.json
├── vite.config.ts
└── README.md
```

## Key Implementation Details

### 1. Protected Routes
```typescript
// src/components/auth/ProtectedRoute.tsx
const ProtectedRoute = ({ children }) => {
  const { isAuthenticated } = useAuthStore();
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  
  return children;
};

// Usage in router
<Route path="/dashboard" element={
  <ProtectedRoute>
    <Dashboard />
  </ProtectedRoute>
} />
```

### 2. Error Handling
```typescript
// src/components/common/ErrorBoundary.tsx
class ErrorBoundary extends React.Component {
  state = { hasError: false, error: null };
  
  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }
  
  render() {
    if (this.state.hasError) {
      return <ErrorPage error={this.state.error} />;
    }
    return this.props.children;
  }
}
```

### 3. Toast Notifications
```typescript
// Using react-hot-toast
import toast from 'react-hot-toast';

// Success
toast.success('API key created successfully!');

// Error
toast.error('Failed to load data');

// Loading
const toastId = toast.loading('Creating API key...');
// Later:
toast.success('Created!', { id: toastId });
```

### 4. Data Caching
```typescript
// Simple cache for measurements
const measurementCache = new Map();

const getCachedMeasurements = (key, fetchFn, ttl = 30000) => {
  const cached = measurementCache.get(key);
  if (cached && Date.now() - cached.timestamp < ttl) {
    return cached.data;
  }
  
  const data = await fetchFn();
  measurementCache.set(key, { data, timestamp: Date.now() });
  return data;
};
```

## Running the Application

### Development
```bash
# Install dependencies
npm install

# Start dev server
npm run dev

# Available at http://localhost:5173
```

### Production Build
```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

### Docker
```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f frontend

# Stop
docker-compose down
```

## Testing

### Unit Tests (Vitest)
```bash
npm run test
```

### E2E Tests (Playwright - optional)
```bash
npm run test:e2e
```

## Performance Optimization

### Code Splitting
```typescript
// Lazy load pages
const Dashboard = lazy(() => import('./pages/Dashboard'));
const NodeDetails = lazy(() => import('./pages/NodeDetails'));

// Wrap in Suspense
<Suspense fallback={<Spinner />}>
  <Dashboard />
</Suspense>
```

### Memoization
```typescript
// Memoize expensive calculations
const chartData = useMemo(() => {
  return processChartData(measurements);
}, [measurements]);

// Memoize components
const NodeCard = memo(({ node }) => {
  // ...
});
```

### Virtual Scrolling
For large measurement tables, use react-virtualized or react-window.

## Accessibility

- Semantic HTML
- ARIA labels for interactive elements
- Keyboard navigation support
- Focus management in modals
- Color contrast compliance (WCAG AA)

## Security Considerations

1. **XSS Prevention**: React escapes by default
2. **CSRF**: Not needed (JWT in header, not cookie)
3. **Token Storage**: localStorage (acceptable for this use case)
4. **HTTPS**: Always use in production
5. **API Key Display**: Only show once, copy to clipboard

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Troubleshooting

**Problem**: Can't connect to API
- Check VITE_API_URL in .env
- Verify CORS settings on server
- Check network tab in devtools

**Problem**: Charts not rendering
- Verify ECharts is installed
- Check console for errors
- Ensure data format matches schema

**Problem**: Auto-refresh not working
- Check VITE_REFRESH_INTERVAL
- Verify useAutoRefresh hook
- Check browser console for errors

## Future Enhancements

- Dark mode
- Export data (CSV, JSON, PDF reports)
- Custom date range picker
- Alert configuration UI
- User preferences (chart colors, refresh rate)
- Mobile responsive improvements
- PWA support (offline mode)
- WebSocket for real-time updates (instead of polling)
- Notification center
- Advanced filtering and search
