import { ReactNode } from 'react';
import { FileQuestion } from 'lucide-react';

interface EmptyStateProps {
  title: string;
  message?: string;
  icon?: ReactNode;
  action?: ReactNode;
}

/**
 * Empty state component for displaying when no data is available
 */
export default function EmptyState({ 
  title, 
  message, 
  icon,
  action 
}: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center p-12 text-center">
      <div className="text-base-content/30 mb-4">
        {icon || <FileQuestion size={64} />}
      </div>
      <h3 className="text-lg font-semibold mb-2">{title}</h3>
      {message && (
        <p className="text-base-content/60 mb-4 max-w-md">{message}</p>
      )}
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}
