package cmd

import (
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// EVMNodePresenter implements TableRenderer for an EVMNodeResource.
type EVMNodePresenter struct {
	presenters.EVMNodeResource
}

// ToRow presents the EVMNodeResource as a slice of strings.
func (p *EVMNodePresenter) ToRow() []string {
	return []string{p.GetID(), p.Name, p.ChainID, p.State, p.Config}
}

// RenderTable implements TableRenderer
func (p EVMNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

// EVMNodePresenters implements TableRenderer for a slice of EVMNodePresenter.
type EVMNodePresenters []EVMNodePresenter

// RenderTable implements TableRenderer
func (ps EVMNodePresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

func NewEVMNodeClient(c *Client) NodeClient {
	return newNodeClient[EVMNodePresenters](c, "evm")
}
