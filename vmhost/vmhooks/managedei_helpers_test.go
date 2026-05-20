package vmhooks

import (
	"encoding/hex"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountDRWASyncOperations_ValidPayload(t *testing.T) {
	t.Parallel()

	count, ok := countDRWASyncOperations(buildValidDRWASyncPayloadV1(3))
	require.True(t, ok)
	require.Equal(t, 3, count)
}

func TestCountDRWASyncOperations_ValidRecoveryPayload(t *testing.T) {
	t.Parallel()

	count, ok := countDRWASyncOperations(buildValidDRWASyncPayloadV2(2))
	require.True(t, ok)
	require.Equal(t, 2, count)
}

func TestCountDRWASyncOperations_AcceptsRustGeneratedFixtures(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		fixture       string
		expectedCount int
	}{
		{
			name:          "schema v1",
			fixture:       "sync-envelope-v1.hex",
			expectedCount: 1,
		},
		{
			name:          "schema v2 recovery",
			fixture:       "sync-envelope-v2-recovery.hex",
			expectedCount: 2,
		},
		{
			name:          "schema v1 all operation tags",
			fixture:       "sync-envelope-v1-all-op-tags.hex",
			expectedCount: 9,
		},
		{
			name:          "schema v1 near payload cap",
			fixture:       "sync-envelope-v1-near-cap.hex",
			expectedCount: maxDRWASyncOps,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			payload := readRustDRWASyncFixture(t, testCase.fixture)
			count, ok := countDRWASyncOperations(payload)
			require.True(t, ok)
			require.Equal(t, testCase.expectedCount, count)
		})
	}
}

func TestCountDRWASyncOperations_RejectsUnsupportedSchemaVersion(t *testing.T) {
	t.Parallel()

	payload := makeDRWAHeader(3)

	count, ok := countDRWASyncOperations(payload)
	require.False(t, ok)
	require.Zero(t, count)
}

func TestCountDRWASyncOperations_MalformedPayload(t *testing.T) {
	t.Parallel()

	count, ok := countDRWASyncOperations(buildMalformedDRWASyncPayload())
	require.False(t, ok)
	require.Zero(t, count)
}

func TestCountDRWASyncOperations_RejectsTooManyOperations(t *testing.T) {
	t.Parallel()

	count, ok := countDRWASyncOperations(buildValidDRWASyncPayloadV1(maxDRWASyncOps + 1))
	require.False(t, ok)
	require.Zero(t, count)
}

func TestCountDRWASyncOperations_RejectsOversizedField(t *testing.T) {
	t.Parallel()

	count, ok := countDRWASyncOperations(buildOversizedFieldDRWASyncPayload())
	require.False(t, ok)
	require.Zero(t, count)
}

func TestCountDRWASyncOperations_RejectsInvalidOperationTag(t *testing.T) {
	t.Parallel()

	payload := makeDRWAHeader(1)
	payload = append(payload, 9)
	payload = appendLenPrefixed(payload, []byte("TOKEN-123"))
	payload = appendLenPrefixed(payload, []byte("erd1holder"))
	payload = append(payload, 0, 0, 0, 0, 0, 0, 0, 1)
	payload = appendLenPrefixed(payload, []byte("body"))

	count, ok := countDRWASyncOperations(payload)
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

func buildValidDRWASyncPayloadV1(numOps int) []byte {
	payload := makeDRWAHeader(1)
	for i := 0; i < numOps; i++ {
		payload = appendDRWASyncOperation(payload, byte(i%9), byte(i+1))
	}

	return payload
}

func buildValidDRWASyncPayloadV2(numOps int) []byte {
	payload := makeDRWAHeader(2)
	payload = appendLenPrefixed(payload, []byte("pre-recovery-state-hash"))
	payload = append(payload, 0, 2)
	payload = appendLenPrefixed(payload, []byte("TOKEN-111111"))
	payload = appendLenPrefixed(payload, []byte("ASSET-222222"))
	payload = append(payload, byte(numOps>>8), byte(numOps))
	for i := 0; i < numOps; i++ {
		payload = appendDRWASyncOperation(payload, byte(i%9), byte(i+1))
	}

	return payload
}

func buildMalformedDRWASyncPayload() []byte {
	payload := makeDRWAHeader(1)
	payload = append(payload, 0x01)
	payload = append(payload, 0, 0, 0, 5, 'b', 'a')
	return payload
}

func buildOversizedFieldDRWASyncPayload() []byte {
	payload := makeDRWAHeader(1)
	payload = append(payload, 0x01)
	payload = appendLenPrefixed(payload, bytesOfLen(maxDRWASyncFieldLen+1))
	payload = appendLenPrefixed(payload, []byte("erd1holder"))
	payload = append(payload, 0, 0, 0, 0, 0, 0, 0, 1)
	payload = appendLenPrefixed(payload, []byte("body"))
	return payload
}

func makeDRWAHeader(schemaVersion byte) []byte {
	payload := make([]byte, 32)
	payload = append(payload, 0, schemaVersion)
	payload = append(payload, 0)
	return payload
}

func appendDRWASyncOperation(payload []byte, opTag byte, version byte) []byte {
	payload = append(payload, opTag)
	payload = appendLenPrefixed(payload, []byte("TOKEN-123"))
	payload = appendLenPrefixed(payload, []byte("erd1holder"))
	payload = append(payload, 0, 0, 0, 0, 0, 0, 0, version)
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

func readRustDRWASyncFixture(t *testing.T, fixtureName string) []byte {
	t.Helper()

	fixtureDir := os.Getenv("DRWA_SYNC_FIXTURE_DIR")
	if fixtureDir == "" {
		fixtureDir = filepath.Join(
			"..", "..", "..",
			"mx-sdk-rs", "contracts", "drwa", "common", "testdata", "drwa-sync-fixtures",
		)
	}
	path := filepath.Join(fixtureDir, fixtureName)
	contents, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("DRWA sync fixture %q not available at %s; set DRWA_SYNC_FIXTURE_DIR in split-repo CI", fixtureName, path)
	}
	require.NoError(t, err)

	payload, err := hex.DecodeString(strings.Join(strings.Fields(string(contents)), ""))
	require.NoError(t, err)

	return payload
}
