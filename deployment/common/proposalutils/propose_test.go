package proposalutils_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	solanasdk "github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestBuildProposalFromBatchesV2(t *testing.T) {
	lggr := logger.TestLogger(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1,
		SolChains: 1,
		Chains:    2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	chainSelectorSolana := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
	config := proposalutils.SingleGroupMCMSV2(t)

	changeset.SetPreloadedSolanaAddresses(t, env, chainSelectorSolana)
	env, err := changeset.Apply(t, env, nil,
		changeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2),
			map[uint64]commontypes.MCMSWithTimelockConfigV2{
				chainSelector: {
					Canceller:        config,
					Bypasser:         config,
					Proposer:         config,
					TimelockMinDelay: big.NewInt(0),
				},
				chainSelectorSolana: {
					Canceller:        config,
					Bypasser:         config,
					Proposer:         config,
					TimelockMinDelay: big.NewInt(0),
				},
			},
		),
	)
	require.NoError(t, err)

	chain := env.BlockChains.EVMChains()[chainSelector]
	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)
	timelockAddress := mcmsState.Timelock.Address()
	require.NoError(t, err)

	solChain := env.BlockChains.SolanaChains()[chainSelectorSolana]
	addrs, err = env.ExistingAddresses.AddressesForChain(chainSelectorSolana)
	require.NoError(t, err)
	solState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addrs)
	require.NoError(t, err)

	solpk := solanasdk.NewWallet().PublicKey()

	timelockAddressPerChain := map[uint64]string{
		chainSelector:       timelockAddress.Hex(),
		chainSelectorSolana: solana.ContractAddress(solState.TimelockProgram, solana.PDASeed(solState.TimelockSeed)),
	}
	proposerAddressPerChain := map[uint64]string{
		chainSelector:       mcmsState.ProposerMcm.Address().Hex(),
		chainSelectorSolana: solana.ContractAddress(solState.McmProgram, solana.PDASeed(solState.ProposerMcmSeed)),
	}
	inspectorPerChain, err := proposalutils.McmsInspectors(env)
	require.NoError(t, err)

	description := "Test Proposal"
	minDelay := 24 * time.Hour

	solTx, err := solana.NewTransaction(solpk.String(), []byte("data1"), big.NewInt(0), []*solanasdk.AccountMeta{}, "", []string{})
	require.NoError(t, err)

	solMetadata, err := solana.NewChainMetadata(
		0,
		solState.McmProgram,
		solana.PDASeed(solState.ProposerMcmSeed),
		solState.ProposerAccessControllerAccount,
		solState.CancellerAccessControllerAccount,
		solState.BypasserAccessControllerAccount)
	require.NoError(t, err)

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
					Transactions: []types.Transaction{
						{
							To:               "0xRecipient1",
							Data:             []byte("data1"),
							AdditionalFields: json.RawMessage(`{"value": 0}`),
						},
					},
				},
				{
					ChainSelector: types.ChainSelector(chainSelectorSolana),
					Transactions:  []types.Transaction{solTx},
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
			proposal, err := proposalutils.BuildProposalFromBatchesV2(env, timelockAddressPerChain,
				proposerAddressPerChain, inspectorPerChain, tt.batches, description, proposalutils.TimelockConfig{MinDelay: minDelay})
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
				assert.Equal(t, map[types.ChainSelector]types.ChainMetadata{
					types.ChainSelector(chainSelector): {
						StartingOpCount: 0x0,
						MCMAddress:      proposerAddressPerChain[chainSelector],
					},
					types.ChainSelector(chainSelectorSolana): {
						StartingOpCount:  0x0,
						MCMAddress:       proposerAddressPerChain[chainSelectorSolana],
						AdditionalFields: solMetadata.AdditionalFields,
					},
				}, proposal.ChainMetadata)
				assert.Equal(t, map[types.ChainSelector]string{
					types.ChainSelector(chainSelector):       timelockAddressPerChain[chainSelector],
					types.ChainSelector(chainSelectorSolana): timelockAddressPerChain[chainSelectorSolana],
				}, proposal.TimelockAddresses)
				assert.Equal(t, tt.batches, proposal.Operations)
			}
		})
	}
}
