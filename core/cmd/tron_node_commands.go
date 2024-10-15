package cmd

import (
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// TronNodePresenter implements TableRenderer for a TronNodeResource.
type TronNodePresenter struct {
	presenters.TronNodeResource
}

// ToRow presents the TronNodeResource as a slice of strings.
func (p *TronNodePresenter) ToRow() []string {
	return []string{p.Name, p.ChainID, p.State, p.Config}
}

// RenderTable implements TableRenderer
func (p TronNodePresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

// TronNodePresenters implements TableRenderer for a slice of TronNodePresenter.
type TronNodePresenters []TronNodePresenter

// RenderTable implements TableRenderer
func (ps TronNodePresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(nodeHeaders, rows, rt.Writer)

	return nil
}

func NewTronNodeClient(s *Shell) NodeClient {
	return newNodeClient[TronNodePresenters](s, "tron")
}
