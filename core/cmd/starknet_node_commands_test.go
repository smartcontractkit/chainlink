package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncfg "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func starknetStartNewApplication(t *testing.T, cfgs ...*config.TOMLConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Starknet = cfgs
		c.EVM = nil
		c.Solana = nil
	})
}

func TestShell_IndexStarkNetNodes(t *testing.T) {
	t.Parallel()

	id := "starknet chain ID"
	node1 := config.Node{
		Name: ptr("first"),
		URL:  commoncfg.MustParseURL("https://starknet1.example"),
	}
	node2 := config.Node{
		Name: ptr("second"),
		URL:  commoncfg.MustParseURL("https://starknet2.example"),
	}
	chain := config.TOMLConfig{
		ChainID: &id,
		Nodes:   config.Nodes{&node1, &node2},
	}
	app := starknetStartNewApplication(t, &chain)
	client, r := app.NewShellAndRenderer()

	require.Nil(t, cmd.NewStarkNetNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.StarkNetNodePresenters)
	require.Len(t, nodes, 2)
	n1 := nodes[0]
	n2 := nodes[1]
	assert.Equal(t, id, n1.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(id, *node1.Name), n1.ID)
	assert.Equal(t, *node1.Name, n1.Name)
	wantConfig, err := toml.Marshal(node1)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n1.Config)
	assert.Equal(t, id, n2.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(id, *node2.Name), n2.ID)
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
	assert.Equal(t, 17, len(renderLines))
	assert.Contains(t, renderLines[2], "Name")
	assert.Contains(t, renderLines[2], n1.Name)
	assert.Contains(t, renderLines[3], "Chain ID")
	assert.Contains(t, renderLines[3], n1.ChainID)
	assert.Contains(t, renderLines[4], "State")
	assert.Contains(t, renderLines[4], n1.State)
	assert.Contains(t, renderLines[9], "Name")
	assert.Contains(t, renderLines[9], n2.Name)
	assert.Contains(t, renderLines[10], "Chain ID")
	assert.Contains(t, renderLines[10], n2.ChainID)
	assert.Contains(t, renderLines[11], "State")
	assert.Contains(t, renderLines[11], n2.State)
}
