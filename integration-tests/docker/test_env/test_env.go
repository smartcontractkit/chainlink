package test_env

import (
	"sync"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	tc "github.com/testcontainers/testcontainers-go"
	"go.uber.org/multierr"
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
	CLNodes    []*ClNode
	Geth       *Geth
	MockServer *MockServer

	/* common contracts */
	LinkToken       contracts.LinkToken
	MockETHLinkFeed contracts.MockETHLINKFeed

	/* VRFv2 */
	CoordinatorV2    contracts.VRFCoordinatorV2
	LoadTestConsumer contracts.VRFv2LoadTestConsumer
	BHSV2            contracts.BlockHashStore
	/* VRFv1 */
	CoordinatorV1 contracts.VRFCoordinator
	ConsumerV1    contracts.VRFConsumer
	BHSV1         contracts.BlockHashStore
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

func (m *CLClusterTestEnv) ParallelTransactions(enabled bool) {
	m.Geth.EthClient.ParallelTransactions(enabled)
}

func (m *CLClusterTestEnv) StartGeth() error {
	return m.Geth.StartContainer()
}

func (m *CLClusterTestEnv) StartMockServer() error {
	return m.MockServer.StartContainer()
}

func (m *CLClusterTestEnv) GetAPIs() []*client.ChainlinkClient {
	clients := make([]*client.ChainlinkClient, 0)
	for _, c := range m.CLNodes {
		clients = append(clients, c.API)
	}
	return clients
}

// StartClNodes start one bootstrap node and {count} OCR nodes
func (m *CLClusterTestEnv) StartClNodes(nodeConfig *chainlink.Config, count int) error {
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
			if m.Cfg != nil {
				nodeContainerName = m.Cfg.Nodes[i].NodeContainerName
				dbContainerName = m.Cfg.Nodes[i].DbContainerName
			}
			n := NewClNode([]string{m.Network.Name}, nodeConfig,
				WithNodeContainerName(nodeContainerName),
				WithDbContainerName(dbContainerName),
				WithLogWatch(m.LogWatch))
			err := n.StartContainer()
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
			return errors.Wrap(err, ErrFundCLNode)
		}
	}
	return m.Geth.EthClient.WaitForEvents()
}

func (m *CLClusterTestEnv) GetNodeCSAKeys() ([]string, error) {
	var keys []string
	for _, n := range m.CLNodes {
		csaKeys, err := n.GetNodeCSAKeys()
		if err != nil {
			return nil, errors.Wrap(err, ErrGetNodeCSAKeys)
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
