package smoke

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	mercury_server "github.com/smartcontractkit/chainlink-env/pkg/helm/mercury-server"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/store/models"
	networks "github.com/smartcontractkit/chainlink/integration-tests"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func TestMercuryDONReportsGeneration(t *testing.T) {
	t.Parallel()

	testEnvironment, testNetwork, chainlinkNodes, mercuryServerInternalUrl := setupMercuryEnvironment(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err, "Error connecting to blockchain")

	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	// Setup random mock server response for mercury price feed
	mockServerClient, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Error connecting to mock server")
	err = mockServerClient.SetRandomValuePath("/variable")
	require.NoError(t, err, "Setting mockserver value path shouldn't fail")

	mercuryFeedID := "ETH-USD-Optimism-Goerli-1"
	var mercuryFeedIDBytes [32]byte
	copy(mercuryFeedIDBytes[:], mercuryFeedID)

	// Setup contracts
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")
	accessController, err := contractDeployer.DeployReadAccessController()
	require.NoError(t, err, "Error deploying ReadAccessController contract")
	verifierProxy, err := contractDeployer.DeployVerifierProxy(accessController.Address())
	require.NoError(t, err, "Error deploying VerifierProxy contract")
	verifier, err := contractDeployer.DeployVerifier(mercuryFeedIDBytes, verifierProxy.Address())
	require.NoError(t, err, "Error deploying Verifier contract")
	nodesWithoutBootstrap := chainlinkNodes[1:]
	ocrConfig := actions.BuildMercuryOCR2Config(t, nodesWithoutBootstrap)
	verifier.SetConfig(ocrConfig)
	latestConfigDetails, err := verifier.LatestConfigDetails()
	require.NoError(t, err, "Error getting Verifier.LatestConfigDetails()")
	verifierProxy.InitializeVerifier(latestConfigDetails.ConfigDigest, verifier.Address())

	// Create node jobs
	osTemplate := `
		ds1          [type=http method=GET url="%s" allowunrestrictednetworkaccess="true"];
		ds1_parse    [type=jsonparse path="data,result"];
		ds1_multiply [type=multiply times=100];
		ds1 -> ds1_parse -> ds1_multiply -> answer1;

		answer1 [type=median index=0 allowedFaults=4];
	`
	os := fmt.Sprintf(string(osTemplate), mockServerClient.Config.ClusterURL+"/variable")
	network := networks.SelectedNetwork
	createMercuryJobs(t, chainlinkNodes, verifier.Address(),
		mercuryFeedID, mercuryServerInternalUrl, network.ChainID, 0, os)

	// Set OCR2 config in the contract
	verifier.SetConfig(ocrConfig)

	// Wait for the DON to start generating reports
	time.Sleep(160 * time.Second)

	// Get a report from mercury server
	mercuryServerRemoteUrl := testEnvironment.URLs[mercury_server.URLsKey][1]
	msClient := client.NewMercuryServer(mercuryServerRemoteUrl)
	latestBlockNum, err := chainClient.LatestBlockNumber(context.Background())
	_ = latestBlockNum
	require.NoError(t, err, "Err getting latest block number")
	report, _, err := msClient.GetReports(mercuryFeedID, "5554794")
	require.NoError(t, err, "Error getting report from Mercury Server")
	require.NotEmpty(t, report.ChainlinkBlob, "Report response does not contain chainlinkBlob")
}

func setupMercuryEnvironment(t *testing.T) (*environment.Environment, blockchain.EVMNetwork, []*client.Chainlink, string) {
	testNetwork := networks.SelectedNetwork
	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}

	mercuryServerInternalUrl := fmt.Sprintf("http://%s:3000", mercury_server.URLsKey)
	secretsToml := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		"[[Mercury.Credentials]]",
		fmt.Sprintf(`URL = "%s/reports"`, mercuryServerInternalUrl),
		`Username = "node"`,
		`Password = "nodepass"`,
	)

	testEnvironment := environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-mercury-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(mercury_server.New(map[string]interface{}{
			"imageRepo": os.Getenv("MERCURY_SERVER_IMAGE"),
			"imageTag":  os.Getenv("MERCURY_SERVER_TAG"),
		})).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "5",
			"toml": client.AddNetworksConfig(
				config.BaseMercuryTomlConfig,
				testNetwork),
			"secretsToml": secretsToml,
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	return testEnvironment, testNetwork, chainlinkNodes, mercuryServerInternalUrl
}

func createMercuryJobs(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
	contractID string,
	feedID string,
	mercuryServerUrl string,
	chainID int64,
	keyIndex int,
	observationSource string,
) {
	bootstrapNode := chainlinkNodes[0]
	bootstrapNode.RemoteIP()
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	require.NoError(t, err, "Shouldn't fail reading P2P keys from bootstrap node")
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID

	bootstrapSpec := &client.OCR2TaskJobSpec{
		Name:    "ocr2 bootstrap node",
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: contractID,
			Relay:      "evm",
			RelayConfig: map[string]interface{}{
				"chainID": int(chainID),
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
		},
	}
	_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
	require.NoError(t, err, "Shouldn't fail creating bootstrap job on bootstrap node")
	P2Pv2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PId, bootstrapNode.RemoteIP(), 6690)

	for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
		nodeTransmitterAddress, err := chainlinkNodes[nodeIndex].EthAddresses()
		require.NoError(t, err, "Shouldn't fail getting primary ETH address from OCR node %d", nodeIndex+1)
		nodeOCRKeys, err := chainlinkNodes[nodeIndex].MustReadOCR2Keys()
		require.NoError(t, err, "Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
		var nodeOCRKeyId []string
		for _, key := range nodeOCRKeys.Data {
			if key.Attributes.ChainType == string(chaintype.EVM) {
				nodeOCRKeyId = append(nodeOCRKeyId, key.ID)
				break
			}
		}

		// Convert feedID to hex with zeros appended
		// Example: 0x4554482d5553442d4f7074696d69736d2d476f65726c692d3100000000000000
		var feedIDBytes [32]byte
		copy(feedIDBytes[:], feedID)
		feedIDHex := fmt.Sprintf("0x%x", feedIDBytes)

		autoOCR2JobSpec := client.OCR2TaskJobSpec{
			Name:            "ocr2",
			JobType:         "offchainreporting2",
			MaxTaskDuration: "1s",
			OCR2OracleSpec: job.OCR2OracleSpec{
				PluginType: "median",
				PluginConfig: map[string]interface{}{
					"juelsPerFeeCoinSource": `"""
						bn1          [type=ethgetblock];
						bn1_lookup   [type=lookup key="number"];
						bn1 -> bn1_lookup;
					"""`,
				},
				Relay: "evm",
				RelayConfig: map[string]interface{}{
					"chainID": int(chainID),
				},
				RelayConfigMercuryConfig: map[string]interface{}{
					"feedID": feedIDHex,
					"url":    fmt.Sprintf("%s/reports", mercuryServerUrl),
				},
				ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
				ContractID:                        contractID,                                        // registryAddr
				OCRKeyBundleID:                    null.StringFrom(nodeOCRKeyId[keyIndex]),           // get node ocr2config.ID
				TransmitterID:                     null.StringFrom(nodeTransmitterAddress[keyIndex]), // node addr
				P2PV2Bootstrappers:                pq.StringArray{P2Pv2Bootstrapper},                 // bootstrap node key and address <p2p-key>@bootstrap:8000
			},
			ObservationSource: observationSource,
		}

		_, err = chainlinkNodes[nodeIndex].MustCreateJob(&autoOCR2JobSpec)
		require.NoError(t, err, "Shouldn't fail creating OCR Task job on OCR node %d", nodeIndex+1)
	}
	log.Info().Msg("Done creating OCR automation jobs")
}
