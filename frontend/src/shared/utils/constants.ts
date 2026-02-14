/**
 * Application-wide constants
 */

// Node status types
export const NODE_STATUS = {
  ACTIVE: 'active',
  UNREACHABLE: 'unreachable',
  INACTIVE: 'inactive',
} as const;

export type NodeStatus = typeof NODE_STATUS[keyof typeof NODE_STATUS];

// Time range options
export const TIME_RANGES = [
  { label: 'Last Day', value: 'day' as const, hours: 24 },
  { label: 'Last Week', value: 'week' as const, hours: 168 },
  { label: 'Last Month', value: 'month' as const, hours: 720 },
] as const;

export type TimeRange = typeof TIME_RANGES[number]['value'];

// Chart colors for nodes (cycle through these)
export const CHART_COLORS = [
  '#3b82f6', // blue
  '#8b5cf6', // purple
  '#10b981', // green
  '#f59e0b', // amber
  '#ef4444', // red
  '#06b6d4', // cyan
  '#ec4899', // pink
  '#84cc16', // lime
] as const;

// Status badge classes (daisyUI)
export const STATUS_BADGE_CLASS: Record<NodeStatus, string> = {
  [NODE_STATUS.ACTIVE]: 'badge-success',
  [NODE_STATUS.UNREACHABLE]: 'badge-warning',
  [NODE_STATUS.INACTIVE]: 'badge-ghost',
} as const;

// Local storage keys
export const STORAGE_KEYS = {
  TOKEN: 'token',
  USERNAME: 'username',
  THEME: 'theme',
} as const;
