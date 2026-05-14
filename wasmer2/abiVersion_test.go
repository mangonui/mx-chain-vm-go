package wasmer2

import "testing"

// TestCWasmerAPIVersion_MatchesExpected — ISSUE-020. The single
// authoritative liveness check for the FFI ABI handshake. If the linked
// .so / .dylib reports a different version than the Go bridge expects,
// this fails immediately and the build refresh story (header + lib +
// expectedAPIVersion bump) is the obvious diff.
func TestCWasmerAPIVersion_MatchesExpected(t *testing.T) {
	t.Parallel()

	got := cWasmerAPIVersion()
	if got != expectedAPIVersion {
		t.Fatalf(
			"libvmexeccapi reports ABI v%d, Go bridge expected v%d. "+
				"Either the linked .so/.dylib is stale (refresh from "+
				"mx-vm-executor-rs/target/release/) or the bridge constant "+
				"`expectedAPIVersion` needs bumping to match the lib.",
			got, expectedAPIVersion,
		)
	}
}

// TestCheckAPIVersion_NilOnMatch — when the lib version matches, the
// guarded handshake returns nil and is safe to call repeatedly.
func TestCheckAPIVersion_NilOnMatch(t *testing.T) {
	t.Parallel()

	if err := checkAPIVersion(); err != nil {
		t.Fatalf("checkAPIVersion returned %v with matching lib; expected nil", err)
	}
	// Idempotency: a second call returns the same result without
	// re-invoking the underlying cgo call (sync.Once-guarded).
	if err := checkAPIVersion(); err != nil {
		t.Fatalf("second checkAPIVersion call returned %v; expected nil", err)
	}
}
