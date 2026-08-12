# Docker

## Overview

The full stack runs locally via a single `docker-compose.yml` at `infra/docker-compose.yml`. Four services join one Docker network: `frontend`, `backend`, `prometheus`, `grafana`. This is the reference topology that [DEPLOYMENT.md](DEPLOYMENT.md) maps onto Render.

## Services

| Service | Image/Build | Published Port | Internal Port |
|---------|--------------|------------------|-----------------|
| `frontend` | Built from `frontend/` (multi-stage: Vite build served by nginx) | `5173` | `80` |
| `backend` | Built from `backend/` | `8080` | `8080` |
| `prometheus` | `prom/prometheus` | `9090` | `9090` |
| `grafana` | `grafana/grafana` | `3000` | `3000` |

The frontend's internal port is nginx's default (`80`), not the Vite dev server's `5173` — the Docker build produces a static, already-built app; nothing is running `vite dev` inside the container.

`backend`'s port is published in local dev for convenience (direct `curl`/browser access to `/metrics` or `/health`), even though only `frontend` and `prometheus` need to reach it over the internal network in principle.

## Network

All services join a single user-defined bridge network (see the diagram in [ARCHITECTURE.md](ARCHITECTURE.md)). Docker Compose's default DNS makes each service name resolvable as a hostname from any other service — this is how Prometheus resolves `backend:8080` and Grafana resolves `prometheus:9090` without hardcoded IPs.

## Volumes

| Volume | Mounted Into | Purpose |
|--------|---------------|---------|
| `infra/prometheus/prometheus.yml` | `prometheus` (read-only) | Scrape configuration |
| `infra/grafana/provisioning/` | `grafana` (read-only) | Data source and dashboard provisioning |
| `infra/grafana/dashboards/` | `grafana` (read-only) | Dashboard JSON definitions |
| named volume `prometheus-data` | `prometheus` | Time series storage (survives container restarts, not `docker compose down -v`) |
| named volume `grafana-data` | `grafana` | Grafana's own settings/session storage |

## Environment Variable Wiring

| Service | Variable | Value |
|---------|----------|-------|
| `frontend` | `VITE_BACKEND_URL` | `http://localhost:8080` (browser-facing, not container-facing — see note below) |
| `frontend` | `VITE_PROMETHEUS_URL` | `http://localhost:9090` |
| `frontend` | `VITE_GRAFANA_URL` | `http://localhost:3000` |
| `backend` | `PORT` | `8080` |
| `prometheus` | — | configured via mounted `prometheus.yml`, not env vars |
| `grafana` | `GF_SECURITY_ADMIN_PASSWORD` | set for local dev convenience, not meant for production use |

**Note:** the frontend's env vars point at `localhost`, not Docker service names, because they're consumed by the user's browser (outside the Docker network), not by the frontend container itself.

## Local Dev Commands

```bash
# Start the full stack
docker compose -f infra/docker-compose.yml up --build

# Tail logs for one service
docker compose -f infra/docker-compose.yml logs -f backend

# Stop and remove containers (keep volumes)
docker compose -f infra/docker-compose.yml down

# Stop and remove containers and volumes (full reset, including Prometheus data)
docker compose -f infra/docker-compose.yml down -v
```

## Related Documents

- [ARCHITECTURE.md](ARCHITECTURE.md) — the network topology this Compose file implements
- [OBSERVABILITY.md](OBSERVABILITY.md) — what's inside `prometheus.yml` and the Grafana provisioning files
- [DEPLOYMENT.md](DEPLOYMENT.md) — how this topology maps onto Render
