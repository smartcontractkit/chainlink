package smoke

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func TestForwarderOCRBasic(t *testing.T) {
	t.Parallel()
	testEnvironment, testNetwork := setupForwarderOCRTest(t)
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
	nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkNodes)
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

	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.05))
	require.NoError(t, err, "Error funding Chainlink nodes")

	operators, authorizedForwarders, _ := actions.DeployForwarderContracts(
		t, contractDeployer, linkTokenContract, chainClient, len(chainlinkNodes)-1,
	)
	forwarderNodes := chainlinkNodes[1:]
	forwarderNodesAddresses := nodeAddresses[1:]
	for i := range forwarderNodes {
		actions.AcceptAuthorizedReceiversOperator(
			t, operators[i], authorizedForwarders[i], []common.Address{forwarderNodesAddresses[i]}, chainClient, contractLoader,
		)
		require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
		actions.TrackForwarder(t, chainClient, authorizedForwarders[i], forwarderNodes[i])
		err = chainClient.WaitForEvents()
	}
	ocrInstances := actions.DeployOCRContractsForwarderFlow(
		t, 1, linkTokenContract, contractDeployer, chainlinkNodes, authorizedForwarders, chainClient,
	)
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	err = actions.SetAllAdapterResponsesToTheSameValue(5, ocrInstances, chainlinkNodes, mockServer)
	require.NoError(t, err)
	actions.CreateOCRJobsWithForwarder(t, ocrInstances, chainlinkNodes, mockServer)
	err = actions.StartNewRound(1, ocrInstances, chainClient)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	answer, err := ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

	err = actions.SetAllAdapterResponsesToTheSameValue(10, ocrInstances, chainlinkNodes, mockServer)
	require.NoError(t, err)
	err = actions.StartNewRound(2, ocrInstances, chainClient)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	answer, err = ocrInstances[0].GetLatestAnswer(context.Background())
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}

func setupForwarderOCRTest(t *testing.T) (testEnvironment *environment.Environment, testNetwork blockchain.EVMNetwork) {
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

[Feature]
LogPoller = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690`
	networkDetailTOML := `[EVM.Transactions]
ForwardersEnabled = true`
	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-ocr-forwarder-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"toml":     client.AddNetworkDetailedConfig(baseTOML, networkDetailTOML, testNetwork),
			"replicas": 6,
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")
	return testEnvironment, testNetwork
}
