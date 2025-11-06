package proposalutils_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	solanasdk "github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

func TestBuildProposalFromBatchesV2(t *testing.T) {
	quarantine.Flaky(t, "DX-1824")
	t.Parallel()

	evmSelector := chain_selectors.TEST_90000001.Selector
	solSelector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath, programIDs, ab := soltestutils.PreloadMCMS(t, solSelector)
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{evmSelector}),
		environment.WithSolanaContainer(t, []uint64{solSelector}, programsPath, programIDs),
		environment.WithAddressBook(ab),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	evmChain := rt.Environment().BlockChains.EVMChains()[evmSelector]
	solChain := rt.Environment().BlockChains.SolanaChains()[solSelector]

	config := proposalutils.SingleGroupMCMSV2(t)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2), map[uint64]commontypes.MCMSWithTimelockConfigV2{
			evmSelector: {
				Canceller:        config,
				Bypasser:         config,
				Proposer:         config,
				TimelockMinDelay: big.NewInt(0),
			},
			solSelector: {
				Canceller:        config,
				Bypasser:         config,
				Proposer:         config,
				TimelockMinDelay: big.NewInt(0),
			},
		}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().AddressBook.AddressesForChain(evmSelector)
	require.NoError(t, err)
	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockChainState(evmChain, addrs)
	require.NoError(t, err)

	addrs, err = rt.State().AddressBook.AddressesForChain(solSelector)
	require.NoError(t, err)
	solState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addrs)
	require.NoError(t, err)

	solpk := solanasdk.NewWallet().PublicKey()

	timelockAddressPerChain := map[uint64]string{
		evmSelector: mcmsState.Timelock.Address().Hex(),
		solSelector: solana.ContractAddress(solState.TimelockProgram, solana.PDASeed(solState.TimelockSeed)),
	}
	proposerAddressPerChain := map[uint64]string{
		evmSelector: mcmsState.ProposerMcm.Address().Hex(),
		solSelector: solana.ContractAddress(solState.McmProgram, solana.PDASeed(solState.ProposerMcmSeed)),
	}
	inspectorPerChain, err := proposalutils.McmsInspectors(rt.Environment())
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
					ChainSelector: types.ChainSelector(evmSelector),
					Transactions: []types.Transaction{
						{
							To:               "0xRecipient1",
							Data:             []byte("data1"),
							AdditionalFields: json.RawMessage(`{"value": 0}`),
						},
					},
				},
				{
					ChainSelector: types.ChainSelector(solSelector),
					Transactions:  []types.Transaction{solTx},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid fields: missing required AdditionalFields",
			batches: []types.BatchOperation{
				{
					ChainSelector: types.ChainSelector(evmSelector),
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
			proposal, err := proposalutils.BuildProposalFromBatchesV2(rt.Environment(), timelockAddressPerChain,
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
					types.ChainSelector(evmSelector): {
						StartingOpCount: 0x0,
						MCMAddress:      proposerAddressPerChain[evmSelector],
					},
					types.ChainSelector(solSelector): {
						StartingOpCount:  0x0,
						MCMAddress:       proposerAddressPerChain[solSelector],
						AdditionalFields: solMetadata.AdditionalFields,
					},
				}, proposal.ChainMetadata)
				assert.Equal(t, map[types.ChainSelector]string{
					types.ChainSelector(evmSelector): timelockAddressPerChain[evmSelector],
					types.ChainSelector(solSelector): timelockAddressPerChain[solSelector],
				}, proposal.TimelockAddresses)
				assert.Equal(t, tt.batches, proposal.Operations)
			}
		})
	}
}
