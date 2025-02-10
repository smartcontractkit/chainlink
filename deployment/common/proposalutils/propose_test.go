package proposalutils_test

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	types2 "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestBuildProposalFromBatchesV2(t *testing.T) {
	lggr := logger.TestLogger(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.AllChainSelectors()[0]
	config := proposalutils.SingleGroupMCMS(t)

	env, err := changeset.ApplyChangesets(t, env, nil, []changeset.ChangesetApplication{
		{
			Changeset: changeset.WrapChangeSet(changeset.DeployMCMSWithTimelock),
			Config: map[uint64]types2.MCMSWithTimelockConfig{
				chainSelector: {
					Canceller:        config,
					Bypasser:         config,
					Proposer:         config,
					TimelockMinDelay: big.NewInt(0),
				},
			},
		},
	})
	require.NoError(t, err)

	chain := env.Chains[chainSelector]
	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)
	timelockAddress := mcmsState.Timelock.Address()
	require.NoError(t, err)

	timelockAddressPerChain := map[uint64]string{
		chainSelector: timelockAddress.Hex(),
	}
	proposerAddressPerChain := map[uint64]string{
		chainSelector: mcmsState.ProposerMcm.Address().Hex(),
	}
	inspectorPerChain := map[uint64]sdk.Inspector{
		chainSelector: evm.NewInspector(chain.Client),
	}

	description := "Test Proposal"
	minDelay := 24 * time.Hour

	tests := []struct {
		name    string
		batches []types.BatchOperation
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			batches: []types.BatchOperation{
				{
					ChainSelector: types.ChainSelector(chainSelector),
					Transactions:  []types.Transaction{{To: "0xRecipient1", Data: []byte("data1"), AdditionalFields: json.RawMessage(`{"value": 0}`)}},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid fields: missing required AdditionalFields",
			batches: []types.BatchOperation{
				{
					ChainSelector: types.ChainSelector(chainSelector),
					Transactions:  []types.Transaction{{To: "0xRecipient1", Data: []byte("data1")}},
				},
			},
			wantErr: true,
			errMsg:  "Key: 'TimelockProposal.Operations[0].Transactions[0].AdditionalFields' Error:Field validation for 'AdditionalFields' failed on the 'required' tag",
		},
		{
			name:    "empty batches",
			batches: []types.BatchOperation{},
			wantErr: true,
			errMsg:  "no operations in batch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal, err := proposalutils.BuildProposalFromBatchesV2(context.Background(), timelockAddressPerChain,
				proposerAddressPerChain, inspectorPerChain, tt.batches, description, minDelay)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, proposal)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, proposal)
				assert.Equal(t, "v1", proposal.Version)
				assert.Equal(t, string(types.TimelockActionSchedule), string(proposal.Action))
				//nolint:gosec // G115
				assert.InEpsilon(t, uint32(time.Now().Unix()+int64(proposalutils.DefaultValidUntil.Seconds())), proposal.ValidUntil, 1)
				assert.Equal(t, description, proposal.Description)
				assert.InEpsilon(t, minDelay.Seconds(), proposal.Delay.Seconds(), 0)
				assert.Equal(t, map[types.ChainSelector]types.ChainMetadata{0xc9f9284461c852b: {StartingOpCount: 0x0, MCMAddress: mcmsState.ProposerMcm.Address().String()}}, proposal.ChainMetadata)
				assert.Equal(t, timelockAddress.String(), proposal.TimelockAddresses[types.ChainSelector(chainSelector)])
				assert.Equal(t, tt.batches, proposal.Operations)
			}
		})
	}
}
