import { Node } from '../types/node.types';
import { STATUS_BADGE_CLASS } from '@/shared/utils/constants';
import { formatRelativeTime } from '@/shared/utils/date';
import { formatSpeed, formatLatency } from '@/shared/utils/format';
import Badge from '@/shared/components/ui/Badge';
import Card from '@/shared/components/ui/Card';
import Button from '@/shared/components/ui/Button';
import { Link } from 'react-router-dom';
import { Activity, Clock } from 'lucide-react';

interface NodeCardProps {
  node: Node;
}

/**
 * Node card component displaying node summary
 */
export default function NodeCard({ node }: NodeCardProps) {
  const badgeVariant = node.status === 'active' ? 'success' : 
                       node.status === 'unreachable' ? 'warning' : 'ghost';

  return (
    <Card>
      <div className="flex justify-between items-start mb-3">
        <div>
          <h3 className="text-lg font-semibold">{node.name}</h3>
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

      <div className="mt-4">
        <Link to={`/nodes/${node.id}`}>
          <Button variant="outline" size="sm" className="w-full">
            View Details →
          </Button>
        </Link>
      </div>
    </Card>
  );
}
