package v1_5_1_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/stretchr/testify/require"
)

func TestDeployTokenPoolFactoryChangeset(t *testing.T) {
	t.Parallel()

	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
		testCfg.Chains = 2
		testCfg.PrerequisiteDeploymentOnly = true
	})
	e := deployedEnvironment.Env
	selectors := e.AllChainSelectors()

	e, err := commonchangeset.Apply(t, e, nil, commonchangeset.Configure(
		v1_5_1.DeployTokenPoolFactoryChangeset,
		v1_5_1.DeployTokenPoolFactoryConfig{
			Chains: selectors,
		},
	))
	require.NoError(t, err, "failed to apply DeployTokenPoolFactoryChangeset")

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err, "failed to load onchain state")

	for _, chainSel := range selectors {
		tpf := state.Chains[chainSel].TokenPoolFactory
		require.NotNil(t, tpf, "token pool factory should be deployed on chain %d", chainSel)
		typeAndVersion, err := tpf.TypeAndVersion(nil)
		require.NoError(t, err, "failed to get type and version of token pool factory on chain %d", chainSel)
		require.Equal(t, "TokenPoolFactory 1.5.1", typeAndVersion, "unexpected type and version of token pool factory on chain %d", chainSel)
	}
}
