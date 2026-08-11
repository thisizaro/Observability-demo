package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics holds every metric from docs/OBSERVABILITY.md's catalog,
// plus the dedicated Registry they're all registered to. Using a
// per-instance Registry (rather than prometheus.DefaultRegisterer,
// the package-level global) means multiple Metrics instances — e.g.
// one per test — never collide with a "duplicate metrics collector
// registration attempted" panic.
type Metrics struct {
	Registry *prometheus.Registry

	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	LoadGenerationStatus prometheus.Gauge
	CPULoadActive        prometheus.Gauge
	MemoryLoadBytes      prometheus.Gauge
	// backend_uptime_seconds has no field here: it's a CounterFunc
	// that computes time.Since(startTime) at scrape time, so nothing
	// ever needs to call a method on it directly — see New() below.
}

// New constructs a Metrics with its own Registry and registers every
// metric from the catalog. startTime is used to compute
// backend_uptime_seconds dynamically at scrape time.
//
// TODO: for each metric, use promauto.With(registry).New___(...) —
// this creates AND registers in one call, instead of two separate
// steps. You'll need something like:
//
//	registry := prometheus.NewRegistry()
//	factory := promauto.With(registry)
//
//	httpRequestsTotal := factory.NewCounterVec(
//		prometheus.CounterOpts{Name: "http_requests_total", Help: "..."},
//		[]string{"result"},
//	)
//
//	httpRequestDuration := factory.NewHistogramVec(
//		prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "..."},
//		[]string{"result"},
//	)
//
//	loadGenerationStatus := factory.NewGauge(
//		prometheus.GaugeOpts{Name: "load_generation_status", Help: "..."},
//	)
//	// same pattern for CPULoadActive and MemoryLoadBytes
//
//	factory.NewCounterFunc(
//		prometheus.CounterOpts{Name: "backend_uptime_seconds", Help: "..."},
//		func() float64 { return time.Since(startTime).Seconds() },
//	)
//	// no variable needed for this one — nothing calls methods on it
//
// Then build and return the *Metrics struct with the Registry and the
// five named fields above.
func New(startTime time.Time) *Metrics {
	panic("TODO: implement New")
}
