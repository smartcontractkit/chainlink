package cmd

import (
	"encoding/json"
	"strconv"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// StarkNetChainPresenter implements TableRenderer for a StarkNetChainResource
type StarkNetChainPresenter struct {
	presenters.StarkNetChainResource
}

// ToRow presents the StarkNetChainResource as a slice of strings.
func (p *StarkNetChainPresenter) ToRow() []string {
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

func StarkNetChainClient(client *Client) ChainClient {
	return newChainClient[*db.ChainCfg, presenters.StarkNetChainResource, StarkNetChainPresenter, StarkNetChainPresenters](client, "starknet")
}
