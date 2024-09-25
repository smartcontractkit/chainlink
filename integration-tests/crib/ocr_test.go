package crib

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	ocr_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/ocr"
)

// TestCRIBChaos an example of how we can run chaos tests with havoc and core CRIB
func TestCRIBChaos(t *testing.T) {
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig([]string{"Crib"}, tc.OCR)
	require.NoError(t, err)

	sethClient, msClient, bootstrapNode, workerNodes, _, err := ConnectRemote()
	require.NoError(t, err)

	lta, err := actions.SetupOCRv1Cluster(l, sethClient, config.OCR, workerNodes)
	require.NoError(t, err)
	ocrInstances, err := actions.SetupOCRv1Feed(l, sethClient, lta, config.OCR, msClient, bootstrapNode, workerNodes)
	require.NoError(t, err)

	err = actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstances, workerNodes, msClient)
	require.NoError(t, err)
	actions.SimulateOCRv1EAActivity(l, 3*time.Second, ocrInstances, workerNodes, msClient)

	err = actions.WatchNewOCRRound(l, sethClient, 1, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), 5*time.Minute)
	require.NoError(t, err, "Error watching for new OCR round")

	if os.Getenv("TEST_PERSISTENCE") != "" {
		ch, err := rebootCLNamespace(
			1*time.Second,
			os.Getenv("CRIB_NAMESPACE"),
		)
		require.NoError(t, err, "Error rebooting CL namespace")
		ch.Create(context.Background())
		ch.AddListener(havoc.NewChaosLogger(l))
		t.Cleanup(func() {
			err := ch.Delete(context.Background())
			require.NoError(t, err, "Error deleting chaos")
		})
		require.Eventually(t, func() bool {
			err = actions.WatchNewOCRRound(l, sethClient, 3, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), 5*time.Minute)
			if err != nil {
				l.Info().Err(err).Msg("OCR round is not there yet")
				return false
			}
			return true
		}, 20*time.Minute, 5*time.Second)
	}
}

// TestCRIBRPCChaos and example of how we can run RPC chaos with Geth or Anvil
func TestCRIBRPCChaos(t *testing.T) {
	l := logging.GetTestLogger(t)

	sethClient, msClient, bootstrapNode, workerNodes, vars, err := ConnectRemote()
	require.NoError(t, err)

	ocrConfig := &ocr_config.Config{
		Contracts: &ocr_config.Contracts{
			ShouldBeUsed: ptr.Ptr(false),
		},
	}

	lta, err := actions.SetupOCRv1Cluster(l, sethClient, ocrConfig, workerNodes)
	require.NoError(t, err)
	ocrInstances, err := actions.SetupOCRv1Feed(l, sethClient, lta, ocrConfig, msClient, bootstrapNode, workerNodes)
	require.NoError(t, err)

	err = actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstances, workerNodes, msClient)
	require.NoError(t, err)
	actions.SimulateOCRv1EAActivity(l, 3*time.Second, ocrInstances, workerNodes, msClient)

	err = actions.WatchNewOCRRound(l, sethClient, 1, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), 5*time.Minute)
	require.NoError(t, err, "Error watching for new OCR round")

	ac := client.NewRPCClient(sethClient.URL, vars.BlockchainNodeHeaders)
	err = ac.GethSetHead(10)
	require.NoError(t, err)

	err = actions.WatchNewOCRRound(l, sethClient, 3, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), 5*time.Minute)
	require.NoError(t, err, "Error watching for new OCR round")
}
