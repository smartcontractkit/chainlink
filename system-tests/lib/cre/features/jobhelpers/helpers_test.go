package jobhelpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeSpecsByIndexPreservesInputOrder(t *testing.T) {
	t.Parallel()

	merged, err := MergeSpecsByIndex([]map[string][]string{
		{"node-a": {"first-a"}, "node-b": {"first-b"}},
		{"node-a": {"second-a"}},
		{"node-b": {"second-b"}},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"first-a", "second-a"}, merged["node-a"])
	require.Equal(t, []string{"first-b", "second-b"}, merged["node-b"])
}
