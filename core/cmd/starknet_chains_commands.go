package cmd

import (
	"strconv"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// StarknetChainPresenter implements TableRenderer for a StarknetChainResource
type StarknetChainPresenter struct {
	presenters.StarknetChainResource
}

// ToRow presents the StarknetChainResource as a slice of strings.
func (p *StarknetChainPresenter) ToRow() []string {
	return []string{p.GetID(), strconv.FormatBool(p.Enabled), p.Config}
}

// RenderTable implements TableRenderer
// Just renders a single row
func (p StarknetChainPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{}
	rows = append(rows, p.ToRow())

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

// StarknetChainPresenters implements TableRenderer for a slice of StarknetChainPresenters.
type StarknetChainPresenters []StarknetChainPresenter

// RenderTable implements TableRenderer
func (ps StarknetChainPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

func StarknetChainClient(client *Client) ChainClient {
	return newChainClient[StarknetChainPresenters](client, "starknet")
}
