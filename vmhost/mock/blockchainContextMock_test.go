package mock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlockchainContextMock_ApplyDRWASyncEnvelopeBytesRequiresExplicitMocking(t *testing.T) {
	t.Parallel()

	err := (&BlockchainContextMock{}).ApplyDRWASyncEnvelopeBytes([]byte("payload"), []byte("caller"))
	require.ErrorIs(t, err, ErrApplyDRWASyncEnvelopeBytesNotMocked)
}

func TestBlockchainContextMock_ApplyDRWASyncEnvelopeBytesUsesCallback(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected DRWA sync failure")
	mock := &BlockchainContextMock{
		ApplyDRWASyncEnvelopeBytesCalled: func(payload []byte, callerAddress []byte) error {
			require.Equal(t, []byte("payload"), payload)
			require.Equal(t, []byte("caller"), callerAddress)
			return expectedErr
		},
	}

	err := mock.ApplyDRWASyncEnvelopeBytes([]byte("payload"), []byte("caller"))
	require.ErrorIs(t, err, expectedErr)
}
