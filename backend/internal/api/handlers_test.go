package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"observability-demo/backend/internal/config"
	"observability-demo/backend/internal/loadgen"
)

// These tests describe the HTTP contract from docs/API_SPEC.md before
// it's written. They will fail to compile until handlers.go defines:
//   - Server, NewServer(mgr *loadgen.Manager) *Server
//   - (*Server) Routes() http.Handler
//
// Each test drives the server the same way a real client would: build
// a request, send it through Routes(), inspect the recorded response.
// Short config durations keep the suite fast.

func newTestServer() *Server {
	cfg := config.Config{
		Port:                 8080,
		LoadIntervalMs:       1000,
		CPULoadDurationMs:    20,
		MemoryLoadDurationMs: 20,
		MemoryLoadMB:         1,
	}
	mgr := loadgen.NewManager(cfg)
	return NewServer(mgr)
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(v); err != nil {
		t.Fatalf("failed to decode JSON response body: %v", err)
	}
}

func TestHealthHandler(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body struct {
		Status string `json:"status"`
	}
	decodeJSON(t, rec, &body)
	if body.Status != "ok" {
		t.Errorf("status field = %q, want %q", body.Status, "ok")
	}
}

func TestStatusHandler_InitialState(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body struct {
		LoadGeneration   string `json:"load_generation"`
		CPULoadActive    bool   `json:"cpu_load_active"`
		MemoryLoadActive bool   `json:"memory_load_active"`
		UptimeSeconds    int    `json:"uptime_seconds"`
	}
	decodeJSON(t, rec, &body)
	if body.LoadGeneration != "idle" {
		t.Errorf("load_generation = %q, want %q", body.LoadGeneration, "idle")
	}
	if body.CPULoadActive {
		t.Error("cpu_load_active = true on a fresh server, want false")
	}
	if body.MemoryLoadActive {
		t.Error("memory_load_active = true on a fresh server, want false")
	}
	if body.UptimeSeconds < 0 {
		t.Errorf("uptime_seconds = %d, want >= 0", body.UptimeSeconds)
	}
}

func TestLoadStart_AcceptsThenRejectsDouble(t *testing.T) {
	srv := newTestServer()

	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/load/start", nil))
	if rec.Code != http.StatusAccepted {
		t.Fatalf("first start status = %d, want %d", rec.Code, http.StatusAccepted)
	}

	rec2 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec2, httptest.NewRequest(http.MethodPost, "/api/load/start", nil))
	if rec2.Code != http.StatusConflict {
		t.Fatalf("second start status = %d, want %d", rec2.Code, http.StatusConflict)
	}
	var body struct {
		Error string `json:"error"`
	}
	decodeJSON(t, rec2, &body)
	if body.Error != "already_running" {
		t.Errorf("error = %q, want %q", body.Error, "already_running")
	}
}

func TestLoadStop_AfterStart(t *testing.T) {
	srv := newTestServer()
	srv.Routes().ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/load/start", nil))

	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/load/stop", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("stop status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestCPULoadHandler(t *testing.T) {
	srv := newTestServer()
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/load/cpu", nil))

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
}

func TestMemoryLoadHandler(t *testing.T) {
	srv := newTestServer()
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/load/memory", nil))

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
}

func TestTrafficSuccessHandler(t *testing.T) {
	srv := newTestServer()
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/traffic/success", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body struct {
		Result string `json:"result"`
	}
	decodeJSON(t, rec, &body)
	if body.Result != "success" {
		t.Errorf("result = %q, want %q", body.Result, "success")
	}
}

func TestTrafficFailHandler(t *testing.T) {
	srv := newTestServer()
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/traffic/fail", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body struct {
		Result string `json:"result"`
	}
	decodeJSON(t, rec, &body)
	if body.Result != "failure" {
		t.Errorf("result = %q, want %q", body.Result, "failure")
	}
}

func TestTrafficRandomHandler(t *testing.T) {
	srv := newTestServer()
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/traffic/random", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body struct {
		RequestsGenerated int `json:"requests_generated"`
	}
	decodeJSON(t, rec, &body)
	if body.RequestsGenerated <= 0 {
		t.Errorf("requests_generated = %d, want > 0", body.RequestsGenerated)
	}
}

func TestResetHandler_ReturnsStateToIdle(t *testing.T) {
	srv := newTestServer()
	srv.Routes().ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/load/start", nil))

	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/reset", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("reset status = %d, want %d", rec.Code, http.StatusOK)
	}

	statusRec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(statusRec, httptest.NewRequest(http.MethodGet, "/api/status", nil))
	var body struct {
		LoadGeneration string `json:"load_generation"`
	}
	decodeJSON(t, statusRec, &body)
	if body.LoadGeneration != "idle" {
		t.Errorf("load_generation after reset = %q, want %q", body.LoadGeneration, "idle")
	}
}
