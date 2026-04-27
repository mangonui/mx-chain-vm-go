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
