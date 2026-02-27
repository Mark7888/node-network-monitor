import { apiClient } from '@/core/api/apiClient';
import { LoginRequest, LoginResponse } from '../types/auth.types';

/**
 * Authentication API service
 */

/**
 * Login with username and password
 */
export async function login(credentials: LoginRequest): Promise<LoginResponse> {
  return apiClient.login(credentials);
}

/**
 * Logout (client-side only — clear token)
 */
export function logout(): void {
  localStorage.removeItem('token');
  localStorage.removeItem('username');
}
