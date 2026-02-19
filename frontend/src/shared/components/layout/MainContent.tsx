import { ReactNode } from 'react';

interface MainContentProps {
  children: ReactNode;
}

/**
 * Main content area wrapper
 * Responsive: Adds bottom padding on mobile for bottom navigation
 */
export default function MainContent({ children }: MainContentProps) {
  return (
    <main className="flex-1 bg-base-200 overflow-auto">
      <div className="container mx-auto p-3 md:p-6 max-w-7xl pb-20 md:pb-6">
        {children}
      </div>
    </main>
  );
}
