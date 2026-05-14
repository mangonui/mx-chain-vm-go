package wasmer2

import "testing"

// TestCWasmerCacheFree_NullPtrIsSafeNoOp — ISSUE-009. The Rust side's
// `vm_exec_cache_free` no-ops on null/zero inputs (defensive, see
// capi_instance_cache.rs); the Go wrapper must round-trip that without
// crashing. Real end-to-end coverage of the (alloc, hand to Go, free)
// round trip lives in the scenario tests that exercise Cache() via
// runtime.go:335 — this is the cheap unit-test stub.
func TestCWasmerCacheFree_NullPtrIsSafeNoOp(t *testing.T) {
	t.Parallel()

	// Must not panic / crash. The Rust impl's null-and-zero guards
	// make this a defined no-op; if a future refactor on the Rust
	// side removes those guards, this test fires immediately.
	cWasmerCacheFree(nil, 0)
	cWasmerCacheFree(nil, 100) // ptr-null short-circuit
}
