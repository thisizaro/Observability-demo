package metrics

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// These tests describe New() before it's written. See metrics.go's
// doc comments for the exact fields and prometheus/promauto calls
// expected. testutil.ToFloat64 reads a single metric's current value
// directly (no text-format parsing needed) — it's part of
// client_golang itself, not something we're adding.

func scrape(t *testing.T, m *Metrics) string {
	t.Helper()
	handler := promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec.Body.String()
}

func TestNew_HTTPRequestsTotalCountsByLabel(t *testing.T) {
	m := New(time.Now())
	m.HTTPRequestsTotal.WithLabelValues("success").Inc()
	m.HTTPRequestsTotal.WithLabelValues("success").Inc()
	m.HTTPRequestsTotal.WithLabelValues("failure").Inc()

	if got := testutil.ToFloat64(m.HTTPRequestsTotal.WithLabelValues("success")); got != 2 {
		t.Errorf("success count = %v, want 2", got)
	}
	if got := testutil.ToFloat64(m.HTTPRequestsTotal.WithLabelValues("failure")); got != 1 {
		t.Errorf("failure count = %v, want 1", got)
	}
}

func TestNew_HTTPRequestDurationRecordsObservations(t *testing.T) {
	m := New(time.Now())
	m.HTTPRequestDuration.WithLabelValues("success").Observe(0.25)

	output := scrape(t, m)
	if !strings.Contains(output, "http_request_duration_seconds_bucket") {
		t.Errorf("expected histogram buckets in output, got:\n%s", output)
	}
	if !strings.Contains(output, `http_request_duration_seconds_count{result="success"} 1`) {
		t.Errorf("expected one observation counted, got:\n%s", output)
	}
}

func TestNew_LoadGenerationStatusGauge(t *testing.T) {
	m := New(time.Now())
	m.LoadGenerationStatus.Set(1)

	if got := testutil.ToFloat64(m.LoadGenerationStatus); got != 1 {
		t.Errorf("load_generation_status = %v, want 1", got)
	}
}

func TestNew_CPULoadActiveGauge(t *testing.T) {
	m := New(time.Now())
	m.CPULoadActive.Set(1)

	if got := testutil.ToFloat64(m.CPULoadActive); got != 1 {
		t.Errorf("cpu_load_active = %v, want 1", got)
	}
}

func TestNew_MemoryLoadBytesGauge(t *testing.T) {
	m := New(time.Now())
	m.MemoryLoadBytes.Set(1048576)

	if got := testutil.ToFloat64(m.MemoryLoadBytes); got != 1048576 {
		t.Errorf("memory_load_bytes = %v, want 1048576", got)
	}
}

func TestNew_BackendUptimeSecondsIsExposed(t *testing.T) {
	start := time.Now().Add(-5 * time.Second)
	m := New(start)

	output := scrape(t, m)
	if !strings.Contains(output, "backend_uptime_seconds") {
		t.Errorf("expected backend_uptime_seconds to be exposed, got:\n%s", output)
	}
}

func TestNew_EachCallGetsAnIndependentRegistry(t *testing.T) {
	// If New() used the global default registry instead of creating
	// its own, this would panic with "duplicate metrics collector
	// registration attempted" on the second call.
	New(time.Now())
	New(time.Now())
}
