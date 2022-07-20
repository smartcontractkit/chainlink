package cmd

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type TerraKeyPresenter struct {
	JAID
	presenters.TerraKeyResource
}

// RenderTable implements TableRenderer
func (p TerraKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Terra Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *TerraKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
	}

	return row
}

type TerraKeyPresenters []TerraKeyPresenter

// RenderTable implements TableRenderer
func (ps TerraKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Terra Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewTerraKeysClient(c *Client) KeysClient {
	return newKeysClient[terrakey.Key, TerraKeyPresenter, TerraKeyPresenters]("Terra", c)
}
