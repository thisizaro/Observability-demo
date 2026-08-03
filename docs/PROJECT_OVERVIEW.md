# Project Overview

## What This Is

The Observability Demo Platform is a small, self-contained system that shows how a modern application is monitored in production. It pairs a backend service that generates configurable traffic and load with a live Prometheus and Grafana stack, all controlled from a single web-based control panel.

The point is not the application logic itself (there isn't much) — the point is watching metrics move in real time as you trigger different conditions: success, failure, CPU load, memory load, and traffic bursts.

## Problem Statement

Most observability demos either:

- Show a static Grafana dashboard with no way to change what it's displaying, or
- Bury the interesting parts inside a large, unrelated real-world application

This project isolates the "cause and effect" loop: click a button, watch a metric change, see it appear on a dashboard. That loop is the entire teaching value of the project, so everything else is kept minimal.

## Goals

- Provide a control panel where a user can trigger specific backend behaviors on demand.
- Expose Prometheus metrics from the backend that reflect those behaviors accurately.
- Provide Grafana dashboards that make those metrics easy to read at a glance.
- Run the entire stack locally with a single Docker Compose command.
- Deploy the same stack to Render with minimal changes, and document a path to AWS/GCP later.
- Be readable and understandable by someone encountering the repo for the first time, including non-experts in Go, Prometheus, or Grafana.

## Non-Goals

- **Production-grade reliability.** State is in-memory and resets on backend restart. There's no database, no persistence, no multi-instance coordination.
- **Authentication or multi-tenancy.** The control panel is unauthenticated by design; it's a demo, not a product.
- **Load testing at scale.** The "load generation" features simulate traffic patterns for dashboard purposes, not to stress-test infrastructure.
- **Comprehensive alerting.** Alertmanager and paging integrations are out of scope; the focus is metrics and visualization, not incident response tooling.

## Target Audience

- Engineers evaluating this repo as a portfolio piece, assessing system design and observability fluency.
- Anyone learning how Prometheus and Grafana fit together with a real (if small) service, rather than reading docs in the abstract.

## Feature List (Control Panel Actions)

The frontend control panel exposes the following actions, each mapped to a backend behavior (see [API_SPEC.md](API_SPEC.md) for the exact contract):

| # | Action | Effect |
|---|--------|--------|
| 1 | Start load generation | Backend begins generating traffic on an interval |
| 2 | Stop load generation | Backend stops generating traffic |
| 3 | Generate random API traffic | Fires a single burst of randomized requests |
| 4 | Trigger successful request | Records one synthetic success (2xx) |
| 5 | Trigger failed request | Records one synthetic failure (4xx/5xx) |
| 6 | Generate CPU load | Spins up CPU-bound work for a short duration |
| 7 | Generate memory load | Allocates and holds memory for a short duration |
| 8 | Reset demo state | Clears counters and returns the backend to idle |
| 9 | Perform health check | Calls the health endpoint and displays the result |
| 10 | View current backend status | Displays whether load generation is running, idle, etc. |
| 11 | Open Prometheus | Opens the Prometheus UI in a new tab |
| 12 | Open Grafana | Opens the Grafana UI in a new tab |

## Success Criteria

- A reviewer can clone the repo, run one command, and see traffic flowing through a dashboard within a few minutes.
- Every control panel action visibly affects at least one metric within one Prometheus scrape interval.
- The `/docs` folder fully explains the system without requiring the reader to read source code first.

## Related Documents

- [ARCHITECTURE.md](ARCHITECTURE.md) — how the components fit together
- [API_SPEC.md](API_SPEC.md) — the exact contract behind the feature list above
- [ROADMAP.md](ROADMAP.md) — phased delivery plan
