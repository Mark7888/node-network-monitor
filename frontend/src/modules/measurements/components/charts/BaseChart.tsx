import ReactECharts from 'echarts-for-react';
import { EChartsOption } from 'echarts';

interface BaseChartProps {
  option: EChartsOption;
  height?: number | string;
  className?: string;
}

/**
 * Base chart component wrapper for ECharts
 */
export default function BaseChart({ 
  option, 
  height = 400,
  className = '' 
}: BaseChartProps) {
  return (
    <div className={`chart-container ${className}`}>
      <ReactECharts
        option={option}
        style={{ height, width: '100%' }}
        opts={{ renderer: 'canvas' }}
      />
    </div>
  );
}
