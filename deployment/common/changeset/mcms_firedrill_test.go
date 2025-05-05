package changeset_test

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// setupFiredrillTestEnv deploys all required contracts for the firedrill proposal execution
func setupFiredrillTestEnv(t *testing.T) deployment.Environment {
	lggr := logger.TestLogger(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains:    2,
		SolChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.AllChainSelectors()[0]
	chainSelector2 := env.AllChainSelectors()[1]
	chainSelectorSolana := env.AllChainSelectorsSolana()[0]

	commonchangeset.SetPreloadedSolanaAddresses(t, env, chainSelectorSolana)
	config := proposalutils.SingleGroupTimelockConfigV2(t)
	// Deploy MCMS and Timelock
	env, err := commonchangeset.Apply(t, env, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			map[uint64]commontypes.MCMSWithTimelockConfigV2{
				chainSelector:       config,
				chainSelector2:      config,
				chainSelectorSolana: config,
			},
		),
	)
	require.NoError(t, err)
	//nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
	addresses, err := env.ExistingAddresses.AddressesForChain(chainSelectorSolana)
	require.NoError(t, err)
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(env.SolChains[env.AllChainSelectorsSolana()[0]], addresses)
	require.NoError(t, err)
	timelockSigner := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	mcmSigner := state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed)
	mcmSignerBypasser := state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.BypasserMcmSeed)
	solChain := env.SolChains[chainSelectorSolana]
	memory.FundSolanaAccounts(env.GetContext(), t, []solana.PublicKey{timelockSigner, mcmSigner, mcmSignerBypasser, solChain.DeployerKey.PublicKey()}, 150, solChain.Client)
	require.NoError(t, err)
	return env
}

func TestMCMSSignFireDrillChangeset(t *testing.T) {
	t.Parallel()
	env := setupFiredrillTestEnv(t)
	chainSelector := env.AllChainSelectors()[0]
	chainSelector2 := env.AllChainSelectors()[1]
	chainSelectorSolana := env.AllChainSelectorsSolana()[0]
	// Add the timelock as a signer to check state changes
	for _, tc := range []struct {
		name       string
		changeSets func() []commonchangeset.ConfiguredChangeSet
	}{
		{
			name: "MCMS Firedrill execution",
			changeSets: func() []commonchangeset.ConfiguredChangeSet {
				return []commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.MCMSSignFireDrillChangeset),
						commonchangeset.FireDrillConfig{
							Selectors: []uint64{chainSelector, chainSelector2, chainSelectorSolana},
							TimelockCfg: proposalutils.TimelockConfig{
								MCMSAction: mcmsTypes.TimelockActionBypass,
							},
						},
					),
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			changesetsToApply := tc.changeSets()
			_, _, err := commonchangeset.ApplyChangesetsV2(t, env, changesetsToApply)
			require.NoError(t, err)
		})
	}
}
