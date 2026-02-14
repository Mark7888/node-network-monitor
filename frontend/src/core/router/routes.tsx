import { lazy, Suspense } from 'react';
import { createBrowserRouter, Navigate, RouterProvider } from 'react-router-dom';
import Layout from '@/shared/components/layout/Layout';
import ProtectedRoute from '@/modules/auth/components/ProtectedRoute';
import Spinner from '@/shared/components/ui/Spinner';

// Lazy load pages for code splitting
const LoginPage = lazy(() => import('@/modules/auth/components/LoginPage'));
const DashboardPage = lazy(() => import('@/modules/dashboard/components/DashboardPage'));
const NodesPage = lazy(() => import('@/modules/nodes/components/NodesPage'));
const NodeDetailsPage = lazy(() => import('@/modules/nodes/components/NodeDetailsPage'));
const APIKeysPage = lazy(() => import('@/modules/api-keys/components/APIKeysPage'));

/**
 * Application routes configuration
 */
const router = createBrowserRouter([
  {
    path: '/login',
    element: (
      <Suspense fallback={<Spinner fullScreen />}>
        <LoginPage />
      </Suspense>
    ),
  },
  {
    path: '/',
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
    ],
  },
]);

/**
 * Router component wrapper
 */
export function AppRouter() {
  return <RouterProvider router={router} />;
}
