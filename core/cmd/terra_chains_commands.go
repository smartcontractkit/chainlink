package cmd

import (
	"encoding/json"
	"strconv"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// TerraChainPresenter implements TableRenderer for a TerraChainResource
type TerraChainPresenter struct {
	presenters.TerraChainResource
}

// ToRow presents the TerraChainResource as a slice of strings.
func (p *TerraChainPresenter) ToRow() []string {
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
func (p TerraChainPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{}
	rows = append(rows, p.ToRow())

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

// TerraChainPresenters implements TableRenderer for a slice of TerraChainPresenters.
type TerraChainPresenters []TerraChainPresenter

// RenderTable implements TableRenderer
func (ps TerraChainPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

func TerraChainClient(client *Client) ChainClient[*db.ChainCfg, presenters.TerraChainResource, TerraChainPresenter, TerraChainPresenters] {
	return newChainClient[*db.ChainCfg, presenters.TerraChainResource, TerraChainPresenter, TerraChainPresenters](client, "terra")
}
