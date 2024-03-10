package cmd

import (
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// CosmosNodePresenter implements TableRenderer for a CosmosNodeResource.
type CosmosNodePresenter struct {
	presenters.CosmosNodeResource
}

// ToRow presents the CosmosNodeResource as a slice of strings.
func (p *CosmosNodePresenter) ToRow() []string {
	return []string{p.Name, p.ChainID, p.State, p.Config}
}

// RenderTable implements TableRenderer
func (p CosmosNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(nodeHeaders, rows, rt.Writer)

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

	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

func NewCosmosNodeClient(s *Shell) NodeClient {
	return newNodeClient[CosmosNodePresenters](s, "cosmos")
}
