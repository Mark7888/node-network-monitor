import api from '@/core/api/axiosConfig';
import {
  APIKeysResponse,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
  UpdateAPIKeyRequest,
  APIKey,
} from '../types/apiKey.types';

/**
 * API Key service
 */

/**
 * Get all API keys
 */
export async function getAPIKeys(): Promise<APIKeysResponse> {
  const response = await api.get<APIKeysResponse>('/api/v1/admin/api-keys');
  return response.data;
}

/**
 * Create a new API key
 */
export async function createAPIKey(data: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse> {
  const response = await api.post<CreateAPIKeyResponse>('/api/v1/admin/api-keys', data);
  return response.data;
}

/**
 * Update an API key (enable/disable)
 */
export async function updateAPIKey(
  id: string,
  data: UpdateAPIKeyRequest
): Promise<APIKey> {
  const response = await api.patch<APIKey>(`/api/v1/admin/api-keys/${id}`, data);
  return response.data;
}

/**
 * Delete an API key
 */
export async function deleteAPIKey(id: string): Promise<void> {
  await api.delete(`/api/v1/admin/api-keys/${id}`);
}
