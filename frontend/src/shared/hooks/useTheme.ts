import { useThemeStore } from '@/shared/store/themeStore';

/**
 * Custom hook for theme management
 */
export function useTheme() {
  const { theme, setTheme, toggleTheme, initTheme } = useThemeStore();

  return {
    theme,
    setTheme,
    toggleTheme,
    initTheme,
    isDark: theme === 'dark',
    isLight: theme === 'light',
  };
}
