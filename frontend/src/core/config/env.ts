/**
 * Application environment configuration
 * Reads from environment variables set in .env files
 */

export const env = {
  apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  refreshInterval: parseInt(import.meta.env.VITE_REFRESH_INTERVAL || '10000', 10),
  enableChartAnimation: import.meta.env.VITE_ENABLE_CHART_ANIMATION === 'true',
  debug: import.meta.env.VITE_DEBUG === 'true',
  isDevelopment: import.meta.env.MODE === 'development',
  isProduction: import.meta.env.MODE === 'production',
} as const;

export default env;
