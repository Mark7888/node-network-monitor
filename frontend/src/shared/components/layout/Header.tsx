import { useAuthStore } from '@/modules/auth/store/authStore';
import { useTheme } from '@/shared/hooks/useTheme';
import { User, Sun, Moon } from 'lucide-react';

/**
 * Header component
 * Responsive: Compact on mobile, full on desktop
 */
export default function Header() {
  const { username } = useAuthStore();
  const { theme, toggleTheme } = useTheme();

  return (
    <header className="bg-base-100 shadow-sm px-3 md:px-6 py-3 md:py-4">
      <div className="flex justify-between items-center">
        <div className="md:hidden">
          {/* Mobile: Show app name */}
          <h1 className="text-base font-bold text-primary">
            Speedtest Monitor
          </h1>
        </div>
        <div className="hidden md:block">
          {/* Desktop: Breadcrumb or page title can go here */}
        </div>
        
        <div className="flex items-center gap-2 md:gap-3">
          {/* Theme Toggle */}
          <button
            onClick={toggleTheme}
            className="btn btn-ghost btn-circle btn-sm md:btn-md"
            aria-label="Toggle theme"
          >
            {theme === 'light' ? (
              <Moon size={18} className="md:w-5 md:h-5" />
            ) : (
              <Sun size={18} className="md:w-5 md:h-5" />
            )}
          </button>

          {/* User Info - Hidden on small mobile, shown on larger screens */}
          <div className="hidden sm:flex items-center gap-2 px-3 py-2 bg-base-200 rounded-lg">
            <User size={16} className="text-base-content/60 md:w-[18px] md:h-[18px]" />
            <span className="text-xs md:text-sm font-medium truncate max-w-[100px] md:max-w-none">
              {username || 'Admin'}
            </span>
          </div>
        </div>
      </div>
    </header>
  );
}
