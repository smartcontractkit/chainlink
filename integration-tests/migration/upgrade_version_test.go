package migration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"

	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

func TestVersionUpgrade(t *testing.T) {
	t.Parallel()
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithStandardCleanup().
		WithGeth().
		WithCLNodes(1).
		Build()
	require.NoError(t, err)

	upgradeImage, err := osutil.GetEnv("UPGRADE_IMAGE")
	require.NoError(t, err, "Error getting upgrade image")
	upgradeVersion, err := osutil.GetEnv("UPGRADE_VERSION")
	require.NoError(t, err, "Error getting upgrade version")
	// [Database]
	// MigrateOnStartup = true
	//
	// by default
	err = env.ClCluster.Nodes[0].UpgradeVersion(upgradeImage, upgradeVersion)
	require.NoError(t, err)

}
