import { apiClient } from '@/core/api/apiClient';
import { DashboardSummary } from '../types/dashboard.types';

/**
 * Dashboard API service
 */

export async function getDashboardSummary(): Promise<DashboardSummary> {
  return apiClient.getDashboardSummary();
}
