package test_env

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/seth"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/runid"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	d "github.com/smartcontractkit/chainlink/integration-tests/docker"
)

var (
	ErrFundCLNode = "failed to fund CL node"
)

type CLClusterTestEnv struct {
	Cfg           *TestEnvConfig
	DockerNetwork *tc.DockerNetwork
	LogStream     *logstream.LogStream
	TestConfig    ctf_config.GlobalTestConfig

	/* components */
	ClCluster              *ClCluster
	MockAdapter            *test_env.Killgrave
	evmClients             map[int64]blockchain.EVMClient
	sethClients            map[int64]*seth.Client
	ContractDeployer       contracts.ContractDeployer
	ContractLoader         contracts.ContractLoader
	PrivateEthereumConfigs []*ctf_config.EthereumNetworkConfig
	EVMNetworks            []*blockchain.EVMNetwork
	rpcProviders           map[int64]*test_env.RpcProvider
	l                      zerolog.Logger
	t                      *testing.T
	isSimulatedNetwork     bool
}

func NewTestEnv() (*CLClusterTestEnv, error) {
	log.Logger = logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL")
	network, err := docker.CreateNetwork(log.Logger)
	if err != nil {
		return nil, err
	}
	return &CLClusterTestEnv{
		DockerNetwork: network,
		l:             log.Logger,
	}, nil
}

// WithTestEnvConfig sets the test environment cfg.
// Sets up private ethereum chain and MockAdapter containers with the provided cfg.
func (te *CLClusterTestEnv) WithTestEnvConfig(cfg *TestEnvConfig) *CLClusterTestEnv {
	te.Cfg = cfg
	if cfg.MockAdapter.ContainerName != "" {
		n := []string{te.DockerNetwork.Name}
		te.MockAdapter = test_env.NewKillgrave(n, te.Cfg.MockAdapter.ImpostersPath, test_env.WithContainerName(te.Cfg.MockAdapter.ContainerName), test_env.WithLogStream(te.LogStream))
	}
	return te
}

func (te *CLClusterTestEnv) WithTestInstance(t *testing.T) *CLClusterTestEnv {
	te.t = t
	te.l = logging.GetTestLogger(t)
	if te.MockAdapter != nil {
		te.MockAdapter.WithTestInstance(t)
	}
	return te
}

func (te *CLClusterTestEnv) ParallelTransactions(enabled bool) {
	for _, evmClient := range te.evmClients {
		evmClient.ParallelTransactions(enabled)
	}
}

func (te *CLClusterTestEnv) StartEthereumNetwork(cfg *ctf_config.EthereumNetworkConfig) (blockchain.EVMNetwork, test_env.RpcProvider, error) {
	// if environment is being restored from a previous state, use the existing config
	// this might fail terribly if temporary folders with chain data on the host machine were removed
	if te.Cfg != nil && te.Cfg.EthereumNetworkConfig != nil {
		cfg = te.Cfg.EthereumNetworkConfig
	}

	te.l.Info().
		Str("Execution Layer", string(*cfg.ExecutionLayer)).
		Str("Ethereum Version", string(*cfg.EthereumVersion)).
		Str("Custom Docker Images", fmt.Sprintf("%v", cfg.CustomDockerImages)).
		Msg("Starting Ethereum network")

	builder := test_env.NewEthereumNetworkBuilder()
	c, err := builder.WithExistingConfig(*cfg).
		WithTest(te.t).
		WithLogStream(te.LogStream).
		Build()
	if err != nil {
		return blockchain.EVMNetwork{}, test_env.RpcProvider{}, err
	}

	n, rpc, err := c.Start()

	if err != nil {
		return blockchain.EVMNetwork{}, test_env.RpcProvider{}, err
	}

	return n, rpc, nil
}

func (te *CLClusterTestEnv) StartMockAdapter() error {
	return te.MockAdapter.StartContainer()
}

// pass config here
func (te *CLClusterTestEnv) StartClCluster(nodeConfig *chainlink.Config, count int, secretsConfig string, testconfig ctf_config.GlobalTestConfig, opts ...ClNodeOption) error {
	if te.Cfg != nil && te.Cfg.ClCluster != nil {
		te.ClCluster = te.Cfg.ClCluster
	} else {
		// prepend the postgres version option from the toml config
		if testconfig.GetChainlinkImageConfig().PostgresVersion != nil && *testconfig.GetChainlinkImageConfig().PostgresVersion != "" {
			opts = append([]func(c *ClNode){
				func(c *ClNode) {
					c.PostgresDb.EnvComponent.ContainerVersion = *testconfig.GetChainlinkImageConfig().PostgresVersion
				},
			}, opts...)
		}
		opts = append(opts, WithSecrets(secretsConfig), WithLogStream(te.LogStream))
		te.ClCluster = &ClCluster{}
		for i := 0; i < count; i++ {
			ocrNode, err := NewClNode([]string{te.DockerNetwork.Name}, *testconfig.GetChainlinkImageConfig().Image, *testconfig.GetChainlinkImageConfig().Version, nodeConfig, opts...)
			if err != nil {
				return err
			}
			te.ClCluster.Nodes = append(te.ClCluster.Nodes, ocrNode)
		}
	}

	// Set test logger
	if te.t != nil {
		for _, n := range te.ClCluster.Nodes {
			n.SetTestLogger(te.t)
		}
	}

	// Start/attach node containers
	return te.ClCluster.Start()
}

// FundChainlinkNodes will fund all the provided Chainlink nodes with a set amount of native currency
func (te *CLClusterTestEnv) FundChainlinkNodes(amount *big.Float) error {
	if len(te.sethClients) == 0 && len(te.evmClients) == 0 {
		return fmt.Errorf("both EVMClients and SethClient are nil, unable to fund chainlink nodes")
	}

	if len(te.sethClients) > 0 && len(te.evmClients) > 0 {
		return fmt.Errorf("both EVMClients and SethClient are set, you can't use both at the same time")
	}

	if len(te.sethClients) > 0 {
		for _, sethClient := range te.sethClients {
			if err := actions_seth.FundChainlinkNodesFromRootAddress(te.l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(te.ClCluster.NodeAPIs()), amount); err != nil {
				return err
			}
		}
	}

	if len(te.evmClients) > 0 {
		for _, evmClient := range te.evmClients {
			for _, cl := range te.ClCluster.Nodes {
				if err := cl.Fund(evmClient, amount); err != nil {
					return fmt.Errorf("%s, err: %w", ErrFundCLNode, err)
				}
			}
			err := evmClient.WaitForEvents()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (te *CLClusterTestEnv) Terminate() error {
	// TESTCONTAINERS_RYUK_DISABLED=false by default so ryuk will remove all
	// the containers and the Network
	return nil
}

type CleanupOpts struct {
	TestName string
}

// Cleanup cleans the environment up after it's done being used, mainly for returning funds when on live networks and logs.
func (te *CLClusterTestEnv) Cleanup(opts CleanupOpts) error {
	te.l.Info().Msg("Cleaning up test environment")

	runIdErr := runid.RemoveLocalRunId(te.TestConfig.GetLoggingConfig().RunId)
	if runIdErr != nil {
		te.l.Warn().Msgf("Failed to remove .run.id file due to: %s (not a big deal, you can still remove it manually)", runIdErr.Error())
	}

	if te.t == nil {
		return fmt.Errorf("cannot cleanup test environment without a testing.T")
	}

	if te.ClCluster == nil || len(te.ClCluster.Nodes) == 0 {
		return fmt.Errorf("chainlink nodes are nil, unable cleanup chainlink nodes")
	}

	te.logWhetherAllContainersAreRunning()

	if len(te.evmClients) == 0 && len(te.sethClients) == 0 {
		return fmt.Errorf("both EVMClients and SethClient are nil, unable to return funds from chainlink nodes during cleanup")
	} else if te.isSimulatedNetwork {
		te.l.Info().
			Msg("Network is a simulated network. Skipping fund return.")
	} else {
		if err := te.returnFunds(); err != nil {
			return err
		}
	}

	err := te.handleNodeCoverageReports(opts.TestName)
	if err != nil {
		te.l.Error().Err(err).Msg("Error handling node coverage reports")
	}

	// close EVMClient connections
	for _, evmClient := range te.evmClients {
		err := evmClient.Close()
		return err
	}

	for _, sethClient := range te.sethClients {
		sethClient.Client.Close()
	}

	return nil
}

// handleNodeCoverageReports handles the coverage reports for the chainlink nodes
func (te *CLClusterTestEnv) handleNodeCoverageReports(testName string) error {
	testName = strings.ReplaceAll(testName, "/", "_")
	showHTMLCoverageReport := te.TestConfig.GetLoggingConfig().ShowHTMLCoverageReport != nil && *te.TestConfig.GetLoggingConfig().ShowHTMLCoverageReport
	isCI := os.Getenv("CI") != ""

	te.l.Info().
		Bool("showCoverageReportFlag", showHTMLCoverageReport).
		Bool("isCI", isCI).
		Bool("show", showHTMLCoverageReport || isCI).
		Msg("Checking if coverage report should be shown")

	var covHelper *d.NodeCoverageHelper

	if showHTMLCoverageReport || isCI {
		// Stop all nodes in the chainlink cluster.
		// This is needed to get go coverage profile from the node containers https://go.dev/doc/build-cover#FAQ
		// TODO: fix this as it results in: ERR LOG AFTER TEST ENDED ... INF 🐳 Stopping container
		err := te.ClCluster.Stop()
		if err != nil {
			return err
		}

		clDir, err := getChainlinkDir()
		if err != nil {
			return err
		}

		var coverageRootDir string
		if os.Getenv("GO_COVERAGE_DEST_DIR") != "" {
			coverageRootDir = filepath.Join(os.Getenv("GO_COVERAGE_DEST_DIR"), testName)
		} else {
			coverageRootDir = filepath.Join(clDir, ".covdata", testName)
		}

		var containers []tc.Container
		for _, node := range te.ClCluster.Nodes {
			containers = append(containers, node.Container)
		}

		covHelper, err = d.NewNodeCoverageHelper(context.Background(), containers, clDir, coverageRootDir)
		if err != nil {
			return err
		}
	}

	// Show html coverage report when flag is set (local runs)
	if showHTMLCoverageReport {
		path, err := covHelper.SaveMergedHTMLReport()
		if err != nil {
			return err
		}
		te.l.Info().Str("testName", testName).Str("filePath", path).Msg("Chainlink node coverage html report saved")
	}

	// Save percentage coverage report when running in CI
	if isCI {
		// Save coverage percentage to a file to show in the CI
		path, err := covHelper.SaveMergedCoveragePercentage()
		if err != nil {
			te.l.Error().Err(err).Str("testName", testName).Msg("Failed to save coverage percentage for test")
		} else {
			te.l.Info().Str("testName", testName).Str("filePath", path).Msg("Chainlink node coverage percentage report saved")
		}
	}

	return nil
}

// getChainlinkDir returns the path to the chainlink directory
func getChainlinkDir() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("cannot determine the path of the calling file")
	}
	dir := filepath.Dir(filename)
	chainlinkDir := filepath.Clean(filepath.Join(dir, "../../.."))
	return chainlinkDir, nil
}

func (te *CLClusterTestEnv) logWhetherAllContainersAreRunning() {
	for _, node := range te.ClCluster.Nodes {
		if node.Container == nil {
			continue
		}

		isCLRunning := node.Container.IsRunning()
		isDBRunning := node.PostgresDb.Container.IsRunning()

		if !isCLRunning {
			te.l.Warn().Str("Node", node.ContainerName).Msg("Chainlink node was not running, when test ended")
		}

		if !isDBRunning {
			te.l.Warn().Str("Node", node.ContainerName).Msg("Postgres DB is not running, when test ended")
		}
	}
}

func (te *CLClusterTestEnv) returnFunds() error {
	te.l.Info().Msg("Attempting to return Chainlink node funds to default network wallets")

	if len(te.evmClients) == 0 && len(te.sethClients) == 0 {
		return fmt.Errorf("both EVMClients and SethClient are nil, unable to return funds from chainlink nodes")
	}

	for _, evmClient := range te.evmClients {
		for _, chainlinkNode := range te.ClCluster.Nodes {
			fundedKeys, err := chainlinkNode.API.ExportEVMKeysForChain(evmClient.GetChainID().String())
			if err != nil {
				return err
			}
			for _, key := range fundedKeys {
				keyToDecrypt, err := json.Marshal(key)
				if err != nil {
					return err
				}
				// This can take up a good bit of RAM and time. When running on the remote-test-runner, this can lead to OOM
				// issues. So we avoid running in parallel; slower, but safer.
				decryptedKey, err := keystore.DecryptKey(keyToDecrypt, client.ChainlinkKeyPassword)
				if err != nil {
					return err
				}
				if te.evmClients[0] != nil {
					te.l.Debug().
						Str("ChainId", evmClient.GetChainID().String()).
						Msg("Returning funds from chainlink node")
					if err = evmClient.ReturnFunds(decryptedKey.PrivateKey); err != nil {
						// If we fail to return funds from one, go on to try the others anyway
						te.l.Error().Err(err).Str("Node", chainlinkNode.ContainerName).Msg("Error returning funds from node")
					}
				}
			}
		}
	}

	for _, sethClient := range te.sethClients {
		if err := actions_seth.ReturnFundsFromNodes(te.l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(te.ClCluster.NodeAPIs())); err != nil {
			te.l.Error().Err(err).Msg("Error returning funds from node")
		}
	}

	te.l.Info().Msg("Returned funds from Chainlink nodes")
	return nil
}

func (te *CLClusterTestEnv) GetEVMClient(chainId int64) (blockchain.EVMClient, error) {
	if len(te.sethClients) > 0 {
		return nil, fmt.Errorf("Environment is using Seth clients, not EVM clients")
	}

	if evmClient, ok := te.evmClients[chainId]; ok {
		return evmClient, nil
	}

	return nil, fmt.Errorf("no EVMClient available for chain ID %d", chainId)
}

func (te *CLClusterTestEnv) GetSethClient(chainId int64) (*seth.Client, error) {
	if len(te.evmClients) > 0 {
		return nil, fmt.Errorf("Environment is using EVMClients, not Seth clients")
	}
	if sethClient, ok := te.sethClients[chainId]; ok {
		return sethClient, nil
	}

	return nil, fmt.Errorf("no Seth client available for chain ID %d", chainId)
}

func (te *CLClusterTestEnv) GetSethClientForSelectedNetwork() (*seth.Client, error) {
	n := networks.MustGetSelectedNetworkConfig(te.TestConfig.GetNetworkConfig())[0]
	return te.GetSethClient(n.ChainID)
}

func (te *CLClusterTestEnv) GetRpcProvider(chainId int64) (*test_env.RpcProvider, error) {
	if rpc, ok := te.rpcProviders[chainId]; ok {
		return rpc, nil
	}

	return nil, fmt.Errorf("no RPC provider available for chain ID %d", chainId)
}
