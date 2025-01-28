package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHead_LatestFinalizedHead(t *testing.T) {
	t.Parallel()
	newFinalizedHead := func(num int64) *Head {
		result := &Head{Number: num}
		result.IsFinalized.Store(true)
		return result
	}
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
			Head:      sliceToChain(&Head{}, &Head{}, &Head{}),
			Finalized: nil,
		},
		{
			Name:      "Returns head if it's finalized",
			Head:      sliceToChain(newFinalizedHead(2), newFinalizedHead(1)),
			Finalized: &Head{Number: 2},
		},
		{
			Name:      "Returns first block in chain if it's finalized",
			Head:      sliceToChain(&Head{Number: 3}, newFinalizedHead(2), newFinalizedHead(1)),
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

func TestHead_ChainString(t *testing.T) {
	cases := []struct {
		Name           string
		Chain          *Head
		ExpectedResult string
	}{
		{
			Name:           "Empty chain",
			ExpectedResult: "->nil",
		},
		{
			Name:           "Single head",
			Chain:          &Head{Number: 1},
			ExpectedResult: "Head{Number: 1, Hash: 0x0000000000000000000000000000000000000000000000000000000000000000, ParentHash: 0x0000000000000000000000000000000000000000000000000000000000000000}->nil",
		},
		{
			Name:           "Multiple heads",
			Chain:          sliceToChain(&Head{Number: 1}, &Head{Number: 2}, &Head{Number: 3}),
			ExpectedResult: "Head{Number: 1, Hash: 0x0000000000000000000000000000000000000000000000000000000000000000, ParentHash: 0x0000000000000000000000000000000000000000000000000000000000000000}->Head{Number: 2, Hash: 0x0000000000000000000000000000000000000000000000000000000000000000, ParentHash: 0x0000000000000000000000000000000000000000000000000000000000000000}->Head{Number: 3, Hash: 0x0000000000000000000000000000000000000000000000000000000000000000, ParentHash: 0x0000000000000000000000000000000000000000000000000000000000000000}->nil",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.ExpectedResult, testCase.Chain.ChainString())
		})
	}
}

func sliceToChain(heads ...*Head) *Head {
	if len(heads) == 0 {
		return nil
	}

	for i := 1; i < len(heads); i++ {
		heads[i-1].Parent.Store(heads[i])
	}

	return heads[0]
}
