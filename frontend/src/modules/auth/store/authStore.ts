import { create } from 'zustand';
import { LoginRequest } from '../types/auth.types';
import * as authService from '../services/authService';
import { STORAGE_KEYS } from '@/shared/utils/constants';

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
 * Authentication store using Zustand
 */
export const useAuthStore = create<AuthStore>((set) => ({
  token: localStorage.getItem(STORAGE_KEYS.TOKEN),
  username: localStorage.getItem(STORAGE_KEYS.USERNAME),
  isAuthenticated: !!localStorage.getItem(STORAGE_KEYS.TOKEN),
  isLoading: false,
  error: null,

  login: async (credentials) => {
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
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'Login failed. Please try again.';
      set({
        isLoading: false,
        error: errorMessage,
      });
      throw error;
    }
  },

  logout: () => {
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
    const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
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
