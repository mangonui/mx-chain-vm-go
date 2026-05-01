package wasmer2

import "testing"

func TestManagedDRWASyncMirror_IsRegisteredInWasmerImports(t *testing.T) {
	t.Parallel()

	if _, ok := functionNames["managedDRWASyncMirror"]; !ok {
		t.Fatalf("managedDRWASyncMirror missing from exported Wasmer import names")
	}
	if _, ok := functionNames["managedDRWANativeGovernanceQuery"]; !ok {
		t.Fatalf("managedDRWANativeGovernanceQuery missing from exported Wasmer import names")
	}

	pointers := populateCgoFunctionPointers()
	if pointers.managed_drwa_sync_mirror_func_ptr == nil {
		t.Fatalf("managed_drwa_sync_mirror_func_ptr was not populated")
	}
	if pointers.managed_drwa_native_governance_query_func_ptr == nil {
		t.Fatalf("managed_drwa_native_governance_query_func_ptr was not populated")
	}
}
