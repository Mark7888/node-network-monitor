import { useState } from 'react';
import { useRouteError, isRouteErrorResponse, useNavigate } from 'react-router-dom';
import { Home, RefreshCw } from 'lucide-react';
import Button from './ui/Button';

// ---------------------------------------------------------------------------
// Animated network topology graphic — fits the "network monitor" app theme
// ---------------------------------------------------------------------------
function NetworkGraphic() {
  return (
    <svg
      viewBox="0 0 220 210"
      className="w-52 h-52 mx-auto"
      fill="none"
      aria-hidden="true"
    >
      {/* ── Active connection lines (animated travelling dashes) ── */}
      <line
        stroke="currentColor"
        strokeWidth="1.5"
        strokeDasharray="4 3"
        className="text-primary/50"
        x1="110" y1="108" x2="50" y2="55"
      >
        <animate attributeName="stroke-dashoffset" from="0" to="-14" dur="1s" repeatCount="indefinite" />
      </line>
      <line
        stroke="currentColor"
        strokeWidth="1.5"
        strokeDasharray="4 3"
        className="text-primary/50"
        x1="110" y1="108" x2="170" y2="55"
      >
        <animate attributeName="stroke-dashoffset" from="0" to="-14" dur="1.4s" repeatCount="indefinite" />
      </line>
      <line
        stroke="currentColor"
        strokeWidth="1.5"
        strokeDasharray="4 3"
        className="text-primary/50"
        x1="110" y1="108" x2="50" y2="162"
      >
        <animate attributeName="stroke-dashoffset" from="0" to="-14" dur="0.9s" repeatCount="indefinite" />
      </line>

      {/* ── Broken / missing connection ── */}
      <line
        stroke="currentColor"
        strokeWidth="1.5"
        strokeDasharray="3 6"
        className="text-error/70"
        x1="110" y1="108" x2="170" y2="162"
      />

      {/* ── Active peripheral nodes ── */}
      <circle cx="50"  cy="55"  r="9" fill="currentColor" className="text-primary/80" />
      <circle cx="170" cy="55"  r="9" fill="currentColor" className="text-primary/80" />
      <circle cx="50"  cy="162" r="9" fill="currentColor" className="text-primary/80" />

      {/* ── Center hub — pulsing ripple ring ── */}
      <circle cx="110" cy="108" r="20" stroke="currentColor" strokeWidth="1" fill="none" className="text-primary/20">
        <animate attributeName="r"       values="20;26;20" dur="2s" repeatCount="indefinite" />
        <animate attributeName="opacity" values="0.2;0;0.2" dur="2s" repeatCount="indefinite" />
      </circle>
      <circle cx="110" cy="108" r="15" stroke="currentColor" strokeWidth="2" fill="none" className="text-primary/40" />
      <circle cx="110" cy="108" r="9"  fill="currentColor"  className="text-primary" />

      {/* ── Broken / unreachable node ── */}
      <circle
        cx="170" cy="162" r="12"
        stroke="currentColor" strokeWidth="1.5" fill="none"
        strokeDasharray="3 2"
        className="text-error/60"
      >
        <animate attributeName="opacity" values="0.6;0.2;0.6" dur="1.5s" repeatCount="indefinite" />
      </circle>
      {/* X mark */}
      <line x1="164" y1="156" x2="176" y2="168" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" className="text-error" />
      <line x1="176" y1="156" x2="164" y2="168" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" className="text-error" />
    </svg>
  );
}

// ---------------------------------------------------------------------------
// Router-level error page — used as the `errorElement` on routes
// ---------------------------------------------------------------------------
export default function ErrorPage() {
  const error = useRouteError();
  const navigate = useNavigate();
  const [showDetails, setShowDetails] = useState(false);

  let status: number | null = null;
  let title = 'Something went wrong';
  let description = 'An unexpected error occurred. Please try again.';
  let details: string | null = null;

  // `error` is undefined when this component is rendered as a normal route
  // element (e.g. the catch-all `*` route) rather than as an errorElement.
  if (!error) {
    status = 404;
    title = 'Node not found';
    description =
      "We scanned the entire network but couldn't locate this page. It may have been moved or removed.";
  } else if (isRouteErrorResponse(error)) {
    status = error.status;
    switch (error.status) {
      case 404:
        title = 'Node not found';
        description =
          "We scanned the entire network but couldn't locate this page. It may have been moved or removed.";
        break;
      case 401:
      case 403:
        title = 'Access denied';
        description = "You don't have permission to access this resource.";
        break;
      default:
        title = error.statusText || title;
        description = 'The server returned an error. Please try again later.';
    }
    if (error.data) details = String(error.data);
  } else if (error instanceof Error) {
    title = 'Application error';
    description = error.message;
    details = error.stack ?? null;
  }

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
        <NetworkGraphic />

        {/* Big status code */}
        {status !== null && (
          <div className="text-[7rem] font-black leading-none tracking-tighter text-transparent bg-clip-text bg-gradient-to-br from-primary to-accent select-none -mt-4 mb-2">
            {status}
          </div>
        )}

        <h1 className="text-2xl font-bold mb-3 text-base-content">{title}</h1>
        <p className="text-base-content/60 mb-8 leading-relaxed max-w-sm">{description}</p>

        <div className="flex flex-wrap gap-3 justify-center">
          <Button onClick={() => navigate('/dashboard')} variant="primary">
            <Home size={16} />
            Go to Dashboard
          </Button>
          <Button onClick={() => window.location.reload()} variant="ghost">
            <RefreshCw size={16} />
            Try Again
          </Button>
        </div>

        {/* Technical details (collapsed by default) */}
        {details && (
          <div className="mt-8 w-full text-left">
            <button
              onClick={() => setShowDetails(v => !v)}
              className="text-xs text-base-content/40 hover:text-base-content/70 transition-colors underline-offset-2 hover:underline"
            >
              {showDetails ? 'Hide' : 'Show'} technical details
            </button>
            {showDetails && (
              <pre className="mt-2 p-4 bg-base-300 rounded-lg text-xs font-mono overflow-auto max-h-48 text-base-content/70 whitespace-pre-wrap break-all">
                {details}
              </pre>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
