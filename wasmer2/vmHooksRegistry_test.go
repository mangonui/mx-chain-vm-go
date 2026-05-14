package wasmer2

import (
	"sync"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	executorwrapper "github.com/multiversx/mx-chain-vm-go/executor/wrapper"
)

// newStub returns a non-nil executor.VMHooks interface value backed by
// a typed-nil *WrapperVMHooks pointer. This lets the registry tests
// round-trip a real interface value through Register / Lookup without
// having to implement every method of the (very large) VMHooks
// interface — none of the tests invoke methods on the returned
// interface, they only verify registry mechanics.
//
// Note the `_ int` parameter is intentionally unused; it documents
// the test's intent that each call site is conceptually a distinct
// hooks instance even though the underlying typed-nil values compare
// equal. The registry doesn't depend on identity, only on the handle.
func newStub(_ int) executor.VMHooks {
	var w *executorwrapper.WrapperVMHooks
	return w
}

func TestVMHooksRegistry_RegisterLookup(t *testing.T) {
	t.Parallel()

	r := &vmHooksRegistry{entries: make(map[uint64]executor.VMHooks)}
	h := newStub(7)
	id := r.Register(h)
	if id == 0 {
		t.Fatal("expected non-zero handle")
	}
	got := r.Lookup(id)
	if got == nil {
		t.Fatal("Lookup returned nil for a freshly-registered handle")
	}
}

func TestVMHooksRegistry_ReleaseStaleHandleReturnsNil(t *testing.T) {
	t.Parallel()

	r := &vmHooksRegistry{entries: make(map[uint64]executor.VMHooks)}
	id := r.Register(newStub(1))
	r.Release(id)
	if got := r.Lookup(id); got != nil {
		t.Fatalf("Lookup of released handle returned %v, want nil", got)
	}
}

func TestVMHooksRegistry_DoubleReleaseIsNoOp(t *testing.T) {
	t.Parallel()

	r := &vmHooksRegistry{entries: make(map[uint64]executor.VMHooks)}
	id := r.Register(newStub(1))
	r.Release(id)
	// Must not panic; map.delete on a missing key is a no-op.
	r.Release(id)
	r.Release(id)
}

func TestVMHooksRegistry_HandlesNeverReused(t *testing.T) {
	t.Parallel()

	r := &vmHooksRegistry{entries: make(map[uint64]executor.VMHooks)}
	first := r.Register(newStub(1))
	r.Release(first)
	second := r.Register(newStub(2))
	if second <= first {
		t.Fatalf("expected monotonically increasing handles; got first=%d second=%d", first, second)
	}
}

// TestVMHooksRegistry_ConcurrentRegisterLookupRelease — exercises the
// RWMutex under load. Run with `go test -race`.
func TestVMHooksRegistry_ConcurrentRegisterLookupRelease(t *testing.T) {
	t.Parallel()

	r := &vmHooksRegistry{entries: make(map[uint64]executor.VMHooks)}

	const goroutines = 32
	const opsPerGoroutine = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(seed int) {
			defer wg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				id := r.Register(newStub(seed*1000 + i))
				_ = r.Lookup(id)
				r.Release(id)
			}
		}(g)
	}
	wg.Wait()

	// After all goroutines finish, all handles released; the map
	// should be empty.
	r.mu.RLock()
	remaining := len(r.entries)
	r.mu.RUnlock()
	if remaining != 0 {
		t.Fatalf("registry leaked %d entries after symmetric register/release", remaining)
	}
}

// TestLookupVMHooksOrPanic_PanicsOnMissingHandle confirms the
// "fail-loud" boundary in getVMHooksFromContextRawPtr.
func TestLookupVMHooksOrPanic_PanicsOnMissingHandle(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for missing handle, got none")
		}
	}()
	// Use a handle value that is guaranteed not to exist in the global
	// registry — the registry uses an atomic monotonic counter starting
	// from 1, so 1 << 63 is safely beyond any plausible live handle.
	lookupVMHooksOrPanic(uint64(1) << 63)
}
