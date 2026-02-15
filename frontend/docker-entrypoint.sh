#!/bin/sh

# Generate runtime configuration from environment variables
cat > /app/dist/config.js << EOF
window.runtimeConfig = {
  apiUrl: '${API_URL:-http://localhost:8080}',
  refreshInterval: ${REFRESH_INTERVAL:-10000},
  enableChartAnimation: ${ENABLE_CHART_ANIMATION:-false},
  debug: ${DEBUG:-false}
};
EOF

echo "Runtime configuration generated:"
cat /app/dist/config.js

# Start the server
exec serve -s dist -l 3000
