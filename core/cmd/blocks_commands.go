package cmd

import (
	"bytes"
	"fmt"
	"net/http"
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

	forceBroadcast := c.Bool("force")

	buf := bytes.NewBufferString("{}")
	resp, err := cli.HTTP.Post(
		fmt.Sprintf(
			"/v2/replay_from_block/%v?force=%s",
			blockNumber,
			strconv.FormatBool(forceBroadcast),
		), buf)
	if err != nil {
		return cli.errorOut(err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bytes, err2 := cli.parseResponse(resp)
		if err2 != nil {
			return errors.Wrap(err2, "parseResponse error")
		}
		return cli.errorOut(errors.New(string(bytes)))
	}

	err = cli.printResponseBody(resp)
	return err
}
