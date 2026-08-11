const BACKEND_URL =
  import.meta.env.VITE_BACKEND_URL ?? 'http://localhost:8080';

export type Status = {
  load_generation: 'idle' | 'running';
  cpu_load_active: boolean;
  memory_load_active: boolean;
  uptime_seconds: number;
};

type ApiError = {
  error: string;
  message: string;
};

type HealthResponse = { status: string };

type LoadStartResponse = { status: string };
type LoadStopResponse = { status: string };
type CpuLoadResponse = { status: string };
type MemoryLoadResponse = { status: string };

type TrafficRandomResponse = { requests_generated: number };
type TrafficRecordResponse = { status: string; result: string };

type ResetResponse = { status: string };

async function parseError(res: Response): Promise<Error> {
  let body: ApiError | null = null;
  try {
    body = (await res.json()) as ApiError;
  } catch {
    // body wasn't JSON
  }
  if (body && typeof body.message === 'string') {
    return new Error(body.message);
  }
  return new Error(`Request failed: ${res.status} ${res.statusText}`);
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BACKEND_URL}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
  });
  if (!res.ok) {
    throw await parseError(res);
  }
  // 202/200 with JSON body
  return (await res.json()) as T;
}

export function getHealth() {
  return request<HealthResponse>('/health', { method: 'GET' });
}

export function getStatus() {
  return request<Status>('/api/status', { method: 'GET' });
}

export function startLoad() {
  return request<LoadStartResponse>('/api/load/start', { method: 'POST' });
}

export function stopLoad() {
  return request<LoadStopResponse>('/api/load/stop', { method: 'POST' });
}

export function generateCpuLoad() {
  return request<CpuLoadResponse>('/api/load/cpu', { method: 'POST' });
}

export function generateMemoryLoad() {
  return request<MemoryLoadResponse>('/api/load/memory', { method: 'POST' });
}

export function generateRandomTraffic() {
  return request<TrafficRandomResponse>('/api/traffic/random', {
    method: 'POST',
  });
}

export function triggerSuccessRequest() {
  return request<TrafficRecordResponse>('/api/traffic/success', {
    method: 'POST',
  });
}

export function triggerFailedRequest() {
  return request<TrafficRecordResponse>('/api/traffic/fail', {
    method: 'POST',
  });
}

export function resetDemoState() {
  return request<ResetResponse>('/api/reset', { method: 'POST' });
}
