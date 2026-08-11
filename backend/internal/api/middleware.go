package api

import (
	"log/slog"
	"net/http"
	"time"
)

// withCORS allows the frontend (a different origin — e.g. localhost:5173
// vs the backend's localhost:8080) to call this API from a browser.
// Permissive by design: there's no auth on this API to protect (see
// docs/DECISIONS.md ADR-006), so a wildcard origin doesn't weaken
// anything that wasn't already open.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// statusRecorder wraps http.ResponseWriter to capture the status code
// a handler actually sent, so withLogging can log it afterward.
// Overriding Write (not just WriteHeader) matters: if a handler never
// calls WriteHeader explicitly and just calls Write, Go implicitly
// sends 200 — and that implicit call goes through Write, bypassing our
// WriteHeader override entirely unless Write also routes through it.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.wroteHeader {
		return
	}
	r.status = status
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

// withLogging logs one line per request via log/slog: method, path,
// status, and duration. Using slog (rather than fmt.Println or
// log.Printf) means every log call already carries structured
// key-value fields — switching the output format to JSON later (for
// shipping to Loki, Datadog, or adding OpenTelemetry trace/span IDs)
// is a one-line change where the logger is configured, not a rewrite
// of every log statement in this file.
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		slog.Info("request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
