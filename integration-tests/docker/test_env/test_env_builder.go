package test_env

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
)

type CLTestEnvBuilder struct {
	hasLogWatch          bool
	hasGeth              bool
	hasMockServer        bool
	hasForwarders        bool
	clNodeConfig         *chainlink.Config
	nonDevGethNetworks   []blockchain.EVMNetwork
	clNodesCount         int
	externalAdapterCount int
	customNodeCsaKeys    []string
	defaultNodeCsaKeys   []string
	l                    zerolog.Logger
	t                    *testing.T

	/* funding */
	ETHFunds *big.Float
}

func NewCLTestEnvBuilder() *CLTestEnvBuilder {
	return &CLTestEnvBuilder{
		externalAdapterCount: 1,
		l:                    log.Logger,
	}
}

func (b *CLTestEnvBuilder) WithTestLogger(t *testing.T) *CLTestEnvBuilder {
	b.t = t
	b.l = logging.GetTestLogger(t)
	return b
}

func (b *CLTestEnvBuilder) WithLogWatcher() *CLTestEnvBuilder {
	b.hasLogWatch = true
	return b
}

func (b *CLTestEnvBuilder) WithCLNodes(clNodesCount int) *CLTestEnvBuilder {
	b.clNodesCount = clNodesCount
	return b
}

func (b *CLTestEnvBuilder) WithForwarders() *CLTestEnvBuilder {
	b.hasForwarders = true
	return b
}

func (b *CLTestEnvBuilder) WithFunding(eth *big.Float) *CLTestEnvBuilder {
	b.ETHFunds = eth
	return b
}

func (b *CLTestEnvBuilder) WithGeth() *CLTestEnvBuilder {
	b.hasGeth = true
	return b
}

func (b *CLTestEnvBuilder) WithPrivateGethChains(evmNetworks []blockchain.EVMNetwork) *CLTestEnvBuilder {
	b.nonDevGethNetworks = evmNetworks
	return b
}

func (b *CLTestEnvBuilder) WithCLNodeConfig(cfg *chainlink.Config) *CLTestEnvBuilder {
	b.clNodeConfig = cfg
	return b
}

func (b *CLTestEnvBuilder) WithMockServer(externalAdapterCount int) *CLTestEnvBuilder {
	b.hasMockServer = true
	b.externalAdapterCount = externalAdapterCount
	return b
}

func (b *CLTestEnvBuilder) Build() (*CLClusterTestEnv, error) {
	envConfigPath, isSet := os.LookupEnv("TEST_ENV_CONFIG_PATH")
	if isSet {
		cfg, err := NewTestEnvConfigFromFile(envConfigPath)
		if err != nil {
			return nil, err
		}
		_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		return b.buildNewEnv(cfg)
	}
	return b.buildNewEnv(nil)
}

func (b *CLTestEnvBuilder) buildNewEnv(cfg *TestEnvConfig) (*CLClusterTestEnv, error) {
	b.l.Info().
		Bool("hasGeth", b.hasGeth).
		Bool("hasMockServer", b.hasMockServer).
		Int("externalAdapterCount", b.externalAdapterCount).
		Int("clNodesCount", b.clNodesCount).
		Strs("customNodeCsaKeys", b.customNodeCsaKeys).
		Strs("defaultNodeCsaKeys", b.defaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")

	var te *CLClusterTestEnv
	var err error
	if cfg != nil {
		te, err = NewTestEnvFromCfg(b.l, cfg)
		if err != nil {
			return nil, err
		}
	} else {
		te, err = NewTestEnv()
		if err != nil {
			return nil, err
		}
	}

	if b.t != nil {
		te.WithTestLogger(b.t)
	}

	if b.hasLogWatch {
		te.LogWatch, err = logwatch.NewLogWatch(nil, nil)
		if err != nil {
			return nil, err
		}
	}

	if b.hasMockServer {
		err = te.StartMockServer()
		if err != nil {
			return nil, err
		}
		err = te.MockServer.SetExternalAdapterMocks(b.externalAdapterCount)
		if err != nil {
			return nil, err
		}
	}
	if b.nonDevGethNetworks != nil {
		te.WithPrivateChain(b.nonDevGethNetworks)
		err := te.StartPrivateChain()
		if err != nil {
			return te, err
		}
		var nonDevNetworks []blockchain.EVMNetwork
		for i, n := range te.PrivateChain {
			primaryNode := n.GetPrimaryNode()
			if primaryNode == nil {
				return te, errors.WithStack(fmt.Errorf("Primary node is nil in PrivateChain interface"))
			}
			nonDevNetworks = append(nonDevNetworks, *n.GetNetworkConfig())
			nonDevNetworks[i].URLs = []string{primaryNode.GetInternalWsUrl()}
			nonDevNetworks[i].HTTPURLs = []string{primaryNode.GetInternalHttpUrl()}
		}
		if nonDevNetworks == nil {
			return nil, errors.New("cannot create nodes with custom config without nonDevNetworks")
		}

		err = te.StartClNodes(b.clNodeConfig, b.clNodesCount)
		if err != nil {
			return nil, err
		}
		return te, nil
	}
	networkConfig := networks.SelectedNetwork
	var internalDockerUrls test_env.InternalDockerUrls
	if b.hasGeth && networkConfig.Simulated {
		networkConfig, internalDockerUrls, err = te.StartGeth()
		if err != nil {
			return nil, err
		}

	}

	bc, err := blockchain.NewEVMClientFromNetwork(networkConfig, b.l)
	if err != nil {
		return nil, err
	}

	te.EVMClient = bc

	cd, err := contracts.NewContractDeployer(bc, b.l)
	if err != nil {
		return nil, err
	}
	te.ContractDeployer = cd

	cl, err := contracts.NewContractLoader(bc, b.l)
	if err != nil {
		return nil, err
	}
	te.ContractLoader = cl

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if b.clNodesCount > 0 {
		var cfg *chainlink.Config
		if b.clNodeConfig != nil {
			cfg = b.clNodeConfig
		} else {
			cfg = node.NewConfig(node.NewBaseConfig(),
				node.WithOCR1(),
				node.WithP2Pv1(),
			)
		}
		//node.SetDefaultSimulatedGeth(cfg, te.Geth.InternalWsUrl, te.Geth.InternalHttpUrl, b.hasForwarders)

		var httpUrls []string
		var wsUrls []string
		if networkConfig.Simulated {
			httpUrls = []string{internalDockerUrls.HttpUrl}
			wsUrls = []string{internalDockerUrls.WsUrl}
		} else {
			httpUrls = networkConfig.HTTPURLs
			wsUrls = networkConfig.URLs
		}

		node.SetChainConfig(cfg, wsUrls, httpUrls, networkConfig, b.hasForwarders)

		err := te.StartClNodes(cfg, b.clNodesCount)
		if err != nil {
			return nil, err
		}

		nodeCsaKeys, err = te.GetNodeCSAKeys()
		if err != nil {
			return nil, err
		}
		b.defaultNodeCsaKeys = nodeCsaKeys
	}

	if b.hasGeth && b.clNodesCount > 0 && b.ETHFunds != nil {
		te.ParallelTransactions(true)
		defer te.ParallelTransactions(false)
		if err := te.FundChainlinkNodes(b.ETHFunds); err != nil {
			return nil, err
		}
	}

	return te, nil
}
