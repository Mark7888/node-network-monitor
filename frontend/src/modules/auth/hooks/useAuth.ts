import { useAuthStore } from '../store/authStore';

/**
 * Custom hook for authentication
 */
export function useAuth() {
  const {
    isAuthenticated,
    isLoading,
    error,
    username,
    login,
    logout,
    checkAuth,
    clearError,
  } = useAuthStore();

  return {
    isAuthenticated,
    isLoading,
    error,
    username,
    login,
    logout,
    checkAuth,
    clearError,
  };
}
