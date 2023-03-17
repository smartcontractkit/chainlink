package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

var chainHeaders = []string{"ID", "Enabled", "Config"}

// chainCommand returns a cli.Command with subcommands for the given ChainClient.
// The chainId cli.Flag must be named "id", but may be String or Int.
func chainCommand(typ string, client ChainClient, chainID cli.Flag) cli.Command {
	if flagName := chainID.GetName(); flagName != "id" {
		panic(fmt.Errorf("chainID flag name must be 'id', got: %s", flagName))
	}
	lower := strings.ToLower(typ)
	return cli.Command{
		Name:  lower,
		Usage: fmt.Sprintf("Commands for handling %s chains", typ),
		Subcommands: cli.Commands{
			{
				Name:   "list",
				Usage:  fmt.Sprintf("List all existing %s chains", typ),
				Action: client.IndexChains,
			},
		},
	}
}

// ChainClient is a generic client interface for any type of chain.
type ChainClient interface {
	IndexChains(c *cli.Context) error
}

type chainClient[C chains.Config, R presenters.ChainResource[C], P TableRenderer, P2 ~[]P] struct {
	*Client
	path string
}

// newChainClient returns a new ChainClient for a particular type of chains.Config.
// P is a TableRenderer corresponding to R, and P2 is the slice variant (type P2 []P).
func newChainClient[C chains.Config, R presenters.ChainResource[C], P TableRenderer, P2 ~[]P](c *Client, name string) ChainClient {
	return &chainClient[C, R, P, P2]{
		Client: c,
		path:   "/v2/chains/" + name,
	}
}

// IndexChains returns all chains.
func (cli *chainClient[C, R, P, P2]) IndexChains(c *cli.Context) (err error) {
	var p P2
	return cli.getPage(cli.path, c.Int("page"), &p)
}
