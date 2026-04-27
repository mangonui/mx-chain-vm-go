package wasmer2

import "unsafe"

func allocateVMHookPointers() *cWasmerVmHookPointers {
	vmHookPointers := (*cWasmerVmHookPointers)(cMalloc(unsafe.Sizeof(cWasmerVmHookPointers{})))
	*vmHookPointers = *populateCgoFunctionPointers()
	return vmHookPointers
}
