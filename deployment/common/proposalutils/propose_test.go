package proposalutils_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	solanasdk "github.com/gagliardetto/solana-go"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/quarantine"
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
	solpk := solanasdk.NewWallet().PublicKey()

	evmTx := types.Transaction{To: "0xRecipient1", Data: []byte("data1"), AdditionalFields: json.RawMessage(`{"value": 0}`)}
	solTx, err := solana.NewTransaction(solpk.String(), []byte("data1"), big.NewInt(0), []*solanasdk.AccountMeta{}, "", []string{})
	require.NoError(t, err)
	batches := []types.BatchOperation{
		{
			ChainSelector: types.ChainSelector(evmSelector),
			Transactions:  []types.Transaction{evmTx},
		}, {
			ChainSelector: types.ChainSelector(solSelector),
			Transactions:  []types.Transaction{solTx},
		},
	}

	evmMetadata := types.ChainMetadata{
		StartingOpCount: 0,
		MCMAddress:      mcmsState.ProposerMcm.Address().Hex(),
	}
	solMetadata, err := solana.NewChainMetadata(
		0,
		solState.McmProgram,
		solana.PDASeed(solState.ProposerMcmSeed),
		solState.ProposerAccessControllerAccount,
		solState.CancellerAccessControllerAccount,
		solState.BypasserAccessControllerAccount)
	require.NoError(t, err)

	wantProposal := &mcms.TimelockProposal{
		BaseProposal: mcms.BaseProposal{
			Version:    "v1",
			Kind:       "TimelockProposal",
			ValidUntil: 1234,
			ChainMetadata: map[types.ChainSelector]types.ChainMetadata{
				types.ChainSelector(evmSelector): evmMetadata,
				types.ChainSelector(solSelector): solMetadata,
			},
			Description: description,
		},
		Action: "schedule",
		Delay:  types.NewDuration(minDelay),
		TimelockAddresses: func(addrs map[uint64]string) map[types.ChainSelector]string {
			copiedAddrs := make(map[types.ChainSelector]string, len(addrs))
			for k, v := range addrs {
				copiedAddrs[types.ChainSelector(k)] = v
			}
			return copiedAddrs
		}(timelockAddressPerChain),
		Operations: []types.BatchOperation{
			{ChainSelector: types.ChainSelector(evmSelector), Transactions: []types.Transaction{evmTx}},
			{ChainSelector: types.ChainSelector(solSelector), Transactions: []types.Transaction{solTx}},
		},
	}

	tests := []struct {
		name       string
		batches    []types.BatchOperation
		inspectors map[uint64]mcmssdk.Inspector
		options    []proposalutils.BuildProposalOption
		want       *mcms.TimelockProposal
		wantErr    string
	}{
		{
			name:       "success: explicit inspectors",
			batches:    batches,
			inspectors: inspectorPerChain,
			options:    []proposalutils.BuildProposalOption{},
			want:       wantProposal,
		},
		{
			name:       "success: implicit inspectors",
			batches:    batches,
			inspectors: nil,
			options:    []proposalutils.BuildProposalOption{},
			want:       wantProposal,
		},
		{
			name:       "success: extra chain metadata",
			batches:    batches,
			inspectors: nil,
			options: []proposalutils.BuildProposalOption{
				proposalutils.WithChainMetadata(proposalutils.ChainMetadata{
					evmSelector: map[string]any{
						"gasLimit": 100,
						"gasPrice": 50,
					},
				}),
			},
			want: func() *mcms.TimelockProposal {
				proposal := *wantProposal
				proposal.ChainMetadata = map[types.ChainSelector]types.ChainMetadata{
					types.ChainSelector(evmSelector): {
						StartingOpCount:  0,
						MCMAddress:       mcmsState.ProposerMcm.Address().Hex(),
						AdditionalFields: json.RawMessage(`{"gasLimit":100,"gasPrice":50}`),
					},
					types.ChainSelector(solSelector): solMetadata,
				}
				t.Logf("PROPOSAL1.ChainMetadata:     %#v", proposal.ChainMetadata)
				t.Logf("WANTPROPOSAL1.ChainMetadata: %#v", wantProposal.ChainMetadata)
				return &proposal
			}(),
		},
		{
			name: "invalid fields: missing required AdditionalFields",
			batches: []types.BatchOperation{
				{
					ChainSelector: types.ChainSelector(evmSelector),
					Transactions:  []types.Transaction{{To: "0xRecipient1", Data: []byte("data1")}},
				},
			},
			options: []proposalutils.BuildProposalOption{},
			wantErr: "Key: 'TimelockProposal.Operations[0].Transactions[0].AdditionalFields' Error:Field validation for 'AdditionalFields' failed on the 'required' tag",
		},
		{
			name:    "empty batches",
			batches: []types.BatchOperation{},
			options: []proposalutils.BuildProposalOption{},
			wantErr: "no operations in batch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal, err := proposalutils.BuildProposalFromBatchesV2(rt.Environment(),
				timelockAddressPerChain, proposerAddressPerChain, tt.inspectors, tt.batches, description,
				proposalutils.TimelockConfig{MinDelay: minDelay}, tt.options...)

			if tt.wantErr == "" {
				require.NoError(t, err)
				require.Empty(t, cmp.Diff(tt.want, proposal,
					cmpopts.IgnoreFields(mcms.BaseProposal{}, "useSimulatedBackend", "ValidUntil")))
			} else {
				require.Nil(t, proposal)
				require.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}
