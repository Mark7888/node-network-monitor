import { apiClient } from '@/core/api/apiClient';
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

export async function getAPIKeys(): Promise<APIKeysResponse> {
  return apiClient.getAPIKeys();
}

export async function createAPIKey(data: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse> {
  return apiClient.createAPIKey(data);
}

export async function updateAPIKey(
  id: string,
  data: UpdateAPIKeyRequest,
): Promise<APIKey> {
  return apiClient.updateAPIKey(id, data);
}

export async function deleteAPIKey(id: string): Promise<void> {
  return apiClient.deleteAPIKey(id);
}
