package migration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/test-go/testify/require"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
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
	testEnvironment := environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("upgrade-version-%s", strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(evmConfig).
		AddHelmCharts(chainlink.NewDeployment(3, map[string]any{
			"toml": client.AddNetworksConfig("", network),
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")

	_, err = client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	// testEnvironment.Client.ListPods(testEnvironment.Cfg.Namespace)
	for _, chart := range testEnvironment.Charts {
		log.Warn().Str("Path", chart.GetPath()).Str("Name", chart.GetName()).Msg("Chart")
	}
}
