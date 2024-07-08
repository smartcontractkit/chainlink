package crib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func TestCRIB(t *testing.T) {
	l := logging.GetTestLogger(t)

	sethClient, msClient, bootstrapNode, workerNodes, err := ConnectRemote(true)
	require.NoError(t, err)

	lta, err := actions.SetupOCRv1Cluster(l, sethClient, workerNodes)
	require.NoError(t, err)
	ocrInstances, err := actions.SetupOCRv1Feed(l, sethClient, lta, msClient, bootstrapNode, workerNodes)
	require.NoError(t, err)

	err = actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstances, workerNodes, msClient)
	require.NoError(t, err)

	err = actions.WatchNewOCRRound(l, sethClient, 2, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(3*time.Minute))
	require.NoError(t, err, "Error watching for new OCR round")

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}
