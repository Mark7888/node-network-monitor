import BaseChart from './BaseChart';
import { ChartSeries } from '../../utils/dataTransform';
import { generateChartOption, getNodeColor } from '../../utils/chartConfig';

interface PingChartProps {
  data: ChartSeries[];
  height?: number | string;
}

/**
 * Ping/latency chart component
 */
export default function PingChart({ data, height }: PingChartProps) {
  const series = data.map((nodeSeries, index) => ({
    name: nodeSeries.node_name,
    type: 'line',
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
    'Ping / Latency',
    'ms',
    series,
    (value: number) => `${value} ms`
  );

  return <BaseChart option={option} height={height} />;
}
