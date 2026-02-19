/**
 * Measurement related types
 */

export interface Measurement {
  id: number;
  node_id: string;
  timestamp: string;
  created_at: string;
  
  // Ping metrics
  ping_jitter?: number;
  ping_latency?: number;
  ping_low?: number;
  ping_high?: number;
  
  // Download metrics
  download_bandwidth?: number;      // bytes per second
  download_bytes?: number;
  download_elapsed?: number;        // milliseconds
  download_latency_iqm?: number;
  download_latency_low?: number;
  download_latency_high?: number;
  download_latency_jitter?: number;
  
  // Upload metrics
  upload_bandwidth?: number;        // bytes per second
  upload_bytes?: number;
  upload_elapsed?: number;          // milliseconds
  upload_latency_iqm?: number;
  upload_latency_low?: number;
  upload_latency_high?: number;
  upload_latency_jitter?: number;
  
  // Network info
  packet_loss?: number;             // percentage (0-100)
  isp?: string;
  interface_internal_ip?: string;
  interface_name?: string;
  interface_mac?: string;
  interface_is_vpn?: boolean;
  interface_external_ip?: string;
  
  // Server info
  server_id?: number;
  server_host?: string;
  server_port?: number;
  server_name?: string;
  server_location?: string;
  server_country?: string;
  server_ip?: string;
  
  // Result info
  result_id?: string;
  result_url?: string;

  // Failed measurement info
  is_failed: boolean;
  error_message?: string;
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
  page?: number;
  limit?: number;
}

export interface AggregatedDataResponse {
  data: AggregatedMeasurement[];
}

export interface MeasurementQueryParams {
  node_ids?: string[] | string;  // Array or comma-separated string
  from?: string;
  to?: string;
  interval?: string;
  page?: number;
  limit?: number;
  offset?: number;
  status?: 'all' | 'successful' | 'failed';  // Filter by measurement status
  hide_archived?: boolean;  // Hide archived nodes from results
}
