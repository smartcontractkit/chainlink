package soak

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var (
	automationBaseTOML = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[Keeper]
TurnLookBack = 0

[Keeper.Registry]
SyncInterval = '5m'
PerformGasOverhead = 150_000

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`
	defaultOCRRegistryConfig = contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(200000000),
		FlatFeeMicroLINK:     uint32(0),
		BlockCountPerTurn:    big.NewInt(10),
		CheckGasLimit:        uint32(2500000),
		StalenessSeconds:     big.NewInt(90000),
		GasCeilingMultiplier: uint16(1),
		MinUpkeepSpend:       big.NewInt(0),
		MaxPerformGas:        uint32(5000000),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
		MaxCheckDataSize:     uint32(5000),
		MaxPerformDataSize:   uint32(5000),
	}
	automationDefaultRegistryConfig = contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(200000000),
		FlatFeeMicroLINK:     uint32(0),
		BlockCountPerTurn:    big.NewInt(10),
		CheckGasLimit:        uint32(2500000),
		StalenessSeconds:     big.NewInt(90000),
		GasCeilingMultiplier: uint16(1),
		MinUpkeepSpend:       big.NewInt(0),
		MaxPerformGas:        uint32(5000000),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
		MaxCheckDataSize:     uint32(5000),
		MaxPerformDataSize:   uint32(5000),
	}
)

func TestOCR2Soak(t *testing.T) {
	t.Parallel()

	network := networks.SelectedNetwork
	evmConfig := eth.New(nil)
	if !network.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
	}

	testEnvironment := environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("soak-ocr2-basic-%s", strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")),
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "5",
			"toml":     client.AddNetworksConfig(automationBaseTOML, network),
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error setting up test environment")

	chainClient, err := blockchain.NewEVMClient(network, testEnvironment)
	require.NoError(t, err, "Error connecting to blockchain")
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Error building contract deployer")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	chainClient.ParallelTransactions(true)

	// Register cleanup for any test
	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	txCost, err := chainClient.EstimateCostForChainlinkOperations(1000)
	require.NoError(t, err, "Error estimating cost for Chainlink Operations")
	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, txCost)
	require.NoError(t, err, "Error funding Chainlink nodes")

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Error deploying LINK token")

	// Deploy access controller contracts required by the aggregator contract
	accessController, err := contractDeployer.DeployAccessController()
	require.NoError(t, err, "Error deploying billing access controller contract")
	// requesterAccessController, err := contractDeployer.DeployAccessController()
	// require.NoError(t, err, "Error deploying requester access controller contract")

	// for i, cl := range chainlinkNodes {
	// 	clAddress, err := cl.PrimaryEthAddress()
	// 	require.NoError(t, err, "Shouldn't fail getting primary ETH address from OCR node: index %d", i)
	// 	err = writeAccessController.AddAccess(clAddress)
	// 	require.NoError(t, err, "Error adding write access for OCR node: index %d", i)
	// }

	// Deploy aggregator contract
	ocr2Aggregator, err := contractDeployer.DeployOCR2Aggregator(
		linkToken.Address(),
		accessController.Address(),
		accessController.Address(),
	)
	require.NoError(t, err, "Error deploying ocr2 aggregator contract")

	// Setup mock server response
	mockServerClient, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Error connecting to mock server")
	err = mockServerClient.SetValuePath("/variable", 5)
	require.NoError(t, err, "Setting mockserver value path shouldn't fail")

	// Create jobs
	osTemplate := `
		ds1          [type=http method=GET url="%s" allowunrestrictednetworkaccess="true"];
		ds1_parse    [type=jsonparse path="answer"];
		ds1_multiply [type=multiply times=100];
		ds1 -> ds1_parse -> ds1_multiply;
	`
	os := fmt.Sprintf(string(osTemplate), mockServerClient.Config.ClusterURL+"/variable")
	actions.CreateOCR2Jobs(t, chainlinkNodes, ocr2Aggregator.Address(), network.ChainID, 0, os)

	// Set the aggregator configuration

	// TODO: check if KeyBundleID is set correctly

	nodesWithoutBootstrap := chainlinkNodes[1:]
	ocrConfig := actions.BuildOCR2AggregatorConfig(t, nodesWithoutBootstrap, 5*time.Second)
	fmt.Printf("%+v", ocrConfig)
	err = ocr2Aggregator.SetConfig(ocrConfig)
	require.NoError(t, err, "OCR2 aggregator config should be be set successfully")
	require.NoError(t, chainClient.WaitForEvents(), "Waiting for config to be set")
}
