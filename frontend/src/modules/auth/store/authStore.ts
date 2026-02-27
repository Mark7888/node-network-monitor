import { create } from 'zustand';
import { LoginRequest } from '../types/auth.types';
import * as authService from '../services/authService';
import { STORAGE_KEYS } from '@/shared/utils/constants';
import { isMockMode } from '@/core/config/env';

interface AuthStore {
  token: string | null;
  username: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;

  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
  checkAuth: () => void;
  clearError: () => void;
}

/**
 * Authentication store using Zustand.
 * In mock mode the user is automatically considered authenticated (no login required).
 */
export const useAuthStore = create<AuthStore>((set) => ({
  token:           isMockMode ? 'mock-token' : localStorage.getItem(STORAGE_KEYS.TOKEN),
  username:        isMockMode ? 'demo'       : localStorage.getItem(STORAGE_KEYS.USERNAME),
  isAuthenticated: isMockMode ? true         : !!localStorage.getItem(STORAGE_KEYS.TOKEN),
  isLoading: false,
  error: null,

  login: async (credentials) => {
    if (isMockMode) {
      set({ token: 'mock-token', username: 'demo', isAuthenticated: true, isLoading: false, error: null });
      return;
    }
    set({ isLoading: true, error: null });
    try {
      const data = await authService.login(credentials);
      localStorage.setItem(STORAGE_KEYS.TOKEN, data.token);
      localStorage.setItem(STORAGE_KEYS.USERNAME, data.username);
      set({
        token: data.token,
        username: data.username,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const errorMessage = err.response?.data?.error || 'Login failed. Please try again.';
      set({
        isLoading: false,
        error: errorMessage,
      });
      throw error;
    }
  },

  logout: () => {
    if (isMockMode) {
      // In mock mode, skip backend logout and navigation but ensure any real
      // credentials (if present) are cleared from storage and state.
      localStorage.removeItem(STORAGE_KEYS.TOKEN);
      localStorage.removeItem(STORAGE_KEYS.USERNAME);
      set({
        token: 'mock-token',
        username: 'demo',
        isAuthenticated: true,
        error: null,
      });
      return;
    }
    authService.logout();
    set({
      token: null,
      username: null,
      isAuthenticated: false,
      error: null,
    });
    window.location.href = '/login';
  },

  checkAuth: () => {
    if (isMockMode) return;
    const token    = localStorage.getItem(STORAGE_KEYS.TOKEN);
    const username = localStorage.getItem(STORAGE_KEYS.USERNAME);
    set({
      token,
      username,
      isAuthenticated: !!token,
    });
  },

  clearError: () => {
    set({ error: null });
  },
}));
