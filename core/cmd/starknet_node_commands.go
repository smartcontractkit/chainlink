package cmd

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// StarkNetNodePresenter implements TableRenderer for a StarkNetNodeResource.
type StarkNetNodePresenter struct {
	presenters.StarkNetNodeResource
}

// ToRow presents the StarkNetNodeResource as a slice of strings.
func (p *StarkNetNodePresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Name,
		p.ChainID,
		p.URL,
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
}

var starknetNodeHeaders = []string{"ID", "Name", "Chain ID", "URL", "Created", "Updated"}

// RenderTable implements TableRenderer
func (p StarkNetNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(starknetNodeHeaders, rows, rt.Writer)

	return nil
}

// StarkNetNodePresenters implements TableRenderer for a slice of StarkNetNodePresenter.
type StarkNetNodePresenters []StarkNetNodePresenter

// RenderTable implements TableRenderer
func (ps StarkNetNodePresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(starknetNodeHeaders, rows, rt.Writer)

	return nil
}

func NewStarkNetNodeClient(c *Client) NodeClient {
	createNode := func(c *cli.Context) (node db.Node, err error) {
		node.Name = c.String("name")
		node.ChainID = c.String("chain-id")
		node.URL = c.String("url")

		if node.Name == "" {
			err = errors.New("missing --name")
			return
		}
		if node.ChainID == "" {
			err = errors.New("missing --chain-id")
			return
		}

		if _, err2 := url.Parse(node.URL); err2 != nil {
			err = errors.Errorf("invalid url: %v", err2)
			return
		}
		return
	}
	return newNodeClient[db.Node, StarkNetNodePresenter, StarkNetNodePresenters](c, "starknet", createNode)
}
