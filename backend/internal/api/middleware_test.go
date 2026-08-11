package api

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// captureLogs redirects slog's default logger to a buffer for the
// duration of one test, restoring the original afterward. withLogging
// (in middleware.go) logs via the slog package-level functions, so
// tests observe its output by swapping the default handler rather than
// injecting a logger instance.
func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(original) })
	return &buf
}

// These tests describe withLogging before it's written. It should wrap
// a handler and, after the handler runs, log one line containing at
// least the request method, path, and the response's actual status
// code (including the 200 Go implies when a handler never calls
// WriteHeader explicitly — see the Write()-without-WriteHeader case
// below, which is the part worth getting right).

func TestWithLogging_LogsMethodPathAndExplicitStatus(t *testing.T) {
	buf := captureLogs(t)

	handler := withLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/load/cpu", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	output := buf.String()
	for _, want := range []string{"POST", "/api/load/cpu", "202"} {
		if !strings.Contains(output, want) {
			t.Errorf("log output missing %q, got: %q", want, output)
		}
	}
}

func TestWithLogging_DefaultsTo200WhenHandlerOnlyWrites(t *testing.T) {
	buf := captureLogs(t)

	// This handler never calls WriteHeader — Go implicitly sends 200
	// on the first Write(). The logging middleware needs to observe
	// that implicit status, not just ones set explicitly.
	handler := withLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	output := buf.String()
	for _, want := range []string{"GET", "/health", "200"} {
		if !strings.Contains(output, want) {
			t.Errorf("log output missing %q, got: %q", want, output)
		}
	}
}
