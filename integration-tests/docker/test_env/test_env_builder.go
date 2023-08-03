package test_env

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"math/big"
)

type CLTestEnvBuilder struct {
	HasLogWatch          bool
	HasGeth              bool
	HasMockServer        bool
	CLNodesCount         int
	ExternalAdapterCount int
	CustomNodeCsaKeys    []string
	DefaultNodeCsaKeys   []string

	/* funding */
	ETHFunds *big.Float
}

func NewCLTestEnvBuilder() *CLTestEnvBuilder {
	return &CLTestEnvBuilder{
		ExternalAdapterCount: 1,
	}
}

func (b *CLTestEnvBuilder) WithLogWatcher() *CLTestEnvBuilder {
	b.HasLogWatch = true
	return b
}

func (b *CLTestEnvBuilder) WithCLNodes(clNodesCount int) *CLTestEnvBuilder {
	b.CLNodesCount = clNodesCount
	return b
}

func (b *CLTestEnvBuilder) WithFunding(eth *big.Float) *CLTestEnvBuilder {
	b.ETHFunds = eth
	return b
}

func (b *CLTestEnvBuilder) WithGeth() *CLTestEnvBuilder {
	b.HasGeth = true
	return b
}

func (b *CLTestEnvBuilder) WithMockServer(externalAdapterCount int) *CLTestEnvBuilder {
	b.HasMockServer = true
	b.ExternalAdapterCount = externalAdapterCount
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
		Bool("HasGeth", b.HasGeth).
		Bool("HasMockServer", b.HasMockServer).
		Int("ExternalAdapterCount", b.ExternalAdapterCount).
		Int("CLNodesCount", b.CLNodesCount).
		Strs("CustomNodeCsaKeys", b.CustomNodeCsaKeys).
		Strs("DefaultNodeCsaKeys", b.DefaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")
	te, err := NewTestEnvFromCfg(cfg)
	if err != nil {
		return te, err
	}

	if b.HasLogWatch {
		lw, err := logwatch.NewLogWatch(nil, nil)
		if err != nil {
			return te, err
		}
		te.LogWatch = lw
	}

	if b.HasMockServer {
		err := te.StartMockServer()
		if err != nil {
			return te, err
		}
		err = te.MockServer.SetExternalAdapterMocks(b.ExternalAdapterCount)
		if err != nil {
			return te, err
		}
	}

	if b.HasGeth {
		err := te.StartGeth()
		if err != nil {
			return te, err
		}
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if b.CLNodesCount > 0 {
		// Create nodes
		nodeConfig := node.NewConfig(node.BaseConf,
			node.WithOCR1(),
			node.WithP2Pv1(),
			node.WithSimulatedEVM(te.Geth.InternalHttpUrl, te.Geth.InternalWsUrl),
		)
		err = te.StartClNodes(nodeConfig, b.CLNodesCount)
		if err != nil {
			return te, err
		}

		nodeCsaKeys, err = te.GetNodeCSAKeys()
		if err != nil {
			return te, err
		}
		b.DefaultNodeCsaKeys = nodeCsaKeys
	}

	return te, nil
}

func (b *CLTestEnvBuilder) buildNewEnv() (*CLClusterTestEnv, error) {
	log.Info().
		Bool("HasGeth", b.HasGeth).
		Bool("HasMockServer", b.HasMockServer).
		Int("ExternalAdapterCount", b.ExternalAdapterCount).
		Int("CLNodesCount", b.CLNodesCount).
		Strs("CustomNodeCsaKeys", b.CustomNodeCsaKeys).
		Strs("DefaultNodeCsaKeys", b.DefaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")

	te, err := NewTestEnv()
	if err != nil {
		return te, err
	}

	if b.HasLogWatch {
		lw, err := logwatch.NewLogWatch(nil, nil)
		if err != nil {
			return te, err
		}
		te.LogWatch = lw
	}

	if b.HasMockServer {
		err := te.StartMockServer()
		if err != nil {
			return te, err
		}
		err = te.MockServer.SetExternalAdapterMocks(b.ExternalAdapterCount)
		if err != nil {
			return te, err
		}
	}

	if b.HasGeth {
		err := te.StartGeth()
		if err != nil {
			return te, err
		}
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if b.CLNodesCount > 0 {
		// Create nodes
		nodeConfig := node.NewConfig(node.BaseConf,
			node.WithOCR1(),
			node.WithP2Pv1(),
			node.WithSimulatedEVM(te.Geth.InternalHttpUrl, te.Geth.InternalWsUrl),
		)
		err = te.StartClNodes(nodeConfig, b.CLNodesCount)
		if err != nil {
			return te, err
		}

		nodeCsaKeys, err = te.GetNodeCSAKeys()
		if err != nil {
			return te, err
		}
		b.DefaultNodeCsaKeys = nodeCsaKeys
	}

	if b.HasGeth && b.CLNodesCount > 0 && b.ETHFunds != nil {
		te.ParallelTransactions(true)
		defer te.ParallelTransactions(false)
		if err = te.FundChainlinkNodes(b.ETHFunds); err != nil {
			return te, err
		}
	}

	return te, nil
}
