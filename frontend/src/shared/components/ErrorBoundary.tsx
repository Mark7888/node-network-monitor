import { Component, ReactNode, ErrorInfo } from 'react';
import Button from './ui/Button';
import { AlertTriangle } from 'lucide-react';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

/**
 * Error boundary component to catch and display errors
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
      return (
        <div className="min-h-screen flex items-center justify-center bg-base-200 p-4">
          <div className="card bg-base-100 shadow-xl max-w-md w-full">
            <div className="card-body text-center">
              <div className="flex justify-center mb-4">
                <AlertTriangle size={64} className="text-error" />
              </div>
              <h2 className="card-title justify-center text-2xl">
                Something went wrong
              </h2>
              <p className="text-base-content/70 my-4">
                An unexpected error occurred. Please try refreshing the page or contact support if the problem persists.
              </p>
              {this.state.error && (
                <div className="bg-base-200 p-3 rounded text-left text-sm font-mono overflow-auto max-h-32">
                  {this.state.error.message}
                </div>
              )}
              <div className="card-actions justify-center mt-6">
                <Button onClick={this.handleReset} variant="primary">
                  Go to Dashboard
                </Button>
                <Button onClick={() => window.location.reload()} variant="ghost">
                  Refresh Page
                </Button>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
