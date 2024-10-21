package cmd

import (
	"errors"
	"strconv"

	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initBrideSubCmds(s *Shell) []cli.Command {
	return []cli.Command{
		{
			Name:   "create",
			Usage:  "Create a new Bridge to an External Adapter",
			Action: s.CreateBridge,
		},
		{
			Name:   "destroy",
			Usage:  "Destroys the Bridge for an External Adapter",
			Action: s.RemoveBridge,
		},
		{
			Name:   "list",
			Usage:  "List all Bridges to External Adapters",
			Action: s.IndexBridges,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page",
					Usage: "page of results to display",
				},
			},
		},
		{
			Name:   "show",
			Usage:  "Show a Bridge's details",
			Action: s.ShowBridge,
		},
	}
}

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
func (s *Shell) IndexBridges(c *cli.Context) (err error) {
	return s.getPage("/v2/bridge_types", c.Int("page"), &BridgePresenters{})
}

// ShowBridge returns the info for the given Bridge name.
func (s *Shell) ShowBridge(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("must pass the name of the bridge to be shown"))
	}
	bridgeName := c.Args().First()
	resp, err := s.HTTP.Get(s.ctx(), "/v2/bridge_types/"+bridgeName)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &BridgePresenter{})
}

// CreateBridge adds a new bridge to the chainlink node
func (s *Shell) CreateBridge(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("must pass in the bridge's parameters [JSON blob | JSON filepath]"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return s.errorOut(err)
	}

	resp, err := s.HTTP.Post(s.ctx(), "/v2/bridge_types", buf)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &BridgePresenter{})
}

func (s *Shell) UpdateBridge(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("must pass the name of the bridge to be updated"))
	}
	bridgeName := c.Args().First()
	buf, err := getBufferFromJSON(c.Args().Get(1))
	if err != nil {
		return s.errorOut(err)
	}

	resp, err := s.HTTP.Patch(s.ctx(), "/v2/bridge_types/"+bridgeName, buf)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &BridgePresenter{})
}

// RemoveBridge removes a specific Bridge by name.
func (s *Shell) RemoveBridge(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("must pass the name of the bridge to be removed"))
	}
	bridgeName := c.Args().First()
	resp, err := s.HTTP.Delete(s.ctx(), "/v2/bridge_types/"+bridgeName)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &BridgePresenter{})
}
