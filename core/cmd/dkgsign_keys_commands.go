package cmd

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type DKGSignKeyPresenter struct {
	JAID
	presenters.DKGSignKeyResource
}

var _ TableRenderer = DKGSignKeyPresenter{}
var _ TableRenderer = DKGSignKeyPresenters{}

// RenderTable implements TableRenderer
func (p DKGSignKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 DKGSign Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *DKGSignKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PublicKey,
	}

	return row
}

type DKGSignKeyPresenters []DKGSignKeyPresenter

// RenderTable implements TableRenderer
func (ps DKGSignKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 DKGSign Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewDKGSignKeysClient(c *Client) KeysClient {
	return newKeysClient[dkgsignkey.Key, DKGSignKeyPresenter, DKGSignKeyPresenters]("DKGSign", c)
}
