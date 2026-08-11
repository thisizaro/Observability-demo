import { useCallback, useState } from 'react';
import {
  type LucideIcon,
  Activity,
  Square,
  Zap,
  MemoryStick,
  Shuffle,
  CheckCircle2,
  XCircle,
  RotateCcw,
} from 'lucide-react';
import ActionButton from './components/ActionButton';
import StatusPanel from './components/StatusPanel';
import ExternalLinks from './components/ExternalLinks';
import { useStatus } from './useStatus';
import {
  type Status,
  startLoad,
  stopLoad,
  generateCpuLoad,
  generateMemoryLoad,
  generateRandomTraffic,
  triggerSuccessRequest,
  triggerFailedRequest,
  resetDemoState,
} from './api';

type ActionResult = { ok: true; message: string } | { ok: false; message: string };

type ActionConfig = {
  id: string;
  label: string;
  Icon: LucideIcon;
  variant: 'default' | 'success' | 'danger';
  run: () => Promise<ActionResult>;
  isDisabled: (status: Status | null) => boolean;
  confirm?: string;
};

export default function App() {
  const { status, error, loading, refresh } = useStatus();
  const [busyId, setBusyId] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<ActionResult | null>(null);

  const runAction = useCallback(
    async (cfg: ActionConfig) => {
      if (cfg.confirm && !window.confirm(cfg.confirm)) return;
      setBusyId(cfg.id);
      setFeedback(null);
      try {
        const result = await cfg.run();
        setFeedback(result);
        if (result.ok) {
          await refresh();
        }
      } finally {
        setBusyId(null);
      }
    },
    [refresh],
  );

  const handleResult = (ok: boolean, label: string): ActionResult =>
    ok
      ? { ok: true, message: `${label} succeeded.` }
      : { ok: false, message: `${label} failed.` };

  const actions: ActionConfig[] = [
    {
      id: 'start-load',
      label: 'Start Load Generation',
      Icon: Activity,
      variant: 'success',
      isDisabled: (s) => s?.load_generation === 'running',
      run: async () => {
        try {
          await startLoad();
          return handleResult(true, 'Start Load Generation');
        } catch (e) {
          return { ok: false, message: e instanceof Error ? e.message : 'Start Load Generation failed.' };
        }
      },
    },
    {
      id: 'stop-load',
      label: 'Stop Load Generation',
      Icon: Square,
      variant: 'default',
      isDisabled: (s) => !s || s.load_generation !== 'running',
      run: async () => {
        try {
          await stopLoad();
          return handleResult(true, 'Stop Load Generation');
        } catch (e) {
          return { ok: false, message: e instanceof Error ? e.message : 'Stop Load Generation failed.' };
        }
      },
    },
    {
      id: 'cpu-load',
      label: 'Generate CPU Load',
      Icon: Zap,
      variant: 'default',
      isDisabled: (s) => Boolean(s?.cpu_load_active),
      run: async () => {
        try {
          await generateCpuLoad();
          return handleResult(true, 'Generate CPU Load');
        } catch (e) {
          return { ok: false, message: e instanceof Error ? e.message : 'Generate CPU Load failed.' };
        }
      },
    },
    {
      id: 'memory-load',
      label: 'Generate Memory Load',
      Icon: MemoryStick,
      variant: 'default',
      isDisabled: (s) => Boolean(s?.memory_load_active),
      run: async () => {
        try {
          await generateMemoryLoad();
          return handleResult(true, 'Generate Memory Load');
        } catch (e) {
          return { ok: false, message: e instanceof Error ? e.message : 'Generate Memory Load failed.' };
        }
      },
    },
    {
      id: 'random-traffic',
      label: 'Generate Random API Traffic',
      Icon: Shuffle,
      variant: 'default',
      isDisabled: () => false,
      run: async () => {
        try {
          const res = await generateRandomTraffic();
          return { ok: true, message: `Generated ${res.requests_generated} requests.` };
        } catch (e) {
          return { ok: false, message: e instanceof Error ? e.message : 'Generate Random API Traffic failed.' };
        }
      },
    },
    {
      id: 'success-request',
      label: 'Trigger Successful Request',
      Icon: CheckCircle2,
      variant: 'success',
      isDisabled: () => false,
      run: async () => {
        try {
          await triggerSuccessRequest();
          return handleResult(true, 'Trigger Successful Request');
        } catch (e) {
          return { ok: false, message: e instanceof Error ? e.message : 'Trigger Successful Request failed.' };
        }
      },
    },
    {
      id: 'fail-request',
      label: 'Trigger Failed Request',
      Icon: XCircle,
      variant: 'danger',
      isDisabled: () => false,
      run: async () => {
        try {
          await triggerFailedRequest();
          return handleResult(true, 'Trigger Failed Request');
        } catch (e) {
          return { ok: false, message: e instanceof Error ? e.message : 'Trigger Failed Request failed.' };
        }
      },
    },
    {
      id: 'reset',
      label: 'Reset Demo State',
      Icon: RotateCcw,
      variant: 'danger',
      isDisabled: () => false,
      confirm: 'Reset all demo state? This clears load generation, CPU/memory load, and traffic records.',
      run: async () => {
        try {
          await resetDemoState();
          return handleResult(true, 'Reset Demo State');
        } catch (e) {
          return { ok: false, message: e instanceof Error ? e.message : 'Reset Demo State failed.' };
        }
      },
    },
  ];

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-5xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-slate-900 text-white">
              <Activity className="h-5 w-5" />
            </div>
            <div>
              <h1 className="text-base font-semibold leading-tight">
                Observability Control Panel
              </h1>
              <p className="text-xs text-slate-500">
                Trigger backend behaviors and monitor status
              </p>
            </div>
          </div>
          <ExternalLinks />
        </div>
      </header>

      <main className="mx-auto max-w-5xl space-y-6 px-6 py-8">
        <StatusPanel status={status} error={error} loading={loading} />

        {feedback && (
          <div
            className={`flex items-center justify-between rounded-lg border px-4 py-3 text-sm ${
              feedback.ok
                ? 'border-emerald-200 bg-emerald-50 text-emerald-700'
                : 'border-rose-200 bg-rose-50 text-rose-700'
            }`}
          >
            <span>{feedback.message}</span>
            <button
              type="button"
              onClick={() => setFeedback(null)}
              className="text-xs font-medium uppercase opacity-60 hover:opacity-100"
            >
              Dismiss
            </button>
          </div>
        )}

        <section>
          <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-slate-500">
            Actions
          </h2>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
            {actions.map((cfg) => {
              return (
                <ActionButton
                  key={cfg.id}
                  label={cfg.label}
                  Icon={cfg.Icon}
                  disabled={cfg.isDisabled(status)}
                  loading={busyId === cfg.id}
                  variant={cfg.variant}
                  onClick={() => runAction(cfg)}
                />
              );
            })}
          </div>
          <p className="mt-4 text-xs text-slate-400">
            {actions.some((a) => a.id === busyId)
              ? 'Working…'
              : 'Buttons disable automatically based on live backend state.'}
          </p>
        </section>
      </main>
    </div>
  );
}
