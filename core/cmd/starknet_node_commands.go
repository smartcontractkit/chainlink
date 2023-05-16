package cmd

import (
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// StarkNetNodePresenter implements TableRenderer for a StarkNetNodeResource.
type StarkNetNodePresenter struct {
	presenters.StarkNetNodeResource
}

// ToRow presents the StarkNetNodeResource as a slice of strings.
func (p *StarkNetNodePresenter) ToRow() []string {
	return []string{p.GetID(), p.Name, p.ChainID, p.State, p.Config}
}

// RenderTable implements TableRenderer
func (p StarkNetNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(nodeHeaders, rows, rt.Writer)

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

	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

func NewStarkNetNodeClient(c *Client) NodeClient {
	return newNodeClient[StarkNetNodePresenters](c, "starknet")
}
