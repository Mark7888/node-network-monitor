import { useState, useEffect } from 'react';
import { NodeDetails } from '../types/node.types';
import { getNodeDetails } from '../services/nodeService';
import toast from 'react-hot-toast';

/**
 * Hook for fetching node details
 */
export function useNodeDetails(nodeId: string) {
  const [node, setNode] = useState<NodeDetails | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await getNodeDetails(nodeId);
      setNode(data);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Failed to fetch node details';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (nodeId) {
      fetchData();
    }
  }, [nodeId]);

  return {
    node,
    isLoading,
    error,
    refetch: fetchData,
  };
}
