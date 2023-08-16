package evm

import (
	"bytes"
	"math/big"
	"sort"
	"testing"

	coreTypes "github.com/ethereum/go-ethereum/core/types"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	clientmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

func TestUpkeepProvider_GetActiveUpkeeps(t *testing.T) {
	t.Skip()
	ctx := testutils.Context(t)
	c := new(clientmocks.Client)

	bs := &BlockSubscriber{}
	var lp logpoller.LogPoller

	tests := []struct {
		name        string
		active      ActiveUpkeepList
		blockHeader coreTypes.Header
		want        []ocr2keepers.UpkeepPayload
		wantErr     bool
	}{
		{
			"empty",
			&mockActiveUpkeepList{
				ViewFn: func(upkeepType ...ocr2keepers.UpkeepType) []*big.Int {
					return []*big.Int{}
				},
			},
			coreTypes.Header{Number: big.NewInt(0)},
			nil,
			false,
		},
		{
			"happy flow",
			&mockActiveUpkeepList{
				ViewFn: func(upkeepType ...ocr2keepers.UpkeepType) []*big.Int {
					return []*big.Int{
						big.NewInt(1),
						big.NewInt(2),
					}
				},
			},
			coreTypes.Header{Number: big.NewInt(1)},
			[]ocr2keepers.UpkeepPayload{
				{
					UpkeepID: core.UpkeepIDFromInt("10"),
				},
				{
					UpkeepID: core.UpkeepIDFromInt("15"),
				},
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := coreTypes.NewBlockWithHeader(&tc.blockHeader)
			c.On("BlockByNumber", mock.Anything, mock.Anything).Return(b, nil)

			p := NewUpkeepProvider(tc.active, bs, lp)

			got, err := p.GetActiveUpkeeps(ctx)
			require.NoError(t, err)
			require.Len(t, got, len(tc.want))
			sort.Slice(got, func(i, j int) bool {
				return bytes.Compare(got[i].UpkeepID[:], got[j].UpkeepID[:]) < 0
			})
			for i, payload := range got {
				expected := tc.want[i]
				// require.Equal(t, expected.ID, payload.ID) // TODO: uncomment once we change to workID
				require.Equal(t, expected.UpkeepID, payload.UpkeepID)
				require.Equal(t, b.Number().Int64(), payload.Trigger.BlockNumber)
			}
		})
	}
}

type mockActiveUpkeepList struct {
	ActiveUpkeepList
	ViewFn func(...ocr2keepers.UpkeepType) []*big.Int
}

func (l *mockActiveUpkeepList) View(u ...ocr2keepers.UpkeepType) []*big.Int {
	return l.ViewFn(u...)
}
