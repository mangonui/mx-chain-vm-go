package wasmer2

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestWasmer2Instance_IDUsesDecimalOpaqueHandle(t *testing.T) {
	t.Parallel()

	backing := new(byte)
	instance := &Wasmer2Instance{
		cgoInstance: (*cWasmerInstanceT)(unsafe.Pointer(backing)),
	}

	require.NotContains(t, instance.ID(), "0x")
}

func TestWasmer2Instance_VMHooksPtrIsDocumentedNoOp(t *testing.T) {
	t.Parallel()

	instance := &Wasmer2Instance{}

	instance.SetVMHooksPtr(12345)

	require.Equal(t, uintptr(0), instance.GetVMHooksPtr())
}

func TestWasmer2Instance_CleanAlreadyCleanedIsConcurrentSafe(t *testing.T) {
	t.Parallel()

	instance := &Wasmer2Instance{AlreadyClean: true}

	const goroutines = 32
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			require.False(t, instance.Clean())
			require.True(t, instance.IsAlreadyCleaned())
		}()
	}
	wg.Wait()
}
