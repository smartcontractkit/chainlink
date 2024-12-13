package v1_5

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type TestConfig struct {
	Chains             int
	NumOfUsersPerChain int
	Nodes              int
	Bootstraps         int
}

func NewEnvironment(t *testing.T, tc *changeset.TestConfigs, tEnv changeset.TestEnvironment) changeset.DeployedEnv {
	var err error
	lggr := logger.Test(t)
	tEnv.StartChains(t, tc)
	e := tEnv.DeployedEnvironment()
	require.NotEmpty(t, e.Env.Chains)
	ab := deployment.NewMemoryAddressBook()
	tEnv.StartNodes(t, tc, deployment.CapabilityRegistryConfig{})
	e = tEnv.DeployedEnvironment()
	allChains := e.Env.AllChainSelectors()

	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, c := range e.Env.AllChainSelectors() {
		mcmsCfg[c] = commontypes.MCMSWithTimelockConfig{
			Canceller:        commonchangeset.SingleGroupMCMS(t),
			Bypasser:         commonchangeset.SingleGroupMCMS(t),
			Proposer:         commonchangeset.SingleGroupMCMS(t),
			TimelockMinDelay: big.NewInt(0),
		}
	}
	var prereqCfg []changeset.DeployPrerequisiteConfigPerChain
	for _, chain := range allChains {
		var opts []changeset.PrerequisiteOpt
		if tc != nil {
			if tc.IsUSDC {
				opts = append(opts, changeset.WithUSDCEnabled())
			}
			if tc.IsMultiCall3 {
				opts = append(opts, changeset.WithMulticall3())
			}
		}
		opts = append(opts, changeset.WithLegacyDeployment(changeset.LegacyDeploymentConfig{
			PriceRegStalenessThreshold: 60 * 60 * 24 * 14, // two weeks
		}))
		prereqCfg = append(prereqCfg, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
			Opts:          opts,
		})
	}

	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    allChains,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployPrerequisites),
			Config: changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
	})
	require.NoError(t, err)
}
