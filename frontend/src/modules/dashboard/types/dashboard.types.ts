/**
 * Dashboard related types
 */

export interface DashboardStats24h {
  download_mbps: number;
  upload_mbps: number;
  ping_ms: number;
  jitter_ms: number;
  packet_loss: number;
}

export interface DashboardSummary {
  total_nodes: number;
  active_nodes: number;
  unreachable_nodes: number;
  total_measurements: number;
  measurements_last_24h: number;
  last_measurement?: string;
  average_stats_24h?: DashboardStats24h;
}
