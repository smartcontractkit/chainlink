package chainlink

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/assert"
)

// TestSortNodeStatuses is a regression test for #18263. The /nodes dashboard
// page was reordering every render because per-relayer node iteration is not
// guaranteed to be stable. The aggregator now sorts by (ChainID, Name) so
// repeated calls with the same inputs return the same page.
func TestSortNodeStatuses(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   []types.NodeStatus
		want []types.NodeStatus
	}{
		{
			name: "shuffled within a chain sorts by name",
			in: []types.NodeStatus{
				{ChainID: "1", Name: "zebra"},
				{ChainID: "1", Name: "alpha"},
				{ChainID: "1", Name: "mango"},
			},
			want: []types.NodeStatus{
				{ChainID: "1", Name: "alpha"},
				{ChainID: "1", Name: "mango"},
				{ChainID: "1", Name: "zebra"},
			},
		},
		{
			name: "groups by chain id then sorts by name",
			in: []types.NodeStatus{
				{ChainID: "2", Name: "bear"},
				{ChainID: "1", Name: "zebra"},
				{ChainID: "2", Name: "ant"},
				{ChainID: "1", Name: "alpha"},
			},
			want: []types.NodeStatus{
				{ChainID: "1", Name: "alpha"},
				{ChainID: "1", Name: "zebra"},
				{ChainID: "2", Name: "ant"},
				{ChainID: "2", Name: "bear"},
			},
		},
		{
			name: "lexicographic chain ids (string compare, not numeric)",
			in: []types.NodeStatus{
				{ChainID: "10", Name: "n1"},
				{ChainID: "2", Name: "n2"},
				{ChainID: "1", Name: "n3"},
			},
			// ChainID is a string in NodeStatus, so "10" sorts before "2"
			// lexicographically. The sort intentionally mirrors that — chain
			// id space is opaque (EVM, Solana, dummy, ...) so a numeric sort
			// would be wrong for non-EVM chains.
			want: []types.NodeStatus{
				{ChainID: "1", Name: "n3"},
				{ChainID: "10", Name: "n1"},
				{ChainID: "2", Name: "n2"},
			},
		},
		{
			name: "already sorted is a no-op",
			in: []types.NodeStatus{
				{ChainID: "1", Name: "a"},
				{ChainID: "1", Name: "b"},
				{ChainID: "2", Name: "c"},
			},
			want: []types.NodeStatus{
				{ChainID: "1", Name: "a"},
				{ChainID: "1", Name: "b"},
				{ChainID: "2", Name: "c"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]types.NodeStatus(nil), tc.in...)
			sortNodeStatuses(got)
			assert.Equal(t, tc.want, got)
		})
	}

	t.Run("empty and nil slices are no-ops", func(t *testing.T) {
		empty := []types.NodeStatus{}
		sortNodeStatuses(empty)
		assert.Empty(t, empty)
		sortNodeStatuses(nil)
	})
}

// TestSortNodeStatusesStable asserts equal keys preserve input order so the
// aggregator output stays predictable even when two relayers happen to use
// the same (ChainID, Name) pair (e.g. dev fixtures, multi-relayer dummy
// chains).
func TestSortNodeStatusesStable(t *testing.T) {
	t.Parallel()

	in := []types.NodeStatus{
		{ChainID: "1", Name: "same", Config: "first"},
		{ChainID: "1", Name: "same", Config: "second"},
		{ChainID: "1", Name: "same", Config: "third"},
	}
	got := append([]types.NodeStatus(nil), in...)
	sortNodeStatuses(got)
	assert.Equal(t, in, got, "equal-key entries must keep their original order")
}
