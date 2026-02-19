import { useState, useEffect, useCallback } from 'react';
import { AggregatedMeasurement } from '../types/measurement.types';
import { getAggregatedMeasurements } from '../services/measurementService';
import { getTimeRange, getAggregationInterval } from '@/shared/utils/date';
import { TimeRange } from '@/shared/utils/constants';
import { showToast } from '@/shared/services/toastService';

/**
 * Hook for fetching and managing measurements data
 */
export function useMeasurements(timeRange: TimeRange, nodeIds?: string[], hideArchived?: boolean) {
  const [data, setData] = useState<AggregatedMeasurement[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const { from, to } = getTimeRange(timeRange);
      const interval = getAggregationInterval(timeRange);

      const response = await getAggregatedMeasurements({
        node_ids: nodeIds,
        from,
        to,
        interval,
        hide_archived: hideArchived,
      });

      setData(response.data || []);
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const errorMessage = err.response?.data?.error || 'Failed to fetch measurements';
      setError(errorMessage);
      showToast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  }, [timeRange, nodeIds, hideArchived]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  };
}
