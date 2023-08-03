package test_env

import (
	"sync"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/docker"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	tc "github.com/testcontainers/testcontainers-go"
	"go.uber.org/multierr"
)

type CLClusterTestEnv struct {
	cfg        *TestEnvConfig
	Network    *tc.DockerNetwork
	LogWatch   *logwatch.LogWatch
	CLNodes    []*ClNode
	Geth       *Geth
	MockServer *MockServer
}

func NewTestEnv() (*CLClusterTestEnv, error) {
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
	network, err := docker.CreateNetwork()
	if err != nil {
		return nil, err
	}
	networks := []string{network.Name}
	log.Info().Interface("Cfg", cfg).Send()
	return &CLClusterTestEnv{
		cfg:        cfg,
		Network:    network,
		Geth:       NewGeth(networks, WithContainerName(cfg.Geth.ContainerName)),
		MockServer: NewMockServer(networks, WithContainerName(cfg.MockServer.ContainerName)),
	}, nil
}

func (m *CLClusterTestEnv) StartGeth() error {
	return m.Geth.StartContainer(m.LogWatch)
}

func (m *CLClusterTestEnv) StartMockServer() error {
	return m.MockServer.StartContainer(m.LogWatch)
}

// StartClNodes start one bootstrap node and {count} OCR nodes
func (m *CLClusterTestEnv) StartClNodes(nodeConfig chainlink.Config, count int) error {
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
			if m.cfg != nil {
				nodeContainerName = m.cfg.Nodes[i].NodeContainerName
				dbContainerName = m.cfg.Nodes[i].DbContainerName
			}
			n := NewClNode([]string{m.Network.Name}, nodeConfig,
				WithNodeContainerName(nodeContainerName),
				WithDbContainerName(dbContainerName))
			err := n.StartContainer(m.LogWatch)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			} else {
				mu.Lock()
				m.CLNodes = append(m.CLNodes, n)
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
func (m *CLClusterTestEnv) ChainlinkNodeAddresses() ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, n := range m.CLNodes {
		primaryAddress, err := n.ChainlinkNodeAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, primaryAddress)
	}
	return addresses, nil
}

// FundChainlinkNodes will fund all the provided Chainlink nodes with a set amount of native currency
func (m *CLClusterTestEnv) FundChainlinkNodes(amount *big.Float) error {
	for _, cl := range m.CLNodes {
		if err := cl.Fund(m.Geth.EthClient, amount); err != nil {
			return err
		}
	}
	return m.Geth.EthClient.WaitForEvents()
}

func (m *CLClusterTestEnv) GetNodeCSAKeys() ([]string, error) {
	var keys []string
	for _, n := range m.CLNodes {
		csaKeys, err := n.GetNodeCSAKeys()
		if err != nil {
			return nil, err
		}
		keys = append(keys, csaKeys.Data[0].ID)
	}
	return keys, nil
}

func (m *CLClusterTestEnv) Terminate() error {
	// TESTCONTAINERS_RYUK_DISABLED=false by defualt so ryuk will remove all
	// the containers and the network
	return nil
}
