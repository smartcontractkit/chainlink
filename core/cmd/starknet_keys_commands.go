package cmd

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type StarknetKeyPresenter struct {
	JAID
	presenters.StarknetKeyResource
}

// RenderTable implements TableRenderer
func (p StarknetKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Starknet Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *StarknetKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
	}

	return row
}

type StarknetKeyPresenters []StarknetKeyPresenter

// RenderTable implements TableRenderer
func (ps StarknetKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Starknet Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewStarknetKeysClient(c *Client) KeysClient {
	return newKeysClient[starkkey.Key, StarknetKeyPresenter, StarknetKeyPresenters]("Starknet", c)
}
