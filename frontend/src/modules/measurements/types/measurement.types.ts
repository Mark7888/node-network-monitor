/**
 * Measurement related types
 */

export interface Measurement {
  id: string;
  node_id: string;
  timestamp: string;
  download_bandwidth: number;  // bytes per second
  upload_bandwidth: number;    // bytes per second
  ping: number;                // milliseconds
  jitter: number;              // milliseconds
  packet_loss: number;         // percentage (0-100)
  server_host: string;
  server_location?: string;
  result_url?: string;
}

export interface MeasurementStats {
  avg_download: number;
  avg_upload: number;
  avg_ping: number;
  avg_jitter: number;
  avg_packet_loss: number;
  min_download: number;
  min_upload: number;
  max_download: number;
  max_upload: number;
  count: number;
}

export interface AggregatedMeasurement {
  timestamp: string;
  node_id: string;
  node_name: string;
  avg_download_mbps: number;
  avg_upload_mbps: number;
  avg_ping_ms: number;
  avg_jitter_ms: number;
  avg_packet_loss: number;
  min_download_mbps?: number;
  max_download_mbps?: number;
  sample_count: number;
}

export interface MeasurementsResponse {
  measurements: Measurement[];
  total: number;
}

export interface AggregatedDataResponse {
  data: AggregatedMeasurement[];
}

export interface MeasurementQueryParams {
  node_ids?: string[] | string;  // Array or comma-separated string
  from?: string;
  to?: string;
  interval?: string;
  limit?: number;
  offset?: number;
}
