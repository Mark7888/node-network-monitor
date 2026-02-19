import { useState, useCallback } from 'react';
import {
  getAPIKeys,
  createAPIKey,
  updateAPIKey,
  deleteAPIKey,
} from '../services/apiKeyService';
import { APIKey, CreateAPIKeyResponse } from '../types/apiKey.types';
import { showToast } from '@/shared/services/toastService';

/**
 * Hook for managing API keys
 */
export function useAPIKeys() {
  const [apiKeys, setAPIKeys] = useState<APIKey[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAPIKeys = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await getAPIKeys();
      setAPIKeys(data.api_keys || []);
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const errorMessage = err.response?.data?.error || 'Failed to fetch API keys';
      setError(errorMessage);
      showToast.error(errorMessage);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createKey = useCallback(async (name: string): Promise<CreateAPIKeyResponse | null> => {
    try {
      const data = await createAPIKey({ name });
      showToast.success('API key created successfully!');
      await fetchAPIKeys(); // Refresh list
      return data;
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const errorMessage = err.response?.data?.error || 'Failed to create API key';
      showToast.error(errorMessage);
      return null;
    }
  }, [fetchAPIKeys]);

  const toggleKey = useCallback(async (id: string, enabled: boolean) => {
    try {
      await updateAPIKey(id, { enabled });
      showToast.success(`API key ${enabled ? 'enabled' : 'disabled'}`);
      await fetchAPIKeys(); // Refresh list
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const errorMessage = err.response?.data?.error || 'Failed to update API key';
      showToast.error(errorMessage);
    }
  }, [fetchAPIKeys]);

  const deleteKey = useCallback(async (id: string) => {
    try {
      await deleteAPIKey(id);
      showToast.success('API key deleted');
      await fetchAPIKeys(); // Refresh list
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      const errorMessage = err.response?.data?.error || 'Failed to delete API key';
      showToast.error(errorMessage);
    }
  }, [fetchAPIKeys]);

  return {
    apiKeys,
    isLoading,
    error,
    fetchAPIKeys,
    createKey,
    toggleKey,
    deleteKey,
  };
}
