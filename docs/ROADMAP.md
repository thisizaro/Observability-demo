# Roadmap

## Overview

This is a living document. The phase table gives the bird's-eye sequencing; the task checklist below it is the granular, per-phase to-do list, updated as work progresses. The two are kept in one file rather than split across `ROADMAP.md`/`TASKS.md` so they can't drift out of sync in a project this size.

## Phases

| Phase | Goal | Deliverables | Status |
|-------|------|---------------|--------|
| 1. Documentation | Fully specify the system before writing code | All docs in `/docs` finalized and internally consistent, `DECISIONS.md` seeded | Complete |
| 2. Backend | Implement the API contract and load-generator logic | Go service implementing every endpoint in [API_SPEC.md](API_SPEC.md) | Complete |
| 3. Frontend | Implement the control panel | React/TypeScript app implementing every action in [FRONTEND.md](FRONTEND.md) against the live backend | Complete |
| 4. Prometheus Integration | Instrument the backend | `/metrics` exposes every metric in [OBSERVABILITY.md](OBSERVABILITY.md) | Complete |
| 5. Grafana Dashboards | Build the dashboards | Dashboard JSON for all panels in [OBSERVABILITY.md](OBSERVABILITY.md) | Not started |
| 6. Docker | Wire the full stack together locally | `docker-compose.yml` per [DOCKER.md](DOCKER.md), one-command startup | Complete |
| 7. Deployment | Ship to Render | Live deployment per [DEPLOYMENT.md](DEPLOYMENT.md) | Not started |
| 8. Polish | Final validation and presentation | README screenshots, troubleshooting notes validated against the live stack | Not started |

Phases are sequenced so each one builds on a working, testable artifact from the last: the backend is built and testable via `curl` before the frontend exists to call it; metrics are instrumented before Prometheus needs anything to scrape; Docker wiring happens last among the local-dev phases because it depends on all four components already working independently.

## Task Checklist

### Phase 1 — Documentation
- [x] Write PROJECT_OVERVIEW.md
- [x] Write ARCHITECTURE.md
- [x] Write API_SPEC.md
- [x] Write BACKEND.md
- [x] Write OBSERVABILITY.md
- [x] Write FRONTEND.md
- [x] Write DOCKER.md
- [x] Write DEPLOYMENT.md
- [x] Write ROADMAP.md
- [x] Write DECISIONS.md
- [x] Write README.md

### Phase 2 — Backend
- [x] Scaffold Go module and package layout per [BACKEND.md](BACKEND.md)
- [x] Implement load-generator state machine (`internal/loadgen`, race-clean)
- [x] Implement all `/api` endpoints per [API_SPEC.md](API_SPEC.md)
- [x] Implement `/health`
- [x] Add unit tests for the state machine and handlers (21 tests passing)

### Phase 3 — Frontend
Built against the TypeScript + single-config-array spec (React + Vite + TypeScript + Tailwind; built externally, then verified end-to-end against the running backend).
- [x] Scaffold Vite + React + TypeScript app
- [x] Implement `api.ts` (typed fetch functions per [API_SPEC.md](API_SPEC.md))
- [x] Implement `useStatus.ts` polling hook
- [x] Implement `ActionButton.tsx` and the action-config array in `App.tsx`
- [x] Implement `StatusPanel.tsx` and `ExternalLinks.tsx`
- [x] Wire environment variables per [FRONTEND.md](FRONTEND.md)

### Phase 4 — Prometheus Integration
- [x] Add `prometheus/client_golang` dependency
- [x] Instrument all metrics from the catalog in [OBSERVABILITY.md](OBSERVABILITY.md) (`internal/metrics`, wired into every handler via `syncGauges`/`recordTraffic`)
- [x] Verify `/metrics` output manually via curl

### Phase 5 — Grafana Dashboards
Phase 6's Compose stack now provides a real, running Prometheus + Grafana to build against, so the "temporary local Prometheus" originally planned here isn't needed — use `docker compose -f infra/docker-compose.yml up` instead.
- [ ] Build Traffic Overview, Resource Load, and System Status dashboards
- [ ] Export dashboards as JSON into `infra/grafana/dashboards`

### Phase 6 — Docker
- [x] Write `infra/docker-compose.yml`
- [x] Write `infra/prometheus/prometheus.yml`
- [x] Write Grafana provisioning config (datasource auto-configured, dashboard provider watching `infra/grafana/dashboards`)
- [x] Verify full stack starts with one command — confirmed: all 4 containers up, Prometheus target `backend` shows `up`, Grafana datasource auto-provisioned, and a full trigger → scrape → PromQL-query pipeline verified end to end. Dashboards themselves are still Phase 5 (the folder is watched but empty).

### Phase 7 — Deployment
- [ ] Create Render services per [DEPLOYMENT.md](DEPLOYMENT.md)
- [ ] Configure environment variables and secrets
- [ ] Verify live deployment end-to-end

### Phase 8 — Polish
- [ ] Validate troubleshooting notes against the live deployment
- [ ] Add README screenshots/GIF
- [ ] Final consistency pass across all docs
- [ ] Add rate-limiting middleware on the backend (deferred from Phase 2; the public Render deployment has no auth per [DECISIONS.md](DECISIONS.md) ADR-006, so basic abuse protection on load-triggering endpoints is worth adding before Phase 8 wraps up)

## Related Documents

- [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md) — the goals this roadmap serves
- [DECISIONS.md](DECISIONS.md) — architectural decisions made along the way
