package config

import (
	"os"
	"strconv"
)

// Config holds all backend configuration, loaded from environment
// variables. See the env var table in docs/BACKEND.md for the full list
// of variables and their defaults.
type Config struct {
	Port                 int
	LoadIntervalMs       int
	CPULoadDurationMs    int
	MemoryLoadDurationMs int
	MemoryLoadMB         int
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

// Load reads configuration from environment variables, falling back to
// the documented default for any variable that is unset OR fails to
// parse as an integer (see config_test.go's TestLoad_InvalidValueFallsBackToDefault
// for that exact behavior).
//
// TODO: implement this using os.Getenv + strconv.Atoi for each field.
// Consider a small local helper like:
//
//	func envInt(key string, fallback int) int { ... }
//
// so you don't repeat the same get-parse-fallback logic five times.
func Load() Config {
	// panic("TODO: implement Load()")
	return Config{
		Port:                 envInt("PORT", 8080),
		LoadIntervalMs:       envInt("LOAD_INTERVAL_MS", 1000),
		CPULoadDurationMs:    envInt("CPU_LOAD_DURATION_MS", 5000),
		MemoryLoadDurationMs: envInt("MEMORY_LOAD_DURATION_MS", 5000),
		MemoryLoadMB:         envInt("MEMORY_LOAD_MB", 100),
	}
}
