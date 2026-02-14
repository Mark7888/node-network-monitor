import { useState, useEffect } from 'react';
import { useDashboard } from '../hooks/useDashboard';
import { useMeasurements } from '@/modules/measurements/hooks/useMeasurements';
import { useChartData } from '@/modules/measurements/hooks/useChartData';
import { useNodes } from '@/modules/nodes/hooks/useNodes';
import { useAutoRefresh } from '@/shared/hooks/useAutoRefresh';
import { TIME_RANGES, TimeRange } from '@/shared/utils/constants';
import { formatMbps, formatLatency, formatNumber } from '@/shared/utils/format';
import Card from '@/shared/components/ui/Card';
import Spinner from '@/shared/components/ui/Spinner';
import ErrorMessage from '@/shared/components/ui/ErrorMessage';
import DownloadChart from '@/modules/measurements/components/charts/DownloadChart';
import UploadChart from '@/modules/measurements/components/charts/UploadChart';
import PingChart from '@/modules/measurements/components/charts/PingChart';
import JitterChart from '@/modules/measurements/components/charts/JitterChart';
import PacketLossChart from '@/modules/measurements/components/charts/PacketLossChart';
import env from '@/core/config/env';

/**
 * Dashboard page component
 */
export default function DashboardPage() {
  const [timeRange, setTimeRange] = useState<TimeRange>('day');
  const { summary, isLoading: summaryLoading, refetch: refetchSummary } = useDashboard();
  const { fetchNodes } = useNodes();
  const { data: measurements, isLoading: measurementsLoading, refetch: refetchMeasurements } = useMeasurements(timeRange);
  const chartData = useChartData(measurements);

  // Initial fetch
  useEffect(() => {
    fetchNodes();
  }, []);

  // Auto-refresh data
  useAutoRefresh(() => {
    refetchSummary();
    refetchMeasurements();
  }, env.refreshInterval);

  if (summaryLoading && !summary) {
    return <Spinner message="Loading dashboard..." />;
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-base-content/60 mt-1">
          Overview of all speedtest nodes and measurements
        </p>
      </div>

      {/* Summary Stats */}
      {summary && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="stats shadow bg-base-100">
            <div className="stat">
              <div className="stat-title">Total Nodes</div>
              <div className="stat-value text-primary">{summary.total_nodes}</div>
              <div className="stat-desc">{summary.active_nodes} active</div>
            </div>
          </div>

          <div className="stats shadow bg-base-100">
            <div className="stat">
              <div className="stat-title">Avg Download</div>
              <div className="stat-value text-success">
                {formatMbps(summary.average_stats_24h?.download_mbps)}
              </div>
              <div className="stat-desc">Last 24 hours average</div>
            </div>
          </div>

          <div className="stats shadow bg-base-100">
            <div className="stat">
              <div className="stat-title">Avg Upload</div>
              <div className="stat-value text-info">
                {formatMbps(summary.average_stats_24h?.upload_mbps)}
              </div>
              <div className="stat-desc">Last 24 hours average</div>
            </div>
          </div>

          <div className="stats shadow bg-base-100">
            <div className="stat">
              <div className="stat-title">Avg Ping</div>
              <div className="stat-value text-warning">
                {formatLatency(summary.average_stats_24h?.ping_ms)}
              </div>
              <div className="stat-desc">Average latency</div>
            </div>
          </div>
        </div>
      )}

      {/* Time Range Filter */}
      <div className="flex gap-2">
        {TIME_RANGES.map((range) => (
          <button
            key={range.value}
            onClick={() => setTimeRange(range.value)}
            className={`btn ${
              timeRange === range.value ? 'btn-primary' : 'btn-ghost'
            }`}
          >
            {range.label}
          </button>
        ))}
      </div>

      {/* Charts */}
      {measurementsLoading && measurements.length === 0 ? (
        <Spinner message="Loading charts..." />
      ) : measurements.length === 0 ? (
        <Card>
          <div className="text-center py-8 text-base-content/60">
            No measurement data available for the selected time range
          </div>
        </Card>
      ) : (
        <div className="space-y-6">
          <Card>
            <DownloadChart data={chartData.downloadData} />
          </Card>

          <Card>
            <UploadChart data={chartData.uploadData} />
          </Card>

          <Card>
            <PingChart data={chartData.pingData} />
          </Card>

          <Card>
            <JitterChart data={chartData.jitterData} />
          </Card>

          <Card>
            <PacketLossChart data={chartData.packetLossData} />
          </Card>
        </div>
      )}
    </div>
  );
}
