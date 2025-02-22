package solana_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	cs_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestDeployChainContractsChangesetSolana(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.AllChainSelectorsSolana()
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	contractParams := make(map[uint64]v1_6.ChainContractParams)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfigV2(t)
		contractParams[chain] = v1_6.ChainContractParams{
			FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
			OffRampParams:   v1_6.DefaultOffRampParams(),
		}
	}
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	feeAggregatorPrivKey, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey := feeAggregatorPrivKey.PublicKey()
	ci := os.Getenv("CI") == "true"
	// we can't upgrade in place locally so we have to change where we build
	buildCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(cs_solana.BuildSolanaChangeset),
		cs_solana.BuildSolanaConfig{
			ChainSelector:       solChainSelectors[0],
			GitCommitSha:        "0863d8fed5fbada9f352f33c405e1753cbb7d72c",
			DestinationDir:      e.SolChains[solChainSelectors[0]].ProgramsPath,
			CleanDestinationDir: true,
		},
	)
	if ci {
		testhelpers.SavePreloadedSolAddresses(t, e, solChainSelectors[0])
	} else {
		e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			buildCs,
		})
		require.NoError(t, err)
	}

	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			e.AllChainSelectors(),
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			e.AllChainSelectorsSolana(),
		),

		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			cfg,
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
			changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			v1_6.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(cs_solana.DeployChainContractsChangesetSolana),
			cs_solana.DeployChainContractsConfigSolana{
				DeployChainContractsConfig: v1_6.DeployChainContractsConfig{
					HomeChainSelector: homeChainSel,
					ContractParamsPerChain: map[uint64]v1_6.ChainContractParams{
						solChainSelectors[0]: {
							FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
							OffRampParams:   v1_6.DefaultOffRampParams(),
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(cs_solana.SetFeeAggregator),
			cs_solana.SetFeeAggregatorConfig{
				ChainSelector: solChainSelectors[0],
				FeeAggregator: feeAggregatorPubKey.String(),
			},
		),
	})
	require.NoError(t, err)
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
	timelockSignerPDA, _ := testhelpers.TransferOwnershipSolana(t, &e, solChainSelectors[0], true, true, true, true)
	upgradeAuthority := timelockSignerPDA

	// make sure idempotency works and setting the upgrade authority
	upgradeAuthorityCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(cs_solana.DeployChainContractsChangesetSolana),
		cs_solana.DeployChainContractsConfigSolana{
			DeployChainContractsConfig: v1_6.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]v1_6.ChainContractParams{
					solChainSelectors[0]: {
						FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
						OffRampParams:   v1_6.DefaultOffRampParams(),
					},
				},
			},
			NewUpgradeAuthority: &upgradeAuthority,
		},
	)
	upgradeCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(cs_solana.DeployChainContractsChangesetSolana),
		cs_solana.DeployChainContractsConfigSolana{
			DeployChainContractsConfig: v1_6.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]v1_6.ChainContractParams{
					solChainSelectors[0]: {
						FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
						OffRampParams:   v1_6.DefaultOffRampParams(),
					},
				},
			},
			UpgradeConfig: cs_solana.UpgradeConfigSolana{
				NewFeeQuoterVersion: &deployment.Version1_1_0,
				NewRouterVersion:    &deployment.Version1_1_0,
				UpgradeAuthority:    upgradeAuthority,
				SpillAddress:        upgradeAuthority,
				MCMS: &cs.MCMSConfig{
					MinDelay: 1 * time.Second,
				},
			},
		},
	)
	// because we cannot upgrade in place locally, we can't redeploy offramp
	offRampCs := commonchangeset.Configure(
		deployment.CreateLegacyChangeSet(cs_solana.DeployChainContractsChangesetSolana),
		cs_solana.DeployChainContractsConfigSolana{
			DeployChainContractsConfig: v1_6.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]v1_6.ChainContractParams{
					solChainSelectors[0]: {
						FeeQuoterParams: v1_6.DefaultFeeQuoterParams(),
						OffRampParams:   v1_6.DefaultOffRampParams(),
					},
				},
			},
			UpgradeConfig: cs_solana.UpgradeConfigSolana{
				NewOffRampVersion: &deployment.Version1_1_0,
				UpgradeAuthority:  upgradeAuthority,
				SpillAddress:      upgradeAuthority,
				MCMS: &cs.MCMSConfig{
					MinDelay: 1 * time.Second,
				},
			},
		},
	)
	if ci {
		e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			upgradeAuthorityCs,
			buildCs,
			upgradeCs,
		})
		require.NoError(t, err)
		state, err := changeset.LoadOnchainStateSolana(e)
		require.NoError(t, err)
		testhelpers.ValidateSolanaState(t, e, solChainSelectors)
		oldOffRampAddress := state.SolChains[solChainSelectors[0]].OffRamp
		// add a second offramp address
		e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			offRampCs,
		})
		require.NoError(t, err)
		// verify the offramp address is different
		state, err = changeset.LoadOnchainStateSolana(e)
		require.NoError(t, err)
		newOffRampAddress := state.SolChains[solChainSelectors[0]].OffRamp
		require.NotEqual(t, oldOffRampAddress, newOffRampAddress)
	} else {
		e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			upgradeAuthorityCs,
			upgradeCs,
		})
		require.NoError(t, err)
	}
	// Verify router and fee quoter upgraded in place
	// and offramp had 2nd address added
	addresses, err := e.ExistingAddresses.AddressesForChain(solChainSelectors[0])
	require.NoError(t, err)
	numRouters := 0
	numFeeQuoters := 0
	numOffRamps := 0
	for _, address := range addresses {
		if address.Type == changeset.Router {
			numRouters++
		}
		if address.Type == changeset.FeeQuoter {
			numFeeQuoters++
		}
		if address.Type == changeset.OffRamp {
			numOffRamps++
		}
	}
	require.Equal(t, 1, numRouters)
	require.Equal(t, 1, numFeeQuoters)
	if ci {
		require.Equal(t, 2, numOffRamps)
	} else {
		require.Equal(t, 1, numOffRamps)
	}
	require.NoError(t, err)
	// solana verification
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
}
