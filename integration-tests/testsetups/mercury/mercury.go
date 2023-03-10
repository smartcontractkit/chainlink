package mercury

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	mshelm "github.com/smartcontractkit/chainlink-env/pkg/helm/mercury-server"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/store/models"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	testconfig "github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"
)

type MercuryTestEnv struct {
	T                     *testing.T
	Config                mercuryTestConfig
	Env                   *environment.Environment
	ChainlinkNodes        []*client.Chainlink
	MSClient              *client.MercuryServer // Mercury server client authenticated with admin role
	IsExistingTestEnv     bool                  // true if config in MERCURY_ENV_CONFIG_PATH contains namespace
	KeepEnv               bool                  // Set via MERCURY_KEEP_ENV=true env
	EnvTTL                time.Duration         // Set via MERCURY_ENV_TTL_MINS env
	EvmClient             blockchain.EVMClient
	VerifierContract      contracts.Verifier
	VerifierProxyContract contracts.VerifierProxy
	ExchangerContract     contracts.Exchanger
}

type mercuryTestConfig struct {
	K8Namespace          string `json:"k8Namespace"`
	ChainId              int64  `json:"chainId"`
	FeedId               string `json:"feedId"`
	VerifierAddress      string `json:"verifierAddress"`
	VerifierProxyAddress string `json:"verifierProxyAddress"`
	ExchangerAddress     string `json:"exchangerAddress"`
	MSRemoteUrl          string `json:"mercuryServerRemoteUrl"`
	MSLocalUrl           string `json:"mercuryServerLocalUrl"`
	MSAdminId            string `json:"mercuryServerAdminId"`
	MSAdminKey           string `json:"mercuryServerAdminKey"`
	MSAdminEncryptedKey  string `json:"mercuryServerAdminEncryptedKey"`
}

func configFromEnv() mercuryTestConfig {
	c := mercuryTestConfig{}
	if os.Getenv("MERCURY_ENV_CONFIG_PATH") != "" {
		jsonFile, err := os.Open(os.Getenv("MERCURY_ENV_CONFIG_PATH"))
		if err != nil {
			return c
		}
		defer jsonFile.Close()
		b, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(b, &c)
		if err == nil {
			log.Info().Msgf("Using existing mercury env config from: %s\n%s",
				os.Getenv("MERCURY_ENV_CONFIG_PATH"), c.Json())
		}
	}
	return c
}

func NewMercuryTestEnv(t *testing.T) *MercuryTestEnv {
	testEnv := &MercuryTestEnv{}

	// Re-use existing env when MERCURY_ENV_CONFIG_PATH env with json c specified
	c := configFromEnv()
	if c.K8Namespace != "" {
		testEnv.IsExistingTestEnv = true
		// Set env variables for chainlink-env to reuse existing deployment
		os.Setenv("ENV_NAMESPACE", c.K8Namespace)
		os.Setenv("NO_MANIFEST_UPDATE", "true")
	} else {
		testEnv.IsExistingTestEnv = false
	}

	testEnv.T = t
	testEnv.KeepEnv = os.Getenv("MERCURY_KEEP_ENV") == "true"
	envTTL, err := strconv.ParseUint(os.Getenv("MERCURY_ENV_TTL_MINS"), 10, 64)
	if err == nil {
		testEnv.EnvTTL = time.Duration(envTTL) * time.Minute
	} else {
		testEnv.EnvTTL = 20 * time.Minute
	}
	testEnv.Config = c
	testEnv.Config.MSAdminId = os.Getenv("MS_DATABASE_FIRST_ADMIN_ID")
	testEnv.Config.MSAdminKey = os.Getenv("MS_DATABASE_FIRST_ADMIN_KEY")
	testEnv.Config.MSAdminEncryptedKey = os.Getenv("MS_DATABASE_FIRST_ADMIN_ENCRYPTED_KEY")

	return testEnv
}

func (c *mercuryTestConfig) Json() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *mercuryTestConfig) Save() (string, error) {
	// Create mercury env log dir if necessary
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	confDir := fmt.Sprintf("%s/logs", pwd)
	if _, err := os.Stat(confDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(confDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Save mercury env config to disk
	confId, err := uuid.NewV4()
	if err != nil {
		return "", nil
	}
	confPath := fmt.Sprintf("%s/%s.json", confDir, confId)
	f, _ := json.MarshalIndent(c, "", " ")
	err = ioutil.WriteFile(confPath, f, 0644)

	return confPath, err
}

// Setup DON, Mercury Server and all mercury contracts
func (e *MercuryTestEnv) SetupFullMercuryEnv(dbSettings map[string]interface{}, serverResources map[string]interface{}) {
	testNetwork := networks.SelectedNetwork
	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}

	// Fail when existing env is different than current chain
	if e.IsExistingTestEnv {
		require.Equal(e.T, e.Config.ChainId, testNetwork.ChainID,
			"Chain set in SELECTED_NETWORKS is different than chain id set in config provided by MERCURY_ENV_CONFIG_PATH")
	}
	e.Config.ChainId = testNetwork.ChainID

	env := e.SetupDON(e.T, testNetwork, evmConfig)
	e.Env = env
	e.Config.K8Namespace = env.Cfg.Namespace

	chainlinkNodes, err := client.ConnectChainlinkNodes(env)
	e.ChainlinkNodes = chainlinkNodes
	require.NoError(e.T, err, "Error connecting to Chainlink nodes")

	evmClient, err := blockchain.NewEVMClient(testNetwork, env)
	e.EvmClient = evmClient
	require.NoError(e.T, err, "Error connecting to blockchain")

	e.T.Cleanup(func() {
		if e.KeepEnv {
			envConfFile, err := e.Config.Save()
			require.NoError(e.T, err, "Could not save mercury env conf file")
			log.Info().Msgf("Keep mercury environment running."+
				" Chain: %d. Initial TTL: %s", e.Config.ChainId, env.Cfg.TTL)
			log.Info().Msgf("To reuse this env in next test on chain %d, set:\n"+
				"\"MERCURY_ENV_CONFIG_PATH\"=\"%s\"", e.Config.ChainId, envConfFile)
		} else {
			log.Info().Msgf("Destroy this mercury env because MERCURY_KEEP_ENV not set to \"true\"")
			err := actions.TeardownSuite(e.T, env, utils.ProjectRoot,
				chainlinkNodes, nil, zapcore.PanicLevel, evmClient)
			require.NoError(e.T, err, "Error tearing down environment")
		}
	})

	msRpcPubKey := e.SetupMercuryServer(env, dbSettings, serverResources,
		e.Config.MSAdminId, e.Config.MSAdminEncryptedKey)

	// Setup random mock server response for mercury price feed
	mockserverClient, err := ctfClient.ConnectMockServer(env)
	require.NoError(e.T, err, "Error connecting to mock server")

	msLocalUrl := env.URLs[mshelm.URLsKey][1]
	msClient := client.NewMercuryServerClient(msLocalUrl, e.Config.MSAdminId, e.Config.MSAdminKey)
	e.Config.MSLocalUrl = msLocalUrl
	e.MSClient = msClient

	if e.IsExistingTestEnv {
		// Load existing contracts
		contractDeployer, err := contracts.NewContractDeployer(evmClient)
		require.NoError(e.T, err)
		e.VerifierContract, err = contractDeployer.LoadVerifier(common.HexToAddress(e.Config.VerifierAddress))
		require.NoError(e.T, err)
		e.VerifierProxyContract, err = contractDeployer.LoadVerifierProxy(common.HexToAddress(e.Config.VerifierProxyAddress))
		require.NoError(e.T, err)
		e.ExchangerContract, err = contractDeployer.LoadExchanger(common.HexToAddress(e.Config.ExchangerAddress))
		require.NoError(e.T, err)
	} else {
		// Deploy new contracts and setup jobs
		var feedIdStr = "ETH-USD-1"
		var feedId = StringToByte32(feedIdStr)
		e.Config.FeedId = feedIdStr

		nodesWithoutBootstrap := chainlinkNodes[1:]
		ocrConfig, err := BuildMercuryOCRConfig(nodesWithoutBootstrap)
		require.NoError(e.T, err)

		verifier, verifierProxy, exchanger, _, _ := SetupMercuryContracts(
			evmClient, "", feedId, *ocrConfig)
		e.VerifierContract = verifier
		e.VerifierProxyContract = verifierProxy
		e.ExchangerContract = exchanger
		e.Config.VerifierAddress = verifier.Address()
		e.Config.VerifierProxyAddress = verifierProxy.Address()
		e.Config.ExchangerAddress = exchanger.Address()

		latestBlockNum, err := evmClient.LatestBlockNumber(context.Background())
		require.NoError(e.T, err)
		msRemoteUrl := env.URLs[mshelm.URLsKey][0]
		e.SetupMercuryNodeJobs(chainlinkNodes, mockserverClient, verifier.Address(),
			feedId, latestBlockNum, msRemoteUrl, msRpcPubKey, testNetwork.ChainID, 0)
		e.Config.MSRemoteUrl = msRemoteUrl

		verifier.SetConfig(feedId, *ocrConfig)

		e.WaitForDONReports()
	}
}

func (e *MercuryTestEnv) WaitForDONReports() {
	// Wait for the DON to start generating reports
	// TODO: use gomega Eventually to check reports in node logs or mercury server or mercury db
	d := 160 * time.Second
	log.Info().Msgf("Sleeping for %s to wait for Mercury env to be ready..", d)
	time.Sleep(d)
}

func (e *MercuryTestEnv) SetupDON(t *testing.T, evmNetwork blockchain.EVMNetwork, evmConfig environment.ConnectedChart) *environment.Environment {
	testEnv := environment.New(&environment.Config{
		TTL:             e.EnvTTL,
		NamespacePrefix: fmt.Sprintf("smoke-mercury-%s", strings.ReplaceAll(strings.ToLower(evmNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(map[string]interface{}{
			"app": map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "4000m",
						"memory": "4048Mi",
					},
					"limits": map[string]interface{}{
						"cpu":    "4000m",
						"memory": "4048Mi",
					},
				},
			},
		})).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "5",
			"toml": client.AddNetworksConfig(
				testconfig.BaseMercuryTomlConfig,
				evmNetwork),
			// "secretsToml": secretsToml,
			"prometheus": "true",
		}))
	err := testEnv.Run()
	require.NoError(t, err, "Error running test environment")

	return testEnv
}

func (e *MercuryTestEnv) SetupMercuryNodeJobs(
	chainlinkNodes []*client.Chainlink,
	mockserverClient *ctfClient.MockserverClient,
	contractID string,
	feedId [32]byte,
	fromBlock uint64,
	msRemoteUrl string,
	mercuryServerPubKey ed25519.PublicKey,
	chainID int64,
	keyIndex int,
) {
	err := mockserverClient.SetRandomValuePath("/variable")
	require.NoError(e.T, err, "Setting mockserver value path shouldn't fail")

	observationSource := fmt.Sprintf(`
// Benchmark Price
price1          [type=http method=GET url="%[1]s" allowunrestrictednetworkaccess="true"];
price1_parse    [type=jsonparse path="data,result"];
price1_multiply [type=multiply times=100000000 index=0];

price1 -> price1_parse -> price1_multiply;

// Bid
bid          [type=http method=GET url="%[1]s" allowunrestrictednetworkaccess="true"];
bid_parse    [type=jsonparse path="data,result"];
bid_multiply [type=multiply times=100000000 index=1];

bid -> bid_parse -> bid_multiply;

// Ask
ask          [type=http method=GET url="%[1]s" allowunrestrictednetworkaccess="true"];
ask_parse    [type=jsonparse path="data,result"];
ask_multiply [type=multiply times=100000000 index=2];

ask -> ask_parse -> ask_multiply;	

// Block Num + Hash
b1                 [type=ethgetblock];
bnum_lookup        [type=lookup key="number" index=3];
bhash_lookup       [type=lookup key="hash" index=4];

b1 -> bnum_lookup;
b1 -> bhash_lookup;`, mockserverClient.Config.ClusterURL+"/variable")

	bootstrapNode := chainlinkNodes[0]
	bootstrapNode.RemoteIP()
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	require.NoError(e.T, err, "Shouldn't fail reading P2P keys from bootstrap node")
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID

	bootstrapSpec := &client.OCR2TaskJobSpec{
		Name:    "ocr2 bootstrap node",
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: contractID,
			Relay:      "evm",
			RelayConfig: map[string]interface{}{
				"chainID": int(chainID),
				"feedID":  fmt.Sprintf("\"0x%x\"", feedId),
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
		},
	}
	_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
	require.NoError(e.T, err, "Shouldn't fail creating bootstrap job on bootstrap node")
	P2Pv2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PId, bootstrapNode.RemoteIP(), 6690)

	for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
		nodeOCRKeys, err := chainlinkNodes[nodeIndex].MustReadOCR2Keys()
		require.NoError(e.T, err, "Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
		csaKeys, _, err := chainlinkNodes[nodeIndex].ReadCSAKeys()
		require.NoError(e.T, err)
		// csaKeyId := csaKeys.Data[0].ID
		csaPubKey := csaKeys.Data[0].Attributes.PublicKey

		var nodeOCRKeyId []string
		for _, key := range nodeOCRKeys.Data {
			if key.Attributes.ChainType == string(chaintype.EVM) {
				nodeOCRKeyId = append(nodeOCRKeyId, key.ID)
				break
			}
		}

		jobSpec := client.OCR2TaskJobSpec{
			Name:            "ocr2",
			JobType:         "offchainreporting2",
			MaxTaskDuration: "1s",
			OCR2OracleSpec: job.OCR2OracleSpec{
				PluginType: "mercury",
				// PluginConfig: map[string]interface{}{
				// 	"juelsPerFeeCoinSource": `"""
				// 		bn1          [type=ethgetblock];
				// 		bn1_lookup   [type=lookup key="number"];
				// 		bn1 -> bn1_lookup;
				// 	"""`,
				// },
				PluginConfig: map[string]interface{}{
					// "serverHost":   fmt.Sprintf("\"%s:1338\"", mercury_server.URLsKey),
					"serverURL":    fmt.Sprintf("\"%s:1338\"", msRemoteUrl[7:len(msRemoteUrl)-5]),
					"serverPubKey": fmt.Sprintf("\"%s\"", hex.EncodeToString(mercuryServerPubKey)),
				},
				Relay: "evm",
				RelayConfig: map[string]interface{}{
					"chainID":   int(chainID),
					"feedID":    fmt.Sprintf("\"0x%x\"", feedId),
					"fromBlock": fromBlock,
				},
				// RelayConfigMercuryConfig: map[string]interface{}{
				// 	"clientPrivKeyID": csaKeyId,
				// },
				ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
				ContractID:                        contractID,
				OCRKeyBundleID:                    null.StringFrom(nodeOCRKeyId[keyIndex]),
				TransmitterID:                     null.StringFrom(csaPubKey),
				P2PV2Bootstrappers:                pq.StringArray{P2Pv2Bootstrapper},
			},
			ObservationSource: observationSource,
		}

		_, err = chainlinkNodes[nodeIndex].MustCreateJob(&jobSpec)
		require.NoError(e.T, err, "Shouldn't fail creating OCR Task job on OCR node %d", nodeIndex+1)
	}
	log.Info().Msg("Done creating OCR automation jobs")
}

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

// Build config with nodes for Mercury server
func BuildRpcNodesJsonConf(chainlinkNodes []*client.Chainlink) ([]byte, error) {
	var msRpcNodesConf []*oracle
	for i, chainlinkNode := range chainlinkNodes {
		nodeName := fmt.Sprint(i)
		nodeAddress, err := chainlinkNode.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		csaKeys, _, err := chainlinkNode.ReadCSAKeys()
		if err != nil {
			return nil, err
		}
		csaPubKey := csaKeys.Data[0].Attributes.PublicKey
		ocr2Keys, _, err := chainlinkNode.ReadOCR2Keys()
		if err != nil {
			return nil, err
		}
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
					PublicKey:   csaPubKey,
				},
			},
			Ocr2ConfigPublicKey:   []string{ocr2ConfigPublicKey},
			Ocr2OffchainPublicKey: []string{ocr2OffchainPublicKey},
			Ocr2OnchainPublicKey:  []string{ocr2OnchainPublicKey},
		}
		msRpcNodesConf = append(msRpcNodesConf, node)
	}
	return json.Marshal(msRpcNodesConf)
}

func buildInitialDbSql(adminId string, adminEncryptedKey string) (string, error) {
	data := struct {
		UserId       string
		UserRole     string
		EncryptedKey string
	}{
		UserId:       adminId,
		UserRole:     "admin",
		EncryptedKey: adminEncryptedKey,
	}

	// Get file path to the sql
	_, filename, _, _ := runtime.Caller(0)
	tmplPath := path.Join(path.Dir(filename), "/mercury_db_init_sql_template")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (e *MercuryTestEnv) SetupMercuryServer(
	testEnv *environment.Environment,
	dbSettings map[string]interface{},
	serverSettings map[string]interface{},
	adminId string,
	adminEncryptedKey string,
) ed25519.PublicKey {
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnv)
	require.NoError(e.T, err, "Error connecting to Chainlink nodes")

	rpcNodesJsonConf, _ := BuildRpcNodesJsonConf(chainlinkNodes)
	log.Info().Msgf("RPC nodes conf for mercury server: %s", rpcNodesJsonConf)

	// Generate keys for Mercury RPC server
	// rpcPrivKey, rpcPubKey, err := generateEd25519Keys()
	rpcPubKey, rpcPrivKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(e.T, err)

	initDbSql, err := buildInitialDbSql(adminId, adminEncryptedKey)
	require.NoError(e.T, err)
	log.Info().Msgf("Initialize mercury server db with:\n%s", initDbSql)

	settings := map[string]interface{}{
		"image": map[string]interface{}{
			"repository": os.Getenv("MERCURY_SERVER_IMAGE"),
			"tag":        os.Getenv("MERCURY_SERVER_TAG"),
		},
		"postgresql": map[string]interface{}{
			"enabled": true,
		},
		"qa": map[string]interface{}{
			"rpcPrivateKey": hex.EncodeToString(rpcPrivKey),
			"enabled":       true,
			"initDbSql":     initDbSql,
		},
		"rpcNodesConf": string(rpcNodesJsonConf),
		"prometheus":   "true",
	}

	if dbSettings != nil {
		settings["db"] = dbSettings
	}
	if serverSettings != nil {
		settings["resources"] = serverSettings
	}

	testEnv.AddHelm(mshelm.New(settings)).Run()

	return rpcPubKey
}

func BuildMercuryOCRConfig(chainlinkNodes []*client.Chainlink) (*contracts.MercuryOCRConfig, error) {
	// Build onchain config
	c := relaymercury.OnchainConfig{Min: big.NewInt(0), Max: big.NewInt(math.MaxInt64)}
	onchainConfig, err := (relaymercury.StandardOnchainConfigCodec{}).Encode(c)
	if err != nil {
		return nil, err
	}

	_, oracleIdentities := getOracleIdentities(chainlinkNodes)
	if err != nil {
		return nil, err
	}
	signerOnchainPublicKeys, _, f, onchainConfig,
		offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress time.Duration,
		20*time.Second,       // deltaResend time.Duration,
		100*time.Millisecond, // deltaRound time.Duration,
		0,                    // deltaGrace time.Duration,
		1*time.Minute,        // deltaStage time.Duration,
		100,                  // rMax uint8,
		[]int{len(chainlinkNodes)},
		oracleIdentities,
		[]byte{},
		0*time.Millisecond, // maxDurationQuery time.Duration,
		// TODO: fix mockserver and switch to 250ms?
		500*time.Millisecond, // maxDurationObservation time.Duration,
		500*time.Millisecond, // maxDurationReport time.Duration,
		500*time.Millisecond, // maxDurationShouldAcceptFinalizedReport time.Duration,
		500*time.Millisecond, // maxDurationShouldTransmitAcceptedReport time.Duration,
		// 250*time.Millisecond, // maxDurationObservation time.Duration,
		// 250*time.Millisecond, // maxDurationReport time.Duration,
		// 250*time.Millisecond, // maxDurationShouldAcceptFinalizedReport time.Duration,
		// 250*time.Millisecond, // maxDurationShouldTransmitAcceptedReport time.Duration,
		1, // f int,
		onchainConfig,
	)
	if err != nil {
		return nil, err
	}

	// Convert signers to addresses
	var signers []common.Address
	for _, signer := range signerOnchainPublicKeys {
		signers = append(signers, common.BytesToAddress(signer))
	}

	// Use node CSA pub key as transmitter
	transmitters := make([][32]byte, len(chainlinkNodes))
	for i, n := range chainlinkNodes {
		csaKeys, _, err := n.ReadCSAKeys()
		if err != nil {
			return nil, err
		}
		csaPubKey, err := hex.DecodeString(csaKeys.Data[0].Attributes.PublicKey)
		if err != nil {
			return nil, err
		}
		transmitters[i] = [32]byte(csaPubKey)
	}

	return &contracts.MercuryOCRConfig{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}, nil
}

func getOracleIdentities(chainlinkNodes []*client.Chainlink) ([]int, []confighelper.OracleIdentityExtra) {
	S := make([]int, len(chainlinkNodes))
	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(chainlinkNodes))
	sharedSecretEncryptionPublicKeys := make([]types.ConfigEncryptionPublicKey, len(chainlinkNodes))
	var wg sync.WaitGroup
	for i, cl := range chainlinkNodes {
		wg.Add(1)
		go func(i int, cl *client.Chainlink) error {
			defer wg.Done()

			ocr2Keys, err := cl.MustReadOCR2Keys()
			if err != nil {
				return err
			}
			var ocr2Config client.OCR2KeyAttributes
			for _, key := range ocr2Keys.Data {
				if key.Attributes.ChainType == string(chaintype.EVM) {
					ocr2Config = key.Attributes
					break
				}
			}

			keys, err := cl.MustReadP2PKeys()
			if err != nil {
				return err
			}
			p2pKeyID := keys.Data[0].Attributes.PeerID

			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_"))
			if err != nil {
				return err
			}

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			copy(offchainPkBytesFixed[:], offchainPkBytes)
			if err != nil {
				return err
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				return err
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			copy(configPkBytesFixed[:], configPkBytes)
			if err != nil {
				return err
			}

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_"))
			if err != nil {
				return err
			}

			csaKeys, _, err := cl.ReadCSAKeys()
			if err != nil {
				return err
			}

			sharedSecretEncryptionPublicKeys[i] = configPkBytesFixed
			oracleIdentities[i] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   types.Account(csaKeys.Data[0].Attributes.PublicKey),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[i] = 1

			return nil
		}(i, cl)
	}
	wg.Wait()
	log.Info().Msgf("Done fetching oracle identities")
	return S, oracleIdentities
}

func StringToByte32(str string) [32]byte {
	var bytes [32]byte
	copy(bytes[:], str)
	return bytes
}
