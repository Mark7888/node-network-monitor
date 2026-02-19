import { useEffect, useState, useMemo } from 'react';
import { useNodes } from '../hooks/useNodes';
import { useAutoRefresh } from '@/shared/hooks/useAutoRefresh';
import NodeCard from './NodeCard';
import Spinner from '@/shared/components/ui/Spinner';
import ErrorMessage from '@/shared/components/ui/ErrorMessage';
import EmptyState from '@/shared/components/ui/EmptyState';
import { Server, ArrowUp, ArrowDown } from 'lucide-react';
import env from '@/core/config/env';
import { Node } from '../types/node.types';

type TabType = 'all' | 'archived';
type StatusFilter = 'all' | 'active' | 'unreachable' | 'inactive';
type SortOption = 'name' | 'last_activity' | 'status';
type SortDirection = 'asc' | 'desc';

/**
 * Nodes list page component
 */
export default function NodesPage() {
  const { nodes, isLoading, error, fetchNodes } = useNodes();
  const [activeTab, setActiveTab] = useState<TabType>('all');
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('all');
  const [sortBy, setSortBy] = useState<SortOption>('last_activity');
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc');
  const [showFavoritesSeparately, setShowFavoritesSeparately] = useState(true);

  // Handle sort option change with direction toggle
  const handleSortChange = (newSortBy: SortOption) => {
    if (sortBy === newSortBy) {
      // Toggle direction if clicking the same sort option
      setSortDirection(prev => prev === 'asc' ? 'desc' : 'asc');
    } else {
      // Set new sort option with default direction
      setSortBy(newSortBy);
      // Default direction based on sort type
      if (newSortBy === 'name') {
        setSortDirection('asc');
      } else {
        setSortDirection('desc');
      }
    }
  };

  // Initial fetch
  useEffect(() => {
    fetchNodes();
  }, [fetchNodes]);

  // Auto-refresh
  useAutoRefresh(fetchNodes, env.refreshInterval);

  // Filter and sort logic
  const { favoriteNodes, regularNodes, displayNodes } = useMemo(() => {
    if (!nodes) return { favoriteNodes: [], regularNodes: [], displayNodes: [] };

    // Filter by archived status based on active tab
    let filtered = nodes.filter(node => 
      activeTab === 'archived' ? node.archived : !node.archived
    );

    // Filter by status
    if (statusFilter !== 'all') {
      filtered = filtered.filter(node => node.status === statusFilter);
    }

    // Sort nodes
    const sortNodes = (nodesToSort: Node[]) => {
      return [...nodesToSort].sort((a, b) => {
        let comparison = 0;
        
        switch (sortBy) {
          case 'name':
            comparison = a.name.localeCompare(b.name);
            break;
          case 'status':
            comparison = a.status.localeCompare(b.status);
            break;
          case 'last_activity':
          default:
            // For dates, positive means a is newer (higher timestamp)
            comparison = new Date(a.last_alive).getTime() - new Date(b.last_alive).getTime();
            break;
        }
        
        // Apply sort direction
        return sortDirection === 'asc' ? comparison : -comparison;
      });
    };

    // Separate favorites and regular nodes (only for non-archived tab)
    if (activeTab === 'all' && showFavoritesSeparately) {
      const favorites = sortNodes(filtered.filter(node => node.favorite));
      const regular = sortNodes(filtered.filter(node => !node.favorite));
      return {
        favoriteNodes: favorites,
        regularNodes: regular,
        displayNodes: filtered,
      };
    }

    // If not showing favorites separately, just return sorted nodes
    return {
      favoriteNodes: [],
      regularNodes: sortNodes(filtered),
      displayNodes: filtered,
    };
  }, [nodes, activeTab, statusFilter, sortBy, sortDirection, showFavoritesSeparately]);

  const hasFavorites = nodes?.some(node => node.favorite && !node.archived) || false;

  if (isLoading && (!nodes || nodes.length === 0)) {
    return <Spinner message="Loading nodes..." />;
  }

  if (error) {
    return <ErrorMessage message={error} onRetry={fetchNodes} />;
  }

  return (
    <div className="space-y-4 md:space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl md:text-3xl font-bold">Nodes</h1>
        <p className="text-sm md:text-base text-base-content/60 mt-1">
          Manage and monitor all speedtest nodes
        </p>
      </div>

      {/* Tabs */}
      <div className="tabs tabs-boxed w-fit">
        <button 
          className={`tab tab-sm md:tab-md ${activeTab === 'all' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('all')}
        >
          All Nodes
        </button>
        <button 
          className={`tab tab-sm md:tab-md ${activeTab === 'archived' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('archived')}
        >
          Archived
        </button>
      </div>

      {/* Filters and Controls */}
      <div className="flex flex-col md:flex-row md:flex-wrap items-start md:items-center gap-3 md:gap-4">
        {/* Status Filter */}
        <div className="flex items-center gap-2 w-full md:w-auto">
          <span className="text-xs md:text-sm font-medium whitespace-nowrap">Status:</span>
          <div className="join overflow-x-auto w-full">
            <button 
              className={`btn btn-xs md:btn-sm join-item ${statusFilter === 'all' ? 'btn-primary' : 'btn-ghost'}`}
              onClick={() => setStatusFilter('all')}
            >
              All
            </button>
            <button 
              className={`btn btn-xs md:btn-sm join-item ${statusFilter === 'active' ? 'btn-primary' : 'btn-ghost'}`}
              onClick={() => setStatusFilter('active')}
            >
              Active
            </button>
            <button 
              className={`btn btn-xs md:btn-sm join-item ${statusFilter === 'unreachable' ? 'btn-primary' : 'btn-ghost'}`}
              onClick={() => setStatusFilter('unreachable')}
            >
              Unreachable
            </button>
            <button 
              className={`btn btn-xs md:btn-sm join-item ${statusFilter === 'inactive' ? 'btn-primary' : 'btn-ghost'}`}
              onClick={() => setStatusFilter('inactive')}
            >
              Inactive
            </button>
          </div>
        </div>

        {/* Sort By */}
        <div className="flex items-center gap-2 w-full md:w-auto overflow-x-auto">
          <span className="text-xs md:text-sm font-medium whitespace-nowrap">Sort by:</span>
          <div className="join">
            <button 
              className={`btn btn-xs md:btn-sm join-item gap-1 whitespace-nowrap ${sortBy === 'last_activity' ? 'btn-primary' : 'btn-ghost'}`}
              onClick={() => handleSortChange('last_activity')}
            >
              <span className="hidden sm:inline">Last Activity</span>
              <span className="sm:hidden">Activity</span>
              {sortBy === 'last_activity' && (
                sortDirection === 'asc' ? <ArrowUp size={12} className="md:w-[14px] md:h-[14px]" /> : <ArrowDown size={12} className="md:w-[14px] md:h-[14px]" />
              )}
            </button>
            <button 
              className={`btn btn-xs md:btn-sm join-item gap-1 whitespace-nowrap ${sortBy === 'name' ? 'btn-primary' : 'btn-ghost'}`}
              onClick={() => handleSortChange('name')}
            >
              Name
              {sortBy === 'name' && (
                sortDirection === 'asc' ? <ArrowUp size={12} className="md:w-[14px] md:h-[14px]" /> : <ArrowDown size={12} className="md:w-[14px] md:h-[14px]" />
              )}
            </button>
            <button 
              className={`btn btn-xs md:btn-sm join-item gap-1 whitespace-nowrap ${sortBy === 'status' ? 'btn-primary' : 'btn-ghost'}`}
              onClick={() => handleSortChange('status')}
            >
              Status
              {sortBy === 'status' && (
                sortDirection === 'asc' ? <ArrowUp size={12} className="md:w-[14px] md:h-[14px]" /> : <ArrowDown size={12} className="md:w-[14px] md:h-[14px]" />
              )}
            </button>
          </div>
        </div>

        {/* Show Favorites Separately - Only show if there are favorites and on all tab */}
        {hasFavorites && activeTab === 'all' && (
          <div className="flex items-center gap-2">
            <label className="label cursor-pointer gap-2">
              <input
                type="checkbox"
                className="checkbox checkbox-xs md:checkbox-sm"
                checked={showFavoritesSeparately}
                onChange={(e) => setShowFavoritesSeparately(e.target.checked)}
              />
              <span className="label-text text-xs md:text-sm">Show favorites separately</span>
            </label>
          </div>
        )}
      </div>

      {/* Empty State */}
      {!displayNodes || displayNodes.length === 0 ? (
        <EmptyState
          title={activeTab === 'archived' ? 'No archived nodes' : 'No nodes found'}
          message={
            activeTab === 'archived'
              ? 'You have not archived any nodes yet.'
              : 'No speedtest nodes have registered yet. Deploy a node to start collecting measurements.'
          }
          icon={<Server size={64} />}
        />
      ) : (
        <>
          {/* Favorites Section */}
          {showFavoritesSeparately && favoriteNodes.length > 0 && activeTab === 'all' && (
            <div className="space-y-3 md:space-y-4">
              <h2 className="text-lg md:text-xl font-semibold flex items-center gap-2">
                ⭐ Favorite Nodes
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3 md:gap-4">
                {favoriteNodes.map((node) => (
                  <NodeCard key={node.id} node={node} />
                ))}
              </div>
            </div>
          )}

          {/* Regular Nodes Section */}
          {regularNodes.length > 0 && (
            <div className="space-y-3 md:space-y-4">
              {showFavoritesSeparately && favoriteNodes.length > 0 && activeTab === 'all' && (
                <h2 className="text-lg md:text-xl font-semibold">Other Nodes</h2>
              )}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3 md:gap-4">
                {regularNodes.map((node) => (
                  <NodeCard key={node.id} node={node} isArchived={activeTab === 'archived'} />
                ))}
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
