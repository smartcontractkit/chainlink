package cmd

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

type DKGEncryptKeyPresenter struct {
	JAID
	presenters.DKGEncryptKeyResource
}

var _ TableRenderer = DKGEncryptKeyPresenter{}
var _ TableRenderer = DKGEncryptKeyPresenters{}

// RenderTable implements TableRenderer
func (p DKGEncryptKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 DKGEncrypt Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *DKGEncryptKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PublicKey,
	}

	return row
}

type DKGEncryptKeyPresenters []DKGEncryptKeyPresenter

// RenderTable implements TableRenderer
func (ps DKGEncryptKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 DKGEncrypt Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewDKGEncryptKeysClient(c *Client) KeysClient {
	return newKeysClient[dkgsignkey.Key, DKGEncryptKeyPresenter, DKGEncryptKeyPresenters]("DKGEncrypt", c)
}
