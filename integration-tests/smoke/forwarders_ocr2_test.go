package smoke

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
)

func TestForwarderOCR2Basic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithGeth().
		WithMockAdapter().
		WithCLNodeConfig(node.NewConfig(node.NewBaseConfig(),
			node.WithOCR2(),
			node.WithP2Pv2(),
		)).
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
		actions.AcceptAuthorizedReceiversOperator(t, operators[i], authorizedForwarders[i], []common.Address{workerNodeAddresses[i]}, env.EVMClient, env.ContractLoader)
		require.NoError(t, err, "Accepting Authorized Receivers on Operator shouldn't fail")
		err = actions.TrackForwarderLocal(env.EVMClient, authorizedForwarders[i], workerNodes[i], l)
		require.NoError(t, err, "failed to track forwarders")
		err = env.EVMClient.WaitForEvents()
		require.NoError(t, err, "Error waiting for events")
	}

	// Gather transmitters
	var transmitters []string
	for _, forwarderCommonAddress := range authorizedForwarders {
		transmitters = append(transmitters, forwarderCommonAddress.Hex())
	}

	ocrOffchainOptions := contracts.DefaultOffChainAggregatorOptions()
	ocrInstances, err := actions.DeployOCRv2Contracts(1, linkTokenContract, env.ContractDeployer, transmitters, env.EVMClient, ocrOffchainOptions)
	require.NoError(t, err, "Error deploying OCRv2 contracts with forwarders")
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	err = actions.CreateOCRv2JobsLocal(ocrInstances, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 5, env.EVMClient.GetChainID().Uint64(), true)
	require.NoError(t, err, "Error creating OCRv2 jobs with forwarders")
	err = env.EVMClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	ocrv2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
	require.NoError(t, err, "Error building OCRv2 config")
	ocrv2Config.Transmitters = authorizedForwarders

	err = actions.ConfigureOCRv2AggregatorContracts(env.EVMClient, ocrv2Config, ocrInstances)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

	err = actions.StartNewOCR2Round(1, ocrInstances, env.EVMClient, time.Minute*10, l)
	require.NoError(t, err)

	answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Getting latest answer from OCRv2 contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCRw contract to be 5 but got %d", answer.Int64())

	for i := 2; i <= 3; i++ {
		ocrRoundVal := (5 + i) % 10
		err = env.MockAdapter.SetAdapterBasedIntValuePath("ocr2", []string{http.MethodGet, http.MethodPost}, ocrRoundVal)
		require.NoError(t, err)
		err = actions.StartNewOCR2Round(int64(i), ocrInstances, env.EVMClient, time.Minute*10, l)
		require.NoError(t, err)

		answer, err = ocrInstances[0].GetLatestAnswer(context.Background())
		require.NoError(t, err, "Error getting latest OCRv2 answer")
		require.Equal(t, int64(ocrRoundVal), answer.Int64(), fmt.Sprintf("Expected latest answer from OCRv2 contract to be %d but got %d", ocrRoundVal, answer.Int64()))
	}
}
