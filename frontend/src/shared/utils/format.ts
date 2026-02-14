/**
 * Formatting utilities for displaying data
 */

/**
 * Convert bytes/second to Mbps
 */
export function bytesToMbps(bytes: number): number {
  return (bytes / 1000000) * 8;
}

/**
 * Format bytes/second to Mbps string
 */
export function formatSpeed(megaBitsPerSecond: number | null | undefined): string {
  if (megaBitsPerSecond == null || isNaN(megaBitsPerSecond)) {
    return 'N/A';
  }
  return `${megaBitsPerSecond.toFixed(2)} Mbps`;
}

/**
 * Format Mbps value directly (when already in Mbps)
 */
export function formatMbps(mbps: number | null | undefined): string {
  if (mbps == null || isNaN(mbps)) {
    return 'N/A';
  }
  return `${mbps.toFixed(2)} Mbps`;
}

/**
 * Format milliseconds to display string
 */
export function formatLatency(ms: number | null | undefined): string {
  if (ms == null || isNaN(ms)) {
    return 'N/A';
  }
  return `${ms.toFixed(2)} ms`;
}

/**
 * Format percentage
 */
export function formatPercent(value: number | null | undefined): string {
  if (value == null || isNaN(value)) {
    return 'N/A';
  }
  return `${value.toFixed(2)}%`;
}

/**
 * Format number with commas
 */
export function formatNumber(num: number | null | undefined): string {
  if (num == null || isNaN(num)) {
    return 'N/A';
  }
  return num.toLocaleString();
}
