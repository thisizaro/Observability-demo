package loadgen

import (
	"errors"
	"sync"
	"time"

	"observability-demo/backend/internal/config"
)

// State represents whether continuous background load generation is
// running. See the state diagram in docs/BACKEND.md — CPU load and
// memory load are tracked separately from this (see Status below)
// since they're self-terminating bursts, not states of their own.
type State string

const (
	StateIdle    State = "idle"
	StateRunning State = "running"
)

// Sentinel errors returned by Start/Stop when the requested transition
// doesn't apply to the current state. The API layer maps these to a
// 409 Conflict response — see docs/API_SPEC.md's Error Responses section.
var (
	ErrAlreadyRunning = errors.New("load generation already running")
	ErrNotRunning      = errors.New("load generation not running")
)

// Status is a snapshot of the Manager's current state, returned by
// Status(). Field names correspond to (but aren't identical to) the
// GET /api/status response shape in docs/API_SPEC.md — the api package
// is responsible for translating between the two (e.g. MemoryLoadBytes
// here becomes a memory_load_active bool in the API response).
type Status struct {
	LoadGeneration  State
	CPULoadActive   bool
	MemoryLoadBytes int
}

// Manager owns all load-generator state in memory and is safe for
// concurrent use — HTTP handlers and background goroutines will call
// into it from different goroutines.
//
// TODO: think about what fields you need. You'll likely want:
//   - a sync.Mutex (or sync.RWMutex) guarding everything below
//   - the current State
//   - a bool for CPU load active
//   - an int for current memory load bytes
//   - the config.Config so Start()/TriggerCPULoad()/TriggerMemoryLoad()
//     know what interval/durations to use
//   - some way to stop the background goroutine Start() spawns (e.g. a
//     chan struct{} you close in Stop())
type Manager struct {
	mu              sync.Mutex
	state           State
	cpuLoadActive   bool
	memoryLoadBytes int
	cfg             config.Config
	stopCh          chan struct{}
}

// NewManager constructs a Manager in the idle state, ready to use.
func NewManager(cfg config.Config) *Manager {
	return &Manager{
		state: StateIdle,
		cfg:   cfg,
	}
}

// Start transitions the manager to Running and spawns a background
// goroutine that does... something on cfg.LoadIntervalMs — see
// docs/ARCHITECTURE.md's request lifecycle section for the general
// "background work + gauge" pattern used elsewhere in this package.
// Returns ErrAlreadyRunning if already Running.
func (m *Manager) Start() error {
	m.mu.Lock()
	if m.state == StateRunning {
		m.mu.Unlock()
		return ErrAlreadyRunning
	}
	m.state = StateRunning
	stopCh := make(chan struct{})
	m.stopCh = stopCh
	interval := time.Duration(m.cfg.LoadIntervalMs) * time.Millisecond
	m.mu.Unlock()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				// Background traffic tick — later phases (see
				// docs/OBSERVABILITY.md) will record a metric here.
			}
		}
	}()

	return nil
}

// Stop transitions the manager back to Idle and stops the background
// goroutine started by Start(). Returns ErrNotRunning if already Idle.
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state != StateRunning {
		return ErrNotRunning
	}
	close(m.stopCh)
	m.stopCh = nil
	m.state = StateIdle
	return nil
}

// Reset clears all state back to its zero value: Idle, no CPU load,
// no memory load. Unlike Stop(), Reset() never errors — it's valid to
// call from any state.
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state == StateRunning && m.stopCh != nil {
		close(m.stopCh)
		m.stopCh = nil
	}
	m.state = StateIdle
	m.cpuLoadActive = false
	m.memoryLoadBytes = 0
}

// Status returns a snapshot of the current state. Safe to call
// concurrently with Start/Stop/Trigger*.
func (m *Manager) Status() Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	return Status{
		LoadGeneration:  m.state,
		CPULoadActive:   m.cpuLoadActive,
		MemoryLoadBytes: m.memoryLoadBytes,
	}
}

// TriggerCPULoad spins up CPU-bound work in a background goroutine for
// cfg.CPULoadDurationMs, setting CPULoadActive true for that duration
// and false once it completes. Returns immediately — the caller
// (an HTTP handler) should not block on the load finishing.
func (m *Manager) TriggerCPULoad() {
	m.mu.Lock()
	m.cpuLoadActive = true
	duration := time.Duration(m.cfg.CPULoadDurationMs) * time.Millisecond
	m.mu.Unlock()

	go func() {
		deadline := time.Now().Add(duration)
		for time.Now().Before(deadline) {
			// Busy-loop to actually consume CPU rather than sleeping.
		}

		m.mu.Lock()
		m.cpuLoadActive = false
		m.mu.Unlock()
	}()
}

// TriggerMemoryLoad allocates cfg.MemoryLoadMB megabytes and holds them
// for cfg.MemoryLoadDurationMs, setting MemoryLoadBytes accordingly and
// back to 0 once released. Returns immediately, same as TriggerCPULoad.
func (m *Manager) TriggerMemoryLoad() {
	m.mu.Lock()
	sizeBytes := m.cfg.MemoryLoadMB * 1024 * 1024
	duration := time.Duration(m.cfg.MemoryLoadDurationMs) * time.Millisecond
	m.memoryLoadBytes = sizeBytes
	m.mu.Unlock()

	// Actually allocate and touch the memory so it's real resident
	// memory, not just a number stored in a variable.
	buf := make([]byte, sizeBytes)
	for i := range buf {
		buf[i] = 1
	}

	go func() {
		time.Sleep(duration)
		_ = buf // keep the allocation alive until the sleep completes

		m.mu.Lock()
		m.memoryLoadBytes = 0
		m.mu.Unlock()
	}()
}
