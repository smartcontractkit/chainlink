package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func assertTableRenders(t *testing.T, r *cltest.RendererMock) {
	// Should be no error rendering any of the responses as tables
	b := bytes.NewBuffer([]byte{})
	tb := cmd.RendererTable{b}
	for _, rn := range r.Renders {
		require.NoError(t, tb.Render(rn))
	}
}

func TestShell_IndexEVMNodes(t *testing.T) {
	t.Parallel()

	chainID := newRandChainID()
	node1 := evmcfg.Node{
		Name:     ptr("Test node 1"),
		WSURL:    commonconfig.MustParseURL("ws://localhost:8546"),
		HTTPURL:  commonconfig.MustParseURL("http://localhost:8546"),
		SendOnly: ptr(false),
		Order:    ptr(int32(15)),
	}
	node2 := evmcfg.Node{
		Name:     ptr("Test node 2"),
		WSURL:    commonconfig.MustParseURL("ws://localhost:8547"),
		HTTPURL:  commonconfig.MustParseURL("http://localhost:8547"),
		SendOnly: ptr(false),
		Order:    ptr(int32(36)),
	}
	chain := evmcfg.EVMConfig{
		ChainID: chainID,
		Chain:   evmcfg.Defaults(chainID),
		Nodes:   evmcfg.EVMNodes{&node1, &node2},
	}
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM = evmcfg.EVMConfigs{&chain}
	})
	client, r := app.NewShellAndRenderer()

	require.Nil(t, cmd.NewEVMNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.EVMNodePresenters)
	require.Len(t, nodes, 2)
	n1 := nodes[0]
	n2 := nodes[1]
	assert.Equal(t, chainID.String(), n1.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(chainID.String(), *node1.Name), n1.ID)
	assert.Equal(t, *node1.Name, n1.Name)
	wantConfig, err := toml.Marshal(node1)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n1.Config)
	assert.Equal(t, chainID.String(), n2.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(chainID.String(), *node2.Name), n2.ID)
	assert.Equal(t, *node2.Name, n2.Name)
	wantConfig2, err := toml.Marshal(node2)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig2), n2.Config)
	assertTableRenders(t, r)

	//Render table and check the fields order
	b := new(bytes.Buffer)
	rt := cmd.RendererTable{b}
	require.NoError(t, nodes.RenderTable(rt))
	renderLines := strings.Split(b.String(), "\n")
	assert.Equal(t, 23, len(renderLines))
	assert.Contains(t, renderLines[2], "Name")
	assert.Contains(t, renderLines[2], n1.Name)
	assert.Contains(t, renderLines[3], "Chain ID")
	assert.Contains(t, renderLines[3], n1.ChainID)
	assert.Contains(t, renderLines[4], "State")
	assert.Contains(t, renderLines[4], n1.State)
	assert.Contains(t, renderLines[12], "Name")
	assert.Contains(t, renderLines[12], n2.Name)
	assert.Contains(t, renderLines[13], "Chain ID")
	assert.Contains(t, renderLines[13], n2.ChainID)
	assert.Contains(t, renderLines[14], "State")
	assert.Contains(t, renderLines[14], n2.State)
}
