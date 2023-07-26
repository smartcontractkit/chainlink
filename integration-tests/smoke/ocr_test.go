package smoke

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker"
	"github.com/stretchr/testify/require"
)

func TestOCRBasic(t *testing.T) {
	t.Parallel()
	env, err := docker.NewChainlinkCluster(t, 6)
	require.NoError(t, err)
	clients, err := docker.ConnectClients(env)
	require.NoError(t, err)

	bootstrapNode, workerNodes := clients.Chainlink[0], clients.Chainlink[1:]
	clients.Networks[0].ParallelTransactions(true)

	linkTokenContract, err := clients.NetworkDeployers[0].DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	err = actions.FundChainlinkNodes(workerNodes, clients.Networks[0], big.NewFloat(.05))
	require.NoError(t, err, "Error funding Chainlink nodes")

	ocrInstances, err := actions.DeployOCRContracts(1, linkTokenContract, clients.NetworkDeployers[0], bootstrapNode, workerNodes, clients.Networks[0])
	require.NoError(t, err)
	err = clients.Networks[0].WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	err = actions.CreateOCRJobs(ocrInstances, bootstrapNode, workerNodes, 5, clients.Mockserver)
	require.NoError(t, err)
	err = actions.StartNewRound(1, ocrInstances, clients.Networks[0])
	require.NoError(t, err)

	answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

	err = actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstances, workerNodes, clients.Mockserver)
	require.NoError(t, err)
	err = actions.StartNewRound(2, ocrInstances, clients.Networks[0])
	require.NoError(t, err)

	answer, err = ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}
