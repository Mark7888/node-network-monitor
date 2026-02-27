import { useState, useEffect } from 'react';
import { X, User as UserIcon, LogOut, Key } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/modules/auth/store/authStore';
import { getHealth } from '@/shared/services/healthService';

interface ProfilePanelProps {
  isOpen: boolean;
  onClose: () => void;
}

/**
 * Full-screen profile panel for mobile
 * Shows user profile and logout option
 */
export default function ProfilePanel({ isOpen, onClose }: ProfilePanelProps) {
  const { username, logout } = useAuthStore();
  const navigate = useNavigate();
  const [version, setVersion] = useState<string>('');

  useEffect(() => {
    const fetchVersion = async () => {
      try {
        const health = await getHealth();
        setVersion(health.version);
      } catch {
        // Silently fail
      }
    };
    fetchVersion();
  }, []);

  if (!isOpen) return null;

  const handleLogout = () => {
    onClose();
    logout();
  };

  const handleAPIKeys = () => {
    onClose();
    navigate('/api-keys');
  };

  return (
    <div className="fixed inset-0 z-50 bg-base-100 md:hidden">
      {/* Header */}
      <div className="flex justify-between items-center px-4 py-4 border-b border-base-300">
        <h2 className="text-lg font-bold">Profile</h2>
        <button
          onClick={onClose}
          className="btn btn-ghost btn-circle btn-sm"
          aria-label="Close profile"
        >
          <X size={20} />
        </button>
      </div>

      {/* Content */}
      <div className="flex flex-col items-center pt-12 px-6">
        {/* Profile Avatar */}
        <div className="w-24 h-24 rounded-full bg-primary flex items-center justify-center mb-4">
          <UserIcon size={48} className="text-primary-content" />
        </div>

        {/* Username */}
        <h3 className="text-xl font-semibold mb-2">{username || 'Admin'}</h3>

        {/* Actions Section */}
        <div className="w-full max-w-sm space-y-4">
          <div className="divider"/>
          <button
            onClick={handleAPIKeys}
            className="btn btn-outline btn-block justify-start gap-3"
          >
            <Key size={20} />
            API Keys
          </button>          
          <button
            onClick={handleLogout}
            className="btn btn-outline btn-error btn-block justify-start gap-3"
          >
            <LogOut size={20} />
            Logout
          </button>
        </div>

        {version && (
          <p className="mt-8 text-xs text-base-content/40 font-mono">
            Version: {version}
          </p>
        )}
      </div>
    </div>
  );
}
