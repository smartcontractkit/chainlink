package cmd

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initNodeSubCmds(s *Shell) []cli.Command {
	cmds := []cli.Command{}
	for _, network := range slices.Sorted(maps.Keys(relay.SupportedNetworks)) {
		if network == relay.NetworkDummy {
			continue
		}
		cmds = append(cmds, nodeCommand(network, NewNodeClient(s, network)))
	}
	return cmds
}

// nodeCommand returns a cli.Command with subcommands for the given NodeClient.
// A string cli.Flag for "name" is automatically included.
func nodeCommand(typ string, client NodeClient) cli.Command {
	lower := strings.ToLower(typ)
	return cli.Command{
		Name:  lower,
		Usage: fmt.Sprintf("Commands for handling %s node configuration", typ),
		Subcommands: cli.Commands{
			{
				Name:   "list",
				Usage:  fmt.Sprintf("List all existing %s nodes", typ),
				Action: client.IndexNodes,
			},
		},
	}
}

// EVMNodePresenter implements TableRenderer for an EVMNodeResource.
type NodePresenter struct {
	presenters.NodeResource
}

// ToRow presents the EVMNodeResource as a slice of strings.
func (p *NodePresenter) ToRow() []string {
	return []string{p.Name, p.ChainID, p.State, p.Config}
}

// RenderTable implements TableRenderer
func (p NodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

// NodePresenters implements TableRenderer for a slice of NodePresenter.
type NodePresenters []NodePresenter

// RenderTable implements TableRenderer
func (ps NodePresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

// NodeClient is a generic client interface for any of node.
type NodeClient interface {
	IndexNodes(c *cli.Context) error
}

type nodeClient[P TableRenderer] struct {
	*Shell
	path string
}

// newNodeClient returns a new NodeClient for a particular type of NodeStatus.
// P is a TableRenderer for []types.NodeStatus.
func NewNodeClient(s *Shell, network string) NodeClient {
	return &nodeClient[NodePresenters]{
		Shell: s,
		path:  "/v2/nodes/" + network,
	}
}

// IndexNodes returns all nodes.
func (cli *nodeClient[P2]) IndexNodes(c *cli.Context) (err error) {
	var p P2
	return cli.getPage(cli.path, c.Int("page"), &p)
}

var nodeHeaders = []string{"Name", "Chain ID", "State", "Config"}
