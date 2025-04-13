package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployPrerequisites(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	newChain := e.AllChainSelectors()[0]
	cfg := changeset.DeployPrerequisiteConfig{
		Configs: []changeset.DeployPrerequisiteConfigPerChain{
			{
				ChainSelector: newChain,
				Opts: []changeset.PrerequisiteOpt{
					changeset.WithTokenPoolFactoryEnabled(),
				},
			},
		},
	}
	output, err := changeset.DeployPrerequisitesChangeset(e, cfg)
	require.NoError(t, err)
	err = e.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err)
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[newChain].Weth9)
	require.NotNil(t, state.Chains[newChain].TokenAdminRegistry)
	require.NotNil(t, state.Chains[newChain].TokenPoolFactory)
	require.NotNil(t, state.Chains[newChain].FactoryBurnMintERC20Token)
	require.NotNil(t, state.Chains[newChain].RegistryModules1_6)
	require.NotNil(t, state.Chains[newChain].Router)
}
