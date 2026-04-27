package vmhookstest

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	test "github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks"
	"github.com/stretchr/testify/require"
)

const (
	drwaSyncOpsCapForTests      = 256
	drwaSyncFieldLenCapForTests = 64 * 1024
)

func TestManagedDRWASyncMirror_NodeHookFailure_Reverts(t *testing.T) {
	t.Parallel()

	stubErr := errors.New("node hook: sync envelope rejected")

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("testFunction", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						payloadHandle := managedTypes.NewManagedBufferFromBytes(buildValidDRWASyncPayload(1))
						result := hooks.ManagedDRWASyncMirror(payloadHandle)
						require.Equal(t, int32(1), result)

						return instanceMock
					})
				}),
		).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.ProvidedBlockchainHook = &contextmock.BlockchainHookStub{
				ApplyDRWASyncEnvelopeBytesCalled: func(_ []byte, _ []byte) error {
					return stubErr
				},
			}
		}).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("testFunction").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				ReturnMessage(stubErr.Error())
		})
	require.NoError(t, err)
}

func TestManagedDRWANativeGovernanceQuery_ForwardsToNodeHook(t *testing.T) {
	t.Parallel()

	expectedResult := []byte(`{"version":1}`)
	queryKey := []byte("TOKEN-abcdef")

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("testFunction", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						keyHandle := managedTypes.NewManagedBufferFromBytes(queryKey)
						destHandle := managedTypes.NewManagedBufferFromBytes(nil)
						result := hooks.ManagedDRWANativeGovernanceQuery(0, keyHandle, destHandle)
						require.Equal(t, int32(0), result)

						actual, err := managedTypes.GetBytes(destHandle)
						require.NoError(t, err)
						require.Equal(t, expectedResult, actual)

						return instanceMock
					})
				}),
		).
		WithSetup(func(host vmhost.VMHost, world *worldmock.MockWorld) {
			world.ProvidedBlockchainHook = &contextmock.BlockchainHookStub{
				QueryDRWANativeGovernanceCalled: func(queryType uint32, key []byte) ([]byte, error) {
					require.Equal(t, uint32(0), queryType)
					require.Equal(t, queryKey, key)
					return expectedResult, nil
				},
			}
		}).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("testFunction").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.Ok()
		})
	require.NoError(t, err)
}

func TestManagedDRWASyncMirror_MalformedPayload_Reverts(t *testing.T) {
	t.Parallel()

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("testFunction", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						payloadHandle := managedTypes.NewManagedBufferFromBytes(buildMalformedDRWASyncPayload())
						result := hooks.ManagedDRWASyncMirror(payloadHandle)
						require.Equal(t, int32(1), result)

						return instanceMock
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("testFunction").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				ReturnMessage(vmhost.ErrInvalidArgument.Error())
		})
	require.NoError(t, err)
}

func TestManagedDRWASyncMirror_ZeroHashPrefix_Reverts(t *testing.T) {
	t.Parallel()

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("testFunction", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						payloadHandle := managedTypes.NewManagedBufferFromBytes(buildZeroHashDRWASyncPayload())
						result := hooks.ManagedDRWASyncMirror(payloadHandle)
						require.Equal(t, int32(1), result)

						return instanceMock
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("testFunction").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				ReturnMessage(vmhost.ErrInvalidArgument.Error())
		})
	require.NoError(t, err)
}

func TestManagedDRWASyncMirror_TooManyOperations_Reverts(t *testing.T) {
	t.Parallel()

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("testFunction", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						payloadHandle := managedTypes.NewManagedBufferFromBytes(buildValidDRWASyncPayload(drwaSyncOpsCapForTests + 1))
						result := hooks.ManagedDRWASyncMirror(payloadHandle)
						require.Equal(t, int32(1), result)

						return instanceMock
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("testFunction").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				ReturnMessage(vmhost.ErrInvalidArgument.Error())
		})
	require.NoError(t, err)
}

func TestManagedDRWASyncMirror_OversizedField_Reverts(t *testing.T) {
	t.Parallel()

	_, err := test.BuildMockInstanceCallTest(t).
		WithContracts(
			test.CreateMockContract(test.ParentAddress).
				WithBalance(1000).
				WithMethods(func(instanceMock *contextmock.InstanceMock, config interface{}) {
					instanceMock.AddMockMethod("testFunction", func() *contextmock.InstanceMock {
						host := instanceMock.Host
						managedTypes := host.ManagedTypes()
						hooks := vmhooks.NewVMHooksImpl(host)

						payloadHandle := managedTypes.NewManagedBufferFromBytes(buildOversizedFieldDRWASyncPayload())
						result := hooks.ManagedDRWASyncMirror(payloadHandle)
						require.Equal(t, int32(1), result)

						return instanceMock
					})
				}),
		).
		WithInput(test.CreateTestContractCallInputBuilder().
			WithRecipientAddr(test.ParentAddress).
			WithGasProvided(100000).
			WithFunction("testFunction").
			Build()).
		AndAssertResults(func(world *worldmock.MockWorld, verify *test.VMOutputVerifier) {
			verify.ExecutionFailed().
				ReturnMessage(vmhost.ErrInvalidArgument.Error())
		})
	require.NoError(t, err)
}

func buildValidDRWASyncPayload(numOps int) []byte {
	payload := make([]byte, 33)
	for i := 0; i < 32; i++ {
		payload[i] = byte(i + 1)
	}
	for i := 0; i < numOps; i++ {
		payload = append(payload, byte(i%3))
		payload = appendLenPrefixed(payload, []byte("TOKEN-123"))
		payload = appendLenPrefixed(payload, []byte("erd1holder"))
		payload = append(payload, 0, 0, 0, 0, 0, 0, 0, byte(i+1))
		payload = appendLenPrefixed(payload, []byte("body"))
	}

	return payload
}

func buildZeroHashDRWASyncPayload() []byte {
	payload := make([]byte, 33)
	payload[32] = 0x01
	payload = appendLenPrefixed(payload, []byte("TOKEN-123"))
	payload = appendLenPrefixed(payload, []byte("erd1holder"))
	payload = append(payload, 0, 0, 0, 0, 0, 0, 0, 1)
	payload = appendLenPrefixed(payload, []byte("body"))
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
	payload = appendLenPrefixed(payload, bytesOfLen(drwaSyncFieldLenCapForTests+1))
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
