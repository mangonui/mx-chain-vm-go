package hostCore

import (
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

func TestNewVMHost_ExecutionTimeoutClampedToMaximum(t *testing.T) {
	blockchainHook := worldmock.NewMockWorld()
	bfc := builtInFunctions.NewBuiltInFunctionContainer()
	epochNotifier := &mock.EpochNotifierStub{}
	epochsHandler := &worldmock.EnableEpochsHandlerStub{
		IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
			return true
		},
	}
	esdtTransferParser, err := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)
	require.NoError(t, err)

	host, err := NewVMHost(blockchainHook, &vmhost.VMHostParameters{
		VMType:                              []byte("vmType"),
		ESDTTransferParser:                  esdtTransferParser,
		BuiltInFuncContainer:                bfc,
		EpochNotifier:                       epochNotifier,
		EnableEpochsHandler:                 epochsHandler,
		GasSchedule:                         config.MakeGasMapForTests(),
		ProtectedKeyPrefix:                  []byte("ELROND"),
		Hasher:                              worldmock.DefaultHasher,
		MapOpcodeAddressIsAllowed:           map[string]map[string]struct{}{},
		TimeOutForSCExecutionInMilliseconds: uint32((maxExecutionTimeout + 5*time.Second) / time.Millisecond),
	})
	require.NoError(t, err)

	concreteHost, ok := host.(*vmHost)
	require.True(t, ok)
	require.Equal(t, maxExecutionTimeout, concreteHost.executionTimeout)
}
