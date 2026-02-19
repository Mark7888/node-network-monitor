import { useState, useCallback, useEffect } from 'react';
import { getNodeMeasurements } from '@/modules/nodes/services/nodeService';
import { Measurement } from '../types/measurement.types';
import { showToast } from '@/shared/services/toastService';

/**
 * Hook for fetching node measurements with pagination and filtering
 */
export function useNodeMeasurements(nodeId: string | undefined, status: 'all' | 'successful' | 'failed' = 'all') {
  const [measurements, setMeasurements] = useState<Measurement[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [limit] = useState(50); // Fixed page size

  const fetchData = useCallback(async (currentPage: number) => {
    if (!nodeId) return;

    setIsLoading(true);
    setError(null);

    try {
      const response = await getNodeMeasurements(nodeId, {
        page: currentPage,
        limit,
        status,
      });

      setMeasurements(response.measurements || []);
      setTotal(response.total || 0);
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const errorMessage = err.response?.data?.error || 'Failed to fetch measurements';
      setError(errorMessage);
      showToast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  }, [nodeId, limit, status]);

  useEffect(() => {
    fetchData(page);
  }, [fetchData, page]);

  const goToPage = useCallback((newPage: number) => {
    setPage(newPage);
  }, []);

  const nextPage = useCallback(() => {
    const maxPage = Math.ceil(total / limit);
    if (page < maxPage) {
      setPage(page + 1);
    }
  }, [page, total, limit]);

  const prevPage = useCallback(() => {
    if (page > 1) {
      setPage(page - 1);
    }
  }, [page]);

  const refetch = useCallback(() => {
    fetchData(page);
  }, [fetchData, page]);

  return {
    measurements,
    isLoading,
    error,
    total,
    page,
    limit,
    totalPages: Math.ceil(total / limit),
    goToPage,
    nextPage,
    prevPage,
    refetch,
  };
}
