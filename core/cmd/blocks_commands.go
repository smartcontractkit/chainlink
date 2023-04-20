package cmd

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
)

func initBlocksSubCmds(client *Client) []cli.Command {
	return []cli.Command{
		{
			Name:   "replay",
			Usage:  "Replays block data from the given number",
			Action: client.ReplayFromBlock,
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
	}
}

// ReplayFromBlock replays chain data from the given block number until the most recent
func (cli *Client) ReplayFromBlock(c *clipkg.Context) (err error) {
	blockNumber := c.Int64("block-number")
	if blockNumber <= 0 {
		return cli.errorOut(errors.New("Must pass a positive value in '--block-number' parameter"))
	}

	v := url.Values{}
	v.Add("force", strconv.FormatBool(c.Bool("force")))

	if c.IsSet("evm-chain-id") {
		v.Add("evmChainID", fmt.Sprintf("%d", c.Int64("evm-chain-id")))
	}

	buf := bytes.NewBufferString("{}")
	resp, err := cli.HTTP.Post(
		fmt.Sprintf(
			"/v2/replay_from_block/%v?%s",
			blockNumber,
			v.Encode(),
		), buf)
	if err != nil {
		return cli.errorOut(err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	_, err = cli.parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}
	fmt.Println("Replay started")
	return nil
}
