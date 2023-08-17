package migration

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"os"
)

func TestVersionUpgrade(t *testing.T) {
	t.Parallel()
	env, err := test_env.NewCLTestEnvBuilder().
		WithGeth().
		WithCLNodes(1).
		Build()
	require.NoError(t, err)

	upgradeImage, err := utils.GetEnv("UPGRADE_IMAGE")
	require.NoError(t, err, "Error getting upgrade image")
	upgradeVersion, err := utils.GetEnv("UPGRADE_VERSION")
	require.NoError(t, err, "Error getting upgrade version")

	_ = os.Setenv("CHAINLINK_IMAGE", upgradeImage)
	_ = os.Setenv("CHAINLINK_VERSION", upgradeVersion)

	// just restarting CL container with the same name, DB is still the same
	//
	// [Database]
	// MigrateOnStartup = true
	//
	// by default
	err = env.CLNodes[0].Restart(env.CLNodes[0].NodeConfig)
	require.NoError(t, err)
}
