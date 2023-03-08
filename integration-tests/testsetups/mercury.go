package testsetups

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
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

type csaKey struct {
	NodeName    string `json:"nodeName"`
	NodeAddress string `json:"nodeAddress"`
	PublicKey   string `json:"publicKey"`
}

type oracle struct {
	Id                    string   `json:"id"`
	Website               string   `json:"website"`
	Name                  string   `json:"name"`
	Status                string   `json:"status"`
	NodeAddress           []string `json:"nodeAddress"`
	OracleAddress         string   `json:"oracleAddress"`
	CsaKeys               []csaKey `json:"csaKeys"`
	Ocr2ConfigPublicKey   []string `json:"ocr2ConfigPublicKey"`
	Ocr2OffchainPublicKey []string `json:"ocr2OffchainPublicKey"`
	Ocr2OnchainPublicKey  []string `json:"ocr2OnchainPublicKey"`
}

func ValidateReport(r map[string]interface{}) error {
	feedIdInterface, ok := r["feedId"]
	if !ok {
		return errors.Errorf("unpacked report has no 'feedId'")
	}
	feedID, ok := feedIdInterface.([32]byte)
	if !ok {
		return errors.Errorf("cannot cast feedId to [32]byte, type is %T", feedID)
	}
	log.Trace().Str("FeedID", string(feedID[:])).Msg("Feed ID")

	priceInterface, ok := r["median"]
	if !ok {
		return errors.Errorf("unpacked report has no 'median'")
	}
	medianPrice, ok := priceInterface.(*big.Int)
	if !ok {
		return errors.Errorf("cannot cast median to *big.Int, type is %T", medianPrice)
	}
	log.Trace().Int64("Price", medianPrice.Int64()).Msg("Median price")

	observationsBlockNumberInterface, ok := r["observationsBlocknumber"]
	if !ok {
		return errors.Errorf("unpacked report has no 'observationsBlocknumber'")
	}
	observationsBlockNumber, ok := observationsBlockNumberInterface.(uint64)
	if !ok {
		return errors.Errorf("cannot cast observationsBlocknumber to uint64, type is %T", observationsBlockNumber)
	}
	log.Trace().Uint64("Block", observationsBlockNumber).Msg("Observation block number")

	observationsTimestampInterface, ok := r["observationsTimestamp"]
	if !ok {
		return errors.Errorf("unpacked report has no 'observationsTimestamp'")
	}
	observationsTimestamp, ok := observationsTimestampInterface.(uint32)
	if !ok {
		return errors.Errorf("cannot cast observationsTimestamp to uint32, type is %T", observationsTimestamp)
	}
	log.Trace().Uint32("Timestamp", observationsTimestamp).Msg("Observation timestamp")

	return nil
}

func generateEd25519Keys() (string, string, error) {
	var (
		err  error
		pub  ed25519.PublicKey
		priv ed25519.PrivateKey
	)
	pub, priv, err = ed25519.GenerateKey(rand.Reader)
	return hex.EncodeToString(priv), hex.EncodeToString(pub), err
}

func SetupMercuryServer(
	t *testing.T,
	testEnv *environment.Environment,
	dbSettings map[string]interface{},
	serverSettings map[string]interface{},
) string {
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnv)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	// Build rpc config for Mercury server
	var msRpcNodesConf []*oracle
	for i, chainlinkNode := range chainlinkNodes {
		nodeName := fmt.Sprint(i)
		csaKeys, resp, err := chainlinkNode.ReadCSAKeys()
		require.NoError(t, err)
		csaKeyId := csaKeys.Data[0].ID
		ocr2Keys, resp, err := chainlinkNode.ReadOCR2Keys()
		_ = ocr2Keys
		_ = resp
		require.NoError(t, err)
		var ocr2Config client.OCR2KeyAttributes
		for _, key := range ocr2Keys.Data {
			if key.Attributes.ChainType == string(chaintype.EVM) {
				ocr2Config = key.Attributes
				break
			}
		}
		ocr2ConfigPublicKey := strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_")
		ocr2OffchainPublicKey := strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_")
		ocr2OnchainPublicKey := strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_")

		node := &oracle{
			Id:            fmt.Sprint(i),
			Name:          nodeName,
			Status:        "active",
			NodeAddress:   []string{nodeAddress},
			OracleAddress: "0x0000000000000000000000000000000000000000",
			CsaKeys: []csaKey{
				{
					NodeName:    nodeName,
					NodeAddress: nodeAddress,
					PublicKey:   csaKeyId,
				},
			},
			Ocr2ConfigPublicKey:   []string{ocr2ConfigPublicKey},
			Ocr2OffchainPublicKey: []string{ocr2OffchainPublicKey},
			Ocr2OnchainPublicKey:  []string{ocr2OnchainPublicKey},
		}
		msRpcNodesConf = append(msRpcNodesConf, node)
	}
	rpcNodesJsonConf, _ := json.Marshal(msRpcNodesConf)
	// result := []interface{}{}
	// err = json.Unmarshal([]byte(msRpcNodesJsonConf), &result)
	// err = mapstructure.Decode(msRpcNodesConf[0], &result)

	// Generate keys for Mercury RPC server
	rpcPrivKey, rpcPubKey, err := generateEd25519Keys()
	require.NoError(t, err)

	settings := map[string]interface{}{
		"imageRepo":     os.Getenv("MERCURY_SERVER_IMAGE"),
		"imageTag":      os.Getenv("MERCURY_SERVER_TAG"),
		"rpcPrivateKey": rpcPrivKey,
		"rpcNodesConf":  string(rpcNodesJsonConf),
		"prometheus":    "true",
	}

	if dbSettings != nil {
		settings["db"] = dbSettings
	}
	if serverSettings != nil {
		settings["resources"] = serverSettings
	}

	testEnv.AddHelm(mercury_server.New(settings)).Run()

	return rpcPubKey
}

func SetupMercuryEnv(t *testing.T, dbSettings map[string]interface{}, serverResources map[string]interface{}) (
	*environment.Environment, bool, blockchain.EVMNetwork, []*client.Chainlink, string,
	blockchain.EVMClient, *ctfClient.MockserverClient, *client.MercuryServer, string) {
	testNetwork := networks.SelectedNetwork
	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}

	secretsToml := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		"[[Mercury.Credentials]]",
		fmt.Sprintf(`URL = "%s"`, fmt.Sprintf("%s:1338", mercury_server.URLsKey)),
		`Username = "node"`,
		`Password = "nodepass"`,
	)

	testEnvironment := environment.New(&environment.Config{
		TTL:             12 * time.Hour,
		NamespacePrefix: fmt.Sprintf("smoke-mercury-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "5",
			"toml": client.AddNetworksConfig(
				config.BaseMercuryTomlConfig,
				testNetwork),
			"secretsToml": secretsToml,
			"prometheus":  "true",
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")

	msRpcPubKey := SetupMercuryServer(t, testEnvironment, dbSettings, serverResources)

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
		mercuryServerRemoteUrl, evmClient, mockserverClient, mercuryServerClient, msRpcPubKey
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
	mercuryServerPubKey string,
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
		csaKeys, _, err := chainlinkNodes[nodeIndex].ReadCSAKeys()
		require.NoError(t, err)
		csaKeyId := csaKeys.Data[0].ID
		_ = csaKeyId

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
					"clientPrivKeyID": csaKeyId,
					"serverPubKey":    fmt.Sprintf("0x%s", mercuryServerPubKey),
					"feedID":          feedIdHex,
					"url":             fmt.Sprintf("%s:1338", mercury_server.URLsKey),
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
