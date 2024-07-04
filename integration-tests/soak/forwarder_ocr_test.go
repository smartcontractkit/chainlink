package soak

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestForwarderOCRv1Soak(t *testing.T) {
	//nolint:revive
	t.Fatalf("This test is disabled because the implementation is broken")
	config, err := tc.GetConfig([]string{"Soak"}, tc.ForwarderOcr)
	require.NoError(t, err, "Error getting config")

	executeForwarderOCRSoakTest(t, &config)
}

func TestForwarderOCRv2Soak(t *testing.T) {
	//nolint:revive
	t.Fatalf("This test is disabled because the implementation is broken")
	config, err := tc.GetConfig([]string{"Soak"}, tc.ForwarderOcr2)
	require.NoError(t, err, "Error getting config")

	executeForwarderOCRSoakTest(t, &config)
}

func executeForwarderOCRSoakTest(t *testing.T, config *tc.TestConfig) {
	l := logging.GetTestLogger(t)
	// Use this variable to pass in any custom EVM specific TOML values to your Chainlink nodes
	customNetworkTOML := `[EVM.Transactions]
ForwardersEnabled = true`
	// Uncomment below for debugging TOML issues on the node
	// fmt.Println("Using Chainlink TOML\n---------------------")
	// fmt.Println(networks.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, customNetworkTOML, network))
	// fmt.Println("---------------------")

	config, err := tc.GetConfig("Soak", tc.OCR)
	require.NoError(t, err, "Error getting config")

	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config, testsetups.WithForwarderFlow(true))
	require.NoError(t, err, "Error creating soak test")
	ocrSoakTest.DeployEnvironment(customNetworkTOML, &config)
	if ocrSoakTest.Environment().WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	ocrSoakTest.Setup(&config)
	ocrSoakTest.Run()
}
