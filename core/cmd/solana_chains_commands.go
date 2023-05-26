package cmd

import (
	"strconv"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// SolanaChainPresenter implements TableRenderer for a SolanaChainResource
type SolanaChainPresenter struct {
	presenters.SolanaChainResource
}

// ToRow presents the SolanaChainResource as a slice of strings.
func (p *SolanaChainPresenter) ToRow() []string {
	return []string{p.GetID(), strconv.FormatBool(p.Enabled), p.Config}
}

// RenderTable implements TableRenderer
// Just renders a single row
func (p SolanaChainPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{}
	rows = append(rows, p.ToRow())

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

// SolanaChainPresenters implements TableRenderer for a slice of SolanaChainPresenters.
type SolanaChainPresenters []SolanaChainPresenter

// RenderTable implements TableRenderer
func (ps SolanaChainPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

func SolanaChainClient(client *Client) ChainClient {
	return newChainClient[SolanaChainPresenters](client, "solana")
}
