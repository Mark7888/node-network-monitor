import { Node } from '../types/node.types';
import { formatRelativeTime } from '@/shared/utils/date';
import Badge from '@/shared/components/ui/Badge';
import Card from '@/shared/components/ui/Card';
import Button from '@/shared/components/ui/Button';
import { Link } from 'react-router-dom';
import { Clock, Star, Archive, Trash2, ArchiveRestore } from 'lucide-react';
import { useNodes } from '../hooks/useNodes';
import { useState } from 'react';

interface NodeCardProps {
  node: Node;
  isArchived?: boolean;
}

/**
 * Node card component displaying node summary
 */
export default function NodeCard({ node, isArchived = false }: NodeCardProps) {
  const { toggleFavorite, archiveNode, deleteNode } = useNodes();
  const [isDeleting, setIsDeleting] = useState(false);

  const badgeVariant = node.status === 'active' ? 'success' : 
                       node.status === 'unreachable' ? 'warning' : 'ghost';

  const handleToggleFavorite = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    await toggleFavorite(node.id);
  };

  const handleArchiveToggle = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    await archiveNode(node.id, !node.archived);
  };

  const handleDelete = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    const confirmed = window.confirm(
      `Are you sure you want to delete "${node.name}"? This will permanently delete all measurements and cannot be undone.`
    );
    
    if (!confirmed) return;
    
    setIsDeleting(true);
    try {
      await deleteNode(node.id);
    } catch (error) {
      // Error is already handled in store
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <Card>
      <div className="flex justify-between items-start mb-3">
        <div className="flex-1">
          <div className="flex items-center gap-2">
            <h3 className="text-lg font-semibold">{node.name}</h3>
            {node.favorite && !isArchived && (
              <Star size={16} className="text-warning fill-warning" />
            )}
          </div>
          <p className="text-sm text-base-content/60">ID: {node.id.substring(0, 8)}...</p>
        </div>
        <Badge variant={badgeVariant}>
          {node.status.charAt(0).toUpperCase() + node.status.slice(1)}
        </Badge>
      </div>

      <div className="space-y-2 text-sm">
        <div className="flex items-center gap-2 text-base-content/70">
          <Clock size={16} />
          <span>Last seen: {formatRelativeTime(node.last_seen)}</span>
        </div>
        {node.location && (
          <div className="text-base-content/70">
            📍 {node.location}
          </div>
        )}
      </div>

      <div className="mt-4 flex items-center gap-2">
        <Link to={`/nodes/${node.id}`} className="flex-1">
          <Button variant="outline" size="sm" className="w-full">
            View Details →
          </Button>
        </Link>
        
        {!isArchived && (
          <Button
            variant="ghost"
            size="sm"
            onClick={handleToggleFavorite}
            className="px-2"
            title={node.favorite ? 'Remove from favorites' : 'Add to favorites'}
          >
            <Star size={18} className={node.favorite ? 'fill-warning text-warning' : ''} />
          </Button>
        )}
        
        <Button
          variant="ghost"
          size="sm"
          onClick={handleArchiveToggle}
          className="px-2"
          title={node.archived ? 'Unarchive node' : 'Archive node'}
        >
          {node.archived ? <ArchiveRestore size={18} /> : <Archive size={18} />}
        </Button>
        
        <Button
          variant="ghost"
          size="sm"
          onClick={handleDelete}
          className="px-2 text-error hover:bg-error/10"
          loading={isDeleting}
          disabled={isDeleting}
          title="Delete node"
        >
          <Trash2 size={18} />
        </Button>
      </div>
    </Card>
  );
}
