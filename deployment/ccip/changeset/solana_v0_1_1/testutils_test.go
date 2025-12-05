package solana_test

import (
	"os"
	"testing"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/stretchr/testify/require"
)

func skipInCI(t *testing.T) {
	t.Helper()

	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}

func initialDeployCS(t *testing.T, e cldf.Environment, buildConfig *ccipChangesetSolana.BuildSolanaConfig) []commonchangeset.ConfiguredChangeSet {
	evmSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	homeChainSel := evmSelectors[0]
	solChainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	feeAggregatorPrivKey, _ := solana.NewRandomPrivateKey()
	feeAggregatorPubKey := feeAggregatorPrivKey.PublicKey()
	mcmsConfig := proposalutils.SingleGroupTimelockConfigV2(t)
	solLinkTokenPrivKey, _ := solana.NewRandomPrivateKey()
	return []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.BlockChains.EVMChains()[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeploySolanaLinkToken),
			commonchangeset.DeploySolanaLinkTokenConfig{
				ChainSelector: solChainSelectors[0],
				TokenPrivKey:  solLinkTokenPrivKey,
				TokenDecimals: 9,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ChainSelector:     solChainSelectors[0],
				ContractParamsPerChain: ccipChangesetSolana.ChainContractParams{
					FeeQuoterParams: ccipChangesetSolana.FeeQuoterParams{
						DefaultMaxFeeJuelsPerMsg: solBinary.Uint128{Lo: 300000000, Hi: 0, Endianness: nil},
					},
					OffRampParams: ccipChangesetSolana.OffRampParams{
						EnableExecutionAfter: int64(globals.PermissionLessExecutionThreshold.Seconds()),
					},
				},
				MCMSWithTimelockConfig: &mcmsConfig,
				BuildConfig:            buildConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployReceiverForTest),
			ccipChangesetSolana.DeployForTestConfig{
				ChainSelector: solChainSelectors[0],
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
			ccipChangesetSolana.SetFeeAggregatorConfig{
				ChainSelector: solChainSelectors[0],
				FeeAggregator: feeAggregatorPubKey.String(),
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.ExtendGlobalLookupTableChangeset),
			ccipChangesetSolana.ExtendGlobalLookupTableConfig{
				ChainSelector:   solChainSelectors[0],
				LookupTableKeys: []solana.PublicKey{feeAggregatorPubKey, solLinkTokenPrivKey.PublicKey()}, // just add some random keys
			},
		),
	}
}
