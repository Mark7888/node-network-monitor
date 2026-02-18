/**
 * Common types used across the application
 */

export type LoadingState = 'idle' | 'loading' | 'success' | 'error';

export interface PaginationParams {
  page?: number;
  limit?: number;
}

export interface TimeRangeParams {
  from?: string;
  to?: string;
}

export interface ApiError {
  message: string;
  code?: string;
  details?: unknown;
}

export interface ApiResponse<T> {
  data?: T;
  error?: ApiError;
  success: boolean;
}

export interface HealthResponse {
  status: string;
  database: string;
  uptime_seconds: number;
  version: string;
}
