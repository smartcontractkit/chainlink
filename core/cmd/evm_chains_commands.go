package cmd

import (
	"strconv"

	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// EVMChainPresenter implements TableRenderer for an EVMChainResource.
type EVMChainPresenter struct {
	presenters.EVMChainResource
}

// ToRow presents the EVMChainResource as a slice of strings.
func (p *EVMChainPresenter) ToRow() []string {
	return []string{p.GetID(), strconv.FormatBool(p.Enabled), p.Config}
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
	return newChainClient[EVMChainPresenters](client, "evm")
}
