import { useAuthStore } from '@/modules/auth/store/authStore';
import { useTheme } from '@/shared/hooks/useTheme';
import { User, Sun, Moon } from 'lucide-react';

/**
 * Header component
 */
export default function Header() {
  const { username } = useAuthStore();
  const { theme, toggleTheme } = useTheme();

  return (
    <header className="bg-base-100 shadow-sm px-6 py-4">
      <div className="flex justify-between items-center">
        <div>
          {/* Breadcrumb or page title can go here */}
        </div>
        
        <div className="flex items-center gap-3">
          {/* Theme Toggle */}
          <button
            onClick={toggleTheme}
            className="btn btn-ghost btn-circle"
            aria-label="Toggle theme"
          >
            {theme === 'light' ? (
              <Moon size={20} />
            ) : (
              <Sun size={20} />
            )}
          </button>

          {/* User Info */}
          <div className="flex items-center gap-2 px-3 py-2 bg-base-200 rounded-lg">
            <User size={18} className="text-base-content/60" />
            <span className="text-sm font-medium">{username || 'Admin'}</span>
          </div>
        </div>
      </div>
    </header>
  );
}
