import { useCallback, useEffect, useRef, useState } from 'react';
import { getStatus, type Status } from './api';

const POLL_INTERVAL_MS = 2000;

export type UseStatusResult = {
  status: Status | null;
  error: string | null;
  loading: boolean;
  refresh: () => void;
};

export function useStatus(): UseStatusResult {
  const [status, setStatus] = useState<Status | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const mounted = useRef(true);

  const refresh = useCallback(async () => {
    try {
      const next = await getStatus();
      if (!mounted.current) return;
      setStatus(next);
      setError(null);
    } catch (e) {
      if (!mounted.current) return;
      setError(e instanceof Error ? e.message : 'Failed to fetch status');
    } finally {
      if (mounted.current) setLoading(false);
    }
  }, []);

  useEffect(() => {
    mounted.current = true;
    refresh();
    const id = window.setInterval(refresh, POLL_INTERVAL_MS);
    return () => {
      mounted.current = false;
      window.clearInterval(id);
    };
  }, [refresh]);

  return { status, error, loading, refresh };
}
