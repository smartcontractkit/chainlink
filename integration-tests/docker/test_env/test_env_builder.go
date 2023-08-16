package test_env

import (
	"math/big"
	"os"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/networks"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
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

	/* funding */
	ETHFunds *big.Float
}

type InternalDockerUrls struct {
	HttpUrl string
	WsUrl   string
}

func NewCLTestEnvBuilder() *CLTestEnvBuilder {
	return &CLTestEnvBuilder{
		externalAdapterCount: 1,
	}
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

func (m *CLTestEnvBuilder) WithPrivateGethChains(evmNetworks []blockchain.EVMNetwork) *CLTestEnvBuilder {
	m.nonDevGethNetworks = evmNetworks
	return m
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
	log.Info().
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
		te, err = NewTestEnvFromCfg(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		te, err = NewTestEnv()
		if err != nil {
			return nil, err
		}
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
		te.WithPrivateGethChain(b.nonDevGethNetworks)
		err := te.StartPrivateGethChain()
		if err != nil {
			return te, err
		}
		var nonDevGethNetworks []blockchain.EVMNetwork
		for i, n := range te.PrivateGethChain {
			nonDevGethNetworks = append(nonDevGethNetworks, *n.NetworkConfig)
			nonDevGethNetworks[i].URLs = []string{n.PrimaryNode.InternalWsUrl}
			nonDevGethNetworks[i].HTTPURLs = []string{n.PrimaryNode.InternalHttpUrl}
		}
		if nonDevGethNetworks == nil {
			return nil, errors.New("cannot create nodes with custom config without nonDevGethNetworks")
		}

		err = te.StartClNodes(b.clNodeConfig, b.clNodesCount)
		if err != nil {
			return nil, err
		}
		return te, nil
	}
	networkConfig := networks.SelectedNetwork
	var internalDockerUrls InternalDockerUrls
	if b.hasGeth && networkConfig.Simulated {
		networkConfig, internalDockerUrls, err = te.StartGeth()
		if err != nil {
			return nil, err
		}

	}

	bc, err := blockchain.NewEVMClientFromNetwork(networkConfig)
	if err != nil {
		return nil, err
	}

	te.EVMClient = bc

	cd, err := contracts.NewContractDeployer(bc)
	if err != nil {
		return nil, err
	}
	te.ContractDeployer = cd

	cl, err := contracts.NewContractLoader(bc)
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
			cfg = node.NewConfig(node.BaseConf,
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
