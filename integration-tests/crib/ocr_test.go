package crib

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/havoc"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

func TestCRIB(t *testing.T) {
	l := logging.GetTestLogger(t)

	sethClient, msClient, bootstrapNode, workerNodes, err := ConnectRemote()
	require.NoError(t, err)

	lta, err := actions.SetupOCRv1Cluster(l, sethClient, workerNodes)
	require.NoError(t, err)
	ocrInstances, err := actions.SetupOCRv1Feed(l, sethClient, lta, msClient, bootstrapNode, workerNodes)
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
		ch.Create(context.Background())
		ch.AddListener(havoc.NewChaosLogger(l))
		t.Cleanup(func() {
			err := ch.Delete(context.Background())
			require.NoError(t, err)
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
