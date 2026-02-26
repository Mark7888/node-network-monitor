import { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/modules/auth/store/authStore';
import { useTheme } from '@/shared/hooks/useTheme';
import { User, Sun, Moon, ChevronDown, LogOut, Key } from 'lucide-react';

/**
 * Header component
 * Responsive: Compact on mobile, full on desktop
 */
export default function Header() {
  const { username, logout } = useAuthStore();
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false);
      }
    };

    if (isDropdownOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isDropdownOpen]);

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
          <div className="hidden sm:block relative" ref={dropdownRef}>
            <button
              onClick={() => setIsDropdownOpen(!isDropdownOpen)}
              className="flex items-center gap-2 px-3 py-2 bg-base-200 rounded-lg hover:bg-base-300 transition-colors"
            >
              <User size={16} className="text-base-content/60 md:w-[18px] md:h-[18px]" />
              <span className="text-xs md:text-sm font-medium truncate max-w-[100px] md:max-w-none">
                {username || 'Admin'}
              </span>
              <ChevronDown 
                size={16} 
                className={`transition-transform ${isDropdownOpen ? 'rotate-180' : ''}`}
              />
            </button>

            {/* Dropdown Menu */}
            {isDropdownOpen && (
              <div className="absolute right-0 mt-2 w-48 bg-base-100 rounded-lg shadow-lg border border-base-300 py-2 z-50">
                <button
                  onClick={() => {
                    setIsDropdownOpen(false);
                    navigate('/api-keys');
                  }}
                  className="flex items-center gap-3 w-full px-4 py-2 text-sm hover:bg-base-200 transition-colors"
                >
                  <Key size={18} />
                  API Keys
                </button>
                <div className="border-t border-base-300 my-1" />
                <button
                  onClick={() => {
                    setIsDropdownOpen(false);
                    logout();
                  }}
                  className="flex items-center gap-3 w-full px-4 py-2 text-sm hover:bg-base-200 transition-colors text-error"
                >
                  <LogOut size={18} />
                  Logout
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
}
