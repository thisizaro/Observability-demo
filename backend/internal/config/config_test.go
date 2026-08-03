package config

import "testing"

// These tests describe the contract for Load() before it's written.
// Values and defaults come from the env var table in docs/BACKEND.md.
//
// Run with: go test ./internal/config/...
// They will fail to even compile until you add a Config struct and a
// Load() function in config.go — that's expected, start there.

func TestLoad_Defaults(t *testing.T) {
	// No env vars set — every field should fall back to its documented default.
	cfg := Load()

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want default 8080", cfg.Port)
	}
	if cfg.LoadIntervalMs != 1000 {
		t.Errorf("LoadIntervalMs = %d, want default 1000", cfg.LoadIntervalMs)
	}
	if cfg.CPULoadDurationMs != 5000 {
		t.Errorf("CPULoadDurationMs = %d, want default 5000", cfg.CPULoadDurationMs)
	}
	if cfg.MemoryLoadDurationMs != 5000 {
		t.Errorf("MemoryLoadDurationMs = %d, want default 5000", cfg.MemoryLoadDurationMs)
	}
	if cfg.MemoryLoadMB != 100 {
		t.Errorf("MemoryLoadMB = %d, want default 100", cfg.MemoryLoadMB)
	}
}

func TestLoad_OverridesFromEnv(t *testing.T) {
	// t.Setenv automatically restores the previous value after this test finishes.
	t.Setenv("PORT", "9000")
	t.Setenv("LOAD_INTERVAL_MS", "250")
	t.Setenv("CPU_LOAD_DURATION_MS", "1500")
	t.Setenv("MEMORY_LOAD_DURATION_MS", "2500")
	t.Setenv("MEMORY_LOAD_MB", "50")

	cfg := Load()

	if cfg.Port != 9000 {
		t.Errorf("Port = %d, want 9000", cfg.Port)
	}
	if cfg.LoadIntervalMs != 250 {
		t.Errorf("LoadIntervalMs = %d, want 250", cfg.LoadIntervalMs)
	}
	if cfg.CPULoadDurationMs != 1500 {
		t.Errorf("CPULoadDurationMs = %d, want 1500", cfg.CPULoadDurationMs)
	}
	if cfg.MemoryLoadDurationMs != 2500 {
		t.Errorf("MemoryLoadDurationMs = %d, want 2500", cfg.MemoryLoadDurationMs)
	}
	if cfg.MemoryLoadMB != 50 {
		t.Errorf("MemoryLoadMB = %d, want 50", cfg.MemoryLoadMB)
	}
}

func TestLoad_InvalidValueFallsBackToDefault(t *testing.T) {
	// A malformed env var (not a number) should not crash Load().
	// Documented behavior: fall back to the default instead.
	t.Setenv("PORT", "not-a-number")

	cfg := Load()

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want default 8080 when PORT is invalid", cfg.Port)
	}
}
