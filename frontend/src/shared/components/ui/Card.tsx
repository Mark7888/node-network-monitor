import { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  title?: string;
  actions?: ReactNode;
  className?: string;
  compact?: boolean;
}

/**
 * Reusable Card component using daisyUI styles
 */
export default function Card({ 
  children, 
  title, 
  actions, 
  className = '',
  compact = false 
}: CardProps) {
  return (
    <div className={`card bg-base-100 shadow-md ${className}`}>
      <div className={`card-body ${compact ? 'p-4' : ''}`}>
        {title && (
          <h2 className="card-title text-lg font-semibold">{title}</h2>
        )}
        {children}
        {actions && (
          <div className="card-actions justify-end mt-4">
            {actions}
          </div>
        )}
      </div>
    </div>
  );
}
