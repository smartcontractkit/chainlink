package persistent_test

import (
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
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
	testCfg, err := chainlink_test_config.GetConfig([]string{"Smoke"}, chainlink_test_config.OCR2)
	require.NoError(t, err, "Error getting test config")

	envConfig, err := persistent.EnvironmentConfigFromChainlinkTestConfig(t, testCfg, true, nil, nil)
	require.NoError(t, err, "Error creating chain config from test config")

	_, err = persistent.NewEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new persistent environment")

	//deploy something here to chain?
}

// ok
func TestStartExistingChainlinkEnvironmentFromTestConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	testCfg, err := chainlink_test_config.GetConfig([]string{"Smoke"}, chainlink_test_config.OCR)
	require.NoError(t, err, "Error getting test config")

	// adjust this to match the existing cluster
	existingCluster := &ccip_test_config.CLCluster{
		Name:      ptr.Ptr("crib-bartek"),
		NoOfNodes: ptr.Ptr(5),
		NodeConfigs: []*client.ChainlinkConfig{
			{
				URL:        "https://crib-bartek-node1.main.stage.cldev.sh",
				Email:      "notreal@fakeemail.ch",
				Password:   "fj293fbBnlQ!f9vNs",
				InternalIP: "app-node-1",
			},
			{
				URL:        "https://crib-bartek-node2.main.stage.cldev.sh",
				Email:      "notreal@fakeemail.ch",
				Password:   "fj293fbBnlQ!f9vNs",
				InternalIP: "app-node-2",
			},
			{
				URL:        "https://crib-bartek-node3.main.stage.cldev.sh",
				Email:      "notreal@fakeemail.ch",
				Password:   "fj293fbBnlQ!f9vNs",
				InternalIP: "app-node-3",
			},
			{
				URL:        "https://crib-bartek-node4.main.stage.cldev.sh",
				Email:      "notreal@fakeemail.ch",
				Password:   "fj293fbBnlQ!f9vNs",
				InternalIP: "app-node-4",
			},
			{
				URL:        "https://crib-bartek-node5.main.stage.cldev.sh",
				Email:      "notreal@fakeemail.ch",
				Password:   "fj293fbBnlQ!f9vNs",
				InternalIP: "app-node-5",
			},
		},
	}

	envConfig, err := persistent.EnvironmentConfigFromChainlinkTestConfig(t, testCfg, false, existingCluster, ptr.Ptr("https://crib-bartek-mockserver.main.stage.cldev.sh"))
	require.NoError(t, err, "Error creating chain config from test config")

	_, err = persistent.NewEnvironment(lggr, envConfig)
	require.NoError(t, err, "Error creating new persistent environment")

	//deploy something here to chain?
}
