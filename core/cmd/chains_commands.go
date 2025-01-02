package cmd

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
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

type chainClient[P TableRenderer] struct {
	*Shell
	path string
}

// NewChainClient returns a new ChainClient for a particular type of chains.Config.
// P is a TableRenderer corresponding to R, and P2 is the slice variant (type P2 []P).
func NewChainClient(s *Shell, network string) ChainClient {
	return &chainClient[ChainPresenters]{
		Shell: s,
		path:  "/v2/chains/" + network,
	}
}

// IndexChains returns all chains.
func (cli *chainClient[P]) IndexChains(c *cli.Context) (err error) {
	var p P
	return cli.getPage(cli.path, c.Int("page"), &p)
}

// ChainPresenter implements TableRenderer for a ChainResource
type ChainPresenter struct {
	presenters.ChainResource
}

// ToRow presents the ChainResource as a slice of strings.
func (p *ChainPresenter) ToRow() []string {
	return []string{p.GetID(), strconv.FormatBool(p.Enabled), p.Config}
}

// RenderTable implements TableRenderer
// Just renders a single row
func (p ChainPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{}
	rows = append(rows, p.ToRow())

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

// ChainPresenters implements TableRenderer for a slice of ChainPresenters.
type ChainPresenters []ChainPresenter

// RenderTable implements TableRenderer
func (ps ChainPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

func initChainSubCmds(s *Shell) []cli.Command {
	cmds := []cli.Command{}
	for _, network := range slices.Sorted(maps.Keys(relay.SupportedNetworks)) {
		if network == relay.NetworkDummy {
			continue
		}
		cmds = append(cmds, chainCommand(network, NewChainClient(s, network), cli.StringFlag{Name: "id", Usage: "chain ID"}))
	}
	return cmds
}
