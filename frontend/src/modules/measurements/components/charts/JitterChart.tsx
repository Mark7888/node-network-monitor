import BaseChart from './BaseChart';
import { ChartSeries } from '../../utils/dataTransform';
import { generateChartOption, getNodeColor } from '../../utils/chartConfig';

interface JitterChartProps {
  data: ChartSeries[];
  height?: number | string;
}

/**
 * Jitter chart component
 */
export default function JitterChart({ data, height }: JitterChartProps) {
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
    'Jitter',
    'ms',
    series,
    (value: number) => `${value} ms`
  );

  return <BaseChart option={option} height={height} />;
}
