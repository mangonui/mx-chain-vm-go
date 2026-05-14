package vmhooks

// Integration regression tests for ISSUE-076, ISSUE-081, ISSUE-085 fixes.
//
// These tests exercise the actual vmhook entry points (not just helper
// arithmetic) to verify that each fix's bounds check fires and calls
// FailExecution BEFORE any dangerous downstream code path runs (allocation,
// MemLoad, etc.).
//
// The math-pinning unit tests in baseOps_helpers_test.go and
// metering_test.go (TestWriteLog_NumTopicsHashLenInt32OverflowDocumented,
// TestWriteEventLog_TopicDataTotalLenInt32OverflowDocumented,
// TestMeteringContext_BoundGasLimit_NegativeReturnsZero) prove the
// arithmetic is fixed in isolation. The integration tests in THIS file
// prove the vmhook entry points actually invoke the fix on the call path
// — closing the regression-prevention gap a future refactor could open.
//
// Mock setup follows the pattern in managedei_drwa_test.go: VMHostStub
// + RuntimeContextWrapper (with overridable FailExecutionFunc) +
// MeteringContextMock. Output and other contexts are unset (returning
// nil) because the failure paths under test do NOT dereference them.

import (
	"testing"

	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

// newCapturingHost builds a minimal host whose only observable side
// effect is recording the error that FailExecution is called with.
// Returns the host, and a pointer to the captured-error slot.
//
// The host wires up Runtime + Metering only. Other context getters
// return nil. Tests that exercise vmhook failure paths must NOT touch
// any other context, by construction (the fix-under-test rejects
// before those contexts are read).
func newCapturingHost(t *testing.T, gasLeft uint64) (*contextmock.VMHostStub, *error) {
	t.Helper()

	metering := &contextmock.MeteringContextMock{
		GasLeftMock: gasLeft,
	}

	baseRuntime := &contextmock.RuntimeContextMock{}
	baseRuntimeIfc := vmhost.RuntimeContext(baseRuntime)
	runtime := contextmock.NewRuntimeContextWrapper(&baseRuntimeIfc)

	var failedErr error
	runtime.FailExecutionFunc = func(err error) {
		failedErr = err
	}

	host := &contextmock.VMHostStub{
		RuntimeCalled: func() vmhost.RuntimeContext {
			return runtime
		},
		MeteringCalled: func() vmhost.MeteringContext {
			return metering
		},
	}

	return host, &failedErr
}

// === ISSUE-076 — WriteLog OOM ===

// TestWriteLog_RejectsNumTopicsAtOverflowExploitValue is the integration
// regression pin for ISSUE-076. The exploit value numTopics = 2^30 was
// the exact input that, pre-fix, wrapped the int32 multiplication
// `numTopics * vmhost.HashLen` to 0, bypassed the gas check, and
// triggered a ~24 GiB `make([][]byte, numTopics)` allocation.
//
// Post-fix: WriteLog must reject this input at the explicit
// `numTopics > maxNumArgumentsFromMemory` guard at the top of the
// function, BEFORE any arithmetic on numTopics or any allocation.
//
// Mutation check (manual): comment out the `if numTopics > maxNumArgumentsFromMemory`
// guard in baseOps.go::WriteLog → this test must FAIL with no
// FailExecution called. Restore the guard → test passes.
func TestWriteLog_RejectsNumTopicsAtOverflowExploitValue(t *testing.T) {
	t.Parallel()

	host, failedErr := newCapturingHost(t, 1_000_000)

	// The exact ISSUE-076 attack value.
	const exploitNumTopics = int32(1 << 30) // 2^30 = 1,073,741,824

	NewVMHooksImpl(host).WriteLog(0, 0, 0, exploitNumTopics)

	require.Error(t, *failedErr,
		"WriteLog MUST call FailExecution for numTopics=2^30 — the ISSUE-076 OOM exploit value")
	require.Contains(t, (*failedErr).Error(), "exceeds maximum",
		"FailExecution error MUST come from the explicit `exceeds maximum` guard, not a downstream check — proves the new bounds check at the top of WriteLog fired before any allocation could occur")
}

// TestWriteLog_RejectsNumTopicsAtCapPlusOne verifies the OVER-cap branch
// of the bound check: numTopics = maxNumArgumentsFromMemory + 1 must be
// rejected with the "exceeds maximum" error.
//
// SCOPE NOTE: this test does NOT exercise the at-cap acceptance branch
// (numTopics == maxNumArgumentsFromMemory). The bound is implemented as
// strict `>`, so at-cap should be accepted, but verifying that requires
// running the full WriteLog body downstream of the bound check (gas
// charge → MemLoad → make([][]byte, 16M) → 16M-iteration topic-load
// loop). With the mock metering returning success and zero-length
// MemLoad short-circuiting, the make would still allocate ~384 MiB of
// slice headers — feasible but slow and noisy for a unit test. The
// at-cap acceptance is verified by code review of the strict `>`
// operator at the bound check; this test pins only the rejection
// behavior at cap+1.
func TestWriteLog_RejectsNumTopicsAtCapPlusOne(t *testing.T) {
	t.Parallel()

	host, failedErr := newCapturingHost(t, 1_000_000)

	// One beyond the cap.
	tooMany := int32(maxNumArgumentsFromMemory + 1)

	NewVMHooksImpl(host).WriteLog(0, 0, 0, tooMany)

	require.Error(t, *failedErr,
		"WriteLog MUST reject numTopics > maxNumArgumentsFromMemory")
	require.Contains(t, (*failedErr).Error(), "exceeds maximum")
}

// TestWriteLog_RejectsNegativeNumTopics pins the original negative-input
// guard (which existed pre-fix but ran AFTER the gas calc; ISSUE-076's
// fix moved it BEFORE). This test ensures negative input is still
// rejected with the right error after the reordering.
func TestWriteLog_RejectsNegativeNumTopics(t *testing.T) {
	t.Parallel()

	host, failedErr := newCapturingHost(t, 1_000_000)

	NewVMHooksImpl(host).WriteLog(0, 0, 0, int32(-1))

	require.ErrorIs(t, *failedErr, vmhost.ErrNegativeLength,
		"negative numTopics MUST be rejected with ErrNegativeLength")
}

// === ISSUE-081 — WriteEventLog ungated MemLoad + int32 overflow ===

// TestWriteEventLog_RejectsNegativeNumTopics pins the new
// `if numTopics < 0` guard at the top of WriteEventLog. Pre-fix, no
// such guard existed; getArgumentsFromMemory's own check would catch it
// later, but my ISSUE-081 fix added an explicit upfront guard so the
// uint64 cast in the pre-load gas calc never sign-extends.
func TestWriteEventLog_RejectsNegativeNumTopics(t *testing.T) {
	t.Parallel()

	host, failedErr := newCapturingHost(t, 1_000_000)

	NewVMHooksImpl(host).WriteEventLog(int32(-1), 0, 0, 0, 0)

	require.ErrorIs(t, *failedErr, vmhost.ErrNegativeLength,
		"WriteEventLog MUST reject negative numTopics with ErrNegativeLength at the new upfront guard")
}

// === ISSUE-085 — Crypto vmhooks negative-length cast ===

// TestCryptoVMHooks_RejectNegativeLengthBeforeUint64Cast is a
// table-driven integration test covering all six crypto vmhooks fixed
// by ISSUE-085. Each MUST reject negative length / messageLength with
// ErrNegativeLength via the new guard added at the function head,
// BEFORE any uint64 cast in the gas calculation.
//
// Mutation check (per row): comment out the `if length < 0` (or
// `if messageLength < 0`) guard in the corresponding cryptoei.go function
// → that table row must FAIL. Restore → passes.
//
// Note: the test exercises each vmhook with the relevant negative
// argument; offsets/handles are 0 because the failure path doesn't
// dereference them.
func TestCryptoVMHooks_RejectNegativeLengthBeforeUint64Cast(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		call func(v *VMHooksImpl)
	}{
		{
			name: "Sha256",
			call: func(v *VMHooksImpl) { v.Sha256(0, -1, 0) },
		},
		{
			name: "Keccak256",
			call: func(v *VMHooksImpl) { v.Keccak256(0, -1, 0) },
		},
		{
			name: "Ripemd160",
			call: func(v *VMHooksImpl) { v.Ripemd160(0, -1, 0) },
		},
		{
			name: "VerifyBLS",
			call: func(v *VMHooksImpl) { v.VerifyBLS(0, 0, -1, 0) },
		},
		{
			name: "VerifyEd25519",
			call: func(v *VMHooksImpl) { v.VerifyEd25519(0, 0, -1, 0) },
		},
		{
			name: "VerifyCustomSecp256k1",
			call: func(v *VMHooksImpl) { v.VerifyCustomSecp256k1(0, 0, 0, -1, 0, 0) },
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			host, failedErr := newCapturingHost(t, 1_000_000)
			tc.call(NewVMHooksImpl(host))

			require.ErrorIs(t, *failedErr, vmhost.ErrNegativeLength,
				"%s MUST reject negative length/messageLength with ErrNegativeLength at the new upfront guard — proves the ISSUE-085 fix is wired correctly at this entry point",
				tc.name)
		})
	}
}
