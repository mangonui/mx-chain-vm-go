package wasmer2

import (
	"sync"
	"sync/atomic"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

// ISSUE-011: typed handle registry for executor.VMHooks references.
//
// Background. The cgo callback chain stores a uintptr in C-allocated
// memory (`vmHooksPtrStorage`) and reads it back when wasmer triggers a
// host import. The previous scheme stored `uintptr(unsafe.Pointer(&exec.vmHooks))`
// — the address of a Go interface field on a heap-allocated struct.
// That worked in practice because:
//   - Go's GC is non-moving (heap addresses are stable while reachable)
//   - `Wasmer2Executor` is held alive by VM machinery throughout
//     contract execution
// But the contract was implicit and `go vet`'s `unsafeptr` rule (correctly)
// flagged the `uintptr → unsafe.Pointer → executor.VMHooks` round-trip.
// Any future VM-lifecycle change toward short-lived executors, or a
// future Go runtime change tightening unsafe rules, would break this
// silently.
//
// This registry replaces the address-passing scheme with a stable
// `uint64` handle. The registry holds a strong reference to the
// `executor.VMHooks` interface (keeping it alive regardless of executor
// GC-eligibility) and resolves handle → hooks under an RWMutex. Stale
// handles look up to nil; the read-side helper panics with a clear
// message rather than dereferencing arbitrary memory.
//
// Concurrency. RWMutex over a map is the simplest correct primitive.
// Hook callbacks are read-heavy (Lookup); register/release happens once
// per executor lifecycle (rare). If lookup contention shows up in
// benchmarks the next step is `sync.Map` or sharded maps.

type vmHooksRegistry struct {
	mu      sync.RWMutex
	nextID  uint64
	entries map[uint64]executor.VMHooks
}

// globalVMHooksRegistry is the singleton consumed by the cgo callback
// helper. It MUST stay package-global so the callback path can resolve
// the handle without a per-call closure or context plumbing.
var globalVMHooksRegistry = &vmHooksRegistry{
	entries: make(map[uint64]executor.VMHooks),
}

// Register inserts the hooks and returns a stable handle. Handles are
// monotonically increasing and never reused; a released handle's value
// will never be reissued, so a stale-handle lookup always returns nil
// (the safest possible failure mode).
func (r *vmHooksRegistry) Register(hooks executor.VMHooks) uint64 {
	id := atomic.AddUint64(&r.nextID, 1)
	r.mu.Lock()
	r.entries[id] = hooks
	r.mu.Unlock()
	return id
}

// Lookup resolves a handle. Returns nil for unregistered or
// already-released handles. Caller-side nil handling is the security
// boundary — see lookupVMHooksOrPanic below.
func (r *vmHooksRegistry) Lookup(id uint64) executor.VMHooks {
	r.mu.RLock()
	h := r.entries[id]
	r.mu.RUnlock()
	return h
}

// Release frees a handle. Subsequent Lookup(id) returns nil. Idempotent
// (delete from a Go map on a missing key is a no-op), so concurrent
// double-release is safe.
func (r *vmHooksRegistry) Release(id uint64) {
	r.mu.Lock()
	delete(r.entries, id)
	r.mu.Unlock()
}

// lookupVMHooksOrPanic resolves a handle and panics with a clear message
// on miss. Used by getVMHooksFromContextRawPtr — every wasmer host-import
// hook calls that helper expecting a non-nil interface, so returning nil
// would just defer the panic to a cryptic nil-interface dispatch
// elsewhere. Failing here with a clear message preserves debuggability
// while still terminating the contract execution cleanly (the panic
// propagates up through the cgo callback into the wasmer execution
// frame, which is the documented path for host-import errors).
//
// In practice this panic should not fire unless:
//   - The Wasmer2Executor was destroyed (handle released) but the
//     wasmer instance still has a stale handle in its context — that's
//     a lifecycle bug elsewhere.
//   - The C-side data slot got corrupted between SetContextData and the
//     callback — unrelated memory-safety bug.
// Either way, panic-with-message is the right termination.
func lookupVMHooksOrPanic(handle uint64) executor.VMHooks {
	hooks := globalVMHooksRegistry.Lookup(handle)
	if hooks == nil {
		panic(vmHooksRegistryMissPanicMessage(handle))
	}
	return hooks
}

func vmHooksRegistryMissPanicMessage(handle uint64) string {
	// Inline because errors.New / fmt.Errorf would allocate on a hot
	// path even for the success branch (escape analysis); a constant
	// string + the handle keeps the success path allocation-free.
	return "wasmer2: vmHooks handle not found in registry; executor lifecycle bug or stale wasmer context (handle=" +
		uint64ToDecimalString(handle) + ")"
}

// uint64ToDecimalString avoids strconv import bloat on what is otherwise
// a tight package. The helper is allocation-free for the success path
// (it's only called inside the panic branch).
func uint64ToDecimalString(v uint64) string {
	if v == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
