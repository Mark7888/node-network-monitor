import { ReactNode } from 'react';

interface MainContentProps {
  children: ReactNode;
}

/**
 * Main content area wrapper
 */
export default function MainContent({ children }: MainContentProps) {
  return (
    <main className="flex-1 bg-base-200 overflow-auto">
      <div className="container mx-auto p-6 max-w-7xl">
        {children}
      </div>
    </main>
  );
}
