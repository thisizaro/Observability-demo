import { Activity, Cpu, MemoryStick, Timer } from 'lucide-react';
import type { Status } from '../api';

type StatusPanelProps = {
  status: Status | null;
  error: string | null;
  loading: boolean;
};

function formatUptime(seconds: number): string {
  const s = Math.max(0, Math.floor(seconds));
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  const sec = s % 60;
  return `${h}h ${m}m ${sec}s`;
}

export default function StatusPanel({ status, error, loading }: StatusPanelProps) {
  if (loading && !status) {
    return (
      <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
        <p className="text-sm text-slate-500">Loading status…</p>
      </div>
    );
  }

  if (error && !status) {
    return (
      <div className="rounded-xl border border-rose-200 bg-rose-50 p-6 shadow-sm">
        <p className="text-sm font-medium text-rose-700">
          Unable to reach backend
        </p>
        <p className="mt-1 text-sm text-rose-600">{error}</p>
      </div>
    );
  }

  if (!status) return null;

  const cards = [
    {
      icon: Activity,
      label: 'Load Generation',
      value: status.load_generation,
      active: status.load_generation === 'running',
    },
    {
      icon: Cpu,
      label: 'CPU Load',
      value: status.cpu_load_active ? 'active' : 'inactive',
      active: status.cpu_load_active,
    },
    {
      icon: MemoryStick,
      label: 'Memory Load',
      value: status.memory_load_active ? 'active' : 'inactive',
      active: status.memory_load_active,
    },
    {
      icon: Timer,
      label: 'Uptime',
      value: formatUptime(status.uptime_seconds),
      active: false,
      neutral: true,
    },
  ];

  return (
    <div className="rounded-xl border border-slate-200 bg-white shadow-sm">
      <div className="border-b border-slate-100 px-6 py-4">
        <div className="flex items-center gap-2">
          <span className="relative flex h-2.5 w-2.5">
            {error ? (
              <span className="absolute inline-flex h-full w-full rounded-full bg-rose-400" />
            ) : (
              <>
                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75" />
                <span className="relative inline-flex h-2.5 w-2.5 rounded-full bg-emerald-500" />
              </>
            )}
          </span>
          <h2 className="text-sm font-semibold text-slate-700">Backend Status</h2>
        </div>
        {error && (
          <p className="mt-1 text-xs text-rose-500">Last poll failed: {error}</p>
        )}
      </div>
      <div className="grid grid-cols-2 gap-px overflow-hidden rounded-b-xl bg-slate-100 lg:grid-cols-4">
        {cards.map((card) => {
          const Icon = card.icon;
          return (
            <div
              key={card.label}
              className="flex flex-col gap-2 bg-white px-5 py-4"
            >
              <div className="flex items-center gap-2 text-slate-500">
                <Icon className="h-4 w-4" />
                <span className="text-xs font-medium uppercase tracking-wide">
                  {card.label}
                </span>
              </div>
              <span
                className={`text-lg font-semibold capitalize ${
                  card.neutral
                    ? 'text-slate-700'
                    : card.active
                      ? 'text-emerald-600'
                      : 'text-slate-400'
                }`}
              >
                {card.value}
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
