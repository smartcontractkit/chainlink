package evm

import (
	"context"
	"math/big"
	"testing"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
)

func TestNewPayloadBuilder(t *testing.T) {
	for _, tc := range []struct {
		name         string
		activeList   ActiveUpkeepList
		recoverer    logprovider.LogRecoverer
		proposals    []types.CoordinatedBlockProposal
		wantPayloads []types.UpkeepPayload
	}{
		{
			name: "for log trigger upkeeps, new payloads are created",
			activeList: &mockActiveUpkeepList{
				IsActiveFn: func(id *big.Int) bool {
					return true
				},
			},
			proposals: []types.CoordinatedBlockProposal{
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "abc"),
					WorkID:   "workID1",
					Trigger: types.Trigger{
						BlockNumber: 1,
						BlockHash:   [32]byte{1},
					},
				},
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "def"),
					WorkID:   "workID2",
					Trigger: types.Trigger{
						BlockNumber: 2,
						BlockHash:   [32]byte{2},
					},
				},
			},
			recoverer: &mockLogRecoverer{
				BuildPayloadFn: func(ctx context.Context, proposal types.CoordinatedBlockProposal) (types.UpkeepPayload, error) {
					return types.UpkeepPayload{
						UpkeepID: proposal.UpkeepID,
						WorkID:   proposal.WorkID,
						Trigger:  proposal.Trigger,
					}, nil
				},
			},
			wantPayloads: []types.UpkeepPayload{
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "abc"),
					WorkID:   "workID1",
					Trigger: types.Trigger{
						BlockNumber: 1,
						BlockHash:   [32]byte{1},
					},
				},
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "def"),
					WorkID:   "workID2",
					Trigger: types.Trigger{
						BlockNumber: 2,
						BlockHash:   [32]byte{2},
					},
				},
			},
		},
		{
			name: "for an inactive log trigger upkeep, an empty payload is created",
			activeList: &mockActiveUpkeepList{
				IsActiveFn: func(id *big.Int) bool {
					if core.GenUpkeepID(types.LogTrigger, "ghi").BigInt().Cmp(id) == 0 {
						return false
					}
					return true
				},
			},
			proposals: []types.CoordinatedBlockProposal{
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "abc"),
					WorkID:   "workID1",
					Trigger: types.Trigger{
						BlockNumber: 1,
						BlockHash:   [32]byte{1},
					},
				},
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "def"),
					WorkID:   "workID2",
					Trigger: types.Trigger{
						BlockNumber: 2,
						BlockHash:   [32]byte{2},
					},
				},
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "ghi"),
					WorkID:   "workID3",
					Trigger: types.Trigger{
						BlockNumber: 3,
						BlockHash:   [32]byte{3},
					},
				},
			},
			recoverer: &mockLogRecoverer{
				BuildPayloadFn: func(ctx context.Context, proposal types.CoordinatedBlockProposal) (types.UpkeepPayload, error) {
					return types.UpkeepPayload{
						UpkeepID: proposal.UpkeepID,
						WorkID:   proposal.WorkID,
						Trigger:  proposal.Trigger,
					}, nil
				},
			},
			wantPayloads: []types.UpkeepPayload{
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "abc"),
					WorkID:   "workID1",
					Trigger: types.Trigger{
						BlockNumber: 1,
						BlockHash:   [32]byte{1},
					},
				},
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "def"),
					WorkID:   "workID2",
					Trigger: types.Trigger{
						BlockNumber: 2,
						BlockHash:   [32]byte{2},
					},
				},
				{},
			},
		},
		{
			name: "when the recoverer errors, an empty payload is returned",
			activeList: &mockActiveUpkeepList{
				IsActiveFn: func(id *big.Int) bool {
					return true
				},
			},
			proposals: []types.CoordinatedBlockProposal{
				{
					UpkeepID: core.GenUpkeepID(types.LogTrigger, "abc"),
					WorkID:   "workID1",
					Trigger: types.Trigger{
						BlockNumber: 1,
						BlockHash:   [32]byte{1},
					},
				},
			},
			recoverer: &mockLogRecoverer{
				BuildPayloadFn: func(ctx context.Context, proposal types.CoordinatedBlockProposal) (types.UpkeepPayload, error) {
					return types.UpkeepPayload{}, errors.New("recoverer boom")
				},
			},
			wantPayloads: []types.UpkeepPayload{
				{},
			},
		},
		{
			name: "currently a conditional upkeep does not have a new payload built, and an empty payload is added",
			activeList: &mockActiveUpkeepList{
				IsActiveFn: func(id *big.Int) bool {
					return true
				},
			},
			proposals: []types.CoordinatedBlockProposal{
				{
					UpkeepID: core.GenUpkeepID(types.ConditionTrigger, "def"),
					WorkID:   "workID1",
					Trigger: types.Trigger{
						BlockNumber: 1,
						BlockHash:   [32]byte{1},
					},
				},
			},
			wantPayloads: []types.UpkeepPayload{
				{},
			},
		},
		{
			name: "an unknown upkeep type does not have a new payload built, and an empty payload is added",
			activeList: &mockActiveUpkeepList{
				IsActiveFn: func(id *big.Int) bool {
					return true
				},
			},
			proposals: []types.CoordinatedBlockProposal{
				{
					UpkeepID: types.UpkeepIdentifier([32]byte{1}),
					WorkID:   "workID1",
					Trigger: types.Trigger{
						BlockNumber: 1,
						BlockHash:   [32]byte{1},
					},
				},
			},
			wantPayloads: []types.UpkeepPayload{
				{},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lggr, _ := logger.NewLogger()
			builder := NewPayloadBuilder(tc.activeList, tc.recoverer, lggr)
			payloads, err := builder.BuildPayloads(context.Background(), tc.proposals...)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantPayloads, payloads)
		})
	}
}

type mockLogRecoverer struct {
	logprovider.LogRecoverer
	BuildPayloadFn func(context.Context, types.CoordinatedBlockProposal) (types.UpkeepPayload, error)
}

func (r *mockLogRecoverer) BuildPayload(ctx context.Context, p types.CoordinatedBlockProposal) (types.UpkeepPayload, error) {
	return r.BuildPayloadFn(ctx, p)
}
