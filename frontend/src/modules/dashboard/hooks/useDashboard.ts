import { useState, useEffect } from 'react';
import { DashboardSummary } from '../types/dashboard.types';
import { getDashboardSummary } from '../services/dashboardService';
import toast from 'react-hot-toast';

/**
 * Hook for fetching dashboard summary
 */
export function useDashboard() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchSummary = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await getDashboardSummary();
      setSummary(data);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Failed to fetch dashboard data';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchSummary();
  }, []);

  return {
    summary,
    isLoading,
    error,
    refetch: fetchSummary,
  };
}
