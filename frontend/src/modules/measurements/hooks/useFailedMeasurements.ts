import { useState, useEffect, useCallback } from 'react';
import { Measurement } from '../types/measurement.types';
import { getNodeMeasurements } from '@/modules/nodes/services/nodeService';
import { getTimeRange } from '@/shared/utils/date';
import { TimeRange } from '@/shared/utils/constants';

/**
 * Hook for fetching failed measurements for a specific node
 * Returns just the timestamps of failed measurements for marking on charts
 */
export function useFailedMeasurements(timeRange: TimeRange, nodeId?: string) {
  const [failedTimestamps, setFailedTimestamps] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    if (!nodeId) {
      setFailedTimestamps([]);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const { from, to } = getTimeRange(timeRange);

      const response = await getNodeMeasurements(nodeId, {
        status: 'failed',
        from,
        to,
        limit: 10000, // Get all failed measurements within the time range
      });

      // Extract just the timestamps
      const timestamps = (response.measurements || []).map((m: Measurement) => m.timestamp);
      setFailedTimestamps(timestamps);
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const errorMessage = err.response?.data?.error || 'Failed to fetch failed measurements';
      setError(errorMessage);
      // Don't show toast for failed measurements as it's not critical
      console.error('Failed to fetch failed measurements:', errorMessage);
    } finally {
      setIsLoading(false);
    }
  }, [timeRange, nodeId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    failedTimestamps,
    isLoading,
    error,
    refetch: fetchData,
  };
}
