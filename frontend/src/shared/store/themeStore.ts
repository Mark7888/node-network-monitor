import { create } from 'zustand';
import { STORAGE_KEYS } from '@/shared/utils/constants';

export type Theme = 'light' | 'dark';

interface ThemeStore {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  toggleTheme: () => void;
  initTheme: () => void;
}

/**
 * Theme store using Zustand
 */
export const useThemeStore = create<ThemeStore>((set, get) => ({
  theme: 'light',

  setTheme: (theme: Theme) => {
    localStorage.setItem(STORAGE_KEYS.THEME, theme);
    document.documentElement.setAttribute('data-theme', theme);
    set({ theme });
  },

  toggleTheme: () => {
    const currentTheme = get().theme;
    const newTheme: Theme = currentTheme === 'light' ? 'dark' : 'light';
    get().setTheme(newTheme);
  },

  initTheme: () => {
    const savedTheme = localStorage.getItem(STORAGE_KEYS.THEME) as Theme | null;
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const theme: Theme = savedTheme || (prefersDark ? 'dark' : 'light');
    
    document.documentElement.setAttribute('data-theme', theme);
    set({ theme });
  },
}));
