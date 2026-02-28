import api from './axiosConfig';
import { IApiClient } from './IApiClient';
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
 * Real API client — delegates all calls to the backend over HTTP via axios.
 */
export const realApiClient: IApiClient = {
  // ── Auth ──────────────────────────────────────────────────────────────────

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await api.post<LoginResponse>('/api/v1/admin/login', credentials);
    return response.data;
  },

  // ── Nodes ─────────────────────────────────────────────────────────────────

  async getNodes(): Promise<NodesResponse> {
    const response = await api.get<NodesResponse>('/api/v1/admin/nodes');
    return response.data;
  },

  async getNodeDetails(nodeId: string): Promise<NodeDetails> {
    const response = await api.get<NodeDetails>(`/api/v1/admin/nodes/${nodeId}`);
    return response.data;
  },

  async archiveNode(nodeId: string, archived: boolean): Promise<void> {
    await api.patch(`/api/v1/admin/nodes/${nodeId}/archive`, { archived });
  },

  async setNodeFavorite(nodeId: string, favorite: boolean): Promise<void> {
    await api.patch(`/api/v1/admin/nodes/${nodeId}/favorite`, { favorite });
  },

  async deleteNode(nodeId: string): Promise<void> {
    await api.delete(`/api/v1/admin/nodes/${nodeId}`);
  },

  // ── Measurements ──────────────────────────────────────────────────────────

  async getMeasurements(params?: MeasurementQueryParams): Promise<MeasurementsResponse> {
    const response = await api.get<MeasurementsResponse>('/api/v1/admin/measurements', { params });
    return response.data;
  },

  async getNodeMeasurements(
    nodeId: string,
    params?: MeasurementQueryParams,
  ): Promise<MeasurementsResponse> {
    const response = await api.get<MeasurementsResponse>(
      `/api/v1/admin/nodes/${nodeId}/measurements`,
      { params },
    );
    return response.data;
  },

  async getAggregatedMeasurements(
    params: MeasurementQueryParams,
  ): Promise<AggregatedDataResponse> {
    const queryParams: Record<string, string | string[] | number> = {};

    if (params.node_ids) {
      queryParams.node_ids = Array.isArray(params.node_ids) ? params.node_ids : [params.node_ids];
    }
    if (params.from)            queryParams.from           = params.from;
    if (params.to)              queryParams.to             = params.to;
    if (params.interval)        queryParams.interval       = params.interval;
    if (params.limit)           queryParams.limit          = params.limit;
    if (params.offset)          queryParams.offset         = params.offset;
    if (params.hide_archived !== undefined) {
      queryParams.hide_archived = params.hide_archived ? 'true' : 'false';
    }

    const response = await api.get<AggregatedDataResponse>(
      '/api/v1/admin/measurements/aggregate',
      {
        params: queryParams,
        paramsSerializer: { indexes: null },
      },
    );
    return response.data;
  },

  // ── API Keys ──────────────────────────────────────────────────────────────

  async getAPIKeys(): Promise<APIKeysResponse> {
    const response = await api.get<APIKeysResponse>('/api/v1/admin/api-keys');
    return response.data;
  },

  async createAPIKey(data: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse> {
    const response = await api.post<CreateAPIKeyResponse>('/api/v1/admin/api-keys', data);
    return response.data;
  },

  async updateAPIKey(id: string, data: UpdateAPIKeyRequest): Promise<APIKey> {
    const response = await api.patch<APIKey>(`/api/v1/admin/api-keys/${id}`, data);
    return response.data;
  },

  async deleteAPIKey(id: string): Promise<void> {
    await api.delete(`/api/v1/admin/api-keys/${id}`);
  },

  // ── Dashboard ─────────────────────────────────────────────────────────────

  async getDashboardSummary(): Promise<DashboardSummary> {
    const response = await api.get<DashboardSummary>('/api/v1/admin/dashboard');
    return response.data;
  },

  // ── Health ────────────────────────────────────────────────────────────────

  async getHealth(): Promise<HealthResponse> {
    const response = await api.get<HealthResponse>('/health');
    return response.data;
  },
};
