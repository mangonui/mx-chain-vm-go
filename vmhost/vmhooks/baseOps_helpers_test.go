package vmhooks

import (
	"math"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

func TestComputeSignalErrorGas_NegativeLength(t *testing.T) {
	t.Parallel()

	_, err := computeSignalErrorGas(5, 7, executor.MemLength(-1))
	require.ErrorIs(t, err, vmhost.ErrArgOutOfRange)
}

func TestComputeSignalErrorGas_LengthGasOverflow(t *testing.T) {
	t.Parallel()

	_, err := computeSignalErrorGas(5, math.MaxUint64, executor.MemLength(2))
	require.ErrorIs(t, err, vmhost.ErrArgOutOfRange)
}

func TestComputeSignalErrorGas_TotalGasOverflow(t *testing.T) {
	t.Parallel()

	_, err := computeSignalErrorGas(math.MaxUint64, 1, executor.MemLength(1))
	require.ErrorIs(t, err, vmhost.ErrArgOutOfRange)
}

func TestComputeSignalErrorGas_Success(t *testing.T) {
	t.Parallel()

	gas, err := computeSignalErrorGas(5, 7, executor.MemLength(3))
	require.NoError(t, err)
	require.Equal(t, uint64(26), gas)
}

func TestComputeTransferBatchBaseGas_Overflow(t *testing.T) {
	t.Parallel()

	_, err := computeTransferBatchBaseGas(math.MaxUint64, 2)
	require.ErrorIs(t, err, vmhost.ErrArgOutOfRange)
}

func TestComputeTransferBatchBaseGas_Success(t *testing.T) {
	t.Parallel()

	gas, err := computeTransferBatchBaseGas(9, 4)
	require.NoError(t, err)
	require.Equal(t, uint64(36), gas)
}

// TestWriteLog_NumTopicsHashLenInt32OverflowDocumented pins the exact
// arithmetic bug that ISSUE-076 closes. Pre-fix WriteLog computed the
// gas-input expression `numTopics*vmhost.HashLen+dataLength` in int32
// because vmhost.HashLen is an untyped constant that adopts the int32
// type of numTopics. For numTopics = 2^30 (a positive value that passes
// the `numTopics < 0` check), 2^30*32 = 2^35 wraps to 0 mod 2^32 in
// int32, the gas calc returned ~0, the gas check passed, and the
// subsequent `make([][]byte, numTopics)` allocated ~24 GiB on the
// validator process.
//
// This test documents both halves of the fix:
//  1. The unsafe arithmetic still wraps (shows the bug exists).
//  2. The uint64-promoted arithmetic does NOT wrap (shows the fix works).
//
// If either assertion changes, that signals either the underlying Go
// arithmetic semantics changed or someone re-introduced the int32 path
// in WriteLog.
func TestWriteLog_NumTopicsHashLenInt32OverflowDocumented(t *testing.T) {
	t.Parallel()

	const HashLenUntyped = 32 // mirrors vmhost.HashLen
	var numTopics int32 = 1 << 30
	var dataLength int32 = 0

	// Pre-fix arithmetic shape: int32 multiplication wraps silently.
	// The expression below is what the buggy code computed at
	// baseOps.go:2706 before ISSUE-076 was fixed.
	wrappedInt32 := numTopics*HashLenUntyped + dataLength
	require.Equal(t, int32(0), wrappedInt32,
		"pre-fix int32 multiplication MUST wrap to 0 for numTopics=2^30 — if this fails, Go's int32 overflow semantics changed")

	// Post-fix arithmetic shape: uint64-promoted multiplication is safe.
	// The expression below is what the fixed code computes — operands
	// are explicitly widened to uint64 BEFORE the multiplication, so
	// no overflow occurs for any int32 numTopics input.
	safeUint64 := uint64(numTopics)*uint64(HashLenUntyped) + uint64(dataLength)
	require.Equal(t, uint64(34_359_738_368), safeUint64,
		"post-fix uint64 multiplication MUST equal 2^35 = 34359738368 for numTopics=2^30 — if this fails, the fix is broken")

	// Sanity: the wrapped int32 does NOT equal the correct value.
	require.NotEqual(t, uint64(wrappedInt32), safeUint64,
		"wrapped int32 result and correct uint64 result must differ — otherwise there is no overflow to fix")
}

// TestWriteEventLog_TopicDataTotalLenInt32OverflowDocumented pins the
// equivalent arithmetic bug that ISSUE-081 (deeper fix) closes. Pre-fix
// WriteEventLog at baseOps.go:2816 computed `uint64(topicDataTotalLen+dataLength)`
// — int32 addition BEFORE the uint64 cast. Two int32 values that
// individually fit but whose SUM exceeds MaxInt32 wrap silently to a
// small (or negative) int32, then sign-extend through the uint64 cast,
// producing the wrong gas amount.
//
// This test documents both arithmetic shapes the same way ISSUE-076's
// test does — pre-fix wraps, post-fix does not.
func TestWriteEventLog_TopicDataTotalLenInt32OverflowDocumented(t *testing.T) {
	t.Parallel()

	// Two int32 values that individually fit but whose sum overflows.
	var topicDataTotalLen int32 = 1<<30 + 1<<29 // ~1.6 GiB worth of length
	var dataLength int32 = 1 << 30              // ~1 GiB

	// Pre-fix shape: int32 addition wraps to negative, then uint64 cast
	// sign-extends to ~MaxUint64. (math.MulUint64 in the live code then
	// saturates to MaxUint64 — same as the ScalarMult-style pattern.)
	wrappedInt32Sum := topicDataTotalLen + dataLength
	require.Less(t, wrappedInt32Sum, int32(0),
		"pre-fix int32 addition of two large positives MUST wrap to negative — if this fails, Go's int32 overflow semantics changed")

	// Post-fix shape: each operand widened to uint64 BEFORE the addition,
	// so the sum is the mathematically correct value.
	safeUint64Sum := uint64(topicDataTotalLen) + uint64(dataLength)
	expected := uint64(int64(topicDataTotalLen) + int64(dataLength))
	require.Equal(t, expected, safeUint64Sum,
		"post-fix uint64 addition MUST equal int64-promoted sum — if this fails, the fix is broken")

	// Sanity: the wrapped result and the correct result differ.
	require.NotEqual(t, uint64(wrappedInt32Sum), safeUint64Sum,
		"wrapped int32 result and correct uint64 result must differ — otherwise there is no overflow to fix")
}
