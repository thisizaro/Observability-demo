package api

import (
	"net/http"

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
}

// NewServer constructs a Server around the given Manager.
func NewServer(mgr *loadgen.Manager) *Server {
	panic("TODO: implement NewServer")
}

// Routes builds the request router and registers every endpoint from
// docs/API_SPEC.md. http.NewServeMux (Go 1.22+ pattern syntax, e.g.
// "POST /api/load/start") is enough here — no third-party router needed
// for this API's size.
func (s *Server) Routes() http.Handler {
	panic("TODO: implement Routes — register all handlers below on a mux and return it")
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
	panic("TODO: implement handleHealth")
}

// GET /api/status -> 200 with load_generation, cpu_load_active,
// memory_load_active, uptime_seconds. Note: memory_load_active here is
// a bool — translate it from the Manager's MemoryLoadBytes (> 0 means
// active). See loadgen.Status's doc comment for why that translation
// belongs here, not in loadgen.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement handleStatus")
}

// POST /api/load/start -> 202 {"status": "running"}
// On loadgen.ErrAlreadyRunning -> 409 {"error": "already_running", ...}
func (s *Server) handleLoadStart(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement handleLoadStart")
}

// POST /api/load/stop -> 200 {"status": "idle"}
func (s *Server) handleLoadStop(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement handleLoadStop")
}

// POST /api/load/cpu -> 202 {"status": "cpu_load_active"}
func (s *Server) handleCPULoad(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement handleCPULoad")
}

// POST /api/load/memory -> 202 {"status": "memory_load_active"}
func (s *Server) handleMemoryLoad(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement handleMemoryLoad")
}

// POST /api/traffic/random -> 200 {"requests_generated": <int > 0>}
// Design decision left to you: pick a small random count (e.g. 1-5)
// however you like, there's no fixed contract on the exact number.
func (s *Server) handleTrafficRandom(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement handleTrafficRandom")
}

// POST /api/traffic/success -> 200 {"status": "recorded", "result": "success"}
func (s *Server) handleTrafficSuccess(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement handleTrafficSuccess")
}

// POST /api/traffic/fail -> 200 {"status": "recorded", "result": "failure"}
func (s *Server) handleTrafficFail(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement handleTrafficFail")
}

// POST /api/reset -> 200 {"status": "reset"}, and the manager's state
// (per loadgen.Manager.Reset()) goes back to idle/no load.
func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	panic("TODO: implement handleReset")
}
