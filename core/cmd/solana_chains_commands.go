package cmd

import (
	"encoding/json"
	"strconv"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// SolanaChainPresenter implements TableRenderer for a SolanaChainResource
type SolanaChainPresenter struct {
	presenters.SolanaChainResource
}

// ToRow presents the SolanaChainResource as a slice of strings.
func (p *SolanaChainPresenter) ToRow() []string {
	// NOTE: it's impossible to omitempty null fields when serializing to JSON: https://github.com/golang/go/issues/11939
	config, err := json.MarshalIndent(p.Config, "", "    ")
	if err != nil {
		panic(err)
	}

	row := []string{
		p.GetID(),
		strconv.FormatBool(p.Enabled),
		string(config),
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
	}
	return row
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

func SolanaChainClient(client *Client) ChainClient[*db.ChainCfg, presenters.SolanaChainResource, SolanaChainPresenter, SolanaChainPresenters] {
	return newChainClient[*db.ChainCfg, presenters.SolanaChainResource, SolanaChainPresenter, SolanaChainPresenters](client, "solana")
}
