/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  readonly VITE_REFRESH_INTERVAL: string;
  readonly VITE_ENABLE_CHART_ANIMATION: string;
  readonly VITE_DEBUG: string;
  readonly MODE: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
