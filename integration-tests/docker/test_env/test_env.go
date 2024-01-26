package test_env

import (
	"encoding/json"
	"fmt"
	"math/big"
	"runtime/debug"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/runid"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	core_testconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	ErrFundCLNode = "failed to fund CL node"
)

type CLClusterTestEnv struct {
	Cfg       *TestEnvConfig
	Network   *tc.DockerNetwork
	LogStream *logstream.LogStream

	/* components */
	ClCluster             *ClCluster
	PrivateChain          []test_env.PrivateChain // for tests using non-dev networks -- unify it with new approach
	MockAdapter           *test_env.Killgrave
	EVMClient             blockchain.EVMClient
	ContractDeployer      contracts.ContractDeployer
	ContractLoader        contracts.ContractLoader
	RpcProvider           test_env.RpcProvider
	PrivateEthereumConfig *test_env.EthereumNetwork // new approach to private chains, supporting eth1 and eth2
	l                     zerolog.Logger
	t                     *testing.T
}

func NewTestEnv() (*CLClusterTestEnv, error) {
	log.Logger = logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL")
	network, err := docker.CreateNetwork(log.Logger)
	if err != nil {
		return nil, err
	}
	return &CLClusterTestEnv{
		Network: network,
		l:       log.Logger,
	}, nil
}

// WithTestEnvConfig sets the test environment cfg.
// Sets up private ethereum chain and MockAdapter containers with the provided cfg.
func (te *CLClusterTestEnv) WithTestEnvConfig(cfg *TestEnvConfig) *CLClusterTestEnv {
	te.Cfg = cfg
	if cfg.MockAdapter.ContainerName != "" {
		n := []string{te.Network.Name}
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
	te.EVMClient.ParallelTransactions(enabled)
}

func (te *CLClusterTestEnv) WithPrivateChain(evmNetworks []blockchain.EVMNetwork) *CLClusterTestEnv {
	var chains []test_env.PrivateChain
	for _, evmNetwork := range evmNetworks {
		n := evmNetwork
		pgc := test_env.NewPrivateGethChain(&n, []string{te.Network.Name})
		if te.t != nil {
			pgc.GetPrimaryNode().WithTestInstance(te.t)
		}
		chains = append(chains, pgc)
		var privateChain test_env.PrivateChain
		switch n.SimulationType {
		case "besu":
			privateChain = test_env.NewPrivateBesuChain(&n, []string{te.Network.Name})
		default:
			privateChain = test_env.NewPrivateGethChain(&n, []string{te.Network.Name})
		}
		chains = append(chains, privateChain)
	}
	te.PrivateChain = chains
	return te
}

func (te *CLClusterTestEnv) StartPrivateChain() error {
	for _, chain := range te.PrivateChain {
		primaryNode := chain.GetPrimaryNode()
		if primaryNode == nil {
			return fmt.Errorf("primary node is nil in PrivateChain interface, stack: %s", string(debug.Stack()))
		}
		err := primaryNode.Start()
		if err != nil {
			return err
		}
		err = primaryNode.ConnectToClient()
		if err != nil {
			return err
		}
	}
	return nil
}

func (te *CLClusterTestEnv) StartEthereumNetwork(cfg *test_env.EthereumNetwork) (blockchain.EVMNetwork, test_env.RpcProvider, error) {
	// if environment is being restored from a previous state, use the existing config
	// this might fail terribly if temporary folders with chain data on the host machine were removed
	if te.Cfg != nil && te.Cfg.EthereumNetwork != nil {
		builder := test_env.NewEthereumNetworkBuilder()
		c, err := builder.WithExistingConfig(*te.Cfg.EthereumNetwork).
			WithTest(te.t).
			Build()
		if err != nil {
			return blockchain.EVMNetwork{}, test_env.RpcProvider{}, err
		}
		cfg = &c
	}
	n, rpc, err := cfg.Start()

	if err != nil {
		return blockchain.EVMNetwork{}, test_env.RpcProvider{}, err
	}

	return n, rpc, nil
}

func (te *CLClusterTestEnv) StartMockAdapter() error {
	return te.MockAdapter.StartContainer()
}

// pass config here
func (te *CLClusterTestEnv) StartClCluster(nodeConfig *chainlink.Config, count int, secretsConfig string, testconfig core_testconfig.GlobalTestConfig, opts ...ClNodeOption) error {
	if te.Cfg != nil && te.Cfg.ClCluster != nil {
		te.ClCluster = te.Cfg.ClCluster
	} else {
		opts = append(opts, WithSecrets(secretsConfig), WithLogStream(te.LogStream))
		te.ClCluster = &ClCluster{}
		for i := 0; i < count; i++ {
			ocrNode, err := NewClNode([]string{te.Network.Name}, *testconfig.GetChainlinkImageConfig().Image, *testconfig.GetChainlinkImageConfig().Version, nodeConfig, opts...)
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
	for _, cl := range te.ClCluster.Nodes {
		if err := cl.Fund(te.EVMClient, amount); err != nil {
			return fmt.Errorf("%s, err: %w", ErrFundCLNode, err)
		}
	}
	return te.EVMClient.WaitForEvents()
}

func (te *CLClusterTestEnv) Terminate() error {
	// TESTCONTAINERS_RYUK_DISABLED=false by default so ryuk will remove all
	// the containers and the Network
	return nil
}

// Cleanup cleans the environment up after it's done being used, mainly for returning funds when on live networks and logs.
func (te *CLClusterTestEnv) Cleanup() error {
	te.l.Info().Msg("Cleaning up test environment")

	runIdErr := runid.RemoveLocalRunId()
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

	if te.EVMClient == nil {
		return fmt.Errorf("evm client is nil, unable to return funds from chainlink nodes during cleanup")
	} else if te.EVMClient.NetworkSimulated() {
		te.l.Info().
			Str("Network Name", te.EVMClient.GetNetworkName()).
			Msg("Network is a simulated network. Skipping fund return.")
	} else {
		if err := te.returnFunds(); err != nil {
			return err
		}
	}

	// close EVMClient connections
	if te.EVMClient != nil {
		err := te.EVMClient.Close()
		return err
	}

	return nil
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
	for _, chainlinkNode := range te.ClCluster.Nodes {
		fundedKeys, err := chainlinkNode.API.ExportEVMKeysForChain(te.EVMClient.GetChainID().String())
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
			if err = te.EVMClient.ReturnFunds(decryptedKey.PrivateKey); err != nil {
				// If we fail to return funds from one, go on to try the others anyway
				te.l.Error().Err(err).Str("Node", chainlinkNode.ContainerName).Msg("Error returning funds from node")
			}
		}
	}

	te.l.Info().Msg("Returned funds from Chainlink nodes")
	return nil
}
