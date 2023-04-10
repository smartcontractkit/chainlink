package cmd_test

import (
	"bytes"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func assertTableRenders(t *testing.T, r *cltest.RendererMock) {
	// Should be no error rendering any of the responses as tables
	b := bytes.NewBuffer([]byte{})
	tb := cmd.RendererTable{b}
	for _, rn := range r.Renders {
		require.NoError(t, tb.Render(rn))
	}
}

func TestClient_IndexEVMNodes(t *testing.T) {
	t.Parallel()

	chainID := newRandChainID()
	node1 := evmcfg.Node{
		Name:     ptr("Test node 1"),
		WSURL:    models.MustParseURL("ws://localhost:8546"),
		HTTPURL:  models.MustParseURL("http://localhost:8546"),
		SendOnly: ptr(false),
	}
	node2 := evmcfg.Node{
		Name:     ptr("Test node 2"),
		WSURL:    models.MustParseURL("ws://localhost:8547"),
		HTTPURL:  models.MustParseURL("http://localhost:8547"),
		SendOnly: ptr(false),
	}
	chain := evmcfg.EVMConfig{
		ChainID: chainID,
		Chain:   evmcfg.Defaults(chainID),
		Nodes:   evmcfg.EVMNodes{&node1, &node2},
	}
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM = evmcfg.EVMConfigs{&chain}
	})
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.NewEVMNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.EVMNodePresenters)
	require.Len(t, nodes, 2)
	n1 := nodes[0]
	n2 := nodes[1]
	assert.Equal(t, chainID.String(), n1.ChainID)
	assert.Equal(t, *node1.Name, n1.ID)
	assert.Equal(t, *node1.Name, n1.Name)
	wantConfig, err := toml.Marshal(node1)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n1.Config)
	assert.Equal(t, chainID.String(), n2.ChainID)
	assert.Equal(t, *node2.Name, n2.ID)
	assert.Equal(t, *node2.Name, n2.Name)
	wantConfig2, err := toml.Marshal(node2)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig2), n2.Config)
	assertTableRenders(t, r)
}
