import { useState } from 'react';
import { useNodeMeasurements } from '../hooks/useNodeMeasurements';
import { formatTimestamp } from '@/shared/utils/date';
import { formatMbps, formatLatency, formatPercent } from '@/shared/utils/format';
import Card from '@/shared/components/ui/Card';
import Spinner from '@/shared/components/ui/Spinner';
import ErrorMessage from '@/shared/components/ui/ErrorMessage';
import Badge from '@/shared/components/ui/Badge';
import { ChevronLeft, ChevronRight, ExternalLink, AlertCircle } from 'lucide-react';

interface MeasurementsListProps {
  nodeId: string;
}

type StatusFilter = 'all' | 'successful' | 'failed';

/**
 * Component to display detailed measurements list with pagination
 */
export default function MeasurementsList({ nodeId }: MeasurementsListProps) {
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('all');

  const {
    measurements,
    isLoading,
    error,
    total,
    page,
    limit,
    totalPages,
    nextPage,
    prevPage,
    goToPage,
    refetch,
  } = useNodeMeasurements(nodeId, statusFilter);

  if (isLoading && measurements.length === 0) {
    return <Spinner message="Loading measurements..." />;
  }

  if (error) {
    return <ErrorMessage message={error} onRetry={refetch} />;
  }

  if (measurements.length === 0) {
    return (
      <Card>
        <div className="text-center py-8 text-base-content/60">
          No measurements available
        </div>
      </Card>
    );
  }

  // Helper to convert bandwidth bytes/sec to Mbps
  const formatBandwidth = (bandwidth?: number) => {
    if (!bandwidth) return 'N/A';
    return formatMbps((bandwidth / 1000000) * 8);
  };

  return (
    <div className="space-y-4">
      {/* Filter and Summary */}
      <div className="flex justify-between items-center">
        <div className="text-sm text-base-content/70">
          Showing {(page - 1) * limit + 1}-{Math.min(page * limit, total)} of {total} measurements
        </div>
        
        {/* Status Filter Dropdown */}
        <select 
          className="select select-bordered select-sm"
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value as StatusFilter)}
        >
          <option value="all">All Measurements</option>
          <option value="successful">Successful Only</option>
          <option value="failed">Failed Only</option>
        </select>
      </div>

      {/* Measurements Cards */}
      <div className="space-y-4">
        {measurements.map((m) => (
          <Card key={m.id}>
            <div className="space-y-4">
              {/* Header */}
              <div className="flex justify-between items-start">
                <div>
                  <div className="flex items-center gap-2">
                    <h3 className="text-lg font-semibold">
                      {formatTimestamp(m.timestamp, 'MMM dd, yyyy HH:mm:ss')}
                    </h3>
                    {m.is_failed && (
                      <Badge variant="error" size="sm">
                        Failed
                      </Badge>
                    )}
                  </div>
                  <p className="text-sm text-base-content/60">ID: {m.id}</p>
                </div>
                {m.result_url && !m.is_failed && (
                  <a
                    href={m.result_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="btn btn-ghost btn-sm gap-2"
                  >
                    View Result <ExternalLink size={14} />
                  </a>
                )}
              </div>

              {/* Failed Measurement Display */}
              {m.is_failed ? (
                <div className="alert alert-error">
                  <AlertCircle size={20} />
                  <div>
                    <h4 className="font-semibold">Measurement Failed</h4>
                    <p className="text-sm">{m.error_message || 'No error details available'}</p>
                  </div>
                </div>
              ) : (
                <>
                  {/* Main Metrics Grid */}
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div>
                      <div className="text-xs text-base-content/60">Download</div>
                      <div className="text-lg font-semibold text-success">
                        {formatBandwidth(m.download_bandwidth)}
                      </div>
                    </div>
                    <div>
                      <div className="text-xs text-base-content/60">Upload</div>
                      <div className="text-lg font-semibold text-info">
                        {formatBandwidth(m.upload_bandwidth)}
                      </div>
                    </div>
                    <div>
                      <div className="text-xs text-base-content/60">Ping Latency</div>
                      <div className="text-lg font-semibold text-warning">
                        {formatLatency(m.ping_latency)}
                      </div>
                    </div>
                    <div>
                      <div className="text-xs text-base-content/60">Packet Loss</div>
                      <div className={`text-lg font-semibold ${(m.packet_loss ?? 0) < 5 ? 'text-info' : 'text-error'}`}>
                        {formatPercent(m.packet_loss)}
                      </div>
                    </div>
                  </div>

                  {/* Detailed Metrics */}
                  <div className="collapse collapse-arrow bg-base-200">
                    <input type="checkbox" />
                    <div className="collapse-title text-sm font-medium">
                      Show detailed metrics
                    </div>
                <div className="collapse-content">
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 text-sm">
                    {/* Ping Details */}
                    <div>
                      <h4 className="font-semibold mb-2">Ping Details</h4>
                      <div className="space-y-1">
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Jitter:</span>
                          <span>{formatLatency(m.ping_jitter)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Low:</span>
                          <span>{formatLatency(m.ping_low)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">High:</span>
                          <span>{formatLatency(m.ping_high)}</span>
                        </div>
                      </div>
                    </div>

                    {/* Download Details */}
                    <div>
                      <h4 className="font-semibold mb-2">Download Details</h4>
                      <div className="space-y-1">
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Bytes:</span>
                          <span>{m.download_bytes?.toLocaleString() || 'N/A'}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Elapsed:</span>
                          <span>{m.download_elapsed ? `${m.download_elapsed}ms` : 'N/A'}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Latency IQM:</span>
                          <span>{formatLatency(m.download_latency_iqm)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Latency Low:</span>
                          <span>{formatLatency(m.download_latency_low)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Latency High:</span>
                          <span>{formatLatency(m.download_latency_high)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Latency Jitter:</span>
                          <span>{formatLatency(m.download_latency_jitter)}</span>
                        </div>
                      </div>
                    </div>

                    {/* Upload Details */}
                    <div>
                      <h4 className="font-semibold mb-2">Upload Details</h4>
                      <div className="space-y-1">
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Bytes:</span>
                          <span>{m.upload_bytes?.toLocaleString() || 'N/A'}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Elapsed:</span>
                          <span>{m.upload_elapsed ? `${m.upload_elapsed}ms` : 'N/A'}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Latency IQM:</span>
                          <span>{formatLatency(m.upload_latency_iqm)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Latency Low:</span>
                          <span>{formatLatency(m.upload_latency_low)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Latency High:</span>
                          <span>{formatLatency(m.upload_latency_high)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Latency Jitter:</span>
                          <span>{formatLatency(m.upload_latency_jitter)}</span>
                        </div>
                      </div>
                    </div>

                    {/* Network Interface */}
                    {(m.interface_name || m.interface_internal_ip) && (
                      <div>
                        <h4 className="font-semibold mb-2">Network Interface</h4>
                        <div className="space-y-1">
                          {m.interface_name && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">Name:</span>
                              <span>{m.interface_name}</span>
                            </div>
                          )}
                          {m.interface_internal_ip && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">Internal IP:</span>
                              <span>{m.interface_internal_ip}</span>
                            </div>
                          )}
                          {m.interface_external_ip && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">External IP:</span>
                              <span>{m.interface_external_ip}</span>
                            </div>
                          )}
                          {m.interface_mac && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">MAC:</span>
                              <span className="font-mono text-xs">{m.interface_mac}</span>
                            </div>
                          )}
                          {m.interface_is_vpn !== undefined && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">VPN:</span>
                              <span>{m.interface_is_vpn ? 'Yes' : 'No'}</span>
                            </div>
                          )}
                        </div>
                      </div>
                    )}

                    {/* Server Info */}
                    {m.server_host && (
                      <div>
                        <h4 className="font-semibold mb-2">Server Info</h4>
                        <div className="space-y-1">
                          <div className="flex justify-between">
                            <span className="text-base-content/60">Host:</span>
                            <span>{m.server_host}</span>
                          </div>
                          {m.server_name && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">Name:</span>
                              <span>{m.server_name}</span>
                            </div>
                          )}
                          {m.server_location && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">Location:</span>
                              <span>{m.server_location}</span>
                            </div>
                          )}
                          {m.server_country && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">Country:</span>
                              <span>{m.server_country}</span>
                            </div>
                          )}
                          {m.server_ip && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">IP:</span>
                              <span>{m.server_ip}</span>
                            </div>
                          )}
                          {m.server_port && (
                            <div className="flex justify-between">
                              <span className="text-base-content/60">Port:</span>
                              <span>{m.server_port}</span>
                            </div>
                          )}
                        </div>
                      </div>
                    )}

                    {/* ISP */}
                    {m.isp && (
                      <div>
                        <h4 className="font-semibold mb-2">ISP</h4>
                        <div>{m.isp}</div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
                </>
              )}
            </div>
          </Card>
        ))}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex justify-center items-center gap-2 mt-6">
          <button
            onClick={prevPage}
            disabled={page === 1}
            className="btn btn-sm"
          >
            <ChevronLeft size={16} />
            Previous
          </button>

          <div className="flex gap-1">
            {/* Show first page */}
            {page > 3 && (
              <>
                <button onClick={() => goToPage(1)} className="btn btn-sm">
                  1
                </button>
                {page > 4 && <span className="px-2 py-1">...</span>}
              </>
            )}

            {/* Show pages around current */}
            {Array.from({ length: totalPages }, (_, i) => i + 1)
              .filter((p) => p >= page - 2 && p <= page + 2)
              .map((p) => (
                <button
                  key={p}
                  onClick={() => goToPage(p)}
                  className={`btn btn-sm ${p === page ? 'btn-primary' : ''}`}
                >
                  {p}
                </button>
              ))}

            {/* Show last page */}
            {page < totalPages - 2 && (
              <>
                {page < totalPages - 3 && <span className="px-2 py-1">...</span>}
                <button onClick={() => goToPage(totalPages)} className="btn btn-sm">
                  {totalPages}
                </button>
              </>
            )}
          </div>

          <button
            onClick={nextPage}
            disabled={page === totalPages}
            className="btn btn-sm"
          >
            Next
            <ChevronRight size={16} />
          </button>
        </div>
      )}
    </div>
  );
}
