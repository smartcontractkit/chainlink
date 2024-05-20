package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHead_LatestFinalizedHead(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Name      string
		Head      *Head
		Finalized *Head
	}{
		{
			Name:      "Empty chain returns nil on finalized",
			Head:      nil,
			Finalized: nil,
		},
		{
			Name:      "Chain without finalized returns nil",
			Head:      &Head{Parent: &Head{Parent: &Head{}}},
			Finalized: nil,
		},
		{
			Name:      "Returns head if it's finalized",
			Head:      &Head{Number: 2, IsFinalized: true, Parent: &Head{Number: 1, IsFinalized: true}},
			Finalized: &Head{Number: 2},
		},
		{
			Name:      "Returns first block in chain if it's finalized",
			Head:      &Head{Number: 3, IsFinalized: false, Parent: &Head{Number: 2, IsFinalized: true, Parent: &Head{Number: 1, IsFinalized: true}}},
			Finalized: &Head{Number: 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := tc.Head.LatestFinalizedHead()
			if tc.Finalized == nil {
				assert.Nil(t, actual)
			} else {
				require.NotNil(t, actual)
				assert.Equal(t, tc.Finalized.Number, actual.BlockNumber())
			}
		})
	}
}
