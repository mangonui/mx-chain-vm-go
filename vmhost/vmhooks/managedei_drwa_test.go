package vmhooks

import (
	"testing"

	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

func TestManagedDRWASyncMirrorRejectsUnauthorizedOpcode(t *testing.T) {
	t.Parallel()

	metering := &contextmock.MeteringContextMock{
		GasLeftMock: 1000,
	}

	var failedErr error
	baseRuntime := &contextmock.RuntimeContextMock{}
	baseRuntimeIfc := vmhost.RuntimeContext(baseRuntime)
	runtime := contextmock.NewRuntimeContextWrapper(&baseRuntimeIfc)
	runtime.FailExecutionFunc = func(err error) {
		failedErr = err
	}

	host := &contextmock.VMHostStub{
		IsAllowedToExecuteCalled: func(opcode string) bool {
			require.Equal(t, managedDRWASyncMirrorName, opcode)
			return false
		},
		RuntimeCalled: func() vmhost.RuntimeContext {
			return runtime
		},
		MeteringCalled: func() vmhost.MeteringContext {
			return metering
		},
	}

	result := NewVMHooksImpl(host).ManagedDRWASyncMirror(1)
	require.Equal(t, int32(1), result)
	require.ErrorIs(t, failedErr, vmhost.ErrOpcodeIsNotAllowed)
}
