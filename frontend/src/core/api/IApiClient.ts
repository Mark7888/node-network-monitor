import { LoginRequest, LoginResponse } from '@/modules/auth/types/auth.types';
import { NodesResponse, NodeDetails } from '@/modules/nodes/types/node.types';
import {
  MeasurementsResponse,
  AggregatedDataResponse,
  MeasurementQueryParams,
} from '@/modules/measurements/types/measurement.types';
import {
  APIKey,
  APIKeysResponse,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
  UpdateAPIKeyRequest,
} from '@/modules/api-keys/types/apiKey.types';
import { DashboardSummary } from '@/modules/dashboard/types/dashboard.types';
import { HealthResponse } from '@/shared/types/common.types';

/**
 * Unified API client interface.
 * The real implementation talks to the backend over HTTP.
 * The mock implementation uses an in-memory state loaded from a static JSON file.
 */
export interface IApiClient {
  // Auth
  login(credentials: LoginRequest): Promise<LoginResponse>;

  // Nodes
  getNodes(): Promise<NodesResponse>;
  getNodeDetails(nodeId: string): Promise<NodeDetails>;
  archiveNode(nodeId: string, archived: boolean): Promise<void>;
  setNodeFavorite(nodeId: string, favorite: boolean): Promise<void>;
  deleteNode(nodeId: string): Promise<void>;

  // Measurements
  getMeasurements(params?: MeasurementQueryParams): Promise<MeasurementsResponse>;
  getNodeMeasurements(nodeId: string, params?: MeasurementQueryParams): Promise<MeasurementsResponse>;
  getAggregatedMeasurements(params: MeasurementQueryParams): Promise<AggregatedDataResponse>;

  // API Keys
  getAPIKeys(): Promise<APIKeysResponse>;
  createAPIKey(data: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse>;
  updateAPIKey(id: string, data: UpdateAPIKeyRequest): Promise<APIKey>;
  deleteAPIKey(id: string): Promise<void>;

  // Dashboard
  getDashboardSummary(): Promise<DashboardSummary>;

  // Health
  getHealth(): Promise<HealthResponse>;
}
