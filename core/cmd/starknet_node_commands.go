package cmd

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// StarknetNodePresenter implements TableRenderer for a StarknetNodeResource.
type StarknetNodePresenter struct {
	presenters.StarknetNodeResource
}

// ToRow presents the StarknetNodeResource as a slice of strings.
func (p *StarknetNodePresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Name,
		p.ChainID,
		p.URL,
	}
	return row
}

var starknetNodeHeaders = []string{"ID", "Name", "Chain ID", "URL"}

// RenderTable implements TableRenderer
func (p StarknetNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(starknetNodeHeaders, rows, rt.Writer)

	return nil
}

// StarknetNodePresenters implements TableRenderer for a slice of StarknetNodePresenter.
type StarknetNodePresenters []StarknetNodePresenter

// RenderTable implements TableRenderer
func (ps StarknetNodePresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(starknetNodeHeaders, rows, rt.Writer)

	return nil
}

func NewStarknetNodeClient(c *Client) NodeClient {
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
	return newNodeClient[db.Node, StarknetNodePresenter, StarknetNodePresenters](c, "starknet", createNode)
}
