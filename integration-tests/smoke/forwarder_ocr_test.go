package smoke

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestForwarderOCRBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.ForwarderOcr)
	if err != nil {
		t.Fatal(err)
	}

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithGeth().
		WithMockAdapter().
		WithForwarders().
		WithCLNodes(6).
		WithFunding(big.NewFloat(.1)).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	env.ParallelTransactions(true)

	nodeClients := env.ClCluster.NodeAPIs()
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
		err = actions.TrackForwarderLocal(env.EVMClient, authorizedForwarders[i], workerNodes[i], l)
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
	require.NoError(t, err, "Error deploying OCR contracts")

	err = actions.CreateOCRJobsWithForwarderLocal(ocrInstances, bootstrapNode, workerNodes, 5, env.MockAdapter, env.EVMClient.GetChainID().String())
	require.NoError(t, err, "failed to setup forwarder jobs")
	err = actions.WatchNewRound(1, ocrInstances, env.EVMClient, l)
	require.NoError(t, err)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

	err = actions.SetAllAdapterResponsesToTheSameValueLocal(10, ocrInstances, workerNodes, env.MockAdapter)
	require.NoError(t, err)
	err = actions.WatchNewRound(2, ocrInstances, env.EVMClient, l)
	require.NoError(t, err)
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	answer, err = ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}
