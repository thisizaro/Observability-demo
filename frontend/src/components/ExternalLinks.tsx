import { ExternalLink, BarChart3, LineChart } from 'lucide-react';

const PROMETHEUS_URL =
  import.meta.env.VITE_PROMETHEUS_URL ?? 'http://localhost:9090';
const GRAFANA_URL =
  import.meta.env.VITE_GRAFANA_URL ?? 'http://localhost:3000';

const links = [
  {
    label: 'Open Prometheus',
    href: PROMETHEUS_URL,
    Icon: BarChart3,
  },
  {
    label: 'Open Grafana',
    href: GRAFANA_URL,
    Icon: LineChart,
  },
];

export default function ExternalLinks() {
  return (
    <div className="flex flex-col gap-2 sm:flex-row">
      {links.map(({ label, href, Icon }) => (
        <a
          key={label}
          href={href}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center justify-center gap-2 rounded-lg border border-slate-200 bg-white px-4 py-2.5 text-sm font-medium text-slate-700 shadow-sm transition-all duration-150 hover:border-slate-300 hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2"
        >
          <Icon className="h-4 w-4 text-slate-500" />
          {label}
          <ExternalLink className="h-3.5 w-3.5 text-slate-400" />
        </a>
      ))}
    </div>
  );
}
