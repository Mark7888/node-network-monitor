import { useEffect } from 'react';
import { useNodes } from '../hooks/useNodes';
import { useAutoRefresh } from '@/shared/hooks/useAutoRefresh';
import NodeCard from './NodeCard';
import Spinner from '@/shared/components/ui/Spinner';
import ErrorMessage from '@/shared/components/ui/ErrorMessage';
import EmptyState from '@/shared/components/ui/EmptyState';
import { Server } from 'lucide-react';
import env from '@/core/config/env';

/**
 * Nodes list page component
 */
export default function NodesPage() {
  const { nodes, isLoading, error, fetchNodes } = useNodes();

  // Initial fetch
  useEffect(() => {
    fetchNodes();
  }, [fetchNodes]);

  // Auto-refresh
  useAutoRefresh(fetchNodes, env.refreshInterval);

  if (isLoading && nodes.length === 0) {
    return <Spinner message="Loading nodes..." />;
  }

  if (error) {
    return <ErrorMessage message={error} onRetry={fetchNodes} />;
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold">Nodes</h1>
        <p className="text-base-content/60 mt-1">
          Manage and monitor all speedtest nodes
        </p>
      </div>

      {/* Nodes Grid */}
      {nodes.length === 0 ? (
        <EmptyState
          title="No nodes found"
          message="No speedtest nodes have registered yet. Deploy a node to start collecting measurements."
          icon={<Server size={64} />}
        />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {nodes.map((node) => (
            <NodeCard key={node.id} node={node} />
          ))}
        </div>
      )}
    </div>
  );
}
