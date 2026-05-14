package wasmer2

import (
	"sync"
	"testing"
)

// TestWasmer2ExecutorDestroyIsIdempotent verifies that Destroy can be
// called multiple times on the same executor without crashing or
// double-freeing. Uses a zero-value executor (all C pointers nil) so the
// test does not depend on the wasmer C library being initialised — the
// nil-checks inside Destroy short-circuit every cFree, and what is being
// exercised here is the sync.Once gate added for ISSUE-014.
func TestWasmer2ExecutorDestroyIsIdempotent(t *testing.T) {
	t.Parallel()

	exec := &Wasmer2Executor{}

	// Single thread, multiple calls — must not panic.
	exec.Destroy()
	exec.Destroy()
	exec.Destroy()
}

// TestWasmer2ExecutorDestroyConcurrent is the goroutine-race scenario
// from ISSUE-014: N goroutines call Destroy simultaneously on the same
// executor. Must complete without panic and without tripping the Go
// race detector (run via `go test -race ./wasmer2/`).
func TestWasmer2ExecutorDestroyConcurrent(t *testing.T) {
	t.Parallel()

	exec := &Wasmer2Executor{}

	const goroutines = 32
	var wg sync.WaitGroup
	wg.Add(goroutines)
	start := make(chan struct{})
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			<-start
			exec.Destroy()
		}()
	}
	close(start)
	wg.Wait()
}

// TestWasmer2ExecutorDestroyNilReceiver guards the documented "safe on
// nil receiver" behaviour. Going through sync.Once on a nil receiver
// would deref a zero-value field; the early return above the .Do call
// prevents that.
func TestWasmer2ExecutorDestroyNilReceiver(t *testing.T) {
	t.Parallel()

	var exec *Wasmer2Executor
	exec.Destroy()
}
