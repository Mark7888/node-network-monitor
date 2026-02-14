import api from '@/core/api/axiosConfig';
import { DashboardSummary } from '../types/dashboard.types';

/**
 * Dashboard API service
 */

/**
 * Get dashboard summary statistics
 */
export async function getDashboardSummary(): Promise<DashboardSummary> {
  const response = await api.get<DashboardSummary>('/api/v1/admin/dashboard');
  return response.data;
}
