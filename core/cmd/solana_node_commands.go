package cmd

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// SolanaNodePresenter implements TableRenderer for a SolanaNodeResource.
type SolanaNodePresenter struct {
	presenters.SolanaNodeResource
}

// ToRow presents the SolanaNodeResource as a slice of strings.
func (p *SolanaNodePresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Name,
		p.SolanaChainID,
		p.SolanaURL,
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
}

var solanaNodeHeaders = []string{"ID", "Name", "Chain ID", "URL", "Created", "Updated"}

// RenderTable implements TableRenderer
func (p SolanaNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(solanaNodeHeaders, rows, rt.Writer)

	return nil
}

// SolanaNodePresenters implements TableRenderer for a slice of SolanaNodePresenter.
type SolanaNodePresenters []SolanaNodePresenter

// RenderTable implements TableRenderer
func (ps SolanaNodePresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(solanaNodeHeaders, rows, rt.Writer)

	return nil
}

func NewSolanaNodeClient(c *Client) NodeClient[db.Node, presenters.TerraNodeResource, TerraNodePresenter, TerraNodePresenters] {
	createNode := func(c *cli.Context) (any, error) {
		name := c.String("name")
		chainID := c.String("chain-id")
		urlStr := c.String("url")

		if name == "" {
			return nil, errors.New("missing --name")
		}
		if chainID == "" {
			return nil, errors.New("missing --chain-id")
		}

		if _, err2 := url.Parse(urlStr); err2 != nil {
			return nil, errors.Errorf("invalid url: %v", err2)
		}

		params := db.NewNode{
			Name:          name,
			SolanaChainID: chainID,
			SolanaURL:     urlStr,
		}
		return params, nil
	}
	return newNodeClient[db.Node, presenters.SolanaNodeResource, SolanaNodePresenter, SolanaNodePresenters](c, "solana", createNode)
}
