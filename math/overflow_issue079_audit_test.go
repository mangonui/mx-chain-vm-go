package math

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var issue079ProductionCallerPattern = regexp.MustCompile(`\b(math|vmMath)\.(AddUint64|MulUint64)\b`)

func TestIssue079AuditMentionsEveryProductionCallerFile(t *testing.T) {
	t.Parallel()

	auditBytes, err := os.ReadFile("ISSUE079_SATURATING_ARITHMETIC_AUDIT.md")
	require.NoError(t, err)

	auditText := string(auditBytes)
	productionCallerFiles := make(map[string]struct{})

	err = filepath.WalkDir("..", func(path string, d os.DirEntry, walkErr error) error {
		require.NoError(t, walkErr)

		if d.IsDir() {
			switch d.Name() {
			case ".git", "test", "testcommon", "mock":
				return filepath.SkipDir
			}

			return nil
		}

		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if filepath.Clean(path) == filepath.Clean("overflow.go") {
			return nil
		}

		contents, readErr := os.ReadFile(path)
		require.NoError(t, readErr)

		if issue079ProductionCallerPattern.Match(contents) {
			relativeToRepo := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(path)), "../")
			productionCallerFiles[relativeToRepo] = struct{}{}
		}

		return nil
	})
	require.NoError(t, err)

	require.NotEmpty(t, productionCallerFiles)
	for file := range productionCallerFiles {
		require.Contains(t, auditText, file)
	}
}
