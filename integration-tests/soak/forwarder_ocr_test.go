package soak

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestForwarderOCRv1Soak(t *testing.T) {
	config, err := tc.GetConfig("Soak", tc.ForwarderOcr)
	require.NoError(t, err, "Error getting config")

	executeForwarderOCRSoakTest(t, &config)
}

func TestForwarderOCRv2Soak(t *testing.T) {
	config, err := tc.GetConfig("Soak", tc.ForwarderOcr2)
	require.NoError(t, err, "Error getting config")

	executeForwarderOCRSoakTest(t, &config)
}

func executeForwarderOCRSoakTest(t *testing.T, config *tc.TestConfig) {
	l := logging.GetTestLogger(t)

	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, config, true)
	require.NoError(t, err, "Error creating soak test")
	ocrSoakTest.DeployEnvironment(config)
	if ocrSoakTest.Environment().WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		if err := actions_seth.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	ocrSoakTest.Setup(config)
	ocrSoakTest.Run()
}
