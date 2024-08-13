package crib

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/havoc/k8schaos"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestCRIB(t *testing.T) {
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig([]string{"Crib"}, tc.OCR)
	require.NoError(t, err)

	sethClient, msClient, bootstrapNode, workerNodes, err := ConnectRemote()
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
		ch.AddListener(k8schaos.NewChaosLogger(l))
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
