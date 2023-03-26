package performance

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/stretchr/testify/require"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestOCRBasic(t *testing.T) {
	t.Parallel()
	testEnvironment, testNetwork := setupOCRTest(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	mockServer, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Creating mockserver clients shouldn't fail")

	chainClient.ParallelTransactions(true)

	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.05))
	require.NoError(t, err, "Error funding Chainlink nodes")

	ocrInstances, err := actions.DeployOCRContracts(1, linkTokenContract, contractDeployer, chainlinkNodes, chainClient)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	profileFunction := func(chainlinkNode *client.Chainlink) {
		err = actions.SetAllAdapterResponsesToTheSameValue(5, ocrInstances, chainlinkNodes, mockServer)
		require.NoError(t, err)
		err = actions.CreateOCRJobs(ocrInstances, chainlinkNodes, mockServer)
		require.NoError(t, err)
		err = actions.StartNewRound(1, ocrInstances, chainClient)
		require.NoError(t, err)

		answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
		require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
		require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

		err = actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstances, chainlinkNodes, mockServer)
		require.NoError(t, err)
		err = actions.StartNewRound(2, ocrInstances, chainClient)
		require.NoError(t, err)

		answer, err = ocrInstances[0].GetLatestAnswer(context.Background())
		require.NoError(t, err, "Error getting latest OCR answer")
		require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
	}
	profileTest := testsetups.NewChainlinkProfileTest(testsetups.ChainlinkProfileTestInputs{
		ProfileFunction: profileFunction,
		ProfileDuration: time.Minute,
		ChainlinkNodes:  chainlinkNodes,
	})
	t.Cleanup(func() {
		CleanupPerformanceTest(t, testEnvironment, chainlinkNodes, profileTest.TestReporter, chainClient)
	})
	profileTest.Setup(testEnvironment)
	profileTest.Run()
}

func setupOCRTest(t *testing.T) (testEnvironment *environment.Environment, testNetwork blockchain.EVMNetwork) {
	testNetwork = networks.SelectedNetwork
	evmConfig := ethereum.New(nil)
	if !testNetwork.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}
	baseTOML := `[OCR]
Enabled = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690`
	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("performance-ocr-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"toml":     client.AddNetworksConfig(baseTOML, testNetwork),
			"replicas": 6,
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")
	return testEnvironment, testNetwork
}
