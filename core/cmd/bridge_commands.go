package cmd

import (
	"errors"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

// IndexBridges returns all bridges.
func (cli *Client) IndexBridges(c *cli.Context) (err error) {
	return cli.getPage("/v2/bridge_types", c.Int("page"), &[]presenters.BridgeResource{})
}

// ShowBridge returns the info for the given Bridge name.
func (cli *Client) ShowBridge(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the name of the bridge to be shown"))
	}
	bridgeName := c.Args().First()
	resp, err := cli.HTTP.Get("/v2/bridge_types/" + bridgeName)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var resource presenters.BridgeResource
	return cli.renderAPIResponse(resp, &resource)
}

// CreateBridge adds a new bridge to the chainlink node
func (cli *Client) CreateBridge(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass in the bridge's parameters [JSON blob | JSON filepath]"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post("/v2/bridge_types", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var resource presenters.BridgeResource
	err = cli.renderAPIResponse(resp, &resource)
	return err
}

// RemoveBridge removes a specific Bridge by name.
func (cli *Client) RemoveBridge(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the name of the bridge to be removed"))
	}
	bridgeName := c.Args().First()
	resp, err := cli.HTTP.Delete("/v2/bridge_types/" + bridgeName)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var resource presenters.BridgeResource
	err = cli.renderAPIResponse(resp, &resource)
	return err
}
