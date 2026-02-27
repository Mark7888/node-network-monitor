import { apiClient } from '@/core/api/apiClient';
import { NodesResponse, NodeDetails } from '../types/node.types';
import { MeasurementsResponse, MeasurementQueryParams } from '@/modules/measurements/types/measurement.types';

/**
 * Node API service
 */

export async function getNodes(): Promise<NodesResponse> {
  return apiClient.getNodes();
}

export async function getNodeDetails(nodeId: string): Promise<NodeDetails> {
  return apiClient.getNodeDetails(nodeId);
}

export async function getNodeMeasurements(
  nodeId: string,
  params?: MeasurementQueryParams,
): Promise<MeasurementsResponse> {
  return apiClient.getNodeMeasurements(nodeId, params);
}

export async function archiveNode(nodeId: string, archived: boolean): Promise<void> {
  return apiClient.archiveNode(nodeId, archived);
}

export async function setNodeFavorite(nodeId: string, favorite: boolean): Promise<void> {
  return apiClient.setNodeFavorite(nodeId, favorite);
}

export async function deleteNode(nodeId: string): Promise<void> {
  return apiClient.deleteNode(nodeId);
}
