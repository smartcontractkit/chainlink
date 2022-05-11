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
	createNode := func(c *cli.Context) (node db.Node, err error) {
		node.Name = c.String("name")
		node.SolanaChainID = c.String("chain-id")
		node.SolanaURL = c.String("url")

		if node.Name == "" {
			err = errors.New("missing --name")
			return
		}
		if node.SolanaChainID == "" {
			err = errors.New("missing --chain-id")
			return
		}

		if _, err2 := url.Parse(node.SolanaURL); err2 != nil {
			err = errors.Errorf("invalid url: %v", err2)
			return
		}
		return
	}
	return newNodeClient[db.Node, presenters.SolanaNodeResource, SolanaNodePresenter, SolanaNodePresenters](c, "solana", createNode)
}
