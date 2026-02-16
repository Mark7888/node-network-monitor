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
  
  // Create markLine data for failed measurements
  const markLineData = failedTimestamps?.map(timestamp => ({
    xAxis: timestamp,
    lineStyle: {
      color: 'rgba(239, 68, 68, 0.35)', // Tailwind red-500 with lower opacity
      type: 'solid' as const,
      width: 2,
    },
    label: {
      show: false, // Hide labels to avoid clutter
    },
  })) || [];

  // Add markLine to each series if there are failed timestamps
  const enhancedSeries = Array.isArray(series) 
    ? series.map((s: any) => ({
        ...s,
        markLine: failedTimestamps && failedTimestamps.length > 0 ? {
          silent: true,
          symbol: 'none', // No symbols at the ends
          data: markLineData,
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
