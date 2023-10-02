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
	"github.com/ethereum/go-ethereum/common"
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

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
)

var (
	ErrFundCLNode     = "failed to fund CL node"
	ErrGetNodeCSAKeys = "failed get CL node CSA keys"
)

type CLClusterTestEnv struct {
	Cfg      *TestEnvConfig
	Network  *tc.DockerNetwork
	LogWatch *logwatch.LogWatch

	/* components */
	CLNodes          []*ClNode
	Geth             *test_env.Geth          // for tests using --dev networks
	PrivateChain     []test_env.PrivateChain // for tests using non-dev networks
	MockServer       *test_env.MockServer
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
		Geth:       test_env.NewGeth(n),
		MockServer: test_env.NewMockServer(n),
		Network:    network,
		l:          log.Logger,
	}, nil
}

// WithTestEnvConfig sets the test environment cfg.
// Sets up the Geth and MockServer containers with the provided cfg.
func (te *CLClusterTestEnv) WithTestEnvConfig(cfg *TestEnvConfig) *CLClusterTestEnv {
	te.Cfg = cfg
	n := []string{te.Network.Name}
	te.Geth = test_env.NewGeth(n, test_env.WithContainerName(cfg.Geth.ContainerName))
	te.MockServer = test_env.NewMockServer(n, test_env.WithContainerName(cfg.MockServer.ContainerName))
	return te
}

func (te *CLClusterTestEnv) WithTestLogger(t *testing.T) *CLClusterTestEnv {
	te.t = t
	te.l = logging.GetTestLogger(t)
	te.Geth.WithTestLogger(t)
	te.MockServer.WithTestLogger(t)
	return te
}

func (te *CLClusterTestEnv) ParallelTransactions(enabled bool) {
	te.EVMClient.ParallelTransactions(enabled)
}

func (te *CLClusterTestEnv) WithPrivateChain(evmNetworks []blockchain.EVMNetwork) *CLClusterTestEnv {
	var chains []test_env.PrivateChain
	for _, evmNetwork := range evmNetworks {
		n := evmNetwork
		var privateChain test_env.PrivateChain
		switch n.SimulationType {
		case "besu":
			privateChain = test_env.NewPrivateBesuChain(&n, []string{te.Network.Name})
			if te.t != nil {
				privateChain.GetPrimaryNode().WithTestLogger(te.t)
			}
		default:
			privateChain = test_env.NewPrivateGethChain(&n, []string{te.Network.Name})
			if te.t != nil {
				privateChain.GetPrimaryNode().WithTestLogger(te.t)
			}
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

func (te *CLClusterTestEnv) StartMockServer() error {
	return te.MockServer.StartContainer()
}

func (te *CLClusterTestEnv) GetAPIs() []*client.ChainlinkClient {
	clients := make([]*client.ChainlinkClient, 0)
	for _, c := range te.CLNodes {
		clients = append(clients, c.API)
	}
	return clients
}

// StartClNodes start one bootstrap node and {count} OCR nodes
func (te *CLClusterTestEnv) StartClNodes(nodeConfig *chainlink.Config, count int, secretsConfig string) error {
	eg := &errgroup.Group{}
	nodes := make(chan *ClNode, count)

	// Start nodes
	for i := 0; i < count; i++ {
		nodeIndex := i
		eg.Go(func() error {
			var nodeContainerName, dbContainerName string
			if te.Cfg != nil {
				nodeContainerName = te.Cfg.Nodes[nodeIndex].NodeContainerName
				dbContainerName = te.Cfg.Nodes[nodeIndex].DbContainerName
			}
			n := NewClNode([]string{te.Network.Name}, nodeConfig,
				WithSecrets(secretsConfig),
				WithNodeContainerName(nodeContainerName),
				WithDbContainerName(dbContainerName),
			)
			if te.t != nil {
				n.WithTestLogger(te.t)
			}
			err := n.StartContainer()
			if err != nil {
				return err
			}
			nodes <- n
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	close(nodes)

	for node := range nodes {
		te.CLNodes = append(te.CLNodes, node)
	}

	return nil
}

// ChainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func (te *CLClusterTestEnv) ChainlinkNodeAddresses() ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, n := range te.CLNodes {
		primaryAddress, err := n.ChainlinkNodeAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, primaryAddress)
	}
	return addresses, nil
}

// FundChainlinkNodes will fund all the provided Chainlink nodes with a set amount of native currency
func (te *CLClusterTestEnv) FundChainlinkNodes(amount *big.Float) error {
	for _, cl := range te.CLNodes {
		if err := cl.Fund(te.EVMClient, amount); err != nil {
			return errors.Wrap(err, ErrFundCLNode)
		}
	}
	return te.EVMClient.WaitForEvents()
}

func (te *CLClusterTestEnv) GetNodeCSAKeys() ([]string, error) {
	var keys []string
	for _, n := range te.CLNodes {
		csaKeys, err := n.GetNodeCSAKeys()
		if err != nil {
			return nil, errors.Wrap(err, ErrGetNodeCSAKeys)
		}
		keys = append(keys, csaKeys.Data[0].ID)
	}
	return keys, nil
}

func (te *CLClusterTestEnv) Terminate() error {
	// TESTCONTAINERS_RYUK_DISABLED=false by default so ryuk will remove all
	// the containers and the Network
	return nil
}

// Cleanup cleans the environment up after it's done being used, mainly for returning funds when on live networks.
// Intended to be used as part of t.Cleanup() in tests.
func (te *CLClusterTestEnv) Cleanup(t *testing.T) error {
	if te.EVMClient == nil {
		return errors.New("blockchain client is nil, unable to return funds from chainlink nodes")
	}
	if te.CLNodes == nil {
		return errors.New("chainlink nodes are nil, unable to return funds from chainlink nodes")
	}

	// TODO: This is an imperfect and temporary solution, see TT-590 for a more sustainable solution
	// Collect logs if the test fails, or if we just want them
	if t.Failed() || os.Getenv("TEST_LOG_COLLECT") == "true" {
		folder := fmt.Sprintf("./logs/%s-%s", t.Name(), time.Now().Format("2006-01-02T15-04-05"))
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return err
		}

		te.l.Info().Msg("Collecting test logs")
		eg := &errgroup.Group{}
		for _, n := range te.CLNodes {
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
	}

	// Check if we need to return funds
	if te.EVMClient.NetworkSimulated() {
		te.l.Info().Str("Network Name", te.EVMClient.GetNetworkName()).
			Msg("Network is a simulated network. Skipping fund return.")
	} else {
		te.l.Info().Msg("Attempting to return Chainlink node funds to default network wallets")
		for _, chainlinkNode := range te.CLNodes {
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
					return err
				}
			}
		}
	}

	return nil
}
