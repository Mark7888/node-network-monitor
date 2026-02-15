/**
 * Application environment configuration
 * Reads from environment variables set in .env files (dev)
 * or from runtime config injected by Docker (production)
 */

// TypeScript declaration for runtime config
declare global {
  interface Window {
    runtimeConfig?: {
      apiUrl: string;
      refreshInterval: number;
      enableChartAnimation: boolean;
      debug: boolean;
    };
  }
}

// Check if runtime config exists (Docker production), otherwise use build-time env vars (dev)
const runtimeConfig = typeof window !== 'undefined' ? window.runtimeConfig : undefined;

export const env = {
  apiUrl: runtimeConfig?.apiUrl || import.meta.env.VITE_API_URL || 'http://localhost:8080',
  refreshInterval: runtimeConfig?.refreshInterval || parseInt(import.meta.env.VITE_REFRESH_INTERVAL || '10000', 10),
  enableChartAnimation: runtimeConfig?.enableChartAnimation ?? (import.meta.env.VITE_ENABLE_CHART_ANIMATION === 'true'),
  debug: runtimeConfig?.debug ?? (import.meta.env.VITE_DEBUG === 'true'),
  isDevelopment: import.meta.env.MODE === 'development',
  isProduction: import.meta.env.MODE === 'production',
} as const;

export default env;
