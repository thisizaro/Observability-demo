package api

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"observability-demo/backend/internal/loadgen"
	"observability-demo/backend/internal/metrics"
)

// Server holds everything the HTTP handlers need. It wraps a
// loadgen.Manager and doesn't own any state of its own beyond that —
// see docs/BACKEND.md for why state lives in loadgen, not here.
type Server struct {
	mgr       *loadgen.Manager
	startTime time.Time
	metrics   *metrics.Metrics
}

// NewServer constructs a Server around the given Manager. metrics.New
// is called here (not passed in) so startTime is guaranteed to be the
// single shared reference point for both the API's uptime_seconds
// field and the backend_uptime_seconds metric — no risk of the two
// drifting apart from being set separately.
func NewServer(mgr *loadgen.Manager) *Server {
	startTime := time.Now()
	return &Server{
		mgr:       mgr,
		startTime: startTime,
		metrics:   metrics.New(startTime),
	}
}

// Routes builds the request router and registers every endpoint from
// docs/API_SPEC.md. http.NewServeMux (Go 1.22+ pattern syntax, e.g.
// "POST /api/load/start") is enough here — no third-party router needed
// for this API's size.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /api/status", s.handleStatus)
	mux.HandleFunc("POST /api/load/start", s.handleLoadStart)
	mux.HandleFunc("POST /api/load/stop", s.handleLoadStop)
	mux.HandleFunc("POST /api/load/cpu", s.handleCPULoad)
	mux.HandleFunc("POST /api/load/memory", s.handleMemoryLoad)
	mux.HandleFunc("POST /api/traffic/random", s.handleTrafficRandom)
	mux.HandleFunc("POST /api/traffic/success", s.handleTrafficSuccess)
	mux.HandleFunc("POST /api/traffic/fail", s.handleTrafficFail)
	mux.HandleFunc("POST /api/reset", s.handleReset)
	mux.Handle("GET /metrics", promhttp.HandlerFor(s.metrics.Registry, promhttp.HandlerOpts{}))

	// Order matters: withLogging wraps withCORS, so it logs the final
	// status of every request, including CORS preflight OPTIONS
	// requests that withCORS answers directly. See middleware.go.
	return withLogging(withCORS(mux))
}

// syncGauges copies the Manager's current state into the three gauge
// metrics (load_generation_status, cpu_load_active, memory_load_bytes)
// and returns that state so callers needing it for a JSON response
// (handleStatus) don't have to call s.mgr.Status() a second time.
//
// Called after every state-changing action, not just from
// handleStatus: /metrics can be scraped independently of the frontend
// polling /api/status, so gauges need to be current at all times, not
// just right after a poll.
func (s *Server) syncGauges() loadgen.Status {
	status := s.mgr.Status()

	if status.LoadGeneration == loadgen.StateRunning {
		s.metrics.LoadGenerationStatus.Set(1)
	} else {
		s.metrics.LoadGenerationStatus.Set(0)
	}

	if status.CPULoadActive {
		s.metrics.CPULoadActive.Set(1)
	} else {
		s.metrics.CPULoadActive.Set(0)
	}

	s.metrics.MemoryLoadBytes.Set(float64(status.MemoryLoadBytes))

	return status
}

// recordTraffic increments http_requests_total for the given result
// and observes a synthetic latency into http_request_duration_seconds.
// The latency is fabricated (there's no real downstream call here) —
// it exists so the histogram has a realistic-looking spread for
// dashboards later, rather than every observation landing in the same
// bucket.
func (s *Server) recordTraffic(result string) {
	s.metrics.HTTPRequestsTotal.WithLabelValues(result).Inc()
	s.metrics.HTTPRequestDuration.WithLabelValues(result).Observe(syntheticLatencySeconds())
}

func syntheticLatencySeconds() float64 {
	return 0.01 + rand.Float64()*0.4 // 10ms-410ms
}

// writeJSON writes body as a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

// writeError writes the {"error": ..., "message": ...} shape from
// docs/API_SPEC.md's Error Responses section.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{
		"error":   code,
		"message": message,
	})
}

// --- Handlers below, one per docs/API_SPEC.md endpoint. ---
//
// Two helpers you'll probably want to write once and reuse across all
// of these (there's no framework here doing this for you):
//
//	func writeJSON(w http.ResponseWriter, status int, body any)
//	func writeError(w http.ResponseWriter, status int, code, message string)
//
// writeJSON should set Content-Type: application/json, write the
// status code, and json.NewEncoder(w).Encode(body). writeError should
// produce the {"error": ..., "message": ...} shape from API_SPEC.md's
// Error Responses section.

// GET /health -> 200 {"status": "ok"}
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GET /api/status -> 200 with load_generation, cpu_load_active,
// memory_load_active, uptime_seconds. Note: memory_load_active here is
// a bool — translate it from the Manager's MemoryLoadBytes (> 0 means
// active). See loadgen.Status's doc comment for why that translation
// belongs here, not in loadgen.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := s.syncGauges()

	writeJSON(w, http.StatusOK, map[string]any{
		"load_generation":    string(status.LoadGeneration),
		"cpu_load_active":    status.CPULoadActive,
		"memory_load_active": status.MemoryLoadBytes > 0,
		"uptime_seconds":     int(time.Since(s.startTime).Seconds()),
	})
}

// POST /api/load/start -> 202 {"status": "running"}
// On loadgen.ErrAlreadyRunning -> 409 {"error": "already_running", ...}
func (s *Server) handleLoadStart(w http.ResponseWriter, r *http.Request) {
	if err := s.mgr.Start(); err != nil {
		writeError(w, http.StatusConflict, "already_running", err.Error())
		return
	}
	s.syncGauges()
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "running"})
}

// POST /api/load/stop -> 200 {"status": "idle"}
func (s *Server) handleLoadStop(w http.ResponseWriter, r *http.Request) {
	if err := s.mgr.Stop(); err != nil {
		writeError(w, http.StatusConflict, "not_running", err.Error())
		return
	}
	s.syncGauges()
	writeJSON(w, http.StatusOK, map[string]string{"status": "idle"})
}

// POST /api/load/cpu -> 202 {"status": "cpu_load_active"}
func (s *Server) handleCPULoad(w http.ResponseWriter, r *http.Request) {
	s.mgr.TriggerCPULoad()
	s.syncGauges()
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "cpu_load_active"})
}

// POST /api/load/memory -> 202 {"status": "memory_load_active"}
func (s *Server) handleMemoryLoad(w http.ResponseWriter, r *http.Request) {
	s.mgr.TriggerMemoryLoad()
	s.syncGauges()
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "memory_load_active"})
}

// POST /api/traffic/random -> 200 {"requests_generated": <int > 0>}
// Design decision left to you: pick a small random count (e.g. 1-5)
// however you like, there's no fixed contract on the exact number.
func (s *Server) handleTrafficRandom(w http.ResponseWriter, r *http.Request) {
	count := rand.Intn(5) + 1 // 1-5 requests, see docs/API_SPEC.md — no fixed contract on the count
	for i := 0; i < count; i++ {
		result := "success"
		if rand.Float64() < 0.2 { // ~20% failure rate, just for a visible-but-minor error rate on dashboards
			result = "failure"
		}
		s.recordTraffic(result)
	}
	writeJSON(w, http.StatusOK, map[string]int{"requests_generated": count})
}

// POST /api/traffic/success -> 200 {"status": "recorded", "result": "success"}
func (s *Server) handleTrafficSuccess(w http.ResponseWriter, r *http.Request) {
	s.recordTraffic("success")
	writeJSON(w, http.StatusOK, map[string]string{"status": "recorded", "result": "success"})
}

// POST /api/traffic/fail -> 200 {"status": "recorded", "result": "failure"}
func (s *Server) handleTrafficFail(w http.ResponseWriter, r *http.Request) {
	s.recordTraffic("failure")
	writeJSON(w, http.StatusOK, map[string]string{"status": "recorded", "result": "failure"})
}

// POST /api/reset -> 200 {"status": "reset"}, and the manager's state
// (per loadgen.Manager.Reset()) goes back to idle/no load. Note this
// only resets loadgen's state, not the Prometheus counters
// (HTTPRequestsTotal, HTTPRequestDuration) — counters are meant to
// only ever increase for a process's lifetime; deliberately resetting
// one mid-process is a Prometheus anti-pattern that breaks rate()
// calculations, which assume a drop in value means the process
// restarted, not that someone reset it by hand.
func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	s.mgr.Reset()
	s.syncGauges()
	writeJSON(w, http.StatusOK, map[string]string{"status": "reset"})
}
