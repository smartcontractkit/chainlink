package cmd

import (
	"strconv"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// StarkNetChainPresenter implements TableRenderer for a StarkNetChainResource
type StarkNetChainPresenter struct {
	presenters.StarkNetChainResource
}

// ToRow presents the StarkNetChainResource as a slice of strings.
func (p *StarkNetChainPresenter) ToRow() []string {
	return []string{p.GetID(), strconv.FormatBool(p.Enabled), p.Config}
}

// RenderTable implements TableRenderer
// Just renders a single row
func (p StarkNetChainPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{}
	rows = append(rows, p.ToRow())

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

// StarkNetChainPresenters implements TableRenderer for a slice of StarkNetChainPresenters.
type StarkNetChainPresenters []StarkNetChainPresenter

// RenderTable implements TableRenderer
func (ps StarkNetChainPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

func StarkNetChainClient(s *Shell) ChainClient {
	return newChainClient[StarkNetChainPresenters](s, "starknet")
}
