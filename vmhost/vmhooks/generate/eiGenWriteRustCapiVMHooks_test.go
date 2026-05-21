package vmhooksgenerate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteRustCapiVMHooksEmitsCheckedMemoryBridgeConversions(t *testing.T) {
	t.Parallel()

	outputDir := t.TempDir()
	outputFile := "capi_vm_hooks.rs"
	out := NewEIGenWriter(outputDir, outputFile)
	WriteRustCapiVMHooks(out, &EIMetadata{})
	out.Close()

	generated, err := os.ReadFile(filepath.Join(outputDir, outputFile))
	if err != nil {
		t.Fatal(err)
	}

	generatedCode := string(generated)
	requiredSnippets := []string{
		"use crate::capi_mem_conversion::{mem_length_to_i32, mem_ptr_to_i32};",
		"mem_ptr_to_i32(mem_ptr)",
		"mem_length_to_i32(mem_length)",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(generatedCode, snippet) {
			t.Fatalf("generated C-API VM hooks missing checked conversion snippet %q", snippet)
		}
	}

	forbiddenSnippets := []string{
		"mem_ptr as i32",
		"mem_length as i32",
	}
	for _, snippet := range forbiddenSnippets {
		if strings.Contains(generatedCode, snippet) {
			t.Fatalf("generated C-API VM hooks still contain truncating conversion %q", snippet)
		}
	}
}
