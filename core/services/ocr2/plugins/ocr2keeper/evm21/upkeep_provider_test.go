package evm

import (
	"bytes"
	"math/big"
	"sort"
	"sync"
	"testing"

	coreTypes "github.com/ethereum/go-ethereum/core/types"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	clientmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestUpkeepProvider_GetActiveUpkeeps(t *testing.T) {
	t.Skip()
	ctx := testutils.Context(t)
	c := new(clientmocks.Client)

	r := &EvmRegistry{
		mu:     sync.RWMutex{},
		active: map[string]activeUpkeep{},
		client: c,
	}

	p := NewUpkeepProvider(r)

	tests := []struct {
		name        string
		active      map[string]activeUpkeep
		blockHeader coreTypes.Header
		want        []ocr2keepers.UpkeepPayload
		wantErr     bool
	}{
		{
			"empty",
			map[string]activeUpkeep{},
			coreTypes.Header{Number: big.NewInt(0)},
			nil,
			false,
		},
		{
			"happy flow",
			map[string]activeUpkeep{
				"1": {
					ID: big.NewInt(1),
				},
				"2": {
					ID: big.NewInt(2),
				},
			},
			coreTypes.Header{Number: big.NewInt(1)},
			[]ocr2keepers.UpkeepPayload{
				{
					Upkeep: ocr2keepers.ConfiguredUpkeep{
						ID: ocr2keepers.UpkeepIdentifier(big.NewInt(1).String()),
					},
				}, {
					Upkeep: ocr2keepers.ConfiguredUpkeep{
						ID: ocr2keepers.UpkeepIdentifier(big.NewInt(2).String()),
					},
				},
			},
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := coreTypes.NewBlockWithHeader(&tc.blockHeader)
			c.On("BlockByNumber", mock.Anything, mock.Anything).Return(b, nil)

			r.mu.Lock()
			r.active = tc.active
			r.mu.Unlock()

			got, err := p.GetActiveUpkeeps(ctx, BlockKeyHelper[int64]{}.MakeBlockKey(b.Number().Int64()))
			require.NoError(t, err)
			require.Len(t, got, len(tc.want))
			sort.Slice(got, func(i, j int) bool {
				return bytes.Compare(got[i].Upkeep.ID, got[j].Upkeep.ID) < 0
			})
			for i, payload := range got {
				expected := tc.want[i]
				// require.Equal(t, expected.ID, payload.ID) // TODO: uncomment once we change to workID
				require.Equal(t, expected.Upkeep.ID, payload.Upkeep.ID)
				require.Equal(t, b.Number().Int64(), payload.Trigger.BlockNumber)
			}
		})
	}
}
