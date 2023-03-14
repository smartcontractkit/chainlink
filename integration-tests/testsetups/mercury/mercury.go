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
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"
)

type MercuryTestEnv struct {
	Namespace             string
	NsPrefix              string
	Env                   *environment.Environment
	ChainlinkNodes        []*client.Chainlink
	MockserverClient      *ctfClient.MockserverClient
	FeedId                string                // feed id configured in Mercury
	MSClient              *client.MercuryServer // Mercury server client authenticated with admin role
	MSInfo                mercuryServerInfo
	IsExistingTestEnv     bool          // true if config in MERCURY_ENV_CONFIG_PATH contains namespace
	KeepEnv               bool          // Set via MERCURY_KEEP_ENV=true env
	EnvTTL                time.Duration // Set via MERCURY_ENV_TTL_MINS env
	ChainId               int64
	EvmClient             blockchain.EVMClient
	VerifierContract      contracts.Verifier
	VerifierProxyContract contracts.VerifierProxy
	ExchangerContract     contracts.Exchanger
	ContractInfo          mercuryContractInfo
}

type mercuryTestConfig struct {
	K8Namespace   string              `json:"k8Namespace"`
	ChainId       int64               `json:"chainId"`
	FeedId        string              `json:"feedId"`
	ContractsInfo mercuryContractInfo `json:"contracts"`
	MSInfo        mercuryServerInfo   `json:"mercuryServer"`
}

type mercuryContractInfo struct {
	VerifierAddress      string `json:"verifierAddress"`
	VerifierProxyAddress string `json:"verifierProxyAddress"`
	ExchangerAddress     string `json:"exchangerAddress"`
}

type mercuryServerInfo struct {
	RemoteUrl         string `json:"remoteUrl"`
	LocalUrl          string `json:"localUrl"`
	AdminId           string `json:"adminId"`
	AdminKey          string `json:"adminKey"`
	AdminEncryptedKey string `json:"adminEncryptedKey"`
}

// Fetch mercury environment config from local json file
func configFromFile(path string) (*mercuryTestConfig, error) {
	c := &mercuryTestConfig{}
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	b, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
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
	confPath := fmt.Sprintf("%s/%s.json", confDir, c.K8Namespace)
	f, _ := json.MarshalIndent(c, "", " ")
	err = ioutil.WriteFile(confPath, f, 0644)

	return confPath, err
}

// Setup new mercury env
// Required envs:
// MS_DATABASE_FIRST_ADMIN_ID: mercury server admin id
// MS_DATABASE_FIRST_ADMIN_KEY: mercury server admin key
// MS_DATABASE_FIRST_ADMIN_ENCRYPTED_KEY: mercury server admin encrypted key
// Optional envs:
// MERCURY_ENV_CONFIG_PATH: path to saved mercury test env config
// MERCURY_KEEP_ENV: Env config file will be generated and the env will not be destroyed when true
// MERCURY_ENV_TTL_MINS: Env ttl in min
func SetupMercuryTestEnv(
	namespacePrefix string,
	msDbSettings map[string]interface{},
	msResources map[string]interface{}) (*MercuryTestEnv, error) {

	var (
		feedId                    string
		namespace                 string
		isExistingTestEnv         bool
		keepEnv                   bool
		envTTL                    time.Duration
		chainId                   int64
		msLocalUrl                string
		msRemoteUrl               string
		msAdminId                 string
		msAdminKey                string
		msAdminEncryptedKey       string
		existingVerifierAddr      string
		existingVerifierProxyAddr string
		existingExchangerAddr     string
		verifierContract          contracts.Verifier
		verifierProxyContract     contracts.VerifierProxy
		exchangerContract         contracts.Exchanger
	)

	keepEnv = os.Getenv("MERCURY_KEEP_ENV") == "true"
	ttl, err := strconv.ParseUint(os.Getenv("MERCURY_ENV_TTL_MINS"), 10, 64)
	if err == nil {
		envTTL = time.Duration(ttl) * time.Minute
	} else {
		// Set default TTL for k8 environment
		envTTL = 20 * time.Minute
	}

	// Load mercury env info from a config file if it exists
	configPath := os.Getenv("MERCURY_ENV_CONFIG_PATH")
	if configPath != "" {
		c, err := configFromFile(configPath)
		if err != nil {
			return nil, err
		}
		// Fail when chain on env loaded from config is different than currently selected chain
		if c.ChainId != networks.SelectedNetwork.ChainID {
			return nil, fmt.Errorf("chain set in SELECTED_NETWORKS is" +
				" different than chain id set in config provided by MERCURY_ENV_CONFIG_PATH")
		}

		log.Info().Msgf("Using existing mercury env config from: %s\n%s",
			configPath, c.Json())

		namespace = c.K8Namespace
		feedId = c.FeedId
		msAdminId = c.MSInfo.AdminId
		msAdminKey = c.MSInfo.AdminKey
		msAdminEncryptedKey = c.MSInfo.AdminEncryptedKey
		existingVerifierAddr = c.ContractsInfo.VerifierAddress
		existingVerifierProxyAddr = c.ContractsInfo.VerifierProxyAddress
		existingExchangerAddr = c.ContractsInfo.ExchangerAddress

		isExistingTestEnv = true
	} else {
		// Feed id can have max 32 characters
		feedId = "feed-1234"
		msAdminId = os.Getenv("MS_DATABASE_FIRST_ADMIN_ID")
		msAdminKey = os.Getenv("MS_DATABASE_FIRST_ADMIN_KEY")
		msAdminEncryptedKey = os.Getenv("MS_DATABASE_FIRST_ADMIN_ENCRYPTED_KEY")

		isExistingTestEnv = false
	}

	chainId = networks.SelectedNetwork.ChainID
	evmNetwork, evmConfig := setupEvmNetwork()

	env, chainlinkNodes, err := setupDON(envTTL, namespace, namespacePrefix, isExistingTestEnv,
		evmNetwork, evmConfig)
	if err != nil {
		return nil, err
	}

	evmClient, err := blockchain.NewEVMClient(evmNetwork, env)
	if err != nil {
		return nil, err
	}

	msRpcPubKey, msLocalUrl, msRemoteUrl, msClient, err := setupMercuryServer(
		env, msDbSettings, msResources,
		msAdminId, msAdminKey, msAdminEncryptedKey)
	if err != nil {
		return nil, err
	}

	// Setup random mock server response for mercury price feed
	mockserverClient, err := ctfClient.ConnectMockServer(env)
	if err != nil {
		return nil, err
	}

	if isExistingTestEnv {
		verifierContract, verifierProxyContract, exchangerContract, err = LoadMercuryContracts(
			evmClient,
			existingVerifierAddr,
			existingVerifierProxyAddr,
			existingExchangerAddr,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// Build OCR config
		nodesWithoutBootstrap := chainlinkNodes[1:]
		ocrConfig, err := buildMercuryOCRConfig(nodesWithoutBootstrap)
		if err != nil {
			return nil, err
		}

		// Deploy contracts
		feedId := StringToByte32(feedId)
		verifierContract, verifierProxyContract, exchangerContract, _, err = DeployMercuryContracts(
			evmClient, "", *ocrConfig)
		if err != nil {
			return nil, err
		}

		// Setup feed verifier contract
		verifierContract.SetConfig(feedId, *ocrConfig)
		c, err := verifierContract.LatestConfigDetails(feedId)
		if err != nil {
			return nil, err
		}
		log.Info().Msgf("Latest Verifier config digest: %x", c.ConfigDigest)
		verifierProxyContract.InitializeVerifier(c.ConfigDigest, verifierContract.Address())

		// Setup jobs on the nodes
		if err != nil {
			return nil, err
		}
		setupMercuryNodeJobs(chainlinkNodes, mockserverClient, verifierContract.Address(),
			feedId, c.BlockNumber, msRemoteUrl, msRpcPubKey, evmNetwork.ChainID, 0)
	}

	err = waitForReportsInMercuryDb(feedId, evmClient, msClient)
	if err != nil {
		return nil, err
	}

	return &MercuryTestEnv{
		Namespace:         namespace,
		NsPrefix:          namespacePrefix,
		Env:               env,
		EnvTTL:            envTTL,
		KeepEnv:           keepEnv,
		ChainId:           chainId,
		IsExistingTestEnv: isExistingTestEnv,
		FeedId:            feedId,
		MSInfo: mercuryServerInfo{
			LocalUrl:          msLocalUrl,
			RemoteUrl:         msRemoteUrl,
			AdminId:           msAdminId,
			AdminKey:          msAdminKey,
			AdminEncryptedKey: msAdminEncryptedKey,
		},
		MSClient:              msClient,
		EvmClient:             evmClient,
		MockserverClient:      mockserverClient,
		ChainlinkNodes:        chainlinkNodes,
		VerifierContract:      verifierContract,
		VerifierProxyContract: verifierProxyContract,
		ExchangerContract:     exchangerContract,
	}, nil
}

// Build config of the current mercury env
func (e *MercuryTestEnv) Config() *mercuryTestConfig {
	return &mercuryTestConfig{
		K8Namespace: e.Env.Cfg.Namespace,
		ChainId:     e.ChainId,
		FeedId:      e.FeedId,
		ContractsInfo: mercuryContractInfo{
			VerifierAddress:      e.VerifierContract.Address(),
			VerifierProxyAddress: e.VerifierProxyContract.Address(),
			ExchangerAddress:     e.ExchangerContract.Address(),
		},
		MSInfo: e.MSInfo,
	}
}

func (e *MercuryTestEnv) Cleanup(t *testing.T) error {
	if !e.IsExistingTestEnv && e.KeepEnv {
		envConfPath, err := e.Config().Save()
		if err == nil {
			log.Info().Msgf("Keep mercury environment running."+
				" Chain: %d. Initial TTL: %s", e.ChainId, e.EnvTTL)
			log.Info().Msgf("To reuse this env in next test on chain %d, set:\n"+
				"\"MERCURY_ENV_CONFIG_PATH\"=\"%s\"", e.ChainId, envConfPath)
		} else {
			log.Error().Msgf("Could not save mercury env config to file. Err: %v", err)
		}
	}
	if !e.KeepEnv {
		log.Info().Msgf("Destroy this mercury env because MERCURY_KEEP_ENV not set to \"true\"")
		err := actions.TeardownSuite(t, e.Env, utils.ProjectRoot,
			e.ChainlinkNodes, nil, zapcore.PanicLevel, e.EvmClient)
		return err
	}
	return nil
}

// Wait for the DON to start generating reports and storing them in mercury server db
func waitForReportsInMercuryDb(
	feedId string, evmClient blockchain.EVMClient, msClient *client.MercuryServer) error {
	log.Info().Msg("Wait for mercury server to have at least one report in the db..")

	latestBlockNum, err := evmClient.LatestBlockNumber(context.Background())
	if err != nil {
		return err
	}

	timeout := time.Minute * 3
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	to := time.NewTimer(timeout)
	defer to.Stop()
	for {
		select {
		case <-to.C:
			return fmt.Errorf("no reports found in mercury db after %s", timeout)
		case <-ticker.C:
			report, _, _ := msClient.GetReports(feedId, latestBlockNum)
			if report != nil && report.ChainlinkBlob != "" {
				return nil
			}
		}
	}
}

func setupDON(envTTL time.Duration, namespace string, namespacePrefix string, isExistingTestEnv bool,
	evmNetwork blockchain.EVMNetwork, evmConfig environment.ConnectedChart) (*environment.Environment, []*client.Chainlink, error) {
	env := environment.New(&environment.Config{
		TTL:              envTTL,
		NamespacePrefix:  fmt.Sprintf("%s-mercury-%s", namespacePrefix, strings.ReplaceAll(strings.ToLower(evmNetwork.Name), " ", "-")),
		Namespace:        namespace,
		NoManifestUpdate: isExistingTestEnv,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(map[string]interface{}{
			"app": map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "8000m",
						"memory": "8048Mi",
					},
					"limits": map[string]interface{}{
						"cpu":    "8000m",
						"memory": "8048Mi",
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
			"prometheus": "true",
		}))
	err := env.Run()
	if err != nil {
		return nil, nil, err
	}

	nodes, err := client.ConnectChainlinkNodes(env)
	if err != nil {
		return env, nil, err
	}

	return env, nodes, nil
}

// Setup node jobs for Mercury OCR
// For 'fromBlock', use the block number in which the config was set. Or latest block number if
// the config is not set yet
func setupMercuryNodeJobs(
	chainlinkNodes []*client.Chainlink,
	mockserverClient *ctfClient.MockserverClient,
	contractID string,
	feedId [32]byte,
	fromBlock uint32,
	msRemoteUrl string,
	mercuryServerPubKey ed25519.PublicKey,
	chainID int64,
	keyIndex int,
) error {
	err := mockserverClient.SetRandomValuePath("/variable")
	if err != nil {
		return err
	}

	observationSource := fmt.Sprintf(`
// Benchmark Price
price1          [type=http method=GET url="%[1]s" allowunrestrictednetworkaccess="true"];
price1_parse    [type=jsonparse path="data,result"];
price1_multiply [type=multiply times=10 index=0];

price1 -> price1_parse -> price1_multiply;

// Bid
bid          [type=http method=GET url="%[1]s" allowunrestrictednetworkaccess="true"];
bid_parse    [type=jsonparse path="data,result"];
bid_multiply [type=multiply times=10 index=1];

bid -> bid_parse -> bid_multiply;

// Ask
ask          [type=http method=GET url="%[1]s" allowunrestrictednetworkaccess="true"];
ask_parse    [type=jsonparse path="data,result"];
ask_multiply [type=multiply times=10 index=2];

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
	if err != nil {
		return err
	}
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID

	bootstrapSpec := &client.OCR2TaskJobSpec{
		Name:    "ocr2 bootstrap node",
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: contractID,
			Relay:      "evm",
			RelayConfig: map[string]interface{}{
				"chainID":   int(chainID),
				"feedID":    fmt.Sprintf("\"0x%x\"", feedId),
				"fromBlock": fromBlock,
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
		},
	}
	_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
	if err != nil {
		return err
	}
	P2Pv2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PId, bootstrapNode.RemoteIP(), 6690)

	for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
		nodeOCRKeys, err := chainlinkNodes[nodeIndex].MustReadOCR2Keys()
		if err != nil {
			return err
		}
		csaKeys, _, err := chainlinkNodes[nodeIndex].ReadCSAKeys()
		if err != nil {
			return err
		}
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
		if err != nil {
			return err
		}
	}
	log.Info().Msg("Done creating OCR automation jobs")
	return nil
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
func buildRpcNodesJsonConf(chainlinkNodes []*client.Chainlink) ([]byte, error) {
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

func setupMercuryServer(
	env *environment.Environment,
	dbSettings map[string]interface{},
	serverSettings map[string]interface{},
	adminId string,
	adminKey string,
	adminEncryptedKey string,
) (ed25519.PublicKey, string, string, *client.MercuryServer, error) {
	chainlinkNodes, err := client.ConnectChainlinkNodes(env)
	if err != nil {
		return nil, "", "", nil, err
	}

	rpcNodesJsonConf, _ := buildRpcNodesJsonConf(chainlinkNodes)
	log.Info().Msgf("RPC nodes conf for mercury server: %s", rpcNodesJsonConf)

	// Generate keys for Mercury RPC server
	// rpcPrivKey, rpcPubKey, err := generateEd25519Keys()
	rpcPubKey, rpcPrivKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, "", "", nil, err
	}

	initDbSql, err := buildInitialDbSql(adminId, adminEncryptedKey)
	if err != nil {
		return nil, "", "", nil, err
	}
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

	env.AddHelm(mshelm.New(settings)).Run()

	msRemoteUrl := env.URLs[mshelm.URLsKey][0]
	msLocalUrl := env.URLs[mshelm.URLsKey][1]
	msClient := client.NewMercuryServerClient(msLocalUrl, adminId, adminKey)

	return rpcPubKey, msLocalUrl, msRemoteUrl, msClient, nil
}

func buildMercuryOCRConfig(chainlinkNodes []*client.Chainlink) (*contracts.MercuryOCRConfig, error) {
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
		// 500*time.Millisecond, // maxDurationObservation time.Duration,
		// 500*time.Millisecond, // maxDurationReport time.Duration,
		// 500*time.Millisecond, // maxDurationShouldAcceptFinalizedReport time.Duration,
		// 500*time.Millisecond, // maxDurationShouldTransmitAcceptedReport time.Duration,
		250*time.Millisecond, // maxDurationObservation time.Duration,
		250*time.Millisecond, // maxDurationReport time.Duration,
		250*time.Millisecond, // maxDurationShouldAcceptFinalizedReport time.Duration,
		250*time.Millisecond, // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,                    // f int,
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

func setupEvmNetwork() (blockchain.EVMNetwork, environment.ConnectedChart) {
	network := networks.SelectedNetwork
	var evmChart environment.ConnectedChart
	if network.Simulated {
		evmChart = eth.New(nil)
	} else {
		evmChart = eth.New(&eth.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
	}

	return network, evmChart
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
