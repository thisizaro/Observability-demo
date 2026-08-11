# Frontend

## Overview

The frontend is a single-page control panel built with React, Vite, and TypeScript. It has one job: let a user trigger backend behaviors and see the resulting state, then send them to Prometheus or Grafana to see the resulting metrics. It does not render any charts or dashboards itself — that's Grafana's job (see [ARCHITECTURE.md](ARCHITECTURE.md)).

TypeScript is used specifically so the request/response shapes in `api.ts` can be typed directly against [API_SPEC.md](API_SPEC.md) — a mismatch between what the frontend sends/expects and what the backend actually returns becomes a compile error instead of a silent runtime bug.

## File Structure

Deliberately minimal — no state management library, no CSS framework, no routing, no per-feature component groups. The 8 action buttons are structurally identical (a label and a click handler), so they're rendered from one config array through a single reusable component rather than one hand-written component per feature group.

```
frontend/src/
├── main.tsx               — entry point, mounts <App />
├── App.tsx                — top-level layout; renders the action buttons, StatusPanel, ExternalLinks
├── api.ts                 — every backend call, one typed function per API_SPEC.md endpoint
├── useStatus.ts            — one hook: polls GET /api/status, returns current state
└── components/
    ├── ActionButton.tsx    — one generic button (label, onClick, disabled); reused for all 8 actions
    ├── StatusPanel.tsx     — renders whatever useStatus returns
    └── ExternalLinks.tsx   — the 2 "open Prometheus/Grafana" links
```

| File | Responsibility |
|------|-----------------|
| `api.ts` | Single source of truth for the backend contract — every fetch call lives here, typed request/response shapes matching [API_SPEC.md](API_SPEC.md) |
| `useStatus.ts` | Polls `GET /api/status` on an interval, exposes the latest status to any component that needs it |
| `App.tsx` | Declares the action-button config array (label + which `api.ts` function to call) and lays out the page |
| `ActionButton.tsx` | Dumb, reusable: takes a label and an onClick, has no knowledge of which endpoint it's wired to |
| `StatusPanel.tsx` | Displays `load_generation`, `cpu_load_active`, `memory_load_active`, `uptime_seconds` from `useStatus` |
| `ExternalLinks.tsx` | Renders the Prometheus/Grafana links from env vars |

## Action-to-Endpoint Mapping

Every button in the control panel maps directly to one backend call, per [API_SPEC.md](API_SPEC.md):

| UI Action | Endpoint Called |
|-----------|-------------------|
| Start load generation | `POST /api/load/start` |
| Stop load generation | `POST /api/load/stop` |
| Generate random API traffic | `POST /api/traffic/random` |
| Trigger successful request | `POST /api/traffic/success` |
| Trigger failed request | `POST /api/traffic/fail` |
| Generate CPU load | `POST /api/load/cpu` |
| Generate memory load | `POST /api/load/memory` |
| Reset demo state | `POST /api/reset` |
| Perform health check | `GET /health` |
| View current backend status | `GET /api/status` (polled automatically, no button required) |
| Open Prometheus | Direct link, no backend call |
| Open Grafana | Direct link, no backend call |

All 8 backend-calling actions are rendered from a single config array in `App.tsx` (`{ label, onClick }` pairs) through `ActionButton` — there's no per-feature component for any of them. "Reset demo state" is the one exception with extra behavior: its `onClick` wraps the `api.ts` call in a `window.confirm(...)` before firing, since it clears all state — that confirmation lives inline in `App.tsx`'s config entry, not in a dedicated component.

## Status Polling

`StatusPanel` (via the `useStatus` hook) polls `GET /api/status` on a fixed interval (2 seconds) rather than using a WebSocket or SSE connection. For a control panel with this few state fields, polling is simpler to implement and reason about, and the added latency (at most one polling interval) is imperceptible next to the 5-second Prometheus scrape interval it's ultimately reflecting (see [OBSERVABILITY.md](OBSERVABILITY.md)).

## Environment Variables

| Variable | Purpose | Local Default |
|----------|---------|-----------------|
| `VITE_BACKEND_URL` | Base URL for all API calls | `http://localhost:8080` |
| `VITE_PROMETHEUS_URL` | Link target for "Open Prometheus" | `http://localhost:9090` |
| `VITE_GRAFANA_URL` | Link target for "Open Grafana" | `http://localhost:3000` |

These are read at build time by Vite and must be re-supplied per environment (local, Render) — see [DEPLOYMENT.md](DEPLOYMENT.md) for how they're set in each.

## Related Documents

- [API_SPEC.md](API_SPEC.md) — the contract this UI consumes
- [ARCHITECTURE.md](ARCHITECTURE.md) — why the frontend doesn't talk to Prometheus/Grafana directly
- [DEPLOYMENT.md](DEPLOYMENT.md) — environment variable configuration per deployment target
