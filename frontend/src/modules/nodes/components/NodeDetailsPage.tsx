import { useState, useMemo } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useNodeDetails } from '../hooks/useNodeDetails';
import { useMeasurements } from '@/modules/measurements/hooks/useMeasurements';
import { useChartData } from '@/modules/measurements/hooks/useChartData';
import { useFailedMeasurements } from '@/modules/measurements/hooks/useFailedMeasurements';
import { useAutoRefresh } from '@/shared/hooks/useAutoRefresh';
import { TIME_RANGES, TimeRange } from '@/shared/utils/constants';
import { formatSpeed, formatLatency, formatPercent } from '@/shared/utils/format';
import { formatTimestamp } from '@/shared/utils/date';
import Card from '@/shared/components/ui/Card';
import Badge from '@/shared/components/ui/Badge';
import Spinner from '@/shared/components/ui/Spinner';
import ErrorMessage from '@/shared/components/ui/ErrorMessage';
import DownloadChart from '@/modules/measurements/components/charts/DownloadChart';
import UploadChart from '@/modules/measurements/components/charts/UploadChart';
import PingChart from '@/modules/measurements/components/charts/PingChart';
import JitterChart from '@/modules/measurements/components/charts/JitterChart';
import PacketLossChart from '@/modules/measurements/components/charts/PacketLossChart';
import MeasurementsList from '@/modules/measurements/components/MeasurementsList';
import { ArrowLeft } from 'lucide-react';
import env from '@/core/config/env';

type TabType = 'charts' | 'measurements';

/**
 * Node details page component
 */
export default function NodeDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const [timeRange, setTimeRange] = useState<TimeRange>('day');
  const [activeTab, setActiveTab] = useState<TabType>('charts');
  
  // Memoize nodeIds array to prevent infinite re-renders
  const nodeIds = useMemo(() => id ? [id] : undefined, [id]);
  
  const { node, isLoading: nodeLoading, error: nodeError, refetch: refetchNode } = useNodeDetails(id!);
  const { data: measurements, isLoading: measurementsLoading, refetch: refetchMeasurements } = useMeasurements(timeRange, nodeIds);
  const { failedTimestamps, refetch: refetchFailedMeasurements } = useFailedMeasurements(timeRange, id);
  const chartData = useChartData(measurements);

  // Auto-refresh
  useAutoRefresh(() => {
    refetchNode();
    if (activeTab === 'charts') {
      refetchMeasurements();
      refetchFailedMeasurements();
    }
  }, env.refreshInterval);

  if (nodeLoading && !node) {
    return <Spinner message="Loading node details..." />;
  }

  if (nodeError || !node) {
    return (
      <div className="space-y-4">
        <Link to="/nodes" className="btn btn-ghost btn-sm gap-2">
          <ArrowLeft size={16} />
          Back to Nodes
        </Link>
        <ErrorMessage message={nodeError || 'Node not found'} onRetry={refetchNode} />
      </div>
    );
  }

  const badgeVariant = node.status === 'active' ? 'success' : 
                       node.status === 'unreachable' ? 'warning' : 'ghost';

  console.log(node);

  return (
    <div className="space-y-6">
      {/* Back Button */}
      <Link to="/nodes" className="btn btn-ghost btn-sm gap-2">
        <ArrowLeft size={16} />
        Back to Nodes
      </Link>

      {/* Node Header */}
      <Card>
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">{node.name}</h1>
            <p className="text-base-content/60 mt-1">ID: {node.id}</p>
            <div className="mt-3 space-y-1 text-sm">
              <p>First seen: {formatTimestamp(node.first_seen, 'MMM dd, yyyy HH:mm')}</p>
              <p>Last seen: {formatTimestamp(node.last_seen, 'MMM dd, yyyy HH:mm')}</p>
              {node.location && <p>Location: {node.location}</p>}
            </div>
          </div>
          <Badge variant={badgeVariant} size="lg">
            {node.status.charAt(0).toUpperCase() + node.status.slice(1)}
          </Badge>
        </div>
      </Card>

      {/* Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
        <div className="stats shadow bg-base-100">
          <div className="stat">
            <div className="stat-title">Avg Download</div>
            <div className="stat-value text-success text-2xl">
              {formatSpeed(node.statistics.avg_download_mbps)}
            </div>
          </div>
        </div>

        <div className="stats shadow bg-base-100">
          <div className="stat">
            <div className="stat-title">Avg Upload</div>
            <div className="stat-value text-info text-2xl">
              {formatSpeed(node.statistics.avg_upload_mbps)}
            </div>
          </div>
        </div>

        <div className="stats shadow bg-base-100">
          <div className="stat">
            <div className="stat-title">Avg Ping</div>
            <div className="stat-value text-warning text-2xl">
              {formatLatency(node.statistics.avg_ping_ms)}
            </div>
          </div>
        </div>

        <div className="stats shadow bg-base-100">
          <div className="stat">
            <div className="stat-title">Avg Packet Loss</div>
            <div className={`stat-value text-2xl ${node.statistics.avg_packet_loss < 5 ? 'text-info' : 'text-error'}`}>
              {formatPercent(node.statistics.avg_packet_loss)}
            </div>
          </div>
        </div>

        <div className="stats shadow bg-base-100">
          <div className="stat">
            <div className="stat-title">Success Rate (24h)</div>
            <div className={`stat-value text-2xl ${
              node.statistics.success_rate_24h >= 80 ? 'text-success' : 
              node.statistics.success_rate_24h >= 50 ? 'text-warning' : 
              'text-error'
            }`}>
              {node.statistics.success_rate_24h.toFixed(1)}%
            </div>
            <div className="stat-desc">
              {node.statistics.success_count_24h} succeeded, {node.statistics.failed_count_24h} failed
            </div>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="tabs tabs-boxed bg-base-200 w-fit">
        <button
          className={`tab ${activeTab === 'charts' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('charts')}
        >
          Charts
        </button>
        <button
          className={`tab ${activeTab === 'measurements' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('measurements')}
        >
          Measurements
        </button>
      </div>

      {/* Charts Tab */}
      {activeTab === 'charts' && (
        <>
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
                <DownloadChart data={chartData.downloadData} failedTimestamps={failedTimestamps} />
              </Card>

              <Card>
                <UploadChart data={chartData.uploadData} failedTimestamps={failedTimestamps} />
              </Card>

              <Card>
                <PingChart data={chartData.pingData} failedTimestamps={failedTimestamps} />
              </Card>

              <Card>
                <JitterChart data={chartData.jitterData} failedTimestamps={failedTimestamps} />
              </Card>

              <Card>
                <PacketLossChart data={chartData.packetLossData} failedTimestamps={failedTimestamps} />
              </Card>
            </div>
          )}
        </>
      )}

      {/* Measurements List Tab */}
      {activeTab === 'measurements' && id && (
        <MeasurementsList nodeId={id} />
      )}
    </div>
  );
}
