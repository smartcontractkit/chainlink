package soak

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestForwarderOCRSoak(t *testing.T) {
	l := logging.GetTestLogger(t)
	// Use this variable to pass in any custom EVM specific TOML values to your Chainlink nodes
	customNetworkTOML := `[EVM.Transactions]
ForwardersEnabled = true`
	// Uncomment below for debugging TOML issues on the node
	// fmt.Println("Using Chainlink TOML\n---------------------")
	// fmt.Println(client.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, customNetworkTOML, network))
	// fmt.Println("---------------------")

	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, true)
	require.NoError(t, err, "Error creating soak test")
	ocrSoakTest.DeployEnvironment(customNetworkTOML)
	if ocrSoakTest.Environment().WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	ocrSoakTest.Setup()
	ocrSoakTest.Run()
}
