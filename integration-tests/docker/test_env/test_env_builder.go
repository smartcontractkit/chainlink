package test_env

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"math/big"
)

type CLTestEnvBuilder struct {
	hasLogWatch          bool
	hasGeth              bool
	hasMockServer        bool
	clNodesCount         int
	externalAdapterCount int
	customNodeCsaKeys    []string
	defaultNodeCsaKeys   []string

	/* funding */
	ETHFunds *big.Float
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

func (b *CLTestEnvBuilder) WithFunding(eth *big.Float) *CLTestEnvBuilder {
	b.ETHFunds = eth
	return b
}

func (b *CLTestEnvBuilder) WithGeth() *CLTestEnvBuilder {
	b.hasGeth = true
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
		return b.connectExistingEnv(cfg)
	}
	return b.buildNewEnv()
}

func (b *CLTestEnvBuilder) connectExistingEnv(cfg *TestEnvConfig) (*CLClusterTestEnv, error) {
	log.Info().
		Bool("hasGeth", b.hasGeth).
		Bool("hasMockServer", b.hasMockServer).
		Int("externalAdapterCount", b.externalAdapterCount).
		Int("clNodesCount", b.clNodesCount).
		Strs("customNodeCsaKeys", b.customNodeCsaKeys).
		Strs("defaultNodeCsaKeys", b.defaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")
	te, err := NewTestEnvFromCfg(cfg)
	if err != nil {
		return te, err
	}

	if b.hasLogWatch {
		lw, err := logwatch.NewLogWatch(nil, nil)
		if err != nil {
			return te, err
		}
		te.LogWatch = lw
	}

	if b.hasMockServer {
		err := te.StartMockServer()
		if err != nil {
			return te, err
		}
		err = te.MockServer.SetExternalAdapterMocks(b.externalAdapterCount)
		if err != nil {
			return te, err
		}
	}

	if b.hasGeth {
		err := te.StartGeth()
		if err != nil {
			return te, err
		}
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if b.clNodesCount > 0 {
		// Create nodes
		nodeConfig := node.NewConfig(node.BaseConf,
			node.WithOCR1(),
			node.WithP2Pv1(),
			node.WithSimulatedEVM(te.Geth.InternalHttpUrl, te.Geth.InternalWsUrl),
		)
		err = te.StartClNodes(nodeConfig, b.clNodesCount)
		if err != nil {
			return te, err
		}

		nodeCsaKeys, err = te.GetNodeCSAKeys()
		if err != nil {
			return te, err
		}
		b.defaultNodeCsaKeys = nodeCsaKeys
	}

	return te, nil
}

func (b *CLTestEnvBuilder) buildNewEnv() (*CLClusterTestEnv, error) {
	log.Info().
		Bool("hasGeth", b.hasGeth).
		Bool("hasMockServer", b.hasMockServer).
		Int("externalAdapterCount", b.externalAdapterCount).
		Int("clNodesCount", b.clNodesCount).
		Strs("customNodeCsaKeys", b.customNodeCsaKeys).
		Strs("defaultNodeCsaKeys", b.defaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")

	te, err := NewTestEnv()
	if err != nil {
		return te, err
	}

	if b.hasLogWatch {
		lw, err := logwatch.NewLogWatch(nil, nil)
		if err != nil {
			return te, err
		}
		te.LogWatch = lw
	}

	if b.hasMockServer {
		err := te.StartMockServer()
		if err != nil {
			return te, err
		}
		err = te.MockServer.SetExternalAdapterMocks(b.externalAdapterCount)
		if err != nil {
			return te, err
		}
	}

	if b.hasGeth {
		err := te.StartGeth()
		if err != nil {
			return te, err
		}
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if b.clNodesCount > 0 {
		// Create nodes
		nodeConfig := node.NewConfig(node.BaseConf,
			node.WithOCR1(),
			node.WithP2Pv1(),
			node.WithSimulatedEVM(te.Geth.InternalHttpUrl, te.Geth.InternalWsUrl),
		)
		err = te.StartClNodes(nodeConfig, b.clNodesCount)
		if err != nil {
			return te, err
		}

		nodeCsaKeys, err = te.GetNodeCSAKeys()
		if err != nil {
			return te, err
		}
		b.defaultNodeCsaKeys = nodeCsaKeys
	}

	if b.hasGeth && b.clNodesCount > 0 && b.ETHFunds != nil {
		te.ParallelTransactions(true)
		defer te.ParallelTransactions(false)
		if err = te.FundChainlinkNodes(b.ETHFunds); err != nil {
			return te, err
		}
	}

	return te, nil
}
