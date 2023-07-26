package migration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/networks"
)

func TestVersionUpgrade(t *testing.T) {
	t.Parallel()
	testEnvironment, testNetwork := setupUpgradeTest(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	upgradeImage, err := utils.GetEnv("UPGRADE_IMAGE")
	require.NoError(t, err, "Error getting upgrade image")
	upgradeVersion, err := utils.GetEnv("UPGRADE_VERSION")
	require.NoError(t, err, "Error getting upgrade version")

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")

	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	err = actions.UpgradeChainlinkNodeVersions(testEnvironment, upgradeImage, upgradeVersion, chainlinkNodes[0])
	require.NoError(t, err, "Upgrading chainlink nodes shouldn't fail")
}

func setupUpgradeTest(t *testing.T) (
	testEnvironment *environment.Environment,
	testNetwork blockchain.EVMNetwork,
) {
	testNetwork = networks.SelectedNetwork
	evmConfig := ethereum.New(nil)
	if !testNetwork.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}
	charts, err := chainlink.NewDeployment(1, map[string]any{
		"toml": client.AddNetworksConfig(config.BaseOCRP2PV1Config, testNetwork),
		"db": map[string]any{
			"stateful": true,
		},
	})
	require.NoError(t, err, "Error creating chainlink deployments")
	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("upgrade-version-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(evmConfig).
		AddHelmCharts(charts)
	err = testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")
	return testEnvironment, testNetwork
}
