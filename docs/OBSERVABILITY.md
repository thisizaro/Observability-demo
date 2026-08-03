# Observability: Metrics, Prometheus, and Grafana

## Overview

This document covers the full observability pipeline in one place: what the backend measures, how Prometheus collects it, and how Grafana displays it. These three concerns are merged into a single doc rather than split across `METRICS.md`, `PROMETHEUS.md`, and `GRAFANA.md` because they describe one continuous pipeline — instrument, scrape, visualize — and keeping them together makes it easier to trace a metric from its definition to its dashboard panel.

## Metric Catalog

| Name | Type | Labels | Description |
|------|------|--------|--------------|
| `http_requests_total` | Counter | `result` (`success`/`failure`) | Total synthetic requests recorded, by outcome |
| `http_request_duration_seconds` | Histogram | `result` | Distribution of synthetic request latencies |
| `load_generation_status` | Gauge | — | `1` if background load generation is running, `0` if idle |
| `cpu_load_active` | Gauge | — | `1` while a CPU load burst is running, `0` otherwise |
| `memory_load_bytes` | Gauge | — | Bytes currently held by an active memory load burst, `0` when idle |
| `backend_uptime_seconds` | Counter | — | Seconds since the backend process started |

### Naming Conventions

Metric names follow Prometheus's own conventions: `_total` suffix for counters, base units (`seconds`, `bytes`) rather than derived ones (`ms`, `MB`), and a noun-first name that reads naturally in a sentence ("http requests total"). Labels are kept low-cardinality on purpose — `result` has exactly two values, and no label is derived from unbounded input (like a raw path or user ID), since unbounded label values are the most common cause of Prometheus cardinality blowups.

## Prometheus Configuration

Prometheus scrapes the backend's `/metrics` endpoint on a fixed interval using static configuration (no service discovery needed for a four-service demo):

```yaml
scrape_configs:
  - job_name: 'backend'
    scrape_interval: 5s
    static_configs:
      - targets: ['backend:8080']
```

- **Scrape interval:** 5 seconds — short enough that control panel actions show up on dashboards within a few seconds, without generating meaningful load for a demo-scale service.
- **Target discovery:** static config using the Docker Compose service name `backend` as the DNS hostname (see [DOCKER.md](DOCKER.md) for the network that makes this resolvable).
- **Retention:** default local retention is sufficient; this is a demo, not a long-term monitoring deployment.

### Example PromQL Queries

| Question | Query |
|----------|-------|
| Request rate by outcome | `rate(http_requests_total[1m])` |
| Error ratio | `rate(http_requests_total{result="failure"}[1m]) / rate(http_requests_total[1m])` |
| p95 latency | `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))` |
| Is load generation running | `load_generation_status` |
| Current memory load | `memory_load_bytes` |

## Grafana Configuration

Grafana's Prometheus data source and dashboards are provisioned as code (files checked into `infra/grafana/provisioning` and `infra/grafana/dashboards`) rather than configured by hand through the UI, so the dashboards are reproducible from a fresh `docker compose up`.

### Dashboard Inventory

| Dashboard | Panels |
|-----------|--------|
| **Traffic Overview** | Request rate by outcome (line graph), error ratio (line graph), p95/p50 latency (line graph) |
| **Resource Load** | CPU load active (state timeline), memory load bytes (line graph) |
| **System Status** | Load generation status (stat panel), backend uptime (stat panel) |

Each panel's query is one of the PromQL examples above or a close variant.

## Opening Prometheus and Grafana

The frontend control panel does not embed either UI — it provides direct links (`http://localhost:9090` for Prometheus, `http://localhost:3000` for Grafana in local dev) that open in a new tab. See [FRONTEND.md](FRONTEND.md) for how these URLs are configured per environment.

## Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| Grafana panels show "No Data" | Prometheus data source misconfigured, or Prometheus itself has no data | Check the data source URL in Grafana settings; confirm `http://prometheus:9090/graph` shows data for the query directly |
| Prometheus target shows "DOWN" | Backend isn't reachable on the Docker network, or `/metrics` isn't responding | Check `docker compose logs backend`; confirm `curl http://backend:8080/metrics` from another container on the network |
| Metrics don't update after an action | Waiting on the next scrape interval (up to 5s), or the action didn't actually update the gauge/counter | Wait one scrape interval; check backend logs for the action's handler |

## Related Documents

- [BACKEND.md](BACKEND.md) — where these metrics are instrumented in code
- [DOCKER.md](DOCKER.md) — the network that makes `backend`, `prometheus`, `grafana` resolvable hostnames
- [ARCHITECTURE.md](ARCHITECTURE.md) — how this pipeline fits into the whole system
