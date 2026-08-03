package loadgen

import (
	"testing"
	"time"

	"observability-demo/backend/internal/config"
)

// These tests describe the load-generator state machine from docs/BACKEND.md
// before it's written. They will fail to compile until loadgen.go defines:
//   - State, StateIdle, StateRunning
//   - Status{ LoadGeneration State, CPULoadActive bool, MemoryLoadBytes int }
//   - Manager, NewManager(cfg config.Config) *Manager
//   - (*Manager) Start() error, Stop() error, Reset(), Status() Status
//   - (*Manager) TriggerCPULoad(), TriggerMemoryLoad()
//   - ErrAlreadyRunning, ErrNotRunning sentinel errors
//
// Durations in testConfig are kept short (20ms) so the CPU/memory load
// tests don't slow down the suite.
func testConfig() config.Config {
	return config.Config{
		Port:                 8080,
		LoadIntervalMs:       1000,
		CPULoadDurationMs:    20,
		MemoryLoadDurationMs: 20,
		MemoryLoadMB:         1,
	}
}

func TestNewManager_StartsIdle(t *testing.T) {
	m := NewManager(testConfig())
	status := m.Status()

	if status.LoadGeneration != StateIdle {
		t.Errorf("LoadGeneration = %q, want %q", status.LoadGeneration, StateIdle)
	}
	if status.CPULoadActive {
		t.Error("CPULoadActive = true, want false on a fresh Manager")
	}
	if status.MemoryLoadBytes != 0 {
		t.Errorf("MemoryLoadBytes = %d, want 0 on a fresh Manager", status.MemoryLoadBytes)
	}
}

func TestStart_TransitionsToRunning(t *testing.T) {
	m := NewManager(testConfig())

	if err := m.Start(); err != nil {
		t.Fatalf("Start() returned error: %v", err)
	}
	if m.Status().LoadGeneration != StateRunning {
		t.Errorf("LoadGeneration = %q, want %q after Start()", m.Status().LoadGeneration, StateRunning)
	}

	m.Stop() // cleanup so the background goroutine doesn't outlive the test
}

func TestStart_ErrorsWhenAlreadyRunning(t *testing.T) {
	m := NewManager(testConfig())
	if err := m.Start(); err != nil {
		t.Fatalf("first Start() returned error: %v", err)
	}

	if err := m.Start(); err != ErrAlreadyRunning {
		t.Errorf("second Start() error = %v, want ErrAlreadyRunning", err)
	}

	m.Stop()
}

func TestStop_TransitionsToIdle(t *testing.T) {
	m := NewManager(testConfig())
	m.Start()

	if err := m.Stop(); err != nil {
		t.Fatalf("Stop() returned error: %v", err)
	}
	if m.Status().LoadGeneration != StateIdle {
		t.Errorf("LoadGeneration = %q, want %q after Stop()", m.Status().LoadGeneration, StateIdle)
	}
}

func TestStop_ErrorsWhenAlreadyIdle(t *testing.T) {
	m := NewManager(testConfig())

	if err := m.Stop(); err != ErrNotRunning {
		t.Errorf("Stop() on idle manager error = %v, want ErrNotRunning", err)
	}
}

func TestTriggerCPULoad_SetsActiveThenClears(t *testing.T) {
	m := NewManager(testConfig())

	m.TriggerCPULoad()
	if !m.Status().CPULoadActive {
		t.Fatal("CPULoadActive = false immediately after TriggerCPULoad(), want true")
	}

	time.Sleep(50 * time.Millisecond) // longer than testConfig's 20ms CPU duration
	if m.Status().CPULoadActive {
		t.Error("CPULoadActive = true after duration elapsed, want false")
	}
}

func TestTriggerMemoryLoad_SetsBytesThenClears(t *testing.T) {
	m := NewManager(testConfig())
	wantBytes := 1 * 1024 * 1024 // MemoryLoadMB: 1 in testConfig

	m.TriggerMemoryLoad()
	if got := m.Status().MemoryLoadBytes; got != wantBytes {
		t.Fatalf("MemoryLoadBytes = %d immediately after trigger, want %d", got, wantBytes)
	}

	time.Sleep(50 * time.Millisecond) // longer than testConfig's 20ms memory duration
	if got := m.Status().MemoryLoadBytes; got != 0 {
		t.Errorf("MemoryLoadBytes = %d after duration elapsed, want 0", got)
	}
}

func TestReset_ClearsAllState(t *testing.T) {
	m := NewManager(testConfig())
	m.Start()
	m.TriggerCPULoad()

	m.Reset()

	status := m.Status()
	if status.LoadGeneration != StateIdle {
		t.Errorf("LoadGeneration = %q after Reset(), want %q", status.LoadGeneration, StateIdle)
	}
	if status.CPULoadActive {
		t.Error("CPULoadActive = true after Reset(), want false")
	}
	if status.MemoryLoadBytes != 0 {
		t.Errorf("MemoryLoadBytes = %d after Reset(), want 0", status.MemoryLoadBytes)
	}
}
