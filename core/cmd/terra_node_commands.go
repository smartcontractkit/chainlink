package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	terratypes "github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
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
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
}

var terraNodeHeaders = []string{"ID", "Name", "Chain ID", "Tendermint URL", "Created", "Updated"}

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

func NewTerraNodeClient(c *Client) NodeClient[db.Node, presenters.TerraNodeResource, TerraNodePresenter, TerraNodePresenters] {
	createNode := func(c *cli.Context) (any, error) {
		name := c.String("name")
		chainID := c.String("chain-id")
		tendermintURL := c.String("tendermint-url")

		if name == "" {
			return nil, errors.New("missing --name")
		}
		if chainID == "" {
			return nil, errors.New("missing --chain-id")
		}

		if _, err2 := url.Parse(tendermintURL); err2 != nil {
			return nil, fmt.Errorf("invalid tendermint-url: %v", err2)
		}

		params := terratypes.NewNode{
			Name:          name,
			TerraChainID:  chainID,
			TendermintURL: tendermintURL,
		}
		return params, nil
	}
	return newNodeClient[db.Node, presenters.TerraNodeResource, TerraNodePresenter, TerraNodePresenters](c, "terra", createNode)
}
