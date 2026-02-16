import BaseChart from './BaseChart';
import { ChartSeries } from '../../utils/dataTransform';
import { generateChartOption, getNodeColor } from '../../utils/chartConfig';

interface PacketLossChartProps {
  data: ChartSeries[];
  height?: number | string;
  failedTimestamps?: string[];
}

/**
 * Packet loss chart component
 */
export default function PacketLossChart({ data, height, failedTimestamps }: PacketLossChartProps) {
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
    'Packet Loss',
    '%',
    series,
    (value: number) => `${value}%`,
    failedTimestamps
  );

  return <BaseChart option={option} height={height} />;
}
