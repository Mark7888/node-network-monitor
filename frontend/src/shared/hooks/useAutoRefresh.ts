import { useEffect, useRef } from 'react';

/**
 * Hook for auto-refreshing data at a specified interval
 */
export function useAutoRefresh(callback: () => void, interval: number, enabled: boolean = true) {
  const savedCallback = useRef(callback);

  // Update saved callback if it changes
  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  // Setup interval
  useEffect(() => {
    if (!enabled) return;

    const tick = () => savedCallback.current();
    const id = setInterval(tick, interval);
    
    return () => clearInterval(id);
  }, [interval, enabled]);
}
