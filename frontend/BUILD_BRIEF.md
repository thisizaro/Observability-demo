# Frontend Build Brief

This file is self-contained so another AI/tool can build the frontend
without needing to read the rest of the repo. Full source of truth for
anything not covered here: [`docs/FRONTEND.md`](../docs/FRONTEND.md)
and [`docs/API_SPEC.md`](../docs/API_SPEC.md).

## What this is

A single-page control panel (React + Vite + TypeScript) with buttons
that trigger a Go backend's behaviors, and a status readout that
polls the backend. It does not render charts — Prometheus/Grafana
(separate services, not built yet) handle that. This app only needs
to talk to the backend.

## Tech stack

- React + Vite + **TypeScript** (not JS — type request/response shapes
  against the API below so a contract mismatch is a compile error)
- No state management library, no CSS framework, no routing
- Plain `fetch`, no axios/react-query

## File structure to produce

```
frontend/src/
├── main.tsx              — entry point, mounts <App />
├── App.tsx               — top-level layout; action-button config array, StatusPanel, ExternalLinks
├── api.ts                — one typed function per backend endpoint below
├── useStatus.ts           — hook: polls GET /api/status every 2000ms, returns latest status + error
└── components/
    ├── ActionButton.tsx   — generic: { label, onClick, disabled? } — reused for all 8 action buttons
    ├── StatusPanel.tsx    — renders whatever useStatus returns
    └── ExternalLinks.tsx  — 2 links: "Open Prometheus", "Open Grafana"
```

Key design constraint: the 8 action buttons are structurally
identical (label + click handler), so render them from **one config
array** in `App.tsx` through the single `ActionButton` component —
don't write a separate component per action.

## Environment variables

Read via Vite's `import.meta.env`, with local-dev fallbacks in code:

| Variable | Purpose | Local default |
|----------|---------|-----------------|
| `VITE_BACKEND_URL` | Base URL for all API calls | `http://localhost:8080` |
| `VITE_PROMETHEUS_URL` | "Open Prometheus" link target | `http://localhost:9090` |
| `VITE_GRAFANA_URL` | "Open Grafana" link target | `http://localhost:3000` |

Create a `frontend/.env` with these three lines set to the local
defaults above.

## Backend API — full contract

Base URL: `VITE_BACKEND_URL`. All bodies are JSON. No auth. The
backend already has CORS enabled for all origins, so no special
fetch config is needed.

| Method | Path | Success Response |
|--------|------|-------------------|
| `GET` | `/health` | `200` `{ "status": "ok" }` |
| `GET` | `/api/status` | `200` — see shape below |
| `POST` | `/api/load/start` | `202` `{ "status": "running" }` — `409` `{ "error": "already_running", "message": "..." }` if already running |
| `POST` | `/api/load/stop` | `200` `{ "status": "idle" }` — `409` `{ "error": "not_running", "message": "..." }` if already idle |
| `POST` | `/api/load/cpu` | `202` `{ "status": "cpu_load_active" }` |
| `POST` | `/api/load/memory` | `202` `{ "status": "memory_load_active" }` |
| `POST` | `/api/traffic/random` | `200` `{ "requests_generated": <number> }` |
| `POST` | `/api/traffic/success` | `200` `{ "status": "recorded", "result": "success" }` |
| `POST` | `/api/traffic/fail` | `200` `{ "status": "recorded", "result": "failure" }` |
| `POST` | `/api/reset` | `200` `{ "status": "reset" }` |

**`GET /api/status` response shape:**

```ts
type Status = {
  load_generation: "idle" | "running"
  cpu_load_active: boolean
  memory_load_active: boolean
  uptime_seconds: number
}
```

**Error shape** (any non-2xx response):

```ts
type ApiError = {
  error: string    // machine-readable code, e.g. "already_running"
  message: string  // human-readable description
}
```

`api.ts` functions should check `response.ok`; on failure, parse the
body as `ApiError` and throw `new Error(body.message)`. On success,
parse and return the typed body.

## UI action list (button label → endpoint → notes)

1. **Start Load Generation** → `POST /api/load/start` — disable this button when `status.load_generation === "running"`
2. **Stop Load Generation** → `POST /api/load/stop` — disable when `status.load_generation !== "running"`
3. **Generate Random API Traffic** → `POST /api/traffic/random`
4. **Trigger Successful Request** → `POST /api/traffic/success`
5. **Trigger Failed Request** → `POST /api/traffic/fail`
6. **Generate CPU Load** → `POST /api/load/cpu` — disable when `status.cpu_load_active`
7. **Generate Memory Load** → `POST /api/load/memory` — disable when `status.memory_load_active`
8. **Reset Demo State** → `POST /api/reset` — wrap in `window.confirm(...)` before firing, since it clears everything

Plus two non-API actions:
9. **Open Prometheus** → `<a href={VITE_PROMETHEUS_URL} target="_blank">` (`ExternalLinks.tsx`)
10. **Open Grafana** → `<a href={VITE_GRAFANA_URL} target="_blank">` (`ExternalLinks.tsx`)

After any action succeeds, immediately re-fetch `/api/status` rather
than waiting for the next poll tick, so the UI feels responsive.

## Polling behavior

`useStatus.ts`: on mount, fetch `/api/status` immediately, then every
2000ms via `setInterval`, cleaning up on unmount. Expose
`{ status, error, refresh }` — `refresh` lets `App.tsx` force an
update right after an action.

## Non-goals (don't build these)

- No login/auth
- No embedding of Prometheus/Grafana UI — links only
- No routing — one page
- No charts/graphs in this app

---

## Ready-to-paste build prompt

> Build a React + Vite + TypeScript single-page app in `frontend/`.
> Follow the file structure, environment variables, API contract, and
> UI action list exactly as specified in `frontend/BUILD_BRIEF.md` in
> this repo. Key points: one generic `ActionButton` component driven
> by a config array (not one component per button), a `useStatus`
> hook polling `GET /api/status` every 2 seconds, typed `api.ts`
> functions for every endpoint in the brief's contract table, and a
> confirmation prompt before firing the reset action. No state
> library, no CSS framework, no routing. When done, run `npm run
> build` to confirm it compiles cleanly.
