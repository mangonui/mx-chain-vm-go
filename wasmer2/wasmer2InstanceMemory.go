package wasmer2

import (
	"errors"
	"fmt"
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
//
// DEPRECATED — see issues/ISSUE-003.
//
// The returned slice is an alias over wasmer's linear memory, NOT a
// copy. After any subsequent `memory.Grow()`, wasmer may reallocate the
// backing buffer; the slice returned here then becomes a dangling
// pointer that crashes or returns wrong-but-plausible bytes when the
// caller next reads it. This is silent UAF.
//
// New callers MUST use [`Wasmer2Memory.ReadMemory`] instead, which
// returns a defensive copy and bounds-checks the requested range. Existing
// callers should migrate; the only in-tree user remaining is
// [`Wasmer2Instance.MemDump`] which itself is documented as test-only and
// has been migrated to ReadMemory in the same change set.
//
// This method is kept (rather than deleted) only so out-of-tree consumers
// don't hard-break on import; they should migrate ASAP.
//
// nolint
func (memory *Wasmer2Memory) Data() []byte {
	length := memory.Length()
	data := (*uint8)(cWasmerMemoryData(memory.cgoInstance))
	if data == nil || length == 0 {
		return []byte{}
	}

	return unsafe.Slice(data, length)
}

// ReadMemory returns a stable copy of the requested memory range.
func (memory *Wasmer2Memory) ReadMemory(offset uint32, length uint32) ([]byte, error) {
	end := uint64(offset) + uint64(length)
	totalLen := memory.Length()
	if end > uint64(totalLen) {
		return nil, errors.New("memory range out of bounds")
	}
	if length == 0 {
		return []byte{}, nil
	}

	data := (*uint8)(cWasmerMemoryData(memory.cgoInstance))
	if data == nil {
		return nil, errors.New("memory data pointer is nil")
	}

	copied := make([]byte, length)
	copy(copied, unsafe.Slice(data, totalLen)[offset:end])
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
