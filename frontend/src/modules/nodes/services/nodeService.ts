import api from '@/core/api/axiosConfig';
import { NodesResponse, NodeDetails } from '../types/node.types';
import { MeasurementsResponse, MeasurementQueryParams } from '@/modules/measurements/types/measurement.types';

/**
 * Node API service
 */

/**
 * Get all nodes
 */
export async function getNodes(): Promise<NodesResponse> {
  const response = await api.get<NodesResponse>('/api/v1/admin/nodes');
  return response.data;
}

/**
 * Get single node details
 */
export async function getNodeDetails(nodeId: string): Promise<NodeDetails> {
  const response = await api.get<NodeDetails>(`/api/v1/admin/nodes/${nodeId}`);
  return response.data;
}

/**
 * Get measurements for a specific node
 */
export async function getNodeMeasurements(
  nodeId: string,
  params?: MeasurementQueryParams
): Promise<MeasurementsResponse> {
  const response = await api.get<MeasurementsResponse>(
    `/api/v1/admin/nodes/${nodeId}/measurements`,
    { params }
  );
  return response.data;
}
