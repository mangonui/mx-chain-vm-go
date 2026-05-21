# ISSUE-079 Saturating Arithmetic Audit

Status: fixed-verified audit, no behavior change.

This file records the production call-site audit for `math.AddUint64` and
`math.MulUint64`. The helpers intentionally keep their historical
saturating-to-`MaxUint64` behavior. Changing any of these call sites to
`AddUint64WithErr` or `MulUint64WithErr` can change active VM gas/accounting
semantics and needs an epoch-flag/replay-determinism decision.

## Policy

- Keep `AddUint64` / `MulUint64` only where saturation is an acceptable
  fail-closed or historically-established behavior.
- Use `AddUint64WithErr` / `MulUint64WithErr` for new code when the caller can
  explicitly fail the operation without changing existing consensus behavior.
- Any future behavior change in active VM gas accounting must be epoch-gated.

## Classification Summary

| Class | Meaning | Decision |
| --- | --- | --- |
| `gas-charge` | Computes gas to charge before/inside VM hooks. Saturating to `MaxUint64` forces bounded gas failure or maximum charge. | Keep saturating helper. |
| `gas-lock` | Computes async gas lock/minimum gas. Saturating to `MaxUint64` prevents under-locking. | Keep saturating helper. |
| `gas-accounting` | Accumulates already-consumed/remaining/refunded gas. Changing behavior is replay-visible. | Keep saturating helper; Warn log remains the visibility mechanism. |
| `data-size-accounting` | Computes encoded async/networking byte lengths. Overflow is practically unreachable but saturation avoids wraparound. | Keep saturating helper. |

## Audited Production Call Sites

### `vmhost/asyncCall.go`

- `77` `AddUint64(ac.GasLimit, ac.GasLocked)` — `gas-accounting`.
  Historical async call total gas shape; saturation avoids wraparound and is
  replay-visible if changed.

### `vmhost/contexts/async.go`

- `241` `AddUint64(context.gasAccumulated, prevState.gasAccumulated)` —
  `gas-accounting`.
- `357` `AddUint64(gas, metering.ComputeExtraGasLockedForAsync())` —
  `gas-lock`.
- `519` `AddUint64(call.GasLocked, metering.ComputeExtraGasLockedForAsync())`
  — `gas-lock`.
- `996` `AddUint64(dataLength, separator)` — `data-size-accounting`.
- `997` `MulUint64(uint64(len(argument)), hexSize)` —
  `data-size-accounting`.
- `998` `AddUint64(dataLength, encodedArgumentLength)` —
  `data-size-accounting`.
- `1005` `AddUint64(context.gasAccumulated, gas)` — `gas-accounting`.

### `vmhost/contexts/asyncLocal.go`

- `358` `AddUint64(vmOutput.GasRemaining, asyncCall.GetGasLocked())` —
  `gas-accounting`.
- `362` `MulUint64(copyPerByte, uint64(dataLength))` — `gas-charge`.
- `363` `AddUint64(gasToUse, gas)` — `gas-charge`.

### `vmhost/contexts/managedType.go`

- `268` `MulUint64(uint64(byteLen), CopyPerByteForTooBig)` — `gas-charge`.
- `281` `MulUint64(uint64(len(bytes)), DataCopyPerByte)` — `gas-charge`.

### `vmhost/contexts/metering.go`

- `168` `AddUint64(context.gasUsedByAccounts[address], gas)` —
  `gas-accounting`.
- `199` `AddUint64(account.GasUsed, context.GetGasProvided())` —
  `gas-accounting`.
- `248` `AddUint64(gasUsed, vmOutput.GasRemaining)` — `gas-accounting`.
- `266` `AddUint64(gasUsed, outputAccount.GasUsed)` — `gas-accounting`.
- `267` `AddUint64(gasUsed, gasTransferred)` — `gas-accounting`.
- `284` `AddUint64(gasUsedAndTransferred, gasUsed)` — `gas-accounting`.
- `285` `AddUint64(gasUsedAndTransferred, gasTransferred)` —
  `gas-accounting`.
- `294` `AddUint64(gasUsed, outputTransfer.GasLimit)` — `gas-accounting`.
- `295` `AddUint64(gasUsed, outputTransfer.GasLocked)` — `gas-accounting`.
- `328` `AddUint64(input.GasProvided, input.GasLocked)` — `gas-accounting`.
- `352` `AddUint64(context.host.Runtime().GetPointsUsed(), gas)` —
  `gas-accounting`.
- `387` `AddUint64(context.host.Output().GetRefund(), gas)` —
  `gas-accounting`.
- `409` `AddUint64(context.initialCost, executionGasUsed)` —
  `gas-accounting`.
- `508` `MulUint64(codeSize, costPerByte)` — `gas-lock`.
- `511` `AddUint64(apiGasSchedule.AsyncCallStep, AsyncCallbackGasLock)` —
  `gas-lock`.
- `512` `AddUint64(compilationGasLock, executionGasLock)` — `gas-lock`.
- `565` `MulUint64(codeLength, costPerByte)` — `gas-charge`.
- `566` `AddUint64(baseCost, codeCost)` — `gas-charge`.

### `vmhost/contexts/output.go`

- `518` `AddUint64(account.BytesConsumedByTxAsNetworking, len(transfer.Data))`
  — `data-size-accounting`.

### `vmhost/contexts/storage.go`

- `143` `MulUint64(costPerByte, len(value))` — `gas-charge`.
- `158` `MulUint64(DataCopyPerByte, extraBytes)` — `gas-charge`.
- `412` `MulUint64(PersistPerByte, length)` — `gas-charge`.
- `413` `MulUint64(ReleasePerByte, newValueExtraLength)` — `gas-charge`.
- `419` `MulUint64(PersistPerByte, lengthOldValue)` — `gas-charge`.
- `420` `MulUint64(StorePerByte, newValueExtraLength)` — `gas-charge`.
- `421` `AddUint64(useGas, newValStoreUseGas)` — `gas-charge`.
- `427` `MulUint64(StorePerByte, length)` — `gas-charge`.
- `439` `MulUint64(ReleasePerByte, lengthOldValue)` — `gas-charge`.
- `461` `MulUint64(DataCopyPerByte, length)` — `gas-charge`.
- `494` `MulUint64(DataCopyPerByte, extraBytes)` — `gas-charge`.

### `vmhost/vmhooks/baseOps.go`

- `645`, `768`, `1030`, `1049`, `1112`, `1585`, `1604`, `1667`,
  `1795`, `1895`, `2728`, `2787`, `2825`, `3131`, `3561`, `3638`
  `MulUint64(...)` — `gas-charge`.
- `1729` / `1730` async minimum cost `AddUint64(MulUint64(...), ...)` —
  `gas-lock`.
- `1830` / `1831` async minimum cost `AddUint64(MulUint64(...), ...)` —
  `gas-lock`.
- `2732`, `2827`, `2828`, `3132` `AddUint64(...)` — `gas-charge`.

The `WriteLog` / `WriteEventLog` sites already have ISSUE-076 / ISSUE-081
bounds and pre-charge protections. This audit does not change their gas
semantics.

### `vmhost/vmhooks/bigFloatOps.go`

- `488` `vmMath.MulUint64(BigFloatPowPerIteration, exponent)` —
  `gas-charge`. This is currently added with `+=`; changing this expression is
  consensus-visible and should be handled under the same epoch-flag discipline
  if ever modified.

### `vmhost/vmhooks/bigIntOps.go`

- `102`, `379`, `416`, `446`, `477`, `1262`, `1294`, `1330`
  `MulUint64(...)` — `gas-charge`.

### `vmhost/vmhooks/cryptoei.go`

- `68`, `156`, `243`, `356`, `537`, `687` `MulUint64(...)` —
  `gas-charge`.
- `69`, `157`, `244` `AddUint64(baseCryptoCost, memLoadGas)` —
  `gas-charge`.

### `vmhost/vmhooks/manBufOps.go`

- `220`, `438`, `469`, `504`, `584`, `618`, `931`, `954`
  `MulUint64(...)` — `gas-charge`.
- `955` `AddUint64(baseGasToUse, lengthDependentGasToUse)` — `gas-charge`.

### `vmhost/vmhooks/managedConversions.go`

- `230` `MulUint64(DataCopyPerByte, actualLen)` — `gas-charge`.

### `vmhost/vmhooks/managedei.go`

- `232`, `674`, `1058` `MulUint64(...)` — `gas-charge`.
- `235` `AddUint64(gasToUse, gasForData)` — `gas-charge`.

## No Behavior Change

This audit intentionally does not convert any active VM call site to the
`WithErr` variants. The current saturating helper semantics are retained to
avoid replay-visible gas/accounting changes outside an explicit epoch-gated
feature decision.

