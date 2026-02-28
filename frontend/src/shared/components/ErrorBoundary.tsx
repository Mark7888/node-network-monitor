import { Component, ReactNode, ErrorInfo, useState } from 'react';
import Button from './ui/Button';
import { Home, RefreshCw } from 'lucide-react';

// ---------------------------------------------------------------------------
// Fallback UI — functional so we can use hooks (useState for details toggle)
// ---------------------------------------------------------------------------
function ErrorFallback({
  error,
  onReset,
}: {
  error: Error | null;
  onReset: () => void;
}) {
  const [showDetails, setShowDetails] = useState(false);

  return (
    <div className="min-h-screen bg-base-200 flex items-center justify-center p-4 relative overflow-hidden">
      {/* Subtle dot-grid backdrop */}
      <div
        className="absolute inset-0 opacity-[0.035]"
        style={{
          backgroundImage: 'radial-gradient(circle, currentColor 1px, transparent 0)',
          backgroundSize: '28px 28px',
        }}
      />

      <div className="relative z-10 flex flex-col items-center text-center max-w-lg w-full">
        {/* Fault graphic — broken circuit */}
        <svg
          viewBox="0 0 220 210"
          className="w-52 h-52 mx-auto"
          fill="none"
          aria-hidden="true"
        >
          {/* Stable connections */}
          <line stroke="currentColor" strokeWidth="1.5" strokeDasharray="4 3" className="text-primary/50"
            x1="110" y1="108" x2="50" y2="55">
            <animate attributeName="stroke-dashoffset" from="0" to="-14" dur="1s" repeatCount="indefinite" />
          </line>
          <line stroke="currentColor" strokeWidth="1.5" strokeDasharray="4 3" className="text-primary/50"
            x1="110" y1="108" x2="50" y2="162">
            <animate attributeName="stroke-dashoffset" from="0" to="-14" dur="0.9s" repeatCount="indefinite" />
          </line>

          {/* Two broken connections */}
          <line stroke="currentColor" strokeWidth="1.5" strokeDasharray="3 6" className="text-error/70"
            x1="110" y1="108" x2="170" y2="55" />
          <line stroke="currentColor" strokeWidth="1.5" strokeDasharray="3 6" className="text-error/70"
            x1="110" y1="108" x2="170" y2="162" />

          {/* Active peripheral nodes */}
          <circle cx="50"  cy="55"  r="9" fill="currentColor" className="text-primary/80" />
          <circle cx="50"  cy="162" r="9" fill="currentColor" className="text-primary/80" />

          {/* Center hub — warning pulse */}
          <circle cx="110" cy="108" r="20" stroke="currentColor" strokeWidth="1" fill="none" className="text-warning/20">
            <animate attributeName="r"       values="20;26;20" dur="1.5s" repeatCount="indefinite" />
            <animate attributeName="opacity" values="0.3;0;0.3" dur="1.5s" repeatCount="indefinite" />
          </circle>
          <circle cx="110" cy="108" r="15" stroke="currentColor" strokeWidth="2" fill="none" className="text-warning/50" />
          <circle cx="110" cy="108" r="9"  fill="currentColor" className="text-warning" />

          {/* Fault indicator on hub */}
          <line x1="110" y1="101" x2="110" y2="109" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" className="text-warning-content" />
          <circle cx="110" cy="113" r="1.5" fill="currentColor" className="text-warning-content" />

          {/* Broken nodes */}
          <circle cx="170" cy="55"  r="12" stroke="currentColor" strokeWidth="1.5" fill="none" strokeDasharray="3 2" className="text-error/60">
            <animate attributeName="opacity" values="0.6;0.2;0.6" dur="1.5s" repeatCount="indefinite" />
          </circle>
          <line x1="164" y1="49"  x2="176" y2="61"  stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" className="text-error" />
          <line x1="176" y1="49"  x2="164" y2="61"  stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" className="text-error" />

          <circle cx="170" cy="162" r="12" stroke="currentColor" strokeWidth="1.5" fill="none" strokeDasharray="3 2" className="text-error/60">
            <animate attributeName="opacity" values="0.6;0.2;0.6" dur="2s" repeatCount="indefinite" />
          </circle>
          <line x1="164" y1="156" x2="176" y2="168" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" className="text-error" />
          <line x1="176" y1="156" x2="164" y2="168" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" className="text-error" />
        </svg>

        {/* Label */}
        <div className="text-[7rem] font-black leading-none tracking-tighter text-transparent bg-clip-text bg-gradient-to-br from-warning to-error select-none -mt-4 mb-2">
          Fault
        </div>

        <h1 className="text-2xl font-bold mb-3 text-base-content">System error detected</h1>
        <p className="text-base-content/60 mb-8 leading-relaxed max-w-sm">
          {error?.message
            ? error.message
            : 'An unexpected runtime error occurred. The component tree could not be rendered.'}
        </p>

        <div className="flex flex-wrap gap-3 justify-center">
          <Button onClick={onReset} variant="primary">
            <Home size={16} />
            Go to Dashboard
          </Button>
          <Button onClick={() => window.location.reload()} variant="ghost">
            <RefreshCw size={16} />
            Refresh Page
          </Button>
        </div>

        {error?.stack && (
          <div className="mt-8 w-full text-left">
            <button
              onClick={() => setShowDetails(v => !v)}
              className="text-xs text-base-content/40 hover:text-base-content/70 transition-colors underline-offset-2 hover:underline"
            >
              {showDetails ? 'Hide' : 'Show'} stack trace
            </button>
            {showDetails && (
              <pre className="mt-2 p-4 bg-base-300 rounded-lg text-xs font-mono overflow-auto max-h-48 text-base-content/70 whitespace-pre-wrap break-all">
                {error.stack}
              </pre>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Class component — required by React to catch render-phase errors
// ---------------------------------------------------------------------------
interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

/**
 * Error boundary component — wraps subtrees to catch and display errors.
 */
export default class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  handleReset = () => {
    this.setState({ hasError: false, error: null });
    window.location.href = '/dashboard';
  };

  render() {
    if (this.state.hasError) {
      return <ErrorFallback error={this.state.error} onReset={this.handleReset} />;
    }
    return this.props.children;
  }
}
