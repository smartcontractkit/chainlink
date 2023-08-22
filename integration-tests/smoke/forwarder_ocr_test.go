package smoke

import (
	"math/big"
	"testing"

	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/stretchr/testify/require"
)

func TestForwarderOCRBasic(t *testing.T) {
	t.Parallel()
	env, err := test_env.NewCLTestEnvBuilder().
		WithGeth().
		WithMockServer(1).
		WithForwarders().
		WithCLNodes(6).
		WithFunding(big.NewFloat(10)).
		Build()
	require.NoError(t, err)
	env.ParallelTransactions(true)

	nodeClients := env.GetAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	workerNodeAddresses, err := actions.ChainlinkNodeAddressesLocal(workerNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	linkTokenContract, err := env.ContractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	err = actions.FundChainlinkNodesLocal(workerNodes, env.EVMClient, big.NewFloat(.05))
	require.NoError(t, err, "Error funding Chainlink nodes")

	operators, authorizedForwarders, _ := actions.DeployForwarderContracts(
		t, env.ContractDeployer, linkTokenContract, env.EVMClient, len(workerNodes),
	)
	for i := range workerNodes {
		actions.AcceptAuthorizedReceiversOperator(
			t, operators[i], authorizedForwarders[i], []common.Address{workerNodeAddresses[i]}, env.EVMClient, env.ContractLoader,
		)
		require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
		err = actions.TrackForwarderLocal(env.EVMClient, authorizedForwarders[i], workerNodes[i])
		require.NoError(t, err)
		err = env.EVMClient.WaitForEvents()
	}
	ocrInstances, err := actions.DeployOCRContractsForwarderFlowLocal(
		1,
		linkTokenContract,
		env.ContractDeployer,
		workerNodes,
		authorizedForwarders,
		env.EVMClient,
	)

	err = actions.CreateOCRJobsWithForwarderLocal(ocrInstances, bootstrapNode, workerNodes, 5, env.MockServer.Client)
	require.NoError(t, err, "failed to setup forwarder jobs")
	err = actions.StartNewRound(1, ocrInstances, env.EVMClient)
	require.NoError(t, err)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")
	//time.Sleep(999 * time.Second)

	answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

	err = actions.SetAllAdapterResponsesToTheSameValueLocal(10, ocrInstances, workerNodes, env.MockServer.Client)
	require.NoError(t, err)
	err = actions.StartNewRound(2, ocrInstances, env.EVMClient)
	require.NoError(t, err)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	answer, err = ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}
