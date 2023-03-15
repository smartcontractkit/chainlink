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
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"
)

type TestEnv struct {
	Namespace             string
	NsPrefix              string
	Chart                 string
	Env                   *environment.Environment
	ChainlinkNodes        []*client.Chainlink
	MockserverClient      *ctfClient.MockserverClient
	FeedIds               []string              // feed id configured in Mercury
	MSClient              *client.MercuryServer // Mercury server client authenticated with admin role
	MSInfo                mercuryServerInfo
	IsExistingTestEnv     bool          // true if config in MERCURY_ENV_CONFIG_PATH contains namespace
	KeepEnv               bool          // Set via MERCURY_KEEP_ENV=true env
	EnvTTL                time.Duration // Set via MERCURY_ENV_TTL_MINS env
	ChainId               int64
	EvmNetwork            *blockchain.EVMNetwork
	EvmChart              *environment.ConnectedChart
	EvmClient             blockchain.EVMClient
	VerifierContract      contracts.Verifier
	VerifierProxyContract contracts.VerifierProxy
	ExchangerContract     contracts.Exchanger
	ContractInfo          mercuryContractInfo
}

type TestConfig struct {
	K8Namespace   string              `json:"k8Namespace"`
	ChainId       int64               `json:"chainId"`
	FeedIds       []string            `json:"feedIds"`
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
func configFromFile(path string) (*TestConfig, error) {
	c := &TestConfig{}
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

func (c *TestConfig) Json() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *TestConfig) Save() (string, error) {
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

func NewEnv(namespacePrefix string) (TestEnv, error) {
	testEnv := TestEnv{}

	c, err := configFromFile(os.Getenv("MERCURY_ENV_CONFIG_PATH"))
	if err != nil {
		// Fail when chain on env loaded from config is different than currently selected chain
		if c.ChainId != networks.SelectedNetwork.ChainID {
			return testEnv, fmt.Errorf("chain set in SELECTED_NETWORKS is" +
				" different than chain id set in config provided by MERCURY_ENV_CONFIG_PATH")
		}

		log.Info().Msgf("Using existing mercury environment based on config: %s\n%s",
			os.Getenv("MERCURY_ENV_CONFIG_PATH"), c.Json())

		testEnv.Namespace = c.K8Namespace
		testEnv.FeedIds = c.FeedIds
		testEnv.MSInfo = c.MSInfo
		testEnv.ContractInfo = c.ContractsInfo
		testEnv.IsExistingTestEnv = true
	} else {
		// Feed id can have max 32 characters
		testEnv.FeedIds = []string{"feed-1", "feed-2"}
		testEnv.MSInfo = mercuryServerInfo{
			AdminId:           os.Getenv("MS_DATABASE_FIRST_ADMIN_ID"),
			AdminKey:          os.Getenv("MS_DATABASE_FIRST_ADMIN_KEY"),
			AdminEncryptedKey: os.Getenv("MS_DATABASE_FIRST_ADMIN_ENCRYPTED_KEY"),
		}
		testEnv.IsExistingTestEnv = false
	}

	testEnv.KeepEnv = os.Getenv("MERCURY_KEEP_ENV") == "true"
	ttl, err := strconv.ParseUint(os.Getenv("MERCURY_ENV_TTL_MINS"), 10, 64)
	if err == nil {
		testEnv.EnvTTL = time.Duration(ttl) * time.Minute
	} else {
		// Set default TTL for k8 environment
		testEnv.EnvTTL = 20 * time.Minute
	}
	mschart := os.Getenv("MERCURY_CHART")
	if mschart == "" {
		return testEnv, errors.New("MERCURY_CHART must be provided, a local path or a name of a mercury-server helm chart")
	} else {
		testEnv.Chart = mschart
	}

	return testEnv, nil
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
	msResources map[string]interface{}) (TestEnv, error) {

	testEnv := TestEnv{}

	testEnv.KeepEnv = os.Getenv("MERCURY_KEEP_ENV") == "true"
	ttl, err := strconv.ParseUint(os.Getenv("MERCURY_ENV_TTL_MINS"), 10, 64)
	if err == nil {
		testEnv.EnvTTL = time.Duration(ttl) * time.Minute
	} else {
		// Set default TTL for k8 environment
		testEnv.EnvTTL = 20 * time.Minute
	}
	mschart := os.Getenv("MERCURY_CHART")
	if mschart == "" {
		return testEnv, errors.New("MERCURY_CHART must be provided, a local path or a name of a mercury-server helm chart")
	} else {
		testEnv.Chart = mschart
	}

	// Load mercury env info from a config file if it exists
	configPath := os.Getenv("MERCURY_ENV_CONFIG_PATH")
	if configPath != "" {
		c, err := configFromFile(configPath)
		if err != nil {
			return testEnv, err
		}
		// Fail when chain on env loaded from config is different than currently selected chain
		if c.ChainId != networks.SelectedNetwork.ChainID {
			return testEnv, fmt.Errorf("chain set in SELECTED_NETWORKS is" +
				" different than chain id set in config provided by MERCURY_ENV_CONFIG_PATH")
		}

		log.Info().Msgf("Using existing mercury env config from: %s\n%s",
			configPath, c.Json())

		testEnv.Namespace = c.K8Namespace
		testEnv.FeedIds = c.FeedIds
		testEnv.MSInfo = c.MSInfo
		testEnv.ContractInfo = c.ContractsInfo
		testEnv.IsExistingTestEnv = true
	} else {
		// Feed id can have max 32 characters
		testEnv.FeedIds = []string{"feed-1", "feed-2"}
		testEnv.MSInfo = mercuryServerInfo{
			AdminId:           os.Getenv("MS_DATABASE_FIRST_ADMIN_ID"),
			AdminKey:          os.Getenv("MS_DATABASE_FIRST_ADMIN_KEY"),
			AdminEncryptedKey: os.Getenv("MS_DATABASE_FIRST_ADMIN_ENCRYPTED_KEY"),
		}
		testEnv.IsExistingTestEnv = false
	}

	chainId := networks.SelectedNetwork.ChainID
	evmNetwork, evmConfig := setupEvmNetwork()

	env, chainlinkNodes, err := setupDON(testEnv.EnvTTL, testEnv.Namespace, namespacePrefix,
		testEnv.IsExistingTestEnv, evmNetwork, evmConfig)
	if err != nil {
		return testEnv, err
	}
	testEnv.ChainId = chainId
	testEnv.Env = env
	testEnv.ChainlinkNodes = chainlinkNodes

	evmClient, err := blockchain.NewEVMClient(evmNetwork, env)
	if err != nil {
		return testEnv, err
	}
	testEnv.EvmClient = evmClient

	msRpcPubKey, msLocalUrl, msRemoteUrl, msClient, err := setupMercuryServer(
		env, mschart, "", msDbSettings, msResources,
		testEnv.MSInfo.AdminId, testEnv.MSInfo.AdminKey, testEnv.MSInfo.AdminEncryptedKey)
	if err != nil {
		return testEnv, err
	}
	testEnv.MSInfo.LocalUrl = msLocalUrl
	testEnv.MSInfo.RemoteUrl = msRemoteUrl
	testEnv.MSClient = msClient

	// Setup random mock server response for mercury price feed
	mockserverClient, err := ctfClient.ConnectMockServer(env)
	if err != nil {
		return testEnv, err
	}
	err = mockserverClient.SetRandomValuePath("/variable")
	if err != nil {
		return testEnv, err
	}
	mockserverUrl := mockserverClient.Config.ClusterURL + "/variable"

	if testEnv.IsExistingTestEnv {
		verifierContract, verifierProxyContract, exchangerContract, err := LoadMercuryContracts(
			evmClient,
			testEnv.ContractInfo.VerifierAddress,
			testEnv.ContractInfo.VerifierProxyAddress,
			testEnv.ContractInfo.ExchangerAddress,
		)
		if err != nil {
			return testEnv, err
		}
		testEnv.VerifierContract = verifierContract
		testEnv.VerifierProxyContract = verifierProxyContract
		testEnv.ExchangerContract = exchangerContract
	} else {
		// Build OCR config
		nodesWithoutBootstrap := chainlinkNodes[1:]
		ocrConfig, err := buildMercuryOCRConfig(nodesWithoutBootstrap)
		if err != nil {
			return testEnv, err
		}

		// Deploy contracts
		verifierContract, verifierProxyContract, exchangerContract, _, err := DeployMercuryContracts(
			evmClient, "", *ocrConfig)
		if err != nil {
			return testEnv, err
		}
		testEnv.VerifierContract = verifierContract
		testEnv.VerifierProxyContract = verifierProxyContract
		testEnv.ExchangerContract = exchangerContract

		// Setup single verifier contract for multiple feeds
		err = InitVerifierContract(testEnv.FeedIds, *ocrConfig, verifierContract, verifierProxyContract)
		if err != nil {
			return testEnv, err
		}

		// Setup jobs on the nodes
		latestBlockNum, err := evmClient.LatestBlockNumber(context.Background())
		if err != nil {
			return testEnv, err
		}
		err = setupMercuryNodeJobs(chainlinkNodes, mockserverUrl, verifierContract.Address(),
			testEnv.FeedIds, latestBlockNum, msRemoteUrl, msRpcPubKey, evmNetwork.ChainID, 0)
		if err != nil {
			return testEnv, err
		}
	}

	// TODO: wait for reports for all feeds
	err = waitForReportsInMercuryDb(testEnv.FeedIds[0], evmClient, msClient)
	if err != nil {
		return testEnv, err
	}

	return testEnv, nil
}

// Build config of the current mercury env
func (e *TestEnv) Config() *TestConfig {
	return &TestConfig{
		K8Namespace: e.Env.Cfg.Namespace,
		ChainId:     e.ChainId,
		FeedIds:     e.FeedIds,
		ContractsInfo: mercuryContractInfo{
			VerifierAddress:      e.VerifierContract.Address(),
			VerifierProxyAddress: e.VerifierProxyContract.Address(),
			ExchangerAddress:     e.ExchangerContract.Address(),
		},
		MSInfo: e.MSInfo,
	}
}

func (e *TestEnv) Cleanup(t *testing.T) error {
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
	if !e.KeepEnv && e.Env != nil {
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

func (te *TestEnv) AddDON() error {
	if te.EvmNetwork == nil || te.EvmChart == nil {
		return fmt.Errorf("Setup evm network first")
	}

	env := environment.New(&environment.Config{
		TTL:              te.EnvTTL,
		NamespacePrefix:  fmt.Sprintf("%s-mercury-%s", te.NsPrefix, strings.ReplaceAll(strings.ToLower(te.EvmNetwork.Name), " ", "-")),
		Namespace:        te.Namespace,
		NoManifestUpdate: te.IsExistingTestEnv,
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
		AddHelm(*te.EvmChart).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "5",
			"toml": client.AddNetworksConfig(
				testconfig.BaseMercuryTomlConfig,
				*te.EvmNetwork),
			"prometheus": "true",
		}))
	err := env.Run()
	if err != nil {
		return err
	}
	te.Env = env

	nodes, err := client.ConnectChainlinkNodes(env)
	if err != nil {
		return err
	}
	te.ChainlinkNodes = nodes

	return nil
}

func buildBootstrapSpec(contractID string, chainID int64, fromBlock uint64, feedId string) *client.OCR2TaskJobSpec {
	uuid, _ := uuid.NewV4()
	return &client.OCR2TaskJobSpec{
		Name:    fmt.Sprintf("bootstrap-%s", uuid),
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: contractID,
			Relay:      "evm",
			RelayConfig: map[string]interface{}{
				"chainID":   int(chainID),
				"feedID":    fmt.Sprintf("\"0x%x\"", StringToByte32(feedId)),
				"fromBlock": fromBlock,
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
		},
	}
}

func buildOCRSpec(
	contractID string, chainID int64, fromBlock uint32,
	feedId string, mockserverUrl string,
	csaPubKey string, msRemoteUrl string, msPubKey ed25519.PublicKey,
	nodeOCRKey string, p2pV2Bootstrapper string) *client.OCR2TaskJobSpec {
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
	b1 -> bhash_lookup;`, mockserverUrl)

	uuid, _ := uuid.NewV4()
	return &client.OCR2TaskJobSpec{
		Name:            fmt.Sprintf("ocr2-%s", uuid),
		JobType:         "offchainreporting2",
		MaxTaskDuration: "1s",
		OCR2OracleSpec: job.OCR2OracleSpec{
			PluginType: "mercury",
			PluginConfig: map[string]interface{}{
				// "serverHost":   fmt.Sprintf("\"%s:1338\"", mercury_server.URLsKey),
				"serverURL":    fmt.Sprintf("\"%s:1338\"", msRemoteUrl[7:len(msRemoteUrl)-5]),
				"serverPubKey": fmt.Sprintf("\"%s\"", hex.EncodeToString(msPubKey)),
			},
			Relay: "evm",
			RelayConfig: map[string]interface{}{
				"chainID":   int(chainID),
				"feedID":    fmt.Sprintf("\"0x%x\"", StringToByte32(feedId)),
				"fromBlock": fromBlock,
			},
			// RelayConfigMercuryConfig: map[string]interface{}{
			// 	"clientPrivKeyID": csaKeyId,
			// },
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
			ContractID:                        contractID,
			OCRKeyBundleID:                    null.StringFrom(nodeOCRKey),
			TransmitterID:                     null.StringFrom(csaPubKey),
			P2PV2Bootstrappers:                pq.StringArray{p2pV2Bootstrapper},
		},
		ObservationSource: observationSource,
	}
}

// Setup node jobs for Mercury OCR
// For 'fromBlock', use the block number in which the config was set. Or latest block number if
// the config is not set yet
func setupMercuryNodeJobs(
	chainlinkNodes []*client.Chainlink,
	mockserverUrl string,
	contractID string,
	feedIds []string,
	fromBlock uint64,
	msRemoteUrl string,
	msPubKey ed25519.PublicKey,
	chainID int64,
	keyIndex int,
) error {

	bootstrapNode := chainlinkNodes[0]

	// Create bootstrap spec each feed
	for _, feedId := range feedIds {
		bootstrapSpec := buildBootstrapSpec(contractID, chainID, fromBlock, feedId)
		_, err := bootstrapNode.MustCreateJob(bootstrapSpec)
		if err != nil {
			return err
		}
	}

	bootstrapNode.RemoteIP()
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	if err != nil {
		return err
	}
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
	p2pV2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PId, bootstrapNode.RemoteIP(), 6690)

	// Create ocr jobs for each feed on each node
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

		for _, feedId := range feedIds {
			js := buildOCRSpec(
				contractID, chainID, uint32(fromBlock), feedId, mockserverUrl,
				csaPubKey, msRemoteUrl, msPubKey, nodeOCRKeyId[keyIndex], p2pV2Bootstrapper)
			_, err = chainlinkNodes[nodeIndex].MustCreateJob(js)
			if err != nil {
				return err
			}
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
	chartPath string,
	chartVersion string,
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

	if err = env.AddHelm(mshelm.New(chartPath, "", settings)).Run(); err != nil {
		return rpcPubKey, "", "", nil, nil
	}

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

func (te *TestEnv) SetupEvmNetwork() {
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
	te.EvmNetwork = &network
	te.EvmChart = &evmChart
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
	var b [32]byte
	copy(b[:], str)
	return b
}
