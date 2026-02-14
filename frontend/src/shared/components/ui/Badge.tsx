import { ReactNode } from 'react';

interface BadgeProps {
  children: ReactNode;
  variant?: 'default' | 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'ghost' | 'info';
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

/**
 * Reusable Badge component using daisyUI styles
 */
export default function Badge({ 
  children, 
  variant = 'default', 
  size = 'md',
  className = '' 
}: BadgeProps) {
  const variantClass = {
    default: '',
    primary: 'badge-primary',
    secondary: 'badge-secondary',
    success: 'badge-success',
    warning: 'badge-warning',
    error: 'badge-error',
    ghost: 'badge-ghost',
    info: 'badge-info',
  }[variant];

  const sizeClass = {
    sm: 'badge-sm',
    md: '',
    lg: 'badge-lg',
  }[size];

  return (
    <span className={`badge ${variantClass} ${sizeClass} ${className}`}>
      {children}
    </span>
  );
}
