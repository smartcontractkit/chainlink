package ring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUniqueSorted(t *testing.T) {
	t.Parallel()
	got := uniqueSorted([]string{"c", "a", "b", "a", "c"})
	require.Equal(t, []string{"a", "b", "c"}, got)
}
