package migration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestVersionUpgrade(t *testing.T) {
	t.Parallel()

	config, err := tc.GetConfig(t.Name(), tc.Migration, tc.Node)
	require.NoError(t, err, "Error getting config")

	err = config.ChainlinkUpgradeImage.Validate()
	require.NoError(t, err, "Error validating upgrade image")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestConfig(&config).
		WithTestInstance(t).
		WithGeth().
		WithCLNodes(1).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	// just restarting CL container with the same name, DB is still the same
	//
	// [Database]
	// MigrateOnStartup = true
	//
	// by default
	env.ClCluster.Nodes[0].ContainerImage = *config.ChainlinkUpgradeImage.Image
	env.ClCluster.Nodes[0].ContainerVersion = *config.ChainlinkUpgradeImage.Version
	err = env.ClCluster.Nodes[0].Restart(env.ClCluster.Nodes[0].NodeConfig)
	require.NoError(t, err)
}
