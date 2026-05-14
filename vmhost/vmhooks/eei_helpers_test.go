package vmhooks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetESDTRoles_WellFormedBufferReturnsExpectedRoles is a positive
// control: a buffer with two well-formed roles produces the OR of their
// individual role bits. This locks in the parser's success-path behavior
// so the fail-closed regression test below can be interpreted clearly
// (the failure case is "well-formed prefix + malformed suffix"; without
// the success-path test, "fail-closed returns 0" could be confused with
// "the parser never returns anything").
func TestGetESDTRoles_WellFormedBufferReturnsExpectedRoles(t *testing.T) {
	t.Parallel()

	// Build a well-formed buffer with the "ESDTRoleLocalMint" role
	// (assumed to map to a non-zero bit via roleFromByteArray).
	role := []byte("ESDTRoleLocalMint")
	buf := []byte{'\n', byte(len(role))}
	buf = append(buf, role...)

	got := getESDTRoles(buf, false)
	require.NotEqual(t, int64(0), got,
		"well-formed roles buffer must produce a non-zero role bitmask")
}

// TestGetESDTRoles_MalformedBufferFailsClosed is the ISSUE-084 regression
// pin. A buffer that begins with a well-formed role and trails off into
// a malformed length-byte (length larger than remaining buffer) MUST
// return 0 — the FAIL-CLOSED policy. An earlier version of the fix
// returned the partial result (the well-formed role's bits), which was
// fail-permissive and would have granted privileges from a corrupted
// blockchain-state record.
//
// The exact attack value is: prefix the buffer with the
// "ESDTRoleLocalMint" role (which would otherwise grant a real role
// bit), then append a malformed second-role record whose declared length
// exceeds the remaining buffer. Pre-fix-correct: returns 0 (no roles
// granted). Pre-current-fix (post-Claude-original-fix, "break"): returned
// the LocalMint role bit, which is fail-open. Post-current-fix
// ("return 0"): returns 0 — closed.
func TestGetESDTRoles_MalformedBufferFailsClosed(t *testing.T) {
	t.Parallel()

	// First role: well-formed "ESDTRoleLocalMint"
	role1 := []byte("ESDTRoleLocalMint")
	buf := []byte{'\n', byte(len(role1))}
	buf = append(buf, role1...)
	// Sanity: well-formed-only buffer produces non-zero
	wellFormedResult := getESDTRoles(buf, false)
	require.NotEqual(t, int64(0), wellFormedResult,
		"sanity: well-formed prefix must produce non-zero roles before we tack on the malformed suffix")

	// Append a malformed second record: \n + length=200 (way more than
	// remaining buffer of 0 bytes after the length byte itself).
	malformed := append(buf, '\n', 200)

	got := getESDTRoles(malformed, false)
	require.Equal(t, int64(0), got,
		"FAIL-CLOSED: malformed roles buffer MUST return 0 even if a well-formed role appears first — partial parse would grant privileges from corrupted state")
}

// TestGetESDTRoles_BufferEndingOnDelimiterFailsClosed pins the second
// over-read site fixed by ISSUE-084: a buffer that ends exactly on a \n
// delimiter (so the loop body's `currentIndex += 1` advances past the
// last byte). Pre-fix this would panic on `dataBuffer[currentIndex]`
// when reading the next length byte. Post-fix it must return 0.
func TestGetESDTRoles_BufferEndingOnDelimiterFailsClosed(t *testing.T) {
	t.Parallel()

	// Single \n byte — loop enters, advances past it, then would over-read.
	buf := []byte{'\n'}
	got := getESDTRoles(buf, false)
	require.Equal(t, int64(0), got,
		"FAIL-CLOSED: buffer ending on \\n delimiter must return 0, not panic")
}

// TestGetESDTRoles_EmptyBufferReturnsZero is a positive control for the
// trivial input shape.
func TestGetESDTRoles_EmptyBufferReturnsZero(t *testing.T) {
	t.Parallel()

	got := getESDTRoles([]byte{}, false)
	require.Equal(t, int64(0), got, "empty buffer must return 0")
}
