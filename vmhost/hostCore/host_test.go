package hostCore

import (
	"math"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	wasmConfig "github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	mockcontext "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/stretchr/testify/require"
)

func TestVmHost_IsAllowedToExecute_DRWASyncDynamicAuthorization(t *testing.T) {
	blockchainHook := &mockcontext.BlockchainHookStub{
		IsAuthorizedDRWASyncCallerCalled: func(callerAddress []byte) bool {
			return string(callerAddress) == "drwa-contract-address"
		},
	}
	bfc := builtInFunctions.NewBuiltInFunctionContainer()
	epochNotifier := &mock.EpochNotifierStub{}
	epochsHandler := &worldmock.EnableEpochsHandlerStub{}
	esdtTransferParser, err := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	require.NoError(t, err)

	host, err := NewVMHost(blockchainHook, &vmhost.VMHostParameters{
		VMType:                    []byte("vmType"),
		ESDTTransferParser:        esdtTransferParser,
		BuiltInFuncContainer:      bfc,
		GasSchedule:               wasmConfig.MakeGasMapForTests(),
		ProtectedKeyPrefix:        []byte("ELROND"),
		EpochNotifier:             epochNotifier,
		EnableEpochsHandler:       epochsHandler,
		Hasher:                    worldmock.DefaultHasher,
		MapOpcodeAddressIsAllowed: map[string]map[string]struct{}{},
	})
	require.NoError(t, err)
	require.NotNil(t, host)

	runtimeContext := &mockcontext.RuntimeContextMock{SCAddress: []byte("drwa-contract-address")}
	host.SetRuntimeContext(runtimeContext)

	require.True(t, host.IsAllowedToExecute("managedDRWASyncMirror"))

	runtimeContext.SCAddress = []byte("unauthorized-address")
	require.False(t, host.IsAllowedToExecute("managedDRWASyncMirror"))
}

func TestNewVMHost(t *testing.T) {
	blockchainHook := worldmock.NewMockWorld()
	bfc := builtInFunctions.NewBuiltInFunctionContainer()
	epochNotifier := &mock.EpochNotifierStub{}
	epochsHandler := &worldmock.EnableEpochsHandlerStub{}
	vmType := []byte("vmType")
	esdtTransferParser, err := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	require.Nil(t, err)

	makeHostParameters := func() *vmhost.VMHostParameters {
		return &vmhost.VMHostParameters{
			VMType:                    vmType,
			ESDTTransferParser:        esdtTransferParser,
			BuiltInFuncContainer:      bfc,
			EpochNotifier:             epochNotifier,
			EnableEpochsHandler:       epochsHandler,
			Hasher:                    worldmock.DefaultHasher,
			MapOpcodeAddressIsAllowed: map[string]map[string]struct{}{},
		}
	}

	t.Run("NilBlockchainHook", func(t *testing.T) {
		host, err := NewVMHost(nil, makeHostParameters())
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilBlockChainHook)
	})
	t.Run("NilHostParameters", func(t *testing.T) {
		host, err := NewVMHost(blockchainHook, nil)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilHostParameters)
	})
	t.Run("NilESDTTransferParser", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.ESDTTransferParser = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilESDTTransferParser)
	})
	t.Run("NilBuiltInFunctionsContainer", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.BuiltInFuncContainer = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilBuiltInFunctionsContainer)
	})
	t.Run("NilEpochNotifier", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EpochNotifier = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilEpochNotifier)
	})
	t.Run("NilEnableEpochsHandler", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EnableEpochsHandler = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilEnableEpochsHandler)
	})
	t.Run("InvalidEnableEpochsHandler", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.EnableEpochsHandler = &worldmock.EnableEpochsHandlerStub{
			IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
				return false
			},
		}
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, core.ErrInvalidEnableEpochsHandler)
	})
	t.Run("NilHasher", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.Hasher = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilHasher)
	})
	t.Run("NilVMType", func(t *testing.T) {
		hostParameters := makeHostParameters()
		hostParameters.VMType = nil
		host, err := NewVMHost(blockchainHook, hostParameters)
		require.Nil(t, host)
		require.ErrorIs(t, err, vmhost.ErrNilVMType)
	})
}

func TestValidateVMInput(t *testing.T) {
	vmInput := &vmcommon.VMInput{
		GasProvided: 0,
	}

	vmInput.GasProvided = math.MaxUint64
	err := validateVMInput(vmInput)
	require.ErrorIs(t, err, vmhost.ErrInvalidGasProvided)

	vmInput.GasProvided = math.MaxInt64
	err = validateVMInput(vmInput)
	require.Nil(t, err)
}
