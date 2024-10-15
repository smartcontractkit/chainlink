package cmd

import (
	"strconv"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// TronChainPresenter implements TableRenderer for a TronChainResource
type TronChainPresenter struct {
	presenters.TronChainResource
}

// ToRow presents the TronChainResource as a slice of strings.
func (p *TronChainPresenter) ToRow() []string {
	return []string{p.GetID(), strconv.FormatBool(p.Enabled), p.Config}
}

// RenderTable implements TableRenderer
// Just renders a single row
func (p TronChainPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{}
	rows = append(rows, p.ToRow())

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

// TronChainPresenters implements TableRenderer for a slice of TronChainPresenters.
type TronChainPresenters []TronChainPresenter

// RenderTable implements TableRenderer
func (ps TronChainPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

func TronChainClient(s *Shell) ChainClient {
	return newChainClient[TronChainPresenters](s, "tron")
}
