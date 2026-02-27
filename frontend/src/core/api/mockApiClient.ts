import { IApiClient } from './IApiClient';
import { LoginRequest, LoginResponse } from '@/modules/auth/types/auth.types';
import { Node, NodeDetails, NodesResponse } from '@/modules/nodes/types/node.types';
import {
  Measurement,
  MeasurementsResponse,
  AggregatedDataResponse,
  AggregatedMeasurement,
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

// ── Internal state ────────────────────────────────────────────────────────────

interface MockState {
  nodes: Node[];
  measurements: Record<string, Measurement[]>;
  apiKeys: APIKey[];
  nextApiKeyId: number;
}

let state: MockState = {
  nodes: [],
  measurements: {},
  apiKeys: [],
  nextApiKeyId: 100,
};

// ── Initialisation ────────────────────────────────────────────────────────────

/**
 * Timestamp fields that may be stored as relative second offsets (negative integer = seconds before now).
 * Any other value type (string, null) is left untouched.
 */
const TIMESTAMP_FIELDS = new Set([
  'timestamp', 'created_at', 'updated_at',
  'first_seen', 'last_seen', 'last_alive', 'last_used',
]);

/**
 * Recursively walks the parsed JSON and converts any TIMESTAMP_FIELDS entry that is
 * a plain number (relative offset in seconds from now) into an absolute ISO string.
 * This keeps demo-data.json time-agnostic — it always looks fresh regardless of when it's opened.
 */
function resolveTimestamps(obj: unknown): unknown {
  if (obj === null || typeof obj !== 'object') return obj;
  if (Array.isArray(obj)) return obj.map(resolveTimestamps);

  const now = Date.now();
  const result: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(obj as Record<string, unknown>)) {
    if (TIMESTAMP_FIELDS.has(key) && typeof value === 'number') {
      result[key] = new Date(now + value * 1000).toISOString();
    } else {
      result[key] = resolveTimestamps(value);
    }
  }
  return result;
}

/**
 * Fetches demo-data.json from the public folder, resolves relative timestamps,
 * then populates module-level state. Called once from main.tsx before React mounts.
 */
export async function initializeMockClient(): Promise<void> {
  const base = import.meta.env.BASE_URL ?? '/';
  const url = base.endsWith('/') ? `${base}demo-data.json` : `${base}/demo-data.json`;
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Failed to load mock data: ${response.status} ${response.statusText}`);
  }
  const raw = await response.json();
  const data = resolveTimestamps(raw) as typeof raw;
  state = {
    nodes: data.nodes ?? [],
    measurements: data.measurements ?? {},
    apiKeys: data.api_keys ?? [],
    nextApiKeyId: 100,
  };
}

// ── Helpers ───────────────────────────────────────────────────────────────────

function parseIntervalMs(interval?: string): number {
  if (!interval) return 6 * 3600_000; // default 6 h
  const m = interval.match(/^(\d+)([mhd])$/);
  if (!m) return 6 * 3600_000;
  const n = parseInt(m[1], 10);
  switch (m[2]) {
    case 'm': return n * 60_000;
    case 'h': return n * 3600_000;
    case 'd': return n * 86_400_000;
    default:  return 6 * 3600_000;
  }
}

function bucketMs(tsMs: number, intervalMs: number): number {
  return Math.floor(tsMs / intervalMs) * intervalMs;
}

function bytesPerSecToMbps(bps: number | undefined): number {
  return ((bps ?? 0) * 8) / 1_000_000;
}

/** Running weighted average helper. */
function runAvg(current: number, count: number, next: number): number {
  return (current * count + next) / (count + 1);
}

function nodeById(id: string): Node | undefined {
  return state.nodes.find(n => n.id === id);
}

function uuid(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = (Math.random() * 16) | 0;
    return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16);
  });
}

// ── Aggregation ───────────────────────────────────────────────────────────────

function computeAggregated(params: MeasurementQueryParams): AggregatedDataResponse {
  const fromMs = params.from ? new Date(params.from).getTime() : 0;
  const toMs   = params.to   ? new Date(params.to).getTime()   : Date.now();
  const intervalMs = parseIntervalMs(params.interval);

  const requestedIds = params.node_ids
    ? (Array.isArray(params.node_ids) ? params.node_ids : [params.node_ids])
    : Object.keys(state.measurements);

  // Filter out archived nodes when hide_archived is set
  const nodeIds = requestedIds.filter(id => {
    if (!params.hide_archived) return true;
    const node = nodeById(id);
    return node ? !node.archived : true;
  });

  const buckets = new Map<string, AggregatedMeasurement & { _count: number }>();

  for (const nodeId of nodeIds) {
    const node = nodeById(nodeId);
    if (!node) continue;

    const nodeMeasurements = (state.measurements[nodeId] ?? []).filter(m => {
      if (m.is_failed) return false;
      const ts = new Date(m.timestamp).getTime();
      return ts >= fromMs && ts <= toMs;
    });

    for (const m of nodeMeasurements) {
      const ts  = new Date(m.timestamp).getTime();
      const bkt = bucketMs(ts, intervalMs);
      const key = `${bkt}-${nodeId}`;

      if (!buckets.has(key)) {
        buckets.set(key, {
          timestamp:          new Date(bkt).toISOString(),
          node_id:            nodeId,
          node_name:          node.name,
          avg_download_mbps:  0,
          avg_upload_mbps:    0,
          avg_ping_ms:        0,
          avg_jitter_ms:      0,
          avg_packet_loss:    0,
          min_download_mbps:  Infinity,
          max_download_mbps:  -Infinity,
          sample_count:       0,
          _count:             0,
        });
      }

      const b  = buckets.get(key)!;
      const n  = b._count;
      const dl = bytesPerSecToMbps(m.download_bandwidth);
      const ul = bytesPerSecToMbps(m.upload_bandwidth);

      b.avg_download_mbps  = runAvg(b.avg_download_mbps,  n, dl);
      b.avg_upload_mbps    = runAvg(b.avg_upload_mbps,    n, ul);
      b.avg_ping_ms        = runAvg(b.avg_ping_ms,        n, m.ping_latency ?? 0);
      b.avg_jitter_ms      = runAvg(b.avg_jitter_ms,      n, m.ping_jitter  ?? 0);
      b.avg_packet_loss    = runAvg(b.avg_packet_loss,    n, m.packet_loss   ?? 0);
      b.min_download_mbps  = Math.min(b.min_download_mbps ?? Infinity, dl);
      b.max_download_mbps  = Math.max(b.max_download_mbps ?? -Infinity, dl);
      b._count             = n + 1;
      b.sample_count       = b._count;
    }
  }

  const data: AggregatedMeasurement[] = Array.from(buckets.values())
    .map(({ _count: _c, ...rest }) => {
      void _c;
      return rest as AggregatedMeasurement;
    })
    .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());

  return { data };
}

// ── Node details computation ──────────────────────────────────────────────────

function computeNodeDetails(nodeId: string): NodeDetails {
  const node = nodeById(nodeId);
  if (!node) throw new Error(`Node not found: ${nodeId}`);

  const all          = state.measurements[nodeId] ?? [];
  const successful   = all.filter(m => !m.is_failed);
  const failed       = all.filter(m => m.is_failed);
  const now          = Date.now();
  const ms24h        = 24 * 3600_000;

  const recent = successful.filter(
    m => now - new Date(m.timestamp).getTime() <= ms24h,
  );
  const failedRecent = failed.filter(
    m => now - new Date(m.timestamp).getTime() <= ms24h,
  );

  const avg = (arr: number[]) => arr.length ? arr.reduce((a, b) => a + b, 0) / arr.length : 0;

  const stats = {
    avg_download_mbps:  avg(successful.map(m => bytesPerSecToMbps(m.download_bandwidth))),
    avg_upload_mbps:    avg(successful.map(m => bytesPerSecToMbps(m.upload_bandwidth))),
    avg_ping_ms:        avg(successful.map(m => m.ping_latency  ?? 0)),
    avg_jitter_ms:      avg(successful.map(m => m.ping_jitter   ?? 0)),
    avg_packet_loss:    avg(successful.map(m => m.packet_loss   ?? 0)),
    success_count_24h:  recent.length,
    failed_count_24h:   failedRecent.length,
    success_rate_24h:
      recent.length + failedRecent.length > 0
        ? (recent.length / (recent.length + failedRecent.length)) * 100
        : 100,
  };

  const latest = successful.length > 0 ? successful[successful.length - 1] : undefined;
  const latestMeasurement = latest
    ? {
        timestamp:     latest.timestamp,
        download_mbps: bytesPerSecToMbps(latest.download_bandwidth),
        upload_mbps:   bytesPerSecToMbps(latest.upload_bandwidth),
        ping_ms:       latest.ping_latency ?? 0,
      }
    : undefined;

  return {
    ...node,
    total_measurements: all.length,
    failed_test_count:  failed.length,
    latest_measurement: latestMeasurement,
    statistics:         stats,
  };
}

// ── Dashboard summary computation ─────────────────────────────────────────────

function computeDashboard(): DashboardSummary {
  const visibleNodes = state.nodes;
  const now          = Date.now();
  const ms24h        = 24 * 3600_000;

  const allMeasurements = Object.values(state.measurements).flat();
  const recent24h = allMeasurements.filter(
    m => !m.is_failed && now - new Date(m.timestamp).getTime() <= ms24h,
  );

  const avg = (arr: number[]) => arr.length ? arr.reduce((a, b) => a + b, 0) / arr.length : 0;

  const sortedAll = allMeasurements
    .filter(m => !m.is_failed)
    .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());

  const average_stats_24h =
    recent24h.length > 0
      ? {
          download_mbps: avg(recent24h.map(m => bytesPerSecToMbps(m.download_bandwidth))),
          upload_mbps:   avg(recent24h.map(m => bytesPerSecToMbps(m.upload_bandwidth))),
          ping_ms:       avg(recent24h.map(m => m.ping_latency ?? 0)),
          jitter_ms:     avg(recent24h.map(m => m.ping_jitter  ?? 0)),
          packet_loss:   avg(recent24h.map(m => m.packet_loss  ?? 0)),
        }
      : undefined;

  return {
    total_nodes:           visibleNodes.length,
    active_nodes:          visibleNodes.filter(n => n.status === 'active'      && !n.archived).length,
    unreachable_nodes:     visibleNodes.filter(n => n.status === 'unreachable' && !n.archived).length,
    inactive_nodes:        visibleNodes.filter(n => n.status === 'inactive'    && !n.archived).length,
    total_measurements:    allMeasurements.length,
    measurements_last_24h: recent24h.length,
    last_measurement:      sortedAll[0]?.timestamp,
    average_stats_24h,
  };
}

// ── The mock client object ────────────────────────────────────────────────────

export const mockApiClient: IApiClient = {
  // ── Auth ─────────────────────────────────────────────────────────────────

  async login(_credentials: LoginRequest): Promise<LoginResponse> {
    return { token: 'mock-token', username: 'demo' };
  },

  // ── Nodes ────────────────────────────────────────────────────────────────

  async getNodes(): Promise<NodesResponse> {
    return { nodes: [...state.nodes], total: state.nodes.length };
  },

  async getNodeDetails(nodeId: string): Promise<NodeDetails> {
    return computeNodeDetails(nodeId);
  },

  async archiveNode(nodeId: string, archived: boolean): Promise<void> {
    const node = nodeById(nodeId);
    if (node) node.archived = archived;
  },

  async setNodeFavorite(nodeId: string, favorite: boolean): Promise<void> {
    const node = nodeById(nodeId);
    if (node) node.favorite = favorite;
  },

  async deleteNode(nodeId: string): Promise<void> {
    state.nodes = state.nodes.filter(n => n.id !== nodeId);
    delete state.measurements[nodeId];
  },

  // ── Measurements ─────────────────────────────────────────────────────────

  async getMeasurements(params?: MeasurementQueryParams): Promise<MeasurementsResponse> {
    let all = Object.values(state.measurements).flat();
    if (params?.node_ids) {
      const ids = Array.isArray(params.node_ids) ? params.node_ids : [params.node_ids];
      all = all.filter(m => ids.includes(m.node_id));
    }
    const limit  = params?.limit  ?? 50;
    const offset = params?.offset ?? 0;
    const page   = params?.page   ?? 1;
    const skip   = params?.page !== undefined ? (page - 1) * limit : offset;
    const sliced = all.slice(skip, skip + limit);
    return { measurements: sliced, total: all.length, page, limit };
  },

  async getNodeMeasurements(
    nodeId: string,
    params?: MeasurementQueryParams,
  ): Promise<MeasurementsResponse> {
    const all    = state.measurements[nodeId] ?? [];
    const limit  = params?.limit  ?? 50;
    const offset = params?.offset ?? 0;
    const page   = params?.page   ?? 1;
    const skip   = params?.page !== undefined ? (page - 1) * limit : offset;

    let filtered = [...all];
    if (params?.status === 'failed')     filtered = filtered.filter(m =>  m.is_failed);
    if (params?.status === 'successful') filtered = filtered.filter(m => !m.is_failed);

    const sliced = filtered.slice(skip, skip + limit);
    return { measurements: sliced, total: filtered.length, page, limit };
  },

  async getAggregatedMeasurements(
    params: MeasurementQueryParams,
  ): Promise<AggregatedDataResponse> {
    return computeAggregated(params);
  },

  // ── API Keys ─────────────────────────────────────────────────────────────

  async getAPIKeys(): Promise<APIKeysResponse> {
    return { api_keys: [...state.apiKeys], total: state.apiKeys.length };
  },

  async createAPIKey(data: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse> {
    const id      = uuid();
    const rawKey  = `demo_${Math.random().toString(36).slice(2, 18)}`;
    const newKey: APIKey & { key: string } = {
      id,
      name:       data.name,
      key:        rawKey,
      enabled:    true,
      created_at: new Date().toISOString(),
    };
    state.apiKeys.push(newKey);
    return newKey;
  },

  async updateAPIKey(id: string, data: UpdateAPIKeyRequest): Promise<APIKey> {
    const key = state.apiKeys.find(k => k.id === id);
    if (!key) throw new Error(`API key not found: ${id}`);
    key.enabled = data.enabled;
    return { ...key };
  },

  async deleteAPIKey(id: string): Promise<void> {
    state.apiKeys = state.apiKeys.filter(k => k.id !== id);
  },

  // ── Dashboard ────────────────────────────────────────────────────────────

  async getDashboardSummary(): Promise<DashboardSummary> {
    return computeDashboard();
  },

  // ── Health ───────────────────────────────────────────────────────────────

  async getHealth() {
    return {
      status: 'healthy',
      database: 'connected',
      uptime_seconds: 0,
      version: 'demo',
    };
  },
};
