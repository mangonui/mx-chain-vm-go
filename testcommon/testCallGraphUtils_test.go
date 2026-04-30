package testcommon

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetGraphsFolder_UsesTempDirFallback(t *testing.T) {
	t.Setenv(graphsFolderEnvKey, "")

	require.Equal(t, filepath.Join(os.TempDir(), "mx-chain-vm-go-graphs"), getGraphsFolder())
}

func TestGetGraphsFolder_UsesEnvOverride(t *testing.T) {
	t.Setenv(graphsFolderEnvKey, "/tmp/custom-graphs")

	require.Equal(t, "/tmp/custom-graphs", getGraphsFolder())
}
