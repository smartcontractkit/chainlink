package test_env

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink/integration-tests/types/node"
	"os"
)

type CLTestEnvBuilder struct {
	reusable             bool
	hasLogWatch          bool
	hasGeth              bool
	hasMockServer        bool
	clNodesCount         int
	externalAdapterCount int
	customNodeCsaKeys    []string
	defaultNodeCsaKeys   []string
}

func NewCLTestEnvBuilder() *CLTestEnvBuilder {
	return &CLTestEnvBuilder{
		externalAdapterCount: 1,
	}
}

func (m *CLTestEnvBuilder) WithLogWatcher() *CLTestEnvBuilder {
	m.hasLogWatch = true
	return m
}

func (m *CLTestEnvBuilder) WithCLNodes(clNodesCount int) *CLTestEnvBuilder {
	m.clNodesCount = clNodesCount
	return m
}

func (m *CLTestEnvBuilder) WithGeth() *CLTestEnvBuilder {
	m.hasGeth = true
	return m
}

func (m *CLTestEnvBuilder) WithReusable() *CLTestEnvBuilder {
	m.reusable = true
	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	return m
}

func (m *CLTestEnvBuilder) WithMockServer(externalAdapterCount int) *CLTestEnvBuilder {
	m.hasMockServer = true
	m.externalAdapterCount = externalAdapterCount
	return m
}

func (m *CLTestEnvBuilder) Build() (*CLClusterTestEnv, error) {
	log.Info().
		Bool("hasGeth", m.hasGeth).
		Bool("hasMockServer", m.hasMockServer).
		Int("externalAdapterCount", m.externalAdapterCount).
		Int("clNodesCount", m.clNodesCount).
		Strs("customNodeCsaKeys", m.customNodeCsaKeys).
		Strs("defaultNodeCsaKeys", m.defaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")

	te, err := NewTestEnv(m.reusable)
	if err != nil {
		return te, err
	}

	if m.hasLogWatch {
		lw, err := logwatch.NewLogWatch(nil, nil)
		if err != nil {
			return te, err
		}
		te.LogWatch = lw
	}

	if m.hasMockServer {
		err := te.StartMockServer()
		if err != nil {
			return te, err
		}
		err = te.MockServer.SetExternalAdapterMocks(m.externalAdapterCount)
		if err != nil {
			return te, err
		}
	}

	if m.hasGeth {
		err := te.StartGeth()
		if err != nil {
			return te, err
		}
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if m.clNodesCount > 0 {
		// Create nodes
		nodeConfOpts := node.NodeConfigOpts{
			EVM: struct {
				HttpUrl string
				WsUrl   string
			}{
				HttpUrl: te.Geth.InternalHttpUrl,
				WsUrl:   te.Geth.InternalWsUrl,
			},
		}
		err = te.StartClNodes(nodeConfOpts, m.clNodesCount)
		if err != nil {
			return te, err
		}

		nodeCsaKeys, err = te.GetNodeCSAKeys()
		if err != nil {
			return te, err
		}
		m.defaultNodeCsaKeys = nodeCsaKeys
	}

	return te, nil
}
