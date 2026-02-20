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
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <h3 className="text-base md:text-lg font-semibold truncate">{node.name}</h3>
            {node.favorite && !isArchived && (
              <Star size={14} className="text-warning fill-warning flex-shrink-0 md:w-4 md:h-4" />
            )}
          </div>
          <p className="text-xs md:text-sm text-base-content/60 truncate">ID: {node.id.substring(0, 8)}...</p>
        </div>
        <Badge variant={badgeVariant} className="flex-shrink-0">
          {node.status.charAt(0).toUpperCase() + node.status.slice(1)}
        </Badge>
      </div>

      <div className="space-y-2 text-xs md:text-sm">
        <div className="flex items-center gap-2 text-base-content/70">
          <Clock size={14} className="flex-shrink-0 md:w-4 md:h-4" />
          <span className="truncate">Last seen: {formatRelativeTime(node.last_seen)}</span>
        </div>
        {node.location && (
          <div className="text-base-content/70 truncate">
            📍 {node.location}
          </div>
        )}
      </div>

      <div className="mt-4 flex flex-wrap items-center gap-2">
        <Link to={`/nodes/${node.id}`} className="flex-1 min-w-[130px]">
          <Button variant="outline" size="sm" className="w-full text-xs md:text-sm">
            View Details →
          </Button>
        </Link>
        
        <div className="flex flex-1 items-center gap-2">
          {!isArchived && (
            <Button
              variant="ghost"
              size="sm"
              onClick={handleToggleFavorite}
              className="flex-1 justify-center min-w-[2.5rem]"
              title={node.favorite ? 'Remove from favorites' : 'Add to favorites'}
            >
              <Star size={16} className={node.favorite ? 'fill-warning text-warning' : ''} />
            </Button>
          )}
          
          <Button
            variant="ghost"
            size="sm"
            onClick={handleArchiveToggle}
            className="flex-1 justify-center min-w-[2.5rem]"
            title={node.archived ? 'Unarchive node' : 'Archive node'}
          >
            {node.archived ? <ArchiveRestore size={16} /> : <Archive size={16} />}
          </Button>
          
          <Button
            variant="ghost"
            size="sm"
            onClick={handleDelete}
            className="flex-1 justify-center min-w-[2.5rem] text-error hover:bg-error/10"
            loading={isDeleting}
            disabled={isDeleting}
            title="Delete node"
          >
            <Trash2 size={16} />
          </Button>
        </div>
      </div>
    </Card>
  );
}
