package ring

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUniqueSorted(t *testing.T) {
	t.Parallel()
	got := uniqueSorted([]string{"c", "a", "b", "a", "c"})
	require.Equal(t, []string{"a", "b", "c"}, got)
}

func TestLocateShard(t *testing.T) {
	t.Parallel()

	_, err := LocateShard(nil, "wf")
	require.ErrorIs(t, err, ErrNoHealthyShards)

	counts := map[uint32]int{}
	for i := range 1000 {
		wfID := "workflow-" + strconv.Itoa(i)
		shard, err := LocateShard([]uint32{1, 2}, wfID)
		require.NoError(t, err)
		reversed, err := LocateShard([]uint32{2, 1}, wfID)
		require.NoError(t, err)
		require.Equal(t, shard, reversed, "Shard order should not change the assignment")
		counts[shard]++
	}
	require.Greater(t, counts[1], 350, "Shard 1 should get a fair share of workflows")
	require.Greater(t, counts[2], 350, "Shard 2 should get a fair share of workflows")
}
