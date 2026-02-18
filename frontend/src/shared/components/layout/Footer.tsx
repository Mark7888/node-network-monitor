import { useState, useEffect } from 'react';
import { getHealth } from '@/shared/services/healthService';

/**
 * Footer component displaying version info
 */
export default function Footer() {
  const [version, setVersion] = useState<string>('');

  useEffect(() => {
    const fetchVersion = async () => {
      try {
        const health = await getHealth();
        setVersion(health.version);
      } catch (error) {
        // Silently fail - footer version is not critical
        console.debug('Failed to fetch version:', error);
      }
    };

    fetchVersion();
  }, []);

  return (
    <footer className="bg-base-100 border-t border-base-300 px-6 py-3">
      <div className="flex justify-between items-center text-xs text-base-content/60">
        <div>
          Network Speed Measurement System
        </div>
        {version && (
          <div>
            Version: <span className="font-mono">{version}</span>
          </div>
        )}
      </div>
    </footer>
  );
}
