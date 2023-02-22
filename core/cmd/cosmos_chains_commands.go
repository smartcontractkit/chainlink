package cmd

import (
	"encoding/json"
	"strconv"

	"github.com/smartcontractkit/chainlink-terra/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// CosmosChainPresenter implements TableRenderer for a CosmosChainResource
type CosmosChainPresenter struct {
	presenters.CosmosChainResource
}

// ToRow presents the CosmosChainResource as a slice of strings.
func (p *CosmosChainPresenter) ToRow() []string {
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

func CosmosChainClient(client *Client) ChainClient {
	return newChainClient[*db.ChainCfg, presenters.CosmosChainResource, CosmosChainPresenter, CosmosChainPresenters](client, "cosmos")
}
