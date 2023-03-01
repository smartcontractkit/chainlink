package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// CosmosNodePresenter implements TableRenderer for a CosmosNodeResource.
type CosmosNodePresenter struct {
	presenters.CosmosNodeResource
}

// ToRow presents the CosmosNodeResource as a slice of strings.
func (p *CosmosNodePresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Name,
		p.CosmosChainID,
		p.TendermintURL,
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
}

var cosmosNodeHeaders = []string{"ID", "Name", "Chain ID", "Tendermint URL", "Created", "Updated"}

// RenderTable implements TableRenderer
func (p CosmosNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(cosmosNodeHeaders, rows, rt.Writer)

	return nil
}

// CosmosNodePresenters implements TableRenderer for a slice of CosmosNodePresenter.
type CosmosNodePresenters []CosmosNodePresenter

// RenderTable implements TableRenderer
func (ps CosmosNodePresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(cosmosNodeHeaders, rows, rt.Writer)

	return nil
}

func NewCosmosNodeClient(c *Client) NodeClient {
	createNode := func(c *cli.Context) (node db.Node, err error) {
		node.Name = c.String("name")
		node.CosmosChainID = c.String("chain-id")
		node.TendermintURL = c.String("tendermint-url")

		if node.Name == "" {
			err = errors.New("missing --name")
			return
		}
		if node.CosmosChainID == "" {
			err = errors.New("missing --chain-id")
			return
		}

		if _, err2 := url.Parse(node.TendermintURL); err2 != nil {
			err = fmt.Errorf("invalid tendermint-url: %v", err2)
			return
		}
		return
	}
	return newNodeClient[db.Node, CosmosNodePresenter, CosmosNodePresenters](c, "cosmos", createNode)
}
