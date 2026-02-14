import api from '@/core/api/axiosConfig';
import {
  MeasurementsResponse,
  AggregatedDataResponse,
  MeasurementQueryParams,
} from '../types/measurement.types';

/**
 * Measurement API service
 */

/**
 * Get aggregated measurements across nodes
 */
export async function getAggregatedMeasurements(
  params: MeasurementQueryParams
): Promise<AggregatedDataResponse> {
  // Build query params manually to handle array serialization
  const queryParams: Record<string, string | string[] | number> = {};
  
  // Add node_ids as repeated params if present (Gin expects ?node_ids=uuid1&node_ids=uuid2)
  // Axios will handle the array correctly with paramsSerializer
  if (params.node_ids) {
    queryParams.node_ids = Array.isArray(params.node_ids) ? params.node_ids : [params.node_ids];
  }
  
  // Add other required params
  if (params.from) queryParams.from = params.from;
  if (params.to) queryParams.to = params.to;
  if (params.interval) queryParams.interval = params.interval;
  if (params.limit) queryParams.limit = params.limit;
  if (params.offset) queryParams.offset = params.offset;
  
  const response = await api.get<AggregatedDataResponse>(
    '/api/v1/admin/measurements/aggregate',
    { 
      params: queryParams,
      paramsSerializer: {
        indexes: null, // This makes axios send arrays as ?key=val1&key=val2
      },
    }
  );
  return response.data;
}

/**
 * Get raw measurements
 */
export async function getMeasurements(
  params?: MeasurementQueryParams
): Promise<MeasurementsResponse> {
  const response = await api.get<MeasurementsResponse>(
    '/api/v1/admin/measurements',
    { params }
  );
  return response.data;
}
