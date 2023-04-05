package cmd

import (
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// SolanaNodePresenter implements TableRenderer for a SolanaNodeResource.
type SolanaNodePresenter struct {
	presenters.SolanaNodeResource
}

// ToRow presents the SolanaNodeResource as a slice of strings.
func (p *SolanaNodePresenter) ToRow() []string {
	return []string{p.GetID(), p.Name, p.ChainID, p.State, p.Config}
}

// RenderTable implements TableRenderer
func (p SolanaNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(nodeHeaders, rows, rt.Writer)

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

	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

func NewSolanaNodeClient(c *Client) NodeClient {
	return newNodeClient[SolanaNodePresenters](c, "solana")
}
