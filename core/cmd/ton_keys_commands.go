package cmd

import (
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tonkey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type TONKeyPresenter struct {
	JAID
	presenters.TONKeyResource
}

// RenderTable implements TableRenderer
func (p TONKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "TON Public Key", "Base64 Address", "Raw Address"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 TON Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *TONKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
		p.AddressBase64,
		p.RawAddress,
	}

	return row
}

type TONKeyPresenters []TONKeyPresenter

// RenderTable implements TableRenderer
func (ps TONKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "TON Public Key", "Base64 Address", "Raw Address"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 TON Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewTONKeysClient(s *Shell) KeysClient {
	return newKeysClient[tonkey.Key, TONKeyPresenter, TONKeyPresenters]("TON", s)
}
