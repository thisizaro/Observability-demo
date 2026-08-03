# Decisions (ADR Log)

## Overview

This is a running log of significant architectural decisions, newest first. Each entry captures the context at the time, the decision made, and its consequences, so the reasoning survives even after the decision itself feels obvious in hindsight.

---

## ADR-008: Merge METRICS, PROMETHEUS, and GRAFANA into a single OBSERVABILITY.md

**Date:** 2026-08-03
**Status:** Accepted

**Context:** The initial documentation proposal had three separate documents for metrics, Prometheus configuration, and Grafana dashboards. For a project this size, three files describing one continuous pipeline (instrument → scrape → visualize) added navigation overhead without adding clarity.

**Decision:** Combine them into a single `OBSERVABILITY.md` covering the metric catalog, Prometheus scrape config, Grafana dashboards, and troubleshooting in one place.

**Consequences:** Fewer files to keep in sync; a metric's full lifecycle (defined → scraped → charted) is readable in one pass. Trade-off: the single file is longer than any one of the three would have been individually.

---

## ADR-007: Single DECISIONS.md log instead of one file per ADR

**Date:** 2026-08-03
**Status:** Accepted

**Context:** Standard ADR practice often uses one file per decision (`0001-use-go.md`, `0002-use-react.md`, etc.), which scales well for large teams and long-lived projects with many contributors.

**Decision:** Use a single `DECISIONS.md` with reverse-chronological entries instead.

**Consequences:** Simpler to browse for a project of this size and lifespan; a reader gets the full decision history in one scroll. Trade-off: harder to deep-link to a single decision from elsewhere, and the file will grow long over time — acceptable given the project's scope.

---

## ADR-006: No authentication on the control panel or API

**Date:** 2026-08-03
**Status:** Accepted

**Context:** A real control panel that can trigger CPU/memory load would need access control. This is a demo meant to be run locally or on a disposable Render deployment, not exposed as a shared production tool.

**Decision:** Ship with no authentication on any endpoint.

**Consequences:** Simpler implementation, no session/credential management to document or build. Trade-off: the deployed Render instance should not be treated as a public-facing service without adding auth first — this is called out explicitly in [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md)'s non-goals.

---

## ADR-005: Render for initial deployment

**Date:** 2026-08-03
**Status:** Accepted

**Context:** The project needs a deployment target that's fast to set up, supports both static sites and Docker-based services, and doesn't require managing cloud infrastructure directly for a demo project.

**Decision:** Deploy to Render first; treat AWS/GCP as a documented future path rather than building for them now (see [DEPLOYMENT.md](DEPLOYMENT.md)).

**Consequences:** Fast path to a live, shareable demo. Trade-off: some Render-specific configuration (service types, env var management) will need translation work if/when moving to AWS or GCP.

---

## ADR-004: Grafana for visualization

**Date:** 2026-08-03
**Status:** Accepted

**Context:** Needed a dashboarding tool that pairs naturally with Prometheus and is widely recognized in the industry, since portfolio-readability is a goal.

**Decision:** Use Grafana, with dashboards provisioned as code rather than built by hand through the UI.

**Consequences:** Dashboards are reproducible from a fresh deploy and reviewable as diffs in the repo. Trade-off: hand-authoring dashboard JSON is more tedious than clicking through the UI, but worth it for reproducibility.

---

## ADR-003: Prometheus for metrics collection

**Date:** 2026-08-03
**Status:** Accepted

**Context:** Needed a metrics system that's the de facto standard for this kind of demo, with a simple pull-based model that's easy to reason about and explain.

**Decision:** Use Prometheus with a pull-based scrape model against the backend's `/metrics` endpoint.

**Consequences:** No push infrastructure or message queue needed; Prometheus's own service discovery (static config, in this case) handles target management. Trade-off: pull-based collection requires the backend to always be reachable from Prometheus, which is a non-issue at this scale but would need revisiting for ephemeral/serverless backends.

---

## ADR-002: Docker Compose instead of Kubernetes for local development

**Date:** 2026-08-03
**Status:** Accepted

**Context:** The project needs to run all four services together locally with minimal setup friction for anyone cloning the repo.

**Decision:** Use Docker Compose rather than Kubernetes (e.g. via Kind or Minikube) for local orchestration.

**Consequences:** One command (`docker compose up`) starts the whole stack; no cluster tooling required. Trade-off: Compose doesn't model production-grade orchestration concerns (scaling, rolling deploys) — acceptable since those aren't part of what this demo teaches.

---

## ADR-001: React (with Vite) instead of Next.js for the frontend

**Date:** 2026-08-03
**Status:** Accepted

**Context:** The frontend is a single control panel with no server-rendered pages, no routing complexity, and no SEO requirements.

**Decision:** Use React with Vite rather than Next.js.

**Consequences:** Faster local dev iteration (Vite's dev server), simpler deployment as a static site (see [DEPLOYMENT.md](DEPLOYMENT.md)). Trade-off: gives up Next.js's server-side rendering and API routes, neither of which this project needs.

---

## ADR-000: Go instead of Node for the backend

**Date:** 2026-08-03
**Status:** Accepted

**Context:** The backend needs to run background load-generation work (CPU bursts, memory bursts, interval-based traffic) concurrently with serving HTTP requests, and expose Prometheus metrics.

**Decision:** Use Go, leveraging goroutines for background work and the official `prometheus/client_golang` library for instrumentation.

**Consequences:** Simple concurrency model for background bursts without callback/async complexity; first-class Prometheus tooling. Trade-off: fewer contributors may be immediately familiar with Go versus Node, mitigated by keeping the backend's scope small and well-documented in [BACKEND.md](BACKEND.md). The overall architecture in [ARCHITECTURE.md](ARCHITECTURE.md) does not depend on this choice — a Node or other backend could implement the same [API_SPEC.md](API_SPEC.md) contract.

## Related Documents

- [ARCHITECTURE.md](ARCHITECTURE.md) — the system these decisions shape
- [ROADMAP.md](ROADMAP.md) — where these decisions get implemented
