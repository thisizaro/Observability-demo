# API Specification

## Overview

The backend exposes a small REST API under `/api`, plus two unprefixed operational endpoints (`/health`, `/metrics`). All request and response bodies are JSON. There is no authentication and no versioning — this is a single-consumer (the frontend) demo API, so a `/v1` prefix or auth headers would add ceremony without adding teaching value (see [DECISIONS.md](DECISIONS.md)).

Base URL in local development: `http://localhost:8080`.

## Conventions

- All `POST` action endpoints that kick off background work return `202 Accepted` immediately; they do not wait for the work to finish.
- All `POST` action endpoints that complete synchronously (single traffic events, reset) return `200 OK`.
- Errors use a consistent shape (see [Error Responses](#error-responses)).
- No request bodies are required for any endpoint in the current scope; all actions are parameterless.

## Endpoints

### Load Generation

| Method | Path | Description | Success Response |
|--------|------|-------------|-------------------|
| `POST` | `/api/load/start` | Start continuous background load generation (random traffic on an interval) | `202 Accepted` — `{ "status": "running" }` |
| `POST` | `/api/load/stop` | Stop continuous background load generation | `200 OK` — `{ "status": "idle" }` |
| `POST` | `/api/load/cpu` | Generate a fixed-duration burst of CPU load | `202 Accepted` — `{ "status": "cpu_load_active" }` |
| `POST` | `/api/load/memory` | Generate a fixed-duration burst of memory load | `202 Accepted` — `{ "status": "memory_load_active" }` |

### Traffic Generation

| Method | Path | Description | Success Response |
|--------|------|-------------|-------------------|
| `POST` | `/api/traffic/random` | Fire one burst of randomized requests (mix of success/failure, varied latency) | `200 OK` — `{ "requests_generated": <int> }` |
| `POST` | `/api/traffic/success` | Record one synthetic successful request | `200 OK` — `{ "status": "recorded", "result": "success" }` |
| `POST` | `/api/traffic/fail` | Record one synthetic failed request | `200 OK` — `{ "status": "recorded", "result": "failure" }` |

### State Management

| Method | Path | Description | Success Response |
|--------|------|-------------|-------------------|
| `POST` | `/api/reset` | Clear all counters and return load generation to idle | `200 OK` — `{ "status": "reset" }` |
| `GET` | `/api/status` | Return current backend state | `200 OK` — see below |

**`GET /api/status` response body:**

```json
{
  "load_generation": "idle",
  "cpu_load_active": false,
  "memory_load_active": false,
  "uptime_seconds": 1234
}
```

`load_generation` is one of: `"idle"`, `"running"`.

### Operational Endpoints

| Method | Path | Description | Success Response |
|--------|------|-------------|-------------------|
| `GET` | `/health` | Liveness/readiness check | `200 OK` — `{ "status": "ok" }` |
| `GET` | `/metrics` | Prometheus exposition format | `200 OK` — `text/plain` metrics body (see [OBSERVABILITY.md](OBSERVABILITY.md)) |

## Error Responses

Any endpoint that fails returns a non-2xx status with this body shape:

```json
{
  "error": "short machine-readable code",
  "message": "human-readable description"
}
```

| Status | Meaning | Example `error` |
|--------|---------|------------------|
| `400` | Malformed request | `"invalid_request"` |
| `409` | Action conflicts with current state (e.g. starting load that's already running) | `"already_running"` |
| `500` | Unexpected backend failure | `"internal_error"` |

## Mapping to Control Panel Actions

This table cross-references [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md)'s feature list to the endpoints above:

| Control Panel Action | Endpoint |
|-----------------------|----------|
| Start load generation | `POST /api/load/start` |
| Stop load generation | `POST /api/load/stop` |
| Generate random API traffic | `POST /api/traffic/random` |
| Trigger successful request | `POST /api/traffic/success` |
| Trigger failed request | `POST /api/traffic/fail` |
| Generate CPU load | `POST /api/load/cpu` |
| Generate memory load | `POST /api/load/memory` |
| Reset demo state | `POST /api/reset` |
| Perform health check | `GET /health` |
| View current backend status | `GET /api/status` |

("Open Prometheus" and "Open Grafana" are frontend-only actions — direct links, not backend calls.)

## Related Documents

- [BACKEND.md](BACKEND.md) — implementation of this contract
- [FRONTEND.md](FRONTEND.md) — how the UI consumes this contract
- [OBSERVABILITY.md](OBSERVABILITY.md) — what `/metrics` exposes
