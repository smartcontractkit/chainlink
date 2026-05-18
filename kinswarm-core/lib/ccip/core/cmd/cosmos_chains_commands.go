package cmd

import (
	"strconv"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// CosmosChainPresenter implements TableRenderer for a CosmosChainResource
type CosmosChainPresenter struct {
	presenters.CosmosChainResource
}

// ToRow presents the CosmosChainResource as a slice of strings.
func (p *CosmosChainPresenter) ToRow() []string {
	return []string{p.GetID(), strconv.FormatBool(p.Enabled), p.Config}
}

// RenderTable implements TableRenderer
// Just renders a single row
func (p CosmosChainPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{}
	rows = append(rows, p.ToRow())

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

// CosmosChainPresenters implements TableRenderer for a slice of CosmosChainPresenters.
type CosmosChainPresenters []CosmosChainPresenter

// RenderTable implements TableRenderer
func (ps CosmosChainPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

func CosmosChainClient(s *Shell) ChainClient {
	return newChainClient[CosmosChainPresenters](s, "cosmos")
}
