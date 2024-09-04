package persistent_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	ccip_test_config "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent"
	chainlink_test_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// ok
func TestStartNewCCIPEnvironmentFromTestConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	testCfg := ccip_test_config.GlobalTestConfig()
	require.NoError(t, testCfg.Validate(), "Error validating test config")

	envConfig, err := persistent.EnvironmentConfigFromCCIPTestConfig(t, *testCfg, true)
	require.NoError(t, err, "Error creating chain config from test config")

	_, err = persistent.NewEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new persistent environment")

	//deploy something here to both chains?
}

// ok
func TestStartNewChainlinkEnvironmentFromTestConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	testCfg, err := chainlink_test_config.GetConfig([]string{"Smoke"}, chainlink_test_config.OCR)
	require.NoError(t, err, "Error getting test config")

	envConfig, err := persistent.EnvironmentConfigFromChainlinkTestConfig(t, testCfg, true, nil, nil)
	require.NoError(t, err, "Error creating chain config from test config")

	_, err = persistent.NewEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new persistent environment")

	//deploy something here to chain?
}
