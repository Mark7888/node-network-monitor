import { apiClient } from '@/core/api/apiClient';
import { HealthResponse } from '../types/common.types';

/**
 * Service for health-related API calls
 */

/**
 * Fetch health status and version info
 */
export const getHealth = async (): Promise<HealthResponse> => {
  return apiClient.getHealth();
};

export default {
  getHealth,
};
