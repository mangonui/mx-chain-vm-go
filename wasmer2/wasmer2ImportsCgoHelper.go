package wasmer2

import (
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/executor"
)

// getVMHooksFromContextRawPtr resolves the wasmer-supplied context
// pointer to the VMHooks interface that should service this host import
// callback.
//
// ISSUE-011 (post-fix): the context slot now holds a uint64 registry
// HANDLE, not a Go heap ADDRESS. We read the handle and look it up in
// the global registry. A miss panics with a clear message — this is
// the security boundary; callers (50+ host import implementations in
// wasmer2ImportsCgo.go) all immediately invoke a method on the
// returned interface, so returning nil would just defer the panic to a
// cryptic nil-interface dispatch. Failing here with a clear message is
// strictly better debuggability for the same termination behaviour.
//
// The previous `*(*executor.VMHooks)(unsafe.Pointer(vmHooksPtr))`
// dereference (uintptr → *VMHooks → VMHooks value) is gone. `go vet`'s
// `unsafeptr` warning at this site is now silent.
func getVMHooksFromContextRawPtr(contextRawPtr unsafe.Pointer) executor.VMHooks {
	if contextRawPtr == nil {
		panic("wasmer2: nil context raw pointer in vmHooks lookup")
	}
	handle := *(*uint64)(contextRawPtr)
	return lookupVMHooksOrPanic(handle)
}

func funcPointer(cFuncPtr unsafe.Pointer) *[0]byte {
	return (*[0]byte)(cFuncPtr)
}
