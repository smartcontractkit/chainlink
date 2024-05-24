package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func initCosmosNodeSubCmd(s *Shell) cli.Command {
	return nodeCommand("Cosmos", NewCosmosNodeClient(s))
}

func initStarkNetNodeSubCmd(s *Shell) cli.Command {
	return nodeCommand("StarkNet", NewStarkNetNodeClient(s))
}

func initEVMNodeSubCmd(s *Shell) cli.Command {
	return nodeCommand("EVM", NewEVMNodeClient(s))
}

func initSolanaNodeSubCmd(s *Shell) cli.Command {
	return nodeCommand("Solana", NewSolanaNodeClient(s))
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
func newNodeClient[P TableRenderer](s *Shell, name string) NodeClient {
	return &nodeClient[P]{
		Shell: s,
		path:  "/v2/nodes/" + name,
	}
}

// IndexNodes returns all nodes.
func (cli *nodeClient[P2]) IndexNodes(c *cli.Context) (err error) {
	var p P2
	return cli.getPage(cli.path, c.Int("page"), &p)
}

var nodeHeaders = []string{"Name", "Chain ID", "State", "Config"}
