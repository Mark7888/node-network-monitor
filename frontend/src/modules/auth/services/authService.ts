import api from '@/core/api/axiosConfig';
import { LoginRequest, LoginResponse } from '../types/auth.types';

/**
 * Authentication API service
 */

/**
 * Login with username and password
 */
export async function login(credentials: LoginRequest): Promise<LoginResponse> {
  const response = await api.post<LoginResponse>('/api/v1/admin/login', credentials);
  return response.data;
}

/**
 * Logout (client-side only - clear token)
 */
export function logout(): void {
  // Could call a logout endpoint here if needed
  localStorage.removeItem('token');
  localStorage.removeItem('username');
}
