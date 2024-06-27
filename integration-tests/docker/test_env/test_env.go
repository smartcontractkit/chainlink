package test_env

import (
	"math/big"
	"sync"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	tc "github.com/testcontainers/testcontainers-go"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/docker"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
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
	Geth             *Geth                       // for tests using --dev networks
	PrivateGethChain []test_env.PrivateGethChain // for tests using non-dev networks
	MockServer       *MockServer
	EVMClient        blockchain.EVMClient
	ContractDeployer contracts.ContractDeployer
	ContractLoader   contracts.ContractLoader
}

func NewTestEnv() (*CLClusterTestEnv, error) {
	utils.SetupCoreDockerEnvLogger()
	network, err := docker.CreateNetwork()
	if err != nil {
		return nil, err
	}
	networks := []string{network.Name}
	return &CLClusterTestEnv{
		Network:    network,
		Geth:       NewGeth(networks),
		MockServer: NewMockServer(networks),
	}, nil
}

func NewTestEnvFromCfg(cfg *TestEnvConfig) (*CLClusterTestEnv, error) {
	utils.SetupCoreDockerEnvLogger()
	network, err := docker.CreateNetwork()
	if err != nil {
		return nil, err
	}
	networks := []string{network.Name}
	log.Info().Interface("Cfg", cfg).Send()
	return &CLClusterTestEnv{
		Cfg:        cfg,
		Network:    network,
		Geth:       NewGeth(networks, WithContainerName(cfg.Geth.ContainerName)),
		MockServer: NewMockServer(networks, WithContainerName(cfg.MockServer.ContainerName)),
	}, nil
}

func (te *CLClusterTestEnv) ParallelTransactions(enabled bool) {
	te.EVMClient.ParallelTransactions(enabled)
}

func (m *CLClusterTestEnv) WithPrivateGethChain(evmNetworks []blockchain.EVMNetwork) *CLClusterTestEnv {
	var chains []test_env.PrivateGethChain
	for _, evmNetwork := range evmNetworks {
		n := evmNetwork
		chains = append(chains, test_env.NewPrivateGethChain(&n, []string{m.Network.Name}))
	}
	m.PrivateGethChain = chains
	return m
}

func (m *CLClusterTestEnv) StartPrivateGethChain() error {
	for _, chain := range m.PrivateGethChain {
		err := chain.PrimaryNode.Start()
		if err != nil {
			return err
		}
		err = chain.PrimaryNode.ConnectToClient()
		if err != nil {
			return err
		}
	}
	return nil
}

func (te *CLClusterTestEnv) StartGeth() (blockchain.EVMNetwork, InternalDockerUrls, error) {
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
func (te *CLClusterTestEnv) StartClNodes(nodeConfig *chainlink.Config, count int) error {
	var wg sync.WaitGroup
	var errs = []error{}
	var mu sync.Mutex

	// Start nodes
	for i := 0; i < count; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			var nodeContainerName, dbContainerName string
			if te.Cfg != nil {
				nodeContainerName = te.Cfg.Nodes[i].NodeContainerName
				dbContainerName = te.Cfg.Nodes[i].DbContainerName
			}
			n := NewClNode([]string{te.Network.Name}, nodeConfig,
				WithNodeContainerName(nodeContainerName),
				WithDbContainerName(dbContainerName),
			)
			err := n.StartContainer()
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			} else {
				mu.Lock()
				te.CLNodes = append(te.CLNodes, n)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return multierr.Combine(errs...)
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
	// TESTCONTAINERS_RYUK_DISABLED=false by defualt so ryuk will remove all
	// the containers and the Network
	return nil
}
