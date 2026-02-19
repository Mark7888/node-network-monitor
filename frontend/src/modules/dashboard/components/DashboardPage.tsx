import { useState, useEffect } from 'react';
import { useDashboard } from '../hooks/useDashboard';
import { useMeasurements } from '@/modules/measurements/hooks/useMeasurements';
import { useChartData } from '@/modules/measurements/hooks/useChartData';
import { useNodes } from '@/modules/nodes/hooks/useNodes';
import { useAutoRefresh } from '@/shared/hooks/useAutoRefresh';
import { TIME_RANGES, TimeRange } from '@/shared/utils/constants';
import { formatMbps, formatLatency } from '@/shared/utils/format';
import Card from '@/shared/components/ui/Card';
import Spinner from '@/shared/components/ui/Spinner';
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
  const { data: measurements, isLoading: measurementsLoading, refetch: refetchMeasurements } = useMeasurements(timeRange, undefined, true);
  const chartData = useChartData(measurements);

  // Initial fetch
  useEffect(() => {
    fetchNodes();
  }, [fetchNodes]);

  // Auto-refresh data
  useAutoRefresh(() => {
    refetchSummary();
    refetchMeasurements();
  }, env.refreshInterval);

  if (summaryLoading && !summary) {
    return <Spinner message="Loading dashboard..." />;
  }

  return (
    <div className="space-y-4 md:space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl md:text-3xl font-bold">Dashboard</h1>
        <p className="text-sm md:text-base text-base-content/60 mt-1">
          Overview of all speedtest nodes and measurements
        </p>
      </div>

      {/* Summary Stats */}
      {summary && (
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
          <div className="stats shadow bg-base-100">
            <div className="stat p-3 md:p-4">
              <div className="stat-title text-xs md:text-sm">Total Nodes</div>
              <div className="stat-value text-2xl md:text-3xl text-primary">{summary.total_nodes}</div>
              <div className="stat-desc text-xs hidden sm:block">
                {summary.active_nodes} active, {summary.unreachable_nodes} unreachable, {summary.inactive_nodes} inactive
              </div>
            </div>
          </div>

          <div className="stats shadow bg-base-100">
            <div className="stat p-3 md:p-4">
              <div className="stat-title text-xs md:text-sm">Avg Download</div>
              <div className="stat-value text-2xl md:text-3xl text-success">
                {formatMbps(summary.average_stats_24h?.download_mbps)}
              </div>
              <div className="stat-desc text-xs hidden sm:block">Last 24 hours average</div>
            </div>
          </div>

          <div className="stats shadow bg-base-100">
            <div className="stat p-3 md:p-4">
              <div className="stat-title text-xs md:text-sm">Avg Upload</div>
              <div className="stat-value text-2xl md:text-3xl text-info">
                {formatMbps(summary.average_stats_24h?.upload_mbps)}
              </div>
              <div className="stat-desc text-xs hidden sm:block">Last 24 hours average</div>
            </div>
          </div>

          <div className="stats shadow bg-base-100">
            <div className="stat p-3 md:p-4">
              <div className="stat-title text-xs md:text-sm">Avg Ping</div>
              <div className="stat-value text-2xl md:text-3xl text-warning">
                {formatLatency(summary.average_stats_24h?.ping_ms)}
              </div>
              <div className="stat-desc text-xs hidden sm:block">Average latency</div>
            </div>
          </div>
        </div>
      )}

      {/* Time Range Filter */}
      <div className="flex gap-2 overflow-x-auto pb-2">
        {TIME_RANGES.map((range) => (
          <button
            key={range.value}
            onClick={() => setTimeRange(range.value)}
            className={`btn btn-sm md:btn-md whitespace-nowrap ${
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
