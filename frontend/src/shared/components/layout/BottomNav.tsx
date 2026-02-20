import { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { LayoutDashboard, Server, Key, User } from 'lucide-react';
import ProfilePanel from './ProfilePanel';

/**
 * Bottom navigation bar for mobile devices
 */
export default function BottomNav() {
  const location = useLocation();
  const [isProfileOpen, setIsProfileOpen] = useState(false);

  const navItems = [
    { path: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { path: '/nodes', label: 'Nodes', icon: Server },
    { path: '/api-keys', label: 'Keys', icon: Key },
  ];

  const isActive = (path: string) => location.pathname.startsWith(path);

  return (
    <>
      <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-base-100 border-t border-base-300 safe-area-bottom z-50">
        <div className="flex justify-around items-center h-16">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isActive(item.path);
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`flex flex-col items-center justify-center flex-1 h-full gap-1 transition-colors ${
                  active 
                    ? 'text-primary' 
                    : 'text-base-content/60 active:text-base-content/80'
                }`}
              >
                <Icon size={22} className={active ? 'stroke-[2.5]' : 'stroke-2'} />
                <span className={`text-xs ${active ? 'font-semibold' : 'font-medium'}`}>
                  {item.label}
                </span>
              </Link>
            );
          })}
          <button
            onClick={() => setIsProfileOpen(true)}
            className="flex flex-col items-center justify-center flex-1 h-full gap-1 text-base-content/60 active:text-base-content/80 transition-colors"
          >
            <User size={22} className="stroke-2" />
            <span className="text-xs font-medium">Profile</span>
          </button>
        </div>
      </nav>

      <ProfilePanel isOpen={isProfileOpen} onClose={() => setIsProfileOpen(false)} />
    </>
  );
}
