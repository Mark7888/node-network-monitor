interface SpinnerProps {
  size?: 'xs' | 'sm' | 'md' | 'lg';
  fullScreen?: boolean;
  message?: string;
}

/**
 * Reusable Loading Spinner component using daisyUI styles
 */
export default function Spinner({ 
  size = 'md', 
  fullScreen = false,
  message 
}: SpinnerProps) {
  const sizeClass = {
    xs: 'loading-xs',
    sm: 'loading-sm',
    md: 'loading-md',
    lg: 'loading-lg',
  }[size];

  const spinner = (
    <div className="flex flex-col items-center gap-4">
      <span className={`loading loading-spinner ${sizeClass} text-primary`}></span>
      {message && <p className="text-sm text-base-content/60">{message}</p>}
    </div>
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-base-100">
        {spinner}
      </div>
    );
  }

  return (
    <div className="flex justify-center items-center p-8">
      {spinner}
    </div>
  );
}
