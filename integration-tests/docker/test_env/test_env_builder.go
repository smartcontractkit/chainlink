package test_env

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/utils/templates"
	"math/big"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
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

func (m *CLTestEnvBuilder) WithLogWatcher() *CLTestEnvBuilder {
	m.HasLogWatch = true
	return m
}

func (m *CLTestEnvBuilder) WithCLNodes(clNodesCount int) *CLTestEnvBuilder {
	m.CLNodesCount = clNodesCount
	return m
}

func (m *CLTestEnvBuilder) WithFunding(eth *big.Float) *CLTestEnvBuilder {
	m.ETHFunds = eth
	return m
}

func (m *CLTestEnvBuilder) WithGeth() *CLTestEnvBuilder {
	m.HasGeth = true
	return m
}

func (m *CLTestEnvBuilder) WithMockServer(externalAdapterCount int) *CLTestEnvBuilder {
	m.HasMockServer = true
	m.ExternalAdapterCount = externalAdapterCount
	return m
}

func (m *CLTestEnvBuilder) Build() (*CLClusterTestEnv, error) {
	envConfigPath, isSet := os.LookupEnv("TEST_ENV_CONFIG_PATH")
	if isSet {
		cfg, err := NewTestEnvConfigFromFile(envConfigPath)
		if err != nil {
			return nil, err
		}
		_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		return m.connectExistingEnv(cfg)
	}
	return m.buildNewEnv()
}

func (m *CLTestEnvBuilder) connectExistingEnv(cfg *TestEnvConfig) (*CLClusterTestEnv, error) {
	log.Info().
		Bool("HasGeth", m.HasGeth).
		Bool("HasMockServer", m.HasMockServer).
		Int("ExternalAdapterCount", m.ExternalAdapterCount).
		Int("CLNodesCount", m.CLNodesCount).
		Strs("CustomNodeCsaKeys", m.CustomNodeCsaKeys).
		Strs("DefaultNodeCsaKeys", m.DefaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")
	te, err := NewTestEnvFromCfg(cfg)
	if err != nil {
		return te, err
	}

	if m.HasLogWatch {
		lw, err := logwatch.NewLogWatch(nil, nil)
		if err != nil {
			return te, err
		}
		te.LogWatch = lw
	}

	if m.HasMockServer {
		err := te.StartMockServer()
		if err != nil {
			return te, err
		}
		err = te.MockServer.SetExternalAdapterMocks(m.ExternalAdapterCount)
		if err != nil {
			return te, err
		}
	}

	if m.HasGeth {
		err := te.StartGeth()
		if err != nil {
			return te, err
		}
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if m.CLNodesCount > 0 {
		// Create nodes
		nodeConfig := node.NewConfig(node.BaseConf,
			node.WithOCR1(),
			node.WithP2Pv1(),
			node.WithSimulatedEVM(te.Geth.InternalHttpUrl, te.Geth.InternalWsUrl),
		)
		err = te.StartClNodes(nodeConfig, m.clNodesCount)
		if err != nil {
			return te, err
		}

		nodeCsaKeys, err = te.GetNodeCSAKeys()
		if err != nil {
			return te, err
		}
		m.DefaultNodeCsaKeys = nodeCsaKeys
	}

	return te, nil
}

func (m *CLTestEnvBuilder) buildNewEnv() (*CLClusterTestEnv, error) {
	log.Info().
		Bool("HasGeth", m.HasGeth).
		Bool("HasMockServer", m.HasMockServer).
		Int("ExternalAdapterCount", m.ExternalAdapterCount).
		Int("CLNodesCount", m.CLNodesCount).
		Strs("CustomNodeCsaKeys", m.CustomNodeCsaKeys).
		Strs("DefaultNodeCsaKeys", m.DefaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")

	te, err := NewTestEnv()
	if err != nil {
		return te, err
	}

	if m.HasLogWatch {
		lw, err := logwatch.NewLogWatch(nil, nil)
		if err != nil {
			return te, err
		}
		te.LogWatch = lw
	}

	if m.HasMockServer {
		err := te.StartMockServer()
		if err != nil {
			return te, err
		}
		err = te.MockServer.SetExternalAdapterMocks(m.ExternalAdapterCount)
		if err != nil {
			return te, err
		}
	}

	if m.HasGeth {
		err := te.StartGeth()
		if err != nil {
			return te, err
		}
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if m.CLNodesCount > 0 {
		// Create nodes
		nodeConfig := node.NewConfig(node.BaseConf,
			node.WithOCR1(),
			node.WithP2Pv1(),
			node.WithSimulatedEVM(te.Geth.InternalHttpUrl, te.Geth.InternalWsUrl),
		)
		err = te.StartClNodes(nodeConfig, m.clNodesCount)
		if err != nil {
			return te, err
		}

		nodeCsaKeys, err = te.GetNodeCSAKeys()
		if err != nil {
			return te, err
		}
		m.DefaultNodeCsaKeys = nodeCsaKeys
	}

	if m.HasGeth && m.CLNodesCount > 0 && m.ETHFunds != nil {
		te.ParallelTransactions(true)
		defer te.ParallelTransactions(false)
		if err = te.FundChainlinkNodes(m.ETHFunds); err != nil {
			return te, err
		}
	}

	return te, nil
}
