package smoke

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func TestForwarderOCR2Basic(t *testing.T) {
	t.Parallel()
	testEnvironment, testNetwork := setupOCR2Test(t, true)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")
	contractLoader, err := contracts.NewContractLoader(chainClient)
	require.NoError(t, err, "Loading contracts shouldn't fail")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	bootstrapNode, workerNodes := chainlinkNodes[0], chainlinkNodes[1:]
	workerNodeAddresses, err := actions.ChainlinkNodeAddresses(workerNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
	mockServer, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Creating mockserver clients shouldn't fail")
	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})
	chainClient.ParallelTransactions(true)

	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	err = actions.FundChainlinkNodes(workerNodes, chainClient, big.NewFloat(.05))
	require.NoError(t, err, "Error funding Chainlink nodes")

	operators, authorizedForwarders, _ := actions.DeployForwarderContracts(
		t, contractDeployer, linkTokenContract, chainClient, len(workerNodes),
	)

	for i := range workerNodes {
		actions.AcceptAuthorizedReceiversOperator(t, operators[i], authorizedForwarders[i], []common.Address{workerNodeAddresses[i]}, chainClient, contractLoader)
		require.NoError(t, err, "Accepting Authorized Receivers on Operator shouldn't fail")
		actions.TrackForwarder(t, chainClient, authorizedForwarders[i], workerNodes[i])
		err = chainClient.WaitForEvents()

		require.NoError(t, err, "Error waiting for events")
	}

	// Gather transmitters
	var transmitters []string
	for _, forwarderCommonAddress := range authorizedForwarders {
		transmitters = append(transmitters, forwarderCommonAddress.Hex())
	}

	ocrInstances, err := actions.DeployOCRv2Contracts(1, linkTokenContract, contractDeployer, transmitters, chainClient)

	require.NoError(t, err, "Error deploying OCRv2 contracts with forwarders")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	err = actions.CreateOCRv2Jobs(ocrInstances, bootstrapNode, workerNodes, mockServer, "ocr2", 5, chainClient.GetChainID().Uint64(), true)
	require.NoError(t, err, "Error creating OCRv2 jobs with forwarders")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	ocrv2Config, err := actions.BuildMedianOCR2Config(workerNodes)
	require.NoError(t, err, "Error building OCRv2 config")
	ocrv2Config.Transmitters = authorizedForwarders

	err = actions.ConfigureOCRv2AggregatorContracts(chainClient, ocrv2Config, ocrInstances)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

	err = actions.StartNewOCR2Round(1, ocrInstances, chainClient, time.Minute*10)
	require.NoError(t, err)

	answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Getting latest answer from OCRv2 contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCRw contract to be 5 but got %d", answer.Int64())

	for i := 2; i <= 100; i++ {
		ocrRoundVal := (5 + i) % 10
		err = mockServer.SetValuePath("ocr2", ocrRoundVal)
		require.NoError(t, err)
		err = actions.StartNewOCR2Round(int64(i), ocrInstances, chainClient, time.Minute*10)
		require.NoError(t, err)

		answer, err = ocrInstances[0].GetLatestAnswer(context.Background())
		require.NoError(t, err, "Error getting latest OCRv2 answer")
		require.Equal(t, int64(ocrRoundVal), answer.Int64(), fmt.Sprintf("Expected latest answer from OCRv2 contract to be %d but got %d", ocrRoundVal, answer.Int64()))
	}
}
