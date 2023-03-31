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
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	testconfig "github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"
)

type TestEnv struct {
	Id               string
	Namespace        string
	NsPrefix         string
	MSChartPath      string
	ResourcesConfig  *ResourcesConfig
	Env              *environment.Environment
	ChainlinkNodes   []*client.Chainlink
	MockserverClient *ctfClient.MockserverClient
	MSClient         *client.MercuryServer        // Mercury server client authenticated with admin role
	MSDbClient       *ctfClient.PostgresConnector // Mercury db client
	MSInfo           mercuryServerInfo
	IsExistingEnv    bool // true if config in MERCURY_ENV_CONFIG_PATH contains namespace
	SaveEnv          bool
	EnvTTL           time.Duration // Set via MERCURY_ENV_TTL_MINS env
	ChainId          int64
	EvmNetwork       *blockchain.EVMNetwork
	EvmChart         *environment.ConnectedChart
	EvmClient        blockchain.EVMClient
	ContractDeployer contracts.ContractDeployer
	Contracts        map[string]*contractInfo
	C                *TestEnvConfig // Config used to create this test env
	// Logs of action taken on the test env
	// When existing env is used, the logs are used to skip setting up
	// jobs, contracts, etc. as they are already in place
	ActionLog map[string]*envAction
}

type envAction struct {
	Done bool                   `json:"done"`
	Logs map[string]interface{} `json:"logs"`
}

type contractInfo struct {
	Address  string
	Contract interface{}
}

type mercuryServerInfo struct {
	RemoteUrl        string `json:"remoteUrl"`
	LocalUrl         string `json:"localUrl"`
	RemoteWsrpcUrl   string `json:"remoteWsrpcUrl"`
	LocalWsrpcUrl    string `json:"localWsrpcUrl"`
	RemoteDbUrl      string `json:"remoteDbUrl"`
	LocalDbUrl       string `json:"localDbUrl"`
	UserId           string `json:"userId"`
	UserKey          string `json:"userKey"`
	UserEncryptedKey string `json:"userEncryptedKey"`
	RpcPubKeyString  string `json:"rpcPubKey"`
	RpcPubKey        ed25519.PublicKey
	RpcNodesCsaKeys  []CsaKeyWrapper
}

// Fetch mercury environment config from local json file
func configFromFile(path string) (*TestEnvConfig, error) {
	c := &TestEnvConfig{}
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

func (c *TestEnvConfig) Json() string {
	b, _ := json.Marshal(c)
	return string(b)
}

var (
	EnvConfigPath = os.Getenv("MERCURY_ENV_CONFIG_PATH")
)

// New mercury env
//
// Required envs:
// MS_DATABASE_FIRST_ADMIN_ID: Mercury server admin id
// MS_DATABASE_FIRST_ADMIN_KEY: Mercury server admin key
// MS_DATABASE_FIRST_ADMIN_ENCRYPTED_KEY: Mercury server admin encrypted key
// Optional envs:
// MERCURY_ENV_CONFIG_PATH: Path to saved test env config
// MERCURY_ENV_SAVE: List of test env ids separated by comma that should be saved
// MERCURY_ENV_TTL_MINS: Env ttl in mins
func NewEnv(testEnvId string, namespacePrefix string, r *ResourcesConfig) (TestEnv, error) {
	te := TestEnv{}
	te.Id = testEnvId
	te.NsPrefix = namespacePrefix
	te.ResourcesConfig = r

	savedEnvs := strings.Split(os.Getenv("MERCURY_ENV_SAVE"), ",")
	te.SaveEnv = slices.Contains(savedEnvs, testEnvId)

	c, _ := configFromFile(EnvConfigPath)
	// Load env from config
	if c != nil && c.Id == testEnvId {
		te.C = c
		// Fail when chain on env loaded from config is different than currently selected chain
		if c.ChainId != networks.SelectedNetwork.ChainID {
			return te, fmt.Errorf("chain set in SELECTED_NETWORKS is" +
				" different than chain id set in config provided by MERCURY_ENV_CONFIG_PATH")
		}

		// Load keys for rpc nodes
		var csaKeys []CsaKeyWrapper
		for _, seedStr := range c.MSInfo.RpcNodesCsaPrivKeySeeds {
			b, err := hex.DecodeString(seedStr)
			if err != nil {
				return te, nil
			}
			privKey := ed25519.NewKeyFromSeed(b)
			csaKeys = append(csaKeys, CsaKeyWrapper{
				PrivateKeySeed: seedStr,
				KeyV2:          csakey.Raw(privKey).Key(),
			})
		}

		te.MSInfo = mercuryServerInfo{
			RemoteUrl:        c.MSInfo.RemoteUrl,
			LocalUrl:         c.MSInfo.LocalUrl,
			RemoteWsrpcUrl:   c.MSInfo.RemoteWsrpcUrl,
			LocalWsrpcUrl:    c.MSInfo.LocalWsrpcUrl,
			UserId:           c.MSInfo.UserId,
			UserKey:          c.MSInfo.UserKey,
			UserEncryptedKey: c.MSInfo.UserEncryptedKey,
			RpcPubKey:        c.MSInfo.RpcPubKey,
			RpcNodesCsaKeys:  csaKeys,
		}

		// Load contract addresses
		te.Contracts = map[string]*contractInfo{}
		for k, addr := range c.ContractsInfo {
			te.Contracts[k] = &contractInfo{
				Address: addr,
			}
		}
		te.ActionLog = c.Actions

		te.Namespace = c.K8Namespace

		if te.Namespace != "" {
			te.IsExistingEnv = true
			log.Info().Msgf("Using existing mercury environment in %s. Env config: %s\n%s",
				te.Namespace, EnvConfigPath, c.Json())
		} else {
			te.IsExistingEnv = false
		}

		log.Info().Msgf("Using existing mercury environment based on config: %s\n%s",
			EnvConfigPath, c.Json())
	} else {
		te.MSInfo = mercuryServerInfo{
			UserId:           os.Getenv("MS_DATABASE_FIRST_ADMIN_ID"),
			UserKey:          os.Getenv("MS_DATABASE_FIRST_ADMIN_KEY"),
			UserEncryptedKey: os.Getenv("MS_DATABASE_FIRST_ADMIN_ENCRYPTED_KEY"),
		}
		te.Contracts = map[string]*contractInfo{}
		te.ActionLog = map[string]*envAction{}
		te.IsExistingEnv = false

		log.Info().Msgf("Using a new mercury environment")
	}

	ttl, err := strconv.ParseUint(os.Getenv("MERCURY_ENV_TTL_MINS"), 10, 64)
	if err == nil {
		te.EnvTTL = time.Duration(ttl) * time.Minute
	} else {
		// Set default TTL for k8 environment
		te.EnvTTL = 20 * time.Minute
	}

	mschart := os.Getenv("MERCURY_CHART")
	if mschart == "" {
		return te, errors.New("MERCURY_CHART must be provided, a local path or a name of a mercury-server helm chart")
	} else {
		te.MSChartPath = mschart
	}

	te.ChainId = networks.SelectedNetwork.ChainID

	te.Env = environment.New(&environment.Config{
		TTL:              te.EnvTTL,
		NamespacePrefix:  fmt.Sprintf("%s-mercury", te.NsPrefix),
		Namespace:        te.Namespace,
		NoManifestUpdate: te.IsExistingEnv,
	})

	return te, nil
}

// Build config of the current mercury env
func (te *TestEnv) Config() *TestEnvConfig {
	contractsInfo := map[string]string{}
	for k, c := range te.Contracts {
		contractsInfo[k] = c.Address
	}

	var k8namespace string
	if te.Env != nil {
		k8namespace = te.Env.Cfg.Namespace
	}

	var csaPrivKeySeeds []string
	for _, key := range te.MSInfo.RpcNodesCsaKeys {
		log.Info().Msgf("seed %s", key.PrivateKeySeed)
		csaPrivKeySeeds = append(csaPrivKeySeeds, key.PrivateKeySeed)
	}

	msInfo := MSInfoConf{
		RemoteUrl:               te.MSInfo.RemoteUrl,
		LocalUrl:                te.MSInfo.LocalUrl,
		RemoteWsrpcUrl:          te.MSInfo.RemoteWsrpcUrl,
		LocalWsrpcUrl:           te.MSInfo.LocalWsrpcUrl,
		UserId:                  te.MSInfo.UserId,
		UserKey:                 te.MSInfo.UserKey,
		UserEncryptedKey:        te.MSInfo.UserEncryptedKey,
		RpcPubKey:               te.MSInfo.RpcPubKey,
		RpcNodesCsaPrivKeySeeds: csaPrivKeySeeds,
	}

	return &TestEnvConfig{
		Id:            te.Id,
		K8Namespace:   k8namespace,
		ChainId:       te.ChainId,
		ContractsInfo: contractsInfo,
		MSInfo:        msInfo,
		Actions:       te.ActionLog,
	}
}

// Clean up the env
func (te *TestEnv) Cleanup(t *testing.T) error {
	if !te.IsExistingEnv && te.SaveEnv {
		envConfPath, err := te.Config().Save()
		if err == nil {
			log.Info().Msgf("Keep mercury environment running."+
				" Chain: %d. Initial TTL: %s", te.ChainId, te.EnvTTL)
			log.Info().Msgf("To reuse this env in next test with chain %d, set"+
				" MERCURY_ENV_CONFIG_PATH to \"%s\"", te.ChainId, envConfPath)
		} else {
			log.Error().Msgf("Could not save mercury env config to file. Err: %v", err)
		}
	}
	if te.SaveEnv {
		log.Info().Msgf("Keep this mercury env because MERCURY_ENV_SAVE contains this env id: %s", te.Id)
	} else {
		log.Info().Msgf("Destroy this mercury env because MERCURY_ENV_SAVE does not contain this env id: %s", te.Id)
		if te.Env != nil {
			return actions.TeardownSuite(t, te.Env, utils.ProjectRoot,
				te.ChainlinkNodes, nil, zapcore.PanicLevel, te.EvmClient)
		}
	}
	return nil
}

// Wait for the DON to start generating reports and storing them in mercury server db
func (te *TestEnv) WaitForReportsInMercuryDb(feedIds [][32]byte) error {
	log.Info().Msgf("Wait for mercury server to have at least one report in the db for feeds %s..", feedIds)

	latestBlockNum, err := te.EvmClient.LatestBlockNumber(context.Background())
	if err != nil {
		return err
	}

	timeout := time.Minute * 8
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	to := time.NewTimer(timeout)
	defer to.Stop()
	for {
		select {
		case <-to.C:
			return fmt.Errorf(
				"did not found reports for one of the %s feeds in the mercury db. Tried for %s",
				feedIds, timeout,
			)
		case <-ticker.C:
			var notFound = false
			for _, feedId := range feedIds {
				report, _, _ := te.MSClient.GetReportsByFeedId(Byte32ToString(feedId), latestBlockNum, client.StringFeedId)
				if report == nil || report.ChainlinkBlob == "" {
					log.Debug().Msgf("Report not found for feedId: %s, blockNumber: %d", feedId, latestBlockNum)
					notFound = true
				}
			}
			// Stop if at least one report found for each feed
			if !notFound {
				return nil
			}
		}
	}
}

// Add DON to existing env
func (te *TestEnv) AddDON(mockserverResources map[string]interface{}) error {
	if te.EvmNetwork == nil {
		return fmt.Errorf("setup evm network first")
	}

	te.Env.
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(mockserverResources)).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "5",
			"toml": client.AddNetworksConfig(
				testconfig.BaseMercuryTomlConfig,
				*te.EvmNetwork),
			"chainlink":  te.ResourcesConfig.DONResources,
			"db":         te.ResourcesConfig.DONDBResources,
			"prometheus": "true",
		}))
	err := te.Env.Run()
	if err != nil {
		return err
	}

	mockserverClient, err := ctfClient.ConnectMockServer(te.Env)
	if err != nil {
		return err
	}
	te.MockserverClient = mockserverClient
	err = mockserverClient.SetRandomValuePath("/variable")
	if err != nil {
		return err
	}

	nodes, err := client.ConnectChainlinkNodes(te.Env)
	if err != nil {
		return err
	}
	te.ChainlinkNodes = nodes

	return nil
}

// Deploy or load verifier proxy contract
func (te *TestEnv) AddVerifierProxyContract(contractId string) (contracts.VerifierProxy, error) {
	if te.Contracts[contractId] != nil {
		addr := te.Contracts[contractId].Address
		if addr == "" {
			return nil, fmt.Errorf("no address in config for %s", contractId)
		}
		c, err := te.ContractDeployer.LoadVerifierProxy(common.HexToAddress(addr))
		if err != nil {
			return nil, err
		}
		return c, nil
	} else {
		// Use zero address for access controller disables access control
		c, err := te.ContractDeployer.DeployVerifierProxy("0x0")
		if err != nil {
			return nil, err
		}
		te.EvmClient.WaitForEvents()
		te.Contracts[contractId] = &contractInfo{
			Address:  c.Address(),
			Contract: c,
		}
		return c, err
	}
}

// Deploy or load verifier contract
func (te *TestEnv) AddVerifierContract(contractId string, verifierProxyAddr string) (contracts.Verifier, error) {
	if te.Contracts[contractId] != nil {
		addr := te.Contracts[contractId].Address
		if addr == "" {
			return nil, fmt.Errorf("no address in config for %s", contractId)
		}
		c, err := te.ContractDeployer.LoadVerifier(common.HexToAddress(addr))
		if err != nil {
			return nil, err
		}
		return c, nil
	} else {
		c, err := te.ContractDeployer.DeployVerifier(verifierProxyAddr)
		if err != nil {
			return nil, err
		}
		te.EvmClient.WaitForEvents()
		te.Contracts[contractId] = &contractInfo{
			Address:  c.Address(),
			Contract: c,
		}
		return c, err
	}
}

// Deploy or load exchanger contract
func (te *TestEnv) AddExchangerContract(contractId string, verifierProxyAddr string, lookupURL string, maxDelay uint8) (contracts.Exchanger, error) {
	if te.Contracts[contractId] != nil {
		addr := te.Contracts[contractId].Address
		if addr == "" {
			return nil, fmt.Errorf("no address in config for %s", contractId)
		}
		c, err := te.ContractDeployer.LoadExchanger(common.HexToAddress(addr))
		if err != nil {
			return nil, err
		}
		return c, nil
	} else {
		c, err := te.ContractDeployer.DeployExchanger(verifierProxyAddr, lookupURL, maxDelay)
		if err != nil {
			return nil, err
		}
		te.EvmClient.WaitForEvents()
		te.Contracts[contractId] = &contractInfo{
			Address:  c.Address(),
			Contract: c,
		}
		return c, err
	}
}

func (te *TestEnv) SetConfigAndInitializeVerifierContract(
	actionId string, verifierContractId string, verifierProxyContractId string,
	feedId [32]byte, ocrConfig contracts.MercuryOCRConfig) (uint64, error) {
	if te.IsExistingEnv {
		return uint64(te.ActionLog[actionId].Logs["blockNumber"].(float64)), nil
	} else {
		verifierContract := te.Contracts[verifierContractId].Contract.(contracts.Verifier)
		verifierProxyContract := te.Contracts[verifierProxyContractId].Contract.(contracts.VerifierProxy)

		err := verifierContract.SetConfig(feedId, ocrConfig)
		if err != nil {
			return 0, err
		}
		configDetails, err := verifierContract.LatestConfigDetails(feedId)
		if err != nil {
			return 0, err
		}
		if configDetails.ConfigCount != 1 {
			return 0, fmt.Errorf("verifier config count should be 1 but is %d", configDetails.ConfigCount)
		}
		log.Info().Msgf("Verifier.LatestConfigDetails for feedId: %s: %v\nConfig digest: %x", feedId, configDetails, configDetails.ConfigDigest)

		err = verifierProxyContract.InitializeVerifier(configDetails.ConfigDigest, verifierContract.Address())
		if err != nil {
			return 0, err
		}

		// Use latest block number from L2
		// Don't use block block number from config details as Arbitrum uses block numbers from L1
		latestBlockNum, err := te.EvmClient.LatestBlockNumber(context.Background())
		if err != nil {
			return 0, nil
		}
		log.Info().Msgf("Latest block number: %d", latestBlockNum)

		// TODO: Remove and use bn from config digest when https://smartcontract-it.atlassian.net/browse/MERC-248 resolved
		bnForJobs := uint64(latestBlockNum) - 200
		log.Info().Msgf("Block number used for job specs: %d", bnForJobs)

		te.ActionLog[actionId] = &envAction{
			Done: true,
			Logs: map[string]interface{}{
				"blockNumber": bnForJobs,
			},
		}

		return bnForJobs, nil
	}
}

func (te *TestEnv) errorIfActionNotDone(actionId string) error {
	a := te.ActionLog[actionId]
	if a == nil || !a.Done {
		return fmt.Errorf("action %s not done in the env config provided in %s",
			actionId, EnvConfigPath)
	}
	return nil
}

func (te *TestEnv) saveAction(actionId string, envAction *envAction) {
	te.ActionLog[actionId] = envAction
}

func buildBootstrapSpec(contractID string, chainID int64, fromBlock uint64, feedId [32]byte) *client.OCR2TaskJobSpec {
	uuid, _ := uuid.NewV4()
	return &client.OCR2TaskJobSpec{
		Name:    fmt.Sprintf("bootstrap-%s", uuid),
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: contractID,
			Relay:      "evm",
			FeedID:     common.BytesToHash(feedId[:]),
			RelayConfig: map[string]interface{}{
				"chainID": int(chainID),
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
		},
	}
}

func buildOCRSpec(
	contractID string, chainID int64, fromBlock uint32,
	feedId [32]byte, mockserverUrl string,
	csaPubKey string, msRemoteUrl string, msPubKey string,
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
		Name:              fmt.Sprintf("ocr2-%s", uuid),
		JobType:           "offchainreporting2",
		MaxTaskDuration:   "1s",
		ForwardingAllowed: false,
		OCR2OracleSpec: job.OCR2OracleSpec{
			PluginType: "mercury",
			PluginConfig: map[string]interface{}{
				// "serverHost":   fmt.Sprintf("\"%s:1338\"", mercury_server.URLsKey),
				"serverURL":    fmt.Sprintf("\"%s:1338\"", msRemoteUrl[7:len(msRemoteUrl)-5]),
				"serverPubKey": fmt.Sprintf("\"%s\"", msPubKey),
			},
			Relay: "evm",
			RelayConfig: map[string]interface{}{
				"chainID":   int(chainID),
				"fromBlock": fromBlock,
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
			ContractID:                        contractID,
			FeedID:                            common.BytesToHash(feedId[:]),
			OCRKeyBundleID:                    null.StringFrom(nodeOCRKey),
			TransmitterID:                     null.StringFrom(csaPubKey),
			P2PV2Bootstrappers:                pq.StringArray{p2pV2Bootstrapper},
		},
		ObservationSource: observationSource,
	}
}

func (te *TestEnv) GetBootstrapNode() *client.Chainlink {
	return te.ChainlinkNodes[0]
}

func (te *TestEnv) AddBootstrapJob(actionId, contractId string, fromBlock uint64, feedId [32]byte) error {
	if te.IsExistingEnv {
		return te.errorIfActionNotDone(actionId)
	} else {
		bootstrapSpec := buildBootstrapSpec(contractId, te.ChainId, fromBlock, feedId)
		_, err := te.GetBootstrapNode().MustCreateJob(bootstrapSpec)
		if err != nil {
			return err
		}

		te.saveAction(actionId, &envAction{Done: true})
		return nil
	}
}

// Setup node jobs for Mercury OCR
// For 'fromBlock', use the block number in which the config was set. Or latest block number if
// the config is not set yet
func (te *TestEnv) AddOCRJobs(actionId string, contractId string, fromBlock uint64, feedId [32]byte) error {
	if te.IsExistingEnv {
		return te.errorIfActionNotDone(actionId)
	}

	bootstrapP2PIds, err := te.GetBootstrapNode().MustReadP2PKeys()
	if err != nil {
		return err
	}
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
	p2pV2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PId, te.GetBootstrapNode().RemoteIP(), 6690)
	mockserverUrl := te.MockserverClient.Config.ClusterURL + "/variable"

	// Create ocr jobs for each feed on each node
	for nodeIndex := 1; nodeIndex < len(te.ChainlinkNodes); nodeIndex++ {
		nodeOCRKeys, err := te.ChainlinkNodes[nodeIndex].MustReadOCR2Keys()
		if err != nil {
			return err
		}
		csaKeys, _, err := te.ChainlinkNodes[nodeIndex].ReadCSAKeys()
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

		js := buildOCRSpec(
			contractId, te.ChainId, uint32(fromBlock), feedId, mockserverUrl,
			csaPubKey, te.MSInfo.RemoteUrl, te.MSInfo.RpcPubKeyString, nodeOCRKeyId[0], p2pV2Bootstrapper)
		_, err = te.ChainlinkNodes[nodeIndex].MustCreateJob(js)
		if err != nil {
			return err
		}
	}
	te.saveAction(actionId, &envAction{Done: true})
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

func buildRpcNodesJsonConfMock() ([]byte, error) {
	_, filename, _, _ := runtime.Caller(0)
	p := path.Join(path.Dir(filename), "/rpc_nodes_conf_mock.json")
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func buildMockedRpcNodesConf() ([]RpcNode, []CsaKeyWrapper) {
	_, privKey, _ := ed25519.GenerateKey(rand.Reader)
	csaKeys := []CsaKeyWrapper{
		{
			PrivateKeySeed: hex.EncodeToString(privKey.Seed()),
			KeyV2:          csakey.Raw(privKey).Key(),
		},
	}

	rpcNodeConf := []RpcNode{
		{
			Id:            "0",
			Status:        "active",
			NodeAddress:   []string{"0x9aF03D0296F21f59aB956e83f9d969F544a021Fa"},
			OracleAddress: "0x0000000000000000000000000000000000000000",
			CsaKeys: []CsaKeyInfo{
				{
					NodeName:    "0",
					NodeAddress: "0x9aF03D0296F21f59aB956e83f9d969F544a021Fa",
					PublicKey:   csaKeys[0].KeyV2.PublicKeyString(),
				},
			},
			Ocr2ConfigPublicKey:   []string{"fdff12ced64d6419b432f5096aa9b3de04531cf923b0142095f3e40014e81305"},
			Ocr2OffchainPublicKey: []string{"93400913aedd411ed6ec5d13c83ca7d666636a43dfd1195d62b3f4c0e1e6ce49"},
			Ocr2OnchainPublicKey:  []string{"01f2b0776f613604149579c8aebcf6ccf091b765"},
		},
	}

	return rpcNodeConf, csaKeys
}

// Build config with nodes for Mercury server
func buildRpcNodesConf(chainlinkNodes []*client.Chainlink) ([]*oracle, error) {
	var msRpcNodesConf []*oracle = []*oracle{}
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
	return msRpcNodesConf, nil
}

type User struct {
	Id        string `json:"id"`
	Key       string `json:"key"`
	Secret    string `json:"secret"`
	Role      string `json:"role"`
	Disabled  bool   `json:"disabled"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func buildInitialDbSql(users []User) (string, error) {
	// Get file path to the sql
	_, filename, _, _ := runtime.Caller(0)
	tmplPath := path.Join(path.Dir(filename), "/mercury_db_init_sql_template")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, users)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

type CsaKeyWrapper struct {
	csakey.KeyV2
	PrivateKeySeed string // hex encoded
}

type RpcNode struct {
	Id                    string       `json:"id"`
	Website               string       `json:"website"`
	Name                  string       `json:"name"`
	Status                string       `json:"status"`
	OracleAddress         string       `json:"oracleAddress"`
	NodeAddress           []string     `json:"nodeAddress"`
	Ocr2ConfigPublicKey   []string     `json:"ocr2ConfigPublicKey"`
	Ocr2OffchainPublicKey []string     `json:"ocr2OffchainPublicKey"`
	Ocr2OnchainPublicKey  []string     `json:"ocr2OnchainPublicKey"`
	CsaKeys               []CsaKeyInfo `json:"csaKeys"`
}

type CsaKeyInfo struct {
	NodeName    string `json:"nodeName"`
	NodeAddress string `json:"nodeAddress"`
	PublicKey   string `json:"publicKey"`
}

// Returns rpc pub key and list of node csa keys (when mock conf is used instead of DON)
func (te *TestEnv) AddMercuryServer(users *[]User) (ed25519.PublicKey, []CsaKeyWrapper, error) {
	var rpcPubKey ed25519.PublicKey
	var nodesCsaKeys []CsaKeyWrapper
	var chartSettings map[string]interface{}

	if te.IsExistingEnv {
		rpcPubKey = te.MSInfo.RpcPubKey
		nodesCsaKeys = te.MSInfo.RpcNodesCsaKeys

		// Provide dump values in chart settings for helm to not complain when reusing env
		chartSettings = map[string]interface{}{
			"image": map[string]interface{}{
				"repository": os.Getenv("MERCURY_SERVER_IMAGE"),
				"tag":        os.Getenv("MERCURY_SERVER_TAG"),
			},
			"resources": te.ResourcesConfig.MercuryResources,
			"postgresql": map[string]interface{}{
				"enabled": true,
				"primary": map[string]interface{}{
					"resources": te.ResourcesConfig.MercuryDBResources,
				},
				"mercury": map[string]interface{}{
					"initDbSql": "anything",
				},
			},
			"config": map[string]interface{}{
				"rpc": map[string]interface{}{
					"publicKey": "anything",
				},
			},
			"secrets": map[string]interface{}{
				"config": map[string]interface{}{
					"rpc": map[string]interface{}{
						"privateKey": "anything",
					},
					"database": map[string]interface{}{
						"url":           "postgres://postgres:testpass@mercury-server-postgresql:5432/testdb?sslmode=disable",
						"encryptionKey": "key",
					},
					"bootstrap": map[string]interface{}{
						"adminUsername": "admin2",
						"adminPassword": "admintestkey",
					},
				},
			},
			"chainlinkDONConfig": "anything",
			"ingress": map[string]interface{}{
				"private": map[string]interface{}{
					"enabled": false,
				},
				"public": map[string]interface{}{
					"enabled": false,
				},
			},
		}
	} else {
		// Build conf for rpc nodes
		var rpcNodesConf interface{}
		var err error
		if len(te.ChainlinkNodes) > 0 {
			rpcNodesConf, err = buildRpcNodesConf(te.ChainlinkNodes)
			if err != nil {
				return nil, nil, err
			}
		} else {
			rpcNodesConf, nodesCsaKeys = buildMockedRpcNodesConf()
			te.MSInfo.RpcNodesCsaKeys = nodesCsaKeys

			log.Info().Msg("Use rpc node json mock for mercury server as chainlink nodes not created")
		}
		chainlinkDONConfig, err := json.Marshal(rpcNodesConf)
		if err != nil {
			return nil, nil, err
		}
		log.Info().Msgf("RPC node json conf for mercury server: %s", chainlinkDONConfig)

		// Generate keys for Mercury RPC server
		var rpcPrivKey ed25519.PrivateKey
		rpcPubKey, rpcPrivKey, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		te.MSInfo.RpcPubKey = rpcPubKey
		te.MSInfo.RpcPubKeyString = hex.EncodeToString(rpcPubKey)

		var initDbSql string
		if users != nil {
			initDbSql, err = buildInitialDbSql(*users)
		} else {
			defaultUsers := []User{
				{
					Id:       te.MSInfo.UserId,
					Secret:   te.MSInfo.UserEncryptedKey,
					Role:     "admin",
					Disabled: false,
				},
			}
			initDbSql, err = buildInitialDbSql(defaultUsers)
		}
		if err != nil {
			return nil, nil, err
		}
		log.Info().Msgf("Initialize mercury server db with:\n%s", initDbSql)

		chartSettings = map[string]interface{}{
			"image": map[string]interface{}{
				"repository": os.Getenv("MERCURY_SERVER_IMAGE"),
				"tag":        os.Getenv("MERCURY_SERVER_TAG"),
			},
			"resources": te.ResourcesConfig.MercuryResources,
			"postgresql": map[string]interface{}{
				"enabled": true,
				"primary": map[string]interface{}{
					"resources": te.ResourcesConfig.MercuryDBResources,
				},
				"mercury": map[string]interface{}{
					"initDbSql": initDbSql,
				},
			},
			"config": map[string]interface{}{
				"rpc": map[string]interface{}{
					"publicKey": hex.EncodeToString(rpcPubKey),
				},
			},
			"secrets": map[string]interface{}{
				"config": map[string]interface{}{
					"rpc": map[string]interface{}{
						"privateKey": hex.EncodeToString(rpcPrivKey),
					},
					"database": map[string]interface{}{
						"url":           "postgres://postgres:testpass@mercury-server-postgresql:5432/testdb?sslmode=disable",
						"encryptionKey": "key",
					},
					"bootstrap": map[string]interface{}{
						"adminUsername": "admin2",
						"adminPassword": "admintestkey",
					},
				},
			},
			"chainlinkDONConfig": string(chainlinkDONConfig),
			"ingress": map[string]interface{}{
				"private": map[string]interface{}{
					"enabled": false,
				},
				"public": map[string]interface{}{
					"enabled": false,
				},
			},
		}
	}
	te.Env.
		AddHelm(mshelm.New(te.MSChartPath, "", chartSettings)).
		Run()

	te.MSInfo.RemoteUrl = te.Env.URLs[mshelm.URLsKey][0]
	te.MSInfo.LocalUrl = te.Env.URLs[mshelm.URLsKey][1]
	te.MSInfo.RemoteWsrpcUrl = te.Env.URLs[mshelm.URLsKey][2]
	te.MSInfo.LocalWsrpcUrl = te.Env.URLs[mshelm.URLsKey][3]
	te.MSInfo.RemoteDbUrl = te.Env.URLs[mshelm.URLsKey][4]
	te.MSInfo.LocalDbUrl = te.Env.URLs[mshelm.URLsKey][5]

	te.MSClient = client.NewMercuryServerClient(te.MSInfo.LocalUrl, te.MSInfo.UserId, te.MSInfo.UserKey)

	// Connect to mercury db
	spl := strings.Split(te.MSInfo.LocalDbUrl, ":")
	port := spl[len(spl)-1]
	db, err := ctfClient.NewPostgresConnector(&ctfClient.PostgresConfig{
		Host:     "localhost",
		Port:     port,
		User:     "postgres",
		Password: "testpass",
		DBName:   "testdb",
	})
	if err != nil {
		return rpcPubKey, nodesCsaKeys, err
	}
	te.MSDbClient = db
	log.Info().Msgf("Connected to Mercury db: localhost:%s", port)

	return rpcPubKey, nodesCsaKeys, nil
}

func (te *TestEnv) ClearMercuryReportsInDb() error {
	res, err := te.MSDbClient.Exec("TRUNCATE reports;")
	_ = res
	return err
}

func (te *TestEnv) GetAllReportsFromMercuryDb() ([]ReportTableRow, error) {
	var d []ReportTableRow
	err := te.MSDbClient.Select(&d, "SELECT * FROM reports;")
	return d, err
}

type ReportTableRow struct {
	Id                    int       `db:"id"`
	FeedId                []byte    `db:"feed_id"`
	Price                 float64   `db:"price"`
	Bid                   float64   `db:"bid"`
	Ask                   float64   `db:"ask"`
	OperatorName          string    `db:"operator_name"`
	FullReport            []byte    `db:"full_report"`
	Blob                  []byte    `db:"blob"`
	ValidFromBlockNumber  int64     `db:"valid_from_block_number"`
	ConfigDigest          []byte    `db:"config_digest"`
	Epoch                 int64     `db:"epoch"`
	Round                 int8      `db:"round"`
	ObservationsTimestamp int64     `db:"observations_timestamp"`
	TransmittingOperator  []byte    `db:"transmitting_operator"`
	CreatedAt             time.Time `db:"created_at"`
	CurrentBlockNumber    int64     `db:"current_block_number"`
	CurrentBlockHash      []byte    `db:"current_block_hash"`
}

func (te *TestEnv) BuildOCRConfig() (*contracts.MercuryOCRConfig, error) {
	// Build onchain config
	c := relaymercury.OnchainConfig{Min: big.NewInt(0), Max: big.NewInt(math.MaxInt64)}
	onchainConfig, err := (relaymercury.StandardOnchainConfigCodec{}).Encode(c)
	if err != nil {
		return nil, err
	}

	chainlinkNodes := te.ChainlinkNodes[1:]
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

func (te *TestEnv) AddEvmNetwork() error {
	network := networks.SelectedNetwork
	if network.Simulated {
		evmChart := eth.New(nil)
		te.Env.
			AddHelm(evmChart).
			Run()
	}
	te.EvmNetwork = &network

	evmClient, err := blockchain.NewEVMClient(network, te.Env)
	if err != nil {
		return err
	}
	te.EvmClient = evmClient

	contractDeployer, err := contracts.NewContractDeployer(evmClient)
	if err != nil {
		return err
	}
	te.ContractDeployer = contractDeployer

	return nil
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

func SetupMultiFeedSingleVerifierEnv(
	envId string,
	nsPrefix string,
	feedIDs [][32]byte,
	r *ResourcesConfig,
) (*TestEnv, contracts.VerifierProxy, error) {
	testEnv, err := NewEnv(envId, nsPrefix, r)
	if err != nil {
		return nil, nil, err
	}
	err = testEnv.AddEvmNetwork()
	if err != nil {
		return &testEnv, nil, err
	}
	if err = testEnv.AddDON(GetMockserverResources(len(feedIDs))); err != nil {
		return &testEnv, nil, err
	}
	ocrConfig, err := testEnv.BuildOCRConfig()
	if err != nil {
		return &testEnv, nil, err
	}
	_, _, err = testEnv.AddMercuryServer(nil)
	if err != nil {
		return &testEnv, nil, err
	}
	verifierProxyContract, err := testEnv.AddVerifierProxyContract("verifierProxy1")
	if err != nil {
		return &testEnv, nil, err
	}
	verifierContract, err := testEnv.AddVerifierContract("verifier1", verifierProxyContract.Address())
	if err != nil {
		return &testEnv, verifierProxyContract, err
	}
	// Use single verifier contract for all feeds
	for _, feedId := range feedIDs {
		blockNumber, err := testEnv.SetConfigAndInitializeVerifierContract(
			fmt.Sprintf("setAndInitialize%sVerifier", feedId),
			"verifier1",
			"verifierProxy1",
			feedId,
			*ocrConfig,
		)
		if err != nil {
			return &testEnv, verifierProxyContract, err
		}

		if err = testEnv.AddBootstrapJob(fmt.Sprintf("createBoostrapFor%s", feedId), verifierContract.Address(), uint64(blockNumber), feedId); err != nil {
			return &testEnv, verifierProxyContract, err
		}

		if err = testEnv.AddOCRJobs(fmt.Sprintf("createOcrJobsFor%s", feedId), verifierContract.Address(), uint64(blockNumber), feedId); err != nil {
			return &testEnv, verifierProxyContract, err
		}
	}
	if err = testEnv.WaitForReportsInMercuryDb(feedIDs); err != nil {
		return &testEnv, verifierProxyContract, err
	}
	return &testEnv, verifierProxyContract, nil
}

func StringToByte32(str string) [32]byte {
	var b [32]byte
	copy(b[:], str)
	return b
}

// [32]byte to string without trailing zeros (x00) if byte array not fully filled
func Byte32ToString(b [32]byte) string {
	n := bytes.IndexByte(b[:], 0)
	if n == -1 {
		return string(b[:])
	}
	return string(b[:n])
}
