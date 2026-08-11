import type { LucideIcon } from 'lucide-react';

type ActionButtonProps = {
  label: string;
  onClick: () => void;
  disabled?: boolean;
  loading?: boolean;
  variant?: 'default' | 'danger' | 'success';
  Icon?: LucideIcon;
};

export default function ActionButton({
  label,
  onClick,
  disabled = false,
  loading = false,
  variant = 'default',
  Icon,
}: ActionButtonProps) {
  const base =
    'inline-flex items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium transition-all duration-150 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50';
  const variants = {
    default:
      'bg-slate-800 text-white hover:bg-slate-700 focus:ring-slate-500',
    success:
      'bg-emerald-600 text-white hover:bg-emerald-500 focus:ring-emerald-500',
    danger:
      'bg-rose-600 text-white hover:bg-rose-500 focus:ring-rose-500',
  };

  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled || loading}
      className={`${base} ${variants[variant]}`}
    >
      {loading ? (
        <span className="h-3 w-3 animate-spin rounded-full border-2 border-white/40 border-t-white" />
      ) : (
        Icon && <Icon className="h-4 w-4" />
      )}
      {label}
    </button>
  );
}
