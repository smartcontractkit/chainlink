package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

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
	createNode := func(c *cli.Context) (node db.Node, err error) {
		node.Name = c.String("name")
		node.TerraChainID = c.String("chain-id")
		node.TendermintURL = c.String("tendermint-url")

		if node.Name == "" {
			err = errors.New("missing --name")
			return
		}
		if node.TerraChainID == "" {
			err = errors.New("missing --chain-id")
			return
		}

		if _, err2 := url.Parse(node.TendermintURL); err2 != nil {
			err = fmt.Errorf("invalid tendermint-url: %v", err2)
			return
		}
		return
	}
	return newNodeClient[db.Node, presenters.TerraNodeResource, TerraNodePresenter, TerraNodePresenters](c, "terra", createNode)
}
