package migration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/test-go/testify/require"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

func TestVersionUpgrade(t *testing.T) {
	network := networks.SelectedNetwork
	evmConfig := ethereum.New(nil)
	if !network.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
	}
	charts, err := chainlink.NewDeployment(3, map[string]any{
		"toml": client.AddNetworksConfig("", network),
		"db": map[string]any{
			"stateful": true,
		},
	})
	require.NoError(t, err, "Error creating chainlink deployments")
	testEnvironment := environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("upgrade-version-%s", strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(evmConfig).
		AddHelmCharts(charts)
	err = testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	err = actions.UpgradeChainlinkNodeVersions(testEnvironment, "", "latest", chainlinkNodes[1])
	require.NoError(t, err, "Upgrading chainlink nodes shouldn't fail")
}
