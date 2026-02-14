import { StrictMode, useEffect } from 'react';
import ReactDOM from 'react-dom/client';
import { AppRouter } from './core/router/routes';
import { setupInterceptors } from './core/api/interceptors';
import { useThemeStore } from './shared/store/themeStore';
import { Toaster } from 'react-hot-toast';
import './index.css';

/**
 * Setup API interceptors with router navigation
 */
export function App() {
  useEffect(() => {
    // Initialize theme
    useThemeStore.getState().initTheme();

    // Setup API interceptors
    setupInterceptors(() => {
      // Redirect to login on unauthorized
      window.location.href = '/login';
    });
  }, []);

  return (
    <>
      <AppRouter />
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 4000,
          style: {
            background: '#363636',
            color: '#fff',
          },
          success: {
            duration: 3000,
            iconTheme: {
              primary: '#10b981',
              secondary: '#fff',
            },
          },
          error: {
            duration: 4000,
            iconTheme: {
              primary: '#ef4444',
              secondary: '#fff',
            },
          },
        }}
      />
    </>
  );
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
