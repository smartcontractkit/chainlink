package cmd

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/web"
)

func initBlocksSubCmds(s *Shell) []cli.Command {
	return []cli.Command{
		{
			Name:   "replay",
			Usage:  "Replays block data from the given number",
			Action: s.ReplayFromBlock,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:     "block-number",
					Usage:    "Block number to replay from",
					Required: true,
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "Whether to force broadcasting logs which were already consumed and that would otherwise be skipped",
				},
				cli.Int64Flag{
					Name:     "evm-chain-id",
					Usage:    "Chain ID of the EVM-based blockchain",
					Required: false,
				},
			},
		},
		{
			Name:   "find-lca",
			Usage:  "Find latest common block stored in DB and on chain",
			Action: s.FindLCA,
			Flags: []cli.Flag{
				cli.Int64Flag{
					Name:     "evm-chain-id",
					Usage:    "Chain ID of the EVM-based blockchain",
					Required: true,
				},
			},
		},
	}
}

// ReplayFromBlock replays chain data from the given block number until the most recent
func (s *Shell) ReplayFromBlock(c *cli.Context) (err error) {
	blockNumber := c.Int64("block-number")
	if blockNumber <= 0 {
		return s.errorOut(errors.New("Must pass a positive value in '--block-number' parameter"))
	}

	v := url.Values{}
	v.Add("force", strconv.FormatBool(c.Bool("force")))

	if c.IsSet("evm-chain-id") {
		v.Add("evmChainID", fmt.Sprintf("%d", c.Int64("evm-chain-id")))
	}

	buf := bytes.NewBufferString("{}")
	resp, err := s.HTTP.Post(s.ctx(),
		fmt.Sprintf(
			"/v2/replay_from_block/%v?%s",
			blockNumber,
			v.Encode(),
		), buf)
	if err != nil {
		return s.errorOut(err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	_, err = s.parseResponse(resp)
	if err != nil {
		return s.errorOut(err)
	}
	fmt.Println("Replay started")
	return nil
}

// LCAPresenter implements TableRenderer for an LCAResponse.
type LCAPresenter struct {
	web.LCAResponse
}

// ToRow presents the EVMChainResource as a slice of strings.
func (p *LCAPresenter) ToRow() []string {
	return []string{p.EVMChainID.String(), p.Hash, strconv.FormatInt(p.BlockNumber, 10)}
}

// RenderTable implements TableRenderer
// Just renders a single row
func (p LCAPresenter) RenderTable(rt RendererTable) error {
	renderList([]string{"ChainID", "Block Hash", "Block Number"}, [][]string{p.ToRow()}, rt.Writer)

	return nil
}

// FindLCA finds last common block stored in DB and on chain.
func (s *Shell) FindLCA(c *cli.Context) (err error) {
	v := url.Values{}

	if c.IsSet("evm-chain-id") {
		v.Add("evmChainID", fmt.Sprintf("%d", c.Int64("evm-chain-id")))
	}

	resp, err := s.HTTP.Get(s.ctx(),
		fmt.Sprintf(
			"/v2/find_lca?%s",
			v.Encode(),
		))
	if err != nil {
		return s.errorOut(err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &LCAPresenter{}, "Last Common Ancestor")
}
