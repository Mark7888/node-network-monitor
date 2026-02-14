/**
 * API Key related types
 */

export interface APIKey {
  id: string;
  name: string;
  key?: string;  // Only present when first created
  enabled: boolean;
  created_at: string;
  last_used?: string;
}

export interface CreateAPIKeyRequest {
  name: string;
}

export interface CreateAPIKeyResponse {
  id: string;
  name: string;
  key: string;
  enabled: boolean;
  created_at: string;
}

export interface UpdateAPIKeyRequest {
  enabled: boolean;
}

export interface APIKeysResponse {
  api_keys: APIKey[];
  total: number;
}
