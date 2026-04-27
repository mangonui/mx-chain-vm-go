package vmhooks

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountDRWASyncOperations_ValidPayload(t *testing.T) {
	t.Parallel()

	count, ok := countDRWASyncOperations(buildValidDRWASyncPayload(3))
	require.True(t, ok)
	require.Equal(t, 3, count)
}

func TestCountDRWASyncOperations_MalformedPayload(t *testing.T) {
	t.Parallel()

	count, ok := countDRWASyncOperations(buildMalformedDRWASyncPayload())
	require.False(t, ok)
	require.Zero(t, count)
}

func TestCountDRWASyncOperations_RejectsTooManyOperations(t *testing.T) {
	t.Parallel()

	count, ok := countDRWASyncOperations(buildValidDRWASyncPayload(maxDRWASyncOps + 1))
	require.False(t, ok)
	require.Zero(t, count)
}

func TestCountDRWASyncOperations_RejectsOversizedField(t *testing.T) {
	t.Parallel()

	count, ok := countDRWASyncOperations(buildOversizedFieldDRWASyncPayload())
	require.False(t, ok)
	require.Zero(t, count)
}

func TestSkipLenPrefixed_TruncatedInput(t *testing.T) {
	t.Parallel()

	_, ok := skipLenPrefixed([]byte{0, 0, 0, 4, 1, 2})
	require.False(t, ok)
}

func TestSafeMulUint64_Overflow(t *testing.T) {
	t.Parallel()

	_, overflow := safeMulUint64(math.MaxUint64, 2)
	require.True(t, overflow)
}

func TestSafeMulUint64_NoOverflow(t *testing.T) {
	t.Parallel()

	result, overflow := safeMulUint64(7, 9)
	require.False(t, overflow)
	require.Equal(t, uint64(63), result)
}

func buildValidDRWASyncPayload(numOps int) []byte {
	payload := make([]byte, 33)
	for i := 0; i < numOps; i++ {
		payload = append(payload, byte(i%3))
		payload = appendLenPrefixed(payload, []byte("TOKEN-123"))
		payload = appendLenPrefixed(payload, []byte("erd1holder"))
		payload = append(payload, 0, 0, 0, 0, 0, 0, 0, byte(i+1))
		payload = appendLenPrefixed(payload, []byte("body"))
	}

	return payload
}

func buildMalformedDRWASyncPayload() []byte {
	payload := make([]byte, 33)
	payload = append(payload, 0x01)
	payload = append(payload, 0, 0, 0, 5, 'b', 'a')
	return payload
}

func buildOversizedFieldDRWASyncPayload() []byte {
	payload := make([]byte, 33)
	payload = append(payload, 0x01)
	payload = appendLenPrefixed(payload, bytesOfLen(maxDRWASyncFieldLen+1))
	payload = appendLenPrefixed(payload, []byte("erd1holder"))
	payload = append(payload, 0, 0, 0, 0, 0, 0, 0, 1)
	payload = appendLenPrefixed(payload, []byte("body"))
	return payload
}

func bytesOfLen(length int) []byte {
	value := make([]byte, length)
	for i := range value {
		value[i] = 'a'
	}
	return value
}

func appendLenPrefixed(dst []byte, value []byte) []byte {
	length := len(value)
	dst = append(dst,
		byte(length>>24),
		byte(length>>16),
		byte(length>>8),
		byte(length),
	)
	dst = append(dst, value...)
	return dst
}
