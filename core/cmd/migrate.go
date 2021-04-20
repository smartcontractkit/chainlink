package cmd

import (
	"errors"
	"fmt"

	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
)

// Migrate jobs between the v1 (json) and v2 (toml) format.
func (cli *Client) Migrate(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the job id to be migrated"))
	}
	resp, err := cli.HTTP.Post(fmt.Sprintf("/v2/migrate/%s", c.Args().First()), nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	// TODO: return v2 ID
	return nil
}
