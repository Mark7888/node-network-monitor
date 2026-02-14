import { SelectHTMLAttributes, ReactNode } from 'react';

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  error?: string;
  children: ReactNode;
}

/**
 * Reusable Select component using daisyUI styles
 */
export default function Select({ 
  label, 
  error, 
  children, 
  className = '', 
  ...props 
}: SelectProps) {
  return (
    <div className="form-control w-full">
      {label && (
        <label className="label">
          <span className="label-text">{label}</span>
        </label>
      )}
      <select 
        className={`select select-bordered w-full ${error ? 'select-error' : ''} ${className}`}
        {...props}
      >
        {children}
      </select>
      {error && (
        <label className="label">
          <span className="label-text-alt text-error">{error}</span>
        </label>
      )}
    </div>
  );
}
