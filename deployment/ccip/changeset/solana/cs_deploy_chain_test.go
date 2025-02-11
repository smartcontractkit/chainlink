package solana_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
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
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	contractParams := make(map[uint64]changeset.ChainContractParams)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = proposalutils.SingleGroupTimelockConfig(t)
		contractParams[chain] = changeset.ChainContractParams{
			FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
			OffRampParams:   changeset.DefaultOffRampParams(),
		}
	}
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	for _, chain := range e.AllChainSelectors() {
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
	}

	testhelpers.SavePreloadedSolAddresses(t, e, solChainSelectors[0])
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChainChangeset),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    e.AllChainSelectors(),
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    e.AllChainSelectorsSolana(),
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    cfg,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployPrerequisitesChangeset),
			Config: changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContractsChangeset),
			Config: changeset.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(solana.DeployChainContractsChangesetSolana),
			Config: changeset.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]changeset.ChainContractParams{
					solChainSelectors[0]: {
						FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
						OffRampParams:   changeset.DefaultOffRampParams(),
					},
				},
			},
		},
	})
	require.NoError(t, err)
	// solana verification
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
}
