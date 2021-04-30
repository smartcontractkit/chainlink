package cmd

import (
	"errors"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

type BridgePresenter struct {
	presenters.BridgeResource
}

// FriendlyConfirmations converts the confirmations to a string
func (p *BridgePresenter) FriendlyConfirmations() string {
	return strconv.FormatUint(uint64(p.Confirmations), 10)
}

// RenderTable implements TableRenderer
func (p *BridgePresenter) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"Name", "URL", "Default Confirmations", "Outgoing Token"})
	table.Append([]string{
		p.Name,
		p.URL,
		p.FriendlyConfirmations(),
		p.OutgoingToken,
	})
	render("Bridge", table)
	return nil
}

type BridgePresenters []BridgePresenter

// RenderTable implements TableRenderer
func (ps BridgePresenters) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"Name", "URL", "Confirmations"})
	for _, p := range ps {
		table.Append([]string{
			p.Name,
			p.URL,
			p.FriendlyConfirmations(),
		})
	}

	render("Bridges", table)
	return nil
}

// IndexBridges returns all bridges.
func (cli *Client) IndexBridges(c *cli.Context) (err error) {
	return cli.getPage("/v2/bridge_types", c.Int("page"), &BridgePresenters{})
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

	return cli.renderAPIResponse(resp, &BridgePresenter{})
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

	return cli.renderAPIResponse(resp, &BridgePresenter{})
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

	return cli.renderAPIResponse(resp, &BridgePresenter{})
}
