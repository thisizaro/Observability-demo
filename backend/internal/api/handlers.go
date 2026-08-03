package api

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"observability-demo/backend/internal/loadgen"
)

// Server holds everything the HTTP handlers need. It wraps a
// loadgen.Manager and doesn't own any state of its own beyond that —
// see docs/BACKEND.md for why state lives in loadgen, not here.
//
// TODO: fields you'll likely need:
//   - the *loadgen.Manager passed into NewServer
//   - a startTime time.Time (set in NewServer) so handleStatus can
//     compute uptime_seconds
type Server struct {
	mgr       *loadgen.Manager
	startTime time.Time
}

// NewServer constructs a Server around the given Manager.
func NewServer(mgr *loadgen.Manager) *Server {
	return &Server{
		mgr:       mgr,
		startTime: time.Now(),
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

	return mux
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
	status := s.mgr.Status()

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
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "running"})
}

// POST /api/load/stop -> 200 {"status": "idle"}
func (s *Server) handleLoadStop(w http.ResponseWriter, r *http.Request) {
	if err := s.mgr.Stop(); err != nil {
		writeError(w, http.StatusConflict, "not_running", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "idle"})
}

// POST /api/load/cpu -> 202 {"status": "cpu_load_active"}
func (s *Server) handleCPULoad(w http.ResponseWriter, r *http.Request) {
	s.mgr.TriggerCPULoad()
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "cpu_load_active"})
}

// POST /api/load/memory -> 202 {"status": "memory_load_active"}
func (s *Server) handleMemoryLoad(w http.ResponseWriter, r *http.Request) {
	s.mgr.TriggerMemoryLoad()
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "memory_load_active"})
}

// POST /api/traffic/random -> 200 {"requests_generated": <int > 0>}
// Design decision left to you: pick a small random count (e.g. 1-5)
// however you like, there's no fixed contract on the exact number.
func (s *Server) handleTrafficRandom(w http.ResponseWriter, r *http.Request) {
	count := rand.Intn(5) + 1 // 1-5 requests, see docs/API_SPEC.md — no fixed contract on the count
	writeJSON(w, http.StatusOK, map[string]int{"requests_generated": count})
}

// POST /api/traffic/success -> 200 {"status": "recorded", "result": "success"}
func (s *Server) handleTrafficSuccess(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "recorded", "result": "success"})
}

// POST /api/traffic/fail -> 200 {"status": "recorded", "result": "failure"}
func (s *Server) handleTrafficFail(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "recorded", "result": "failure"})
}

// POST /api/reset -> 200 {"status": "reset"}, and the manager's state
// (per loadgen.Manager.Reset()) goes back to idle/no load.
func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	s.mgr.Reset()
	writeJSON(w, http.StatusOK, map[string]string{"status": "reset"})
}
