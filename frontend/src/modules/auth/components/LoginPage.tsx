import { useState, FormEvent, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import Input from '@/shared/components/ui/Input';
import Button from '@/shared/components/ui/Button';
import { showToast } from '@/shared/services/toastService';

/**
 * Login page component
 */
export default function LoginPage() {
  const navigate = useNavigate();
  const { isAuthenticated, isLoading, error, login, clearError } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard', { replace: true });
    }
  }, [isAuthenticated, navigate]);

  // Clear error when component unmounts or inputs change
  useEffect(() => {
    return () => clearError();
  }, [clearError]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    
    if (!username || !password) {
      showToast.error('Please enter username and password');
      return;
    }

    try {
      await login({ username, password });
      showToast.success('Login successful!');
      navigate('/dashboard', { replace: true });
    } catch (err) {
      // Error is already set in store and displayed
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-base-200">
      <div className="card w-full max-w-md bg-base-100 shadow-2xl">
        <div className="card-body">
          {/* Header */}
          <div className="text-center mb-6">
            <h1 className="text-3xl font-bold text-primary mb-2">
              Speedtest Monitor
            </h1>
            <p className="text-base-content/60">
              Sign in to access the dashboard
            </p>
          </div>

          {/* Error Alert */}
          {error && (
            <div className="alert alert-error mb-4">
              <span>{error}</span>
            </div>
          )}

          {/* Login Form */}
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="Username"
              type="text"
              placeholder="Enter your username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              disabled={isLoading}
              required
              autoComplete="username"
            />

            <Input
              label="Password"
              type="password"
              placeholder="Enter your password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              disabled={isLoading}
              required
              autoComplete="current-password"
            />

            <Button
              type="submit"
              variant="primary"
              className="w-full mt-6"
              loading={isLoading}
              disabled={isLoading}
            >
              {isLoading ? 'Signing in...' : 'Sign In'}
            </Button>
          </form>

          {/* Footer */}
          <div className="text-center mt-6 text-sm text-base-content/60">
            Network Measurement Dashboard v1.0
          </div>
        </div>
      </div>
    </div>
  );
}
