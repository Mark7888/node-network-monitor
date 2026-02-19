import { useState, useEffect, useCallback } from 'react';
import { NodeDetails } from '../types/node.types';
import { getNodeDetails } from '../services/nodeService';
import { showToast } from '@/shared/services/toastService';

/**
 * Hook for fetching node details
 */
export function useNodeDetails(nodeId: string) {
  const [node, setNode] = useState<NodeDetails | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await getNodeDetails(nodeId);
      setNode(data);
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const errorMessage = err.response?.data?.error || 'Failed to fetch node details';
      setError(errorMessage);
      showToast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  }, [nodeId]);

  useEffect(() => {
    if (nodeId) {
      fetchData();
    }
  }, [nodeId, fetchData]);

  return {
    node,
    isLoading,
    error,
    refetch: fetchData,
  };
}
