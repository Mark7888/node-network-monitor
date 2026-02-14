import { useMemo } from 'react';
import { AggregatedMeasurement } from '../types/measurement.types';
import {
  transformDownloadData,
  transformUploadData,
  transformPingData,
  transformJitterData,
  transformPacketLossData,
} from '../utils/dataTransform';

/**
 * Hook for transforming measurement data into chart format
 */
export function useChartData(measurements: AggregatedMeasurement[]) {
  const downloadData = useMemo(() => {
    return transformDownloadData(measurements);
  }, [measurements]);

  const uploadData = useMemo(() => {
    return transformUploadData(measurements);
  }, [measurements]);

  const pingData = useMemo(() => {
    return transformPingData(measurements);
  }, [measurements]);

  const jitterData = useMemo(() => {
    return transformJitterData(measurements);
  }, [measurements]);

  const packetLossData = useMemo(() => {
    return transformPacketLossData(measurements);
  }, [measurements]);

  return {
    downloadData,
    uploadData,
    pingData,
    jitterData,
    packetLossData,
  };
}
