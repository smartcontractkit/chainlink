package test_env

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

var (
	ErrFundCLNode = "failed to fund CL node"
)

type CLClusterTestEnv struct {
	Cfg      *TestEnvConfig
	Network  *tc.DockerNetwork
	LogWatch *logwatch.LogWatch

	/* components */
	ClCluster        *ClCluster
	Geth             *test_env.Geth          // for tests using --dev networks
	PrivateChain     []test_env.PrivateChain // for tests using non-dev networks
	MockAdapter      *test_env.Killgrave
	EVMClient        blockchain.EVMClient
	ContractDeployer contracts.ContractDeployer
	ContractLoader   contracts.ContractLoader
	l                zerolog.Logger
	t                *testing.T
}

func NewTestEnv() (*CLClusterTestEnv, error) {
	utils.SetupCoreDockerEnvLogger()
	network, err := docker.CreateNetwork(log.Logger)
	if err != nil {
		return nil, err
	}
	n := []string{network.Name}
	return &CLClusterTestEnv{
		Geth:        test_env.NewGeth(n),
		MockAdapter: test_env.NewKillgrave(n, ""),
		Network:     network,
		l:           log.Logger,
	}, nil
}

// WithTestEnvConfig sets the test environment cfg.
// Sets up the Geth and MockAdapter containers with the provided cfg.
func (te *CLClusterTestEnv) WithTestEnvConfig(cfg *TestEnvConfig) *CLClusterTestEnv {
	te.Cfg = cfg
	n := []string{te.Network.Name}
	te.Geth = test_env.NewGeth(n, test_env.WithContainerName(te.Cfg.Geth.ContainerName))
	te.MockAdapter = test_env.NewKillgrave(n, te.Cfg.MockAdapter.ImpostersPath, test_env.WithContainerName(te.Cfg.MockAdapter.ContainerName))
	return te
}

func (te *CLClusterTestEnv) WithTestLogger(t *testing.T) *CLClusterTestEnv {
	te.t = t
	te.l = logging.GetTestLogger(t)
	te.Geth.WithTestLogger(t)
	te.MockAdapter.WithTestLogger(t)
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
			pgc.GetPrimaryNode().WithTestLogger(te.t)
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
			return errors.WithStack(fmt.Errorf("primary node is nil in PrivateChain interface"))
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

func (te *CLClusterTestEnv) StartGeth() (blockchain.EVMNetwork, test_env.InternalDockerUrls, error) {
	return te.Geth.StartContainer()
}

func (te *CLClusterTestEnv) StartMockAdapter() error {
	return te.MockAdapter.StartContainer()
}

func (te *CLClusterTestEnv) StartClCluster(nodeConfig *chainlink.Config, count int, secretsConfig string) error {
	if te.Cfg != nil && te.Cfg.ClCluster != nil {
		te.ClCluster = te.Cfg.ClCluster
	} else {
		te.ClCluster = &ClCluster{}
		for i := 0; i < count; i++ {
			ocrNode := NewClNode([]string{te.Network.Name}, nodeConfig,
				WithSecrets(secretsConfig),
			)
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
			return errors.Wrap(err, ErrFundCLNode)
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
	if te.t == nil {
		return errors.New("cannot cleanup test environment without a testing.T")
	}
	if te.ClCluster == nil || len(te.ClCluster.Nodes) == 0 {
		return errors.New("chainlink nodes are nil, unable cleanup chainlink nodes")
	}

	// TODO: This is an imperfect and temporary solution, see TT-590 for a more sustainable solution
	// Collect logs if the test fails, or if we just want them
	if te.t.Failed() || os.Getenv("TEST_LOG_COLLECT") == "true" {
		if err := te.collectTestLogs(); err != nil {
			return err
		}
	}

	if te.EVMClient == nil {
		return errors.New("evm client is nil, unable to return funds from chainlink nodes during cleanup")
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

// collectTestLogs collects the logs from all the Chainlink nodes in the test environment and writes them to local files
func (te *CLClusterTestEnv) collectTestLogs() error {
	te.l.Info().Msg("Collecting test logs")
	folder := fmt.Sprintf("./logs/%s-%s", te.t.Name(), time.Now().Format("2006-01-02T15-04-05"))
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	eg := &errgroup.Group{}
	for _, n := range te.ClCluster.Nodes {
		node := n
		eg.Go(func() error {
			logFileName := filepath.Join(folder, fmt.Sprintf("node-%s.log", node.ContainerName))
			logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer logFile.Close()
			logReader, err := node.Container.Logs(context.Background())
			if err != nil {
				return err
			}
			_, err = io.Copy(logFile, logReader)
			if err != nil {
				return err
			}
			te.l.Info().Str("Node", node.ContainerName).Str("File", logFileName).Msg("Wrote Logs")
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	te.l.Info().Str("Logs Location", folder).Msg("Wrote test logs")
	return nil
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
