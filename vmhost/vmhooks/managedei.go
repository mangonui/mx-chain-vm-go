package vmhooks

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/math"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const (
	managedSCAddressName                         = "managedSCAddress"
	managedOwnerAddressName                      = "managedOwnerAddress"
	managedCallerName                            = "managedCaller"
	managedSignalErrorName                       = "managedSignalError"
	managedWriteLogName                          = "managedWriteLog"
	managedMultiTransferESDTNFTExecuteName       = "managedMultiTransferESDTNFTExecute"
	managedTransferValueExecuteName              = "managedTransferValueExecute"
	managedExecuteOnDestContextName              = "managedExecuteOnDestContext"
	managedExecuteOnDestContextWithReturnName    = "managedExecuteOnDestContextWithReturn"
	managedExecuteOnDestContextByCallerName      = "managedExecuteOnDestContextByCaller"
	managedExecuteOnSameContextName              = "managedExecuteOnSameContext"
	managedExecuteReadOnlyName                   = "managedExecuteReadOnly"
	managedCreateContractName                    = "managedCreateContract"
	managedDeployFromSourceContractName          = "managedDeployFromSourceContract"
	managedUpgradeContractName                   = "managedUpgradeContract"
	managedUpgradeFromSourceContractName         = "managedUpgradeFromSourceContract"
	managedAsyncCallName                         = "managedAsyncCall"
	managedCreateAsyncCallName                   = "managedCreateAsyncCall"
	managedGetCallbackClosure                    = "managedGetCallbackClosure"
	managedGetMultiESDTCallValueName             = "managedGetMultiESDTCallValue"
	managedGetAllTransfersCallValue              = "managedGetAllTransfersCallValue"
	managedGetESDTBalanceName                    = "managedGetESDTBalance"
	managedGetESDTTokenDataName                  = "managedGetESDTTokenData"
	managedGetESDTTokenTypeName                  = "managedGetESDTTokenType"
	managedGetReturnDataName                     = "managedGetReturnData"
	managedGetPrevBlockRandomSeedName            = "managedGetPrevBlockRandomSeed"
	managedGetBlockRandomSeedName                = "managedGetBlockRandomSeed"
	managedGetStateRootHashName                  = "managedGetStateRootHash"
	managedGetOriginalTxHashName                 = "managedGetOriginalTxHash"
	managedIsESDTFrozenName                      = "managedIsESDTFrozen"
	managedIsESDTLimitedTransferName             = "managedIsESDTLimitedTransfer"
	managedIsESDTPausedName                      = "managedIsESDTPaused"
	managedBufferToHexName                       = "managedBufferToHex"
	managedGetCodeMetadataName                   = "managedGetCodeMetadata"
	managedGetCodeHashName                       = "managedGetCodeHash"
	managedIsBuiltinFunction                     = "managedIsBuiltinFunction"
	managedMultiTransferESDTNFTExecuteByUser     = "managedMultiTransferESDTNFTExecuteByUser"
	managedMultiTransferESDTNFTExecuteWithReturn = "managedMultiTransferESDTNFTExecuteByWithReturn"
	managedDRWASyncMirrorName                    = "managedDRWASyncMirror"
	managedDRWANativeGovernanceQueryName         = "managedDRWANativeGovernanceQuery"
	maxDRWASyncPayloadBytes                      = 1 << 20
	maxDRWASyncFieldLen                          = 64 * 1024
	maxDRWASyncOps                               = 256
)

const EGLDTokenName = "EGLD-000000" // TODO: maybe move to core?

// ManagedSCAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedSCAddress(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetSCAddress
	err := metering.UseGasBoundedAndAddTracedGas(managedSCAddressName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	scAddress := runtime.GetContextAddress()

	managedType.SetBytes(destinationHandle, scAddress)
}

// ManagedOwnerAddress VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedOwnerAddress(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOwnerAddress
	err := metering.UseGasBoundedAndAddTracedGas(managedOwnerAddressName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	owner, err := blockchain.GetOwnerAddress()
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(destinationHandle, owner)
}

// ManagedCaller VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedCaller(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	err := metering.UseGasBoundedAndAddTracedGas(managedCallerName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	caller := runtime.GetVMInput().CallerAddr
	managedType.SetBytes(destinationHandle, caller)
}

// ManagedGetOriginalCallerAddr VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetOriginalCallerAddr(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	err := metering.UseGasBoundedAndAddTracedGas(managedCallerName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	caller := runtime.GetVMInput().OriginalCallerAddr
	managedType.SetBytes(destinationHandle, caller)
}

// ManagedGetRelayerAddr VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetRelayerAddr(destinationHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	err := metering.UseGasBoundedAndAddTracedGas(managedCallerName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	caller := runtime.GetVMInput().RelayerAddr
	managedType.SetBytes(destinationHandle, caller)
}

// ManagedSignalError VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedSignalError(errHandle int32) {
	managedType := context.GetManagedTypesContext()
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	metering.StartGasTracing(managedSignalErrorName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.SignalError
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		context.FailExecution(err)
		return
	}

	errBytes, err := managedType.GetBytes(errHandle)
	if err != nil {
		context.FailExecution(err)
		return
	}

	err = managedType.ConsumeGasForBytes(errBytes)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		context.FailExecution(err)
		return
	}

	gasToUse = metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(len(errBytes))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		context.FailExecution(err)
		return
	}

	runtime.SignalUserError(string(errBytes))
}

// ManagedWriteLog VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedWriteLog(
	topicsHandle int32,
	dataHandle int32,
) {
	runtime := context.GetRuntimeContext()
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()
	metering.StartGasTracing(managedWriteLogName)

	topics, sumOfTopicByteLengths, err := managedType.ReadManagedVecOfManagedBuffers(topicsHandle)
	if err != nil {
		context.FailExecution(err)
		return
	}

	dataBytes, err := managedType.GetBytes(dataHandle)
	if err != nil {
		context.FailExecution(err)
		return
	}

	err = managedType.ConsumeGasForBytes(dataBytes)
	if err != nil {
		context.FailExecution(err)
		return
	}

	dataByteLen := uint64(len(dataBytes))

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		sumOfTopicByteLengths+dataByteLen)
	gasToUse = math.AddUint64(gasToUse, gasForData)
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	output.WriteLog(runtime.GetContextAddress(), topics, [][]byte{dataBytes})
}

// ManagedGetOriginalTxHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetOriginalTxHash(resultHandle int32) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOriginalTxHash
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(resultHandle, runtime.GetOriginalTxHash())
}

// ManagedGetStateRootHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetStateRootHash(resultHandle int32) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetStateRootHash
	err := metering.UseGasBoundedAndAddTracedGas(managedGetStateRootHashName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(resultHandle, blockchain.GetStateRootHash())
}

// ManagedGetBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetBlockRandomSeed(resultHandle int32) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	err := metering.UseGasBoundedAndAddTracedGas(managedGetBlockRandomSeedName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(resultHandle, blockchain.CurrentRandomSeed())
}

// ManagedGetPrevBlockRandomSeed VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetPrevBlockRandomSeed(resultHandle int32) {
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	err := metering.UseGasBoundedAndAddTracedGas(managedGetPrevBlockRandomSeedName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(resultHandle, blockchain.LastRandomSeed())
}

// ManagedGetReturnData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetReturnData(resultID int32, resultHandle int32) {
	output := context.GetOutputContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnData
	err := metering.UseGasBoundedAndAddTracedGas(managedGetReturnDataName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		context.FailExecution(vmhost.ErrArgOutOfRange)
		return
	}

	managedType.SetBytes(resultHandle, returnData[resultID])
}

// ManagedGetMultiESDTCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetMultiESDTCallValue(multiCallValueHandle int32) {
	runtime := context.GetRuntimeContext()
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(managedGetMultiESDTCallValueName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	esdtTransfers := runtime.GetVMInput().ESDTTransfers
	multiCallBytes := writeESDTTransfersToBytes(managedType, esdtTransfers)
	err = managedType.ConsumeGasForBytes(multiCallBytes)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(multiCallValueHandle, multiCallBytes)
}

// ManagedGetAllTransfersCallValue VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetAllTransfersCallValue(transferCallValuesListHandle int32) {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()

	allTransfers, err := ManagedGetAllTransfersCallValueTyped(host)
	if err != nil {
		context.FailExecution(err)
		return
	}

	allTransfersBytes := writeESDTTransfersToBytes(managedType, allTransfers)
	err = managedType.ConsumeGasForBytes(allTransfersBytes)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(transferCallValuesListHandle, allTransfersBytes)
}

// ManagedGetAllTransfersCallValueTyped returns a combined list of all transfers (ESDT and EGLD)
func ManagedGetAllTransfersCallValueTyped(
	host vmhost.VMHost,
) ([]*vmcommon.ESDTTransfer, error) {
	runtime := host.Runtime()
	metering := host.Metering()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(managedGetAllTransfersCallValue, gasToUse)
	if err != nil {
		return nil, err
	}

	input := runtime.GetVMInput()
	egldCallValue := input.CallValue
	hasCallValue := egldCallValue.Sign() > 0

	if hasCallValue {
		return []*vmcommon.ESDTTransfer{
			{
				ESDTValue:      egldCallValue,
				ESDTTokenName:  []byte(EGLDTokenName),
				ESDTTokenType:  uint32(core.Fungible),
				ESDTTokenNonce: 0,
			},
		}, nil
	}

	return input.ESDTTransfers, nil
}

// ManagedGetBackTransfers VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetBackTransfers(esdtTransfersValueHandle int32, egldValueHandle int32) {
	metering := context.GetMeteringContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	err := metering.UseGasBoundedAndAddTracedGas(managedGetMultiESDTCallValueName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	esdtTransfers, transferValue := managedType.GetBackTransfers()
	multiCallBytes := writeESDTTransfersToBytes(managedType, esdtTransfers)
	err = managedType.ConsumeGasForBytes(multiCallBytes)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(esdtTransfersValueHandle, multiCallBytes)
	egldValue := managedType.GetBigIntOrCreate(egldValueHandle)
	egldValue.SetBytes(transferValue.Bytes())
}

// ManagedGetESDTBalance VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetESDTBalance(addressHandle int32, tokenIDHandle int32, nonce int64, valueHandle int32) {
	metering := context.GetMeteringContext()
	blockchain := context.GetBlockchainContext()
	managedType := context.GetManagedTypesContext()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(managedGetESDTBalanceName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		context.FailExecution(vmhost.ErrArgOutOfRange)
		return
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		context.FailExecution(vmhost.ErrArgOutOfRange)
		return
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		context.FailExecution(vmhost.ErrArgOutOfRange)
		return
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(esdtToken.Value)
}

// ManagedGetESDTTokenData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetESDTTokenData(
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) {
	host := context.GetVMHost()
	ManagedGetESDTTokenDataWithHost(
		host,
		addressHandle,
		tokenIDHandle,
		nonce,
		valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle)

}

func ManagedGetESDTTokenDataWithHost(
	host vmhost.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedGetESDTTokenDataName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		FailExecution(host, err)
		return
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		FailExecution(host, vmhost.ErrArgOutOfRange)
		return
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		FailExecution(host, vmhost.ErrArgOutOfRange)
		return
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		FailExecution(host, vmhost.ErrArgOutOfRange)
		return
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(esdtToken.Value)

	managedType.SetBytes(propertiesHandle, esdtToken.Properties)
	if esdtToken.TokenMetaData != nil {
		managedType.SetBytes(hashHandle, esdtToken.TokenMetaData.Hash)
		err = managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Hash)
		if err != nil {
			FailExecution(host, err)
			return
		}
		managedType.SetBytes(nameHandle, esdtToken.TokenMetaData.Name)
		err = managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Name)
		if err != nil {
			FailExecution(host, err)
			return
		}
		managedType.SetBytes(attributesHandle, esdtToken.TokenMetaData.Attributes)
		err = managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Attributes)
		if err != nil {
			FailExecution(host, err)
			return
		}
		managedType.SetBytes(creatorHandle, esdtToken.TokenMetaData.Creator)
		err = managedType.ConsumeGasForBytes(esdtToken.TokenMetaData.Creator)
		if err != nil {
			FailExecution(host, err)
			return
		}
		royalties := managedType.GetBigIntOrCreate(royaltiesHandle)
		royalties.SetUint64(uint64(esdtToken.TokenMetaData.Royalties))

		err = managedType.WriteManagedVecOfManagedBuffers(esdtToken.TokenMetaData.URIs, urisHandle)
		if err != nil {
			FailExecution(host, err)
			return
		}
	}

}

// ManagedGetESDTTokenData VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetESDTTokenType(
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	typeHandle int32) {
	host := context.GetVMHost()
	ManagedGetESDTTokenTypeWithHost(
		host,
		addressHandle,
		tokenIDHandle,
		nonce,
		typeHandle)
}

// ManagedGetESDTTokenTypeWithHost implements the ManagedGetESDTTokenType VMHook logic.
func ManagedGetESDTTokenTypeWithHost(
	host vmhost.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64,
	typeHandle int32) {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedGetESDTTokenTypeName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		FailExecution(host, err)
		return
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		FailExecution(host, vmhost.ErrArgOutOfRange)
		return
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		FailExecution(host, vmhost.ErrArgOutOfRange)
		return
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		FailExecution(host, vmhost.ErrArgOutOfRange)
		return
	}

	esdtType := managedType.GetBigIntOrCreate(typeHandle)
	esdtType.SetUint64(uint64(esdtToken.Type))
}

// ManagedAsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedAsyncCall(
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32) {
	host := context.GetVMHost()
	ManagedAsyncCallWithHost(
		host,
		destHandle,
		valueHandle,
		functionHandle,
		argumentsHandle)
}

func ManagedAsyncCallWithHost(
	host vmhost.VMHost,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32) {
	runtime := host.Runtime()
	async := host.Async()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedAsyncCallName)

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.BaseOpsAPICost.AsyncCallStep
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return
	}

	vmInput, err := readDestinationFunctionArguments(host, destHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	data := makeCrossShardCallFromInput(vmInput.function, vmInput.arguments)

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		FailExecution(host, vmhost.ErrArgOutOfRange)
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return
	}

	err = async.RegisterLegacyAsyncCall(vmInput.destination, []byte(data), value.Bytes())
	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}
	if err != nil {
		FailExecution(host, err)
		return
	}
}

// ManagedCreateAsyncCall VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedCreateAsyncCall(
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	successOffset executor.MemPtr,
	successLength executor.MemLength,
	errorOffset executor.MemPtr,
	errorLength executor.MemLength,
	gas int64,
	extraGasForCallback int64,
	callbackClosureHandle int32,
) int32 {

	host := context.GetVMHost()
	managedType := host.ManagedTypes()

	vmInput, err := readDestinationFunctionArguments(host, destHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	data := makeCrossShardCallFromInput(vmInput.function, vmInput.arguments)

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		context.FailExecution(vmhost.ErrArgOutOfRange)
		return 1
	}

	successFunc, err := context.MemLoad(successOffset, successLength)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	errorFunc, err := context.MemLoad(errorOffset, errorLength)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	callbackClosure, err := managedType.GetBytes(callbackClosureHandle)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	return CreateAsyncCallWithTypedArgs(host,
		vmInput.destination,
		value.Bytes(),
		[]byte(data),
		successFunc,
		errorFunc,
		gas,
		extraGasForCallback,
		callbackClosure)
}

// ManagedGetCallbackClosure VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetCallbackClosure(
	callbackClosureHandle int32,
) {
	host := context.GetVMHost()
	GetCallbackClosureWithHost(host, callbackClosureHandle)
}

func GetCallbackClosureWithHost(
	host vmhost.VMHost,
	callbackClosureHandle int32,
) {
	async := host.Async()
	metering := host.Metering()
	managedTypes := host.ManagedTypes()

	metering.StartGasTracing(managedGetCallbackClosure)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallbackClosure
	err := metering.UseGasBounded(gasToUse)
	if err != nil {
		FailExecution(host, err)
		return
	}

	callbackClosure, err := async.GetCallbackClosure()
	if err != nil {
		FailExecution(host, err)
		return
	}

	managedTypes.SetBytes(callbackClosureHandle, callbackClosure)
}

// ManagedUpgradeFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedUpgradeFromSourceContract(
	destHandle int32,
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		context.FailExecution(err)
		return
	}

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	sourceContractAddress, err := managedType.GetBytes(addressHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	lenReturnData := len(host.Output().ReturnData())

	UpgradeFromSourceContractWithTypedArgs(
		host,
		sourceContractAddress,
		vmInput.destination,
		vmInput.value.Bytes(),
		vmInput.arguments,
		gas,
		codeMetadata,
	)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return
	}
}

// ManagedUpgradeContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedUpgradeContract(
	destHandle int32,
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		context.FailExecution(err)
		return
	}

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	code, err := managedType.GetBytes(codeHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	if err != nil {
		FailExecution(host, err)
		return
	}

	lenReturnData := len(host.Output().ReturnData())

	upgradeContract(host, vmInput.destination, code, codeMetadata, vmInput.value.Bytes(), vmInput.arguments, gas)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return
	}
}

// ManagedDeleteContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedDeleteContract(
	destHandle int32,
	gasLimit int64,
	argumentsHandle int32,
) {
	host := context.GetVMHost()
	ManagedDeleteContractWithHost(
		host,
		destHandle,
		gasLimit,
		argumentsHandle,
	)
}

func ManagedDeleteContractWithHost(
	host vmhost.VMHost,
	destHandle int32,
	gasLimit int64,
	argumentsHandle int32,
) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(deleteContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return
	}

	calledSCAddress, err := managedType.GetBytes(destHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	data, _, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	deleteContract(
		host,
		calledSCAddress,
		data,
		gasLimit,
	)
}

// ManagedDeployFromSourceContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedDeployFromSourceContract(
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedDeployFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		context.FailExecution(err)
		return -1
	}

	vmInput, err := readDestinationValueArguments(host, addressHandle, valueHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	lenReturnData := len(host.Output().ReturnData())

	newAddress, err := DeployFromSourceContractWithTypedArgs(
		host,
		vmInput.destination,
		codeMetadata,
		vmInput.value,
		vmInput.arguments,
		gas,
	)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return 1
	}

	return 0
}

// ManagedCreateContract VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedCreateContract(
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedCreateContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	err := metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		context.FailExecution(err)
		return -1
	}

	sender := runtime.GetContextAddress()
	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	data, actualLen, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLen)
	err = metering.UseGasBounded(gasToUse)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		context.FailExecution(err)
		return -1
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	code, err := managedType.GetBytes(codeHandle)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	lenReturnData := len(host.Output().ReturnData())
	newAddress, err := createContract(sender, data, value, gas, code, codeMetadata, host, CreateContract)
	if err != nil {
		FailExecution(host, err)
		return 1
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && runtime.UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return 1
	}

	return 0
}

func setReturnDataIfExists(
	host vmhost.VMHost,
	oldLen int,
	resultHandle int32,
) error {
	returnData := host.Output().ReturnData()
	if len(returnData) > oldLen {
		return host.ManagedTypes().WriteManagedVecOfManagedBuffers(returnData[oldLen:], resultHandle)
	}

	host.ManagedTypes().SetBytes(resultHandle, make([]byte, 0))
	return nil
}

// ManagedExecuteReadOnly VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedExecuteReadOnly(
	gas int64,
	addressHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteReadOnlyName)

	vmInput, err := readDestinationFunctionArguments(host, addressHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteReadOnlyWithTypedArguments(
		host,
		gas,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && host.Runtime().UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return -1
	}

	return returnVal
}

// ManagedExecuteOnSameContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedExecuteOnSameContext(
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnSameContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnSameContextWithTypedArgs(
		host,
		gas,
		vmInput.value,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && host.Runtime().UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return -1
	}

	return returnVal
}

// ManagedExecuteOnDestContext VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedExecuteOnDestContext(
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnDestContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnDestContextWithTypedArgs(
		host,
		gas,
		vmInput.value,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
		true,
	)
	err = setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil && host.Runtime().UseGasBoundedShouldFailExecution() {
		FailExecution(host, err)
		return -1
	}

	return returnVal
}

// ManagedExecuteOnDestContextWithErrorReturn VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedExecuteOnDestContextWithErrorReturn(
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := context.GetVMHost()

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	return ManagedExecuteOnDestContextWithErrorReturnWithHost(host, gas, vmInput.value, vmInput.function, vmInput.destination, vmInput.arguments, resultHandle)
}

// ManagedExecuteOnDestContextWithErrorReturnWithHost - execute on dest context and return error instead of failing execution
func ManagedExecuteOnDestContextWithErrorReturnWithHost(
	host vmhost.VMHost,
	gas int64,
	value *big.Int,
	function string,
	destination []byte,
	arguments [][]byte,
	resultHandle int32,
) int32 {
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnDestContextWithReturnName)

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnDestContextWithTypedArgs(
		host,
		gas,
		value,
		[]byte(function),
		destination,
		arguments,
		false,
	)

	err := setReturnDataIfExists(host, lenReturnData, resultHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	return returnVal
}

// ManagedMultiTransferESDTNFTExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMultiTransferESDTNFTExecute(
	dstHandle int32,
	tokenTransfersHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(managedMultiTransferESDTNFTExecuteName)

	vmInput, err := readDestinationFunctionArguments(host, dstHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	transfers, err := readESDTTransfers(managedType, runtime, tokenTransfersHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	return TransferESDTNFTExecuteWithTypedArgs(
		host,
		vmInput.destination,
		transfers,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	)
}

// ManagedMultiTransferESDTNFTExecuteWithReturn VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMultiTransferESDTNFTExecuteWithReturn(
	dstHandle int32,
	tokenTransfersHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(managedMultiTransferESDTNFTExecuteWithReturn)

	vmInput, err := readDestinationFunctionArguments(host, dstHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	transfers, err := readESDTTransfers(managedType, runtime, tokenTransfersHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	return TransferESDTNFTExecuteWithTypedArgsWithFailure(
		host,
		vmInput.destination,
		transfers,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
		false,
	)
}

// ManagedMultiTransferESDTNFTExecuteByUser VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedMultiTransferESDTNFTExecuteByUser(
	userHandle int32,
	dstHandle int32,
	tokenTransfersHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(managedMultiTransferESDTNFTExecuteByUser)

	if !host.IsAllowedToExecute(managedMultiTransferESDTNFTExecuteByUser) {
		FailExecution(host, vmhost.ErrOpcodeIsNotAllowed)
		return -1
	}

	user, err := managedType.GetBytes(userHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	vmInput, err := readDestinationFunctionArguments(host, dstHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	transfers, err := readESDTTransfers(managedType, runtime, tokenTransfersHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	return TransferESDTNFTExecuteByUserWithTypedArgs(
		host,
		user,
		vmInput.destination,
		transfers,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	)
}

// ManagedTransferValueExecute VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedTransferValueExecute(
	dstHandle int32,
	valueHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := context.GetVMHost()
	metering := host.Metering()
	metering.StartGasTracing(managedTransferValueExecuteName)

	vmInput, err := readDestinationValueFunctionArguments(host, dstHandle, valueHandle, functionHandle, argumentsHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	return TransferValueExecuteWithTypedArgs(
		host,
		vmInput.destination,
		vmInput.value,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	)
}

// ManagedIsESDTFrozen VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedIsESDTFrozen(
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) int32 {
	host := context.GetVMHost()
	return ManagedIsESDTFrozenWithHost(host, addressHandle, tokenIDHandle, nonce)
}

func ManagedIsESDTFrozenWithHost(
	host vmhost.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) int32 {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(managedIsESDTFrozenName, gasToUse)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	esdtToken, err := blockchain.GetESDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	esdtUserData := builtInFunctions.ESDTUserMetadataFromBytes(esdtToken.Properties)
	if esdtUserData.Frozen {
		return 1
	}
	return 0
}

// ManagedIsESDTLimitedTransfer VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedIsESDTLimitedTransfer(tokenIDHandle int32) int32 {
	host := context.GetVMHost()
	return ManagedIsESDTLimitedTransferWithHost(host, tokenIDHandle)
}

func ManagedIsESDTLimitedTransferWithHost(host vmhost.VMHost, tokenIDHandle int32) int32 {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(managedIsESDTLimitedTransferName, gasToUse)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	if blockchain.IsLimitedTransfer(tokenID) {
		return 1
	}

	return 0
}

// ManagedIsESDTPaused VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedIsESDTPaused(tokenIDHandle int32) int32 {
	host := context.GetVMHost()
	return ManagedIsESDTPausedWithHost(host, tokenIDHandle)
}

func ManagedIsESDTPausedWithHost(host vmhost.VMHost, tokenIDHandle int32) int32 {
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	err := metering.UseGasBoundedAndAddTracedGas(managedIsESDTPausedName, gasToUse)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	if blockchain.IsPaused(tokenID) {
		return 1
	}

	return 0
}

// ManagedBufferToHex VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedBufferToHex(sourceHandle int32, destHandle int32) {
	host := context.GetVMHost()
	ManagedBufferToHexWithHost(host, sourceHandle, destHandle)
}

func ManagedBufferToHexWithHost(host vmhost.VMHost, sourceHandle int32, destHandle int32) {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	err := metering.UseGasBoundedAndAddTracedGas(managedBufferToHexName, gasToUse)
	if err != nil {
		FailExecution(host, err)
		return
	}

	mBuff, err := managedType.GetBytes(sourceHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	encoded := hex.EncodeToString(mBuff)
	managedType.SetBytes(destHandle, []byte(encoded))
}

// ManagedGetCodeMetadata VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetCodeMetadata(addressHandle int32, responseHandle int32) {
	host := context.GetVMHost()
	ManagedGetCodeMetadataWithHost(host, addressHandle, responseHandle)
}

func ManagedGetCodeMetadataWithHost(host vmhost.VMHost, addressHandle int32, responseHandle int32) {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCodeMetadata
	err := metering.UseGasBoundedAndAddTracedGas(managedGetCodeMetadataName, gasToUse)
	if err != nil {
		FailExecution(host, err)
		return
	}

	gasToUse = metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	err = metering.UseGasBoundedAndAddTracedGas(managedGetCodeMetadataName, gasToUse)
	if err != nil {
		FailExecution(host, err)
		return
	}

	mBuffAddress, err := managedType.GetBytes(addressHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	contract, err := host.Blockchain().GetUserAccount(mBuffAddress)
	if err != nil || check.IfNil(contract) {
		FailExecution(host, err)
		return
	}

	codeMetadata := contract.GetCodeMetadata()

	managedType.SetBytes(responseHandle, codeMetadata)
}

// ManagedGetCodeHash VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedGetCodeHash(addressHandle int32, codeHashHandle int32) {
	host := context.GetVMHost()
	managedType := host.ManagedTypes()

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		FailExecution(host, err)
		return
	}

	codeHash, err := ManagedGetCodeHashTyped(host, address)
	if err != nil {
		context.FailExecution(err)
		return
	}

	managedType.SetBytes(codeHashHandle, codeHash)
}

// ManagedGetCodeHashWithHost returns the code hash at some address
func ManagedGetCodeHashTyped(
	host vmhost.VMHost,
	address []byte,
) ([]byte, error) {
	metering := host.Metering()
	blockchain := host.Blockchain()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCodeHash
	err := metering.UseGasBoundedAndAddTracedGas(managedGetCodeHashName, gasToUse)
	if err != nil {
		return nil, err
	}

	codeHash := blockchain.GetCodeHash(address)

	return codeHash, nil
}

// ManagedIsBuiltinFunction VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedIsBuiltinFunction(functionNameHandle int32) int32 {
	host := context.GetVMHost()
	return ManagedIsBuiltinFunctionWithHost(host, functionNameHandle)
}

func ManagedIsBuiltinFunctionWithHost(host vmhost.VMHost, functionNameHandle int32) int32 {
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.IsBuiltinFunction
	err := metering.UseGasBoundedAndAddTracedGas(managedIsBuiltinFunction, gasToUse)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	mBuffFunctionName, err := managedType.GetBytes(functionNameHandle)
	if err != nil {
		FailExecution(host, err)
		return -1
	}

	isBuiltinFunction := host.IsBuiltinFunctionName(string(mBuffFunctionName))
	if isBuiltinFunction {
		return 1
	}

	return 0
}

// ManagedDRWASyncMirror VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedDRWASyncMirror(payloadHandle int32) int32 {
	host := context.GetVMHost()
	if !host.IsAllowedToExecute(managedDRWASyncMirrorName) {
		FailExecution(host, vmhost.ErrOpcodeIsNotAllowed)
		return 1
	}

	managedType := context.GetManagedTypesContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()
	runtime := context.GetRuntimeContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetBytes
	err := metering.UseGasBoundedAndAddTracedGas(managedDRWASyncMirrorName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return 1
	}

	payload, err := managedType.GetBytes(payloadHandle)
	if err != nil {
		context.FailExecution(err)
		return 1
	}
	err = managedType.ConsumeGasForBytes(payload)
	if err != nil {
		context.FailExecution(err)
		return 1
	}
	if len(payload) > maxDRWASyncPayloadBytes {
		context.FailExecution(errors.New("DRWA sync payload too large"))
		return 1
	}
	if hasZeroDRWASyncHashPrefix(payload) {
		context.FailExecution(vmhost.ErrInvalidArgument)
		return 1
	}

	// Count operations in the payload to charge per-operation state-write gas.
	// The envelope binary format is: [32-byte hash] || [caller_tag byte] || (op_tag + fields)*
	// We count op-tag bytes (0x00, 0x01, 0x02) after offset 33 as a proxy for
	// operation count.  Each write touches at minimum one trie node, so we
	// charge StorageStore gas per operation.
	opCount, ok := countDRWASyncOperations(payload)
	if !ok {
		context.FailExecution(vmhost.ErrInvalidArgument)
		return 1
	}
	if opCount > 0 {
		// Charge StorePerByte × drwaSyncMinBytesPerOp as the minimum write cost per
		// operation.  This accounts for the trie node update triggered by each
		// token-policy or holder-mirror write; the actual byte-proportional cost is
		// already covered by ConsumeGasForBytes above.
		const drwaSyncMinBytesPerOp = 64
		perOpGas := metering.GasSchedule().BaseOperationCost.StorePerByte * drwaSyncMinBytesPerOp
		writeGas, overflow := safeMulUint64(uint64(opCount), perOpGas)
		if overflow {
			context.FailExecution(errors.New("DRWA sync per-operation gas overflow"))
			return 1
		}
		err = metering.UseGasBoundedAndAddTracedGas(managedDRWASyncMirrorName, writeGas)
		if err != nil {
			context.FailExecution(err)
			return 1
		}
	}

	err = blockchain.ApplyDRWASyncEnvelopeBytes(payload, runtime.GetContextAddress())
	if err != nil {
		context.FailExecution(err)
		return 1
	}

	return 0
}

// ManagedDRWANativeGovernanceQuery VMHooks implementation.
// @autogenerate(VMHooks)
func (context *VMHooksImpl) ManagedDRWANativeGovernanceQuery(queryType int32, keyHandle int32, destHandle int32) int32 {
	managedType := context.GetManagedTypesContext()
	blockchain := context.GetBlockchainContext()
	metering := context.GetMeteringContext()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferGetBytes
	err := metering.UseGasBoundedAndAddTracedGas(managedDRWANativeGovernanceQueryName, gasToUse)
	if err != nil {
		context.FailExecution(err)
		return 1
	}

	if queryType < 0 {
		context.FailExecution(vmhost.ErrInvalidArgument)
		return 1
	}

	key, err := managedType.GetBytes(keyHandle)
	if err != nil {
		context.FailExecution(err)
		return 1
	}
	err = managedType.ConsumeGasForBytes(key)
	if err != nil {
		context.FailExecution(err)
		return 1
	}
	if len(key) == 0 || len(key) > maxDRWASyncFieldLen {
		context.FailExecution(vmhost.ErrInvalidArgument)
		return 1
	}

	result, err := blockchain.QueryDRWANativeGovernance(uint32(queryType), key)
	if err != nil {
		context.FailExecution(err)
		return 1
	}
	if len(result) > maxDRWASyncFieldLen {
		context.FailExecution(errors.New("DRWA native governance query result too large"))
		return 1
	}
	err = managedType.ConsumeGasForBytes(result)
	if err != nil {
		context.FailExecution(err)
		return 1
	}

	managedType.SetBytes(destHandle, result)
	return 0
}

func hasZeroDRWASyncHashPrefix(payload []byte) bool {
	const hashLen = 32
	if len(payload) < hashLen {
		return false
	}

	for i := 0; i < hashLen; i++ {
		if payload[i] != 0 {
			return false
		}
	}

	return true
}

// countDRWASyncOperations counts the number of sync operations encoded in a
// binary DRWA sync payload ([32-byte hash] || [caller_tag] || ops...).
// It walks the canonical serialization format:
//
//	[op_tag u8] [token_id: u32-len-prefixed] [holder: u32-len-prefixed] [version u64] [body: u32-len-prefixed]
func countDRWASyncOperations(payload []byte) (int, bool) {
	const hashLen = 32
	const callerTagLen = 1
	const headerLen = hashLen + callerTagLen

	if len(payload) <= headerLen {
		return 0, true
	}

	data := payload[headerLen:]
	count := 0
	for len(data) > 0 {
		if count >= maxDRWASyncOps {
			return 0, false
		}
		// op_tag (1 byte)
		if len(data) < 1 {
			return 0, false
		}
		data = data[1:]

		// token_id: 4-byte length prefix + body
		var ok bool
		data, ok = skipLenPrefixed(data)
		if !ok {
			return 0, false
		}
		// holder: 4-byte length prefix + body
		data, ok = skipLenPrefixed(data)
		if !ok {
			return 0, false
		}
		// version: 8 bytes
		if len(data) < 8 {
			return 0, false
		}
		data = data[8:]
		// body: 4-byte length prefix + body
		data, ok = skipLenPrefixed(data)
		if !ok {
			return 0, false
		}

		count++
	}

	return count, true
}

func skipLenPrefixed(data []byte) ([]byte, bool) {
	if len(data) < 4 {
		return data, false
	}
	length := int(data[0])<<24 | int(data[1])<<16 | int(data[2])<<8 | int(data[3])
	if length > maxDRWASyncFieldLen {
		return data, false
	}
	data = data[4:]
	if len(data) < length {
		return data, false
	}
	return data[length:], true
}

func safeMulUint64(a, b uint64) (uint64, bool) {
	if a == 0 || b == 0 {
		return 0, false
	}
	result := a * b
	if result/a != b {
		return 0, true
	}
	return result, false
}
