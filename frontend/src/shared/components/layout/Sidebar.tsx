import { Link, useLocation } from 'react-router-dom';
import { LayoutDashboard, Server, Key, LogOut } from 'lucide-react';
import { useAuthStore } from '@/modules/auth/store/authStore';

/**
 * Sidebar navigation component
 * Hidden on mobile (< md breakpoint), shown on desktop
 */
export default function Sidebar() {
  const location = useLocation();
  const { logout } = useAuthStore();

  const navItems = [
    { path: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { path: '/nodes', label: 'Nodes', icon: Server },
    { path: '/api-keys', label: 'API Keys', icon: Key },
  ];

  const isActive = (path: string) => location.pathname.startsWith(path);

  return (
    <aside className="hidden md:flex w-64 bg-base-100 shadow-md min-h-screen flex-col">
      {/* Logo */}
      <div className="p-6 border-b border-base-300">
        <h1 className="text-xl font-bold text-primary">
          Speedtest Monitor
        </h1>
      </div>

      {/* Navigation */}
      <nav className="flex-1 p-4">
        <ul className="menu menu-compact">
          {navItems.map((item) => {
            const Icon = item.icon;
            return (
              <li key={item.path}>
                <Link
                  to={item.path}
                  className={isActive(item.path) ? 'active' : ''}
                >
                  <Icon size={18} />
                  {item.label}
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>

      {/* Logout */}
      <div className="p-4 border-t border-base-300">
        <button
          onClick={logout}
          className="btn btn-ghost btn-block justify-start"
        >
          <LogOut size={18} />
          Logout
        </button>
      </div>
    </aside>
  );
}
