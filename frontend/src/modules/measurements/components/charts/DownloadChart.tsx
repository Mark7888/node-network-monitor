import BaseChart from './BaseChart';
import { ChartSeries } from '../../utils/dataTransform';
import { generateChartOption, getNodeColor } from '../../utils/chartConfig';

interface DownloadChartProps {
  data: ChartSeries[];
  height?: number | string;
  failedTimestamps?: string[];
}

/**
 * Download speed chart component
 */
export default function DownloadChart({ data, height, failedTimestamps }: DownloadChartProps) {
  const series = data.map((nodeSeries, index) => ({
    name: nodeSeries.node_name,
    type: 'line' as const,
    smooth: true,
    data: nodeSeries.data.map((point) => [
      point.timestamp,
      point.value != null && !isNaN(point.value) ? Number(point.value.toFixed(2)) : 0
    ]),
    itemStyle: {
      color: getNodeColor(index),
    },
  }));

  const option = generateChartOption(
    'Download Speed',
    'Mbps',
    series,
    (value: number) => `${value} Mbps`,
    failedTimestamps
  );

  return <BaseChart option={option} height={height} />;
}
