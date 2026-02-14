import { format, formatDistanceToNow, isToday, isYesterday, parseISO } from 'date-fns';

/**
 * Date/time formatting utilities
 */

/**
 * Format timestamp for display
 */
export function formatTimestamp(timestamp: string, formatStr: string = 'MMM dd, HH:mm'): string {
  return format(parseISO(timestamp), formatStr);
}

/**
 * Format relative time (e.g., "2 minutes ago")
 */
export function formatRelativeTime(timestamp: string): string {
  return formatDistanceToNow(parseISO(timestamp), { addSuffix: true });
}

/**
 * Format date for chart axis label based on time range
 */
export function formatChartLabel(timestamp: string, range: 'day' | 'week' | 'month'): string {
  const date = parseISO(timestamp);
  
  if (range === 'day') {
    return format(date, 'HH:mm');
  } else if (range === 'week') {
    if (isToday(date)) return format(date, 'HH:mm');
    if (isYesterday(date)) return `Yesterday ${format(date, 'HH:mm')}`;
    return format(date, 'MMM dd HH:mm');
  } else {
    return format(date, 'MMM dd');
  }
}

/**
 * Get time range boundaries
 */
export function getTimeRange(range: 'day' | 'week' | 'month'): { from: string; to: string } {
  const to = new Date();
  const from = new Date();
  
  switch (range) {
    case 'day':
      from.setHours(from.getHours() - 24);
      break;
    case 'week':
      from.setDate(from.getDate() - 7);
      break;
    case 'month':
      from.setMonth(from.getMonth() - 1);
      break;
  }
  
  return {
    from: from.toISOString(),
    to: to.toISOString(),
  };
}

/**
 * Calculate appropriate interval for aggregation based on time range
 */
export function getAggregationInterval(range: 'day' | 'week' | 'month'): string {
  switch (range) {
    case 'day':
      return '5m';  // 5 minutes for better granularity
    case 'week':
      return '1h';  // 1 hour
    case 'month':
      return '6h';  // 6 hours
    default:
      return '5m';
  }
}
