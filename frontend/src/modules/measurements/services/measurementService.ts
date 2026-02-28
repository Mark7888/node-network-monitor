import { apiClient } from '@/core/api/apiClient';
import {
  MeasurementsResponse,
  AggregatedDataResponse,
  MeasurementQueryParams,
} from '../types/measurement.types';

/**
 * Measurement API service
 */

export async function getAggregatedMeasurements(
  params: MeasurementQueryParams,
): Promise<AggregatedDataResponse> {
  return apiClient.getAggregatedMeasurements(params);
}

export async function getMeasurements(
  params?: MeasurementQueryParams,
): Promise<MeasurementsResponse> {
  return apiClient.getMeasurements(params);
}
