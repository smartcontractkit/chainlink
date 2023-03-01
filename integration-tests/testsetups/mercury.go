package testsetups

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	cl_env_config "github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	mercury_server "github.com/smartcontractkit/chainlink-env/pkg/helm/mercury-server"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/store/models"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"
)

func SetupMercuryEnv(t *testing.T) (
	*environment.Environment, bool, blockchain.EVMNetwork, []*client.Chainlink, string, string,
	blockchain.EVMClient, *ctfClient.MockserverClient, *client.MercuryServer) {
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
		// TTL:             1 * time.Hour,
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

	evmClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err, "Error connecting to blockchain")

	isExistingTestEnv := os.Getenv(cl_env_config.EnvVarNamespace) != "" && os.Getenv(cl_env_config.EnvVarNoManifestUpdate) == "true"

	// Setup random mock server response for mercury price feed
	mockserverClient, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Error connecting to mock server")

	mercuryServerRemoteUrl := testEnvironment.URLs[mercury_server.URLsKey][1]
	mercuryServerClient := client.NewMercuryServer(mercuryServerRemoteUrl)

	t.Cleanup(func() {
		if isExistingTestEnv {
			log.Info().Msg("Do not tear down existing environment")
		} else {
			err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.PanicLevel, evmClient)
			require.NoError(t, err, "Error tearing down environment")
		}
	})

	return testEnvironment, isExistingTestEnv, testNetwork, chainlinkNodes,
		mercuryServerInternalUrl, mercuryServerRemoteUrl, evmClient, mockserverClient, mercuryServerClient
}

func SetupMercuryContracts(t *testing.T, evmClient blockchain.EVMClient, mercuryRemoteUrl string, feedId string, ocrConfig contracts.OCRConfig) (contracts.Verifier, contracts.VerifierProxy, contracts.ReadAccessController, contracts.Exchanger) {
	contractDeployer, err := contracts.NewContractDeployer(evmClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")

	accessController, err := contractDeployer.DeployReadAccessController()
	require.NoError(t, err, "Error deploying ReadAccessController contract")

	// verifierProxy, err := contractDeployer.DeployVerifierProxy(accessController.Address())
	verifierProxy, err := contractDeployer.DeployVerifierProxy("0x0")
	require.NoError(t, err, "Error deploying VerifierProxy contract")

	var feedIdBytes [32]byte
	copy(feedIdBytes[:], feedId)
	verifier, err := contractDeployer.DeployVerifier(feedIdBytes, verifierProxy.Address())
	require.NoError(t, err, "Error deploying Verifier contract")

	verifier.SetConfig(ocrConfig)
	latestConfigDetails, err := verifier.LatestConfigDetails()
	require.NoError(t, err, "Error getting Verifier.LatestConfigDetails()")
	log.Info().Msgf("Latest config digest: %x", latestConfigDetails.ConfigDigest)
	log.Info().Msgf("Latest config details: %v", latestConfigDetails)

	verifierProxy.InitializeVerifier(latestConfigDetails.ConfigDigest, verifier.Address())

	return verifier, verifierProxy, accessController, nil
}

func SetupMercuryNodeJobs(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
	mockserverClient *ctfClient.MockserverClient,
	contractID string,
	feedID string,
	mercuryServerUrl string,
	chainID int64,
	keyIndex int,
) {
	err := mockserverClient.SetRandomValuePath("/variable")
	require.NoError(t, err, "Setting mockserver value path shouldn't fail")

	osTemplate := `
	ds1          [type=http method=GET url="%s" allowunrestrictednetworkaccess="true"];
	ds1_parse    [type=jsonparse path="data,result"];
	ds1_multiply [type=multiply times=100];
	ds1 -> ds1_parse -> ds1_multiply -> answer1;

	answer1 [type=median index=0 allowedFaults=4];
`
	observationSource := fmt.Sprintf(string(osTemplate), mockserverClient.Config.ClusterURL+"/variable")

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
		var feedIdBytes [32]byte
		copy(feedIdBytes[:], feedID)
		feedIdHex := fmt.Sprintf("0x%x", feedIdBytes)
		log.Info().Msgf("Setup feedID, string: %s, hex: %s", feedID, feedIdHex)

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
					"feedID": feedIdHex,
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

func BuildMercuryOCR2Config(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
) contracts.OCRConfig {
	onchainConfig, err := (median.StandardOnchainConfigCodec{}).Encode(median.OnchainConfig{median.MinValue(), median.MaxValue()})
	require.NoError(t, err, "Shouldn't fail encoding config")

	alphaPPB := uint64(1000)

	return actions.BuildGeneralOCR2Config(
		t,
		chainlinkNodes,
		2*time.Second,        // deltaProgress time.Duration,
		20*time.Second,       // deltaResend time.Duration,
		100*time.Millisecond, // deltaRound time.Duration,
		0,                    // deltaGrace time.Duration,
		1*time.Minute,        // deltaStage time.Duration,
		100,                  // rMax uint8,
		[]int{len(chainlinkNodes)},
		median.OffchainConfig{
			false,
			alphaPPB,
			false,
			alphaPPB,
			0,
		}.Encode(),
		0*time.Millisecond,   // maxDurationQuery time.Duration,
		250*time.Millisecond, // maxDurationObservation time.Duration,
		250*time.Millisecond, // maxDurationReport time.Duration,
		250*time.Millisecond, // maxDurationShouldAcceptFinalizedReport time.Duration,
		250*time.Millisecond, // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,                    // f int,
		onchainConfig,
	)
}
