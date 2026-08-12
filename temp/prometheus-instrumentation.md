# How we added Prometheus to the backend

## The goal

Make the backend expose numbers about what it's doing (how many requests, is load generation on, how much CPU/memory load is active) at a URL, `/metrics`, in a text format Prometheus knows how to read.

## The steps we actually took

1. **Installed the client library.** `go get github.com/prometheus/client_golang/...`. This is the official Go library for defining metrics and exposing them.
2. **Created `internal/metrics/metrics.go`.** One file whose only job is to define every metric and register them.
3. **Picked a type for each metric** (Counter, Gauge, or Histogram — see below).
4. **Called each metric from the right place in the HTTP handlers.** e.g. when someone hits "trigger success," increment a counter. When load generation starts, set a gauge to 1.
5. **Added a `/metrics` route** that serves the current values in Prometheus's text format.

That's the whole shape. Everything else is detail.

## The three metric types, and how to pick one

Ask yourself: **can this value go down?**

- **Counter** — no, it only ever goes up (until the app restarts). Use it for counting events: "how many requests have happened total." We used this for `http_requests_total`.
- **Gauge** — yes, it can go up or down. Use it for a current state: "is load generation on right now" (0 or 1), "how many bytes of memory load right now." We used this for `load_generation_status`, `cpu_load_active`, `memory_load_bytes`.
- **Histogram** — use this when you don't just want a count, you want to know the *spread* of values, usually durations. It buckets observations ("how many requests took under 0.25s, under 0.5s, etc.") so you can later ask "what was the p95 latency." We used this for `http_request_duration_seconds`.

If you're stuck: counters count events, gauges show current state, histograms show a distribution.

## The code shape (simplified)

```go
registry := prometheus.NewRegistry()
factory := promauto.With(registry)

requestsTotal := factory.NewCounterVec(
    prometheus.CounterOpts{Name: "http_requests_total", Help: "..."},
    []string{"result"}, // this metric is split by a label, e.g. success/failure
)

cpuActive := factory.NewGauge(
    prometheus.GaugeOpts{Name: "cpu_load_active", Help: "..."},
)
```

- `registry` just keeps track of "every metric that exists in this app." Nothing gets exposed unless it's in here.
- `promauto.With(registry)` is a shortcut so creating a metric *and* registering it happens in one call instead of two.
- `[]string{"result"}` means this isn't one number, it's a family of numbers, one per label value (`result="success"`, `result="failure"`). You get a specific one with `.WithLabelValues("success")`, then call `.Inc()` or `.Set()` on it.

Then, wherever the actual event happens in a handler:

```go
requestsTotal.WithLabelValues("success").Inc()   // counter: went up by 1
cpuActive.Set(1)                                  // gauge: now this value
```

And to expose it all:

```go
mux.Handle("GET /metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
```

`promhttp.HandlerFor` takes your registry and turns it into an `http.Handler` that, when hit, writes out every registered metric's current value in Prometheus's plain-text format.

## One gotcha worth remembering

A gauge only shows what you last `.Set()` it to — it doesn't update itself. If something changes state in the background (like a CPU load burst finishing on its own), nothing tells the gauge unless you explicitly go set it again. We solved this with one small helper function that re-reads the real state and re-sets all the gauges, called after every action *and* every status check, so the gauges don't go stale.

## How the whole pipeline works end to end

1. Something happens in the app (a request comes in, load starts).
2. The handler updates a metric (`.Inc()`, `.Set()`, `.Observe()`).
3. That update just changes a number in memory — nothing is sent anywhere yet.
4. Later, Prometheus (a separate service) calls `GET /metrics` on its own schedule (e.g. every 5 seconds).
5. Prometheus stores whatever value it saw at that moment as one data point in its own database.
6. Grafana asks Prometheus for that stored history and draws a graph.

This is why it's called a "pull" model — the app never pushes anything anywhere, it just always has current numbers ready whenever someone asks.

---

## Interview questions on this topic

**Q: What's the difference between a Counter and a Gauge?**
A: A Counter only increases (resets to 0 only on restart) and is used for counting events. A Gauge can go up or down and represents a current value/state.

**Q: Why would you use `rate()` on a counter instead of just reading its raw value?**
A: The raw value is just "total since the process started," which isn't useful on its own. `rate()` converts that into "events per second over a time window," which is what you actually want to graph.

**Q: What does a Histogram give you that a Counter doesn't?**
A: A Histogram tracks the distribution of values (like request durations), not just a total. That lets you calculate percentiles (p50, p95, p99) later, which a single running total can't give you.

**Q: Is Prometheus push-based or pull-based, and why does that matter?**
A: Pull-based — Prometheus scrapes each target's `/metrics` endpoint on a schedule. This matters because the application never needs to know where Prometheus is or send it anything; the app just needs to expose current values whenever asked. It also means metrics only update as often as the scrape interval, not in true real time.

**Q: Why keep your own `prometheus.Registry` instead of using the global default one?**
A: A dedicated registry avoids shared global state, which makes testing easier (each test can create a fresh one) and avoids "duplicate metric registration" panics if something gets initialized more than once.

**Q: If a background process stops on its own (not triggered by an API call), how does its metric ever get updated?**
A: It doesn't, automatically — nothing pushes a gauge update. You need something (a periodic sync, or the state-reading code path itself) to actively re-check the real state and call `.Set()` again.
