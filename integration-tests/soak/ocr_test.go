package soak

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestOCRv1Soak(t *testing.T) {
	config, err := tc.GetConfig("Soak", tc.OCR)
	require.NoError(t, err, "Error getting config")
	executeOCRSoakTest(t, &config)
}

func TestOCRv2Soak(t *testing.T) {
	config, err := tc.GetConfig("Soak", tc.OCR2)
	require.NoError(t, err, "Error getting config")

	executeOCRSoakTest(t, &config)
}

func TestOCRSoak_GethReorgBelowFinality_FinalityTagDisabled(t *testing.T) {
	config, err := tc.GetConfig(t.Name(), tc.OCR)
	require.NoError(t, err, "Error getting config")
	executeOCRSoakTest(t, &config)
}

func TestOCRSoak_GethReorgBelowFinality_FinalityTagEnabled(t *testing.T) {
	config, err := tc.GetConfig(t.Name(), tc.OCR)
	require.NoError(t, err, "Error getting config")
	executeOCRSoakTest(t, &config)
}

func TestOCRSoak_GasSpike(t *testing.T) {
	config, err := tc.GetConfig(t.Name(), tc.OCR)
	require.NoError(t, err, "Error getting config")
	executeOCRSoakTest(t, &config)
}

// TestOCRSoak_ChangeBlockGasLimit changes next block gas limit and sets it to percentage of last gasUsed in previous block creating congestion
func TestOCRSoak_ChangeBlockGasLimit(t *testing.T) {
	config, err := tc.GetConfig(t.Name(), tc.OCR)
	require.NoError(t, err, "Error getting config")
	executeOCRSoakTest(t, &config)
}

func executeOCRSoakTest(t *testing.T, config *tc.TestConfig) {
	l := logging.GetTestLogger(t)

	// validate Seth config before anything else, but only for live networks (simulated will fail, since there's no chain started yet)
	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	if !network.Simulated {
		_, err := actions_seth.GetChainClient(config, network)
		require.NoError(t, err, "Error creating seth client")
	}

	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, config, false)
	require.NoError(t, err, "Error creating soak test")
	if !ocrSoakTest.Interrupted() {
		ocrSoakTest.DeployEnvironment(config)
	}
	if ocrSoakTest.Environment().WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		if err := actions_seth.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	if ocrSoakTest.Interrupted() {
		err = ocrSoakTest.LoadState()
		require.NoError(t, err, "Error loading state")
		ocrSoakTest.Resume()
	} else {
		ocrSoakTest.Setup(config)
		ocrSoakTest.Run()
	}
}
