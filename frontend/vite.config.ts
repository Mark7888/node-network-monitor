import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [react()],
  /**
   * For the mock/demo build we read VITE_BASE_PATH from the environment so the
   * GitHub Actions workflow can inject the GitHub Pages sub-path at build time
   * (e.g. /network-measure-app/).  Falls back to '/' for local testing.
   */
  base: mode === 'mock' ? (process.env.VITE_BASE_PATH ?? '/') : '/',
  /**
   * Bake VITE_MOCK_MODE directly into the bundle for mock builds so Rollup
   * can statically analyse and tree-shake all real-API code paths.
   * This is more reliable than relying solely on .env.mock being loaded.
   */
  define: mode === 'mock' ? {
    'import.meta.env.VITE_MOCK_MODE': JSON.stringify('true'),
  } : {},
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@/modules': fileURLToPath(new URL('./src/modules', import.meta.url)),
      '@/shared': fileURLToPath(new URL('./src/shared', import.meta.url)),
      '@/core': fileURLToPath(new URL('./src/core', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    host: true,
  },
}))
