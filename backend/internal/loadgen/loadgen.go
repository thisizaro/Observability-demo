package loadgen

import (
	"errors"

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
}

// NewManager constructs a Manager in the idle state, ready to use.
func NewManager(cfg config.Config) *Manager {
	panic("TODO: implement NewManager")
}

// Start transitions the manager to Running and spawns a background
// goroutine that does... something on cfg.LoadIntervalMs — see
// docs/ARCHITECTURE.md's request lifecycle section for the general
// "background work + gauge" pattern used elsewhere in this package.
// Returns ErrAlreadyRunning if already Running.
func (m *Manager) Start() error {
	panic("TODO: implement Start")
}

// Stop transitions the manager back to Idle and stops the background
// goroutine started by Start(). Returns ErrNotRunning if already Idle.
func (m *Manager) Stop() error {
	panic("TODO: implement Stop")
}

// Reset clears all state back to its zero value: Idle, no CPU load,
// no memory load. Unlike Stop(), Reset() never errors — it's valid to
// call from any state.
func (m *Manager) Reset() {
	panic("TODO: implement Reset")
}

// Status returns a snapshot of the current state. Safe to call
// concurrently with Start/Stop/Trigger*.
func (m *Manager) Status() Status {
	panic("TODO: implement Status")
}

// TriggerCPULoad spins up CPU-bound work in a background goroutine for
// cfg.CPULoadDurationMs, setting CPULoadActive true for that duration
// and false once it completes. Returns immediately — the caller
// (an HTTP handler) should not block on the load finishing.
func (m *Manager) TriggerCPULoad() {
	panic("TODO: implement TriggerCPULoad")
}

// TriggerMemoryLoad allocates cfg.MemoryLoadMB megabytes and holds them
// for cfg.MemoryLoadDurationMs, setting MemoryLoadBytes accordingly and
// back to 0 once released. Returns immediately, same as TriggerCPULoad.
func (m *Manager) TriggerMemoryLoad() {
	panic("TODO: implement TriggerMemoryLoad")
}
