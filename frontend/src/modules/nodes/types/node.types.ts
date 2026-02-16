import { NodeStatus } from '@/shared/utils/constants';

/**
 * Node related types
 */

export interface Node {
  id: string;
  name: string;
  status: NodeStatus;
  first_seen: string;
  last_seen: string;
  last_alive: string;
  created_at: string;
  updated_at: string;
  location?: string;
}

export interface LatestMeasurement {
  timestamp: string;
  download_mbps: number;
  upload_mbps: number;
  ping_ms: number;
}

export interface NodeStatistics {
  avg_download_mbps: number;
  avg_upload_mbps: number;
  avg_ping_ms: number;
  avg_jitter_ms: number;
  avg_packet_loss: number;
  success_rate_24h: number;       // Success rate for last 24 hours (0-100)
  success_count_24h: number;      // Successful measurements in last 24h
  failed_count_24h: number;       // Failed measurements in last 24h
}

export interface NodeDetails extends Node {
  total_measurements: number;
  failed_test_count: number;
  latest_measurement?: LatestMeasurement;
  statistics: NodeStatistics;
}

export interface NodesResponse {
  nodes: Node[];
  total: number;
}
