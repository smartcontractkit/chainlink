package migration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestVersionUpgrade(t *testing.T) {
	t.Parallel()

	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Migration", tc.Node)
	require.NoError(t, err, "Error getting config")

	err = config.ChainlinkUpgradeImage.Validate()
	require.NoError(t, err, "Error validating upgrade image")

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestConfig(&config).
		WithTestInstance(t).
		WithStandardCleanup().
		WithPrivateEthereumNetwork(privateNetwork).
		WithCLNodes(1).
		WithStandardCleanup().
		WithSeth().
		Build()
	require.NoError(t, err)

	// just restarting CL container with the same name, DB is still the same
	//
	// [Database]
	// MigrateOnStartup = true
	//
	// by default
	err = env.ClCluster.Nodes[0].UpgradeVersion(*config.ChainlinkUpgradeImage.Image, *config.ChainlinkUpgradeImage.Version)
	require.NoError(t, err)

}
