# Frontend

## Overview

The frontend is a single-page control panel built with React and Vite. It has one job: let a user trigger backend behaviors and see the resulting state, then send them to Prometheus or Grafana to see the resulting metrics. It does not render any charts or dashboards itself — that's Grafana's job (see [ARCHITECTURE.md](ARCHITECTURE.md)).

## Component Inventory

| Component | Responsibility |
|-----------|-----------------|
| `ControlPanel` | Top-level layout; groups actions into sections |
| `LoadGenControls` | Start/Stop load generation buttons, reflects `load_generation` status |
| `TrafficControls` | Buttons for random traffic, success, and failure triggers |
| `ResourceControls` | CPU load and memory load buttons, reflects active state |
| `StatusPanel` | Polls and displays `GET /api/status`, includes a manual health-check button |
| `ExternalLinks` | "Open Prometheus" / "Open Grafana" buttons, links to configured URLs |
| `ResetButton` | Calls `POST /api/reset`, prompts for confirmation since it clears all state |

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

## Status Polling

`StatusPanel` polls `GET /api/status` on a fixed interval (2 seconds) rather than using a WebSocket or SSE connection. For a control panel with this few state fields, polling is simpler to implement and reason about, and the added latency (at most one polling interval) is imperceptible next to the 5-second Prometheus scrape interval it's ultimately reflecting (see [OBSERVABILITY.md](OBSERVABILITY.md)).

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
