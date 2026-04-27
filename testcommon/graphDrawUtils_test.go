package testcommon

import (
	"testing"

	"github.com/awalterschulze/gographviz"
	"github.com/stretchr/testify/require"
)

func TestCreateSvgWithLocation_PathTraversalPanics(t *testing.T) {
	t.Parallel()

	graphviz := gographviz.NewGraph()

	require.PanicsWithValue(t,
		"file name must not contain path traversal sequences",
		func() {
			CreateSvgWithLocation(t.TempDir(), "..", graphviz)
		},
	)
}
