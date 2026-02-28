import { lazy, Suspense } from 'react';
import { createBrowserRouter, createHashRouter, Navigate, RouterProvider } from 'react-router-dom';
import Layout from '@/shared/components/layout/Layout';
import ProtectedRoute from '@/modules/auth/components/ProtectedRoute';
import Spinner from '@/shared/components/ui/Spinner';
import ErrorPage from '@/shared/components/ErrorPage';
import { isMockMode } from '@/core/config/env';

// Lazy load pages for code splitting
const LoginPage = lazy(() => import('@/modules/auth/components/LoginPage'));
const DashboardPage = lazy(() => import('@/modules/dashboard/components/DashboardPage'));
const NodesPage = lazy(() => import('@/modules/nodes/components/NodesPage'));
const NodeDetailsPage = lazy(() => import('@/modules/nodes/components/NodeDetailsPage'));
const APIKeysPage = lazy(() => import('@/modules/api-keys/components/APIKeysPage'));

/**
 * Application routes configuration.
 * In mock mode the /login route redirects straight to the dashboard.
 *
 * Mock builds use createHashRouter so the app works correctly on GitHub Pages
 * (a static file host that can't serve History-API paths). Hash-based URLs
 * (#/dashboard) are never sent to the server, so GitHub Pages always returns
 * index.html regardless of the current route.
 * The real/production build uses the standard browser-history router.
 */
const routeConfig = [
  {
    path: '/login',
    errorElement: <ErrorPage />,
    element: isMockMode
      ? <Navigate to="/dashboard" replace />
      : (
        <Suspense fallback={<Spinner fullScreen />}>
          <LoginPage />
        </Suspense>
      ),
  },
  {
    path: '/',
    errorElement: <ErrorPage />,
    element: (
      <ProtectedRoute>
        <Layout />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <Navigate to="/dashboard" replace />,
      },
      {
        path: 'dashboard',
        element: (
          <Suspense fallback={<Spinner />}>
            <DashboardPage />
          </Suspense>
        ),
      },
      {
        path: 'nodes',
        element: (
          <Suspense fallback={<Spinner />}>
            <NodesPage />
          </Suspense>
        ),
      },
      {
        path: 'nodes/:id',
        element: (
          <Suspense fallback={<Spinner />}>
            <NodeDetailsPage />
          </Suspense>
        ),
      },
      {
        path: 'api-keys',
        element: (
          <Suspense fallback={<Spinner />}>
            <APIKeysPage />
          </Suspense>
        ),
      },
      // Catch-all — renders the error page for any unknown child path
      {
        path: '*',
        element: <ErrorPage />,
      },
    ],
  },
];

const router = isMockMode
  ? createHashRouter(routeConfig)
  : createBrowserRouter(routeConfig);

/**
 * Router component wrapper
 */
export function AppRouter() {
  return <RouterProvider router={router} />;
}
