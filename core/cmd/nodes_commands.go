package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/chains"
)

func initTerraNodeSubCmd(client *Client) cli.Command {
	return nodeCommand("Terra", NewTerraNodeClient(client),
		cli.StringFlag{
			Name:  "chain-id",
			Usage: "chain ID",
		},
		cli.StringFlag{
			Name:  "tendermint-url",
			Usage: "Tendermint URL",
		})
}

func initStarkNetNodeSubCmd(client *Client) cli.Command {
	return nodeCommand("StarkNet", NewStarkNetNodeClient(client),
		cli.StringFlag{
			Name:  "chain-id",
			Usage: "chain ID",
		},
		cli.StringFlag{
			Name:  "url",
			Usage: "URL",
		})
}

func initEVMNodeSubCmd(client *Client) cli.Command {
	return nodeCommand("EVM", NewEVMNodeClient(client),
		cli.StringFlag{
			Name:  "ws-url",
			Usage: "Websocket URL",
		},
		cli.StringFlag{
			Name:  "http-url",
			Usage: "HTTP URL, optional",
		},
		cli.Int64Flag{
			Name:  "chain-id",
			Usage: "chain ID",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "primary|secondary",
		})
}

func initSolanaNodeSubCmd(client *Client) cli.Command {
	return nodeCommand("Solana", NewSolanaNodeClient(client),
		cli.StringFlag{
			Name:  "chain-id",
			Usage: "chain ID, options: [mainnet, testnet, devnet, localnet]",
		},
		cli.StringFlag{
			Name:  "url",
			Usage: "URL",
		})
}

// nodeCommand returns a cli.Command with subcommands for the given NodeClient.
// A string cli.Flag for "name" is automatically included.
func nodeCommand(typ string, client NodeClient, flags ...cli.Flag) cli.Command {
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

type nodeClient[N chains.Node, P TableRenderer, P2 ~[]P] struct {
	*Client
	path       string
	createNode func(c *cli.Context) (N, error)
}

// newNodeClient returns a new NodeClient for a particular type of chains.Node.
// P is a TableRenderer corresponding to R, and P2 is the slice variant (type P2 []P).
func newNodeClient[N chains.Node, P TableRenderer, P2 ~[]P](c *Client, name string, createNode func(*cli.Context) (N, error)) NodeClient {
	return &nodeClient[N, P, P2]{
		Client:     c,
		path:       "/v2/nodes/" + name,
		createNode: createNode,
	}
}

// IndexNodes returns all nodes.
func (cli *nodeClient[N, P, P2]) IndexNodes(c *cli.Context) (err error) {
	var p P2
	return cli.getPage(cli.path, c.Int("page"), &p)
}
