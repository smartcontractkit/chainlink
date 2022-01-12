package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	terratypes "github.com/smartcontractkit/chainlink/core/chains/terra/types"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

// TerraNodePresenter implements TableRenderer for a TerraNodeResource.
type TerraNodePresenter struct {
	presenters.TerraNodeResource
}

// ToRow presents the TerraNodeResource as a slice of strings.
func (p *TerraNodePresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Name,
		p.TerraChainID,
		p.TendermintURL,
		p.FCDURL,
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
}

var terraNodeHeaders = []string{"ID", "Name", "Chain ID", "Tendermint URL", "FCD URL", "Created", "Updated"}

// RenderTable implements TableRenderer
func (p TerraNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(terraNodeHeaders, rows, rt.Writer)

	return nil
}

// TerraNodePresenters implements TableRenderer for a slice of TerraNodePresenter.
type TerraNodePresenters []TerraNodePresenter

// RenderTable implements TableRenderer
func (ps TerraNodePresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(terraNodeHeaders, rows, rt.Writer)

	return nil
}

// IndexTerraNodes returns all Terra nodes.
func (cli *Client) IndexTerraNodes(c *cli.Context) (err error) {
	return cli.getPage("/v2/nodes/terra", c.Int("page"), &TerraNodePresenters{})
}

// CreateTerraNode adds a new node to the nodelink node
func (cli *Client) CreateTerraNode(c *cli.Context) (err error) {
	name := c.String("name")
	chainID := c.String("chain-id")
	tendermintURL := c.String("tendermint-url")
	fcdURL := c.String("fcd-url")

	if name == "" {
		return cli.errorOut(errors.New("missing --name"))
	}
	if chainID == "" {
		return cli.errorOut(errors.New("missing --chain-id"))
	}

	if _, err2 := url.Parse(tendermintURL); err2 != nil {
		return cli.errorOut(fmt.Errorf("invalid tendermint-url: %v", err2))
	}
	if _, err2 := url.Parse(fcdURL); err2 != nil {
		return cli.errorOut(fmt.Errorf("invalid fcd-url: %v", err2))
	}

	params := terratypes.NewNode{
		Name:          name,
		TerraChainID:  chainID,
		TendermintURL: tendermintURL,
		FCDURL:        fcdURL,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post("/v2/nodes/terra", bytes.NewBuffer(body))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &TerraNodePresenter{})
}

// RemoveTerraNode removes a specific Terra Node by name.
func (cli *Client) RemoveTerraNode(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the id of the node to be removed"))
	}
	nodeID := c.Args().First()
	resp, err := cli.HTTP.Delete("/v2/nodes/terra/" + nodeID)
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}

	fmt.Printf("Node %v deleted\n", c.Args().First())
	return nil
}
