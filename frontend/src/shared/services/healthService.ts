import api from '@/core/api/axiosConfig';
import { HealthResponse } from '../types/common.types';

/**
 * Service for health-related API calls
 */

/**
 * Fetch health status and version info
 */
export const getHealth = async (): Promise<HealthResponse> => {
  const response = await api.get<HealthResponse>('/health');
  return response.data;
};

export default {
  getHealth,
};
