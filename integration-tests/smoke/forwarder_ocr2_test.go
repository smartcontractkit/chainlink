package smoke

import (
	"context"
	"log"
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
		require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
		actions.TrackForwarder(t, chainClient, authorizedForwarders[i], workerNodes[i])
		log.Println("Node name ", workerNodes[i].Name())
		log.Println("Forwarder addr  ", authorizedForwarders[i].Hex())
		fwdrs, _, err := workerNodes[i].GetForwarders()
		log.Println("fwdrs   ", fwdrs)
		log.Println("fwdrserr   ", err)
		err = chainClient.WaitForEvents()
		require.NoError(t, err, "Error waiting for events")
	}

	ocrInstances, err := actions.DeployOCRv2ContractsForwardersFlow(1, linkTokenContract, contractDeployer, authorizedForwarders, chainClient)
	require.NoError(t, err, "Error deploying OCRv2 contracts with forwarders")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	ocrv2Config, err := actions.BuildMedianOCR2Config(workerNodes)
	require.NoError(t, err, "Error building OCRv2 config")

	err = actions.CreateOCRv2JobsWithForwarder(ocrInstances, bootstrapNode, workerNodes, mockServer, "ocr2", 5, chainClient.GetChainID().Uint64())
	require.NoError(t, err, "Error creating OCRv2 jobs with forwarders")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	err = actions.ConfigureOCRv2AggregatorContracts(chainClient, ocrv2Config, ocrInstances)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	err = actions.StartNewOCR2Round(1, ocrInstances, chainClient, time.Minute*10)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Getting latest answer from OCRv2 contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCRw contract to be 5 but got %d", answer.Int64())

	err = mockServer.SetValuePath("ocr2", 10)
	require.NoError(t, err)
	err = actions.StartNewOCR2Round(2, ocrInstances, chainClient, time.Minute*10)
	require.NoError(t, err)

	answer, err = ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Error getting latest OCRv2 answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCRv2 contract to be 10 but got %d", answer.Int64())

}
