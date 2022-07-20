package cmd

import (
	"encoding/json"
	"strconv"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// EVMChainPresenter implements TableRenderer for an EVMChainResource.
type EVMChainPresenter struct {
	presenters.EVMChainResource
}

// ToRow presents the EVMChainResource as a slice of strings.
func (p *EVMChainPresenter) ToRow() []string {
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
func (p EVMChainPresenter) RenderTable(rt RendererTable) error {
	rows := [][]string{}
	rows = append(rows, p.ToRow())

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

// EVMChainPresenters implements TableRenderer for a slice of EVMChainPresenters.
type EVMChainPresenters []EVMChainPresenter

// RenderTable implements TableRenderer
func (ps EVMChainPresenters) RenderTable(rt RendererTable) error {
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(chainHeaders, rows, rt.Writer)

	return nil
}

func EVMChainClient(client *Client) ChainClient {
	return newChainClient[*evmtypes.ChainCfg, presenters.EVMChainResource, EVMChainPresenter, EVMChainPresenters](client, "evm")
}
