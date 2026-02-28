import { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { isMockMode } from '@/core/config/env';

interface ProtectedRouteProps {
  children: ReactNode;
}

/**
 * Protected route wrapper that redirects to login if not authenticated.
 * In mock mode authentication is always considered valid.
 */
export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated } = useAuthStore();

  if (!isMockMode && !isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}
