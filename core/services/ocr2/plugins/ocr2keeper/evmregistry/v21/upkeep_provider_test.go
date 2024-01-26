package evm

import (
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

func TestUpkeepProvider_GetActiveUpkeeps(t *testing.T) {
	ctx := testutils.Context(t)

	var lp logpoller.LogPoller

	tests := []struct {
		name        string
		active      ActiveUpkeepList
		latestBlock *ocr2keepers.BlockKey
		want        []ocr2keepers.UpkeepPayload
		wantErr     bool
	}{
		{
			"empty",
			&mockActiveUpkeepList{
				ViewFn: func(upkeepType ...types.UpkeepType) []*big.Int {
					return []*big.Int{}
				},
			},
			&ocr2keepers.BlockKey{Number: 1},
			nil,
			false,
		},
		{
			"happy flow",
			&mockActiveUpkeepList{
				ViewFn: func(upkeepType ...types.UpkeepType) []*big.Int {
					return []*big.Int{
						big.NewInt(1),
						big.NewInt(2),
					}
				},
			},
			&ocr2keepers.BlockKey{Number: 1},
			[]ocr2keepers.UpkeepPayload{
				{
					UpkeepID: core.UpkeepIDFromInt("1"),
					Trigger:  ocr2keepers.NewTrigger(ocr2keepers.BlockNumber(1), [32]byte{}),
					WorkID:   "b10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6",
				},
				{
					UpkeepID: core.UpkeepIDFromInt("2"),
					Trigger:  ocr2keepers.NewTrigger(ocr2keepers.BlockNumber(1), [32]byte{}),
					WorkID:   "405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace",
				},
			},
			false,
		},
		{
			"latest block not found",
			&mockActiveUpkeepList{
				ViewFn: func(upkeepType ...types.UpkeepType) []*big.Int {
					return []*big.Int{
						big.NewInt(1),
						big.NewInt(2),
					}
				},
			},
			nil,
			[]ocr2keepers.UpkeepPayload{},
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := &BlockSubscriber{
				latestBlock: atomic.Pointer[ocr2keepers.BlockKey]{},
			}
			bs.latestBlock.Store(tc.latestBlock)
			p := NewUpkeepProvider(tc.active, bs, lp)

			got, err := p.GetActiveUpkeeps(ctx)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, got, len(tc.want))
			require.Equal(t, tc.want, got)
		})
	}
}

type mockActiveUpkeepList struct {
	ActiveUpkeepList
	ViewFn     func(...types.UpkeepType) []*big.Int
	IsActiveFn func(id *big.Int) bool
}

func (l *mockActiveUpkeepList) View(u ...types.UpkeepType) []*big.Int {
	return l.ViewFn(u...)
}

func (l *mockActiveUpkeepList) IsActive(id *big.Int) bool {
	return l.IsActiveFn(id)
}
