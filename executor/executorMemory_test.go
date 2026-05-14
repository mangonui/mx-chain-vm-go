package executor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type memoryStub struct {
	data []byte
}

func (stub *memoryStub) Length() uint32 {
	return uint32(len(stub.data))
}

func (stub *memoryStub) Data() []byte {
	return stub.data
}

func (stub *memoryStub) Grow(uint32) error {
	return nil
}

func (stub *memoryStub) Destroy() {}

func (stub *memoryStub) IsInterfaceNil() bool {
	return stub == nil
}

func TestMemLoadFromMemory(t *testing.T) {
	mem := &memoryStub{data: []byte{1, 2, 3, 4}}

	t.Run("loads bounded slice", func(t *testing.T) {
		data, err := MemLoadFromMemory(mem, 1, 2)

		require.Nil(t, err)
		require.Equal(t, []byte{2, 3}, data)
	})

	t.Run("clips upper bound", func(t *testing.T) {
		data, err := MemLoadFromMemory(mem, 2, 10)

		require.Nil(t, err)
		require.Equal(t, []byte{3, 4}, data)
	})

	t.Run("negative length", func(t *testing.T) {
		data, err := MemLoadFromMemory(mem, 0, -1)

		require.Nil(t, data)
		require.True(t, errors.Is(err, ErrMemoryNegativeLength))
	})

	t.Run("bad lower bound", func(t *testing.T) {
		data, err := MemLoadFromMemory(mem, -1, 1)

		require.Nil(t, data)
		require.True(t, errors.Is(err, ErrMemoryBadBounds))
	})
}

func TestMemStoreToMemory(t *testing.T) {
	mem := &memoryStub{data: []byte{1, 2, 3, 4}}

	err := MemStoreToMemory(mem, 1, []byte{9, 8})
	require.Nil(t, err)
	require.Equal(t, []byte{1, 9, 8, 4}, mem.data)

	err = MemStoreToMemory(mem, -1, []byte{7})
	require.True(t, errors.Is(err, ErrMemoryBadBoundsLower))

	err = MemStoreToMemory(mem, 3, []byte{7, 6})
	require.True(t, errors.Is(err, ErrMemoryBadBoundsUpper))
}
