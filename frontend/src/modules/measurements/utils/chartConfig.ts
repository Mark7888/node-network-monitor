import { EChartsOption } from 'echarts';
import { CHART_COLORS } from '@/shared/utils/constants';
import env from '@/core/config/env';

/**
 * Base chart configuration for ECharts
 */
export function getBaseChartConfig(): Partial<EChartsOption> {
  return {
    animation: env.enableChartAnimation,
    grid: {
      left: '3%',
      right: '4%',
      bottom: '10%',
      containLabel: true,
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985',
        },
      },
    },
    legend: {
      bottom: 0,
      type: 'scroll',
    },
    xAxis: {
      type: 'time',
      axisLabel: {
        hideOverlap: true,
      },
    },
  };
}

/**
 * Get color for a node (cycles through CHART_COLORS)
 */
export function getNodeColor(index: number): string {
  return CHART_COLORS[index % CHART_COLORS.length];
}

/**
 * Group consecutive failed timestamps into ranges
 */
function groupConsecutiveFailures(timestamps: string[]): Array<[string, string]> {
  if (!timestamps || timestamps.length === 0) return [];
  
  // Sort timestamps
  const sorted = [...timestamps].sort();
  const groups: Array<[string, string]> = [];
  
  let rangeStart = sorted[0];
  let rangeEnd = sorted[0];
  
  for (let i = 1; i < sorted.length; i++) {
    const current = new Date(sorted[i]).getTime();
    const previous = new Date(sorted[i - 1]).getTime();
    
    // If timestamps are within 15 minutes (900000ms), consider them consecutive
    if (current - previous <= 900000) {
      rangeEnd = sorted[i];
    } else {
      // Save the current range and start a new one
      groups.push([rangeStart, rangeEnd]);
      rangeStart = sorted[i];
      rangeEnd = sorted[i];
    }
  }
  
  // Add the last range
  groups.push([rangeStart, rangeEnd]);
  
  return groups;
}

/**
 * Generate chart option for single metric
 */
export function generateChartOption(
  title: string,
  yAxisLabel: string,
  series: EChartsOption['series'],
  formatter?: (value: number) => string,
  failedTimestamps?: string[]
): EChartsOption {
  const baseConfig = getBaseChartConfig();
  
  // Group consecutive failures into ranges
  const failureRanges = failedTimestamps ? groupConsecutiveFailures(failedTimestamps) : [];
  
  // Create markArea data for failed measurement ranges
  const markAreaData = failureRanges.map(([start, end]) => {
    const startTime = new Date(start).getTime();
    const endTime = new Date(end).getTime();
    
    // If it's a single point (or very close), add some width for visibility (5 minutes on each side)
    const padding = startTime === endTime ? 300000 : 0; // 5 minutes in ms
    
    return [
      {
        xAxis: new Date(startTime - padding).toISOString(),
        itemStyle: {
          color: 'rgba(239, 68, 68, 0.15)', // Tailwind red-500 with low opacity for background
        },
      },
      {
        xAxis: new Date(endTime + padding).toISOString(),
      }
    ];
  });

  // Add markArea to each series if there are failed timestamps
  const enhancedSeries = Array.isArray(series) 
    ? series.map((s: any) => ({
        ...s,
        markArea: failedTimestamps && failedTimestamps.length > 0 ? {
          silent: true,
          data: markAreaData,
          animation: false, // Disable animation for better performance
        } : undefined,
      }))
    : series;
  
  return {
    ...baseConfig,
    title: {
      text: title,
      textStyle: {
        fontSize: 16,
        fontWeight: 600,
      },
    },
    yAxis: {
      type: 'value',
      name: yAxisLabel,
      axisLabel: {
        formatter: formatter || undefined,
      },
    },
    series: enhancedSeries,
  } as EChartsOption;
}
