package changeset_solana_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestAddRemoteChain(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	// Default env just has 2 chains with all contracts
	// deployed but no lanes.
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

	evmChain := tenv.Env.AllChainSelectors()[0]
	solChain := tenv.Env.AllChainSelectorsSolana()[0]

	state, err := changeset.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	_, err = commonchangeset.ApplyChangesets(t, tenv.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOnRampsDestsChangeset),
			Config: changeset.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]changeset.OnRampDestinationUpdate{
					evmChain: {
						solChain: {
							IsEnabled:        true,
							TestRouter:       false,
							AllowListEnabled: false,
						},
					},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.AddRemoteChainToSolana),
			Config: changeset_solana.AddRemoteChainToSolanaConfig{
				UpdatesByChain: map[uint64]map[uint64]changeset_solana.RemoteChainConfigSolana{
					solChain: {
						evmChain: {
							EnabledAsSource:          true,
							RemoteChainOnRampAddress: state.Chains[evmChain].OnRamp.Address().String(),
							DestinationConfig: solRouter.DestChainConfig{
								IsEnabled:                   true,
								DefaultTxGasLimit:           200000,
								MaxPerMsgGasLimit:           3000000,
								MaxDataBytes:                30000,
								MaxNumberOfTokensPerMsg:     5,
								DefaultTokenDestGasOverhead: 5000,
								// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
								// TODO: do a similar test for other chain families
								ChainFamilySelector: [4]uint8{40, 18, 213, 44},
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err = changeset.LoadOnchainStateSolana(tenv.Env)
	require.NoError(t, err)

	var sourceChainStateAccount solRouter.SourceChain
	evmSourceChainStatePDA, _ := solState.FindSourceChainStatePDA(evmChain, state.SolChains[solChain].Router)
	err = tenv.Env.SolChains[solChain].GetAccountDataBorshInto(ctx, evmSourceChainStatePDA, &sourceChainStateAccount)
	require.NoError(t, err)
	require.Equal(t, uint64(1), sourceChainStateAccount.State.MinSeqNr)
	require.True(t, sourceChainStateAccount.Config.IsEnabled)

	var destChainStateAccount solRouter.DestChain
	evmDestChainStatePDA, _ := solState.FindDestChainStatePDA(evmChain, state.SolChains[solChain].Router)
	err = tenv.Env.SolChains[solChain].GetAccountDataBorshInto(ctx, evmDestChainStatePDA, &destChainStateAccount)
	require.NoError(t, err)
}

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	testhelpers.DoDeployCCIPContracts(t, 1)
}
