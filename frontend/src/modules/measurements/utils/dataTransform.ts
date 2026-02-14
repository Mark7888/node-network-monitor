import { AggregatedMeasurement } from '../types/measurement.types';
import { formatChartLabel } from '@/shared/utils/date';
import { TimeRange } from '@/shared/utils/constants';

/**
 * Transform aggregated data for chart display
 */

export interface ChartDataPoint {
  timestamp: string;
  value: number;
}

export interface ChartSeries {
  node_id: string;
  node_name: string;
  data: ChartDataPoint[];
}

/**
 * Group aggregated measurements by node
 */
export function groupByNode(measurements: AggregatedMeasurement[]): Map<string, AggregatedMeasurement[]> {
  const grouped = new Map<string, AggregatedMeasurement[]>();
  
  measurements.forEach((m) => {
    if (!grouped.has(m.node_id)) {
      grouped.set(m.node_id, []);
    }
    grouped.get(m.node_id)!.push(m);
  });
  
  return grouped;
}

/**
 * Transform download speed data for charts
 */
export function transformDownloadData(measurements: AggregatedMeasurement[]): ChartSeries[] {
  const grouped = groupByNode(measurements);
  const series: ChartSeries[] = [];
  
  grouped.forEach((data, nodeId) => {
    series.push({
      node_id: nodeId,
      node_name: data[0]?.node_name || nodeId,
      data: data.map((m) => ({
        timestamp: m.timestamp,
        value: m.avg_download_mbps, // Already in Mbps
      })),
    });
  });
  
  return series;
}

/**
 * Transform upload speed data for charts
 */
export function transformUploadData(measurements: AggregatedMeasurement[]): ChartSeries[] {
  const grouped = groupByNode(measurements);
  const series: ChartSeries[] = [];
  
  grouped.forEach((data, nodeId) => {
    series.push({
      node_id: nodeId,
      node_name: data[0]?.node_name || nodeId,
      data: data.map((m) => ({
        timestamp: m.timestamp,
        value: m.avg_upload_mbps, // Already in Mbps
      })),
    });
  });
  
  return series;
}

/**
 * Transform ping data for charts
 */
export function transformPingData(measurements: AggregatedMeasurement[]): ChartSeries[] {
  const grouped = groupByNode(measurements);
  const series: ChartSeries[] = [];
  
  grouped.forEach((data, nodeId) => {
    series.push({
      node_id: nodeId,
      node_name: data[0]?.node_name || nodeId,
      data: data.map((m) => ({
        timestamp: m.timestamp,
        value: m.avg_ping_ms,
      })),
    });
  });
  
  return series;
}

/**
 * Transform jitter data for charts
 */
export function transformJitterData(measurements: AggregatedMeasurement[]): ChartSeries[] {
  const grouped = groupByNode(measurements);
  const series: ChartSeries[] = [];
  
  grouped.forEach((data, nodeId) => {
    series.push({
      node_id: nodeId,
      node_name: data[0]?.node_name || nodeId,
      data: data.map((m) => ({
        timestamp: m.timestamp,
        value: m.avg_jitter_ms,
      })),
    });
  });
  
  return series;
}

/**
 * Transform packet loss data for charts
 */
export function transformPacketLossData(measurements: AggregatedMeasurement[]): ChartSeries[] {
  const grouped = groupByNode(measurements);
  const series: ChartSeries[] = [];
  
  grouped.forEach((data, nodeId) => {
    series.push({
      node_id: nodeId,
      node_name: data[0]?.node_name || nodeId,
      data: data.map((m) => ({
        timestamp: m.timestamp,
        value: m.avg_packet_loss,
      })),
    });
  });
  
  return series;
}
