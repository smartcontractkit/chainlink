package solana_test

import (
	"os"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	solBinary "github.com/gagliardetto/binary"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestDeployChainContractsChangesetSolana(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
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
	prereqCfg := make([]ccipChangeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, ccipChangeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	feeAggregatorPrivKey, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey := feeAggregatorPrivKey.PublicKey()
	ci := os.Getenv("CI") == "true"
	// we can't upgrade in place locally if we preload addresses so we have to change where we build
	// we also don't want to incur two builds in CI, so only do it locally
	if ci {
		testhelpers.SavePreloadedSolAddresses(t, e, solChainSelectors[0])
	} else {
		e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				deployment.CreateLegacyChangeSet(ccipChangesetSolana.BuildSolanaChangeset),
				ccipChangesetSolana.BuildSolanaConfig{
					ChainSelector:       solChainSelectors[0],
					GitCommitSha:        "3da552ac9d30b821310718b8b67e6a298335a485",
					DestinationDir:      e.SolChains[solChainSelectors[0]].ProgramsPath,
					CleanDestinationDir: true,
				},
			),
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
			deployment.CreateLegacyChangeSet(ccipChangeset.DeployPrerequisitesChangeset),
			ccipChangeset.DeployPrerequisiteConfig{
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
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]ccipChangesetSolana.ChainContractParams{
					solChainSelectors[0]: {
						FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
							DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
						},
						OffRampParams: ccipChangesetSolana.OffRampParams{
							EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
			ccipChangesetSolana.SetFeeAggregatorConfig{
				ChainSelector: solChainSelectors[0],
				FeeAggregator: feeAggregatorPubKey.String(),
			},
		),
	})
	require.NoError(t, err)
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
}
