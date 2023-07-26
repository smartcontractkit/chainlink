package smoke

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/docker"
	"github.com/smartcontractkit/chainlink/integration-tests/networks"
	"github.com/stretchr/testify/require"
	"strings"
)

// Tests a basic OCRv2 median feed
func TestOCRv2Basic(t *testing.T) {
	env, err := docker.NewChainlinkCluster(t, 6)
	require.NoError(t, err)
	clients, err := docker.ConnectClients(env)
	require.NoError(t, err)

	//chainClient, err := blockchain.NewEVMClient(clients.Networks[0], testEnvironment)
	//require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	//contractDeployer, err := contracts.NewContractDeployer(chainClient)
	//require.NoError(t, err, "Deploying contracts shouldn't fail")
	//
	//chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	//require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	bootstrapNode, workerNodes := clients.Chainlink[0], clients.Chainlink[1:]
	//mockServer, err := ctfClient.ConnectMockServer(testEnvironment)
	//require.NoError(t, err, "Creating mockserver clients shouldn't fail")
	//t.Cleanup(func() {
	//	err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
	//	require.NoError(t, err, "Error tearing down environment")
	//})
	clients.Networks[0].ParallelTransactions(true)

	linkToken, err := clients.NetworkDeployers[0].DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	err = actions.FundChainlinkNodes(workerNodes, clients.Networks[0], big.NewFloat(.05))
	require.NoError(t, err, "Error funding Chainlink nodes")

	// Gather transmitters
	var transmitters []string
	for _, node := range workerNodes {
		addr, err := node.PrimaryEthAddress()
		if err != nil {
			require.NoError(t, fmt.Errorf("error getting node's primary ETH address: %w", err))
		}
		transmitters = append(transmitters, addr)
	}

	aggregatorContracts, err := actions.DeployOCRv2Contracts(1, linkToken, clients.NetworkDeployers[0], transmitters, clients.Networks[0])
	require.NoError(t, err, "Error deploying OCRv2 aggregator contracts")

	err = actions.CreateOCRv2Jobs(aggregatorContracts, bootstrapNode, workerNodes, clients.Mockserver, "ocr2", 5, clients.Networks[0].GetChainID().Uint64(), false)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	ocrv2Config, err := actions.BuildMedianOCR2Config(workerNodes)
	require.NoError(t, err, "Error building OCRv2 config")

	err = actions.ConfigureOCRv2AggregatorContracts(clients.Networks[0], ocrv2Config, aggregatorContracts)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")
	time.Sleep(999 * time.Second)

	err = actions.StartNewOCR2Round(1, aggregatorContracts, clients.Networks[0], time.Minute*5)
	require.NoError(t, err, "Error starting new OCR2 round")
	roundData, err := aggregatorContracts[0].GetRound(context.Background(), big.NewInt(1))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 5 but got %d",
		roundData.Answer.Int64(),
	)

	err = clients.Mockserver.SetValuePath("ocr2", 10)
	require.NoError(t, err)
	err = actions.StartNewOCR2Round(2, aggregatorContracts, clients.Networks[0], time.Minute*5)
	require.NoError(t, err)

	roundData, err = aggregatorContracts[0].GetRound(context.Background(), big.NewInt(2))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 10 but got %d",
		roundData.Answer.Int64(),
	)
}

func setupOCR2Test(t *testing.T, forwardersEnabled bool) (
	testEnvironment *environment.Environment,
	testNetwork blockchain.EVMNetwork,
) {
	testNetwork = networks.SelectedNetwork
	evmConfig := ethereum.New(nil)
	if !testNetwork.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}

	var toml string
	if forwardersEnabled {
		toml = client.AddNetworkDetailedConfig(config.BaseOCR2Config, config.ForwarderNetworkDetailConfig, testNetwork)
	} else {
		toml = client.AddNetworksConfig(config.BaseOCR2Config, testNetwork)
	}

	chainlinkChart, err := chainlink.NewDeployment(6, map[string]interface{}{
		"toml": toml,
	})
	require.NoError(t, err, "Error creating chainlink deployment")

	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-ocr2-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelmCharts(chainlinkChart)
	err = testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")
	return testEnvironment, testNetwork
}
