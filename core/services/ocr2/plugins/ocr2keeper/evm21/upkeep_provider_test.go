package evm

import (
	"math/big"
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
					ID:      "edc5ec5f1d41b338a9ba6902caa8620a992fd086dcb978a6baa11a08e2e2795f",
					Trigger: ocr2keepers.Trigger{BlockNumber: 1, BlockHash: "0xc3bd2d00745c03048a5616146a96f5ff78e54efb9e5b04af208cdaff6f3830ee"},
				}, {
					ID:      "351362d44977dba636dd8b7429255db2e7b90a93ebddcb5f0f93cd81995ea887",
					Trigger: ocr2keepers.Trigger{BlockNumber: 1, BlockHash: "0xc3bd2d00745c03048a5616146a96f5ff78e54efb9e5b04af208cdaff6f3830ee"},
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

			for i, payload := range got {
				expected := tc.want[i]
				require.Equal(t, expected.ID, payload.ID)
				require.Equal(t, expected.Trigger.BlockNumber, payload.Trigger.BlockNumber)
				require.Equal(t, expected.Trigger.BlockHash, payload.Trigger.BlockHash)
			}
		})
	}
}
