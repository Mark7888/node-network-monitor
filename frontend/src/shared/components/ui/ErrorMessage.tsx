import { AlertCircle } from 'lucide-react';

interface ErrorMessageProps {
  message: string;
  onRetry?: () => void;
}

/**
 * Error message component for displaying errors
 */
export default function ErrorMessage({ message, onRetry }: ErrorMessageProps) {
  return (
    <div className="alert alert-error shadow-lg">
      <div className="flex items-start gap-2">
        <AlertCircle size={20} className="flex-shrink-0 mt-0.5" />
        <div className="flex-1">
          <div className="font-semibold">Error</div>
          <div className="text-sm">{message}</div>
        </div>
      </div>
      {onRetry && (
        <button onClick={onRetry} className="btn btn-sm btn-ghost">
          Retry
        </button>
      )}
    </div>
  );
}
