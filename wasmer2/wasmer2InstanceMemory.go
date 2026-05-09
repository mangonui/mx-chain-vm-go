package wasmer2

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ = (executor.Memory)((*Wasmer2Memory)(nil))

// Wasmer2Instance represents a WebAssembly instance.
type Wasmer2Memory struct {
	// The underlying WebAssembly instance.
	cgoInstance *cWasmerInstanceT
}

// Length calculates the memory length (in bytes).
func (memory *Wasmer2Memory) Length() uint32 {
	return uint32(cWasmerMemoryDataLength(memory.cgoInstance))
}

// Data returns a slice of bytes over the WebAssembly memory.
// nolint
func (memory *Wasmer2Memory) Data() []byte {
	var length = memory.Length()
	var data = (*uint8)(cWasmerMemoryData(memory.cgoInstance))

	var header reflect.SliceHeader
	header = *(*reflect.SliceHeader)(unsafe.Pointer(&header))

	header.Data = uintptr(unsafe.Pointer(data))
	header.Len = int(length)
	header.Cap = int(length)

	return *(*[]byte)(unsafe.Pointer(&header))
}

// ReadMemory returns a stable copy of the requested memory range.
func (memory *Wasmer2Memory) ReadMemory(offset uint32, length uint32) ([]byte, error) {
	data := memory.Data()
	end := uint64(offset) + uint64(length)
	if end > uint64(len(data)) {
		return nil, errors.New("memory range out of bounds")
	}

	copied := make([]byte, length)
	copy(copied, data[offset:uint32(end)])
	return copied, nil
}

// Grow the memory by a number of pages (65kb each).
func (memory *Wasmer2Memory) Grow(numberOfPages uint32) error {
	var growResult = cWasmerMemoryGrow(memory.cgoInstance, cUint32T(numberOfPages))

	if growResult != cWasmerOk {
		var lastError, err = GetLastError()
		var errorMessage = "Failed to grow the memory:\n    %s"

		if err != nil {
			errorMessage = fmt.Sprintf(errorMessage, "(unknown details)")
		} else {
			errorMessage = fmt.Sprintf(errorMessage, lastError)
		}

		return fmt.Errorf("memory grow error: %s", errorMessage)
	}

	return nil
}

// Destroy destroys inner memory. Does nothing in wasmer2.
func (memory *Wasmer2Memory) Destroy() {
}

// IsInterfaceNil returns true if underlying object is nil
func (memory *Wasmer2Memory) IsInterfaceNil() bool {
	return memory == nil
}
