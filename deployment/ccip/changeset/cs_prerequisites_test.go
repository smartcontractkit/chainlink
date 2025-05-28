package changeset_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
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

	testDeployPrerequisitesWithEnv(t, e)
}

func TestDeployPrerequisitesZk(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		ZkChains:   2,
		Nodes:      4,
	})

	testDeployPrerequisitesWithEnv(t, e)
}

func testDeployPrerequisitesWithEnv(t *testing.T, e cldf.Environment) {
	newChain := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
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
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	chainState, _ := state.EVMChainState(newChain)
	require.NotNil(t, chainState.Weth9)
	require.NotNil(t, chainState.TokenAdminRegistry)
	require.NotNil(t, chainState.TokenPoolFactory)
	require.NotNil(t, chainState.FactoryBurnMintERC20Token)
	require.NotNil(t, chainState.RegistryModules1_6)
	require.NotNil(t, chainState.Router)
}
