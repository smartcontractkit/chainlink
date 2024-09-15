package soak

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

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

	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, config, testsetups.WithForwarderFlow(true))
	require.NoError(t, err, "Error creating soak test")
	ocrSoakTest.DeployEnvironment(config)
	if ocrSoakTest.Environment().WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		} else {
			err := ocrSoakTest.Environment().Client.RemoveNamespace(ocrSoakTest.Environment().Cfg.Namespace)
			if err != nil {
				l.Error().Err(err).Msg("Error removing namespace")
			}
		}
	})
	ocrSoakTest.Setup(config)
	ocrSoakTest.Run()
}
