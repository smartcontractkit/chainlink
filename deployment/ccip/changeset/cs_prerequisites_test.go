package changeset_test

import (
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

func TestDeployPrerequisites(t *testing.T) {
	t.Parallel()

	e, err := environment.New(t.Context(),
		environment.WithEVMSimulatedN(t, 2),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	testDeployPrerequisitesWithEnv(t, *e)
}

func TestDeployPrerequisitesZk(t *testing.T) {
	// Timeouts in CI
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-6427")

	t.Parallel()

	e, err := environment.New(t.Context(),
		environment.WithZKSyncContainerN(t, 2),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	testDeployPrerequisitesWithEnv(t, *e)
}

func testDeployPrerequisitesWithEnv(t *testing.T, e cldf.Environment) {
	newChain := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))[0]
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
	require.NotNil(t, chainState.BurnMintTokenPools[shared.FactoryBurnMintERC20Symbol][shared.CurrentTokenPoolVersion])
	require.NotNil(t, chainState.BurnFromMintTokenPools[shared.FactoryBurnMintERC20Symbol][shared.CurrentTokenPoolVersion])
	require.NotNil(t, chainState.BurnWithFromMintTokenPools[shared.FactoryBurnMintERC20Symbol][shared.CurrentTokenPoolVersion])
	require.NotNil(t, chainState.LockReleaseTokenPools[shared.FactoryBurnMintERC20Symbol][shared.CurrentTokenPoolVersion])
	require.NotNil(t, chainState.RegistryModules1_6)
	require.NotNil(t, chainState.Router)
}
